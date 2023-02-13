package assetsfs

import (
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"

	"github.com/owncloud/ocis/v2/ocis-pkg/log"
)

// FileSystem customized to load assets
type FileSystem struct {
	fs        http.FileSystem
	assetPath string
	log       log.Logger
}

// Open checks if assetPath is set and tries to load from there. Falls back to fs if that is not possible
func (f *FileSystem) Open(original string) (http.File, error) {
	if f.assetPath != "" {
		file, err := os.Open(filepath.Join(f.assetPath, original))
		if err == nil {
			return file, nil
		}
	}
	return f.fs.Open(original)
}

func (f *FileSystem) OpenEmbedded(name string) (http.File, error) {
	return f.fs.Open(name)
}

// Create creates a new file in the assetPath
func (f *FileSystem) Create(name string) (*os.File, error) {
	fullPath := f.jailPath(name)
	if err := os.MkdirAll(filepath.Dir(fullPath), 0770); err != nil {
		return nil, err
	}
	return os.Create(fullPath)
}

// jailPath returns the fullPath `<assetPath>/<name>`. It makes sure that the path is
// always under `<assetPath>` to prevent directory traversal.
func (f *FileSystem) jailPath(name string) string {
	return filepath.Join(f.assetPath, filepath.Join("/", name))
}

// New initializes a new FileSystem. Quits on error
func New(embedFS fs.FS, assetPath string, logger log.Logger) *FileSystem {
	f, err := fs.Sub(embedFS, "assets")
	if err != nil {
		fmt.Println("Cannot load subtree fs:", err.Error())
		os.Exit(1)
	}

	return &FileSystem{
		fs:        http.FS(f),
		assetPath: assetPath,
		log:       logger,
	}
}
