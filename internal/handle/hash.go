package handle

import (
	"github.com/kenelite/gedis/internal/response"
	"sync"
)

var (
	HSETs   = map[string]map[string]string{}
	HSETsMu = sync.RWMutex{}
)

func Hset(args []response.Value) response.Value {
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

func Hget(args []response.Value) response.Value {
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

func Hgetall(args []response.Value) response.Value {
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
