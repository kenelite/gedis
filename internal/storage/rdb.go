package storage

import (
	"encoding/gob"
	"github.com/kenelite/gedis/internal/core"
	"os"
	"sync"
)

var rdbMu sync.Mutex

func SaveRDB(path string) error {
	rdbMu.Lock()
	defer rdbMu.Unlock()

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	snapshot := core.Snapshot{
		Strings: core.CopyStrings(),
		Lists:   core.CopyLists(),
		Hsets:   core.CopyHSets(),
		ZSets:   core.CopyZSets(),
		Sets:    core.CopySets(),
	}

	enc := gob.NewEncoder(f)
	return enc.Encode(snapshot)
}

func LoadRDB(path string) error {
	rdbMu.Lock()
	defer rdbMu.Unlock()

	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	var snapshot core.Snapshot
	dec := gob.NewDecoder(f)
	if err := dec.Decode(&snapshot); err != nil {
		return err
	}

	core.RestoreFromSnapshot(snapshot)
	return nil
}
