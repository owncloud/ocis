package cache

import (
	"sync"
	"time"
)

// Entry represents an entry on the cache. You can type assert on V.
type Entry struct {
	V          interface{}
	expiration time.Time
}

// Cache is a barebones cache implementation.
type Cache struct {
	entries map[string]*Entry
	size    int
	m       sync.Mutex
}

// NewCache returns a new instance of Cache.
func NewCache(o ...Option) Cache {
	opts := newOptions(o...)

	return Cache{
		size:    opts.size,
		entries: map[string]*Entry{},
	}
}

// Get gets an entry by given key
func (c *Cache) Get(k string) *Entry {
	c.m.Lock()
	defer c.m.Unlock()

	if _, ok := c.entries[k]; ok {
		if c.expired(c.entries[k]) {
			delete(c.entries, k)
			return nil
		}
		return c.entries[k]
	}
	return nil
}

// Set adds an entry for given key and value
func (c *Cache) Set(k string, val interface{}, expiration time.Time) {
	c.m.Lock()
	defer c.m.Unlock()

	if !c.fits() {
		c.evict()
	}

	c.entries[k] = &Entry{
		val,
		expiration,
	}
}

// Unset removes an entry by given key
func (c *Cache) Unset(k string) bool {
	if _, ok := c.entries[k]; !ok {
		return false
	}

	delete(c.entries, k)
	return true
}

// evict frees memory from the cache by removing entries that exceeded the cache TTL.
func (c *Cache) evict() {
	for i := range c.entries {
		if c.expired(c.entries[i]) {
			delete(c.entries, i)
		}
	}
}

// expired checks if an entry is expired
func (c *Cache) expired(e *Entry) bool {
	return e.expiration.Before(time.Now())
}

// fits returns whether the cache fits more entries.
func (c *Cache) fits() bool {
	return c.size > len(c.entries)
}
