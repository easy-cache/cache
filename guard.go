package cache

import (
	"time"

	"github.com/letsfire/utils"
)

type guardCache struct {
	cache Interface
	guard *utils.Guard
}

func (pc *guardCache) Has(key string) error {
	_, err := pc.guard.Run("has."+key,
		func() (interface{}, error) {
			return nil, pc.cache.Has(key)
		},
	)
	return err
}

func (pc *guardCache) Get(key string, dest interface{}) error {
	_, err := pc.guard.Run("get."+key,
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
	_, err := pc.guard.Run("get.set."+key,
		func() (interface{}, error) {
			return nil, pc.cache.GetOrSet(key, dest, ttl, getter)
		},
	)
	return err
}

func WithGuard(cache Interface) Interface {
	return &guardCache{
		cache: cache,
		guard: utils.NewGuard(),
	}
}
