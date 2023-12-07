package expire

import (
	"container/list"
	"time"
)

type EvictCallback[K any, V any] func(key K, value V)

type Cache struct {
	MaxEntries int

	ll *list.List

	OnEvicted func(key any, value any)

	cache map[any]*list.Element

	ttl time.Duration
}

type entry struct {
	key   any
	value any
}
