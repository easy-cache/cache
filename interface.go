package cache

import (
	"time"
)

// CodecInterface 编码接口
type CodecInterface interface {
	Encode(val interface{}) ([]byte, error)
	Decode(bts []byte, dest interface{}) error
}

// LoggerInterface 日志接口
type LoggerInterface interface {
	Errorf(format string, args ...interface{})
}

// DriverInterface 驱动接口
type DriverInterface interface {
	Get(key string) ([]byte, bool, error)
	Set(key string, val []byte, ttl time.Duration) error
	Del(key string) error
	Has(key string) (bool, error)
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
