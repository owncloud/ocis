package roles

import (
	settings "github.com/owncloud/ocis/settings/pkg/proto/v0"
	"github.com/stretchr/testify/assert"
	"strconv"
	"sync"
	"testing"
	"time"
)

func cacheRunner(size int, ttl time.Duration) (*cache, func(f func(v string))) {
	c := newCache(size, ttl)
	run := func(f func(v string)) {
		wg := sync.WaitGroup{}
		for i := 0; i < size; i++ {
			wg.Add(1)
			go func(i int) {
				f(strconv.Itoa(i))
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
	c, cr := cacheRunner(size, 100 * time.Millisecond)

	cr(func(v string) { c.set(v, &settings.Bundle{})})
	cr(func(v string) { c.get(v)})
}

func TestCache(t *testing.T) {
	size := 1024
	ttl := 100 * time.Millisecond
	c, cr := cacheRunner(size, ttl)

	cr(func(v string) {
		c.set(v, &settings.Bundle{Id: v})
	})

	assert.Equal(t, "50", c.get("50").Id, "it returns the right bundle")
	assert.Nil(t, c.get("unknown"), "unknown bundle ist nil")

	time.Sleep(ttl + 1)
	// roles cache has no access to evict, adding new items triggers a cleanup
	c.set("evict", nil)
	assert.Nil(t, c.get("50"), "old bundles get removed")
}
