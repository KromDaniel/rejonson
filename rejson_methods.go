package rejonson

import (
	"github.com/go-redis/redis"
)

type redisProcessor interface {
	Process(cmd redis.Cmder) error
}

func concatWithCmd(cmdName string, args []interface{}) []interface{} {
	res := make([]interface{}, 1)
	res[0] = cmdName
	for _, v := range args {
		if str, ok := v.(string); ok {
			if len(str) == 0 {
				continue
			}
		}
		res = append(res, v)
	}
	return res
}

func jsonDelExecute(c redisProcessor, args ...interface{}) *redis.IntCmd {
	cmd := redis.NewIntCmd(concatWithCmd("JSON.DEL", args)...)
	c.Process(cmd)
	return cmd
}

func jsonGetExecute(c redisProcessor, args ...interface{}) *redis.StringCmd {
	cmd := redis.NewStringCmd(concatWithCmd("JSON.GET", args)...)
	c.Process(cmd)
	return cmd
}

func jsonSetExecute(c redisProcessor, args ...interface{}) *redis.StatusCmd {
	cmd := redis.NewStatusCmd(concatWithCmd("JSON.SET", args)...)
	c.Process(cmd)
	return cmd
}

func jsonMGetExecute(c redisProcessor, args ...interface{}) *redis.StringSliceCmd {
	cmd := redis.NewStringSliceCmd(concatWithCmd("JSON.MGET", args)...)
	c.Process(cmd)
	return cmd
}

func jsonTypeExecute(c redisProcessor, args ...interface{}) *redis.StringCmd {
	cmd := redis.NewStringCmd(concatWithCmd("JSON.TYPE", args)...)
	c.Process(cmd)
	return cmd
}

func jsonNumIncrByExecute(c redisProcessor, args ...interface{}) *redis.StringCmd {
	cmd := redis.NewStringCmd(concatWithCmd("JSON.NUMINCRBY", args)...)
	c.Process(cmd)
	return cmd
}

func jsonNumMultByExecute(c redisProcessor, args ...interface{}) *redis.StringCmd {
	cmd := redis.NewStringCmd(concatWithCmd("JSON.NUMMULTBY", args)...)
	c.Process(cmd)
	return cmd
}

func jsonStrAppendExecute(c redisProcessor, args ...interface{}) *redis.IntCmd {
	cmd := redis.NewIntCmd(concatWithCmd("JSON.STRAPPEND", args)...)
	c.Process(cmd)
	return cmd
}

func jsonStrLenExecute(c redisProcessor, args ...interface{}) *redis.IntCmd {
	cmd := redis.NewIntCmd(concatWithCmd("JSON.STRLEN", args)...)
	c.Process(cmd)
	return cmd
}

func jsonArrAppendExecute(c redisProcessor, args ...interface{}) *redis.IntCmd {
	cmd := redis.NewIntCmd(concatWithCmd("JSON.ARRAPPEND", args)...)
	c.Process(cmd)
	return cmd
}

func jsoArrIndexExecute(c redisProcessor, args ...interface{}) *redis.IntCmd {
	cmd := redis.NewIntCmd(concatWithCmd("JSON.ARRINDEX", args)...)
	c.Process(cmd)
	return cmd
}

func jsonArrInsertExecute(c redisProcessor, args ...interface{}) *redis.IntCmd {
	cmd := redis.NewIntCmd(concatWithCmd("JSON.ARRINSERT", args)...)
	c.Process(cmd)
	return cmd
}

func jsonArrLenExecute(c redisProcessor, args ...interface{}) *redis.IntCmd {
	cmd := redis.NewIntCmd(concatWithCmd("JSON.ARRLEN", args)...)
	c.Process(cmd)
	return cmd
}

func jsonArrPopExecute(c redisProcessor, args ...interface{}) *redis.StringCmd {
	cmd := redis.NewStringCmd(concatWithCmd("JSON.ARRPOP", args)...)
	c.Process(cmd)
	return cmd
}

func jsonArrTrimExecute(c redisProcessor, args ...interface{}) *redis.IntCmd {
	cmd := redis.NewIntCmd(concatWithCmd("JSON.ARRTRIM", args)...)
	c.Process(cmd)
	return cmd
}

func jsonObjKeysExecute(c redisProcessor, args ...interface{}) *redis.StringSliceCmd {
	cmd := redis.NewStringSliceCmd(concatWithCmd("JSON.OBJKEYS", args)...)
	c.Process(cmd)
	return cmd
}

func jsonObjLen(c redisProcessor, args ...interface{}) *redis.IntCmd {
	cmd := redis.NewIntCmd(concatWithCmd("JSON.OBJLEN", args)...)
	c.Process(cmd)
	return cmd
}



