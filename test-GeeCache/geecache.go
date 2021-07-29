package test_GeeCache

import (
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

type Getter interface {
	Get(key int) ([]byte, error)
}

type Group struct {
	name      string
	getter    Getter
	mainCache cache
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

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

func GetGroup(name string) *Group {
	mu.RLock()
	g := groups[name]
	mu.RUnlock()
	return g
}

/*func (g *Group) Get(key int) (int, error) {
	if key == -1 {
		return -1, fmt.Errorf("key is required")
	}
	if v := g.mainCache.get(key); v != -1 {
		log.Panicln("[GeeCache] hit")
		return v, nil
	}
	return g.load(key)
}

func (g *Group) load(key int) (ByteView, error) {

}

func (g *Group) GetLocally(key int) (ByteView, error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}
}*/
