package doc_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/KromDaniel/rejonson/v9/doc"
)

type User struct {
	Name   string            `json:"name"`
	Visits int               `json:"visits"`
	Tags   []string          `json:"tags"`
	Meta   map[string]string `json:"meta,omitempty"`
}

func redisAddr() string {
	if a := os.Getenv("REDIS_ADDR"); a != "" {
		return a
	}
	return "localhost:6379"
}

// forEachProtocol runs fn against a fresh store for RESP2 and RESP3, with a
// unique key prefix that is cleaned up afterwards.
func forEachProtocol(t *testing.T, fn func(t *testing.T, ctx context.Context, users *doc.Store[User])) {
	for _, proto := range []int{2, 3} {
		t.Run(fmt.Sprintf("resp%d", proto), func(t *testing.T) {
			ctx := context.Background()
			rdb := redis.NewClient(&redis.Options{Addr: redisAddr(), Protocol: proto})
			t.Cleanup(func() { _ = rdb.Close() })
			require.NoError(t, rdb.Ping(ctx).Err())

			prefix := fmt.Sprintf("doctest:%s:%d:%d:", t.Name(), proto, time.Now().UnixNano())
			t.Cleanup(func() {
				keys, _ := rdb.Keys(ctx, prefix+"*").Result()
				if len(keys) > 0 {
					rdb.Del(ctx, keys...)
				}
			})
			fn(t, ctx, doc.NewStore[User](rdb, doc.Prefix(prefix)))
		})
	}
}

func TestSetGet(t *testing.T) {
	forEachProtocol(t, func(t *testing.T, ctx context.Context, users *doc.Store[User]) {
		in := User{Name: "Ada", Visits: 3, Tags: []string{"a", "b"}, Meta: map[string]string{"k": "v"}}
		ok, err := users.Set(ctx, "1", in)
		require.NoError(t, err)
		assert.True(t, ok)

		out, err := users.Get(ctx, "1")
		require.NoError(t, err)
		assert.Equal(t, in, out)

		exists, err := users.Exists(ctx, "1")
		require.NoError(t, err)
		assert.True(t, exists)

		_, err = users.Get(ctx, "missing")
		assert.ErrorIs(t, err, doc.ErrNotFound)
		assert.ErrorIs(t, err, redis.Nil)

		deleted, err := users.Delete(ctx, "1")
		require.NoError(t, err)
		assert.True(t, deleted)
		_, err = users.Get(ctx, "1")
		assert.ErrorIs(t, err, doc.ErrNotFound)
	})
}

func TestSetConditions(t *testing.T) {
	forEachProtocol(t, func(t *testing.T, ctx context.Context, users *doc.Store[User]) {
		ok, err := users.Set(ctx, "1", User{Name: "first"}, doc.IfExists())
		require.NoError(t, err)
		assert.False(t, ok, "XX on a missing key must not write")

		ok, err = users.Set(ctx, "1", User{Name: "first"}, doc.IfNotExists())
		require.NoError(t, err)
		assert.True(t, ok)

		ok, err = users.Set(ctx, "1", User{Name: "second"}, doc.IfNotExists())
		require.NoError(t, err)
		assert.False(t, ok)

		u, err := users.Get(ctx, "1")
		require.NoError(t, err)
		assert.Equal(t, "first", u.Name)
	})
}

func TestTTL(t *testing.T) {
	forEachProtocol(t, func(t *testing.T, ctx context.Context, users *doc.Store[User]) {
		rdb := users.Client()

		_, err := users.Set(ctx, "1", User{Name: "a"}, doc.TTL(time.Hour))
		require.NoError(t, err)
		ttl := rdb.TTL(ctx, users.Key("1")).Val()
		assert.InDelta(t, time.Hour.Seconds(), ttl.Seconds(), 5)

		// NX + TTL on an existing key must not change its expiry.
		ok, err := users.Set(ctx, "1", User{Name: "b"}, doc.IfNotExists(), doc.TTL(time.Minute))
		require.NoError(t, err)
		assert.False(t, ok)
		assert.InDelta(t, time.Hour.Seconds(), rdb.TTL(ctx, users.Key("1")).Val().Seconds(), 5)

		// Plain Set and Update keep the expiry.
		_, err = users.Set(ctx, "1", User{Name: "c"})
		require.NoError(t, err)
		_, err = users.Update(ctx, "1", func(u *User) error { u.Visits++; return nil })
		require.NoError(t, err)
		assert.InDelta(t, time.Hour.Seconds(), rdb.TTL(ctx, users.Key("1")).Val().Seconds(), 5)
	})
}

