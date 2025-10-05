package main

import (
	"github.com/go-redis/redis"
	"reflect"
)

var (
	keyArg = RejsonArg{
		Name: "key",
		Type: reflect.String,
	}
	defaultVariadic = RejsonArg{
		Name:       "args",
		Type:       reflect.Interface,
		IsVariadic: true,
	}
)

var cmds = []RejsonCommand{
	{
		Name: "JsonGet",
		Cmd:  "JSON.GET",
		Args: []RejsonArg{
			keyArg,
			defaultVariadic,
		},
		CommandCtr: redis.NewStringCmd,
	},
	{
		Name: "JsonDel",
		Cmd:  "JSON.DEL",
		Args: []RejsonArg{
			keyArg,
			defaultVariadic,
		},
		CommandCtr: redis.NewIntCmd,
	},
	{
		Name: "JsonSet",
		Cmd:  "JSON.SET",
		Args: []RejsonArg{
			keyArg,
			{
				Name: "path",
				Type: reflect.String,
			},
			{
				Name: "json",
				Type: reflect.String,
			},
			defaultVariadic,
		},
		CommandCtr: redis.NewStatusCmd,
	},
	{
		Name: "JsonMGet",
		Cmd:  "JSON.MGET",
		Args: []RejsonArg{
			keyArg,
			defaultVariadic,
		},
		CommandCtr: redis.NewStringSliceCmd,
	},
	{
		Name: "JsonType",
		Cmd:  "JSON.TYPE",
		Args: []RejsonArg{
			keyArg,
			{
				Name: "path",
				Type: reflect.String,
			},
		},
		CommandCtr: redis.NewStringCmd,
	},
	{
		Name: "JsonNumIncrBy",
		Cmd:  "JSON.NUMINCRBY",
		Args: []RejsonArg{
			keyArg,
			{
				Name: "path",
				Type: reflect.String,
			},
			{
				Name: "value",
				Type: reflect.Int,
			},
		},
		CommandCtr: redis.NewStringCmd,
	},
	{
		Name: "JsonNumMultBy",
		Cmd:  "JSON.NUMMULTBY",
		Args: []RejsonArg{
			keyArg,
			{
				Name: "path",
				Type: reflect.String,
			},
			{
				Name: "value",
				Type: reflect.Int,
			},
		},
		CommandCtr: redis.NewStringCmd,
	},
	{
		Name: "JsonStrAppend",
		Cmd:  "JSON.STRAPPEND",
		Args: []RejsonArg{
			keyArg,
			{
				Name: "path",
				Type: reflect.String,
			},
			{
				Name: "value",
				Type: reflect.String,
			},
		},
		CommandCtr: redis.NewIntCmd,
	},
	{
		Name: "JsonStrLen",
		Cmd:  "JSON.STRLEN",
		Args: []RejsonArg{
			keyArg,
			{
				Name: "path",
				Type: reflect.String,
			},
		},
		CommandCtr: redis.NewIntCmd,
	},
	{
		Name: "JsonArrAppend",
		Cmd:  "JSON.ARRAPPEND",
		Args: []RejsonArg{
			keyArg,
			{
				Name: "path",
				Type: reflect.String,
			},
			defaultVariadic,
		},
		CommandCtr: redis.NewIntCmd,
	},
	{
		Name: "JsonArrIndex",
		Cmd:  "JSON.ARRINDEX",
		Args: []RejsonArg{
			keyArg,
			{
				Name: "path",
				Type: reflect.String,
			},
			{
				Name: "value",
				Type: reflect.Interface,
			},
			{
				Name:       "startAndStop",
				Type:       reflect.Interface,
				IsVariadic: true,
			},
		},
		CommandCtr: redis.NewIntCmd,
	},
	{
		Name: "JsonArrInsert",
		Cmd:  "JSON.ARRINSERT",
		Args: []RejsonArg{
			keyArg,
			{
				Name: "path",
				Type: reflect.String,
			},
			{
				Name: "index",
				Type: reflect.Int,
			},
			{
				Name:       "values",
				Type:       reflect.Interface,
				IsVariadic: true,
			},
		},
		CommandCtr: redis.NewIntCmd,
	},
	{
		Name: "JsonArrLen",
		Cmd:  "JSON.ARRLEN",
		Args: []RejsonArg{
			keyArg,
			{
				Name: "path",
				Type: reflect.String,
			},
		},
		CommandCtr: redis.NewIntCmd,
	},
	{
		Name: "JsonArrPop",
		Cmd:  "JSON.ARRPOP",
		Args: []RejsonArg{
			keyArg,
			{
				Name: "path",
				Type: reflect.String,
			},
			{
				Name: "index",
				Type: reflect.Int,
			},
		},
		CommandCtr: redis.NewStringCmd,
	},
	{
		Name: "JsonArrTrim",
		Cmd:  "JSON.ARRTRIM",
		Args: []RejsonArg{
			keyArg,
			{
				Name: "path",
				Type: reflect.String,
			},
			{
				Name: "start",
				Type: reflect.Int,
			},
			{
				Name: "stop",
				Type: reflect.Int,
			},
		},
		CommandCtr: redis.NewIntCmd,
	},
	{
		Name: "JsonObjKeys",
		Cmd:  "JSON.OBJKEYS",
		Args: []RejsonArg{
			keyArg,
			{
				Name: "path",
				Type: reflect.String,
			},
		},
		CommandCtr: redis.NewStringSliceCmd,
	},
	{
		Name: "JsonObjLen",
		Cmd:  "JSON.OBJLEN",
		Args: []RejsonArg{
			keyArg,
			{
				Name: "path",
				Type: reflect.String,
			},
		},
		CommandCtr: redis.NewIntCmd,
	},
	{
		Name: "JsonClear",
		Cmd:  "JSON.CLEAR",
		Args: []RejsonArg{
			keyArg,
			{
				Name: "path",
				Type: reflect.String,
			},
		},
		CommandCtr: redis.NewIntCmd,
	},
	{
		Name: "JsonForget",
		Cmd:  "JSON.FORGET",
		Args: []RejsonArg{
			keyArg,
			{
				Name: "path",
				Type: reflect.String,
			},
		},
		CommandCtr: redis.NewIntCmd,
	},
	{
		Name: "JsonMerge",
		Cmd:  "JSON.MERGE",
		Args: []RejsonArg{
			keyArg,
			{
				Name: "path",
				Type: reflect.String,
			},
			{
				Name: "value",
				Type: reflect.String,
			},
		},
		CommandCtr: redis.NewStatusCmd,
	},
	{
		Name: "JsonMSet",
		Cmd:  "JSON.MSET",
		Args: []RejsonArg{
			defaultVariadic,
		},
		CommandCtr: redis.NewStatusCmd,
	},
	{
		Name: "JsonSetMode",
		Cmd:  "JSON.SET",
		Args: []RejsonArg{
			keyArg,
			{
				Name: "path",
				Type: reflect.String,
			},
			{
				Name: "value",
				Type: reflect.String,
			},
			{
				Name: "mode",
				Type: reflect.String,
			},
		},
		CommandCtr: redis.NewStatusCmd,
	},
	{
		Name: "JsonToggle",
		Cmd:  "JSON.TOGGLE",
		Args: []RejsonArg{
			keyArg,
			{
				Name: "path",
				Type: reflect.String,
			},
		},
		CommandCtr: redis.NewStringCmd,
	},
}
