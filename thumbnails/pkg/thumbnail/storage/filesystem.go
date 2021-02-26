package storage

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/owncloud/ocis/ocis-pkg/log"
	"github.com/owncloud/ocis/thumbnails/pkg/config"
	"github.com/pkg/errors"
)

const (
	usersDir = "users"
	filesDir = "files"
)

// NewFileSystemStorage creates a new instanz of FileSystem
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
	mux    sync.Mutex
}

// Get loads the image from the file system.
func (s *FileSystem) Get(username string, key string) []byte {
	userDir := s.userDir(username)
	img := filepath.Join(userDir, key)
	content, err := ioutil.ReadFile(img)
	if err != nil {
		s.logger.Debug().Str("err", err.Error()).Str("key", key).Msg("could not load thumbnail from store")
		return nil
	}
	return content
}

// Set writes the image to the file system.
func (s *FileSystem) Set(username string, key string, img []byte) error {
	_, err := s.storeImage(key, img)
	if err != nil {
		return errors.Wrap(err, "could not store image")
	}
	userDir, err := s.createUserDir(username)
	if err != nil {
		return err
	}
	return s.linkImageToUserDir(key, userDir)
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

func (s *FileSystem) storeImage(key string, img []byte) (string, error) {
	s.mux.Lock()
	defer s.mux.Unlock()
	imgPath := filepath.Join(s.root, filesDir, key)
	dir := filepath.Dir(imgPath)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return "", errors.Wrapf(err, "error while creating directory %s", dir)
	}

	if _, err := os.Stat(imgPath); os.IsNotExist(err) {
		f, err := os.Create(imgPath)
		if err != nil {
			return "", errors.Wrapf(err, "could not create file \"%s\"", key)
		}
		defer f.Close()

		_, err = f.Write(img)
		if err != nil {
			return "", errors.Wrapf(err, "could not write to file \"%s\"", key)
		}
	}

	return imgPath, nil
}

// userDir returns the path to the user directory.
// The username is hashed before appending it on the path to prevent bugs caused by invalid folder names.
// Also the hash is then splitted up in three parts that results in a path which looks as follows:
// <filestorage-root>/users/<3 characters>/<3 characters>/<48 characters>/
// This will balance the folders in setups with many users.
func (s *FileSystem) userDir(username string) string {
	mh := md5.New()
	mh.Write([]byte("something"))

	hash := sha256.New224()
	if _, err := hash.Write([]byte(username)); err != nil {
		s.logger.Fatal().Err(err).Msg("failed to create hash")
	}
	unHash := hex.EncodeToString(hash.Sum(nil)) // 224 Bits or 224 / 4 = 56 characters.

	return filepath.Join(s.root, usersDir, unHash[:3], unHash[3:6], unHash[6:])
}

func (s *FileSystem) createUserDir(username string) (string, error) {
	userDir := s.userDir(username)
	if err := os.MkdirAll(userDir, 0700); err != nil {
		return "", errors.Wrapf(err, "could not create userDir: %s", userDir)
	}

	return userDir, nil
}

// linkImageToUserDir links the stored images to the user directory.
// The goal is to minimize disk usage by linking to the images if they already exist and avoid file duplicaiton.
func (s *FileSystem) linkImageToUserDir(key string, userDir string) error {
	imgRootDir := s.rootDir(key)

	s.mux.Lock()
	defer s.mux.Unlock()
	err := os.Symlink(filepath.Join(s.root, filesDir, imgRootDir), filepath.Join(userDir, imgRootDir))
	if err != nil {
		if !os.IsExist(err) {
			return errors.Wrap(err, "could not link image to userdir")
		}
	}
	return nil
}
