package handle

import (
	"github.com/kenelite/gedis/internal/response"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Entry struct {
	Value     string
	ExpiresAt time.Time // zero time means no expiration
}

var (
	SETs   = map[string]Entry{}
	SETsMu = sync.RWMutex{}
)

func set(args []response.Value) response.Value {
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

	SETsMu.Lock()
	entry := Entry{Value: value}
	if ttl > 0 {
		entry.ExpiresAt = time.Now().Add(ttl)
	}
	SETs[key] = entry
	SETsMu.Unlock()

	return response.Value{Typ: "string", Str: "OK"}
}

func get(args []response.Value) response.Value {
	if len(args) != 1 {
		return response.Value{Typ: "error", Str: "ERR wrong number of arguments for 'get' command"}
	}

	key := args[0].Bulk

	SETsMu.RLock()
	entry, ok := SETs[key]
	SETsMu.RUnlock()

	if !ok || (entry.ExpiresAt != (time.Time{}) && time.Now().After(entry.ExpiresAt)) {
		// Optionally remove expired key
		SETsMu.Lock()
		delete(SETs, key)
		SETsMu.Unlock()
		return response.Value{Typ: "null"}
	}

	return response.Value{Typ: "bulk", Bulk: entry.Value}
}

func ttl(args []response.Value) response.Value {
	if len(args) != 1 {
		return response.Value{Typ: "error", Str: "ERR wrong number of arguments for 'ttl' command"}
	}

	key := args[0].Bulk

	SETsMu.RLock()
	entry, ok := SETs[key]
	SETsMu.RUnlock()

	if !ok {
		return response.Value{Typ: "integer", Num: -2} // Key doesn't exist
	}

	if entry.ExpiresAt.IsZero() {
		return response.Value{Typ: "integer", Num: -1} // No expiration set
	}

	remaining := int(time.Until(entry.ExpiresAt).Seconds())
	if remaining <= 0 {
		// Expired — clean up
		SETsMu.Lock()
		delete(SETs, key)
		SETsMu.Unlock()
		return response.Value{Typ: "integer", Num: -2}
	}

	return response.Value{Typ: "integer", Num: remaining}
}
