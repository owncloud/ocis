package thumbnail

import (
	"bytes"
	"image"
	"image/gif"
	"mime"

	"github.com/owncloud/ocis/v2/extensions/thumbnails/pkg/thumbnail/storage"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
)

var (
	// SupportedMimeTypes contains a all mimetypes which are supported by the thumbnailer.
	SupportedMimeTypes = map[string]struct{}{
		"image/png":  {},
		"image/jpg":  {},
		"image/jpeg": {},
		"image/gif":  {},
		"text/plain": {},
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
	// Generate creates a thumbnail and stores it.
	// The function returns a key with which the actual file can be retrieved.
	Generate(Request, interface{}) (string, error)
	// CheckThumbnail checks if a thumbnail with the requested attributes exists.
	// The function will return a status if the file exists and the key to the file.
	CheckThumbnail(Request) (string, bool)
	// GetThumbnail will load the thumbnail from the storage and return its content.
	GetThumbnail(key string) ([]byte, error)
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

func (s SimpleManager) Generate(r Request, img interface{}) (string, error) {
	var match image.Rectangle
	switch m := img.(type) {
	case *gif.GIF:
		match = s.resolutions.ClosestMatch(r.Resolution, m.Image[0].Bounds())
	case image.Image:
		match = s.resolutions.ClosestMatch(r.Resolution, m.Bounds())
	}

	thumbnail, err := r.Generator.GenerateThumbnail(match, img)
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	if err := r.Encoder.Encode(buf, thumbnail); err != nil {
		return "", err
	}

	k := s.storage.BuildKey(mapToStorageRequest(r))
	if err := s.storage.Put(k, buf.Bytes()); err != nil {
		s.logger.Error().Err(err).Msg("could not store thumbnail")
		return "", err
	}
	return k, nil
}

func (s SimpleManager) CheckThumbnail(r Request) (string, bool) {
	k := s.storage.BuildKey(mapToStorageRequest(r))
	return k, s.storage.Stat(k)
}

func (s SimpleManager) GetThumbnail(key string) ([]byte, error) {
	return s.storage.Get(key)
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
	_, supported := SupportedMimeTypes[mimeType]
	return supported
}
