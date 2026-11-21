package doc

import (
	"context"
	"errors"

	"github.com/redis/go-redis/v9"
)

// Future is the result of a command queued on a Pipe. Call Result after
// Pipelined has returned.
type Future[V any] struct {
	result func() (V, error)
}

// Result returns the command's value or error.
func (f *Future[V]) Result() (V, error) { return f.result() }

func failed[V any](err error) *Future[V] {
	return &Future[V]{result: func() (V, error) { var zero V; return zero, err }}
}

// Pipe queues document commands to be sent in one round trip.
type Pipe[T any] struct {
	ctx context.Context
	s   *Store[T]
	p   redis.Pipeliner
}

// Get queues a Get.
func (p *Pipe[T]) Get(id string) *Future[T] {
	cmd := do(p.ctx, p.p, "JSON.GET", p.s.Key(id), "$")
	return &Future[T]{result: func() (T, error) { return decodeFirst[T](p.s.codec, cmd) }}
}

// Set queues a Set.
func (p *Pipe[T]) Set(id string, v T, opts ...SetOption) *Future[bool] {
	data, err := p.s.codec.Marshal(v)
	if err != nil {
		return failed[bool](err)
	}
	cmd := jsonSet(p.ctx, p.p, p.s.Key(id), "$", data, buildSetConfig(opts))
	return &Future[bool]{result: func() (bool, error) { return setResult(cmd) }}
}

// Delete queues a Delete.
func (p *Pipe[T]) Delete(id string) *Future[bool] {
	cmd := p.p.Del(p.ctx, p.s.Key(id))
	return &Future[bool]{result: func() (bool, error) { n, err := cmd.Result(); return n > 0, err }}
}

// Redis exposes the underlying go-redis pipeline, to queue other commands
// in the same round trip.
func (p *Pipe[T]) Redis() redis.Pipeliner { return p.p }

// Pipelined queues the commands fn adds and sends them in one round trip.
// Each Future carries its own result; a missing document is reported by its
// Future as ErrNotFound and is not an error here. The returned error is the
// first other command error.
func (s *Store[T]) Pipelined(ctx context.Context, fn func(p *Pipe[T]) error) error {
	return s.pipelined(ctx, s.rdb.Pipeline(), fn)
}

// TxPipelined is like Pipelined but wraps the commands in MULTI/EXEC.
func (s *Store[T]) TxPipelined(ctx context.Context, fn func(p *Pipe[T]) error) error {
	return s.pipelined(ctx, s.rdb.TxPipeline(), fn)
}

func (s *Store[T]) pipelined(ctx context.Context, pipe redis.Pipeliner, fn func(p *Pipe[T]) error) error {
	if err := fn(&Pipe[T]{ctx: ctx, s: s, p: pipe}); err != nil {
		pipe.Discard()
		return err
	}
	cmds, err := pipe.Exec(ctx)
	for _, c := range cmds {
		if e := c.Err(); e != nil && !errors.Is(e, redis.Nil) {
			return e
		}
	}
	if err != nil && !errors.Is(err, redis.Nil) {
		return err
	}
	return nil
}
