package thumbnail

import (
	"bytes"
	"image"

	"github.com/disintegration/imaging"
	"github.com/owncloud/ocis/ocis-pkg/log"
	"github.com/owncloud/ocis/thumbnails/pkg/thumbnail/storage"
)

// Request bundles information needed to generate a thumbnail for afile
type Request struct {
	Resolution image.Rectangle
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
func NewSimpleManager(resolutions Resolutions, storage storage.Storage, logger log.Logger) SimpleManager {
	return SimpleManager{
		storage:     storage,
		logger:      logger,
		resolutions: resolutions,
	}
}

// SimpleManager is a simple implementation of Manager
type SimpleManager struct {
	storage     storage.Storage
	logger      log.Logger
	resolutions Resolutions
}

// Get implements the Get Method of Manager
func (s SimpleManager) Get(r Request, img image.Image) ([]byte, error) {
	//match := s.resolutions.ClosestMatch(r.Resolution, img.Bounds())
	thumbnail := s.generate(r.Resolution, img)

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

func (s SimpleManager) generate(r image.Rectangle, img image.Image) (thumbnail image.Image) {
	thumbnail = imaging.Fill(img, r.Dx(), r.Dy(), imaging.Center, imaging.Lanczos)
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
