package fsx

import (
	"io/fs"

	"github.com/spf13/afero"
)

var (
	// assert interfaces implemented
	_ afero.Fs = (*BaseFS)(nil)
	_ FS       = (*BaseFS)(nil)
)

// FS is our default interface for filesystems.
type FS interface {
	afero.Fs
	IOFS() fs.FS
}

// BaseFS is our default implementation of the FS interface.
type BaseFS struct {
	afero.Fs
}

// IOFS returns the filesystem as an io/fs.FS.
func (b *BaseFS) IOFS() fs.FS {
	return afero.NewIOFS(b)
}

// FromAfero returns a new BaseFS instance from an afero.Fs.
func FromAfero(fSys afero.Fs) *BaseFS {
	return &BaseFS{Fs: fSys}
}

// FromIOFS returns a new BaseFS instance from an io/fs.FS.
func FromIOFS(fSys fs.FS) *BaseFS {
	return FromAfero(&afero.FromIOFS{FS: fSys})
}

// NewBasePathFs returns a new BaseFS which wraps the given filesystem with a base path.
func NewBasePathFs(fSys FS, basePath string) *BaseFS {
	return FromAfero(afero.NewBasePathFs(fSys, basePath))
}

// NewOsFs returns a new BaseFS which wraps the OS filesystem.
func NewOsFs() *BaseFS {
	return FromAfero(afero.NewOsFs())
}

// NewReadOnlyFs returns a new BaseFS which wraps the given filesystem with a read-only filesystem.
func NewReadOnlyFs(FfSys FS) *BaseFS {
	return FromAfero(afero.NewReadOnlyFs(FfSys))
}

// NewMemMapFs returns a new BaseFS which wraps the memory filesystem.
func NewMemMapFs() *BaseFS {
	return FromAfero(afero.NewMemMapFs())
}
