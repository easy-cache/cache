package cache

import (
	"reflect"
	"time"
)

type Cache struct {
	codec  CodecInterface
	driver DriverInterface
	logger LoggerInterface
}

func (c *Cache) Has(key string) bool {
	ok, err := c.driver.Has(key)
	if err != nil {
		c.logger.Errorf("has key [%s] failed, err = %s", key, err)
	}
	return ok
}

func (c *Cache) Get(key string, dest interface{}) bool {
	bs, ok, err := c.driver.Get(key)
	if err != nil {
		c.logger.Errorf("get key [%s] failed, err = %s", key, err)
	} else if ok {
		if err = c.codec.Decode(bs, dest); err == nil {
			return true
		}
		c.logger.Errorf("get key [%s] failed, bytes = %s, err = %s", key, bs, err)
	}
	return false
}

func (c *Cache) Set(key string, val interface{}, ttl time.Duration) bool {
	mustGtZero(ttl)
	vbs, err := c.codec.Encode(val)
	if err == nil {
		if err = c.driver.Set(key, vbs, ttl); err == nil {
			return true
		}
	}
	c.logger.Errorf("set key [%s] failed, err = %s", key, err)
	return false
}

func (c *Cache) Del(key string) bool {
	if err := c.driver.Del(key); err != nil {
		c.logger.Errorf("del key [%s] failed, err = %s", err)
		return false
	}
	return true
}

func (c *Cache) SetOrDel(key string, val interface{}, ttl time.Duration) bool {
	return c.Set(key, val, ttl) || c.Del(key)
}

func (c *Cache) GetOrSet(key string, dest interface{}, ttl time.Duration, getter func() (interface{}, error)) bool {
	if c.Get(key, dest) {
		return true
	} else if v, err := getter(); err != nil {
		c.logger.Errorf("get or set key [%s] failed, err = %s", key, err)
		return false
	} else {
		reflect.ValueOf(dest).Elem().Set(reflect.ValueOf(v))
		return c.Set(key, v, ttl)
	}
}

func New(args ...interface{}) *Cache {
	c := Cache{
		codec:  JsonCodec(),
		driver: MockDriver(),
		logger: StderrLogger(),
	}
	for _, arg := range args {
		switch v := arg.(type) {
		case CodecInterface:
			c.codec = v
		case DriverInterface:
			c.driver = v
		case LoggerInterface:
			c.logger = v
		default:
			panic("unsupported arg type")
		}
	}
	return &c
}
