package handle

import (
	"github.com/kenelite/gedis/internal/response"
	"strconv"
	"sync"
)

var (
	Lists   = make(map[string][]string)
	ListsMu sync.RWMutex
)

func Lpush(args []response.Value) response.Value {
	if len(args) < 2 {
		return response.Value{Typ: "error", Str: "ERR wrong number of arguments for 'lpush'"}
	}
	key := args[0].Bulk

	ListsMu.Lock()
	defer ListsMu.Unlock()

	for i := len(args) - 1; i >= 1; i-- {
		Lists[key] = append([]string{args[i].Bulk}, Lists[key]...)
	}

	return response.Value{Typ: "integer", Num: len(Lists[key])}
}

func Rpush(args []response.Value) response.Value {
	if len(args) < 2 {
		return response.Value{Typ: "error", Str: "ERR wrong number of arguments for 'rpush'"}
	}
	key := args[0].Bulk

	ListsMu.Lock()
	defer ListsMu.Unlock()

	for i := 1; i < len(args); i++ {
		Lists[key] = append(Lists[key], args[i].Bulk)
	}

	return response.Value{Typ: "integer", Num: len(Lists[key])}
}

func Lpop(args []response.Value) response.Value {
	if len(args) != 1 {
		return response.Value{Typ: "error", Str: "ERR wrong number of arguments for 'lpop'"}
	}
	key := args[0].Bulk

	ListsMu.Lock()
	defer ListsMu.Unlock()

	list, ok := Lists[key]
	if !ok || len(list) == 0 {
		return response.Value{Typ: "null"}
	}

	val := list[0]
	Lists[key] = list[1:]

	return response.Value{Typ: "bulk", Bulk: val}
}

func Rpop(args []response.Value) response.Value {
	if len(args) != 1 {
		return response.Value{Typ: "error", Str: "ERR wrong number of arguments for 'rpop'"}
	}
	key := args[0].Bulk

	ListsMu.Lock()
	defer ListsMu.Unlock()

	list, ok := Lists[key]
	if !ok || len(list) == 0 {
		return response.Value{Typ: "null"}
	}

	val := list[len(list)-1]
	Lists[key] = list[:len(list)-1]

	return response.Value{Typ: "bulk", Bulk: val}
}

func Lrange(args []response.Value) response.Value {
	if len(args) != 3 {
		return response.Value{Typ: "error", Str: "ERR wrong number of arguments for 'lrange'"}
	}
	key := args[0].Bulk
	start, err1 := strconv.Atoi(args[1].Bulk)
	stop, err2 := strconv.Atoi(args[2].Bulk)
	if err1 != nil || err2 != nil {
		return response.Value{Typ: "error", Str: "ERR invalid start or stop"}
	}

	ListsMu.RLock()
	defer ListsMu.RUnlock()

	list, ok := Lists[key]
	if !ok {
		return response.Value{Typ: "array", Array: []response.Value{}}
	}

	if start < 0 {
		start = len(list) + start
	}
	if stop < 0 {
		stop = len(list) + stop
	}

	if start < 0 {
		start = 0
	}
	if stop >= len(list) {
		stop = len(list) - 1
	}

	if start > stop || start >= len(list) {
		return response.Value{Typ: "array", Array: []response.Value{}}
	}

	result := []response.Value{}
	for _, item := range list[start : stop+1] {
		result = append(result, response.Value{Typ: "bulk", Bulk: item})
	}

	return response.Value{Typ: "array", Array: result}
}
