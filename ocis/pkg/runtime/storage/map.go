package storage

import (
	"github.com/owncloud/ocis/ocis/pkg/runtime/process"

	"sync"
)

// Map synchronizes access to extension+pid tuples.
type Map struct {
	c *sync.Map
}

// NewMapStorage initializes a new Storage.
func NewMapStorage() Storage {
	return &Map{
		c: &sync.Map{},
	}
}

// Store a value on the underlying data structure.
func (m *Map) Store(e process.ProcEntry) error {
	m.c.Store(e.Extension, e.Pid)
	return nil
}

// Delete a value on the underlying data structure.
func (m *Map) Delete(e process.ProcEntry) error {
	m.c.Delete(e.Extension)
	return nil
}

// Load a single pid.
func (m *Map) Load(name string) int {
	var val int
	m.c.Range(func(k, v interface{}) bool {
		if k.(string) == name {
			val = v.(int)
			return false
		}
		return true
	})
	return val
}

// LoadAll values from the underlying data structure.
func (m *Map) LoadAll() Entries {
	e := make(map[string]int)
	m.c.Range(func(k, v interface{}) bool {
		ks, ok := k.(string)
		if !ok {
			return false
		}

		vs, ok := v.(int)
		if !ok {
			return false
		}

		e[ks] = vs
		return true
	})
	return e
}
