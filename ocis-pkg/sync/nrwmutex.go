package sync

import (
	"sync"
)

// NRWMutex works the same as RWMutex, the only difference is that it stores mutexes in a map and reuses them.
// It's handy if you want to write-lock, write-unlock, read-lock and read-unlock for specific names only.
type NRWMutex struct {
	pool sync.Pool
	mus  sync.Map
}

// NewNRWMutex returns a new instance of NRWMutex.
func NewNRWMutex() NRWMutex {
	return NRWMutex{pool: sync.Pool{New: func() interface{} {
		return new(sync.RWMutex)
	}}}
}

// Lock locks rw for writing.
func (m *NRWMutex) Lock(name string) {
	m.loadOrStore(name).Lock()
}

// Unlock unlocks rw for writing.
func (m *NRWMutex) Unlock(name string) {
	m.loadOrStore(name).Unlock()
}

// RLock locks rw for reading.
func (m *NRWMutex) RLock(name string) {
	m.loadOrStore(name).RLock()
}

// RUnlock undoes a single RLock call.
func (m *NRWMutex) RUnlock(name string) {
	m.loadOrStore(name).RUnlock()
}

func (m *NRWMutex) loadOrStore(name string) *sync.RWMutex {
	pmu := m.pool.Get()
	mmu, loaded := m.mus.LoadOrStore(name, pmu)
	if loaded {
		m.pool.Put(pmu)
	}

	return mmu.(*sync.RWMutex)
}
