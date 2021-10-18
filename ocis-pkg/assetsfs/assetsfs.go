package assetsfs

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path"

	"github.com/owncloud/ocis/ocis-pkg/log"
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
		file, err := read(f.assetPath, original)
		if err == nil {
			return file, nil
		}
		f.log.Warn().
			Str("path", f.assetPath).
			Str("filename", original).
			Str("error", err.Error()).
			Msg("error reading from assetPath")
	}

	return f.fs.Open(original)
}

// New initializes a new FileSystem. Quits on error
func New(embedFS embed.FS, assetPath string, logger log.Logger) *FileSystem {
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

// tries to read file from disk or errors
func read(assetPath string, fileName string) (http.File, error) {
	if stat, err := os.Stat(assetPath); err != nil || !stat.IsDir() {
		return nil, fmt.Errorf("can't load asset path: %s", err)
	}

	p := path.Join(assetPath, fileName)
	if _, err := os.Stat(p); err != nil {
		return nil, err
	}

	f, err := os.Open(p)
	if err != nil {
		return nil, err
	}

	return f, nil
}
