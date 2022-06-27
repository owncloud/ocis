package rados

/*
#cgo LDFLAGS: -lrados
#include <stdlib.h>
#include <rados/librados.h>
*/
import "C"

import (
	"runtime"
	"unsafe"
)

// OmapKeyValue items are returned by the GetOmapStep's Next call.
type OmapKeyValue struct {
	Key   string
	Value []byte
}

// GetOmapStep values are used to get the results of an GetOmapValues call
// on a WriteOp. Until the Operate method of the WriteOp is called the Next
// call will return an error. After Operate is called, the Next call will
// return valid results.
//
// The life cycle of the GetOmapStep is bound to the ReadOp, if the ReadOp
// Release method is called the public methods of the step must no longer be
// used and may return errors.
type GetOmapStep struct {
	// C returned data:
	iter C.rados_omap_iter_t
	more *C.uchar
	rval *C.int

	// internal state:

	// canIterate is only set after the operation is performed and is
	// intended to prevent premature fetching of data
	canIterate bool
}

func newGetOmapStep() *GetOmapStep {
	gos := &GetOmapStep{
		more: (*C.uchar)(C.malloc(C.sizeof_uchar)),
		rval: (*C.int)(C.malloc(C.sizeof_int)),
	}
	runtime.SetFinalizer(gos, opStepFinalizer)
	return gos
}

func (gos *GetOmapStep) free() {
	gos.canIterate = false
	if gos.iter != nil {
		C.rados_omap_get_end(gos.iter)
	}
	gos.iter = nil
	C.free(unsafe.Pointer(gos.more))
	gos.more = nil
	C.free(unsafe.Pointer(gos.rval))
	gos.rval = nil
}

func (gos *GetOmapStep) update() error {
	err := getError(*gos.rval)
	gos.canIterate = (err == nil)
	return err
}

// Next returns the next key value pair or nil if iteration is exhausted.
func (gos *GetOmapStep) Next() (*OmapKeyValue, error) {
	if !gos.canIterate {
		return nil, ErrOperationIncomplete
	}
	var (
		cKey    *C.char
		cVal    *C.char
		cKeyLen C.size_t
		cValLen C.size_t
	)
	ret := C.rados_omap_get_next2(gos.iter, &cKey, &cVal, &cKeyLen, &cValLen)
	if ret != 0 {
		return nil, getError(ret)
	}
	if cKey == nil {
		return nil, nil
	}
	return &OmapKeyValue{
		Key:   string(C.GoBytes(unsafe.Pointer(cKey), C.int(cKeyLen))),
		Value: C.GoBytes(unsafe.Pointer(cVal), C.int(cValLen)),
	}, nil
}

// More returns true if there are more matching keys available.
func (gos *GetOmapStep) More() bool {
	// tad bit hacky, but go can't automatically convert from
	// unsigned char to bool
	return *gos.more != 0
}

// SetOmap appends the map `pairs` to the omap `oid`
func (ioctx *IOContext) SetOmap(oid string, pairs map[string][]byte) error {
	op := CreateWriteOp()
	defer op.Release()
	op.SetOmap(pairs)
	return op.operateCompat(ioctx, oid)
}

// OmapListFunc is the type of the function called for each omap key
// visited by ListOmapValues
type OmapListFunc func(key string, value []byte)

// ListOmapValues iterates over the keys and values in an omap by way of
// a callback function.
//
// `startAfter`: iterate only on the keys after this specified one
// `filterPrefix`: iterate only on the keys beginning with this prefix
// `maxReturn`: iterate no more than `maxReturn` key/value pairs
// `listFn`: the function called at each iteration
func (ioctx *IOContext) ListOmapValues(oid string, startAfter string, filterPrefix string, maxReturn int64, listFn OmapListFunc) error {

	op := CreateReadOp()
	defer op.Release()
	gos := op.GetOmapValues(startAfter, filterPrefix, uint64(maxReturn))
	err := op.operateCompat(ioctx, oid)
	if err != nil {
		return err
	}

	for {
		kv, err := gos.Next()
		if err != nil {
			return err
		}
		if kv == nil {
			break
		}
		listFn(kv.Key, kv.Value)
	}
	return nil
}

// GetOmapValues fetches a set of keys and their values from an omap and returns then as a map
// `startAfter`: retrieve only the keys after this specified one
// `filterPrefix`: retrieve only the keys beginning with this prefix
// `maxReturn`: retrieve no more than `maxReturn` key/value pairs
func (ioctx *IOContext) GetOmapValues(oid string, startAfter string, filterPrefix string, maxReturn int64) (map[string][]byte, error) {
	omap := map[string][]byte{}

	err := ioctx.ListOmapValues(
		oid, startAfter, filterPrefix, maxReturn,
		func(key string, value []byte) {
			omap[key] = value
		},
	)

	return omap, err
}

// GetAllOmapValues fetches all the keys and their values from an omap and returns then as a map
// `startAfter`: retrieve only the keys after this specified one
// `filterPrefix`: retrieve only the keys beginning with this prefix
// `iteratorSize`: internal number of keys to fetch during a read operation
func (ioctx *IOContext) GetAllOmapValues(oid string, startAfter string, filterPrefix string, iteratorSize int64) (map[string][]byte, error) {
	omap := map[string][]byte{}
	omapSize := 0

	for {
		err := ioctx.ListOmapValues(
			oid, startAfter, filterPrefix, iteratorSize,
			func(key string, value []byte) {
				omap[key] = value
				startAfter = key
			},
		)

		if err != nil {
			return omap, err
		}

		// End of omap
		if len(omap) == omapSize {
			break
		}

		omapSize = len(omap)
	}

	return omap, nil
}

// RmOmapKeys removes the specified `keys` from the omap `oid`
func (ioctx *IOContext) RmOmapKeys(oid string, keys []string) error {
	op := CreateWriteOp()
	defer op.Release()
	op.RmOmapKeys(keys)
	return op.operateCompat(ioctx, oid)
}

// CleanOmap clears the omap `oid`
func (ioctx *IOContext) CleanOmap(oid string) error {
	op := CreateWriteOp()
	defer op.Release()
	op.CleanOmap()
	return op.operateCompat(ioctx, oid)
}
