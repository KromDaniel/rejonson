// Package doc stores typed Go values as RedisJSON documents.
//
// It sits on top of go-redis v9 and works with any redis.UniversalClient
// (Client, ClusterClient, Ring, failover clients), so go-redis hooks,
// OpenTelemetry instrumentation and connection settings all apply.
//
//	type User struct {
//		Name   string   `json:"name"`
//		Visits int      `json:"visits"`
//		Tags   []string `json:"tags"`
//	}
//
//	users := doc.NewStore[User](rdb, doc.Prefix("user:"))
//
//	_, err := users.Set(ctx, "42", User{Name: "Ada"}, doc.TTL(time.Hour))
//	u, err := users.Get(ctx, "42")                     // errors.Is(err, doc.ErrNotFound) if missing
//	name, err := doc.GetPath[string](ctx, users, "42", "$.name")
//	visits, err := users.Incr(ctx, "42", "$.visits", 1)
//	u, err = users.Update(ctx, "42", func(u *User) error { u.Visits++; return nil })
//
// Paths use JSONPath syntax ("$", "$.name", "$.items[0]"). Helpers that return
// a single value use the first match of the path.
//
// A server with JSON support is required: Redis 8+ (built in), or Redis 7.x
// with the RedisJSON module (for example Redis Stack).
package doc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	// ErrNotFound is returned when a document or path does not exist.
	// It wraps redis.Nil, so errors.Is(err, redis.Nil) also holds.
	ErrNotFound = fmt.Errorf("doc: not found: %w", redis.Nil)

	// ErrConflict is returned by Update when the document kept changing
	// concurrently and the retry budget was used up.
	ErrConflict = errors.New("doc: update conflict, retries exhausted")
)

// Codec encodes and decodes values. The default is encoding/json; plug in
// any compatible library (json/v2, sonic, go-json, ...) with WithCodec.
type Codec interface {
	Marshal(v any) ([]byte, error)
	Unmarshal(data []byte, v any) error
}

type stdJSON struct{}

func (stdJSON) Marshal(v any) ([]byte, error)      { return json.Marshal(v) }
func (stdJSON) Unmarshal(data []byte, v any) error { return json.Unmarshal(data, v) }

// DefaultCodec is the encoding/json based codec used when none is given.
var DefaultCodec Codec = stdJSON{}

// Store reads and writes documents of type T under keys "<prefix><id>".
// A Store is safe for concurrent use.
type Store[T any] struct {
	rdb        redis.UniversalClient
	prefix     string
	codec      Codec
	maxRetries int
}

// StoreOption configures a Store.
type StoreOption func(*storeConfig)

type storeConfig struct {
	prefix     string
	codec      Codec
	maxRetries int
}

// Prefix sets the key prefix, e.g. "user:" maps id "42" to key "user:42".
func Prefix(p string) StoreOption { return func(c *storeConfig) { c.prefix = p } }

// WithCodec replaces the default encoding/json codec.
func WithCodec(codec Codec) StoreOption { return func(c *storeConfig) { c.codec = codec } }

// MaxRetries sets how many times Update retries after a concurrent change
// (default 10).
func MaxRetries(n int) StoreOption { return func(c *storeConfig) { c.maxRetries = n } }

// NewStore returns a Store for documents of type T.
func NewStore[T any](rdb redis.UniversalClient, opts ...StoreOption) *Store[T] {
	cfg := storeConfig{codec: DefaultCodec, maxRetries: 10}
	for _, o := range opts {
		o(&cfg)
	}
	if cfg.maxRetries < 1 {
		cfg.maxRetries = 1
	}
	return &Store[T]{rdb: rdb, prefix: cfg.prefix, codec: cfg.codec, maxRetries: cfg.maxRetries}
}

// Key returns the Redis key used for id.
func (s *Store[T]) Key(id string) string { return s.prefix + id }

// Client returns the underlying go-redis client.
func (s *Store[T]) Client() redis.UniversalClient { return s.rdb }

// SetOption configures Set, SetPath and Pipe.Set.
type SetOption func(*setConfig)

type setConfig struct {
	mode string
	ttl  time.Duration
}

// IfNotExists writes only when the key (or path) does not exist yet (NX).
func IfNotExists() SetOption { return func(c *setConfig) { c.mode = "NX" } }

// IfExists writes only when the key (or path) already exists (XX).
func IfExists() SetOption { return func(c *setConfig) { c.mode = "XX" } }

// TTL sets the key's expiry, atomically with the write. It only applies when
// the write happens, so IfNotExists together with TTL never touches the
// expiry of an existing key. Without TTL, an existing key keeps its expiry.
func TTL(d time.Duration) SetOption { return func(c *setConfig) { c.ttl = d } }

func buildSetConfig(opts []SetOption) setConfig {
	var c setConfig
	for _, o := range opts {
		o(&c)
	}
	return c
}

// setWithTTL runs JSON.SET and PEXPIRE atomically, and only sets the expiry
// when JSON.SET actually wrote (it returns nil when NX/XX did not match).
var setWithTTL = redis.NewScript(`
local r
if ARGV[3] ~= '' then
  r = redis.call('JSON.SET', KEYS[1], ARGV[1], ARGV[2], ARGV[3])
else
  r = redis.call('JSON.SET', KEYS[1], ARGV[1], ARGV[2])
end
if r then
  redis.call('PEXPIRE', KEYS[1], ARGV[4])
end
return r
`)

// processor is satisfied by clients, *redis.Tx and pipelines.
type processor interface {
	Process(ctx context.Context, cmd redis.Cmder) error
}

func do(ctx context.Context, p processor, args ...any) *redis.Cmd {
	cmd := redis.NewCmd(ctx, args...)
	_ = p.Process(ctx, cmd)
	return cmd
}

