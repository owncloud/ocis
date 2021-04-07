package storage

import (
	"strings"
)

// NewInMemoryStorage creates a new InMemory instance.
func NewInMemoryStorage() InMemory {
	return InMemory{
		store: make(map[string][]byte),
	}
}

// InMemory represents an in memory storage for thumbnails
// Can be used during development
type InMemory struct {
	store map[string][]byte
}

// Get loads the thumbnail from memory.
func (s InMemory) Get(key string) ([]byte, bool) {
	return s.store[key], true
}

// Set stores the thumbnail in memory.
func (s InMemory) Put(key string, thumbnail []byte) error {
	s.store[key] = thumbnail
	return nil
}

// BuildKey generates a unique key to store and retrieve the thumbnail.
func (s InMemory) BuildKey(r Request) string {
	parts := []string{
		r.Checksum,
		r.Resolution.String(),
		strings.Join(r.Types, ","),
	}
	return strings.Join(parts, "+")
}
