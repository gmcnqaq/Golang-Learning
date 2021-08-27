package lru

import (
	"fmt"
	"testing"
)

type String string

func (s String) Len() int {
	return len(s)
}

func TestCache(t *testing.T) {
	k1, k2, k3, k4 := "key1", "key2", "key3", "key4key4"
	v1, v2, v3, v4 := "value1", "value2", "value3", "value4value4"
	capacity := len(k1 + k2 + v1 + v2)
	fmt.Println("LRU cache capacity:", capacity)
	LRUCache := NewCache(int64(capacity))
	LRUCache.Put(k1, String(v1))
	LRUCache.Put(k2, String(v2))
	res1 := LRUCache.Get(k1)
	LRUCache.Put(k3, String(v3))
	res2 := LRUCache.Get(k2)
	fmt.Println(*LRUCache)
	LRUCache.Put(k4, String(v4))
	res3 := LRUCache.Get(k1)
	res4 := LRUCache.Get(k3)
	LRUCache.Get(k4)
	fmt.Println(*LRUCache)
	switch {
	case res1 != String(v1):
		t.Error("cache hit k1=v1 failed")
	case res2 != nil:
		fmt.Println(res2)
		t.Error("cache miss k2 failed")
	case res3 != nil:
		fmt.Println(res3)
		t.Error("cache miss k1 failed")
	case res4 != nil:
		fmt.Println(res4)
		t.Error("cache miss k3 failed")
	}
}
