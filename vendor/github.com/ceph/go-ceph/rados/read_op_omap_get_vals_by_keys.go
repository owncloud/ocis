package rados

// #cgo LDFLAGS: -lrados
// #include <rados/librados.h>
// #include <stdlib.h>
//
import "C"

import (
	"unsafe"

	"github.com/ceph/go-ceph/internal/cutil"
)

// ReadOpOmapGetValsByKeysStep holds the result of the
// GetOmapValuesByKeys read operation.
// Result is valid only after Operate() was called.
type ReadOpOmapGetValsByKeysStep struct {
	// C arguments

	iter  C.rados_omap_iter_t
	prval *C.int

	// Internal state

	// canIterate is only set after the operation is performed and is
	// intended to prevent premature fetching of data.
	canIterate bool
}

func newReadOpOmapGetValsByKeysStep() *ReadOpOmapGetValsByKeysStep {
	s := &ReadOpOmapGetValsByKeysStep{
		prval: (*C.int)(C.malloc(C.sizeof_int)),
	}

	return s
}

func (s *ReadOpOmapGetValsByKeysStep) free() {
	s.canIterate = false
	C.rados_omap_get_end(s.iter)

	C.free(unsafe.Pointer(s.prval))
	s.prval = nil
}

func (s *ReadOpOmapGetValsByKeysStep) update() error {
	err := getError(*s.prval)
	s.canIterate = (err == nil)

	return err
}

// Next gets the next omap key/value pair referenced by
// ReadOpOmapGetValsByKeysStep's internal iterator.
// If there are no more elements to retrieve, (nil, nil) is returned.
// May be called only after Operate() finished.
func (s *ReadOpOmapGetValsByKeysStep) Next() (*OmapKeyValue, error) {
	if !s.canIterate {
		return nil, ErrOperationIncomplete
	}

	var (
		cKey    *C.char
		cVal    *C.char
		cKeyLen C.size_t
		cValLen C.size_t
	)

	ret := C.rados_omap_get_next2(s.iter, &cKey, &cVal, &cKeyLen, &cValLen)
	if ret != 0 {
		return nil, getError(ret)
	}

	if cKey == nil {
		// Iterator has reached the end of the list.
		return nil, nil
	}

	return &OmapKeyValue{
		Key:   string(C.GoBytes(unsafe.Pointer(cKey), C.int(cKeyLen))),
		Value: C.GoBytes(unsafe.Pointer(cVal), C.int(cValLen)),
	}, nil
}

// GetOmapValuesByKeys starts iterating over specific key/value pairs.
//
// Implements:
//
//	void rados_read_op_omap_get_vals_by_keys2(rados_read_op_t read_op,
//	                                          char const * const * keys,
//	                                          size_t num_keys,
//	                                          const size_t * key_lens,
//	                                          rados_omap_iter_t * iter,
//	                                          int * prval)
func (r *ReadOp) GetOmapValuesByKeys(keys []string) *ReadOpOmapGetValsByKeysStep {
	s := newReadOpOmapGetValsByKeysStep()
	r.steps = append(r.steps, s)

	cKeys := cutil.NewBufferGroupStrings(keys)
	defer cKeys.Free()

	C.rados_read_op_omap_get_vals_by_keys2(
		r.op,
		(**C.char)(cKeys.BuffersPtr()),
		C.size_t(len(keys)),
		(*C.size_t)(cKeys.LengthsPtr()),
		&s.iter,
		s.prval,
	)

	return s
}
