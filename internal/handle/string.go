package handle

import (
	. "github.com/kenelite/gedis/internal/core"
	"github.com/kenelite/gedis/internal/response"
	"strconv"
	"strings"
	"time"
)

func Set(args []response.Value) response.Value {
	if len(args) < 2 {
		return response.Value{Typ: "error", Str: "ERR wrong number of arguments for 'set' command"}
	}

	key := args[0].Bulk
	value := args[1].Bulk

	var ttl time.Duration

	if len(args) >= 4 && strings.ToUpper(args[2].Bulk) == "EX" {
		seconds, err := strconv.Atoi(args[3].Bulk)
		if err != nil {
			return response.Value{Typ: "error", Str: "ERR invalid TTL value"}
		}
		ttl = time.Duration(seconds) * time.Second
	}

	SDSsMu.Lock()
	entry := Entry{Value: value}
	if ttl > 0 {
		entry.ExpiresAt = time.Now().Add(ttl)
	}
	SDSs[key] = entry
	SDSsMu.Unlock()

	return response.Value{Typ: "string", Str: "OK"}
}

func Get(args []response.Value) response.Value {
	if len(args) != 1 {
		return response.Value{Typ: "error", Str: "ERR wrong number of arguments for 'get' command"}
	}

	key := args[0].Bulk

	SDSsMu.RLock()
	entry, ok := SDSs[key]
	SDSsMu.RUnlock()

	if !ok || (entry.ExpiresAt != (time.Time{}) && time.Now().After(entry.ExpiresAt)) {
		// Optionally remove expired key
		SDSsMu.Lock()
		delete(SDSs, key)
		SDSsMu.Unlock()
		return response.Value{Typ: "null"}
	}

	return response.Value{Typ: "bulk", Bulk: entry.Value}
}

func Ttl(args []response.Value) response.Value {
	if len(args) != 1 {
		return response.Value{Typ: "error", Str: "ERR wrong number of arguments for 'ttl' command"}
	}

	key := args[0].Bulk

	SDSsMu.RLock()
	entry, ok := SDSs[key]
	SDSsMu.RUnlock()

	if !ok {
		return response.Value{Typ: "integer", Num: -2} // Key doesn't exist
	}

	if entry.ExpiresAt.IsZero() {
		return response.Value{Typ: "integer", Num: -1} // No expiration set
	}

	remaining := int(time.Until(entry.ExpiresAt).Seconds())
	if remaining <= 0 {
		// Expired — clean up
		SDSsMu.Lock()
		delete(SDSs, key)
		SDSsMu.Unlock()
		return response.Value{Typ: "integer", Num: -2}
	}

	return response.Value{Typ: "integer", Num: remaining}
}

func Del(args []response.Value) response.Value {
	if len(args) < 1 {
		return response.Value{Typ: "error", Str: "ERR wrong number of arguments for 'DEL'"}
	}

	deleted := 0
	SDSsMu.Lock()
	defer SDSsMu.Unlock()

	for _, arg := range args {
		key := arg.Bulk
		if _, exists := SDSs[key]; exists {
			delete(SDSs, key)
			deleted++
		}
	}

	return response.Value{Typ: "integer", Num: deleted}
}

func Incr(args []response.Value) response.Value {
	if len(args) != 1 {
		return response.Value{Typ: "error", Str: "ERR wrong number of arguments for 'INCR'"}
	}
	return incrbyImpl(args[0].Bulk, 1)
}

func Decr(args []response.Value) response.Value {
	if len(args) != 1 {
		return response.Value{Typ: "error", Str: "ERR wrong number of arguments for 'DECR'"}
	}
	return incrbyImpl(args[0].Bulk, -1)
}

func Incrby(args []response.Value) response.Value {
	if len(args) != 2 {
		return response.Value{Typ: "error", Str: "ERR wrong number of arguments for 'INCRBY'"}
	}

	amount, err := strconv.Atoi(args[1].Bulk)
	if err != nil {
		return response.Value{Typ: "error", Str: "ERR value is not an integer or out of range"}
	}

	return incrbyImpl(args[0].Bulk, amount)
}

func Decrby(args []response.Value) response.Value {
	if len(args) != 2 {
		return response.Value{Typ: "error", Str: "ERR wrong number of arguments for 'DECRBY'"}
	}

	amount, err := strconv.Atoi(args[1].Bulk)
	if err != nil {
		return response.Value{Typ: "error", Str: "ERR value is not an integer or out of range"}
	}

	return incrbyImpl(args[0].Bulk, -amount)
}

func incrbyImpl(key string, delta int) response.Value {
	SDSsMu.Lock()
	defer SDSsMu.Unlock()

	entry, exists := SDSs[key]
	if exists && !entry.ExpiresAt.IsZero() && time.Now().After(entry.ExpiresAt) {
		delete(SDSs, key)
		exists = false
	}

	var val int
	if exists {
		var err error
		val, err = strconv.Atoi(entry.Value)
		if err != nil {
			return response.Value{Typ: "error", Str: "ERR value is not an integer or out of range"}
		}
	}

	val += delta
	SDSs[key] = Entry{
		Value:     strconv.Itoa(val),
		ExpiresAt: entry.ExpiresAt,
	}

	return response.Value{Typ: "integer", Num: val}
}
