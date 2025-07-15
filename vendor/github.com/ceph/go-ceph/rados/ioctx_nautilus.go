//go:build nautilus
// +build nautilus

package rados

// #cgo LDFLAGS: -lrados
// #include <rados/librados.h>
//
import "C"

// SetPoolFullTry makes sure to send requests to the cluster despite
// the cluster or pool being marked full; ops will either succeed(e.g., delete)
// or return EDQUOT or ENOSPC.
//
// Implements:
//
//	void rados_set_osdmap_full_try(rados_ioctx_t io);
func (ioctx *IOContext) SetPoolFullTry() error {
	if err := ioctx.validate(); err != nil {
		return err
	}
	C.rados_set_osdmap_full_try(ioctx.ioctx)
	return nil
}

// UnsetPoolFullTry unsets the flag set by SetPoolFullTry()
//
// Implements:
//
//	void rados_unset_osdmap_full_try(rados_ioctx_t io);
func (ioctx *IOContext) UnsetPoolFullTry() error {
	if err := ioctx.validate(); err != nil {
		return err
	}
	C.rados_unset_osdmap_full_try(ioctx.ioctx)
	return nil
}
