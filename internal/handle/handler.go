package handle

import (
	"github.com/kenelite/gedis/internal/response"
)

var Handlers = map[string]func([]response.Value) response.Value{
	"PING":    ping,
	"SET":     set,
	"GET":     get,
	"HSET":    hset,
	"HGET":    hget,
	"HGETALL": hgetall,
	"EXPIRE":  expire,
	"TTL":     ttl,
	"LPUSH":   lpush,
	"RPUSH":   rpush,
	"LPOP":    lpop,
	"RPOP":    rpop,
	"LRANGE":  lrange,
	"ZADD":    zadd,
	"ZRANGE":  zrange,
	"ZREM":    zrem,

	"SADD":     sadd,
	"SREM":     srem,
	"SMEMBERS": smembers,
	"SCARD":    scard,
}

func ping(args []response.Value) response.Value {
	if len(args) == 0 {
		return response.Value{Typ: "string", Str: "PONG"}
	}

	return response.Value{Typ: "string", Str: args[0].Bulk}
}
