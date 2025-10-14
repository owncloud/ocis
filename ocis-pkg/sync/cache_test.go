package sync

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"sync"
	"testing"
	"time"
)

func cacheRunner(size int) (*Cache, func(f func(v string))) {
	c := NewCache(size)
	run := func(f func(v string)) {
		wg := sync.WaitGroup{}
		for i := 0; i < size; i++ {
			wg.Add(1)
			go func(v string) {
				f(v)
				wg.Done()
			}(strconv.Itoa(i))
		}
		wg.Wait()
	}

	return &c, run
}

func BenchmarkCache(b *testing.B) {
	b.ReportAllocs()
	size := 1024
	c, cr := cacheRunner(size)

	cr(func(v string) { c.Store(v, v, time.Now().Add(100*time.Millisecond)) })
	cr(func(v string) { c.Delete(v) })
}

func TestCache(t *testing.T) {
	size := 1024
	c, cr := cacheRunner(size)

	cr(func(v string) { c.Store(v, v, time.Now().Add(100*time.Millisecond)) })
	assert.Equal(t, size, int(c.length), "length is atomic")

	cr(func(v string) { c.Delete(v) })
	assert.Equal(t, 0, int(c.length), "delete is atomic")

	cr(func(v string) {
		time.Sleep(101 * time.Millisecond)
		c.evict()
	})
	assert.Equal(t, 0, int(c.length), "evict is atomic")
}

func TestCache_Load(t *testing.T) {
	size := 1024
	c, cr := cacheRunner(size)

	cr(func(v string) {
		c.Store(v, v, time.Now().Add(10*time.Second))
	})

	cr(func(v string) {
		assert.Equal(t, v, c.Load(v).V, "entry value is the same")
	})

	cr(func(v string) {
		assert.Nil(t, c.Load(v+strconv.Itoa(size)), "entry is nil if unknown")
	})

	cr(func(v string) {
		wait := 100 * time.Millisecond
		c.Store(v, v, time.Now().Add(wait))
		time.Sleep(wait + 1)
		assert.Nil(t, c.Load(v), "entry is nil if it's expired")
	})
}

func TestCache_Store(t *testing.T) {
	c, cr := cacheRunner(1024)

	cr(func(v string) {
		c.Store(v, v, time.Now().Add(100*time.Millisecond))
		assert.Equal(t, v, c.Load(v).V, "new entries can be added")
	})

	cr(func(v string) {
		replacedExpiration := time.Now().Add(10 * time.Minute)
		c.Store(v, "old", time.Now().Add(10*time.Minute))
		c.Store(v, "updated", replacedExpiration)
		assert.Equal(t, "updated", c.Load(v).V, "entry values can be updated")
		assert.Equal(t, replacedExpiration, c.Load(v).expiration, "entry expiration can be updated")
	})
}

func TestCache_Delete(t *testing.T) {
	c, cr := cacheRunner(1024)

	cr(func(v string) {
		c.Store(v, v, time.Now().Add(100*time.Millisecond))
		c.Delete(v)
		assert.Nil(t, c.Load(v), "entries can be deleted")
	})

	assert.Equal(t, 0, int(c.length), "removing a entry decreases the cache size")
}
