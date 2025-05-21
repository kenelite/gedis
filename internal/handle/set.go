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

func sunion(args []response.Value) response.Value {
	if len(args) < 1 {
		return response.Value{Typ: "error", Str: "ERR wrong number of arguments for 'SUNION'"}
	}

	setsMu.RLock()
	defer setsMu.RUnlock()

	union := make(map[string]struct{})
	for _, arg := range args {
		key := arg.Bulk
		if members, exists := sets[key]; exists {
			for member := range members {
				union[member] = struct{}{}
			}
		}
	}

	var result []response.Value
	for member := range union {
		result = append(result, response.Value{Typ: "bulk", Bulk: member})
	}
	return response.Value{Typ: "array", Array: result}
}

func sinter(args []response.Value) response.Value {
	if len(args) < 1 {
		return response.Value{Typ: "error", Str: "ERR wrong number of arguments for 'SINTER'"}
	}

	setsMu.RLock()
	defer setsMu.RUnlock()

	// Initialize result with first set
	firstKey := args[0].Bulk
	base, exists := sets[firstKey]
	if !exists {
		return response.Value{Typ: "array", Array: []response.Value{}}
	}

	intersection := make(map[string]struct{})
	for member := range base {
		intersection[member] = struct{}{}
	}

	// Intersect with all other sets
	for _, arg := range args[1:] {
		key := arg.Bulk
		curr, exists := sets[key]
		if !exists {
			return response.Value{Typ: "array", Array: []response.Value{}}
		}
		for member := range intersection {
			if _, found := curr[member]; !found {
				delete(intersection, member)
			}
		}
	}

	var result []response.Value
	for member := range intersection {
		result = append(result, response.Value{Typ: "bulk", Bulk: member})
	}
	return response.Value{Typ: "array", Array: result}
}

func sdiff(args []response.Value) response.Value {
	if len(args) < 1 {
		return response.Value{Typ: "error", Str: "ERR wrong number of arguments for 'SDIFF'"}
	}

	setsMu.RLock()
	defer setsMu.RUnlock()

	firstKey := args[0].Bulk
	base, exists := sets[firstKey]
	if !exists {
		return response.Value{Typ: "array", Array: []response.Value{}}
	}

	diff := make(map[string]struct{})
	for member := range base {
		diff[member] = struct{}{}
	}

	for _, arg := range args[1:] {
		key := arg.Bulk
		if s, exists := sets[key]; exists {
			for member := range s {
				delete(diff, member)
			}
		}
	}

	var result []response.Value
	for member := range diff {
		result = append(result, response.Value{Typ: "bulk", Bulk: member})
	}
	return response.Value{Typ: "array", Array: result}
}
