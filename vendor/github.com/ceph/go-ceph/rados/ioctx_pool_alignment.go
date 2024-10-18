package rados

// #cgo LDFLAGS: -lrados
// #include <rados/librados.h>
// #include <stdlib.h>
//
import "C"

// Alignment returns the required stripe size in bytes for pools supporting/requiring it, or an error if unsuccessful.
// For an EC pool, a buffer size multiple of its stripe size is required to call Append. To know if the pool requires
// alignment or not, use RequiresAlignment.
//
// Implements:
//
//	int rados_ioctx_pool_required_alignment2(rados_ioctx_t io, uint64_t *alignment)
func (ioctx *IOContext) Alignment() (uint64, error) {
	var alignSizeBytes C.uint64_t
	ret := C.rados_ioctx_pool_required_alignment2(
		ioctx.ioctx,
		&alignSizeBytes)
	if ret != 0 {
		return 0, getError(ret)
	}
	return uint64(alignSizeBytes), nil
}
