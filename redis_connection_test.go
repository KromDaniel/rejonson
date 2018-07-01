package rejonson

import (
	"github.com/go-redis/redis"
)
const redis_tests_prefix = "rejonson::tests::"

func concatKey(key string)string {
	return redis_tests_prefix + key
}
func GetRedisConnection()*Client{
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       5,  // use default DB
	})
	return ExtendClient(client)
}