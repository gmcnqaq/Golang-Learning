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

func (c *cache) put(key string, value ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		c.lru = lru.NewCache(c.cacheBytes)
	}
	c.lru.Put(key, value)
}

func (c *cache) get(key string) (ByteView, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		return *NewByteView(nil), false
	}
	if val := c.lru.Get(key); val != nil {
		return val.(ByteView), true
	}
	return *NewByteView(nil), false
}
