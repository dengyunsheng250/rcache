package policy

import (
	"fmt"
	"testing"
)

type Int int

func (i Int) Len() int {
	return 4
}

func TestGet(t *testing.T) {
	lru := New(10, nil)
	lru.Add("a", Int(1))
	lru.Add("b", Int(2))
	lru.Add("c", Int(3))
	if lru.Len() != 3 {
		t.Fatalf("cache get failed")
	}
}

type String string

func (s String) Len() int {
	return len(s)
}

func TestRemove(t *testing.T) {
	OnEvicted := func(key any, value Value) {
		fmt.Printf("hello %s", key)
	}
	lru := New(100, OnEvicted)
	lru.Add("k1", String("v1"))
	lru.Add("k2", String("v2"))
}
