package cache

import (
	"errors"
	"time"
)

var ErrMiss = errors.New("cache miss")

// CodecInterface 编码接口
type CodecInterface interface {
	Encode(val interface{}) ([]byte, error)
	Decode(bts []byte, dest interface{}) error
}

// DriverInterface 驱动接口
type DriverInterface interface {
	Get(key string) ([]byte, bool, error)
	Set(key string, val []byte, ttl time.Duration) error
	Del(key string) error
}

// Interface 缓存封装
type Interface interface {
	Has(key string) error
	Get(key string, dest interface{}) error
	Set(key string, val interface{}, ttl time.Duration) error
	Del(key string) error
	SetOrDel(key string, val interface{}, ttl time.Duration) error
	GetOrSet(key string, dest interface{}, ttl time.Duration, getter func() (interface{}, error)) error
}

// Item 缓存过期包装
type Item struct {
	Value     []byte    `json:"v"`
	ExpiredAt time.Time `json:"e"`
}

func (item *Item) IsExpired() bool {
	return item.ExpiredAt.Before(time.Now())
}

func (item *Item) GetValue() ([]byte, bool) {
	if item.IsExpired() {
		return nil, false
	}
	return item.Value, true
}

func NewItem(val []byte, ttl time.Duration) *Item {
	return &Item{Value: val, ExpiredAt: time.Now().Add(ttl)}
}

// mustGtZero 缓存必须设置有效期
func mustGtZero(ttl time.Duration) {
	if ttl <= 0 {
		panic("cache: the ttl of cache item must be greater than zero")
	}
}
