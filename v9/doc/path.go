package doc

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// rawReply turns a reply value into JSON bytes. JSON commands reply with a
// JSON string on RESP2; on RESP3 some reply with native arrays/numbers, which
// are re-encoded so both protocols decode the same way.
func rawReply(v any) ([]byte, error) {
	switch x := v.(type) {
	case nil:
		return nil, ErrNotFound
	case string:
		return []byte(x), nil
	case []byte:
		return x, nil
	default:
		return json.Marshal(x)
	}
}

// matches splits a JSONPath ("$...") reply, which is always a JSON array of
// matches, into the raw JSON of each match.
func matches(raw []byte) ([]json.RawMessage, error) {
	raw = bytes.TrimSpace(raw)
	if len(raw) == 0 || raw[0] != '[' {
		// A single scalar (RESP3 replies for a single match).
		return []json.RawMessage{raw}, nil
	}
	var out []json.RawMessage
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func cmdValue(cmd interface{ Result() (any, error) }) (any, error) {
	v, err := cmd.Result()
	if errors.Is(err, redis.Nil) {
		return nil, ErrNotFound
	}
	return v, err
}

func decodeFirst[V any](codec Codec, cmd interface{ Result() (any, error) }) (V, error) {
	v, err := cmdValue(cmd)
	if err != nil {
		var zero V
		return zero, err
	}
	return decodeFirstValue[V](codec, v)
}

func decodeFirstValue[V any](codec Codec, reply any) (V, error) {
	var out V
	all, err := decodeAllValue[V](codec, reply)
	if err != nil {
		return out, err
	}
	if len(all) == 0 {
		return out, ErrNotFound
	}
	return all[0], nil
}

func decodeAllValue[V any](codec Codec, reply any) ([]V, error) {
	raw, err := rawReply(reply)
	if err != nil {
		return nil, err
	}
	ms, err := matches(raw)
	if err != nil {
		return nil, err
	}
	out := make([]V, 0, len(ms))
	for _, m := range ms {
		var v V
		if err := codec.Unmarshal(m, &v); err != nil {
			return nil, err
		}
		out = append(out, v)
	}
	return out, nil
}

// GetPath returns the first value matching a JSONPath in the document for
// id, decoded as V. It returns ErrNotFound if the document or path is missing.
//
//	name, err := doc.GetPath[string](ctx, users, "42", "$.name")
func GetPath[V, T any](ctx context.Context, s *Store[T], id, path string) (V, error) {
	return decodeFirst[V](s.codec, do(ctx, s.rdb, "JSON.GET", s.Key(id), path))
}

// GetPathAll returns every value matching a JSONPath, decoded as V.
//
//	tags, err := doc.GetPathAll[string](ctx, users, "42", "$.tags[*]")
func GetPathAll[V, T any](ctx context.Context, s *Store[T], id, path string) ([]V, error) {
	v, err := cmdValue(do(ctx, s.rdb, "JSON.GET", s.Key(id), path))
	if err != nil {
		return nil, err
	}
	return decodeAllValue[V](s.codec, v)
}

// SetPath sets the value at path inside the document for id. The document
// must exist unless path is "$". It reports false when an IfNotExists or
// IfExists condition prevented the write.
func (s *Store[T]) SetPath(ctx context.Context, id, path string, value any, opts ...SetOption) (bool, error) {
	data, err := s.codec.Marshal(value)
	if err != nil {
		return false, err
	}
	return setResult(jsonSet(ctx, s.rdb, s.Key(id), path, data, buildSetConfig(opts)))
}

// DeletePath removes every value matching path and returns how many were
// removed.
func (s *Store[T]) DeletePath(ctx context.Context, id, path string) (int64, error) {
	return do(ctx, s.rdb, "JSON.DEL", s.Key(id), path).Int64()
}

// Merge applies an RFC 7396 JSON Merge Patch at path: fields in patch are
// added or replaced, and fields set to null are removed.
func (s *Store[T]) Merge(ctx context.Context, id, path string, patch any) error {
	data, err := s.codec.Marshal(patch)
	if err != nil {
		return err
	}
	return do(ctx, s.rdb, "JSON.MERGE", s.Key(id), path, string(data)).Err()
}

// Incr adds by to the number at path and returns the new value of the
// first match.
func (s *Store[T]) Incr(ctx context.Context, id, path string, by float64) (float64, error) {
	cmd := do(ctx, s.rdb, "JSON.NUMINCRBY", s.Key(id), path, by)
	v, err := decodeFirst[*float64](s.codec, cmd)
	if err != nil {
		return 0, err
	}
	if v == nil {
		return 0, fmt.Errorf("doc: incr %q %s: value is not a number", s.Key(id), path)
	}
	return *v, nil
}

// Append adds values to the end of the array at path and returns the new
// length of the first matching array.
func (s *Store[T]) Append(ctx context.Context, id, path string, values ...any) (int64, error) {
	args := make([]any, 0, len(values)+3)
	args = append(args, "JSON.ARRAPPEND", s.Key(id), path)
	for _, v := range values {
		data, err := s.codec.Marshal(v)
		if err != nil {
			return 0, err
		}
		args = append(args, string(data))
	}
	return firstInt(s.Key(id), path, do(ctx, s.rdb, args...))
}

// Len returns the length of the array, string or object at path (first
// match), based on the value's JSON type.
func (s *Store[T]) Len(ctx context.Context, id, path string) (int64, error) {
	key := s.Key(id)
	typ, err := decodeFirst[string](s.codec, jsonType(ctx, s, key, path))
	if err != nil {
		return 0, err
	}
	var cmd string
	switch typ {
	case "array":
		cmd = "JSON.ARRLEN"
	case "string":
		cmd = "JSON.STRLEN"
	case "object":
		cmd = "JSON.OBJLEN"
	default:
		return 0, fmt.Errorf("doc: len %q %s: %s has no length", key, path, typ)
	}
	return firstInt(key, path, do(ctx, s.rdb, cmd, key, path))
}

func jsonType[T any](ctx context.Context, s *Store[T], key, path string) *typeReply {
	return &typeReply{do(ctx, s.rdb, "JSON.TYPE", key, path)}
}

// typeReply re-encodes the JSON.TYPE array of type names as JSON so it can go
// through the regular match decoding.
type typeReply struct {
	cmd interface{ Result() (any, error) }
}

func (t *typeReply) Result() (any, error) {
	v, err := t.cmd.Result()
	if err != nil {
		return nil, err
	}
	b, err := json.Marshal(flatten(v))
	return string(b), err
}

// flatten unwraps RESP3's nested [[type]] form of JSON.TYPE.
func flatten(v any) any {
	arr, ok := v.([]any)
	if !ok {
		return []any{v}
	}
	out := make([]any, 0, len(arr))
	for _, e := range arr {
		if inner, ok := e.([]any); ok {
			out = append(out, inner...)
		} else {
			out = append(out, e)
		}
	}
	return out
}

// firstInt reads the first entry of an integer-per-match reply.
func firstInt(key, path string, cmd interface{ Result() (any, error) }) (int64, error) {
	v, err := cmdValue(cmd)
	if err != nil {
		return 0, err
	}
	switch x := v.(type) {
	case int64:
		return x, nil
	case []any:
		if len(x) == 0 {
			return 0, ErrNotFound
		}
		if n, ok := x[0].(int64); ok {
			return n, nil
		}
		return 0, fmt.Errorf("doc: %q %s: path does not match a value of the right type", key, path)
	default:
		return 0, fmt.Errorf("doc: %q %s: unexpected reply %T", key, path, v)
	}
}
