package rejonson

import "github.com/go-redis/redis"


type ReJsonCommands interface {
	JSONDel(key, path string) *redis.IntCmd
	JSONGet(key string, args ...string) *redis.StringCmd
	JSONSet(key, path, json string, args ...string) *redis.StatusCmd
	//JSONMget(key string, args ...string) *redis.StringCmd
	//JSONType(key, path string) *redis.StringCmd
	//JSONNumIncrBy(key, path string, number int) *redis.StringCmd
	//JSONNumMultBy(key, path string, number int) *redis.StringCmd
}

type Client struct {
	*redis.Client
}

type Pipeline struct {
	redis.Pipeliner
}


func (cl *Client) Pipeline() *Pipeline {
	pip := cl.Client.Pipeline()
	return ExtendPipeline(pip)
}
func (cl *Client) JSONDel(key, path string) *redis.IntCmd {
	return jsonDelExecute(cl, key, path)
}

func (cl *Client) JSONGet(key string, args ...interface{}) *redis.StringCmd {
	return jsonGetExecute(cl, append([]interface{}{key}, args...)...)
}

func (cl *Client) JSONSet(key, path, json string, args ...interface{}) *redis.StatusCmd {
	return jsonSetExecute(cl, append([]interface{}{key, path, json}, args...)...)
}
