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
	layer *BaseFS
	base  *BaseFS
}

// Primary returns the primary filesystem.
func (d *FallbackFS) Primary() *BaseFS {
	return d.base
}

// Secondary returns the secondary filesystem.
func (d *FallbackFS) Secondary() *BaseFS {
	return d.layer
}

// NewFallbackFS returns a new FallbackFS instance.
func NewFallbackFS(base, layer FS) *FallbackFS {
	return &FallbackFS{
		FS:    FromAfero(afero.NewCopyOnWriteFs(layer, base)),
		base:  &BaseFS{Fs: base},
		layer: &BaseFS{Fs: layer},
	}
}
