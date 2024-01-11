package scache

import (
	"fmt"
	"log"
	"reflect"
	"testing"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func TestGetter(t *testing.T) {
	var f Getter = GetterFunc(func(key string) ([]byte, error) {
		return []byte(key), nil
	})
	expect := []byte("key")
	if v, err := f.Get("key"); err != nil || !reflect.DeepEqual(v, expect) {
		t.Errorf("callback failed")
	}
}

var loadCounts = make(map[string]int, len(db))

var f Getter = GetterFunc(func(key string) ([]byte, error) {
	log.Println("db searching:", key)
	if v, ok := db[key]; ok {
		loadCounts[key] += 1
		return []byte(v), nil
	}
	return nil, fmt.Errorf("%s not exist", key)
})

func TestGetGroup(t *testing.T) {
	groupName := "scores"
	NewGroup(groupName, 2<<10, f)
	if group := GetGroup(groupName); group == nil || group.name != groupName {
		t.Fatalf("group name %s get wrong group %s", groupName, group.name)
	}
	if group := GetGroup(groupName + "unknow"); group != nil {
		t.Fatalf("group name %s should be nil ,but get %s", groupName, group.name)
	}
}

func TestGet(t *testing.T) {
	group := NewGroup("scores", 2<<10, f)
	for k, v := range db {
		if view, err := group.Get(k); err != nil || view.String() != v {
			t.Errorf("failed to get value of %s", k)
		}
	}
	if _, err := group.Get("unknown"); err == nil {
		t.Fatalf("expect an error but not get one")
	}
	for k, v := range db {
		if view, err := group.Get(k); err != nil || view.String() != v {
			t.Errorf("failed to get value of %s", k)
		}
		if loadCounts[k] > 1 {
			t.Errorf("cache miss:key %s loaded %d times", k, loadCounts[k])
		}
	}
}
