package test_GeeCache

import (
	"fmt"
	"log"
	"sync"
)

// Group 是 GeeCache 最核心的数据结构，负责与用户的交互，并且控制缓存值存储和获取的流程

/*
接受 key --> 检查是否被缓存 --是--> 返回缓存值(1)
                ｜否
                ｜----> 是否应当从远程节点获取 --是--> 与远程节点交互 --> 返回缓存值(2)
                               ｜否
                               ｜----> 「回调函数」，获取值，并添加到缓存 --> 返回缓存值(3)
*/

// Getter 回调函数，如果缓存不存在，得到源数据。
// Todo: 回调函数 Getter
type Getter interface {
	Get(key string) ([]byte, error)
}

type GetterFunc func(key string) ([]byte, error)

func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

// Group 一个 Group 可以认为是一个缓存的命名空间，每个 Group 拥有一个唯一的名称 name。
// 第二个属性 getter Getter，即缓存未命中时获取源数据的回调
// 第三个属性 mainCache，即实现的并发缓存
type Group struct {
	name      string
	getter    Getter
	mainCache cache
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

// NewGroup 用来实例化 Group，并将 group 存储在全局变量 groups 中
// Todo：单例 group
func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("nil Getter")
	}
	mu.Lock()
	defer mu.Unlock()
	g := &Group{
		name:      name,
		getter:    getter,
		mainCache: cache{cacheBytes: cacheBytes},
	}
	groups[name] = g
	return g
}

// GetGroup 用来特定名称的 Group
func GetGroup(name string) *Group {
	mu.RLock()
	g := groups[name]
	mu.RUnlock()
	return g
}

func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return *NewByteView(nil), fmt.Errorf("key is required")
	}
	if v, exist := g.mainCache.get(key); exist {
		log.Println("[GeeCache] hit")
		return v, nil
	}
	return g.load(key)
}

func (g *Group) load(key string) (ByteView, error) {
	return g.GetLocally(key)
}

func (g *Group) GetLocally(key string) (ByteView, error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return *NewByteView(nil), err
	}
	value := *NewByteView(cloneBytes(bytes))
	g.populateCache(key, value)
	return value, nil
}

func (g *Group) populateCache(key string, value ByteView) {
	g.mainCache.put(key, value)
}
