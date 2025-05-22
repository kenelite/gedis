package handle

import (
	"github.com/kenelite/gedis/internal/response"
)

//var Handlers = map[string]func([]response.Value) response.Value{
//	"PING": ping,
//
//	//kv
//	"SET":    set,
//	"GET":    get,
//	"EXPIRE": expire,
//	"TTL":    ttl,
//
//	//kv
//	"DEL":    del,
//	"INCR":   incr,
//	"DECR":   decr,
//	"INCRBY": incrby,
//	"DECRBY": decrby,
//
//	//hash
//	"HSET":    hset,
//	"HGET":    hget,
//	"HGETALL": hgetall,
//
//	//list
//	"LPUSH":  lpush,
//	"RPUSH":  rpush,
//	"LPOP":   lpop,
//	"RPOP":   rpop,
//	"LRANGE": lrange,
//
//	//set
//	"SADD":     sadd,
//	"SREM":     srem,
//	"SMEMBERS": smembers,
//	"SCARD":    scard,
//	"SUNION":   sunion,
//	"SINTER":   sinter,
//	"SDIFF":    sdiff,
//
//	//zset
//	"ZADD":   zadd,
//	"ZRANGE": zrange,
//	"ZREM":   zrem,
//
//	//save
//	"SAVE": saveRdb,
//}

func Ping(args []response.Value) response.Value {
	if len(args) == 0 {
		return response.Value{Typ: "string", Str: "PONG"}
	}

	return response.Value{Typ: "string", Str: args[0].Bulk}
}
