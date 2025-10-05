# Rejonson

Redis rejson extension built upon [go-redis](https://github.com/go-redis/redis)

[![Tests](https://github.com/KromDaniel/rejonson/actions/workflows/test.yml/badge.svg)](https://github.com/KromDaniel/rejonson/actions/workflows/test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/KromDaniel/rejonson)](https://goreportcard.com/report/github.com/KromDaniel/rejonson)
[![GoDoc](https://pkg.go.dev/badge/github.com/KromDaniel/rejonson)](https://pkg.go.dev/github.com/KromDaniel/rejonson)

## Table of Contents

1. [Quick start](#Quick%20Start)
2. [API](#api)
3. [Dependencies](#dependencies)
4. [Testing](#testing)
5. [License](#license)
6. [Contact](#contact)
 

## Quick start

### Install


#### go-redis v6
```shell
go get github.com/KromDaniel/rejonson
```

#### go-redis v7
```shell
go get github.com/KromDaniel/rejonson/v7
```

#### go-redis v8
```shell
go get github.com/KromDaniel/rejonson/v8
```
#### go-redis v9
```shell
go get github.com/KromDaniel/rejonson/v9
```

## Quick Start


```go
import (
	"github.com/KromDaniel/rejonson"
	"github.com/go-redis/redis"
)

func FunctionalApi(client *redis.Client) {
  // all rejonson.JsonX functions accepts redis.Client or redis.Pipeline
	// notice that some versions of go-redis also require context.Context (which is supported by rejonson)
	jsonStringCmd := rejonson.JsonGet(client, "key")
	if err := jsonStringCmd.Err(); err != nil {
		panic(err)
	}

	pipeline := client.Pipeline()
	rejonson.JsonGet(pipeline, "key")
	cmds, err := pipeline.Exec()
	// do something with cmds, err
}

func ExtendClient(client *redis.Client) *rejonson.Client {
  // You can extend go-redis client to rejonson.Client
  // that will have all JSON API as methods
	rejonsonClient := rejonson.ExtendClient(client)
	pingCmd := rejonsonClient.Ping()
	jsonCmd := rejonsonClient.JsonDel("key")

	return rejonsonClient
}
```

### Functional API
Rejonson exports Json`X` functions, all of them accept `RedisProcessor` as first parameter and `context.Context` (for go-redis versions >= 8) as second parameter, the other parameters are command specific


#### RedisProcessor
RedisProcessor is `interface` with the following definition:

```go
type RedisProcessor interface {
	Process(redis.Cmder) error
}
```
#### go-redis >= 8
```go
type RedisProcessor interface {
	Process(context.Context, redis.Cmder) error
}
```

By default all `*redis.Client`, `redis.Pipeliner`, `*redis.ClusterClient`, `*redis.SentinelClient` implenets that interface, so you can pass any of them to the rejonson functional API

#### example
```go
client := redis.NewClient(&redis.Options{ /*...*/ })

res := rejonson.JsonMGet(client, "key1", "key2", "$..a")
if res.Err() != nil {
	// handle error
}

for _, value := range res.Val() {
	// do something with value
}
```

### Extend Client 
Extends [go-redis](https://github.com/go-redis/redis) client with all ReJSON abilities, so you can use directly the rejson client for all redis usage and commands.

```go
// go redis client
goRedisClient := redis.NewClient(&redis.Options{
  Addr: "localhost:6379",
})

client := rejonson.ExtendClient(goRedisClient)
defer client.Close()

arr := []interface{}{"hello", "world", 1, map[string]interface{}{"key": 12}}
js, err := json.Marshal(arr)
if err != nil {
  // handle
}
// redis "native" command
client.Set("go-redis-cmd", "hello", time.Second)
// rejson command
client.JsonSet("rejson-cmd", ".", string(js))

// int command
arrLen, err := client.JsonArrLen("rejson-cmd", ".").Result()
if err != nil {
  // handle
}

fmt.Printf("Array length: %d", arrLen)
// Output: Array length: 4
```

### Pipeline
Client will also return extended `Pipeline` and `TXPipeline`

```go
goRedisClient := redis.NewClient(&redis.Options{
  Addr: "localhost:6379",
})

client := rejonson.ExtendClient(goRedisClient)

pipeline := client.Pipeline()
pipeline.JsonSet("rejson-cmd-pipeline", ".", "[10]")
pipeline.JsonNumMultBy("rejson-cmd-pipeline", "[0]", 10)
pipeline.Set("go-redis-pipeline-command", "hello from go-redis", time.Second)

_, err := pipeline.Exec()
if err != nil {
  // handle error
}
jsonString, err := client.JsonGet("rejson-cmd-pipeline").Result()
if err != nil {
  // handle error
}

fmt.Printf("Array %s", jsonString)

// Output: Array [100]
```

## API

Rejonson implements all the methods as described at [ReJson Commands](https://oss.redislabs.com/rejson/commands/) except for `JSON.DEBUG` and `JSON.RESP`.

The args will be serialized to redis directly so make sure to read [ReJSON command docs](https://oss.redislabs.com/redisjson/commands/)


All the rejson methods starts with the prefix of `Json` e.g `JsonDel`, `JsonArrIndex`, `JsonMGet`.<br/>Each command returns specific `go-redis.Cmder` by the specific request.

---------
Due to some ReJSON bug - [#issue-76](https://github.com/RedisLabsModules/rejson/issues/76), some empty strings will be ignored.

## Dependencies
Rejonson depends only on [go-redis](https://github.com/go-redis/redis). The [testing](#testing) also depends on assert library
## Test

<b>Rejonson tests must use real redis with ReJson module to run</b>

It is recommended to run the tests when using rejonson.</br>The unit tests will make sure your `go-redis` version is compatible and your `rejson` plugin supports all the methods and working as expected.

The testing library depends on [assert](https://github.com/stretchr/testify/assert) library

## License
Apache 2.0

## Contact
For any question or contribution, feel free to open an issue.

