package test_GeeCache

import (
	"Golang/Learning/geecache/test-GeeCache/singleflight"
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

// Getter 回调函数，如果缓存不存在，从数据源获取数据
type Getter interface {
	Get(key string) ([]byte, error)
}

// GetterFunc 定义函数类型，并实现 Getter 接口的 Get 方法
// 函数类型实现某一个接口，称之为接口型函数，方便使用者在调用时既能够传入函数作为参数，也能够传入实现了该接口的结构体作为参数
type GetterFunc func(key string) ([]byte, error)

func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

// Group 一个 Group 可以认为是一个缓存的命名空间，每个 Group 拥有一个唯一的名称 name。比如缓存学生的成绩命名为 scores，学生名字命名为 names
// 第二个属性 getter Getter，即缓存未命中时获取源数据的回调
// 第三个属性 mainCache，即前面实现的并发缓存
type Group struct {
	name      string
	getter    Getter
	mainCache cache
	peers PeerPicker
	loader *singleflight.Group
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

// NewGroup 用来实例化 Group，并将 group 存储在全局变量 groups 中
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
		loader: &singleflight.Group{},
	}
	groups[name] = g
	return g
}

// GetGroup 用来得到特定名称的 Group
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
	// 如果缓存中存在
	if v, exist := g.mainCache.get(key); exist {
		log.Println("[GeeCache] hit")
		return v, nil
	}
	// 否则回调
	return g.load(key)
}

/*func (g *Group) load(key string) (ByteView, error) {
	return g.GetLocally(key)
}*/
func (g *Group) load(key string) (ByteView, error) {
	viewi, err := g.loader.Do(key, func() (interface{}, error) {
		if g.peers != nil {
			if peer, ok := g.peers.PickPeer(key); ok {
				if value, err := g.getFromPeer(peer, key); err == nil {
					return value, err
				} else {
					log.Println("[GeeCache] Failed to get from peer", err)
				}
			}
		}
		return g.GetLocally(key)
	})
	if err == nil {
		return viewi.(ByteView), nil
	}
	return *NewByteView(nil), nil
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

func (g *Group) RegisterPeers(peers PeerPicker)  {
	if g.peers != nil {
		panic("RegisterPeerPicker called more than once")
	}
	g.peers = peers
}

func (g *Group) getFromPeer(peer PeerGetter, key string) (ByteView, error) {
	bytes, err := peer.Get(g.name, key)
	if err != nil {
		return *NewByteView(nil), err
	}
	return ByteView{b: bytes}, nil
}