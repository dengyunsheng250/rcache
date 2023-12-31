package geecache

import (
	"rcache/policy"
	"sync"
)

type cache struct {
	mu         sync.Mutex
	lru        *policy.TwoQueueCache
	cacheBytes int
}

func (c *cache) add(key string, value ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		c.lru = policy.New2Q(c.cacheBytes, nil)
	}
	c.lru.Add(key, value)
}

func (c *cache) get(key string) (value ByteView, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		return
	}
	if v, ok := c.lru.Get(key); ok {
		return v.(ByteView), ok
	}
	return
}

// func (c *cache) remove(key string) {
// 	c.mu.Lock()
// 	defer c.mu.Unlock()
// 	c.lru.Remove(key)
// }
