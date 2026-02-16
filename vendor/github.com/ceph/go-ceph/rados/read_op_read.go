package rados

// #cgo LDFLAGS: -lrados
// #include <rados/librados.h>
// #include <stdlib.h>
//
import "C"

import (
	"unsafe"
)

// ReadOpReadStep holds the result of the Read read operation.
// Result is valid only after Operate() was called.
type ReadOpReadStep struct {
	// C returned data:
	bytesRead *C.size_t
	prval     *C.int

	BytesRead int64 // Bytes read by this action.
	Result    int   // Result of this action.
}

func (s *ReadOpReadStep) update() error {
	s.BytesRead = (int64)(*s.bytesRead)
	s.Result = (int)(*s.prval)

	return nil
}

func (s *ReadOpReadStep) free() {
	C.free(unsafe.Pointer(s.bytesRead))
	C.free(unsafe.Pointer(s.prval))

	s.bytesRead = nil
	s.prval = nil
}

func newReadOpReadStep() *ReadOpReadStep {
	return &ReadOpReadStep{
		bytesRead: (*C.size_t)(C.malloc(C.sizeof_size_t)),
		prval:     (*C.int)(C.malloc(C.sizeof_int)),
	}
}

// Read bytes from offset into buffer.
// len(buffer) is the maximum number of bytes read from the object.
// buffer[:ReadOpReadStep.BytesRead] then contains object data.
//
// Implements:
//
//	void rados_read_op_read(rados_read_op_t read_op,
//	                        uint64_t offset,
//	                        size_t len,
//	                        char * buffer,
//	                        size_t * bytes_read,
//	                        int * prval)
func (r *ReadOp) Read(offset uint64, buffer []byte) *ReadOpReadStep {
	oe := newReadStep(buffer, offset)
	readStep := newReadOpReadStep()
	r.steps = append(r.steps, oe, readStep)
	C.rados_read_op_read(
		r.op,
		oe.cOffset,
		oe.cReadLen,
		oe.cBuffer,
		readStep.bytesRead,
		readStep.prval,
	)

	return readStep
}
