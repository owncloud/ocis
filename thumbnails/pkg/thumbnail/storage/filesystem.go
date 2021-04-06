package storage

import (
	"github.com/owncloud/ocis/ocis-pkg/log"
	"github.com/owncloud/ocis/thumbnails/pkg/config"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	filesDir = "files"
)

// NewFileSystemStorage creates a new instance of FileSystem
func NewFileSystemStorage(cfg config.FileSystemStorage, logger log.Logger) *FileSystem {
	return &FileSystem{
		root:   cfg.RootDirectory,
		logger: logger,
	}
}

// FileSystem represents a storage for the thumbnails using the local file system.
type FileSystem struct {
	root   string
	logger log.Logger
}

// Get loads the image from the file system.
func (s *FileSystem) Get(key string) ([]byte, bool) {
	img := filepath.Join(s.root, filesDir, key)
	content, err := os.ReadFile(img)
	if err != nil {
		if !os.IsNotExist(err) {
			s.logger.Debug().Str("err", err.Error()).Str("key", key).Msg("could not load thumbnail from store")
		}
		return nil, false
	}
	return content, true
}

// Set writes the image to the file system.
func (s *FileSystem) Put(key string, img []byte) error {
	imgPath := filepath.Join(s.root, filesDir, key)
	dir := filepath.Dir(imgPath)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return errors.Wrapf(err, "error while creating directory %s", dir)
	}

	if _, err := os.Stat(imgPath); os.IsNotExist(err) {
		f, err := os.Create(imgPath)
		if err != nil {
			return errors.Wrapf(err, "could not create file \"%s\"", key)
		}
		defer f.Close()

		if _, err = f.Write(img); err != nil {
			return errors.Wrapf(err, "could not write to file \"%s\"", key)
		}
	}

	return nil
}

// BuildKey generate the unique key for a thumbnail.
// The key is structure as follows:
//
// <first two letters of etag>/<next two letters of etag>/<rest of etag>/<width>x<height>.<filetype>
//
// e.g. 97/9f/4c8db98f7b82e768ef478d3c8612/500x300.png
//
// The key also represents the path to the thumbnail in the filesystem under the configured root directory.
func (s *FileSystem) BuildKey(r Request) string {
	etag := r.ETag
	filetype := r.Types[0]
	filename := strconv.Itoa(r.Resolution.Dx()) + "x" + strconv.Itoa(r.Resolution.Dy()) + "." + filetype

	return filepath.Join(etag[:2], etag[2:4], etag[4:], filename)
}

func (s *FileSystem) rootDir(key string) string {
	p := strings.Split(key, string(os.PathSeparator))
	return p[0]
}
