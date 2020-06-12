# Rejonson

Redis rejson extension built upon [go-redis](https://github.com/go-redis/redis)

[![Build Status](https://travis-ci.org/KromDaniel/rejonson.svg?branch=master)](https://travis-ci.org/KromDaniel/rejonson)
[![Coverage Status](https://coveralls.io/repos/github/KromDaniel/rejonson/badge.svg?branch=master)](https://coveralls.io/github/KromDaniel/rejonson?branch=master)

## Table of Contents

1. [Quick start](#install)
2. [API](#api)
3. [Dependencies](#dependencies)
4. [Testing](#testing)
5. [License](#license)
6. [Contact](#contact)
 

## Quick start

### Install

```shell
go get github.com/KromDaniel/rejonson
```

### Import

```go
import "github.com/KromDaniel/rejonson"
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
Rejonson depends only on [go-redis](https://github.com/go-redis/redis). The [testing](#testing) also depends on assert library.

## Test

<b>Rejonson tests must use real redis with ReJson to run</b>

It is recommended to run the unit tests when using rejonson.</br>The unit tests will make sure your `go-redis` version is compatible and your `rejson` plugin supports all the methods and working as expected.

The testing library depends on [assert](https://github.com/stretchr/testify/assert) library

```

## License
Apache 2.0

## Contact
For any question or contribution, feel free to open an issue.

