package executor

import (
	"strings"

	"github.com/kenelite/gedis/internal/response"
)

type HandlerFunc func([]response.Value) response.Value

var handlers = make(map[string]HandlerFunc)

func Register(name string, fn HandlerFunc) {
	handlers[strings.ToUpper(name)] = fn
}

func Execute(cmd string, args []response.Value) response.Value {
	if fn, ok := handlers[strings.ToUpper(cmd)]; ok {
		return fn(args)
	}
	return response.Value{Typ: "error", Str: "ERR unknown command: " + cmd}
}
