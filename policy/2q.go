package policy

import (
	"os"
	"rcache/logger"
)

type TwoQueueCache struct {
	Sizes   int
	recent  *Cache
	history *FIFOCache
}

func New2Q(sizes int, onEvicted func(key any, value Value)) *TwoQueueCache {
	recent := New(sizes, onEvicted)
	history := NewFIFOCache(sizes*2, onEvicted)
	q := &TwoQueueCache{
		Sizes:   sizes,
		recent:  recent,
		history: history,
	}
	return q
}

func (q *TwoQueueCache) Add(key any, value Value) {
	if q.recent.Contains(key) {
		q.recent.Add(key, value)
		return
	}
	logger.Info("Debug")
	if q.history.Contains(key) {
		os.Exit(-1)
		q.history.Remove(key)
		q.recent.Add(key, value)
		return
	}
	q.history.Add(key, value)
}

func (q *TwoQueueCache) Get(key any) (any, bool) {
	if val, ok := q.recent.Get(key); ok {
		return val, ok
	}
	if val, ok := q.history.Get(key); ok {
		q.history.Remove(key)
		q.recent.Add(key, val.(Value))
		return val, ok
	}
	return nil, false
}

func (q *TwoQueueCache) Contains(key any) bool {
	return q.history.Contains(key) || q.recent.Contains(key)
}

func (q *TwoQueueCache) Remove(key any) {
	q.recent.Remove(key)
	q.history.Remove(key)
}
