package imgsource

import (
	"context"
	"fmt"
	"image"
	"os"
	"path/filepath"

	"github.com/owncloud/ocis-thumbnails/pkg/config"
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
		return nil, fmt.Errorf("failed to load the file %s from %s error %s", file, imgPath, err.Error())
	}

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}

	return img, nil
}
