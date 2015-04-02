package store

import "container/list"

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

type LRUCache struct {
	max int

	onRemove OnRemoveCb

	ll    *list.List
	cache map[string]*list.Element
}

type entry struct {
	key   string
	value interface{}
}

// NewCache returns a new LRU cache for interface{} entries
func NewLRUCache(max int) *LRUCache {
	if max <= 0 {
		return nil
	}

	return &LRUCache{
		max:   max,
		ll:    list.New(),
		cache: make(map[string]*list.Element),
	}
}

func (lru *LRUCache) OnRemove(cb OnRemoveCb) {
	lru.onRemove = cb
}

// MaxEntries update the max allowed entries in cache
func (lru *LRUCache) MaxEntries(max int) {
	lru.max = max
}

// Add new interface{} to LRUCache. If already exists an entry with
// key `key` returns `false` and nothing happens. If the entry
// not exists in cache, then will be added and returns `true`.
// If you only need to rank the given key to top, use Get(key).
func (lru *LRUCache) Add(key string, value interface{}) bool {
	var (
		elem *list.Element
		ok   bool
	)

	if elem, ok = lru.cache[key]; ok {
		// we can't simply use ll.MoveToFront and
		// elem.(*entry).value = value
		// because this way we lost the old ref of
		// interface{}, GC will deallocate the pointer leaving
		// the database locked.
		// Then, be sure of always try lru.Get() and then
		// lru.Add()

		return false
	}

	elem = lru.ll.PushFront(&entry{key, value})
	lru.cache[key] = elem

	if lru.ll.Len() > lru.max {
		lru.removeOldest()
	}

	return true
}

// Get the given `key` from cache. If the key exists, it will be ranked
// to top of the cache.
func (lru *LRUCache) Get(key string) (interface{}, bool) {
	var (
		elem *list.Element
		ok   bool
	)

	if lru.cache == nil || len(lru.cache) == 0 {
		return nil, false
	}

	elem, ok = lru.cache[key]

	if ok {
		lru.ll.MoveToFront(elem)
		return elem.Value.(*entry).value, true
	}

	return nil, false
}

func (lru *LRUCache) Remove(key string) bool {
	var (
		elem *list.Element
		ok   bool
	)

	if lru.cache == nil || len(lru.cache) == 0 {
		return false
	}

	elem, ok = lru.cache[key]

	if ok {
		lru.removeElement(elem)
		return true
	}

	return false
}

func (lru *LRUCache) removeOldest() {
	var elem *list.Element

	if lru.cache == nil || len(lru.cache) == 0 {
		return
	}

	elem = lru.ll.Back()
	if elem != nil {
		lru.removeElement(elem)
	}
}

func (lru *LRUCache) removeElement(elem *list.Element) {
	lru.ll.Remove(elem)

	kv := elem.Value.(*entry)
	delete(lru.cache, kv.key)

	if lru.onRemove != nil {
		lru.onRemove(kv.key, kv.value)
	}
}

// Clean remove all elements of cache calling the OnRemove callback
// when needed!
func (lru *LRUCache) Clean() {
	var (
		cacheLen int = len(lru.cache)
		elem     *list.Element
		key      string
		value    interface{}
	)

	if lru.cache == nil || cacheLen == 0 {
		return
	}

	for elem = lru.ll.Front(); elem != nil; elem = elem.Next() {
		ee := elem.Value.(*entry)
		key = ee.key
		value = ee.value

		if lru.onRemove != nil {
			lru.onRemove(key, value)
		}

		lru.removeElement(elem)
	}
}
