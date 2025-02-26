package rados

// #cgo LDFLAGS: -lrados
// #include <rados/librados.h>
// #include <stdlib.h>
//
import "C"

// SetAllocationHint sets allocation hint for an object. This is an advisory
// operation, it will always succeed (as if it was submitted with a
// LIBRADOS_OP_FLAG_FAILOK flag set) and is not guaranteed to do anything on
// the backend.
//
// Implements:
//
//	void rados_write_op_set_alloc_hint2(rados_write_op_t write_op,
//	                                    uint64_t expected_object_size,
//	                                    uint64_t expected_write_size,
//	                                    uint32_t flags);
func (w *WriteOp) SetAllocationHint(expectedObjectSize uint64, expectedWriteSize uint64, flags AllocHintFlags) {
	C.rados_write_op_set_alloc_hint2(
		w.op,
		C.uint64_t(expectedObjectSize),
		C.uint64_t(expectedWriteSize),
		C.uint32_t(flags))
}
