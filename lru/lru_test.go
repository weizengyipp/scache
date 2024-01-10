package lru

import (
	"reflect"
	"testing"
)

type V struct {
	str string
}

func newV(str string) V {
	return V{str: str}
}

func (v V) Len() int {
	return len(v.str)
}

type String string

func (d String) Len() int {
	return len(d)
}

func TestAddCache(t *testing.T) {
	v1 := newV("1234")
	v2 := newV("12345")
	lru := New(int64(0), nil)
	lru.Add("key1", v1)
	lru.Add("key1", v2)
	if lru.nbytes != int64(len("key1")+v2.Len()) {
		t.Errorf("expect %d, but got %d", int64(len("key1")+v2.Len()), lru.nbytes)
	}
}

func TestGet(t *testing.T) {
	lru := New(int64(0), nil)
	lru.Add("key1", String("1234"))
	if value, ok := lru.Get("key1"); !ok || value.Len() != 4 || string(value.(String)) != "1234" {
		t.Errorf("expect %v, but got %v", true, ok)
	}
	if _, ok := lru.Get("key2"); ok {
		t.Errorf("get unkown key2, expect %v, but got %v", false, ok)
	}
}

func TestGet2(t *testing.T) {
	lru := New(int64(0), nil)
	lru.Add("key1", newV("12345"))
	if value, ok := lru.Get("key1"); !ok || value.Len() != 5 || (value.(V).str) != "12345" {
		t.Errorf("expect %v, but got %v", true, ok)
	}
}

func TestRemoveOldest(t *testing.T) {
	k1, k2, k3 := "key1", "key2", "key3"
	v1, v2, v3 := "value1", "value2", "value3"
	cap := len(k1 + k2 + v1 + v2)
	lru := New(int64(cap), nil)
	lru.Add(k1, newV(v1))
	lru.Add(k2, newV(v2))
	lru.Add(k3, newV(v3))
	if _, ok := lru.Get("key1"); ok || lru.Len() != 2 {
		t.Fatalf("Removeoldest key1 failed")
	}
}

func TestRemoveOldest2(t *testing.T) {
	k1, k2, k3 := "key1", "key2", "key3"
	v1, v2, v3 := "value1", "value2", "value3"
	cap := len(k1 + k2 + k3 + v1 + v2 + v3)
	lru := New(int64(cap), nil)
	lru.Add(k1, newV(v1))
	lru.Add(k2, newV(v2))
	lru.Add(k3, newV(v3))
	lru.RemoveOldest()
	if _, ok := lru.Get(k1); ok || lru.Len() != 2 {
		t.Fatalf("Removeoldest key1 failed, expect %v len %d, but got %v, len %d", false, 2, ok, lru.Len())
	}
}

var evictedKeys = make([]string, 0)

func callback(k string, v Value) {
	evictedKeys = append(evictedKeys, k)
}

func TestOnEvicted(t *testing.T) {
	lru := New(int64(20), callback)
	lru.Add("key1", String("12345678910"))
	lru.Add("key2", String("123"))
	lru.Add("key3", String("123"))
	expect := []string{"key1"}
	if !reflect.DeepEqual(expect, evictedKeys) {
		t.Fatalf("Call OnEvicted failed, expect %v, but got %v", expect, evictedKeys)
	}
	if lru.Len() != 2 {
		t.Fatalf("lru.Len() expected 2, but got %d", lru.Len())
	}
	if _, ok := lru.Get("key1"); ok != false {
		t.Fatalf("lru.Get(key1) expected nil, but got %v", ok)
	}
	if _, ok := lru.Get("key2"); ok != true {
		t.Fatalf("lru.Get(key2) expected true, but got %v", ok)
	}
	if _, ok := lru.Get("key3"); ok != true {
		t.Fatalf("lru.Get(key3) expected true, but got %v", ok)
	}

}
