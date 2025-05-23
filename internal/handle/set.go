package handle

import (
	. "github.com/kenelite/gedis/internal/core"
	"github.com/kenelite/gedis/internal/response"
)

func Sadd(args []response.Value) response.Value {
	if len(args) < 2 {
		return response.Value{Typ: "error", Str: "ERR wrong number of arguments for 'SADD'"}
	}

	key := args[0].Bulk
	SetsMu.Lock()
	defer SetsMu.Unlock()

	if Sets[key] == nil {
		Sets[key] = make(map[string]struct{})
	}

	added := 0
	for _, arg := range args[1:] {
		val := arg.Bulk
		if _, exists := Sets[key][val]; !exists {
			Sets[key][val] = struct{}{}
			added++
		}
	}

	return response.Value{Typ: "integer", Num: added}
}

func Srem(args []response.Value) response.Value {
	if len(args) < 2 {
		return response.Value{Typ: "error", Str: "ERR wrong number of arguments for 'SREM'"}
	}

	key := args[0].Bulk
	SetsMu.Lock()
	defer SetsMu.Unlock()

	members, exists := Sets[key]
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
		delete(Sets, key)
	}

	return response.Value{Typ: "integer", Num: removed}
}

func Smembers(args []response.Value) response.Value {
	if len(args) != 1 {
		return response.Value{Typ: "error", Str: "ERR wrong number of arguments for 'SMEMBERS'"}
	}

	key := args[0].Bulk
	SetsMu.RLock()
	defer SetsMu.RUnlock()

	members, exists := Sets[key]
	if !exists {
		return response.Value{Typ: "array", Array: []response.Value{}}
	}

	var result []response.Value
	for member := range members {
		result = append(result, response.Value{Typ: "bulk", Bulk: member})
	}

	return response.Value{Typ: "array", Array: result}
}

func Scard(args []response.Value) response.Value {
	if len(args) != 1 {
		return response.Value{Typ: "error", Str: "ERR wrong number of arguments for 'SCARD'"}
	}

	key := args[0].Bulk
	SetsMu.RLock()
	defer SetsMu.RUnlock()

	members, exists := Sets[key]
	if !exists {
		return response.Value{Typ: "integer", Num: 0}
	}

	return response.Value{Typ: "integer", Num: len(members)}
}

func Sunion(args []response.Value) response.Value {
	if len(args) < 1 {
		return response.Value{Typ: "error", Str: "ERR wrong number of arguments for 'SUNION'"}
	}

	SetsMu.RLock()
	defer SetsMu.RUnlock()

	union := make(map[string]struct{})
	for _, arg := range args {
		key := arg.Bulk
		if members, exists := Sets[key]; exists {
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

func Sinter(args []response.Value) response.Value {
	if len(args) < 1 {
		return response.Value{Typ: "error", Str: "ERR wrong number of arguments for 'SINTER'"}
	}

	SetsMu.RLock()
	defer SetsMu.RUnlock()

	// Initialize result with first set
	firstKey := args[0].Bulk
	base, exists := Sets[firstKey]
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
		curr, exists := Sets[key]
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

func Sdiff(args []response.Value) response.Value {
	if len(args) < 1 {
		return response.Value{Typ: "error", Str: "ERR wrong number of arguments for 'SDIFF'"}
	}

	SetsMu.RLock()
	defer SetsMu.RUnlock()

	firstKey := args[0].Bulk
	base, exists := Sets[firstKey]
	if !exists {
		return response.Value{Typ: "array", Array: []response.Value{}}
	}

	diff := make(map[string]struct{})
	for member := range base {
		diff[member] = struct{}{}
	}

	for _, arg := range args[1:] {
		key := arg.Bulk
		if s, exists := Sets[key]; exists {
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
