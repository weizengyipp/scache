package singleflight

import (
	"log"
	"sync"
)

type call struct {
	wg  sync.WaitGroup
	val interface{}
	err error
}

type Group struct {
	mu sync.Mutex
	m  map[string]*call
}

func (g *Group) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	log.Println("DO called with key", key)
	g.mu.Lock()
	if g.m == nil {

		g.m = make(map[string]*call)
	}
	log.Println("call map exists")
	if c, ok := g.m[key]; ok {
		log.Println("call for key exists,call waiting", key)
		g.mu.Unlock()
		c.wg.Wait()
		return c.val, c.err
	}
	log.Println("create call for key", key)
	c := new(call)
	log.Println("call waiting")
	c.wg.Add(1)
	g.m[key] = c
	g.mu.Unlock()

	c.val, c.err = fn()
	log.Println("get v from call for key", key)
	c.wg.Done()

	g.mu.Lock()
	delete(g.m, key)
	g.mu.Unlock()

	return c.val, c.err
}
