package rados

// #cgo LDFLAGS: -lrados
// #include <stdlib.h>
// #include <rados/librados.h>
import "C"

import (
	"time"
	"unsafe"

	"github.com/ceph/go-ceph/internal/retry"
)

// CreateSnap creates a pool-wide snapshot.
//
// Implements:
// int rados_ioctx_snap_create(rados_ioctx_t io, const char *snapname)
func (ioctx *IOContext) CreateSnap(snapName string) error {
	if err := ioctx.validate(); err != nil {
		return err
	}

	cSnapName := C.CString(snapName)
	defer C.free(unsafe.Pointer(cSnapName))

	ret := C.rados_ioctx_snap_create(ioctx.ioctx, cSnapName)
	return getError(ret)
}

// RemoveSnap deletes the pool snapshot.
//
// Implements:
//  int rados_ioctx_snap_remove(rados_ioctx_t io, const char *snapname)
func (ioctx *IOContext) RemoveSnap(snapName string) error {
	if err := ioctx.validate(); err != nil {
		return err
	}

	cSnapName := C.CString(snapName)
	defer C.free(unsafe.Pointer(cSnapName))

	ret := C.rados_ioctx_snap_remove(ioctx.ioctx, cSnapName)
	return getError(ret)
}

// SnapID represents the ID of a rados snapshot.
type SnapID C.rados_snap_t

// LookupSnap returns the ID of a pool snapshot.
//
// Implements:
//  int rados_ioctx_snap_lookup(rados_ioctx_t io, const char *name, rados_snap_t *id)
func (ioctx *IOContext) LookupSnap(snapName string) (SnapID, error) {
	var snapID SnapID

	if err := ioctx.validate(); err != nil {
		return snapID, err
	}

	cSnapName := C.CString(snapName)
	defer C.free(unsafe.Pointer(cSnapName))

	ret := C.rados_ioctx_snap_lookup(
		ioctx.ioctx,
		cSnapName,
		(*C.rados_snap_t)(&snapID))
	return snapID, getError(ret)
}

// GetSnapName returns the name of a pool snapshot with the given snapshot ID.
//
// Implements:
//  int rados_ioctx_snap_get_name(rados_ioctx_t io, rados_snap_t id, char *name, int maxlen)
func (ioctx *IOContext) GetSnapName(snapID SnapID) (string, error) {
	if err := ioctx.validate(); err != nil {
		return "", err
	}

	var (
		buf []byte
		err error
	)
	// range from 1k to 64KiB
	retry.WithSizes(1024, 1<<16, func(len int) retry.Hint {
		cLen := C.int(len)
		buf = make([]byte, cLen)
		ret := C.rados_ioctx_snap_get_name(
			ioctx.ioctx,
			(C.rados_snap_t)(snapID),
			(*C.char)(unsafe.Pointer(&buf[0])),
			cLen)
		err = getError(ret)
		return retry.Size(int(cLen)).If(err == errRange)
	})

	if err != nil {
		return "", err
	}
	return C.GoString((*C.char)(unsafe.Pointer(&buf[0]))), nil
}

// GetSnapStamp returns the time of the pool snapshot creation.
//
// Implements:
//  int rados_ioctx_snap_get_stamp(rados_ioctx_t io, rados_snap_t id, time_t *t)
func (ioctx *IOContext) GetSnapStamp(snapID SnapID) (time.Time, error) {
	var cTime C.time_t

	if err := ioctx.validate(); err != nil {
		return time.Unix(int64(cTime), 0), err
	}

	ret := C.rados_ioctx_snap_get_stamp(
		ioctx.ioctx,
		(C.rados_snap_t)(snapID),
		&cTime)
	return time.Unix(int64(cTime), 0), getError(ret)
}

// ListSnaps returns a slice containing the SnapIDs of existing pool snapshots.
//
// Implements:
//  int rados_ioctx_snap_list(rados_ioctx_t io, rados_snap_t *snaps, int maxlen)
func (ioctx *IOContext) ListSnaps() ([]SnapID, error) {
	if err := ioctx.validate(); err != nil {
		return nil, err
	}

	var (
		snapList []SnapID
		cLen     C.int
		err      error
		ret      C.int
	)
	retry.WithSizes(100, 1000, func(maxlen int) retry.Hint {
		cLen = C.int(maxlen)
		snapList = make([]SnapID, cLen)
		ret = C.rados_ioctx_snap_list(
			ioctx.ioctx,
			(*C.rados_snap_t)(unsafe.Pointer(&snapList[0])),
			cLen)
		err = getErrorIfNegative(ret)
		return retry.Size(int(cLen)).If(err == errRange)
	})

	if err != nil {
		return nil, err
	}
	return snapList[:ret], nil
}

// RollbackSnap rollbacks the object with key oID to the pool snapshot.
// The contents of the object will be the same as when the snapshot was taken.
//
// Implements:
//  int rados_ioctx_snap_rollback(rados_ioctx_t io, const char *oid, const char *snapname);
func (ioctx *IOContext) RollbackSnap(oid, snapName string) error {
	if err := ioctx.validate(); err != nil {
		return err
	}

	coid := C.CString(oid)
	defer C.free(unsafe.Pointer(coid))
	cSnapName := C.CString(snapName)
	defer C.free(unsafe.Pointer(cSnapName))

	ret := C.rados_ioctx_snap_rollback(ioctx.ioctx, coid, cSnapName)
	return getError(ret)
}

// SnapHead is the representation of LIBRADOS_SNAP_HEAD from librados.
// SnapHead can be used to reset the IOContext to stop reading from a snapshot.
const SnapHead = SnapID(C.LIBRADOS_SNAP_HEAD)

// SetReadSnap sets the snapshot from which reads are performed.
// Subsequent reads will return data as it was at the time of that snapshot.
// Pass SnapHead for no snapshot (i.e. normal operation).
//
// Implements:
//  void rados_ioctx_snap_set_read(rados_ioctx_t io, rados_snap_t snap);
func (ioctx *IOContext) SetReadSnap(snapID SnapID) error {
	if err := ioctx.validate(); err != nil {
		return err
	}

	C.rados_ioctx_snap_set_read(ioctx.ioctx, (C.rados_snap_t)(snapID))
	return nil
}
