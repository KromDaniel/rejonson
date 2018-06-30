package rejonson

import "github.com/go-redis/redis"


type Client struct {
	*redis.Client
}

type Pipeline struct {
	redis.Pipeliner
}
func extendJsonDel(r redisProcessor) func(key, path string)*redis.IntCmd{
	return func(key, path string) *redis.IntCmd{
		return jsonDelExecute(r, key, path)
	}
}

func (cl *Client) Pipeline() *Pipeline {
	pip := cl.Client.Pipeline()
	return ExtendPipeline(pip)
}
func (cl *Client) JsonDel(key, path string) *redis.IntCmd {
	return jsonDelExecute(cl, key, path)
}

func (cl *Client) JsonGet(key string, args ...interface{}) *redis.StringCmd {
	return jsonGetExecute(cl, append([]interface{}{key}, args...)...)
}

func (cl *Client) JsonSet(key, path, json string, args ...interface{}) *redis.StatusCmd {
	return jsonSetExecute(cl, append([]interface{}{key, path, json}, args...)...)
}

func (cl *Client) JsonMGet(key string, args ...interface{}) *redis.StringSliceCmd {
	return jsonMGetExecute(cl, append([]interface{}{key}, args...)...)
}

func (cl *Client) JsonType(key, path string) *redis.StringCmd {
	return jsonTypeExecute(cl, key, path)
}

func (cl *Client) JsonNumIncrBy(key, path string, num int) *redis.StringCmd {
	return jsonNumIncrByExecute(cl, key, path, num)
}

func (cl *Client) JsonNumMultBy(key, path string, num int) *redis.StringCmd {
	return jsonNumMultByExecute(cl, key, path, num)
}

func (cl *Client) JsonStrAppend(key, path, appendString string) *redis.IntCmd {
	return jsonStrAppendExecute(cl, key, path, appendString)
}

func (cl *Client) JsonStrLen(key, path, appendString string) *redis.IntCmd {
	return jsonStrAppendExecute(cl, key, path, appendString)
}

func (cl *Client) JsonArrAppend(key, path, appendString string) *redis.IntCmd {
	return jsonStrAppendExecute(cl, key, path, appendString)
}