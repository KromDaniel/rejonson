package doc

import (
	"context"
	"iter"
	"strings"
	"sync"

	"github.com/redis/go-redis/v9"
)

const jsonKeyType = "ReJSON-RL"

// Entry is a document returned by Scan.
type Entry[T any] struct {
	ID    string
	Value T
}

// Scan iterates over all documents in the store whose id matches the glob
// pattern match ("*" for all). Documents are fetched in pipelined pages, so
// it is safe to use on large keyspaces. It follows SCAN semantics: documents
// changed during iteration may or may not be seen.
//
// On a ClusterClient or Ring, keys are collected from every master/shard
// first, then fetched page by page.
//
//	for e, err := range users.Scan(ctx, "*") {
//		if err != nil { return err }
//		fmt.Println(e.ID, e.Value.Name)
//	}
func (s *Store[T]) Scan(ctx context.Context, match string) iter.Seq2[Entry[T], error] {
	const page = 100
	pattern := s.prefix + match

	return func(yield func(Entry[T], error) bool) {
		emit := func(keys []string) bool {
			ids := make([]string, len(keys))
			for i, k := range keys {
				ids[i] = strings.TrimPrefix(k, s.prefix)
			}
			stop := false
			err := s.pipelinedGet(ctx, ids, func(id string, v T) {
				if !stop && !yield(Entry[T]{ID: id, Value: v}, nil) {
					stop = true
				}
			})
			if err != nil && !stop {
				yield(Entry[T]{}, err)
				return false
			}
			return !stop
		}

		if forEach := shardIterator(s.rdb); forEach != nil {
			keys, err := shardKeys(ctx, forEach, pattern, page)
			if err != nil {
				yield(Entry[T]{}, err)
				return
			}
			for start := 0; start < len(keys); start += page {
				if !emit(keys[start:min(start+page, len(keys))]) {
					return
				}
			}
			return
		}

		var cursor uint64
		for {
			keys, next, err := s.rdb.ScanType(ctx, cursor, pattern, page, jsonKeyType).Result()
			if err != nil {
				yield(Entry[T]{}, err)
				return
			}
			if len(keys) > 0 && !emit(keys) {
				return
			}
			if next == 0 {
				return
			}
			cursor = next
		}
	}
}

type forEachFunc func(ctx context.Context, fn func(ctx context.Context, node *redis.Client) error) error

// shardIterator returns how to visit every node holding data, or nil for a
// single-node client.
func shardIterator(rdb redis.UniversalClient) forEachFunc {
	switch c := rdb.(type) {
	case *redis.ClusterClient:
		return c.ForEachMaster
	case *redis.Ring:
		return c.ForEachShard
	}
	return nil
}

func shardKeys(ctx context.Context, forEach forEachFunc, pattern string, page int64) ([]string, error) {
	var (
		mu   sync.Mutex
		keys []string
	)
	err := forEach(ctx, func(ctx context.Context, node *redis.Client) error {
		it := node.ScanType(ctx, 0, pattern, page, jsonKeyType).Iterator()
		var local []string
		for it.Next(ctx) {
			local = append(local, it.Val())
		}
		if err := it.Err(); err != nil {
			return err
		}
		mu.Lock()
		keys = append(keys, local...)
		mu.Unlock()
		return nil
	})
	return keys, err
}
