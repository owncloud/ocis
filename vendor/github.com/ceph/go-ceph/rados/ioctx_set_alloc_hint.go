package rados

// #cgo LDFLAGS: -lrados
// #include <rados/librados.h>
// #include <stdlib.h>
//
import "C"

import (
	"unsafe"
)

// SetAllocationHint sets allocation hint for an object. This is an advisory
// operation, it will always succeed (as if it was submitted with a
// LIBRADOS_OP_FLAG_FAILOK flag set) and is not guaranteed to do anything on
// the backend.
//
// Implements:
//
//	int rados_set_alloc_hint2(rados_ioctx_t io,
//	                          const char *o,
//	                          uint64_t expected_object_size,
//	                          uint64_t expected_write_size,
//	                          uint32_t flags);
func (ioctx *IOContext) SetAllocationHint(oid string, expectedObjectSize uint64, expectedWriteSize uint64, flags AllocHintFlags) error {
	coid := C.CString(oid)
	defer C.free(unsafe.Pointer(coid))

	return getError(C.rados_set_alloc_hint2(
		ioctx.ioctx,
		coid,
		(C.uint64_t)(expectedObjectSize),
		(C.uint64_t)(expectedWriteSize),
		(C.uint32_t)(flags),
	))
}
