package core

import (
	"sync"
	"time"
)

type Entry struct {
	Value     string
	ExpiresAt time.Time // zero time means no expiration
}

type ZSetEntry struct {
	Member string
	Score  float64
}

type ZSet struct {
	Entries map[string]float64
}

type Snapshot struct {
	Strings map[string]Entry
	Lists   map[string][]string
	ZSets   map[string]ZSet
	Sets    map[string]map[string]struct{}
}

var (
	Lists   = make(map[string][]string)
	ListsMu sync.RWMutex

	Sets   = make(map[string]map[string]struct{})
	SetsMu sync.RWMutex

	HSETs   = map[string]map[string]string{}
	HSETsMu = sync.RWMutex{}

	SETs   = map[string]Entry{}
	SETsMu = sync.RWMutex{}

	ZSets   = make(map[string]*ZSet)
	ZSetsMu sync.RWMutex
)
