package rados

// #cgo LDFLAGS: -lrados
// #include <rados/librados.h>
// #include <stdlib.h>
//
import "C"

// AssertVersion ensures that the object exists and that its internal version
// number is equal to "ver" before writing. "ver" should be a version number
// previously obtained with IOContext.GetLastVersion().
//
// Implements:
//
//	void rados_read_op_assert_version(rados_read_op_t read_op,
//	                                  uint64_t ver)
func (w *WriteOp) AssertVersion(ver uint64) {
	C.rados_write_op_assert_version(w.op, C.uint64_t(ver))
}
