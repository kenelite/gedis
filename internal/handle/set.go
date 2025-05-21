package handle

import (
	"github.com/kenelite/gedis/internal/response"
	"sync"
)

var (
	sets   = make(map[string]map[string]struct{})
	setsMu sync.RWMutex
)

func sadd(args []response.Value) response.Value {
	if len(args) < 2 {
		return response.Value{Typ: "error", Str: "ERR wrong number of arguments for 'SADD'"}
	}

	key := args[0].Bulk
	setsMu.Lock()
	defer setsMu.Unlock()

	if sets[key] == nil {
		sets[key] = make(map[string]struct{})
	}

	added := 0
	for _, arg := range args[1:] {
		val := arg.Bulk
		if _, exists := sets[key][val]; !exists {
			sets[key][val] = struct{}{}
			added++
		}
	}

	return response.Value{Typ: "integer", Num: added}
}

func srem(args []response.Value) response.Value {
	if len(args) < 2 {
		return response.Value{Typ: "error", Str: "ERR wrong number of arguments for 'SREM'"}
	}

	key := args[0].Bulk
	setsMu.Lock()
	defer setsMu.Unlock()

	members, exists := sets[key]
	if !exists {
		return response.Value{Typ: "integer", Num: 0}
	}

	removed := 0
	for _, arg := range args[1:] {
		val := arg.Bulk
		if _, ok := members[val]; ok {
			delete(members, val)
			removed++
		}
	}

	if len(members) == 0 {
		delete(sets, key)
	}

	return response.Value{Typ: "integer", Num: removed}
}

func smembers(args []response.Value) response.Value {
	if len(args) != 1 {
		return response.Value{Typ: "error", Str: "ERR wrong number of arguments for 'SMEMBERS'"}
	}

	key := args[0].Bulk
	setsMu.RLock()
	defer setsMu.RUnlock()

	members, exists := sets[key]
	if !exists {
		return response.Value{Typ: "array", Array: []response.Value{}}
	}

	var result []response.Value
	for member := range members {
		result = append(result, response.Value{Typ: "bulk", Bulk: member})
	}

	return response.Value{Typ: "array", Array: result}
}

func scard(args []response.Value) response.Value {
	if len(args) != 1 {
		return response.Value{Typ: "error", Str: "ERR wrong number of arguments for 'SCARD'"}
	}

	key := args[0].Bulk
	setsMu.RLock()
	defer setsMu.RUnlock()

	members, exists := sets[key]
	if !exists {
		return response.Value{Typ: "integer", Num: 0}
	}

	return response.Value{Typ: "integer", Num: len(members)}
}
