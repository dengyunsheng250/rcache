package lru

import (
	"fmt"
	"testing"
)

func TestGet(t *testing.T) {
	lru := New(10, nil)
	lru.Add("a", 1)
	lru.Add("b", 2)
	lru.Add("c", 3)
	if lru.Len() != 3 {
		t.Fatalf("cache get failed")
	}
}

func TestRemove(t *testing.T) {
	OnEvicted := func(key any, value any) {
		fmt.Printf("hello %s", key)
	}
	lru := New(1, OnEvicted)
	lru.Add("k1", "v1")
	lru.Add("k2", "v2")
}