func TestPaths(t *testing.T) {
	forEachProtocol(t, func(t *testing.T, ctx context.Context, users *doc.Store[User]) {
		_, err := users.Set(ctx, "1", User{Name: "Ada", Visits: 1, Tags: []string{"x"}})
		require.NoError(t, err)

		name, err := doc.GetPath[string](ctx, users, "1", "$.name")
		require.NoError(t, err)
		assert.Equal(t, "Ada", name)

		_, err = doc.GetPath[string](ctx, users, "1", "$.nope")
		assert.ErrorIs(t, err, doc.ErrNotFound)
		_, err = doc.GetPath[string](ctx, users, "missing", "$.name")
		assert.ErrorIs(t, err, doc.ErrNotFound)

		n, err := users.Incr(ctx, "1", "$.visits", 2)
		require.NoError(t, err)
		assert.Equal(t, 3.0, n)
		f, err := users.Incr(ctx, "1", "$.visits", 0.5)
		require.NoError(t, err)
		assert.Equal(t, 3.5, f)
		_, err = users.Incr(ctx, "1", "$.name", 1)
		assert.Error(t, err)

		// visits is fractional now; reset it so the struct decodes again later.
		_, err = users.SetPath(ctx, "1", "$.visits", 7)
		require.NoError(t, err)

		length, err := users.Append(ctx, "1", "$.tags", "y", "z")
		require.NoError(t, err)
		assert.EqualValues(t, 3, length)

		tags, err := doc.GetPathAll[string](ctx, users, "1", "$.tags[*]")
		require.NoError(t, err)
		assert.Equal(t, []string{"x", "y", "z"}, tags)

		for path, want := range map[string]int64{"$.tags": 3, "$.name": 3, "$": 3} {
			got, err := users.Len(ctx, "1", path)
			require.NoError(t, err, path)
			assert.Equal(t, want, got, path)
		}
		_, err = users.Len(ctx, "1", "$.visits")
		assert.Error(t, err)

		ok, err := users.SetPath(ctx, "1", "$.name", "Grace", doc.IfExists())
		require.NoError(t, err)
		assert.True(t, ok)
		ok, err = users.SetPath(ctx, "1", "$.name", "nope", doc.IfNotExists())
		require.NoError(t, err)
		assert.False(t, ok)

		require.NoError(t, users.Merge(ctx, "1", "$", map[string]any{
			"meta": map[string]string{"role": "admin"},
			"tags": nil, // removes the field
		}))
		u, err := users.Get(ctx, "1")
		require.NoError(t, err)
		assert.Equal(t, User{Name: "Grace", Visits: 7, Meta: map[string]string{"role": "admin"}}, u)

		removed, err := users.DeletePath(ctx, "1", "$.meta")
		require.NoError(t, err)
		assert.EqualValues(t, 1, removed)
	})
}

func TestUpdate(t *testing.T) {
	forEachProtocol(t, func(t *testing.T, ctx context.Context, users *doc.Store[User]) {
		_, err := users.Update(ctx, "missing", func(u *User) error { return nil })
		assert.ErrorIs(t, err, doc.ErrNotFound)

		_, err = users.Set(ctx, "1", User{Name: "Ada"})
		require.NoError(t, err)

		boom := errors.New("boom")
		_, err = users.Update(ctx, "1", func(u *User) error { u.Name = "changed"; return boom })
		assert.ErrorIs(t, err, boom)
		u, _ := users.Get(ctx, "1")
		assert.Equal(t, "Ada", u.Name, "fn error must not write")

		// Concurrent read-modify-write must not lose updates.
		const workers, perWorker = 8, 10
		var wg sync.WaitGroup
		errs := make(chan error, workers*perWorker)
		for range workers {
			wg.Go(func() {
				for range perWorker {
					if _, err := users.Update(ctx, "1", func(u *User) error { u.Visits++; return nil }); err != nil {
						errs <- err
					}
				}
			})
		}
		wg.Wait()
		close(errs)

		conflicts := 0
		for err := range errs {
			require.ErrorIs(t, err, doc.ErrConflict)
			conflicts++
		}
		u, err = users.Get(ctx, "1")
		require.NoError(t, err)
		assert.Equal(t, workers*perWorker-conflicts, u.Visits)
	})
}

