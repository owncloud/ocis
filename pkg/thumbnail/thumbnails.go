package thumbnail

import (
	"bytes"
	"image"

	"github.com/nfnt/resize"
	"github.com/owncloud/ocis-pkg/v2/log"
	"github.com/owncloud/ocis-thumbnails/pkg/thumbnail/resolution"
	"github.com/owncloud/ocis-thumbnails/pkg/thumbnail/storage"
)

// Request bundles information needed to generate a thumbnail for afile
type Request struct {
	Resolution resolution.Resolution
	ImagePath  string
	Encoder    Encoder
	ETag       string
	Username   string
}

// Manager is responsible for generating thumbnails
type Manager interface {
	// Get will return a thumbnail for a file
	Get(Request, image.Image) ([]byte, error)
	// GetStored loads the thumbnail from the storage.
	// It will return nil if no image is stored for the given context.
	GetStored(Request) []byte
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
func (s SimpleManager) Get(r Request, img image.Image) ([]byte, error) {
	thumbnail := s.generate(r, img)

	key := s.storage.BuildKey(mapToStorageRequest(r))

	buf := new(bytes.Buffer)
	err := r.Encoder.Encode(buf, thumbnail)
	if err != nil {
		return nil, err
	}
	bytes := buf.Bytes()
	err = s.storage.Set(r.Username, key, bytes)
	if err != nil {
		s.logger.Warn().Err(err).Msg("could not store thumbnail")
	}
	return bytes, nil
}

// GetStored tries to get the stored thumbnail and return it.
// If there is no cached thumbnail it will return nil
func (s SimpleManager) GetStored(r Request) []byte {
	key := s.storage.BuildKey(mapToStorageRequest(r))
	stored := s.storage.Get(r.Username, key)
	return stored
}

func (s SimpleManager) generate(r Request, img image.Image) image.Image {
	thumbnail := resize.Thumbnail(uint(r.Resolution.Width), uint(r.Resolution.Height), img, resize.Lanczos2)
	return thumbnail
}

func mapToStorageRequest(r Request) storage.Request {
	sR := storage.Request{
		ETag:       r.ETag,
		Resolution: r.Resolution,
		Types:      r.Encoder.Types(),
	}
	return sR
}
