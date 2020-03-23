package storage

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/owncloud/ocis-pkg/v2/log"
	"github.com/owncloud/ocis-thumbnails/pkg/config"
)

// NewFileSystemStorage creates a new instanz of FileSystem
func NewFileSystemStorage(cfg config.FileSystemStorage, logger log.Logger) FileSystem {
	return FileSystem{
		dir:    cfg.RootDirectory,
		logger: logger,
	}
}

// FileSystem represents a storage for the thumbnails using the local file system.
type FileSystem struct {
	dir    string
	logger log.Logger
}

// Get loads the image from the file system.
func (s FileSystem) Get(key string) []byte {
	content, err := ioutil.ReadFile(filepath.Join(s.dir, key))
	if err != nil {
		s.logger.Warn().Err(err).Msgf("could not read file %s", key)
		return nil
	}

	return content
}

// Set writes the image to the file system.
func (s FileSystem) Set(key string, img []byte) error {
	path := filepath.Join(s.dir, key)
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("error while creating directory %s", dir)
	}

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("could not create file \"%s\" error: %s", key, err.Error())
	}
	defer f.Close()
	_, err = f.Write(img)
	if err != nil {
		return fmt.Errorf("could not write to file \"%s\" error: %s", key, err.Error())
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
func (s FileSystem) BuildKey(ctx Context) string {
	etag := ctx.ETag
	filetype := ctx.Types[0]
	filename := ctx.Resolution.String() + "." + filetype

	key := new(bytes.Buffer)
	key.WriteString(etag[:2])
	key.WriteRune('/')
	key.WriteString(etag[2:4])
	key.WriteRune('/')
	key.WriteString(etag[4:])
	key.WriteRune('/')
	key.WriteString(filename)

	return key.String()
}
