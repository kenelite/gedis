package core

func CopyStrings() map[string]Entry {
	SETsMu.RLock()
	defer SETsMu.RUnlock()

	copy := make(map[string]Entry)
	for k, v := range SETs {
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

func CopyZSets() map[string]*ZSet {
	ZSetsMu.RLock()
	defer ZSetsMu.RUnlock()

	copy := make(map[string]*ZSet)
	for k, v := range ZSets {
		copy[k] = v
	}
	return copy
}

//func CopySets() map[string]map[string]struct{} {
//	SETsMu.RLock()
//	defer SETsMu.RUnlock()
//
//	copy := make(map[string]map[string]struct{})
//	for k, v := range SETs {
//		inner := make(map[string]struct{})
//		for member := range v {
//			inner[member] = struct{}{}
//		}
//		copy[k] = inner
//	}
//	return copy
//}

func RestoreFromSnapshot(snapshot Snapshot) {
	SETsMu.Lock()
	SETs = snapshot.Strings
	SETsMu.Unlock()

	ListsMu.Lock()
	Lists = snapshot.Lists
	ListsMu.Unlock()

	ZSetsMu.Lock()
	//ZSets = snapshot.ZSets
	ZSetsMu.Unlock()

}
