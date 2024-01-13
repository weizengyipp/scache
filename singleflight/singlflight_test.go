package singleflight

import (
	"testing"
	"time"
)

func TestDo(t *testing.T) {
	var g Group
	v, err := g.Do("key", func() (interface{}, error) {
		return 1, nil
	})
	if v != 1 || err != nil {
		t.Errorf("Got %v,%v; want 1, nil", v, err)
	}
}

func TestDoConcurrent(t *testing.T) {
	cases := []struct {
		key string
		exp int
	}{
		{"key1", 1},
		{"key1", 1},
		{"key1", 1},
		{"key1", 1},
		{"key1", 1},
		{"key1", 1},
		{"key1", 1},
		{"key1", 1},
		{"key1", 1},
		{"key1", 1},
		{"key1", 1},
		{"key1", 1},
	}
	var g Group

	for _, c := range cases {
		go g.Do(c.key, func() (interface{}, error) {
			t.Log("call", c.key)
			return c.exp, nil
		})

	}

	time.Sleep(30 * time.Millisecond)

}
