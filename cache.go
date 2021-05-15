package cache

import (
	"reflect"
	"time"
)

type cache struct {
	codec  CodecInterface
	driver DriverInterface
}

func (c *cache) Has(key string) error {
	_, ok, err := c.driver.Get(key)
	if err == nil && !ok {
		err = ErrMiss
	}
	return err
}

func (c *cache) Get(key string, dest interface{}) error {
	bs, ok, err := c.driver.Get(key)
	if err == nil {
		if ok {
			err = c.codec.Decode(bs, dest)
		} else {
			err = ErrMiss
		}
	}
	return err
}

func (c *cache) Set(key string, val interface{}, ttl time.Duration) error {
	mustGtZero(ttl)
	vbs, err := c.codec.Encode(val)
	if err == nil {
		err = c.driver.Set(key, vbs, ttl)
	}
	return err
}

func (c *cache) Del(key string) error {
	return c.driver.Del(key)
}

func (c *cache) SetOrDel(key string, val interface{}, ttl time.Duration) error {
	if err := c.Set(key, val, ttl); err == nil {
		return nil
	}
	return c.Del(key)
}

func (c *cache) GetOrSet(key string, dest interface{}, ttl time.Duration, getter func() (interface{}, error)) error {
	if err := c.Get(key, dest); err != ErrMiss {
		return err
	} else if v, err := getter(); err != nil {
		return err
	} else {
		reflect.ValueOf(dest).Elem().Set(reflect.ValueOf(v))
		return c.Set(key, v, ttl)
	}
}

func New(args ...interface{}) Interface {
	c := cache{
		codec:  JsonCodec(),
		driver: NullDriver(),
	}
	for _, arg := range args {
		switch v := arg.(type) {
		case CodecInterface:
			c.codec = v
		case DriverInterface:
			c.driver = v
		default:
			panic("unsupported arg type")
		}
	}
	return &c
}
