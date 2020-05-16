package cache

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGuardCache(t *testing.T) {
	var cache = WithGuard(New(MapDriver()))
	var count int32
	var value string
	var ch = make(chan string)
	var wg = new(sync.WaitGroup)
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			cache.GetOrSet("test", &value, time.Second, func() (i interface{}, err error) {
				atomic.AddInt32(&count, 1)
				return <-ch, nil
			})
		}()
	}
	ch <- "value"
	wg.Wait()
	assert.Equal(t, count, int32(1))
	assert.Equal(t, value, "value")
}
