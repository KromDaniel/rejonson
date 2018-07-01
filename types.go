package rejonson

import "github.com/go-redis/redis"

type redisProcessor struct {
	Process func(cmd redis.Cmder) error
}
type Client struct {
	*redis.Client
 	*redisProcessor
}

type Pipeline struct {
	redis.Pipeliner
	*redisProcessor
}

func (cl *Client) Pipeline() *Pipeline {
	pip := cl.Client.Pipeline()
	return ExtendPipeline(pip)
}

func (cl *Client) TXPipeline() *Pipeline {
	pip := cl.Client.TxPipeline()
	return ExtendPipeline(pip)
}
func (pl *Pipeline) Pipeline() *Pipeline {
	pip := pl.Pipeliner.Pipeline()
	return ExtendPipeline(pip)
}
func (cl *redisProcessor) JsonDel(key, path string) *redis.IntCmd {
	return jsonDelExecute(cl, key, path)
}

func (cl *redisProcessor) JsonGet(key string, args ...interface{}) *redis.StringCmd {
	return jsonGetExecute(cl, append([]interface{}{key}, args...)...)
}

func (cl *redisProcessor) JsonSet(key, path, json string, args ...interface{}) *redis.StatusCmd {
	return jsonSetExecute(cl, append([]interface{}{key, path, json}, args...)...)
}

func (cl *redisProcessor) JsonMGet(key string, args ...interface{}) *redis.StringSliceCmd {
	return jsonMGetExecute(cl, append([]interface{}{key}, args...)...)
}

func (cl *redisProcessor) JsonType(key, path string) *redis.StringCmd {
	return jsonTypeExecute(cl, key, path)
}

func (cl *redisProcessor) JsonNumIncrBy(key, path string, num int) *redis.StringCmd {
	return jsonNumIncrByExecute(cl, key, path, num)
}

func (cl *redisProcessor) JsonNumMultBy(key, path string, num int) *redis.StringCmd {
	return jsonNumMultByExecute(cl, key, path, num)
}

func (cl *redisProcessor) JsonStrAppend(key, path, appendString string) *redis.IntCmd {
	return jsonStrAppendExecute(cl, key, path, appendString)
}

func (cl *redisProcessor) JsonStrLen(key, path string) *redis.IntCmd {
	return jsonStrLenExecute(cl, key, path)
}

func (cl *redisProcessor) JsonArrAppend(key, path string, jsons ...interface{}) *redis.IntCmd {
	return jsonArrAppendExecute(cl, append([]interface{}{key, path}, jsons...)...)
}

func (cl *redisProcessor) JsonArrIndex(key, path string, jsonScalar interface{}, startAndStop ...interface{}) *redis.IntCmd {
	return jsoArrIndexExecute(cl, append([]interface{}{key, path, jsonScalar}, startAndStop...)...)
}

func (cl *redisProcessor) JsonArrInsert(key, path string, index int, jsons ...interface{}) *redis.IntCmd {
	return jsoArrIndexExecute(cl, append([]interface{}{key, path, index}, jsons...)...)
}

func (cl *redisProcessor) JsonArrLen(key, path string) *redis.IntCmd {
	return jsonArrLenExecute(cl, key, path)
}

func (cl *redisProcessor) JsonArrPop(key, path string, index int) *redis.StringCmd {
	return jsonArrPopExecute(cl, key, path, index)
}

func (cl *redisProcessor) JsonArrTrim(key, path string, start, stop int) *redis.IntCmd {
	return jsonArrTrimExecute(cl, key, path, start, stop)
}

func (cl *redisProcessor) JsonObjKeys(key, path string) *redis.StringSliceCmd {
	return jsonObjKeysExecute(cl, key, path)
}

func (cl *redisProcessor) JsonObjLen(key, path string) *redis.IntCmd {
	return jsonObjLen(cl, key, path)
}