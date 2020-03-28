package cache

import (
	"sync"
	"time"
)

type mapDriver struct {
	syncMap *sync.Map
}

func (md mapDriver) Get(key string) ([]byte, bool, error) {
	if v, ok := md.syncMap.Load(key); ok {
		bs, ok := v.(*Item).GetValue()
		return bs, ok, nil
	}
	return nil, false, nil
}

func (md mapDriver) Set(key string, val []byte, ttl time.Duration) error {
	md.syncMap.Store(key, NewItem(val, ttl))
	return nil
}

func (md mapDriver) Has(key string) (bool, error) {
	_, ok := md.syncMap.Load(key)
	return ok, nil
}

func (md mapDriver) Del(key string) error {
	md.syncMap.Delete(key)
	return nil
}

func MapDriver() DriverInterface {
	return mapDriver{syncMap: new(sync.Map)}
}
