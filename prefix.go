package cache

import (
	"fmt"
	"time"
)

type prefixCache struct {
	prefix string
	cache  Interface
}

func (pc *prefixCache) buildKey(key string) string {
	return fmt.Sprintf("%s.%s", pc.prefix, key)
}

func (pc *prefixCache) Has(key string) bool {
	return pc.cache.Has(pc.buildKey(key))
}

func (pc *prefixCache) Get(key string, dest interface{}) bool {
	return pc.cache.Get(pc.buildKey(key), dest)
}

func (pc *prefixCache) Set(key string, val interface{}, ttl time.Duration) bool {
	return pc.cache.Set(pc.buildKey(key), val, ttl)
}

func (pc *prefixCache) Del(key string) bool {
	return pc.cache.Del(pc.buildKey(key))
}

func (pc *prefixCache) SetOrDel(key string, val interface{}, ttl time.Duration) bool {
	return pc.cache.SetOrDel(pc.buildKey(key), val, ttl)
}

func (pc *prefixCache) GetOrSet(key string, dest interface{}, ttl time.Duration, getter func() (interface{}, error)) bool {
	return pc.cache.GetOrSet(pc.buildKey(key), dest, ttl, getter)
}

func WithPrefix(prefix string, cache Interface) Interface {
	if prefix == "" {
		panic("easy-cache: the prefix can not be empty")
	}
	return &prefixCache{
		cache:  cache,
		prefix: prefix,
	}
}
