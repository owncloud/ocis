package storage

import (
	"strings"
)

func NewInMemoryStorage() InMemory {
	return InMemory{
		store: make(map[string][]byte),
	}
}

type InMemory struct {
	store map[string][]byte
}

func (s InMemory) Get(key string) []byte {
	return s.store[key]
}

func (s InMemory) Set(key string, thumbnail []byte) error {
	s.store[key] = thumbnail
	return nil
}

func (s InMemory) BuildKey(ctx StorageContext) string {
	parts := []string{
		ctx.ETag,
		string(ctx.Width) + "x" + string(ctx.Height),
		strings.Join(ctx.Types, ","),
	}
	return strings.Join(parts, "+")
}
