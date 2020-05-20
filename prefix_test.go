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

func TestPrefixCache(t *testing.T) {
	var c = WithPrefix("test", New(MapDriver()))
	ttl := time.Millisecond * 2
	for key, val := range testDataMap {
		var tmp consumer
		assert.NotNil(t, c.Get(key, &tmp))
		assert.Nil(t, c.Set(key, val, ttl))
		assert.Nil(t, c.Get(key, &tmp))
		assert.EqualValues(t, val, tmp)
		time.Sleep(ttl) // 等待过期
		assert.NotNil(t, c.Get(key, &tmp))
	}

	for key, val := range testDataMap {
		var tmp consumer
		var err = c.GetOrSet(key, &tmp, ttl, func() (i interface{}, err error) {
			return val, nil
		})
		assert.Nil(t, err)
		assert.EqualValues(t, val, tmp)
		assert.Nil(t, c.Del(key))
	}
}
