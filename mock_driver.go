package cache

import (
	"time"
)

type mockDriver struct{}

func (mc mockDriver) Get(key string) ([]byte, bool, error) {
	return nil, false, nil
}

func (mc mockDriver) Set(key string, val []byte, ttl time.Duration) error {
	return nil
}

func (mc mockDriver) Has(key string) (bool, error) {
	return false, nil
}

func (mc mockDriver) Del(key string) error {
	return nil
}

func MockDriver() DriverInterface {
	return mockDriver{}
}
