package thumbnail

import (
	"bytes"
	"image"
	"image/gif"
	"mime"
	"strings"

	"github.com/owncloud/ocis/ocis-pkg/log"
	"github.com/owncloud/ocis/thumbnails/pkg/thumbnail/storage"
)

var (
	SupportedMimeTypes = [...]string{
		"image/png",
		"image/jpg",
		"image/jpeg",
		"image/gif",
		"text/plain",
	}
)

// Request bundles information needed to generate a thumbnail for afile
type Request struct {
	Resolution image.Rectangle
	Encoder    Encoder
	Generator  Generator
	Checksum   string
}

// Manager is responsible for generating thumbnails
type Manager interface {
	// Generate will return a thumbnail for a file
	Generate(Request, interface{}) ([]byte, error)
	// Get loads the thumbnail from the storage.
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

// Generate creates a thumbnail and stores it.
// The created thumbnail is also being returned.
func (s SimpleManager) Generate(r Request, img interface{}) ([]byte, error) {
	var match image.Rectangle
	switch m := img.(type) {
	case *gif.GIF:
		match = s.resolutions.ClosestMatch(r.Resolution, m.Image[0].Bounds())
	case image.Image:
		match = s.resolutions.ClosestMatch(r.Resolution, m.Bounds())
	}

	thumbnail, err := r.Generator.GenerateThumbnail(match, img)
	if err != nil {
		return nil, err
	}

	dst := new(bytes.Buffer)
	err = r.Encoder.Encode(dst, thumbnail)
	if err != nil {
		return nil, err
	}

	k := s.storage.BuildKey(mapToStorageRequest(r))
	err = s.storage.Put(k, dst.Bytes())
	if err != nil {
		s.logger.Warn().Err(err).Msg("could not store thumbnail")
	}
	return dst.Bytes(), nil
}

// Get tries to get the stored thumbnail and return it.
// If there is no cached thumbnail it will return nil
func (s SimpleManager) Get(r Request) ([]byte, bool) {
	k := s.storage.BuildKey(mapToStorageRequest(r))
	return s.storage.Get(k)
}

func mapToStorageRequest(r Request) storage.Request {
	return storage.Request{
		Checksum:   r.Checksum,
		Resolution: r.Resolution,
		Types:      r.Encoder.Types(),
	}
}

func IsMimeTypeSupported(m string) bool {
	mimeType, _, err := mime.ParseMediaType(m)
	if err != nil {
		return false
	}
	for _, mt := range SupportedMimeTypes {
		if strings.EqualFold(mt, mimeType) {
			return true
		}
	}
	return false
}
