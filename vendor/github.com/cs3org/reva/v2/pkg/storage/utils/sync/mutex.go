// Copyright 2018-2022 CERN
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// In applying this license, CERN does not waive the privileges and immunities
// granted to it by virtue of its status as an Intergovernmental Organization
// or submit itself to any jurisdiction.

package sync

import (
	"sync"
)

// NamedRWMutex works the same as RWMutex, the only difference is that it stores mutexes in a map and reuses them.
// It's handy if you want to write-lock, write-unlock, read-lock and read-unlock for specific names only.
type NamedRWMutex struct {
	pool sync.Pool
	mus  sync.Map
}

// NewNamedRWMutex returns a new instance of NamedRWMutex.
func NewNamedRWMutex() NamedRWMutex {
	return NamedRWMutex{pool: sync.Pool{New: func() interface{} {
		return new(sync.RWMutex)
	}}}
}

// Lock locks rw for writing.
func (m *NamedRWMutex) Lock(name string) {
	m.loadOrStore(name).Lock()
}

// Unlock unlocks rw for writing.
func (m *NamedRWMutex) Unlock(name string) {
	m.loadOrStore(name).Unlock()
}

// RLock locks rw for reading.
func (m *NamedRWMutex) RLock(name string) {
	m.loadOrStore(name).RLock()
}

// RUnlock undoes a single RLock call.
func (m *NamedRWMutex) RUnlock(name string) {
	m.loadOrStore(name).RUnlock()
}

func (m *NamedRWMutex) loadOrStore(name string) *sync.RWMutex {
	pmu := m.pool.Get()
	mmu, loaded := m.mus.LoadOrStore(name, pmu)
	if loaded {
		m.pool.Put(pmu)
	}

	return mmu.(*sync.RWMutex)
}
