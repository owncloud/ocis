package sync

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"sync"
	"testing"
	"time"
)

func cacheRunner(size int) (*Cache, func(f func(i int))) {
	c := NewCache(size)
	run := func(f func(i int)) {
		wg := sync.WaitGroup{}
		for i := 0; i < size; i++ {
			wg.Add(1)
			go func(i int) {
				f(i)
				wg.Done()
			}(i)
		}
		wg.Wait()
	}

	return &c, run
}

func BenchmarkCache(b *testing.B) {
	b.ReportAllocs()
	size := 1024
	c, cr := cacheRunner(size)

	cr(func(i int) { c.Store(strconv.Itoa(i), i, time.Now().Add(10*time.Millisecond)) })
	cr(func(i int) { c.Delete(strconv.Itoa(i)) })
}

func TestCache(t *testing.T) {
	size := 1024
	c, cr := cacheRunner(size)

	cr(func(i int) { c.Store(strconv.Itoa(i), i, time.Now().Add(10*time.Millisecond)) })
	assert.Equal(t, size, int(c.length), "length is atomic")

	cr(func(i int) { c.Delete(strconv.Itoa(i)) })
	assert.Equal(t, 0, int(c.length), "delete is atomic")

	cr(func(i int) {
		time.Sleep(11 * time.Millisecond)
		c.evict()
	})
	assert.Equal(t, 0, int(c.length), "evict is atomic")
}

func TestCache_Load(t *testing.T) {
	size := 1024
	c, cr := cacheRunner(size)

	cr(func(i int) {
		c.Store(strconv.Itoa(i), i, time.Now().Add(10*time.Second))
	})

	cr(func(i int) {
		assert.Equal(t, i, c.Load(strconv.Itoa(i)).V, "entry value is the same")
	})

	cr(func(i int) {
		assert.Nil(t, c.Load(strconv.Itoa(i+size)), "entry is nil if unknown")
	})

	cr(func(i int) {
		wait := 10 * time.Millisecond
		c.Store(strconv.Itoa(i), i, time.Now().Add(wait))
		time.Sleep(wait + 1)
		assert.Nil(t, c.Load(strconv.Itoa(i)), "entry is nil if it's expired")
	})
}

func TestCache_Store(t *testing.T) {
	c, cr := cacheRunner(1024)

	cr(func(i int) {
		c.Store(strconv.Itoa(i), i, time.Now().Add(10*time.Millisecond))
		assert.Equal(t, i, c.Load(strconv.Itoa(i)).V, "new entries can be added")
	})

	cr(func(i int) {
		replacedExpiration := time.Now().Add(10 * time.Millisecond)
		c.Store(strconv.Itoa(i), "old", time.Now().Add(10*time.Minute))
		c.Store(strconv.Itoa(i), "updated", replacedExpiration)
		assert.Equal(t, "updated", c.Load(strconv.Itoa(i)).V, "entry values can be updated")
		assert.Equal(t, replacedExpiration, c.Load(strconv.Itoa(i)).expiration, "entry expiration can be updated")
	})
}

func TestCache_Delete(t *testing.T) {
	c, cr := cacheRunner(1024)

	cr(func(i int) {
		c.Store(strconv.Itoa(i), i, time.Now().Add(10*time.Millisecond))
		c.Delete(strconv.Itoa(i))
		assert.Nil(t, c.Load(strconv.Itoa(i)), "entries can be deleted")
	})

	assert.Equal(t, 0, int(c.length), "removing a entry decreases the cache size")
}
