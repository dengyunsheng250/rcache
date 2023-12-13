package policy

import (
	"container/list"
)

type Cache struct {
	MaxBytes int
	nbytes   int
	ll       *list.List
	// 回调函数
	OnEvicted func(key any, value Value)

	cache map[any]*list.Element
}

func New(maxBytes int, onEvicted func(key any, value Value)) *Cache {
	c := &Cache{
		MaxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[any]*list.Element),
		OnEvicted: onEvicted,
	}

	return c
}

func (c *Cache) Get(key any) (any, bool) {
	if c.cache == nil {
		return nil, false
	}
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		return ele.Value.(*entry).value, true
	}
	return nil, false
}

func (c *Cache) RemoveOldest() {
	if c.cache == nil {
		return
	}
	ele := c.ll.Back()
	if ele != nil {
		c.removeElement(ele)
	}
}

func (c *Cache) Add(key any, value Value) {
	if c.cache == nil {
		c.cache = make(map[interface{}]*list.Element)
		c.ll = list.New()
	}
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		ele.Value.(*entry).value = value
		return
	}
	ele := c.ll.PushFront(&entry{key, value})
	c.cache[key] = ele
	c.nbytes += value.Len()
	if c.nbytes > c.MaxBytes {
		c.RemoveOldest()
	}
}

func (c *Cache) Len() int {
	if c.cache == nil {
		return 0
	}
	return c.ll.Len()
}

func (c *Cache) Remove(key any) {
	if c.cache == nil {
		return
	}
	if ele, ok := c.cache[key]; ok {
		c.removeElement(ele)
	}
}

func (c *Cache) Contains(key any) bool {
	_, ok := c.cache[key]
	return ok
}

func (c *Cache) removeElement(e *list.Element) {
	c.ll.Remove(e)
	kv := e.Value.(*entry)
	delete(c.cache, kv.key)
	if c.OnEvicted != nil {
		c.OnEvicted(kv.key, kv.value)
	}
	c.nbytes -= kv.value.Len()
}

func (c *Cache) Clear() {
	if c.OnEvicted != nil {
		for _, e := range c.cache {
			kv := e.Value.(*entry)
			c.OnEvicted(kv.key, kv.value)
		}
	}
	c.ll = nil
	c.cache = nil
	c.nbytes = 0
}
