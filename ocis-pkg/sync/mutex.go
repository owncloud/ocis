package sync

import (
	"sync"
)

// NRWMutex works the same as RWMutex, the only difference is that it stores mutexes in a map and reuses them.
// It's handy if you want to write-lock, write-unlock, read-lock and read-unlock for specific names only.
type NRWMutex struct {
	m  sync.Mutex
	mm map[string]*nrw
}

type nrw struct {
	m sync.RWMutex
	c int
}

// NewNRWMutex returns a new instance of NRWMutex.
func NewNRWMutex() NRWMutex {
	return NRWMutex{mm: make(map[string]*nrw)}
}

// Lock locks rw for writing.
func (c *NRWMutex) Lock(k string) {
	c.m.Lock()
	m := c.get(k)
	m.c++
	c.m.Unlock()
	m.m.Lock()
}

// Unlock unlocks rw for writing.
func (c *NRWMutex) Unlock(k string) {
	c.m.Lock()
	defer c.m.Unlock()
	m := c.get(k)
	m.m.Unlock()
	m.c--
	if m.c == 0 {
		delete(c.mm, k)
	}
}

// RLock locks rw for reading.
func (c *NRWMutex) RLock(k string) {
	c.m.Lock()
	m := c.get(k)
	m.c++
	c.m.Unlock()
	m.m.RLock()
}

// RUnlock undoes a single RLock call.
func (c *NRWMutex) RUnlock(k string) {
	c.m.Lock()
	defer c.m.Unlock()
	m := c.get(k)
	m.m.RUnlock()
	m.c--
	if m.c == 0 {
		delete(c.mm, k)
	}
}

func (c *NRWMutex) get(k string) *nrw {
	m, ok := c.mm[k]
	if !ok {
		m = &nrw{}
		c.mm[k] = m
	}

	return m
}
