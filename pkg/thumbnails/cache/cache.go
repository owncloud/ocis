package cache

import (
	"image"
)

// Cache defines the interface for a thumbnail cache.
type Cache interface {
	Get(key string) image.Image
	Set(key string, thumbnail image.Image) (image.Image, error)
}
