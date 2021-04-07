package storage

import (
	"image"
)

// Request combines different attributes needed for storage operations.
type Request struct {
	Checksum   string
	Types      []string
	Resolution image.Rectangle
}

// Storage defines the interface for a thumbnail store.
type Storage interface {
	Get(string) ([]byte, bool)
	Put(string, []byte) error
	BuildKey(Request) string
}
