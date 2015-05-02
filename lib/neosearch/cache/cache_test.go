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

func TestLRUAllowAddInfiniValues(t *testing.T) {
	var (
		vi    interface{}
		value int
		ok    bool
	)

	lru := NewLRUCache(1)

	if lru.Len() != 0 {
		t.Error("OMG... Very piece of buggy code")
	}

	if vi, ok = lru.Get("some thing"); ok == true || vi != nil {
		t.Error("... i dont beliece that ...")
	}

	lru.Add("a", 10)

	if vi, ok = lru.Get("a"); !ok || vi == nil {
		t.Error("Failed to get entry")
	} else {
		value = vi.(int)

		if value != 10 {
			t.Error("Failed to get entry")
		}
	}

	if lru.Len() != 1 {
		t.Error("Invalid length")
	}

	lru.Add("b", 11)

	if vi, ok = lru.Get("b"); !ok || vi == nil {
		t.Error("Failed to get entry")
	} else {
		value = vi.(int)

		if value != 11 {
			t.Error("Failed to get entry")
		}
	}

	if vi, ok = lru.Get("a"); ok || vi != nil {
		t.Error("entry shouldnt exists")
	}

	if lru.Len() != 1 {
		t.Error("Invalid Length")
	}

	lru.Add("c", 12)

	if vi, ok = lru.Get("c"); !ok || vi == nil {
		t.Error("Failed to get entry")
	} else {
		value = vi.(int)

		if value != 12 {
			t.Error("Failed to get entry")
		}
	}

	if vi, ok = lru.Get("a"); ok || vi != nil {
		t.Error("entry shouldnt exists")
	}

	if vi, ok = lru.Get("b"); ok || vi != nil {
		t.Error("entry shouldnt exists")
	}

	if lru.Len() != 1 {
		t.Error("Invalid Length")
	}

	lru = nil

	lru = NewLRUCache(2)

	lru.Add("a", 1)

	if lru.Len() != 1 {
		t.Error("Invalid Length")
	}

	lru.Add("b", 2)

	if lru.Len() != 2 {
		t.Error("Invalid Length")
	}

	lru.Add("c", 3)

	if lru.Len() != 2 {
		t.Error("Invalid Length")
	}
}

func TestLRUOnRemoveCallback(t *testing.T) {
	lru := NewLRUCache(2)
	works := false
	removedKey := ""
	removedVal := 0

	lru.OnRemove(func(key string, value interface{}) {
		v, _ := value.(int)

		removedKey = key
		removedVal = v
		works = true
	})

	lru.Add("teste", 1)
	lru.Add("teste2", 2)
	lru.Add("teste3", 3)

	if works != true || removedKey != "teste" || removedVal != 1 {
		t.Error("OnRemove callback not invoked OR called concurrently")
	}

	works = false
	lru.Add("teste4", 4)

	if works != true || removedKey != "teste2" || removedVal != 2 {
		t.Error("OnRemove callback not invoked OR called concurrently")
	}
}
