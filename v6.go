// Code generated by rejonson. DO NOT EDIT.

package rejonson

import (
	"github.com/go-redis/redis"
)

// RedisProcessor is redis client or pipeline instance that will process a command
type RedisProcessor interface {
	Process(cmd redis.Cmder) error
}

/*
Client is an extended redis.Client, stores a pointer to the original redis.Client
*/
type Client struct {
	*redis.Client
}

/*
Pipeline is an extended redis.Pipeline, stores a pointer to the original redis.Pipeliner
*/
type Pipeline struct {
	redis.Pipeliner
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

func ExtendClient(client *redis.Client) *Client {
	return &Client{
		client,
	}
}

func ExtendPipeline(pipeline redis.Pipeliner) *Pipeline {
	return &Pipeline{
		pipeline,
	}
}

func JsonGet(c RedisProcessor, key string, args ...interface{}) *redis.StringCmd {

	argsSlice := make([]interface{}, 0, 2-1+len(args))
	argsSlice = append(argsSlice, "JSON.GET")
	argsSlice = append(argsSlice, key)
	argsSlice = append(argsSlice, args...)

	cmd := redis.NewStringCmd(argsSlice...)
	// ignore the error since cmd.Error() contains it
	_ = c.Process(cmd)
	return cmd
}

func JsonDel(c RedisProcessor, key string, args ...interface{}) *redis.IntCmd {

	argsSlice := make([]interface{}, 0, 2-1+len(args))
	argsSlice = append(argsSlice, "JSON.DEL")
	argsSlice = append(argsSlice, key)
	argsSlice = append(argsSlice, args...)

	cmd := redis.NewIntCmd(argsSlice...)
	// ignore the error since cmd.Error() contains it
	_ = c.Process(cmd)
	return cmd
}

func JsonSet(c RedisProcessor, key string, path string, json string, args ...interface{}) *redis.StatusCmd {

	argsSlice := make([]interface{}, 0, 4-1+len(args))
	argsSlice = append(argsSlice, "JSON.SET")
	argsSlice = append(argsSlice, key)
	argsSlice = append(argsSlice, path)
	argsSlice = append(argsSlice, json)
	argsSlice = append(argsSlice, args...)

	cmd := redis.NewStatusCmd(argsSlice...)
	// ignore the error since cmd.Error() contains it
	_ = c.Process(cmd)
	return cmd
}

func JsonMGet(c RedisProcessor, key string, args ...interface{}) *redis.StringSliceCmd {

	argsSlice := make([]interface{}, 0, 2-1+len(args))
	argsSlice = append(argsSlice, "JSON.MGET")
	argsSlice = append(argsSlice, key)
	argsSlice = append(argsSlice, args...)

	cmd := redis.NewStringSliceCmd(argsSlice...)
	// ignore the error since cmd.Error() contains it
	_ = c.Process(cmd)
	return cmd
}

func JsonType(c RedisProcessor, key string, path string) *redis.StringCmd {

	argsSlice := make([]interface{}, 0, 2)
	argsSlice = append(argsSlice, "JSON.TYPE")
	argsSlice = append(argsSlice, key)
	argsSlice = append(argsSlice, path)

	cmd := redis.NewStringCmd(argsSlice...)
	// ignore the error since cmd.Error() contains it
	_ = c.Process(cmd)
	return cmd
}

func JsonNumIncrBy(c RedisProcessor, key string, path string, value int) *redis.StringCmd {

	argsSlice := make([]interface{}, 0, 3)
	argsSlice = append(argsSlice, "JSON.NUMINCRBY")
	argsSlice = append(argsSlice, key)
	argsSlice = append(argsSlice, path)
	argsSlice = append(argsSlice, value)

	cmd := redis.NewStringCmd(argsSlice...)
	// ignore the error since cmd.Error() contains it
	_ = c.Process(cmd)
	return cmd
}

func JsonNumMultBy(c RedisProcessor, key string, path string, value int) *redis.StringCmd {

	argsSlice := make([]interface{}, 0, 3)
	argsSlice = append(argsSlice, "JSON.NUMMULTBY")
	argsSlice = append(argsSlice, key)
	argsSlice = append(argsSlice, path)
	argsSlice = append(argsSlice, value)

	cmd := redis.NewStringCmd(argsSlice...)
	// ignore the error since cmd.Error() contains it
	_ = c.Process(cmd)
	return cmd
}

func JsonStrAppend(c RedisProcessor, key string, path string, value string) *redis.IntCmd {

	argsSlice := make([]interface{}, 0, 3)
	argsSlice = append(argsSlice, "JSON.STRAPPEND")
	argsSlice = append(argsSlice, key)
	argsSlice = append(argsSlice, path)
	argsSlice = append(argsSlice, value)

	cmd := redis.NewIntCmd(argsSlice...)
	// ignore the error since cmd.Error() contains it
	_ = c.Process(cmd)
	return cmd
}

func JsonStrLen(c RedisProcessor, key string, path string) *redis.IntCmd {

	argsSlice := make([]interface{}, 0, 2)
	argsSlice = append(argsSlice, "JSON.STRLEN")
	argsSlice = append(argsSlice, key)
	argsSlice = append(argsSlice, path)

	cmd := redis.NewIntCmd(argsSlice...)
	// ignore the error since cmd.Error() contains it
	_ = c.Process(cmd)
	return cmd
}

func JsonArrAppend(c RedisProcessor, key string, path string, args ...interface{}) *redis.IntCmd {

	argsSlice := make([]interface{}, 0, 3-1+len(args))
	argsSlice = append(argsSlice, "JSON.ARRAPPEND")
	argsSlice = append(argsSlice, key)
	argsSlice = append(argsSlice, path)
	argsSlice = append(argsSlice, args...)

	cmd := redis.NewIntCmd(argsSlice...)
	// ignore the error since cmd.Error() contains it
	_ = c.Process(cmd)
	return cmd
}

func JsonArrIndex(c RedisProcessor, key string, path string, value interface{}, startAndStop ...interface{}) *redis.IntCmd {

	argsSlice := make([]interface{}, 0, 4-1+len(startAndStop))
	argsSlice = append(argsSlice, "JSON.ARRINDEX")
	argsSlice = append(argsSlice, key)
	argsSlice = append(argsSlice, path)
	argsSlice = append(argsSlice, value)
	argsSlice = append(argsSlice, startAndStop...)

	cmd := redis.NewIntCmd(argsSlice...)
	// ignore the error since cmd.Error() contains it
	_ = c.Process(cmd)
	return cmd
}

func JsonArrInsert(c RedisProcessor, key string, path string, index int, values ...interface{}) *redis.IntCmd {

	argsSlice := make([]interface{}, 0, 4-1+len(values))
	argsSlice = append(argsSlice, "JSON.ARRINSERT")
	argsSlice = append(argsSlice, key)
	argsSlice = append(argsSlice, path)
	argsSlice = append(argsSlice, index)
	argsSlice = append(argsSlice, values...)

	cmd := redis.NewIntCmd(argsSlice...)
	// ignore the error since cmd.Error() contains it
	_ = c.Process(cmd)
	return cmd
}

func JsonArrLen(c RedisProcessor, key string, path string) *redis.IntCmd {

	argsSlice := make([]interface{}, 0, 2)
	argsSlice = append(argsSlice, "JSON.ARRLEN")
	argsSlice = append(argsSlice, key)
	argsSlice = append(argsSlice, path)

	cmd := redis.NewIntCmd(argsSlice...)
	// ignore the error since cmd.Error() contains it
	_ = c.Process(cmd)
	return cmd
}

func JsonArrPop(c RedisProcessor, key string, path string, index int) *redis.StringCmd {

	argsSlice := make([]interface{}, 0, 3)
	argsSlice = append(argsSlice, "JSON.ARRPOP")
	argsSlice = append(argsSlice, key)
	argsSlice = append(argsSlice, path)
	argsSlice = append(argsSlice, index)

	cmd := redis.NewStringCmd(argsSlice...)
	// ignore the error since cmd.Error() contains it
	_ = c.Process(cmd)
	return cmd
}

func JsonArrTrim(c RedisProcessor, key string, path string, start int, stop int) *redis.IntCmd {

	argsSlice := make([]interface{}, 0, 4)
	argsSlice = append(argsSlice, "JSON.ARRTRIM")
	argsSlice = append(argsSlice, key)
	argsSlice = append(argsSlice, path)
	argsSlice = append(argsSlice, start)
	argsSlice = append(argsSlice, stop)

	cmd := redis.NewIntCmd(argsSlice...)
	// ignore the error since cmd.Error() contains it
	_ = c.Process(cmd)
	return cmd
}

func JsonObjKeys(c RedisProcessor, key string, path string) *redis.StringSliceCmd {

	argsSlice := make([]interface{}, 0, 2)
	argsSlice = append(argsSlice, "JSON.OBJKEYS")
	argsSlice = append(argsSlice, key)
	argsSlice = append(argsSlice, path)

	cmd := redis.NewStringSliceCmd(argsSlice...)
	// ignore the error since cmd.Error() contains it
	_ = c.Process(cmd)
	return cmd
}

func JsonObjLen(c RedisProcessor, key string, path string) *redis.IntCmd {

	argsSlice := make([]interface{}, 0, 2)
	argsSlice = append(argsSlice, "JSON.OBJLEN")
	argsSlice = append(argsSlice, key)
	argsSlice = append(argsSlice, path)

	cmd := redis.NewIntCmd(argsSlice...)
	// ignore the error since cmd.Error() contains it
	_ = c.Process(cmd)
	return cmd
}

func (cmder *Client) JsonGet(key string, args ...interface{}) *redis.StringCmd {
	return JsonGet(cmder, key, args...)
}

func (cmder *Client) JsonDel(key string, args ...interface{}) *redis.IntCmd {
	return JsonDel(cmder, key, args...)
}

func (cmder *Client) JsonSet(key string, path string, json string, args ...interface{}) *redis.StatusCmd {
	return JsonSet(cmder, key, path, json, args...)
}

func (cmder *Client) JsonMGet(key string, args ...interface{}) *redis.StringSliceCmd {
	return JsonMGet(cmder, key, args...)
}

func (cmder *Client) JsonType(key string, path string) *redis.StringCmd {
	return JsonType(cmder, key, path)
}

func (cmder *Client) JsonNumIncrBy(key string, path string, value int) *redis.StringCmd {
	return JsonNumIncrBy(cmder, key, path, value)
}

func (cmder *Client) JsonNumMultBy(key string, path string, value int) *redis.StringCmd {
	return JsonNumMultBy(cmder, key, path, value)
}

func (cmder *Client) JsonStrAppend(key string, path string, value string) *redis.IntCmd {
	return JsonStrAppend(cmder, key, path, value)
}

func (cmder *Client) JsonStrLen(key string, path string) *redis.IntCmd {
	return JsonStrLen(cmder, key, path)
}

func (cmder *Client) JsonArrAppend(key string, path string, args ...interface{}) *redis.IntCmd {
	return JsonArrAppend(cmder, key, path, args...)
}

func (cmder *Client) JsonArrIndex(key string, path string, value interface{}, startAndStop ...interface{}) *redis.IntCmd {
	return JsonArrIndex(cmder, key, path, value, startAndStop...)
}

func (cmder *Client) JsonArrInsert(key string, path string, index int, values ...interface{}) *redis.IntCmd {
	return JsonArrInsert(cmder, key, path, index, values...)
}

func (cmder *Client) JsonArrLen(key string, path string) *redis.IntCmd {
	return JsonArrLen(cmder, key, path)
}

func (cmder *Client) JsonArrPop(key string, path string, index int) *redis.StringCmd {
	return JsonArrPop(cmder, key, path, index)
}

func (cmder *Client) JsonArrTrim(key string, path string, start int, stop int) *redis.IntCmd {
	return JsonArrTrim(cmder, key, path, start, stop)
}

func (cmder *Client) JsonObjKeys(key string, path string) *redis.StringSliceCmd {
	return JsonObjKeys(cmder, key, path)
}

func (cmder *Client) JsonObjLen(key string, path string) *redis.IntCmd {
	return JsonObjLen(cmder, key, path)
}

func (cmder *Pipeline) JsonGet(key string, args ...interface{}) *redis.StringCmd {
	return JsonGet(cmder, key, args...)
}

func (cmder *Pipeline) JsonDel(key string, args ...interface{}) *redis.IntCmd {
	return JsonDel(cmder, key, args...)
}

func (cmder *Pipeline) JsonSet(key string, path string, json string, args ...interface{}) *redis.StatusCmd {
	return JsonSet(cmder, key, path, json, args...)
}

func (cmder *Pipeline) JsonMGet(key string, args ...interface{}) *redis.StringSliceCmd {
	return JsonMGet(cmder, key, args...)
}

func (cmder *Pipeline) JsonType(key string, path string) *redis.StringCmd {
	return JsonType(cmder, key, path)
}

func (cmder *Pipeline) JsonNumIncrBy(key string, path string, value int) *redis.StringCmd {
	return JsonNumIncrBy(cmder, key, path, value)
}

func (cmder *Pipeline) JsonNumMultBy(key string, path string, value int) *redis.StringCmd {
	return JsonNumMultBy(cmder, key, path, value)
}

func (cmder *Pipeline) JsonStrAppend(key string, path string, value string) *redis.IntCmd {
	return JsonStrAppend(cmder, key, path, value)
}

func (cmder *Pipeline) JsonStrLen(key string, path string) *redis.IntCmd {
	return JsonStrLen(cmder, key, path)
}

func (cmder *Pipeline) JsonArrAppend(key string, path string, args ...interface{}) *redis.IntCmd {
	return JsonArrAppend(cmder, key, path, args...)
}

func (cmder *Pipeline) JsonArrIndex(key string, path string, value interface{}, startAndStop ...interface{}) *redis.IntCmd {
	return JsonArrIndex(cmder, key, path, value, startAndStop...)
}

func (cmder *Pipeline) JsonArrInsert(key string, path string, index int, values ...interface{}) *redis.IntCmd {
	return JsonArrInsert(cmder, key, path, index, values...)
}

func (cmder *Pipeline) JsonArrLen(key string, path string) *redis.IntCmd {
	return JsonArrLen(cmder, key, path)
}

func (cmder *Pipeline) JsonArrPop(key string, path string, index int) *redis.StringCmd {
	return JsonArrPop(cmder, key, path, index)
}

func (cmder *Pipeline) JsonArrTrim(key string, path string, start int, stop int) *redis.IntCmd {
	return JsonArrTrim(cmder, key, path, start, stop)
}

func (cmder *Pipeline) JsonObjKeys(key string, path string) *redis.StringSliceCmd {
	return JsonObjKeys(cmder, key, path)
}

func (cmder *Pipeline) JsonObjLen(key string, path string) *redis.IntCmd {
	return JsonObjLen(cmder, key, path)
}
