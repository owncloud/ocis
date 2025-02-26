package rados

// #cgo LDFLAGS: -lrados
// #include <rados/librados.h>
// #include <stdlib.h>
//
import "C"

// Remove object.
//
// Implements:
//
//	void rados_write_op_remove(rados_write_op_t write_op)
func (w *WriteOp) Remove() {
	C.rados_write_op_remove(w.op)
}
