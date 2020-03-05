package storage

import (
	"image"
)

// Storage defines the interface for a thumbnail store.
type Storage interface {
	Get(key string) image.Image
	Set(key string, thumbnail image.Image) (image.Image, error)
}
