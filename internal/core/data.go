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
	Sets    map[string]map[string]struct{}
	Hsets   map[string]map[string]string
	ZSets   map[string]*ZSet
}

var (

	// String
	SDSs   = map[string]Entry{}
	SDSsMu = sync.RWMutex{}

	// List
	Lists   = make(map[string][]string)
	ListsMu sync.RWMutex

	// Set
	Sets   = make(map[string]map[string]struct{})
	SetsMu sync.RWMutex

	// Hash
	HSETs   = map[string]map[string]string{}
	HSETsMu = sync.RWMutex{}

	// Zset
	ZSets   = make(map[string]*ZSet)
	ZSetsMu sync.RWMutex
)
