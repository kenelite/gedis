package handle

import (
	"github.com/kenelite/gedis/internal/response"
	"sync"
)

var Handlers = map[string]func([]response.Value) response.Value{
	"PING":    ping,
	"SET":     set,
	"GET":     get,
	"HSET":    hset,
	"HGET":    hget,
	"HGETALL": hgetall,
}

func ping(args []response.Value) response.Value {
	if len(args) == 0 {
		return response.Value{Typ: "string", Str: "PONG"}
	}

	return response.Value{Typ: "string", Str: args[0].Bulk}
}

var SETs = map[string]string{}
var SETsMu = sync.RWMutex{}

func set(args []response.Value) response.Value {
	if len(args) != 2 {
		return response.Value{Typ: "error", Str: "ERR wrong number of arguments for 'set' command"}
	}

	key := args[0].Bulk
	value := args[1].Bulk

	SETsMu.Lock()
	SETs[key] = value
	SETsMu.Unlock()

	return response.Value{Typ: "string", Str: "OK"}
}

func get(args []response.Value) response.Value {
	if len(args) != 1 {
		return response.Value{Typ: "error", Str: "ERR wrong number of arguments for 'get' command"}
	}

	key := args[0].Bulk

	SETsMu.RLock()
	value, ok := SETs[key]
	SETsMu.RUnlock()

	if !ok {
		return response.Value{Typ: "null"}
	}

	return response.Value{Typ: "bulk", Bulk: value}
}

var HSETs = map[string]map[string]string{}
var HSETsMu = sync.RWMutex{}

func hset(args []response.Value) response.Value {
	if len(args) < 3 || len(args)%2 == 0 {
		return response.Value{Typ: "error", Str: "ERR wrong number of arguments for 'hset' command"}
	}

	hash := args[0].Bulk
	HSETsMu.Lock()

	for i := 1; i < len(args); i += 2 {
		key := args[i].Bulk
		value := args[i+1].Bulk

		if _, ok := HSETs[hash]; !ok {
			HSETs[hash] = map[string]string{}
		}
		HSETs[hash][key] = value

	}
	HSETsMu.Unlock()
	return response.Value{Typ: "string", Str: "OK"}
}

func hget(args []response.Value) response.Value {
	if len(args) != 2 {
		return response.Value{Typ: "error", Str: "ERR wrong number of arguments for 'hget' command"}
	}

	hash := args[0].Bulk
	key := args[1].Bulk

	HSETsMu.RLock()
	value, ok := HSETs[hash][key]
	HSETsMu.RUnlock()

	if !ok {
		return response.Value{Typ: "null"}
	}

	return response.Value{Typ: "bulk", Bulk: value}
}

func hgetall(args []response.Value) response.Value {
	if len(args) != 1 {
		return response.Value{Typ: "error", Str: "ERR wrong number of arguments for 'hgetall' command"}
	}

	hash := args[0].Bulk

	HSETsMu.RLock()
	value, ok := HSETs[hash]
	HSETsMu.RUnlock()

	if !ok {
		return response.Value{Typ: "null"}
	}

	values := []response.Value{}
	for k, v := range value {
		values = append(values, response.Value{Typ: "bulk", Bulk: k})
		values = append(values, response.Value{Typ: "bulk", Bulk: v})
	}

	return response.Value{Typ: "array", Array: values}
}
