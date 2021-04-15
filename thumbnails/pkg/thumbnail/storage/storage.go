package storage

import (
	"image"
)

// Request combines different attributes needed for storage operations.
type Request struct {
	// The checksum of the source file
	// Will be used to determine if a thumbnail exists
	Checksum   string
	// Types provided by the encoder.
	// Contains the mimetypes of the thumbnail.
	// In case of jpg/jpeg it will contain both.
	Types      []string
	// The resolution of the thumbnail
	Resolution image.Rectangle
}

// Storage defines the interface for a thumbnail store.
type Storage interface {
	Get(string) ([]byte, bool)
	Put(string, []byte) error
	BuildKey(Request) string
}
