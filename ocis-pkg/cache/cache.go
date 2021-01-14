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
	sync.Map
	sizeTotal   int
	sizeCurrent int
}

// NewCache returns a new instance of Cache.
func NewCache(sizeTotal int) Cache {
	return Cache{
		sizeTotal: sizeTotal,
	}
}

// Get gets an entry by given key
func (c *Cache) Get(k string) *Entry {
	if sme, ok := c.Load(k); ok {
		e := sme.(*Entry)
		if c.expired(e) {
			c.Delete(k)
			return nil
		}
		return e
	}
	return nil
}

// Set adds an entry for given key and value
func (c *Cache) Set(k string, val interface{}, expiration time.Time) {
	if !c.fits() {
		c.evict()
	}
	c.Store(k, &Entry{
		val,
		expiration,
	})
	c.sizeCurrent++
}

// Unset removes an entry by given key
func (c *Cache) Unset(k string) bool {
	if _, ok := c.Load(k); !ok {
		return false
	}

	c.Delete(k)
	c.sizeCurrent--
	return true
}

// evict frees memory from the cache by removing entries that exceeded the cache TTL.
func (c *Cache) evict() {
	c.Range(func(k, sme interface{}) bool {
		e := sme.(*Entry)
		if c.expired(e) {
			c.Delete(k)
			c.sizeCurrent--
		}
		return true
	})
}

// expired checks if an entry is expired
func (c *Cache) expired(e *Entry) bool {
	return e.expiration.Before(time.Now())
}

// fits returns whether the cache fits more entries.
func (c *Cache) fits() bool {
	return c.sizeTotal > c.sizeCurrent
}
