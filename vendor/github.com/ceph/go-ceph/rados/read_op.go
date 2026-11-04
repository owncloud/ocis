package rados

// #cgo LDFLAGS: -lrados
// #include <errno.h>
// #include <stdlib.h>
// #include <rados/librados.h>
//
import "C"

import (
	"fmt"
	"strings"
	"unsafe"
)

// ReadOp manages a set of discrete object read actions that will be performed
// together atomically.
type ReadOp struct {
	operation
	op C.rados_read_op_t
}

// CreateReadOp returns a newly constructed read operation.
func CreateReadOp() *ReadOp {
	return &ReadOp{
		op: C.rados_create_read_op(),
	}
}

// Release the resources associated with this read operation.
func (r *ReadOp) Release() {
	C.rados_release_read_op(r.op)
	r.op = nil
	r.free()
}

// Operate will perform the operation(s).
func (r *ReadOp) Operate(ioctx *IOContext, oid string, flags OperationFlags) error {
	if err := ioctx.validate(); err != nil {
		return err
	}

	cOid := C.CString(oid)
	defer C.free(unsafe.Pointer(cOid))

	ret := C.rados_read_op_operate(r.op, ioctx.ioctx, cOid, C.int(flags))
	return r.update(readOp, ret)
}

func (r *ReadOp) operateCompat(ioctx *IOContext, oid string) error {
	switch err := r.Operate(ioctx, oid, OperationNoFlag).(type) {
	case nil:
		return nil
	case OperationError:
		// If the op failed, return the bare OpError for backwards
		// compatibility (e.g. ErrNotFound). Otherwise a step failed, so
		// return the whole OperationError so it is detectable via errors.Is.
		if err.OpError != nil {
			return err.OpError
		}
		return err
	default:
		return err
	}
}

// AssertExists assures the object targeted by the read op exists.
//
// Implements:
//
//	void rados_read_op_assert_exists(rados_read_op_t read_op);
func (r *ReadOp) AssertExists() {
	C.rados_read_op_assert_exists(r.op)
}

// GetOmapValues is used to iterate over a set, or sub-set, of omap keys
// as part of a read operation. An GetOmapStep is returned from this
// function. The GetOmapStep may be used to iterate over the key-value
// pairs after the Operate call has been performed.
//
// Note that startAfter and filterPrefix are passed to librados as
// NUL-terminated C strings and therefore must not contain a NUL byte. If
// either argument contains a NUL byte, the resulting Operate call will
// return an error wrapping ErrNulInString instead of silently truncating
// the argument.
func (r *ReadOp) GetOmapValues(startAfter, filterPrefix string, maxReturn uint64) *GetOmapStep {
	gos := newGetOmapStep()
	r.steps = append(r.steps, gos)

	if strings.IndexByte(startAfter, 0) >= 0 {
		gos.err = fmt.Errorf("startAfter: %w", ErrNulInString)
		return gos
	}
	if strings.IndexByte(filterPrefix, 0) >= 0 {
		gos.err = fmt.Errorf("filterPrefix: %w", ErrNulInString)
		return gos
	}

	cStartAfter := C.CString(startAfter)
	cFilterPrefix := C.CString(filterPrefix)
	defer C.free(unsafe.Pointer(cStartAfter))
	defer C.free(unsafe.Pointer(cFilterPrefix))

	C.rados_read_op_omap_get_vals2(
		r.op,
		cStartAfter,
		cFilterPrefix,
		C.uint64_t(maxReturn),
		&gos.iter,
		gos.more,
		gos.rval,
	)
	return gos
}
