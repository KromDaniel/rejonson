package main

import "github.com/KromDaniel/rejonson"
import (
	"github.com/go-redis/redis"
	"fmt"
)

func main(){
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	extendedClient := rejonson.ExtendClient(client)
	b, err := extendedClient.JSONGet("someKey").Bytes()
	fmt.Println(string(b), err)

	c, err := extendedClient.JSONSet("someKey", ".", "{\"key\":89, \"keyB\":70}").Result()

	fmt.Println(c, err)

	e, err := extendedClient.JSONDel("someKey", "").Result()
	fmt.Println(e, err)
}


