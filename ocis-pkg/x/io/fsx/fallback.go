package fsx

import (
	"github.com/spf13/afero"
)

var (
	// assert interfaces implemented
	_ afero.Fs = (*FallbackFS)(nil)
	_ FS       = (*FallbackFS)(nil)
)

// FallbackFS is a filesystem that layers a primary filesystem on top of a secondary filesystem.
type FallbackFS struct {
	FS
	primary   *BaseFS
	secondary *BaseFS
}

// Primary returns the primary filesystem.
func (d *FallbackFS) Primary() *BaseFS {
	return d.primary
}

// Secondary returns the secondary filesystem.
func (d *FallbackFS) Secondary() *BaseFS {
	return d.secondary
}

// NewFallbackFS returns a new FallbackFS instance.
func NewFallbackFS(primary, secondary FS) *FallbackFS {
	return &FallbackFS{
		FS:        FromAfero(afero.NewCopyOnWriteFs(secondary, primary)),
		primary:   &BaseFS{Fs: primary},
		secondary: &BaseFS{Fs: secondary},
	}
}
