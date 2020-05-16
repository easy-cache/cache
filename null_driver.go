package cache

import (
	"time"
)

type nullDriver struct{}

func (nd nullDriver) Get(key string) ([]byte, bool, error) {
	return nil, false, nil
}

func (nd nullDriver) Set(key string, val []byte, ttl time.Duration) error {
	return nil
}

func (nd nullDriver) Del(key string) error {
	return nil
}

func NullDriver() DriverInterface {
	return nullDriver{}
}
