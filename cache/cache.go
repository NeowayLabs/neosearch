package cache

type OnRemoveCb func(key string, value interface{})

type Cache interface {
	// Add new entry to cache
	Add(key string, value interface{}) bool

	// Get the entry with `key`
	Get(key string) (interface{}, bool)

	// Remove `key` from cache
	Remove(key string) bool

	// OnRemove set a callback to be executed when entries are removed
	// from cache
	OnRemove(cb OnRemoveCb)

	// MaxEntries set or update the max entries of cache
	MaxEntries(max int)

	Clean()
}
