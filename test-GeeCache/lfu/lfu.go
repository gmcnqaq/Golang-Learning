package lfu

import "Golang/Learning/test-GeeCache/lru"

type Cache struct {
	capacity   int64
	size       int64
	minFreq    int64
	keyToVal   map[string]lru.Value
	keyToFreq  map[string]int
	freqToKeys map[int][]string
}

func (c *Cache) Get(key string) lru.Value {
	if value, ok := c.keyToVal[key]; ok {
		return value
	} else {
		return nil
	}
}

func (c *Cache) Put(key string) {

}
