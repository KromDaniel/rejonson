package main

import (
	"reflect"
	"strings"
)

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

func (r *RejsonCommand) Last() RejsonArg {
	return r.Args[len(r.Args)-1]
}

func (r *RejsonCommand) HasVariadic() bool {
	for _, arg := range r.Args {
		if arg.IsVariadic {
			return true
		}
	}
	return false
}

type GoredisPackage struct {
	GoRedisImport string
	Output        string
	HasContext    bool
}

func (g *GoredisPackage) TestFile() string {
	return strings.ReplaceAll(g.Output, ".go", "_test.go")
}
