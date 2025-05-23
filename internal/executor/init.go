package executor

import (
	"github.com/kenelite/gedis/internal/handle"
)

func init() {
	// String
	Register("PING", handle.Ping)
	Register("SET", handle.Set)
	Register("GET", handle.Get)
	Register("DEL", handle.Del)
	Register("INCR", handle.Incr)
	Register("DECR", handle.Decr)
	Register("INCRBY", handle.Incrby)
	Register("DECRBY", handle.Decrby)

	Register("EXPIRE", handle.Expire)
	Register("TTL", handle.Ttl)

	// List
	Register("LPUSH", handle.Lpush)
	Register("RPUSH", handle.Rpush)
	Register("LPOP", handle.Lpop)
	Register("RPOP", handle.Rpop)
	Register("LRANGE", handle.Lrange)

	// Set
	Register("SADD", handle.Sadd)
	Register("SREM", handle.Srem)
	Register("SMEMBERS", handle.Smembers)
	Register("SCARD", handle.Scard)
	Register("SINTER", handle.Sinter)
	Register("SUNION", handle.Sunion)
	Register("SDIFF", handle.Sdiff)

	// Hash
	Register("HSET", handle.Hset)
	Register("HGET", handle.Hget)
	Register("HGETALL", handle.Hgetall)

	// Zset
	Register("ZADD", handle.Zadd)
	Register("ZREM", handle.Zrem)
	Register("ZRANGE", handle.Zrange)
	Register("ZCARD", handle.Zcard)

	// Storage
	Register("SAVE", handle.Save)
	Register("BGSAVE", handle.BgSave)
}
