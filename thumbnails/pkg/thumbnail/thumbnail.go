package thumbnail

import (
	"bytes"
	"github.com/owncloud/ocis/ocis-pkg/log"
	"github.com/owncloud/ocis/thumbnails/pkg/thumbnail/storage"
	"golang.org/x/image/draw"
	"image"
)

// Request bundles information needed to generate a thumbnail for afile
type Request struct {
	Resolution image.Rectangle
	Encoder    Encoder
	ETag       string
}

// Manager is responsible for generating thumbnails
type Manager interface {
	// Get will return a thumbnail for a file
	Generate(Request, image.Image) ([]byte, error)
	// GetStored loads the thumbnail from the storage.
	// It will return nil if no image is stored for the given context.
	Get(Request) ([]byte, bool)
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
func (s SimpleManager) Generate(r Request, img image.Image) ([]byte, error) {
	match := s.resolutions.ClosestMatch(r.Resolution, img.Bounds())
	thumbnail := s.generate(match, img)

	dst := new(bytes.Buffer)
	err := r.Encoder.Encode(dst, thumbnail)
	if err != nil {
		return nil, err
	}

	key := s.storage.BuildKey(mapToStorageRequest(r))
	err = s.storage.Put(key, dst.Bytes())
	if err != nil {
		s.logger.Warn().Err(err).Msg("could not store thumbnail")
	}
	return dst.Bytes(), nil
}

// GetStored tries to get the stored thumbnail and return it.
// If there is no cached thumbnail it will return nil
func (s SimpleManager) Get(r Request) ([]byte, bool) {
	key := s.storage.BuildKey(mapToStorageRequest(r))
	return s.storage.Get(key)
}

func (s SimpleManager) generate(r image.Rectangle, img image.Image) image.Image {
	targetResolution := mapRatio(img.Bounds(), r)
	thumbnail := image.NewRGBA(targetResolution)
	draw.ApproxBiLinear.Scale(thumbnail, targetResolution, img, img.Bounds(), draw.Over, nil)
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
