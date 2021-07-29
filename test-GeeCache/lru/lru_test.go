package lru

import (
	"testing"
)

func TestCache(t *testing.T) {
	LRUCache := NewCache(2)
	LRUCache.Put(1, 1)
	LRUCache.Put(2, 2)
	res1 := LRUCache.Get(1)
	LRUCache.Put(3, 3)
	res2 := LRUCache.Get(2)
	LRUCache.Put(4, 4)
	res3 := LRUCache.Get(1)
	LRUCache.Get(3)
	LRUCache.Get(4)
	switch {
	case res1 != 1:
		t.Error("the result of the 1st get() should be 1")
	case res2 != -1:
		t.Error("the result of the 2nd get() should be -1")
	case res3 != -1:
		t.Error("the result of the 3rd get() should be -1")
	}
}
