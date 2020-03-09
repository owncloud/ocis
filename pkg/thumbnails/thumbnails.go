package thumbnails

import (
	"bytes"
	"image"

	"github.com/nfnt/resize"
	"github.com/owncloud/ocis-thumbnails/pkg/thumbnails/storage"
)

// Context bundles information needed to generate a thumbnail for afile
type Context struct {
	Width     int
	Height    int
	ImagePath string
	Encoder   Encoder
	ETag      string
}

// Manager is responsible for generating thumbnails
type Manager interface {
	// Get will return a thumbnail for a file
	Get(Context, image.Image) ([]byte, error)
	// GetStored loads the thumbnail from the storage.
	// It will return nil if no image is stored for the given context.
	GetStored(Context) []byte
}

// SimpleManager is a simple implementation of Manager
type SimpleManager struct {
	Storage storage.Storage
}

// Get implements the Get Method of Manager
func (s SimpleManager) Get(ctx Context, img image.Image) ([]byte, error) {
	thumbnail := s.generate(ctx, img)

	key := s.Storage.BuildKey(mapToStorageContext(ctx))

	buf := new(bytes.Buffer)
	err := ctx.Encoder.Encode(buf, thumbnail)
	if err != nil {
		return nil, err
	}
	bytes := buf.Bytes()
	s.Storage.Set(key, bytes)
	return bytes, nil
}

// GetStored tries to get the stored thumbnail and return it.
// If there is no cached thumbnail it will return nil
func (s SimpleManager) GetStored(ctx Context) []byte {
	key := s.Storage.BuildKey(mapToStorageContext(ctx))
	stored := s.Storage.Get(key)
	return stored
}

func (s SimpleManager) generate(ctx Context, img image.Image) image.Image {
	thumbnail := resize.Thumbnail(uint(ctx.Width), uint(ctx.Height), img, resize.Lanczos2)
	return thumbnail
}

func mapToStorageContext(ctx Context) storage.Context {
	sCtx := storage.Context{
		ETag:   ctx.ETag,
		Width:  ctx.Width,
		Height: ctx.Height,
		Types:  ctx.Encoder.Types(),
	}
	return sCtx
}
