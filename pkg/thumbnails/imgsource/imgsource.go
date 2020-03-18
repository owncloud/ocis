package imgsource

import (
	"context"
	"image"
)

// Source defines the interface for image sources
type Source interface {
	Get(ctx context.Context, path string) (image.Image, error)
}
