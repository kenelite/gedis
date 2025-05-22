package handle

import (
	"github.com/kenelite/gedis/internal/response"
	"sort"
	"strconv"
	"strings"
	"sync"
)

type ZSetEntry struct {
	Member string
	Score  float64
}

type ZSet struct {
	entries map[string]float64
}

var (
	ZSets   = make(map[string]*ZSet)
	ZSetsMu sync.RWMutex
)

func Zadd(args []response.Value) response.Value {
	if len(args) < 3 || len(args)%2 != 1 {
		return response.Value{Typ: "error", Str: "ERR syntax: ZADD key score member [score member ...]"}
	}

	key := args[0].Bulk

	ZSetsMu.Lock()
	defer ZSetsMu.Unlock()

	zset, exists := ZSets[key]
	if !exists {
		zset = &ZSet{entries: make(map[string]float64)}
		ZSets[key] = zset
	}

	added := 0
	for i := 1; i < len(args); i += 2 {
		score, err := strconv.ParseFloat(args[i].Bulk, 64)
		if err != nil {
			return response.Value{Typ: "error", Str: "ERR invalid score"}
		}
		member := args[i+1].Bulk
		if _, exists := zset.entries[member]; !exists {
			added++
		}
		zset.entries[member] = score
	}

	return response.Value{Typ: "integer", Num: added}
}

func Zrange(args []response.Value) response.Value {
	if len(args) < 3 {
		return response.Value{Typ: "error", Str: "ERR wrong number of arguments for 'zrange'"}
	}

	key := args[0].Bulk
	start, err1 := strconv.Atoi(args[1].Bulk)
	stop, err2 := strconv.Atoi(args[2].Bulk)

	if err1 != nil || err2 != nil {
		return response.Value{Typ: "error", Str: "ERR invalid range"}
	}

	withScores := len(args) > 3 && strings.ToUpper(args[3].Bulk) == "WITHSCORES"

	ZSetsMu.RLock()
	defer ZSetsMu.RUnlock()

	zset, exists := ZSets[key]
	if !exists {
		return response.Value{Typ: "array", Array: []response.Value{}}
	}

	sorted := make([]ZSetEntry, 0, len(zset.entries))
	for m, s := range zset.entries {
		sorted = append(sorted, ZSetEntry{Member: m, Score: s})
	}
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Score < sorted[j].Score
	})

	if start < 0 {
		start += len(sorted)
	}
	if stop < 0 {
		stop += len(sorted)
	}
	if start < 0 {
		start = 0
	}
	if stop >= len(sorted) {
		stop = len(sorted) - 1
	}
	if start > stop || start >= len(sorted) {
		return response.Value{Typ: "array", Array: []response.Value{}}
	}

	result := []response.Value{}
	for _, entry := range sorted[start : stop+1] {
		result = append(result, response.Value{Typ: "bulk", Bulk: entry.Member})
		if withScores {
			result = append(result, response.Value{Typ: "bulk", Bulk: strconv.FormatFloat(entry.Score, 'f', -1, 64)})
		}
	}

	return response.Value{Typ: "array", Array: result}
}

func Zrem(args []response.Value) response.Value {
	if len(args) < 2 {
		return response.Value{Typ: "error", Str: "ERR wrong number of arguments for 'zrem'"}
	}

	key := args[0].Bulk
	ZSetsMu.Lock()
	defer ZSetsMu.Unlock()

	zset, exists := ZSets[key]
	if !exists {
		return response.Value{Typ: "integer", Num: 0}
	}

	removed := 0
	for _, arg := range args[1:] {
		if _, exists := zset.entries[arg.Bulk]; exists {
			delete(zset.entries, arg.Bulk)
			removed++
		}
	}

	return response.Value{Typ: "integer", Num: removed}
}

func Zcard(args []response.Value) response.Value {
	if len(args) != 1 {
		return response.Value{Typ: "error", Str: "ERR wrong number of arguments for 'zcard'"}
	}

	key := args[0].Bulk
	ZSetsMu.RLock()
	defer ZSetsMu.RUnlock()

	zset, exists := ZSets[key]
	if !exists {
		return response.Value{Typ: "integer", Num: 0}
	}

	return response.Value{Typ: "integer", Num: len(zset.entries)}
}
