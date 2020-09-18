package roles

import (
	"sync"
	"time"

	settings "github.com/owncloud/ocis/settings/pkg/proto/v0"
)

// entry extends a bundle and adds a TTL
type entry struct {
	*settings.Bundle
	inserted time.Time
}

// cache is a cache implementation for roles, keyed by roleIDs.
type cache struct {
	entries map[string]entry
	size    int
	ttl     time.Duration
	m       sync.Mutex
}

// newCache returns a new instance of Cache.
func newCache(size int, ttl time.Duration) cache {
	return cache{
		size:    size,
		ttl:     ttl,
		entries: map[string]entry{},
	}
}

// get gets a role-bundle by a given `roleID`.
func (c *cache) get(roleID string) *settings.Bundle {
	c.m.Lock()
	defer c.m.Unlock()

	if _, ok := c.entries[roleID]; ok {
		return c.entries[roleID].Bundle
	}
	return nil
}

// set sets a roleID / role-bundle.
func (c *cache) set(roleID string, value *settings.Bundle) {
	c.m.Lock()
	defer c.m.Unlock()

	if !c.fits() {
		c.evict()
	}

	c.entries[roleID] = entry{
		value,
		time.Now(),
	}
}

// evict frees memory from the cache by removing entries that exceeded the cache TTL.
func (c *cache) evict() {
	for i := range c.entries {
		if c.entries[i].inserted.Add(c.ttl).Before(time.Now()) {
			delete(c.entries, i)
		}
	}
}

// fits returns whether the cache fits more entries.
func (c *cache) fits() bool {
	return c.size > len(c.entries)
}
