package test_GeeCache

import (
	"Golang/Learning/test-GeeCache/lru"
	"sync"
)

// 为 LRU Cache 添加并发特性

type cache struct {
	mu         sync.Mutex
	lru        *lru.Cache
	cacheBytes int64
}

func (c *cache) put(key int, value int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		c.lru = lru.NewCache(c.cacheBytes)
	}
	c.lru.Put(key, value)
}

func (c *cache) get(key int) int {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		return -1
	}
	if val := c.lru.Get(key); val != -1 {
		return val
	}
	return -1
}
