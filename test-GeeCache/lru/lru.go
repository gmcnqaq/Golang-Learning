package lru

type DListNode struct {
	key   int
	value int
	next  *DListNode
	prev  *DListNode
}

/*// 只包含一个 len() 方法返回值所占用的内存大小。通用性
type Value interface {
	len() int
}*/

type Cache struct {
	cache    map[int]*DListNode
	head     *DListNode
	tail     *DListNode
	capacity int64
	size     int64
}

// NewCache 为了方便实例化 LRUCache
func NewCache(capacity int64) *Cache {
	cache := Cache{
		capacity: capacity,
		head:     &DListNode{},
		tail:     &DListNode{},
		cache:    make(map[int]*DListNode),
		size:     0,
	}
	cache.head.next = cache.tail
	cache.tail.prev = cache.head
	return &cache
}

func newNode(key int, value int, next *DListNode, prev *DListNode) *DListNode {
	return &DListNode{
		key:   key,
		value: value,
		next:  next,
		prev:  prev,
	}
}

func (c *Cache) Get(key int) int {
	if node, ok := c.cache[key]; ok {
		c.makeRecently(node)
		return node.value
	} else {
		return -1
	}
}

func (c *Cache) Put(key int, value int) {
	if node, ok := c.cache[key]; ok {
		node.value = value
		c.makeRecently(node)
	} else {
		// 容量是否足够
		if c.size >= c.capacity {
			// 删除该节点
			leastRecently := c.removeLeastRecently()
			// 并且从哈希表中移除对应的项
			delete(c.cache, leastRecently.key)
		}
		// 添加新的节点
		node = newNode(key, value, nil, nil)
		c.addRecently(node)
		c.cache[key] = node
	}
}

func (c *Cache) makeRecently(node *DListNode) {
	c.deleteNode(node)
	c.addRecently(node)
}

func (c *Cache) deleteNode(node *DListNode) {
	node.prev.next = node.next
	node.next.prev = node.prev
	c.size -= 1
}

func (c *Cache) addRecently(node *DListNode) {
	node.prev = c.tail.prev
	node.next = c.tail
	c.tail.prev.next = node
	c.tail.prev = node
	c.size += 1
}

func (c *Cache) removeLeastRecently() *DListNode {
	node := c.head.next
	c.deleteNode(node)
	return node
}
