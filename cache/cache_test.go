package cache

import "testing"

func TestLRUCreate(t *testing.T) {
	lru := NewLRUCache(1)

	if lru == nil {
		t.Error("Failed to create LRUCache")
	}

	lru = NewLRUCache(0)

	if lru != nil {
		t.Error("LRUCache should fails for zero length")
	}

	lru = NewLRUCache(-1)

	if lru != nil {
		t.Error("LRUCache should fails for zero length")
	}
}

func TestLRUSmall(t *testing.T) {
	lru := NewLRUCache(1)

	lru.Add("i4k", 1)
}
