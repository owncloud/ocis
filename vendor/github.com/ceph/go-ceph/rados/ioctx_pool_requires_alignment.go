package rados

// #cgo LDFLAGS: -lrados
// #include <rados/librados.h>
// #include <stdlib.h>
//
import "C"

// RequiresAlignment returns true if the pool supports/requires alignment or an error if not successful.
// For an EC pool, a buffer size multiple of its stripe size is required to call Append. See
// Alignment to know how to get the stripe size for pools requiring it.
//
// Implements:
//
//	int rados_ioctx_pool_requires_alignment2(rados_ioctx_t io, int *req)
func (ioctx *IOContext) RequiresAlignment() (bool, error) {
	var alignRequired C.int
	ret := C.rados_ioctx_pool_requires_alignment2(
		ioctx.ioctx,
		&alignRequired)
	if ret != 0 {
		return false, getError(ret)
	}
	return (alignRequired != 0), nil
}
