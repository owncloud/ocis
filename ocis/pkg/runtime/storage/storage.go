package storage

import "github.com/owncloud/ocis/ocis/pkg/runtime/process"

// Entries is a tuple of <extension:pid>
type Entries map[string]int

// Storage defines a basic persistence interface layer.
type Storage interface {
	// Store a representation of a process.
	Store(e process.ProcEntry) error

	// Delete a representation of a process.
	Delete(e process.ProcEntry) error

	// Load a single entry.
	Load(name string) int

	// LoadAll retrieves a set of entries of running processes on the host machine.
	LoadAll() Entries
}
