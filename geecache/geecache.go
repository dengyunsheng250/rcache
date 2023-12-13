package geecache

import (
	"fmt"
	"rcache/logger"
	"sync"
)

// 缓存不存在时，进行回调
type Getter interface {
	Get(key string) ([]byte, error)
}

type GetterFunc func(key string) ([]byte, error)

// 函数类型实现某个接口
// 和给一个基本类型实现一个接口是一个原理
func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

// Group相当于一个缓存的命名空间
type Group struct {
	name      string // 唯一命名空间
	getter    Getter // 缓存未命中时的回调
	mainCache cache  // 并发缓存
	hotCache  cache  // 热点缓存，对于不再本节点的数据，进行本地缓存，避免大量的网络IO
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

// group对应缓存的命名空间
// 一个name对应一个cache
func NewGroup(name string, cacheBytes int, getter Getter) *Group {
	if getter == nil {
		panic("nil getter")
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
	defer mu.RUnlock()
	g := groups[name]
	return g
}

func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("key is required")
	}
	if v, ok := g.mainCache.get(key); ok {
		logger.Info("[Cache] hit")
		return v, nil
	}
	return g.load(key)
}

func (g *Group) load(key string) (value ByteView, err error) {
	return g.getLocally(key)
}

func (g *Group) getLocally(key string) (ByteView, error) {
	mu.Lock()
	defer mu.Unlock()
	b, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}
	value := NewByteView(cloneBytes(b))
	g.populateCache(key, value)
	return value, nil
}

func (g *Group) populateCache(key string, value ByteView) {
	g.mainCache.add(key, value)
}
