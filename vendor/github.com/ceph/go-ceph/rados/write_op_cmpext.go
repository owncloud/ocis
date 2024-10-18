package rados

// #cgo LDFLAGS: -lrados
// #include <rados/librados.h>
// #include <stdlib.h>
//
import "C"

import (
	"unsafe"
)

// WriteOpCmpExtStep holds result of the CmpExt write operation.
// Result is valid only after Operate() was called.
type WriteOpCmpExtStep struct {
	// C returned data:
	prval *C.int

	// Result of the CmpExt write operation.
	Result int
}

func (s *WriteOpCmpExtStep) update() error {
	s.Result = int(*s.prval)
	return nil
}

func (s *WriteOpCmpExtStep) free() {
	C.free(unsafe.Pointer(s.prval))
	s.prval = nil
}

func newWriteOpCmpExtStep() *WriteOpCmpExtStep {
	return &WriteOpCmpExtStep{
		prval: (*C.int)(C.malloc(C.sizeof_int)),
	}
}

// CmpExt ensures that given object range (extent) satisfies comparison.
//
// Implements:
//
//	void rados_write_op_cmpext(rados_write_op_t write_op,
//	                           const char * cmp_buf,
//	                           size_t cmp_len,
//	                           uint64_t off,
//	                           int * prval);
func (w *WriteOp) CmpExt(b []byte, offset uint64) *WriteOpCmpExtStep {
	oe := newWriteStep(b, 0, offset)
	cmpExtStep := newWriteOpCmpExtStep()
	w.steps = append(w.steps, oe, cmpExtStep)
	C.rados_write_op_cmpext(
		w.op,
		oe.cBuffer,
		oe.cDataLen,
		oe.cOffset,
		cmpExtStep.prval)

	return cmpExtStep
}
