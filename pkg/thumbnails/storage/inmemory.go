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
func (s InMemory) Get(key string) []byte {
	return s.store[key]
}

// Set stores the thumbnail in memory.
func (s InMemory) Set(key string, thumbnail []byte) error {
	s.store[key] = thumbnail
	return nil
}

// BuildKey generates a unique key to store and retrieve the thumbnail.
func (s InMemory) BuildKey(ctx Context) string {
	parts := []string{
		ctx.ETag,
		ctx.Resolution.String(),
		strings.Join(ctx.Types, ","),
	}
	return strings.Join(parts, "+")
}
