package cache

import (
	"sync"
	"time"
)

type guardCache struct {
	guard guard
	cache Interface
}

func (pc *guardCache) Has(key string) error {
	_, err := pc.guard.run("has."+key,
		func() (interface{}, error) {
			return nil, pc.cache.Has(key)
		},
	)
	return err
}

func (pc *guardCache) Get(key string, dest interface{}) error {
	_, err := pc.guard.run("get."+key,
		func() (interface{}, error) {
			return nil, pc.cache.Get(key, dest)
		},
	)
	return err
}

func (pc *guardCache) Set(key string, val interface{}, ttl time.Duration) error {
	return pc.cache.Set(key, val, ttl)
}

func (pc *guardCache) Del(key string) error {
	return pc.cache.Del(key)
}

func (pc *guardCache) SetOrDel(key string, val interface{}, ttl time.Duration) error {
	return pc.cache.SetOrDel(key, val, ttl)
}

func (pc *guardCache) GetOrSet(key string, dest interface{}, ttl time.Duration, getter func() (interface{}, error)) error {
	_, err := pc.guard.run("get.set."+key,
		func() (interface{}, error) {
			return nil, pc.cache.GetOrSet(key, dest, ttl, getter)
		},
	)
	return err
}

func WithGuard(cache Interface) Interface {
	return &guardCache{
		cache: cache,
		guard: newGuard(),
	}
}

// guard 防止并发
type guard struct {
	locker    sync.Mutex
	callerMap map[string]*caller
}

func newGuard() guard {
	return guard{
		callerMap: make(map[string]*caller),
	}
}

func (g *guard) run(key string, callback func() (interface{}, error)) (interface{}, error) {
	g.locker.Lock()
	c, ok := g.callerMap[key]
	if ok {
		g.locker.Unlock()
		c.waiter.Wait()
		return c.result()
	} else {
		c = newCall()
		g.callerMap[key] = c
		g.locker.Unlock()
	}
	c.run(callback)
	g.locker.Lock()
	delete(g.callerMap, key)
	g.locker.Unlock()
	return c.result()
}

type caller struct {
	value  interface{}
	error  error
	waiter sync.WaitGroup
}

func newCall() *caller {
	c := new(caller)
	c.waiter.Add(1)
	return c
}

func (c *caller) run(fn func() (interface{}, error)) {
	c.value, c.error = fn()
	c.waiter.Done()
}

func (c *caller) result() (interface{}, error) {
	return c.value, c.error
}
