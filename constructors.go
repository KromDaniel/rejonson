package rejonson

import (
	"github.com/go-redis/redis/v7"
)

func ExtendClient(client *redis.Client) *Client {
	return &Client{
		client,
		&redisProcessor{
			Process: client.Process,
		},
	}
}

func ExtendPipeline(pipeline redis.Pipeliner) *Pipeline {
	return &Pipeline{
		pipeline,
		&redisProcessor{
			Process: pipeline.Process,
		},
	}
}
