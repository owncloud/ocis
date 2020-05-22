package cache

import (
	"fmt"
	"sync"
)

// Entry represents an entry on the cache. You can type assert on V.
type Entry struct {
	V     interface{}
	Valid bool
}

// Cache is a barebones cache implementation.
type Cache struct {
	entries map[string]map[string]Entry
	size    int
	m       sync.Mutex
}

// NewCache returns a new instance of Cache.
func NewCache(o ...Option) Cache {
	opts := newOptions(o...)

	return Cache{
		size:    opts.size,
		entries: map[string]map[string]Entry{},
	}
}

// Get gets an entry on a service `svcKey` by a give `key`.
func (c *Cache) Get(svcKey, key string) (*Entry, error) {
	var value Entry
	ok := true

	c.m.Lock()
	defer c.m.Unlock()

	if value, ok = c.entries[svcKey][key]; !ok {
		return nil, fmt.Errorf("invalid service key: `%v`", key)
	}

	return &value, nil
}

// Set sets a key / value. It lets a service add entries on a request basis.
func (c *Cache) Set(svcKey, key string, val interface{}) error {
	c.m.Lock()
	defer c.m.Unlock()

	if !c.fits() {
		return fmt.Errorf("cache is full")
	}

	if _, ok := c.entries[svcKey]; !ok {
		c.entries[svcKey] = map[string]Entry{}
	}

	if _, ok := c.entries[svcKey][key]; ok {
		return fmt.Errorf("key `%v` already exists", key)
	}

	c.entries[svcKey][key] = Entry{
		V:     val,
		Valid: true,
	}

	return nil
}

// Invalidate invalidates a cache Entry by key.
func (c *Cache) Invalidate(svcKey, key string) error {
	r, err := c.Get(svcKey, key)
	if err != nil {
		return err
	}

	r.Valid = false
	c.entries[svcKey][key] = *r
	return nil
}

// Evict frees memory from the cache by removing invalid keys. It is a noop.
func (c *Cache) Evict() {
	for _, v := range c.entries {
		for k, svcEntry := range v {
			if !svcEntry.Valid {
				delete(v, k)
			}
		}
	}
}

// Length returns the amount of entries per service key.
func (c *Cache) Length(k string) int {
	return len(c.entries[k])
}

func (c *Cache) fits() bool {
	return c.size >= len(c.entries)
}
