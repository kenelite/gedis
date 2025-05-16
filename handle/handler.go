package handle

import (
	"github.com/kenelite/gedis/server"
	"sync"
)

var Handlers = map[string]func([]server.Value) server.Value{
	"PING":    ping,
	"SET":     set,
	"GET":     get,
	"HSET":    hset,
	"HGET":    hget,
	"HGETALL": hgetall,
}

func ping(args []server.Value) server.Value {
	if len(args) == 0 {
		return server.Value{Typ: "string", Str: "PONG"}
	}

	return server.Value{Typ: "string", Str: args[0].Bulk}
}

var SETs = map[string]string{}
var SETsMu = sync.RWMutex{}

func set(args []server.Value) server.Value {
	if len(args) != 2 {
		return server.Value{Typ: "error", Str: "ERR wrong number of arguments for 'set' command"}
	}

	key := args[0].Bulk
	value := args[1].Bulk

	SETsMu.Lock()
	SETs[key] = value
	SETsMu.Unlock()

	return server.Value{Typ: "string", Str: "OK"}
}

func get(args []server.Value) server.Value {
	if len(args) != 1 {
		return server.Value{Typ: "error", Str: "ERR wrong number of arguments for 'get' command"}
	}

	key := args[0].Bulk

	SETsMu.RLock()
	value, ok := SETs[key]
	SETsMu.RUnlock()

	if !ok {
		return server.Value{Typ: "null"}
	}

	return server.Value{Typ: "bulk", Bulk: value}
}

var HSETs = map[string]map[string]string{}
var HSETsMu = sync.RWMutex{}

func hset(args []server.Value) server.Value {
	if len(args) != 3 {
		return server.Value{Typ: "error", Str: "ERR wrong number of arguments for 'hset' command"}
	}

	hash := args[0].Bulk
	key := args[1].Bulk
	value := args[2].Bulk

	HSETsMu.Lock()
	if _, ok := HSETs[hash]; !ok {
		HSETs[hash] = map[string]string{}
	}
	HSETs[hash][key] = value
	HSETsMu.Unlock()

	return server.Value{Typ: "string", Str: "OK"}
}

func hget(args []server.Value) server.Value {
	if len(args) != 2 {
		return server.Value{Typ: "error", Str: "ERR wrong number of arguments for 'hget' command"}
	}

	hash := args[0].Bulk
	key := args[1].Bulk

	HSETsMu.RLock()
	value, ok := HSETs[hash][key]
	HSETsMu.RUnlock()

	if !ok {
		return server.Value{Typ: "null"}
	}

	return server.Value{Typ: "bulk", Bulk: value}
}

func hgetall(args []server.Value) server.Value {
	if len(args) != 1 {
		return server.Value{Typ: "error", Str: "ERR wrong number of arguments for 'hgetall' command"}
	}

	hash := args[0].Bulk

	HSETsMu.RLock()
	value, ok := HSETs[hash]
	HSETsMu.RUnlock()

	if !ok {
		return server.Value{Typ: "null"}
	}

	values := []server.Value{}
	for k, v := range value {
		values = append(values, server.Value{Typ: "bulk", Bulk: k})
		values = append(values, server.Value{Typ: "bulk", Bulk: v})
	}

	return server.Value{Typ: "array", Array: values}
}
