package cache

import (
	"reflect"
	"time"
)

type cache struct {
	codec  CodecInterface
	driver DriverInterface
	logger LoggerInterface
}

func (c *cache) Has(key string) bool {
	_, ok, err := c.driver.Get(key)
	if err != nil {
		c.logger.Errorf("has key [%s] failed, err = %s", key, err)
	}
	return ok
}

func (c *cache) Get(key string, dest interface{}) bool {
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

func (c *cache) Set(key string, val interface{}, ttl time.Duration) bool {
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

func (c *cache) Del(key string) bool {
	if err := c.driver.Del(key); err != nil {
		c.logger.Errorf("del key [%s] failed, err = %s", err)
		return false
	}
	return true
}

func (c *cache) SetOrDel(key string, val interface{}, ttl time.Duration) bool {
	return c.Set(key, val, ttl) || c.Del(key)
}

func (c *cache) GetOrSet(key string, dest interface{}, ttl time.Duration, getter func() (interface{}, error)) error {
	if c.Get(key, dest) {
		return nil
	} else if v, err := getter(); err != nil {
		return err
	} else {
		reflect.ValueOf(dest).Elem().Set(reflect.ValueOf(v))
		c.Set(key, v, ttl) // set cache with ttl
		return nil
	}
}

func New(args ...interface{}) Interface {
	c := cache{
		codec:  JsonCodec(),
		driver: NullDriver(),
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
