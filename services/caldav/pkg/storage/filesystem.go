package storage

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/rs/zerolog/log"

	"github.com/emersion/go-webdav"
	"github.com/emersion/go-webdav/caldav"
	"github.com/emersion/go-webdav/carddav"
)

type filesystemBackend struct {
	webdav.UserPrincipalBackend
	path          string
	caldavPrefix  string
	carddavPrefix string
}

var (
	validFilenameRegex  = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9_-]+(.[a-zA-Z]+)?$`)
	defaultResourceName = "default"
)

func NewFilesystem(fsPath, caldavPrefix, carddavPrefix string, userPrincipalBackend webdav.UserPrincipalBackend) (caldav.Backend, carddav.Backend, error) {
	info, err := os.Stat(fsPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create filesystem backend: %s", err.Error())
	}
	if !info.IsDir() {
		return nil, nil, fmt.Errorf("base path for filesystem backend must be a directory")
	}
	backend := &filesystemBackend{
		UserPrincipalBackend: userPrincipalBackend,
		path:                 fsPath,
		caldavPrefix:         caldavPrefix,
		carddavPrefix:        carddavPrefix,
	}
	return backend, backend, nil
}

func ensureLocalDir(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.MkdirAll(path, 0755)
		if err != nil {
			return fmt.Errorf("error creating '%s': %s", path, err.Error())
		}
	}
	return nil
}

func (b *filesystemBackend) localDir(homeSetPath string, components ...string) (string, error) {
	c := append([]string{b.path}, homeSetPath)
	c = append(c, components...)
	localPath := filepath.Join(c...)
	if err := ensureLocalDir(localPath); err != nil {
		return "", err
	}
	return localPath, nil
}

// don't use this directly, use localCalDAVPath or localCardDAVPath instead.
// note that homesetpath is expected to end in /
func (b *filesystemBackend) safeLocalPath(homeSetPath string, urlPath string) (string, error) {
	localPath := filepath.Join(b.path, homeSetPath)
	if err := ensureLocalDir(localPath); err != nil {
		return "", err
	}

	if urlPath == "" {
		return localPath, nil
	}

	// We are mapping to local filesystem path, so be conservative about what to accept
	if strings.HasSuffix(urlPath, "/") {
		urlPath = path.Clean(urlPath) + "/"
	} else {
		urlPath = path.Clean(urlPath)
	}
	if !strings.HasPrefix(urlPath, homeSetPath) {
		err := fmt.Errorf("access to resource outside of home set: %s", urlPath)
		return "", webdav.NewHTTPError(403, err)
	}
	urlPath = strings.TrimPrefix(urlPath, homeSetPath)

	// only accept simple file names for now
	dir, file := path.Split(urlPath)
	if file != "" && !validFilenameRegex.MatchString(file) {
		log.Debug().Str("file", file).Msg("file name does not match regex")
		err := fmt.Errorf("invalid file name: %s", file)
		return "", webdav.NewHTTPError(400, err)
	}

	return filepath.Join(localPath, dir, file), nil
}

func etagForFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha1.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	csum := h.Sum(nil)

	return base64.StdEncoding.EncodeToString(csum[:]), nil
}
