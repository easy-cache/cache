package cache

import (
	"fmt"
	"time"
)

type PrefixCache struct {
	prefix string
	cache  *Cache
}

func (pc *PrefixCache) buildKey(key string) string {
	return fmt.Sprintf("%s.%s", pc.prefix, key)
}

func (pc *PrefixCache) Has(key string) bool {
	return pc.cache.Has(pc.buildKey(key))
}

func (pc *PrefixCache) Get(key string, dest interface{}) bool {
	return pc.cache.Get(pc.buildKey(key), dest)
}

func (pc *PrefixCache) Set(key string, val interface{}, ttl time.Duration) bool {
	return pc.cache.Set(pc.buildKey(key), val, ttl)
}

func (pc *PrefixCache) Del(key string) bool {
	return pc.cache.Del(pc.buildKey(key))
}

func (pc *PrefixCache) SetOrDel(key string, val interface{}, ttl time.Duration) bool {
	return pc.cache.SetOrDel(pc.buildKey(key), val, ttl)
}

func (pc *PrefixCache) GetOrSet(key string, dest interface{}, ttl time.Duration, getter func() (interface{}, error)) bool {
	return pc.cache.GetOrSet(pc.buildKey(key), dest, ttl, getter)
}

func WithPrefix(prefix string, args ...interface{}) *PrefixCache {
	if prefix == "" {
		panic("easy-cache: the prefix can not be empty")
	}
	return &PrefixCache{
		prefix: prefix,
		cache:  New(args...),
	}
}
