package gowebdav

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type FileInfo interface {
	os.FileInfo
	StatusCode() int
}

// File is our structure for a given file
type File struct {
	path        string
	name        string
	contentType string
	size        int64
	modified    time.Time
	etag        string
	isdir       bool
	propstat    propstat
	status      int
}

func newFile(path, name string, p *propstat) *File {
	f := &File{
		propstat: *p,
	}
	path = FixSlashes(path)

	f.name = filepath.Base(name)
	f.path = filepath.Clean(filepath.Join(path, f.name))
	f.modified = p.Modified()
	f.etag = p.ETag()
	f.contentType = p.ContentType()

	if p.Type() == "collection" {
		f.path = filepath.Clean(f.path + "/")
		f.size = 0
		f.isdir = true
	} else {
		f.size = p.Size()
		f.isdir = false
	}
	return f
}

// Path returns the full path of a file
func (f File) Path() string {
	return f.path
}

// Name returns the name of a file
func (f File) Name() string {
	return f.name
}

// ContentType returns the content type of a file
func (f File) ContentType() string {
	return f.contentType
}

// Size returns the size of a file
func (f File) Size() int64 {
	return f.size
}

// Mode will return the mode of a given file
func (f File) Mode() os.FileMode {
	// TODO check webdav perms
	if f.isdir {
		return 0775 | os.ModeDir
	}

	return 0664
}

// ModTime returns the modified time of a file
func (f File) ModTime() time.Time {
	return f.modified
}

// ETag returns the ETag of a file
func (f File) ETag() string {
	return f.etag
}

// IsDir let us see if a given file is a directory or not
func (f File) IsDir() bool {
	return f.isdir
}

// Sys ????
func (f File) Sys() interface{} {
	return f.propstat.Props
}

func (f File) StatusCode() int {
	return f.propstat.StatusCode()
}

// String lets us see file information
func (f File) String() string {
	if f.isdir {
		return fmt.Sprintf("Dir : '%s' - '%s'", f.path, f.name)
	}

	return fmt.Sprintf("File: '%s' SIZE: %d MODIFIED: %s ETAG: %s CTYPE: %s", f.path, f.size, f.modified.String(), f.etag, f.contentType)
}
