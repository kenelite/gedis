package handle

import (
	. "github.com/kenelite/gedis/internal/core"
	"github.com/kenelite/gedis/internal/response"
	"strconv"
	"time"
)

func Expire(args []response.Value) response.Value {
	if len(args) != 2 {
		return response.Value{Typ: "error", Str: "ERR wrong number of arguments for 'expire' command"}
	}

	key := args[0].Bulk
	seconds, err := strconv.Atoi(args[1].Bulk)
	if err != nil {
		return response.Value{Typ: "error", Str: "ERR invalid expire time"}
	}

	SDSsMu.Lock()
	defer SDSsMu.Unlock()

	entry, exists := SDSs[key]
	if !exists {
		return response.Value{Typ: "integer", Num: 0}
	}

	entry.ExpiresAt = time.Now().Add(time.Duration(seconds) * time.Second)
	SDSs[key] = entry

	return response.Value{Typ: "integer", Num: 1}
}

func StartKeyExpirationLoop() {
	go func() {
		for {
			time.Sleep(time.Second)
			now := time.Now()

			SDSsMu.Lock()
			for k, v := range SDSs {
				if !v.ExpiresAt.IsZero() && now.After(v.ExpiresAt) {
					delete(SDSs, k)
				}
			}
			SDSsMu.Unlock()
		}
	}()
}
