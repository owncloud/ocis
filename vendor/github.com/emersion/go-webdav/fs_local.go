package webdav

import (
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/emersion/go-webdav/internal"
)

type LocalFileSystem string

func (fs LocalFileSystem) localPath(name string) (string, error) {
	if (filepath.Separator != '/' && strings.IndexRune(name, filepath.Separator) >= 0) || strings.Contains(name, "\x00") {
		return "", internal.HTTPErrorf(http.StatusBadRequest, "webdav: invalid character in path")
	}
	name = path.Clean(name)
	if !path.IsAbs(name) {
		return "", internal.HTTPErrorf(http.StatusBadRequest, "webdav: expected absolute path, got %q", name)
	}
	return filepath.Join(string(fs), filepath.FromSlash(name)), nil
}

func (fs LocalFileSystem) externalPath(name string) (string, error) {
	rel, err := filepath.Rel(string(fs), name)
	if err != nil {
		return "", err
	}
	return "/" + filepath.ToSlash(rel), nil
}

func (fs LocalFileSystem) Open(name string) (io.ReadCloser, error) {
	p, err := fs.localPath(name)
	if err != nil {
		return nil, err
	}
	return os.Open(p)
}

func fileInfoFromOS(p string, fi os.FileInfo) *FileInfo {
	return &FileInfo{
		Path:    p,
		Size:    fi.Size(),
		ModTime: fi.ModTime(),
		IsDir:   fi.IsDir(),
		// TODO: fallback to http.DetectContentType?
		MIMEType: mime.TypeByExtension(path.Ext(p)),
		// RFC 2616 section 13.3.3 describes strong ETags. Ideally these would
		// be checksums or sequence numbers, however these are expensive to
		// compute. The modification time with nanosecond granularity is good
		// enough, as it's very unlikely for the same file to be modified twice
		// during a single nanosecond.
		ETag: fmt.Sprintf("%x%x", fi.ModTime().UnixNano(), fi.Size()),
	}
}

func (fs LocalFileSystem) Stat(name string) (*FileInfo, error) {
	p, err := fs.localPath(name)
	if err != nil {
		return nil, err
	}
	fi, err := os.Stat(p)
	if err != nil {
		return nil, err
	}
	return fileInfoFromOS(name, fi), nil
}

func (fs LocalFileSystem) Readdir(name string, recursive bool) ([]FileInfo, error) {
	path, err := fs.localPath(name)
	if err != nil {
		return nil, err
	}

	var l []FileInfo
	err = filepath.Walk(path, func(p string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		href, err := fs.externalPath(p)
		if err != nil {
			return err
		}

		l = append(l, *fileInfoFromOS(href, fi))

		if !recursive && fi.IsDir() && path != p {
			return filepath.SkipDir
		}
		return nil
	})
	return l, err
}

func (fs LocalFileSystem) Create(name string) (io.WriteCloser, error) {
	p, err := fs.localPath(name)
	if err != nil {
		return nil, err
	}
	return os.Create(p)
}

func (fs LocalFileSystem) RemoveAll(name string) error {
	p, err := fs.localPath(name)
	if err != nil {
		return err
	}

	// WebDAV semantics are that it should return a "404 Not Found" error in
	// case the resource doesn't exist. We need to Stat before RemoveAll.
	if _, err = os.Stat(p); err != nil {
		return err
	}

	return os.RemoveAll(p)
}

func (fs LocalFileSystem) Mkdir(name string) error {
	p, err := fs.localPath(name)
	if err != nil {
		return err
	}
	return os.Mkdir(p, 0755)
}

func copyRegularFile(src, dst string, perm os.FileMode) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		// TODO: send http.StatusConflict on os.IsNotExist
		return err
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return err
	}

	return dstFile.Close()
}

func (fs LocalFileSystem) Copy(src, dst string, recursive, overwrite bool) (created bool, err error) {
	srcPath, err := fs.localPath(src)
	if err != nil {
		return false, err
	}
	dstPath, err := fs.localPath(dst)
	if err != nil {
		return false, err
	}

	// TODO: "Note that an infinite-depth COPY of /A/ into /A/B/ could lead to
	// infinite recursion if not handled correctly"

	srcInfo, err := os.Stat(srcPath)
	if err != nil {
		return false, err
	}
	srcPerm := srcInfo.Mode() & os.ModePerm

	if _, err := os.Stat(dstPath); err != nil {
		if !os.IsNotExist(err) {
			return false, err
		}
		created = true
	} else {
		if !overwrite {
			return false, os.ErrExist
		}
		if err := os.RemoveAll(dstPath); err != nil {
			return false, err
		}
	}

	err = filepath.Walk(srcPath, func(p string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if fi.IsDir() {
			if err := os.Mkdir(dstPath, srcPerm); err != nil {
				return err
			}
		} else {
			if err := copyRegularFile(srcPath, dstPath, srcPerm); err != nil {
				return err
			}
		}

		if fi.IsDir() && !recursive {
			return filepath.SkipDir
		}
		return nil
	})
	if err != nil {
		return false, err
	}

	return created, nil
}

func (fs LocalFileSystem) MoveAll(src, dst string, overwrite bool) (created bool, err error) {
	srcPath, err := fs.localPath(src)
	if err != nil {
		return false, err
	}
	dstPath, err := fs.localPath(dst)
	if err != nil {
		return false, err
	}

	if _, err := os.Stat(dstPath); err != nil {
		if !os.IsNotExist(err) {
			return false, err
		}
		created = true
	} else {
		if !overwrite {
			return false, os.ErrExist
		}
		if err := os.RemoveAll(dstPath); err != nil {
			return false, err
		}
	}

	if err := os.Rename(srcPath, dstPath); err != nil {
		return false, err
	}

	return created, nil
}

var _ FileSystem = LocalFileSystem("")
