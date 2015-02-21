package engine

import "bitbucket.org/i4k/neosearch/store"

// StoreEntry have the KVStore and LastAccess of the index
// to implement a simple LRU cache.
// The length of this cache is set by NGConfig.OpenCacheSize.
// Default value is 100.
type StoreEntry struct {
	Store      store.KVStore
	Name       string
	LastAccess int
}

// StoreCache is a hash of the StoreEntry's ids
// At every engine.Open call, a goroutine will
// check the length of the cache and close the least used
// database connections
type StoreCache map[string]*StoreEntry

// ByLastAccess sorts by access time
type ByLastAccess []StoreEntry

func (by ByLastAccess) Len() int {
	return len(by)
}

func (by ByLastAccess) Swap(i, j int) {
	by[i], by[j] = by[j], by[i]
}

func (by ByLastAccess) Less(i, j int) bool {
	return by[i].LastAccess < by[j].LastAccess
}
