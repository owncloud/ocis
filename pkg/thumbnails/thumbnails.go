package thumbnails

import (
	"bytes"
	"image"
	"time"

	"github.com/nfnt/resize"
	"github.com/owncloud/ocis-thumbnails/pkg/thumbnails/storage"
)

// ThumbnailContext bundles information needed to generate a thumbnail for afile
type ThumbnailContext struct {
	Width     int
	Height    int
	ImagePath string
	Encoder   Encoder
	ETag      string
}

// Manager is responsible for generating thumbnails
type Manager interface {
	// Get will return a thumbnail for a file
	Get(ThumbnailContext, image.Image) ([]byte, error)
	GetStored(ThumbnailContext) []byte
}

// SimpleManager is a simple implementation of Manager
type SimpleManager struct {
	Storage storage.Storage
}

// Get implements the Get Method of Manager
func (s SimpleManager) Get(ctx ThumbnailContext, img image.Image) ([]byte, error) {
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
func (s SimpleManager) GetStored(ctx ThumbnailContext) []byte {
	key := s.Storage.BuildKey(mapToStorageContext(ctx))
	stored := s.Storage.Get(key)
	if stored == nil {
		return nil
	}
	return stored
}

func (s SimpleManager) generate(ctx ThumbnailContext, img image.Image) image.Image {
	// TODO: remove, just for demo purposes
	time.Sleep(time.Second * 2)

	thumbnail := resize.Thumbnail(uint(ctx.Width), uint(ctx.Height), img, resize.Lanczos2)
	return thumbnail
}

func mapToStorageContext(ctx ThumbnailContext) storage.StorageContext {
	sCtx := storage.StorageContext{
		ETag:   ctx.ETag,
		Width:  ctx.Width,
		Height: ctx.Height,
		Types:  ctx.Encoder.Types(),
	}
	return sCtx
}
