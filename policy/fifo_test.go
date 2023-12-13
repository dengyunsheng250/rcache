package policy

import "testing"

func TestFIFO(t *testing.T) {

	c := NewFIFOCache(1000, nil)
	c.Add("key1", String("value1"))
	c.Add("key2", String("value2"))
	c.Add("key3", String("value3"))
	c.Add("key4", String("value4"))

	if c.Len() != 4 {
		t.Fatalf("cache get failed")
	}
}
