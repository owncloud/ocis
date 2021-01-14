package sync

import (
	"sync"
)

type NRWMutex struct {
	m  sync.Mutex
	mm map[string]*nrw
}

type nrw struct {
	m sync.RWMutex
	c int
}

func NewNRWMutex() NRWMutex {
	return NRWMutex{mm: make(map[string]*nrw)}
}

func (c *NRWMutex) Lock(k string) {
	c.m.Lock()
	m := c.get(k)
	m.c++
	c.m.Unlock()
	m.m.Lock()
}

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

func (c *NRWMutex) RLock(k string) {
	c.m.Lock()
	m := c.get(k)
	m.c++
	c.m.Unlock()
	m.m.RLock()
}

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
