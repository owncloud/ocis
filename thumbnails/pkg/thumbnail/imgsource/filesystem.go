package imgsource

import (
	"context"
	"image"
	"os"
	"path/filepath"

	"github.com/owncloud/ocis/thumbnails/pkg/config"
	"github.com/pkg/errors"
)

// NewFileSystemSource return a new FileSystem instance
func NewFileSystemSource(cfg config.FileSystemSource) FileSystem {
	return FileSystem{
		basePath: cfg.BasePath,
	}
}

// FileSystem is an image source using the local file system
type FileSystem struct {
	basePath string
}

// Get retrieves an image from the filesystem.
func (s FileSystem) Get(ctx context.Context, file string) (image.Image, error) {
	imgPath := filepath.Join(s.basePath, file)
	f, err := os.Open(filepath.Clean(imgPath))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to load the file %s from %s", file, imgPath)
	}

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, errors.Wrap(err, "Get: Decode:")
	}

	return img, nil
}
