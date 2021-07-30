package lru

import (
	"fmt"
)

type DListNode struct {
	key   string
	value Value
	next  *DListNode
	prev  *DListNode
}

// Cache 是一个 LRU Cache，对并发访问并不安全
// 字典的键是字符串，值是双向链表中对应节点的指针
// capacity 是允许使用的最大内存，size 是当前已使用的内存，单位 byte
type Cache struct {
	cache    map[string]*DListNode
	head     *DListNode
	tail     *DListNode
	capacity int64
	size     int64
}

// Value 为了通用性，我们允许值是实现了 Value 接口的任意类型，该接口只包含了一个方法 Len() int，用于返回值所占用的内存大小
type Value interface {
	Len() int
}

// NewCache 为了方便实例化 LRUCache
func NewCache(capacity int64) *Cache {
	cache := Cache{
		capacity: capacity,
		head:     &DListNode{},
		tail:     &DListNode{},
		cache:    make(map[string]*DListNode),
		size:     0,
	}
	cache.head.next = cache.tail
	cache.tail.prev = cache.head
	return &cache
}

func newNode(key string, value Value, next *DListNode, prev *DListNode) *DListNode {
	return &DListNode{
		key:   key,
		value: value,
		next:  next,
		prev:  prev,
	}
}

func (c *Cache) Get(key string) Value {
	if node, ok := c.cache[key]; ok {
		c.makeRecently(node)
		return node.value
	} else {
		return nil
	}
}

func (c *Cache) Put(key string, value Value) {
	// 先判断是否有足够的空间
	putSize := int64(len(key)) + int64(value.Len())
	if c.capacity > 0 && c.capacity >= putSize {
		for c.size+putSize > c.capacity {
			// 删除最近最少使用的节点，并从哈希标中移除对应的项
			leastRecently := c.removeLeastRecently()
			delete(c.cache, leastRecently.key)
			// 释放相应的空间
			c.size -= int64(len(leastRecently.key)) + int64(leastRecently.value.Len())
		}
		// 判断该值是否在缓存中
		if node, exist := c.cache[key]; exist {
			node.value = value
			c.makeRecently(node)
			// 消耗对应的空间
			c.size += int64(value.Len()) - int64(node.value.Len())
		} else {
			// 添加新的节点
			node = newNode(key, value, nil, nil)
			c.addRecently(node)
			c.cache[key] = node
			c.size += putSize
		}
	} else {
		fmt.Println("error: no enough space")
	}
}

func (c *Cache) makeRecently(node *DListNode) {
	c.deleteNode(node)
	c.addRecently(node)
}

func (c *Cache) deleteNode(node *DListNode) {
	node.prev.next = node.next
	node.next.prev = node.prev
}

func (c *Cache) addRecently(node *DListNode) {
	node.prev = c.tail.prev
	node.next = c.tail
	c.tail.prev.next = node
	c.tail.prev = node
}

func (c *Cache) removeLeastRecently() *DListNode {
	node := c.head.next
	c.deleteNode(node)
	return node
}
