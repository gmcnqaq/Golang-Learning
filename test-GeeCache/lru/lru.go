package lru

type DListNode struct {
	key   string
	value Value
	next  *DListNode
	prev  *DListNode
}

// Value 接口的任意类型，只包含一个 len() 方法返回值所占用的内存大小。通用性
type Value interface {
	len() int
}

type Cache struct {
	cache    map[string]interface{}
	head     *DListNode
	tail     *DListNode
	capacity int64
	size     int64
}

// New 为了方便实例化 Cache
func New(capacity int64) *Cache {
	return &Cache{
		capacity: capacity,
		head:     nil,
		tail:     nil,
		cache:    make(map[string]interface{}),
		size:     0,
	}
}
