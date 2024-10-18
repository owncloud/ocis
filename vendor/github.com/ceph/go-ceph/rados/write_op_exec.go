//go:build ceph_preview
// +build ceph_preview

package rados

// #cgo LDFLAGS: -lrados
// #include <stdlib.h>
// #include <rados/librados.h>
//
import "C"

import (
	"unsafe"
)

// writeOpExecStep - exec step in write operation.
type writeOpExecStep struct {
	withoutFree

	inBuffPtr *C.char
	inBuffLen C.size_t
	prval     C.int
}

// newWriteOpExecStep - init new *writeOpExecStep.
func newWriteOpExecStep(in []byte) *writeOpExecStep {
	es := &writeOpExecStep{
		prval: 0,
	}
	if len(in) > 0 {
		es.inBuffPtr = (*C.char)(unsafe.Pointer(&in[0]))
		es.inBuffLen = C.size_t(len(in))
	}

	return es
}

// update - update state operation.
func (es *writeOpExecStep) update() error {
	return getError(es.prval)
}

// Exec executes an OSD class method on an object.
// See rados_exec() in the RADOS C api documentation for a general description.
//
// Implements:
//
//	void rados_write_op_exec(rados_write_op_t write_op,
//	                         const char *cls,
//	                         const char *method,
//	                         const char *in_buf,
//	                         size_t in_len,
//	                         int *prval)
func (w *WriteOp) Exec(clsName, method string, in []byte) {
	cClsName := C.CString(clsName)
	defer C.free(unsafe.Pointer(cClsName))

	cMethod := C.CString(method)
	defer C.free(unsafe.Pointer(cMethod))

	es := newWriteOpExecStep(in)
	w.steps = append(w.steps, es)
	C.rados_write_op_exec(
		w.op,
		cClsName,
		cMethod,
		es.inBuffPtr,
		es.inBuffLen,
		&es.prval,
	)
}
