package storage

import (
	"strings"
)

// NewInMemoryStorage creates a new InMemory instance.
func NewInMemoryStorage() InMemory {
	return InMemory{
		store: make(map[string]map[string][]byte),
	}
}

// InMemory represents an in memory storage for thumbnails
// Can be used during development
type InMemory struct {
	store map[string]map[string][]byte
}

// Get loads the thumbnail from memory.
func (s InMemory) Get(username string, key string) []byte {
	userImages := s.store[username]
	if userImages == nil {
		return nil
	}
	return s.store[username][key]
}

// Set stores the thumbnail in memory.
func (s InMemory) Set(username string, key string, thumbnail []byte) error {
	if _, ok := s.store[username]; !ok {
		s.store[username] = make(map[string][]byte)
	}
	s.store[username][key] = thumbnail
	return nil
}

// BuildKey generates a unique key to store and retrieve the thumbnail.
func (s InMemory) BuildKey(r Request) string {
	parts := []string{
		r.ETag,
		r.Resolution.String(),
		strings.Join(r.Types, ","),
	}
	return strings.Join(parts, "+")
}
