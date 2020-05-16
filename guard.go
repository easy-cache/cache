package cache

import (
	"time"

	"github.com/letsfire/utils"
)

type guardCache struct {
	cache Interface
	guard *utils.Guard
}

func (pc *guardCache) Has(key string) bool {
	res, _ := pc.guard.Run("has."+key,
		func() (interface{}, error) {
			return pc.cache.Has(key), nil
		},
	)
	return res.(bool)
}

func (pc *guardCache) Get(key string, dest interface{}) bool {
	res, _ := pc.guard.Run("get."+key,
		func() (interface{}, error) {
			return pc.cache.Get(key, dest), nil
		},
	)
	return res.(bool)
}

func (pc *guardCache) Set(key string, val interface{}, ttl time.Duration) bool {
	return pc.cache.Set(key, val, ttl)
}

func (pc *guardCache) Del(key string) bool {
	return pc.cache.Del(key)
}

func (pc *guardCache) SetOrDel(key string, val interface{}, ttl time.Duration) bool {
	return pc.cache.SetOrDel(key, val, ttl)
}

func (pc *guardCache) GetOrSet(key string, dest interface{}, ttl time.Duration, getter func() (interface{}, error)) bool {
	res, _ := pc.guard.Run("get.set."+key,
		func() (interface{}, error) {
			return pc.cache.GetOrSet(key, dest, ttl, getter), nil
		},
	)
	return res.(bool)
}

func WithGuard(cache Interface) Interface {
	return &guardCache{
		cache: cache,
		guard: utils.NewGuard(),
	}
}
