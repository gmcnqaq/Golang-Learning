package consistehash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

// Hash 将字节映射到 uint32
type Hash func(data []byte) uint32

// Map 一致性哈希算法的主数据结构
// Hash 函数 hash，虚拟节点倍数 replicas，哈希环 keys，虚拟节点与真实节点的映射表 hashMap，键是虚拟节点的哈希值，值是真实节点的名称
type Map struct {
	hash     Hash
	replicas int
	keys     []int // 已排序
	hashMap  map[int]string
}

func NewMap(replicas int, fn Hash) *Map {
	m := &Map{
		replicas: replicas,
		hash:     fn,
		hashMap:  make(map[int]string),
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

// Add 添加真实节点/机器的方法
// Add 函数允许传入0 或者多个真实节点的名称
// 对每个真实节点 key，对应创建 m.replicas 个虚拟节点，通过添加编号的方式区分不同虚拟节点
// 使用 m.hash() 计算虚拟节点的哈希值，再 append 添加到环上
func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			m.keys = append(m.keys, hash)
			m.hashMap[hash] = key
		}
	}
	sort.Ints(m.keys)
}

// Get 选择节点
// 根据 key 的哈希值，顺时针找到第一个匹配的虚拟节点的下标 idx，从 m.keys 中获取到对应的哈希值
// 如果 idx == len(m.keys)，则应该选择 m.keys[0]，因为 m.keys 是一个环状结构，所以用取玉树的方式来处理这种情况
func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}
	hash := int(m.hash([]byte(key)))
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})

	return m.hashMap[m.keys[idx%len(m.keys)]]
}
