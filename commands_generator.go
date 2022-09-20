//go:build ignore
// +build ignore

//go:generate go run commands_generator.go

package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"go/format"
	"html/template"
	"io"
	"os"
	"reflect"
	"runtime"
	"strings"
)

var fnsMappings = template.FuncMap{
	"argType": func(v reflect.Kind) string {
		if v == reflect.Interface {
			return v.String() + "{}"
		}
		return v.String()
	},
	"ctr": func(ctr interface{}) string {
		v := reflect.ValueOf(ctr)
		if v.Kind() != reflect.Func {
			panic(fmt.Errorf("CommandCtr has to be func and not %s", v.Kind().String()))
		}
		fullPaths := strings.Split(runtime.FuncForPC(v.Pointer()).Name(), ".")
		return "redis." + fullPaths[len(fullPaths)-1]
	},
	"ctrReturn": func(ctr interface{}) string {
		v := reflect.ValueOf(ctr)
		if v.Kind() != reflect.Func {
			panic(fmt.Errorf("CommandCtr has to be func and not %s", v.Kind().String()))
		}
		return v.Type().Out(0).String()
	},
	"dict": func(values ...interface{}) (map[string]interface{}, error) {
		if len(values)%2 != 0 {
			return nil, errors.New("invalid dict call")
		}
		dict := make(map[string]interface{}, len(values)/2)
		for i := 0; i < len(values); i += 2 {
			key, ok := values[i].(string)
			if !ok {
				return nil, errors.New("dict keys must be strings")
			}
			dict[key] = values[i+1]
		}
		return dict, nil
	},
	"lastArg": func(cmd RejsonCommand) RejsonArg {
		return cmd.Args[len(cmd.Args)-1]
	},
}

type RejsonArg struct {
	Name       string
	Type       reflect.Kind
	IsVariadic bool
}

type RejsonCommand struct {
	Name       string
	Cmd        string
	Args       []RejsonArg
	CommandCtr interface{}
}

func (r *RejsonCommand) HasVariadic() bool {
	for _, arg := range r.Args {
		if arg.IsVariadic {
			return true
		}
	}
	return false
}

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
}

const headerTemplate = `
// Code generated by rejonson. DO NOT EDIT.

package rejonson

import ({{range .}}
"{{ .}}"
{{end }})

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
`

const functionalTemplate = `
{{define "argsTemplate"}} {{range $i,$self := .}} {{if not (eq $i 0)}}, {{end}} {{.Name}} {{if .IsVariadic}}...{{end}}{{argType .Type}} {{end}} {{end}}
{{define "argsCallTemplate"}} {{range $i,$self := .}} {{if not (eq $i 0)}}, {{end}} {{.Name}} {{if .IsVariadic}}...{{end}} {{end}} {{end}}
{{define "createArgsSlice"}}
	argsSlice := make([]interface{}, 0, {{len .Args}} {{if .HasVariadic}} -1 + len({{(lastArg .).Name}}) {{end}} )
	argsSlice = append(argsSlice, "{{.Cmd}}")
	{{range .Args -}}
		argsSlice = append(argsSlice, {{.Name}}{{if .IsVariadic}}...{{end}})
	{{end}}
{{end}}

{{define "functionApi"}}
	func {{.Name}}(c RedisProcessor, {{template "argsTemplate" .Args}} ) {{ctrReturn .CommandCtr}} {
		{{template "createArgsSlice" .}}
		cmd := {{ctr .CommandCtr}}(argsSlice...)
		// ignore the error since cmd.Error() contains it
		_ = c.Process(cmd)
		return cmd
	}
{{end}}

{{define "clientMethods"}}
	func (cmder {{.Caller}}) {{.Name}}({{template "argsTemplate" .Args}} ) {{ctrReturn .CommandCtr}} {
		return {{.Name}}(cmder, {{template "argsCallTemplate" .Args}})
	}
{{end}}

{{range .}}
	{{template "functionApi" .}}
{{end}}

{{range .}}
	{{template "clientMethods" dict "Name" .Name "Args" .Args "CommandCtr" .CommandCtr "Caller" "*Client" }}
{{end}}

{{range .}}
	{{template "clientMethods" dict "Name" .Name "Args" .Args "CommandCtr" .CommandCtr "Caller" "*Pipeline" }}
{{end}}
`

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func createHeader(writer io.Writer) {
	t, err := template.New("header").Parse(headerTemplate)
	must(err)

	must(t.Execute(writer, []string{"github.com/go-redis/redis"}))
}

func createFuncs(writer io.Writer) {
	t, err := template.New("functionalTemplate").Funcs(fnsMappings).Parse(functionalTemplate)

	must(err)
	must(t.Execute(writer, cmds))
}

func main() {

	f, err := os.OpenFile("./v6.go", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	must(err)
	defer func() {
		must(f.Close())
	}()
	out := new(bytes.Buffer)
	createHeader(out)
	createFuncs(out)

	formatted, err := format.Source(out.Bytes())
	if err != nil {
		panic(fmt.Errorf("%s\n\n%w", out.String(), err))
	}
	f.Write(formatted)
}
