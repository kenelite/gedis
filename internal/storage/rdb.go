package storage

import (
	"encoding/json"
	"github.com/kenelite/gedis/internal/handle"
	"os"
	"time"
)

type SnapshotEntry struct {
	Type      string
	Key       string
	Value     any
	ExpiresAt int64 // Unix timestamp
}

func SaveRDB(filename string) error {
	var entries []SnapshotEntry

	handle.SETsMu.RLock()
	for k, v := range handle.SETs {
		entries = append(entries, SnapshotEntry{
			Type:      "string",
			Key:       k,
			Value:     v.Value,
			ExpiresAt: v.ExpiresAt.Unix(),
		})
	}
	handle.SETsMu.RUnlock()

	// Similarly dump lists, sets, zsets, hashes, etc.

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(entries)
}

func LoadRDB(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	var entries []SnapshotEntry
	if err := json.NewDecoder(file).Decode(&entries); err != nil {
		return err
	}

	now := time.Now()
	for _, entry := range entries {
		if entry.Type == "string" {
			expiresAt := time.Time{}
			if entry.ExpiresAt > 0 {
				expiresAt = time.Unix(entry.ExpiresAt, 0)
				if expiresAt.Before(now) {
					continue // skip expired
				}
			}
			handle.SETsMu.Lock()
			handle.SETs[entry.Key] = handle.Entry{
				Value:     entry.Value.(string),
				ExpiresAt: expiresAt,
			}
			handle.SETsMu.Unlock()
		}

		// Handle lists, sets, etc similarly...
	}
	return nil
}