type scriptProcessor interface {
	redis.Scripter
	processor
}

// jsonSet queues or runs JSON.SET with the given options on p.
func jsonSet(ctx context.Context, p scriptProcessor, key, path string, data []byte, cfg setConfig) *redis.Cmd {
	if cfg.ttl > 0 {
		return setWithTTL.Eval(ctx, p, []string{key}, path, string(data), cfg.mode, cfg.ttl.Milliseconds())
	}
	args := []any{"JSON.SET", key, path, string(data)}
	if cfg.mode != "" {
		args = append(args, cfg.mode)
	}
	return do(ctx, p, args...)
}

// setResult maps a JSON.SET reply to (written, error).
func setResult(cmd *redis.Cmd) (bool, error) {
	if err := cmd.Err(); err != nil {
		if errors.Is(err, redis.Nil) {
			return false, nil // NX/XX condition not met
		}
		return false, err
	}
	return true, nil
}

// Set stores v as the whole document for id. It reports false (and no error)
// when an IfNotExists/IfExists condition prevented the write.
func (s *Store[T]) Set(ctx context.Context, id string, v T, opts ...SetOption) (bool, error) {
	data, err := s.codec.Marshal(v)
	if err != nil {
		return false, err
	}
	return setResult(jsonSet(ctx, s.rdb, s.Key(id), "$", data, buildSetConfig(opts)))
}

// Get returns the document stored for id, or ErrNotFound.
func (s *Store[T]) Get(ctx context.Context, id string) (T, error) {
	return s.get(ctx, s.rdb, s.Key(id))
}

func (s *Store[T]) get(ctx context.Context, p processor, key string) (T, error) {
	return decodeFirst[T](s.codec, do(ctx, p, "JSON.GET", key, "$"))
}

// Delete removes the document for id and reports whether it existed.
func (s *Store[T]) Delete(ctx context.Context, id string) (bool, error) {
	n, err := s.rdb.Del(ctx, s.Key(id)).Result()
	return n > 0, err
}

// Exists reports whether a document exists for id.
func (s *Store[T]) Exists(ctx context.Context, id string) (bool, error) {
	n, err := s.rdb.Exists(ctx, s.Key(id)).Result()
	return n > 0, err
}

// Update reads the document for id, applies fn and writes it back atomically
// using optimistic locking (WATCH/MULTI). If the document changes in between,
// it retries, up to MaxRetries times, and then returns ErrConflict. If fn
// returns an error, nothing is written and that error is returned. The key's
// expiry is kept.
func (s *Store[T]) Update(ctx context.Context, id string, fn func(*T) error) (T, error) {
	key := s.Key(id)
	var out T
	for range s.maxRetries {
		err := s.rdb.Watch(ctx, func(tx *redis.Tx) error {
			v, err := s.get(ctx, tx, key)
			if err != nil {
				return err
			}
			if err := fn(&v); err != nil {
				return err
			}
			data, err := s.codec.Marshal(v)
			if err != nil {
				return err
			}
			_, err = tx.TxPipelined(ctx, func(p redis.Pipeliner) error {
				do(ctx, p, "JSON.SET", key, "$", string(data))
				return nil
			})
			if err == nil {
				out = v
			}
			return err
		}, key)
		if errors.Is(err, redis.TxFailedErr) {
			continue
		}
		return out, err
	}
	return out, ErrConflict
}

// MGet returns the documents for ids. Missing ids are left out of the map.
// On a single-node client it uses one JSON.MGET; on a ClusterClient or Ring,
// keys are fetched with a pipeline, so they may live on different shards.
func (s *Store[T]) MGet(ctx context.Context, ids ...string) (map[string]T, error) {
	out := make(map[string]T, len(ids))
	if len(ids) == 0 {
		return out, nil
	}
	if _, single := s.rdb.(*redis.Client); !single {
		return out, s.pipelinedGet(ctx, ids, func(id string, v T) { out[id] = v })
	}

	args := make([]any, 0, len(ids)+2)
	args = append(args, "JSON.MGET")
	for _, id := range ids {
		args = append(args, s.Key(id))
	}
	args = append(args, "$")
	vals, err := do(ctx, s.rdb, args...).Slice()
	if err != nil {
		return nil, err
	}
	for i, raw := range vals {
		if raw == nil || i >= len(ids) {
			continue
		}
		v, err := decodeFirstValue[T](s.codec, raw)
		if errors.Is(err, ErrNotFound) {
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("doc: decode %q: %w", s.Key(ids[i]), err)
		}
		out[ids[i]] = v
	}
	return out, nil
}

func (s *Store[T]) pipelinedGet(ctx context.Context, ids []string, yield func(string, T)) error {
	pipe := s.rdb.Pipeline()
	cmds := make([]*redis.Cmd, len(ids))
	for i, id := range ids {
		cmds[i] = do(ctx, pipe, "JSON.GET", s.Key(id), "$")
	}
	if _, err := pipe.Exec(ctx); err != nil && !errors.Is(err, redis.Nil) {
		// Per-command errors are inspected below; this catches e.g. network errors
		// that no command carries.
		if allNil(cmds) {
			return err
		}
	}
	for i, cmd := range cmds {
		v, err := decodeFirst[T](s.codec, cmd)
		if errors.Is(err, ErrNotFound) {
			continue
		}
		if err != nil {
			return fmt.Errorf("doc: get %q: %w", s.Key(ids[i]), err)
		}
		yield(ids[i], v)
	}
	return nil
}

func allNil(cmds []*redis.Cmd) bool {
	for _, c := range cmds {
		if c.Err() != nil && !errors.Is(c.Err(), redis.Nil) {
			return false
		}
	}
	return true
}
