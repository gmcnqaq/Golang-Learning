package lfu

import (
	"Golang/Learning/geecache/test-GeeCache/lru"
)

type Cache struct {
	capacity   int64
	size       int64
	minFreq    int
	keyToVal   map[string]lru.Value
	keyToFreq  map[string]int
	freqToKeys map[int][]string
}

func (c *Cache) Get(key string) lru.Value {
	if value, ok := c.keyToVal[key]; ok {
		c.increaseKey(key)
		return value
	} else {
		return nil
	}
}

func (c *Cache) Put(key string, value lru.Value) {
	// 先判断是否有足够的空间
	putSize := int64(len(key)) + int64(value.Len())
	if c.capacity > putSize {
		for c.size+putSize > c.capacity {
			// 删除频率最低，并从哈希表中移除对应的项
			c.removeMinFreq()
			// 释放对应的空间
		}
	}
}

func (c *Cache) increaseKey(key string) {
	// 更新 KF 表
	//freq := c.keyToFreq[key]
}

func (c *Cache) removeMinFreq() {
	//keyList := c.freqToKeys[c.minFreq]

}
