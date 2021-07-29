package lru

import (
	"testing"
)

func TestCache(t *testing.T) {
	LRUCache := newCache(2)
	LRUCache.put(1, 1)
	LRUCache.put(2, 2)
	res1 := LRUCache.get(1)
	LRUCache.put(3, 3)
	res2 := LRUCache.get(2)
	LRUCache.put(4, 4)
	res3 := LRUCache.get(1)
	LRUCache.get(3)
	LRUCache.get(4)
	switch {
	case res1 != 1:
		t.Error("the result of the 1st get() should be 1")
	case res2 != -1:
		t.Error("the result of the 2nd get() should be -1")
	case res3 != -1:
		t.Error("the result of the 3rd get() should be -1")
	}
}
