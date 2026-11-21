# Rejonson

> Typed RedisJSON documents for Go, built on [go-redis](https://github.com/redis/go-redis).

[![CI](https://github.com/KromDaniel/rejonson/actions/workflows/ci.yml/badge.svg)](https://github.com/KromDaniel/rejonson/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/KromDaniel/rejonson/v9.svg)](https://pkg.go.dev/github.com/KromDaniel/rejonson/v9)
[![Coverage Status](https://coveralls.io/repos/github/KromDaniel/rejonson/badge.svg?branch=master)](https://coveralls.io/github/KromDaniel/rejonson?branch=master)

Rejonson lets you store Go structs as RedisJSON documents and work with them by JSONPath. You get generics, typed results, atomic updates, pipelines and iterators. It runs on top of the go-redis client you already use.

```bash
go get github.com/KromDaniel/rejonson/v9
```

Requires **go-redis v9** (v9.23+) and **Go 1.26+**. Using go-redis v6, v7 or v8? See [Older go-redis versions](#older-go-redis-versions).

## Contents

- [Quick start](#quick-start)
- [The `doc` package](#the-doc-package)
  - [Documents](#documents)
  - [Paths](#paths)
  - [Atomic updates](#atomic-updates)
  - [Batches and iteration](#batches-and-iteration)
  - [Options](#options)
- [Raw command helpers](#raw-command-helpers)
- [Redis server](#redis-server)
- [Older go-redis versions](#older-go-redis-versions)
- [Testing locally](#testing-locally)
- [Release cadence](#release-cadence)
- [Contributing](#contributing)
- [License](#license)

## Quick start

```go
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/KromDaniel/rejonson/v9/doc"
	"github.com/redis/go-redis/v9"
)

type User struct {
	Name   string   `json:"name"`
	Visits int      `json:"visits"`
	Tags   []string `json:"tags"`
}

func main() {
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})

	users := doc.NewStore[User](rdb, doc.Prefix("user:"))

	if _, err := users.Set(ctx, "42", User{Name: "Ada"}, doc.TTL(time.Hour)); err != nil {
		panic(err)
	}

	visits, _ := users.Incr(ctx, "42", "$.visits", 1)
	name, _ := doc.GetPath[string](ctx, users, "42", "$.name")
	u, _ := users.Get(ctx, "42")

	fmt.Println(name, visits, u.Tags)
}
```

## The `doc` package

`doc.NewStore[T]` accepts any `redis.UniversalClient`: `*redis.Client`, `*redis.ClusterClient`, `*redis.Ring`, or a failover (sentinel) client. Your go-redis hooks, OpenTelemetry instrumentation and connection settings keep working. A store maps id `"42"` to key `"<prefix>42"` and is safe for concurrent use.

### Documents

```go
ok, err := users.Set(ctx, "42", u)          // ok is false if IfNotExists/IfExists prevented the write
u, err := users.Get(ctx, "42")              // errors.Is(err, doc.ErrNotFound) when missing
existed, err := users.Delete(ctx, "42")
found, err := users.Exists(ctx, "42")
byID, err := users.MGet(ctx, "1", "2", "3") // map[string]User; missing ids are left out
```

### Paths

Paths use JSONPath (`$`, `$.name`, `$.items[0]`, `$.tags[*]`). Single-value helpers use the first match.

```go
name, err := doc.GetPath[string](ctx, users, "42", "$.name")
tags, err := doc.GetPathAll[string](ctx, users, "42", "$.tags[*]")

ok, err := users.SetPath(ctx, "42", "$.name", "Grace")
n, err := users.DeletePath(ctx, "42", "$.meta")
v, err := users.Incr(ctx, "42", "$.visits", 1)           // returns the new value
length, err := users.Append(ctx, "42", "$.tags", "ops")  // returns the new array length
size, err := users.Len(ctx, "42", "$.tags")              // array, string or object length
err = users.Merge(ctx, "42", "$", map[string]any{"meta": nil}) // RFC 7396 merge patch
```

### Atomic updates

`Update` is a read-modify-write with optimistic locking (WATCH/MULTI). It retries when the document changes concurrently and returns `doc.ErrConflict` once the retries run out. If your function returns an error, nothing is written.

```go
u, err := users.Update(ctx, "42", func(u *User) error {
	u.Visits++
	return nil
})
```

### Batches and iteration

```go
var a, b *doc.Future[User]
err := users.Pipelined(ctx, func(p *doc.Pipe[User]) error { // or TxPipelined for MULTI/EXEC
	a = p.Get("1")
	b = p.Get("2")
	p.Set("3", User{Name: "Lin"})
	p.Redis().Incr(ctx, "counter") // any go-redis command in the same round trip
	return nil
})
first, err := a.Result()

for e, err := range users.Scan(ctx, "*") { // Go 1.23 range-over-func; pipelined pages of 100
	if err != nil {
		return err
	}
	fmt.Println(e.ID, e.Value.Name)
}
```

On a `ClusterClient`, `MGet` and `Scan` fetch keys with a pipeline, so the keys may live in different hash slots.

### Options

| Option                         | Applies to                         | Effect                                                                                                                       |
| ------------------------------ | ---------------------------------- | ---------------------------------------------------------------------------------------------------------------------------- |
| `doc.Prefix("user:")`          | `NewStore`                         | Key prefix                                                                                                                   |
| `doc.WithCodec(c)`             | `NewStore`                         | Replaces `encoding/json` with your codec (e.g. json/v2, sonic, go-json)                                                      |
| `doc.MaxRetries(n)`            | `NewStore`                         | Number of `Update` retries (default 10)                                                                                      |
| `doc.IfNotExists()` / `IfExists()` | `Set`, `SetPath`, `Pipe.Set`   | `NX` / `XX`                                                                                                                  |
| `doc.TTL(d)`                   | `Set`, `SetPath`, `Pipe.Set`       | Sets the expiry atomically with the write, and only if the write happens. Without it, an existing key keeps its expiry.      |

## Raw command helpers

The root `rejonson/v9` package also has one function per RedisJSON command (`JsonGet`, `JsonSet`, `JsonArrAppend`, …), plus `ExtendClient`/`ExtendPipeline` for method-style calls. They take raw strings and return go-redis command types.

In v9 these helpers are **deprecated** (they still work and aren't going away in v9). go-redis v9 now has native, untyped JSON commands (`rdb.JSONGet`, `rdb.JSONSet`, …), and `doc` covers the typed case. Each helper's deprecation note names its replacement. They remain the main API in the v6–v8 modules, which makes them a stepping stone when migrating.

```go
import rejonson "github.com/KromDaniel/rejonson/v9"

rejonson.JsonSet(rdb, ctx, "example", "$", `{"a":1}`) // deprecated: rdb.JSONSet(ctx, "example", "$", `{"a":1}`)
client := rejonson.ExtendClient(rdb)                  // method style, also on Pipeline()/TXPipeline()
```

## Redis server

You need a server with JSON support:

- **Redis 8+**: JSON is built in. `docker run -p 6379:6379 redis:8`
- **Redis 7.x**: use Redis Stack (`redis/redis-stack-server`) or load the RedisJSON module.

CI tests against Redis 8 and Redis Stack 7.4, with both the RESP2 and RESP3 protocols.

## Older go-redis versions

The modules for go-redis v6–v8 are in maintenance mode: they get bug fixes and dependency updates. They provide the [raw command helpers](#raw-command-helpers) only, without the `doc` package.

| go-redis version | Module                              | Minimum Go | Note                                    |
| ---------------- | ----------------------------------- | ---------- | --------------------------------------- |
| **v9**           | `github.com/KromDaniel/rejonson/v9` | 1.26       | Recommended; includes `doc`             |
| v8               | `github.com/KromDaniel/rejonson/v8` | 1.20       | go-redis v8 is no longer developed      |
| v7               | `github.com/KromDaniel/rejonson/v7` | 1.20       | go-redis v7 is no longer developed      |
| v6               | `github.com/KromDaniel/rejonson`    | 1.20       | go-redis v6 is no longer developed; no `context` |

v8 and v7 take a `context.Context` (v8) or not (v7) exactly like their go-redis counterparts:

```go
import (
	rejonson "github.com/KromDaniel/rejonson/v8"
	"github.com/go-redis/redis/v8"
)

client := rejonson.ExtendClient(redis.NewClient(&redis.Options{Addr: "localhost:6379"}))
client.JsonSet(ctx, "json", ".", `["hello", 42]`)
length, err := client.JsonArrLen(ctx, "json", ".").Result()
```

For v6, drop the `ctx` argument and import `github.com/KromDaniel/rejonson` with `github.com/go-redis/redis`.

Upgrading go-redis to v9? Switch the import to `rejonson/v9`: the helpers have the same names and arguments, so the code compiles unchanged. Their deprecation notes then point you to `doc` or the native go-redis commands, which you can adopt at your own pace.

## Testing locally

```bash
docker run --rm -p 6379:6379 redis:8
./test.sh
```

The script runs `go vet` and `go test -race` on the root module and on each submodule (`v7`, `v8`, `v9`). With `GOTOOLCHAIN=local`, it skips any module that needs a newer Go than the one installed. Set `REDIS_ADDR` to use a server other than `localhost:6379` for the `doc` tests.

## Release cadence

CI runs on every PR against the latest four Go releases (currently 1.24.x–1.27.x; the v9 module runs only on Go 1.26+, which go-redis v9.23 requires). A monthly workflow updates the Go matrix automatically. Each module is tagged separately using semantic versioning.

## Contributing

We love community help! Check out the [CONTRIBUTING guide](./CONTRIBUTING.md) for workflow tips, coding standards, and release steps. Please also review the [Code of Conduct](./CODE_OF_CONDUCT.md) so everyone feels welcome.

## License

Apache 2.0 — have fun building!
