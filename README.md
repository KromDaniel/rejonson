goRedisClient := redis.NewClient(&redis.Options{
goRedisClient := redis.NewClient(&redis.Options{
pipeline.JsonSet("rejson-cmd-pipeline", ".", "[10]")
pipeline.JsonNumMultBy("rejson-cmd-pipeline", "[0]", 10)
pipeline.Set("go-redis-pipeline-command", "hello from go-redis", time.Second)

# Rejonson

> A friendly ReJSON client for Go, powered by the excellent [go-redis](https://github.com/go-redis/redis) ecosystem.

[![CI](https://github.com/KromDaniel/rejonson/actions/workflows/ci.yml/badge.svg)](https://github.com/KromDaniel/rejonson/actions/workflows/ci.yml)
[![Coverage Status](https://coveralls.io/repos/github/KromDaniel/rejonson/badge.svg?branch=master)](https://coveralls.io/github/KromDaniel/rejonson?branch=master)

Rejonson wraps the entire RedisJSON command surface so you can keep using the `go-redis` APIs you already know—clients, pipelines, clusters, and all. Whether you prefer the functional helpers (e.g. `JsonGet`) or method-style access via an extended client, everything feels familiar.

## Highlights

- ✅ Supports go-redis v6, v7, v8, and v9 with generated, type-safe wrappers
- 🧰 Works with clients, pipelines, clusters, and sentinels as drop-in processors
- ⚙️ Keeps your project fresh: CI automatically tests the latest four Go releases and a scheduled job bumps the matrix for you
- 🧪 Ships with integration tests that run against Redis Stack / RedisJSON out of the box
- 📚 Built for open-source collaboration with contributor docs and issue templates in repo

## Contents

- [Compatibility & Install](#compatibility--install)
- [Getting Started](#getting-started)
  - [Functional helpers](#functional-helpers)
  - [Extended client](#extended-client)
- [Testing locally](#testing-locally)
- [Release cadence](#release-cadence)
- [Contributing](#contributing)
- [License](#license)

## Compatibility & Install

| go-redis version | Module path                         | Install command                            |
| ---------------- | ----------------------------------- | ------------------------------------------ |
| v6               | `github.com/KromDaniel/rejonson`    | `go get github.com/KromDaniel/rejonson`    |
| v7               | `github.com/KromDaniel/rejonson/v7` | `go get github.com/KromDaniel/rejonson/v7` |
| v8               | `github.com/KromDaniel/rejonson/v8` | `go get github.com/KromDaniel/rejonson/v8` |
| v9               | `github.com/KromDaniel/rejonson/v9` | `go get github.com/KromDaniel/rejonson/v9` |

Pick the module that matches the `go-redis` major version in your project and `go get` it like any other dependency. Each module is versioned independently so you can upgrade only what you need.

## Getting Started

### Functional helpers

The functional API works with any `go-redis` processor—clients, pipelines, clusters, or sentinels. For go-redis v8 and newer, pass a `context.Context` as the second argument.

```go
package main

import (
	"context"
	"fmt"

	rejonson "github.com/KromDaniel/rejonson/v9"
	"github.com/redis/go-redis/v9"
)

func main() {
	ctx := context.Background()
	client := redis.NewClient(&redis.Options{Addr: "localhost:6379"})

	if err := rejonson.JsonSet(client, ctx, "example", ".", `{"a":1}`).Err(); err != nil {
		panic(err)
	}

	val, err := rejonson.JsonGet(client, ctx, "example").Result()
	if err != nil {
		panic(err)
	}

	fmt.Println(val)
}
```

All helpers follow the same naming convention as their RedisJSON counterparts (`JsonDel`, `JsonArrIndex`, `JsonMGet`, …) and return the corresponding `go-redis` command type (`*redis.StringCmd`, `*redis.IntCmd`, etc.).

### Extended client

Prefer method calls? Wrap an existing client with `ExtendClient` (or `ExtendClusterClient`, `ExtendSentinelClient`) and use the JSON API side-by-side with native Redis commands.

```go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	rejonson "github.com/KromDaniel/rejonson/v8"
	"github.com/go-redis/redis/v8"
)

func main() {
	ctx := context.Background()
	goRedis := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	client := rejonson.ExtendClient(goRedis)
	defer client.Close()

	payload, _ := json.Marshal([]interface{}{"hello", 42, map[string]string{"key": "value"}})

	client.Set(ctx, "native", "hello", time.Minute)
	client.JsonSet(ctx, "json", ".", string(payload))

	length, _ := client.JsonArrLen(ctx, "json", ".").Result()
	fmt.Printf("Array length: %d\n", length)
}
```

Pipelines (`Pipeline`, `TxPipeline`) returned from an extended client also understand the JSON commands, so batching commands stays ergonomic.

## Testing locally

RedisJSON is required to run the integration tests. The quickest setup is Redis Stack in Docker:

```bash
docker run --rm -p 6379:6379 redis/redis-stack-server:latest
./test.sh
```

> The test script walks through the root module and each submodule (`v7`, `v8`, `v9`) to guard against regressions across all supported go-redis versions.

## Release cadence

Continuous integration runs on every PR against the latest four Go toolchains (currently 1.22.x–1.25.x). A scheduled workflow re-generates the matrix once a month so new Go releases are picked up automatically—no manual edits required. We tag releases per module using semantic versioning so `go get` upgrades remain smooth.

## Contributing

We love community help! Check out the [CONTRIBUTING guide](./CONTRIBUTING.md) for workflow tips, coding standards, and release steps. Please also review the [Code of Conduct](./CODE_OF_CONDUCT.md) so everyone feels welcome.

## License

Apache 2.0 — have fun building!
