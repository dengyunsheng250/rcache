package policy

import "container/list"

type FIFOCache struct {
	MaxBytes int
	nbytes   int
	ll       *list.List
	// 回调函数
	OnEvicted func(key any, value Value)

	cache map[any]*list.Element
}

func NewFIFOCache(maxBytes int, onEvicted func(key any, value Value)) *FIFOCache {
	c := &FIFOCache{
		MaxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[any]*list.Element),
		OnEvicted: onEvicted,
	}
	return c
}

func (c *FIFOCache) Get(key any) (any, bool) {
	if c.cache == nil {
		return nil, false
	}
	if ele, ok := c.cache[key]; ok {
		return ele.Value.(*entry).value, true
	}
	return nil, false
}

func (c *FIFOCache) Add(key any, value Value) {
	if ele, ok := c.cache[key]; ok {
		ele.Value.(*entry).value = value
		return
	}
	ele := c.ll.PushBack(&entry{key, value})
	c.cache[key] = ele
	c.nbytes += value.Len()
	if c.nbytes > c.MaxBytes {
		c.Removeold()
	}
}

func (c *FIFOCache) Removeold() {
	ele := c.ll.Front()
	if ele != nil {
		c.ll.Remove(ele)
		delete(c.cache, ele.Value.(*entry).key)
	}
}

func (c *FIFOCache) Contains(key any) bool {
	_, ok := c.cache[key]
	return ok
}

func (c *FIFOCache) Len() int {
	if c.cache == nil {
		return 0
	}
	return c.ll.Len()
}

func (c *FIFOCache) Remove(key any) {
	if c.Contains(key) {
		v := c.cache[key].Value.(*entry).value
		c.nbytes -= v.Len()
		c.ll.Remove(c.cache[key])
		delete(c.cache, key)
		if c.OnEvicted != nil {
			c.OnEvicted(key, v)
		}
	}
}
