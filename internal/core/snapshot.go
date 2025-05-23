package core

func CopyStrings() map[string]Entry {
	SDSsMu.RLock()
	defer SDSsMu.RUnlock()

	copy := make(map[string]Entry)
	for k, v := range SDSs {
		copy[k] = v
	}
	return copy
}

func CopyLists() map[string][]string {
	ListsMu.RLock()
	defer ListsMu.RUnlock()

	copy := make(map[string][]string)
	for k, v := range Lists {
		newList := make([]string, len(v))
		copy[k] = append(newList[:0], v...)
	}
	return copy
}

func CopyHSets() map[string]map[string]string {
	HSETsMu.RLock()
	defer HSETsMu.RUnlock()

	copy := make(map[string]map[string]string)
	for k, v := range HSETs {
		inner := make(map[string]string)
		for field, val := range v {
			inner[field] = val
		}
		copy[k] = inner
	}
	return copy
}

func CopySets() map[string]map[string]struct{} {
	SetsMu.RLock()
	defer SetsMu.RUnlock()

	copy := make(map[string]map[string]struct{})
	for k, v := range Sets {
		inner := make(map[string]struct{})
		for member := range v {
			inner[member] = struct{}{}
		}
		copy[k] = inner
	}
	return copy
}

func CopyZSets() map[string]*ZSet {
	ZSetsMu.RLock()
	defer ZSetsMu.RUnlock()

	copy := make(map[string]*ZSet)
	for k, v := range ZSets {
		copy[k] = v
	}
	return copy
}

func RestoreFromSnapshot(snapshot Snapshot) {
	// String
	SDSsMu.Lock()
	SDSs = snapshot.Strings
	SDSsMu.Unlock()

	// List
	ListsMu.Lock()
	Lists = snapshot.Lists
	ListsMu.Unlock()

	// Hash
	HSETsMu.Lock()
	HSETs = snapshot.Hsets
	HSETsMu.Unlock()

	// Set
	ZSetsMu.Lock()
	//ZSets = snapshot.ZSets
	ZSetsMu.Unlock()

	// ZSet
	ZSetsMu.Lock()
	//ZSets = snapshot.ZSets
	ZSetsMu.Unlock()

}
