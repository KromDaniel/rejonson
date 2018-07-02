package main

import (
	"github.com/go-redis/redis"
	"github.com/KromDaniel/rejonson"
	"github.com/KromDaniel/jonson"
	"time"
	"fmt"
)

func main() {
	goRedisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	client := rejonson.ExtendClient(goRedisClient)
	defer client.Close()
	json := jonson.New([]interface{}{"hello", "world", "rejson", "and", "rejonson", "are", "awesome", 1, 2, 3, 4})
	client.Set("go-redis-command", "hello", time.Second)
	client.JsonSet("rejonson-json-command", ".", json.ToUnsafeJSONString())

	arrLen, err := client.JsonArrLen("rejonson-json-command", ".").Result() // int command
	if err != nil {
		// handle the error
	}

	fmt.Println("The array length is", arrLen) // The array length is 11

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
}
