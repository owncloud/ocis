package thumbnails

import (
	"bytes"
	"image"

	"github.com/nfnt/resize"
	"github.com/owncloud/ocis-pkg/v2/log"
	"github.com/owncloud/ocis-thumbnails/pkg/thumbnails/resolutions"
	"github.com/owncloud/ocis-thumbnails/pkg/thumbnails/storage"
)

// Context bundles information needed to generate a thumbnail for afile
type Context struct {
	Resolution resolutions.Resolution
	ImagePath  string
	Encoder    Encoder
	ETag       string
}

// Manager is responsible for generating thumbnails
type Manager interface {
	// Get will return a thumbnail for a file
	Get(Context, image.Image) ([]byte, error)
	// GetStored loads the thumbnail from the storage.
	// It will return nil if no image is stored for the given context.
	GetStored(Context) []byte
}

// NewSimpleManager creates a new instance of SimpleManager
func NewSimpleManager(storage storage.Storage, logger log.Logger) SimpleManager {
	return SimpleManager{
		storage: storage,
		logger:  logger,
	}
}

// SimpleManager is a simple implementation of Manager
type SimpleManager struct {
	storage storage.Storage
	logger  log.Logger
}

// Get implements the Get Method of Manager
func (s SimpleManager) Get(ctx Context, img image.Image) ([]byte, error) {
	thumbnail := s.generate(ctx, img)

	key := s.storage.BuildKey(mapToStorageContext(ctx))

	buf := new(bytes.Buffer)
	err := ctx.Encoder.Encode(buf, thumbnail)
	if err != nil {
		return nil, err
	}
	bytes := buf.Bytes()
	err = s.storage.Set(key, bytes)
	if err != nil {
		s.logger.Warn().Err(err).Msg("could not store thumbnail")
	}
	return bytes, nil
}

// GetStored tries to get the stored thumbnail and return it.
// If there is no cached thumbnail it will return nil
func (s SimpleManager) GetStored(ctx Context) []byte {
	key := s.storage.BuildKey(mapToStorageContext(ctx))
	stored := s.storage.Get(key)
	return stored
}

func (s SimpleManager) generate(ctx Context, img image.Image) image.Image {
	thumbnail := resize.Thumbnail(uint(ctx.Resolution.Width), uint(ctx.Resolution.Height), img, resize.Lanczos2)
	return thumbnail
}

func mapToStorageContext(ctx Context) storage.Context {
	sCtx := storage.Context{
		ETag:       ctx.ETag,
		Resolution: ctx.Resolution,
		Types:      ctx.Encoder.Types(),
	}
	return sCtx
}
