package handle

import (
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

	SETsMu.Lock()
	defer SETsMu.Unlock()

	entry, exists := SETs[key]
	if !exists {
		return response.Value{Typ: "integer", Num: 0}
	}

	entry.ExpiresAt = time.Now().Add(time.Duration(seconds) * time.Second)
	SETs[key] = entry

	return response.Value{Typ: "integer", Num: 1}
}

func StartKeyExpirationLoop() {
	go func() {
		for {
			time.Sleep(time.Second)
			now := time.Now()

			SETsMu.Lock()
			for k, v := range SETs {
				if !v.ExpiresAt.IsZero() && now.After(v.ExpiresAt) {
					delete(SETs, k)
				}
			}
			SETsMu.Unlock()
		}
	}()
}
