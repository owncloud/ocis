package fsx

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/owncloud/ocis/v2/ocis-pkg/x/path/filepath"
)

// FallbackFS is a FileSystem that tries to load from assetPath first
// and falls back to fs if that is not possible
type FallbackFS struct {
	fsys    fs.FS
	sysPath string
}

// Open checks if assetPath is set and tries to load from there. Falls back to fs if that is not possible
func (f *FallbackFS) Open(name string) (fs.File, error) {
	if f.sysPath != "" {
		file, err := os.Open(filepathx.JailJoin(f.sysPath, name))
		if err == nil {
			return file, nil
		}
	}

	return f.fsys.Open(name)
}

// OpenEmbedded opens a file from the embedded filesystem only
func (f *FallbackFS) OpenEmbedded(name string) (fs.File, error) {
	return f.fsys.Open(name)
}

// Create creates a new file in the assetPath
func (f *FallbackFS) Create(name string) (*os.File, error) {
	fullPath := filepathx.JailJoin(f.sysPath, name)
	if err := os.MkdirAll(filepath.Dir(fullPath), 0770); err != nil {
		return nil, err
	}

	return os.Create(fullPath)
}

// NewFallbackFS return a new EmbeddedFallbackFS instance
func NewFallbackFS(fsys fs.FS, sysPath string) *FallbackFS {
	return &FallbackFS{
		fsys:    fsys,
		sysPath: sysPath,
	}
}
