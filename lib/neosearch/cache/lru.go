package cache

import "container/list"

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

func (lru *LRUCache) Len() int {
	return lru.ll.Len()
}

// Add new interface{} value to LRUCache.
func (lru *LRUCache) Add(key string, value interface{}) {
	var (
		elem *list.Element
		ok   bool
	)

	if elem, ok = lru.cache[key]; ok {
		lru.removeElement(elem)
	}

	elem = lru.ll.PushFront(&entry{key, value})
	lru.cache[key] = elem

	if lru.ll.Len() > lru.max {
		lru.removeOldest()
	}
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
