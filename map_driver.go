package cache

import (
	"sync"
	"time"
)

type mapDriver struct {
	syncMap *sync.Map
}

func (mc mapDriver) Get(key string) ([]byte, bool, error) {
	if v, ok := mc.syncMap.Load(key); ok {
		bs, ok := v.(*Item).GetValue()
		return bs, ok, nil
	}
	return nil, false, nil
}

func (mc mapDriver) Set(key string, val []byte, ttl time.Duration) error {
	mc.syncMap.Store(key, NewItem(val, ttl))
	return nil
}

func (mc mapDriver) Has(key string) (bool, error) {
	_, ok := mc.syncMap.Load(key)
	return ok, nil
}

func (mc mapDriver) Del(key string) error {
	mc.syncMap.Delete(key)
	return nil
}

func MapDriver() DriverInterface {
	return mapDriver{syncMap: new(sync.Map)}
}
