# Rejonson

Redis rejson extension built upon [go-redis](https://github.com/go-redis/redis)


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

*The examples are using [jonson](https://github.com/KromDaniel/jonson) library which is optional (but recommended)*

### Extend Client 
We extend a client to add it all the rejson abilities, extended client will have all the go-redis functionality with all the rejson functionality.<br/>Extended client is very comfortable to use because we don't need to pass a connection each time to some static method that accepts connection + args but instead just use the connection directly

```go
goRedisClient := redis.NewClient(&redis.Options{
    Addr:     "localhost:6379",
    Password: "",
    DB:       0,
})
client := rejonson.ExtendClient(goRedisClient)
defer client.Close()

// client got all go-redis commands with 
json := jonson.New([]interface{}{"hello", "world", "rejson", "and", "rejonson", "are", "awesome", 1,2,3,4})
client.Set("go-redis-command", "hello", time.Second)
client.JsonSet("rejonson-json-command", ".", json.ToUnsafeJSONString())

arrLen, err := client.JsonArrLen("rejonson-json-command", ".").Result() // int command
if err != nil {
    // handle the error
}

fmt.Println("The array length is",  arrLen) // The array length is 11
```

### Pipeline
Client will also return extended `Pipeline` and `TXPipeline`

```go
pipeline := client.Pipeline()
pipeline.JsonNumMultBy("rejonson-json-command", "[7]", 10)
pipeline.Set("go-redis-pipeline-command", "hello from go-redis", time.Second)

_, err = pipeline.Exec()
if err != nil {
    // handle error
}
jsonString, err := client.JsonGet("rejonson-json-command").Result()
if err != nil {
    fmt.Println(err.Error())
    // handle error
}
json = jonson.ParseUnsafe([]byte(jsonString))

fmt.Println(json.At(7).GetUnsafeFloat64()) // 10
```

## API

Rejonson supports all the methods as described at [ReJson Commands](https://oss.redislabs.com/rejson/commands/) except for `JSON.DEBUG` and `JSON.RESP`
The arguments you are sending will be transferred to redis, so please make sure you are working by the rejson documentation.

Due to some rejson behavior [#issue-76](https://github.com/RedisLabsModules/rejson/issues/76), empty strings will be ignored

All the rejson methods starts with the prefix of `Json` e.g `JsonDel`, `JsonArrIndex`, `JsonMGet`.<br/>Each command returns specific `go-redis.Cmder` by the specific request.


## Dependencies
Rejonson got only single dependency which is [go-redis]("https://github.com/go-redis/redis"). The [testing](#testing) got some other dependencies as well

## Testing
It is recommended to run the unit tests when using rejonson.</br>The unit tests will make sure your `go-redis` version is compatible and your `rejson` plugin supports all the methods and working as expected.

The testing got few dependencies of its' own:

* [jonson]("https://github.com/KromDaniel/jonson")
* [assert]("https://github.com/stretchr/testify/assert")

### Configuring Redis
In order to guarantee the code is safe for use **The unit tests will have to use a real redis**.
To configure the connection edit the file `test_config.json`

`redisConnection` will just be passed to [go-redis.Options](https://godoc.org/github.com/go-redis/redis#Options), the keys should confirm with the `Options` keys.</br> 
`redisKeyPrefix` is the prefix the unit tests will add to written test data at the redis (the data is deleted at the end of each test)

```json
{
  "redisConnection": {
    "Addr":  "localhost:6379",
    "Password": "",
    "DB": 0
  },
  "redisKeyPrefix": "rejonson::tests::"
}
```

## License
Apache 2.0

## Contact
For any question or contribution, feel free to contact me at
kromdan@gmail.com