func TestMGetPipelineScan(t *testing.T) {
	forEachProtocol(t, func(t *testing.T, ctx context.Context, users *doc.Store[User]) {
		var a, b *doc.Future[bool]
		err := users.Pipelined(ctx, func(p *doc.Pipe[User]) error {
			a = p.Set("a", User{Name: "A"})
			b = p.Set("b", User{Name: "B"}, doc.TTL(time.Minute))
			p.Redis().Set(ctx, users.Key("plain"), "not json", 0)
			return nil
		})
		require.NoError(t, err)
		okA, errA := a.Result()
		okB, errB := b.Result()
		require.NoError(t, errA)
		require.NoError(t, errB)
		assert.True(t, okA && okB)

		var ga, gMissing *doc.Future[User]
		require.NoError(t, users.TxPipelined(ctx, func(p *doc.Pipe[User]) error {
			ga = p.Get("a")
			gMissing = p.Get("missing")
			return nil
		}))
		va, err := ga.Result()
		require.NoError(t, err)
		assert.Equal(t, "A", va.Name)
		_, err = gMissing.Result()
		assert.ErrorIs(t, err, doc.ErrNotFound)

		got, err := users.MGet(ctx, "a", "missing", "b")
		require.NoError(t, err)
		assert.Equal(t, map[string]User{"a": {Name: "A"}, "b": {Name: "B"}}, got)

		for i := range 250 {
			_, err := users.Set(ctx, fmt.Sprintf("bulk:%03d", i), User{Visits: i})
			require.NoError(t, err)
		}

		var ids []string
		for e, err := range users.Scan(ctx, "*") {
			require.NoError(t, err)
			ids = append(ids, e.ID)
		}
		sort.Strings(ids)
		assert.Len(t, ids, 252, "the plain string key must be skipped")
		assert.Equal(t, []string{"a", "b"}, ids[:2])

		count := 0
		for e, err := range users.Scan(ctx, "bulk:*") {
			require.NoError(t, err)
			require.True(t, strings.HasPrefix(e.ID, "bulk:"))
			if count++; count == 5 {
				break // early break must be safe
			}
		}
		assert.Equal(t, 5, count)
	})
}

type upperNameCodec struct{}

func (upperNameCodec) Marshal(v any) ([]byte, error) {
	if u, ok := v.(User); ok {
		u.Name = strings.ToUpper(u.Name)
		v = u
	}
	return json.Marshal(v)
}
func (upperNameCodec) Unmarshal(data []byte, v any) error { return json.Unmarshal(data, v) }

func TestCodec(t *testing.T) {
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{Addr: redisAddr()})
	defer rdb.Close()
	users := doc.NewStore[User](rdb, doc.Prefix(fmt.Sprintf("doctest:codec:%d:", time.Now().UnixNano())), doc.WithCodec(upperNameCodec{}))
	defer users.Delete(ctx, "1")

	_, err := users.Set(ctx, "1", User{Name: "ada"})
	require.NoError(t, err)
	u, err := users.Get(ctx, "1")
	require.NoError(t, err)
	assert.Equal(t, "ADA", u.Name)
}

func TestRing(t *testing.T) {
	ctx := context.Background()
	ring := redis.NewRing(&redis.RingOptions{Addrs: map[string]string{"s1": redisAddr()}})
	defer ring.Close()
	users := doc.NewStore[User](ring, doc.Prefix(fmt.Sprintf("doctest:ring:%d:", time.Now().UnixNano())))
	defer func() {
		for _, id := range []string{"a", "b"} {
			users.Delete(ctx, id)
		}
	}()

	for _, id := range []string{"a", "b"} {
		_, err := users.Set(ctx, id, User{Name: id})
		require.NoError(t, err)
	}
	got, err := users.MGet(ctx, "a", "b", "missing")
	require.NoError(t, err)
	assert.Len(t, got, 2)

	var ids []string
	for e, err := range users.Scan(ctx, "*") {
		require.NoError(t, err)
		ids = append(ids, e.ID)
	}
	sort.Strings(ids)
	assert.Equal(t, []string{"a", "b"}, ids)
}
