package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type address struct {
	Area   string
	Detail string
}

type consumer struct {
	Id        uint64
	Name      string
	Address   address
	CreatedAt time.Time
}

var testDataMap = map[string]consumer{
	"consumer.01": {
		Id:   1,
		Name: "LiLei",
		Address: address{
			Area:   "110000",
			Detail: "王府井200号",
		},
		CreatedAt: time.Now().Round(0),
	},
	"consumer.02": {
		Id:   2,
		Name: "HanMei",
		Address: address{
			Area:   "110000",
			Detail: "王府井201号",
		},
		CreatedAt: time.Now().Add(time.Second).Round(0),
	},
}

func TestCache(t *testing.T) {
	var c = New(MapDriver())
	ttl := time.Millisecond * 2
	for key, val := range testDataMap {
		var tmp consumer
		assert.False(t, c.Get(key, &tmp))
		assert.True(t, c.Set(key, val, ttl))
		assert.True(t, c.Get(key, &tmp))
		assert.EqualValues(t, val, tmp)
		time.Sleep(ttl) // 等待过期
		assert.False(t, c.Get(key, &tmp))
	}

	for key, val := range testDataMap {
		var tmp consumer
		c.GetOrSet(key, &tmp, ttl, func() (i interface{}, err error) {
			return val, nil
		})
		assert.EqualValues(t, val, tmp)
		assert.True(t, c.Del(key))
	}
}