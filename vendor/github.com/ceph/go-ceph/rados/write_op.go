package rados

// #cgo LDFLAGS: -lrados
// #include <errno.h>
// #include <stdlib.h>
// #include <rados/librados.h>
//
import "C"

import (
	"unsafe"

	"github.com/ceph/go-ceph/internal/cutil"
	ts "github.com/ceph/go-ceph/internal/timespec"
)

// Timespec is a public type for the internal C 'struct timespec'
type Timespec ts.Timespec

// WriteOp manages a set of discrete actions that will be performed together
// atomically.
type WriteOp struct {
	operation
	op C.rados_write_op_t
}

// CreateWriteOp returns a newly constructed write operation.
func CreateWriteOp() *WriteOp {
	return &WriteOp{
		op: C.rados_create_write_op(),
	}
}

// Release the resources associated with this write operation.
func (w *WriteOp) Release() {
	C.rados_release_write_op(w.op)
	w.op = nil
	w.free()
}

func (w WriteOp) operate2(
	ioctx *IOContext, oid string, mtime *Timespec, flags OperationFlags) error {

	if err := ioctx.validate(); err != nil {
		return err
	}

	cOid := C.CString(oid)
	defer C.free(unsafe.Pointer(cOid))
	var cMtime *C.struct_timespec
	if mtime != nil {
		cMtime = &C.struct_timespec{}
		ts.CopyToCStruct(
			ts.Timespec(*mtime),
			ts.CTimespecPtr(cMtime))
	}

	ret := C.rados_write_op_operate2(
		w.op, ioctx.ioctx, cOid, cMtime, C.int(flags))
	return w.update(writeOp, ret)
}

// Operate will perform the operation(s).
func (w *WriteOp) Operate(ioctx *IOContext, oid string, flags OperationFlags) error {
	return w.operate2(ioctx, oid, nil, flags)
}

// OperateWithMtime will perform the operation while setting the modification
// time stamp to the supplied value.
func (w *WriteOp) OperateWithMtime(
	ioctx *IOContext, oid string, mtime Timespec, flags OperationFlags) error {

	return w.operate2(ioctx, oid, &mtime, flags)
}

func (w *WriteOp) operateCompat(ioctx *IOContext, oid string) error {
	switch err := w.Operate(ioctx, oid, OperationNoFlag).(type) {
	case nil:
		return nil
	case OperationError:
		return err.OpError
	default:
		return err
	}
}

// Create a rados object.
func (w *WriteOp) Create(exclusive CreateOption) {
	// category, the 3rd param, is deprecated and has no effect so we do not
	// implement it in go-ceph
	C.rados_write_op_create(w.op, C.int(exclusive), nil)
}

// SetOmap appends the map `pairs` to the omap `oid`.
func (w *WriteOp) SetOmap(pairs map[string][]byte) {
	keys := make([]string, len(pairs))
	values := make([][]byte, len(pairs))
	idx := 0
	for k, v := range pairs {
		keys[idx] = k
		values[idx] = v
		idx++
	}

	cKeys := cutil.NewBufferGroupStrings(keys)
	cValues := cutil.NewBufferGroupBytes(values)
	defer cKeys.Free()
	defer cValues.Free()

	C.rados_write_op_omap_set2(
		w.op,
		(**C.char)(cKeys.BuffersPtr()),
		(**C.char)(cValues.BuffersPtr()),
		(*C.size_t)(cKeys.LengthsPtr()),
		(*C.size_t)(cValues.LengthsPtr()),
		(C.size_t)(len(pairs)))
}

// RmOmapKeys removes the specified `keys` from the omap `oid`.
func (w *WriteOp) RmOmapKeys(keys []string) {
	cKeys := cutil.NewBufferGroupStrings(keys)
	defer cKeys.Free()

	C.rados_write_op_omap_rm_keys2(
		w.op,
		(**C.char)(cKeys.BuffersPtr()),
		(*C.size_t)(cKeys.LengthsPtr()),
		(C.size_t)(len(keys)))
}

// CleanOmap clears the omap `oid`.
func (w *WriteOp) CleanOmap() {
	C.rados_write_op_omap_clear(w.op)
}

// AssertExists assures the object targeted by the write op exists.
//
// Implements:
//
//	void rados_write_op_assert_exists(rados_write_op_t write_op);
func (w *WriteOp) AssertExists() {
	C.rados_write_op_assert_exists(w.op)
}

// Write a given byte slice at the supplied offset.
//
// Implements:
//
//	void rados_write_op_write(rados_write_op_t write_op,
//	                                     const char *buffer,
//	                                     size_t len,
//	                                     uint64_t offset);
func (w *WriteOp) Write(b []byte, offset uint64) {
	oe := newWriteStep(b, 0, offset)
	w.steps = append(w.steps, oe)
	C.rados_write_op_write(
		w.op,
		oe.cBuffer,
		oe.cDataLen,
		oe.cOffset)
}

// WriteFull writes a given byte slice as the whole object,
// atomically replacing it.
//
// Implements:
//
//	void rados_write_op_write_full(rados_write_op_t write_op,
//	                               const char *buffer,
//	                               size_t len);
func (w *WriteOp) WriteFull(b []byte) {
	oe := newWriteStep(b, 0, 0)
	w.steps = append(w.steps, oe)
	C.rados_write_op_write_full(
		w.op,
		oe.cBuffer,
		oe.cDataLen)
}

// WriteSame write a given byte slice to the object multiple times, until
// writeLen is satisfied.
//
// Implements:
//
//	void rados_write_op_writesame(rados_write_op_t write_op,
//	                              const char *buffer,
//	                              size_t data_len,
//	                              size_t write_len,
//	                              uint64_t offset);
func (w *WriteOp) WriteSame(b []byte, writeLen, offset uint64) {
	oe := newWriteStep(b, writeLen, offset)
	w.steps = append(w.steps, oe)
	C.rados_write_op_writesame(
		w.op,
		oe.cBuffer,
		oe.cDataLen,
		oe.cWriteLen,
		oe.cOffset)
}
