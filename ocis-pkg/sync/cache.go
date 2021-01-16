package sync

import (
	"sync"
	"time"
)

// Cache is a barebones cache implementation.
type Cache struct {
	entries  sync.Map
	pool     sync.Pool
	capacity int
	length   int
}

// CacheEntry represents an entry on the cache. You can type assert on V.
type CacheEntry struct {
	V          interface{}
	expiration time.Time
}

// NewCache returns a new instance of Cache.
func NewCache(capacity int) Cache {
	return Cache{
		capacity: capacity,
		pool: sync.Pool{New: func() interface{} {
			return new(CacheEntry)
		}},
	}
}

// Get gets an entry by given key
func (c *Cache) Get(key string) *CacheEntry {
	if mapEntry, ok := c.entries.Load(key); ok {
		entry := mapEntry.(*CacheEntry)
		if c.expired(entry) {
			c.entries.Delete(key)
			return nil
		}
		return entry
	}
	return nil
}

// Set adds an entry for given key and value
func (c *Cache) Set(key string, val interface{}, expiration time.Time) {
	if !c.fits() {
		c.evict()
	}

	poolEntry := c.pool.Get()
	if mapEntry, loaded := c.entries.LoadOrStore(key, poolEntry); loaded {
		entry := mapEntry.(*CacheEntry)
		entry.V = val
		entry.expiration = expiration

		c.pool.Put(poolEntry)
	} else {
		entry := poolEntry.(*CacheEntry)
		entry.V = val
		entry.expiration = expiration

		c.length++
	}
}

// Unset removes an entry by given key
func (c *Cache) Unset(key string) bool {
	if _, loaded := c.entries.LoadAndDelete(key); !loaded {
		return false
	}

	c.length--
	return true
}

// evict frees memory from the cache by removing entries that exceeded the cache TTL.
func (c *Cache) evict() {
	c.entries.Range(func(key, mapEntry interface{}) bool {
		entry := mapEntry.(*CacheEntry)
		if c.expired(entry) {
			c.entries.Delete(key)
			c.length--
		}
		return true
	})
}

// expired checks if an entry is expired
func (c *Cache) expired(e *CacheEntry) bool {
	return e.expiration.Before(time.Now())
}

// fits returns whether the cache fits more entries.
func (c *Cache) fits() bool {
	return c.capacity > c.length
}
