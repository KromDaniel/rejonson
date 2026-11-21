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
	// Native is the go-redis v9 method replacing this command. Empty means
	// "JSON" + the name without the "Json" prefix; "-" means there is none.
	Native string
}

// DeprecationComment returns the "Deprecated:" doc paragraph for the command.
func (r *RejsonCommand) DeprecationComment() string {
	native := r.Native
	if native == "" {
		native = "JSON" + strings.TrimPrefix(r.Name, "Json")
	}
	msg := "Deprecated: use the typed github.com/KromDaniel/rejonson/v9/doc package"
	if native == "-" {
		return msg + " (for example Store.Update). " + r.Cmd + " is deprecated in RedisJSON and has no go-redis equivalent."
	}
	return msg + ", or the native go-redis method " + native + " (see redis.JSONCmdable)."
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
	// Deprecated marks the generated helpers as deprecated in favor of the
	// doc package and go-redis's native JSON commands.
	Deprecated bool
}

func (g *GoredisPackage) TestFile() string {
	return strings.ReplaceAll(g.Output, ".go", "_test.go")
}
