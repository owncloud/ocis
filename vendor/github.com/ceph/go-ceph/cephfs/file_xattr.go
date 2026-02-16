package cephfs

/*
#cgo LDFLAGS: -lcephfs
#cgo CPPFLAGS: -D_FILE_OFFSET_BITS=64
#define _GNU_SOURCE
#include <stdlib.h>
#include <sys/types.h>
#include <sys/xattr.h>
#include <cephfs/libcephfs.h>
*/
import "C"

import (
	"unsafe"

	"github.com/ceph/go-ceph/internal/cutil"
	"github.com/ceph/go-ceph/internal/retry"
)

// XattrFlags are used to control the behavior of set-xattr calls.
type XattrFlags int

const (
	// XattrDefault specifies that set-xattr calls use the default behavior of
	// creating or updating an xattr.
	XattrDefault = XattrFlags(0)
	// XattrCreate specifies that set-xattr calls only set new xattrs.
	XattrCreate = XattrFlags(C.XATTR_CREATE)
	// XattrReplace specifies that set-xattr calls only replace existing xattr
	// values.
	XattrReplace = XattrFlags(C.XATTR_REPLACE)
)

// SetXattr sets an extended attribute on the open file.
//
// NOTE: Attempting to set an xattr value with an empty value may cause the
// xattr to be unset on some older versions of ceph.
// Please refer to https://tracker.ceph.com/issues/46084
//
// Implements:
//
//	int ceph_fsetxattr(struct ceph_mount_info *cmount, int fd, const char *name,
//	                   const void *value, size_t size, int flags);
func (f *File) SetXattr(name string, value []byte, flags XattrFlags) error {
	if err := f.validate(); err != nil {
		return err
	}
	if name == "" {
		return errInvalid
	}
	var vptr unsafe.Pointer
	if len(value) > 0 {
		vptr = unsafe.Pointer(&value[0])
	}
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	ret := C.ceph_fsetxattr(
		f.mount.mount,
		f.fd,
		cName,
		vptr,
		C.size_t(len(value)),
		C.int(flags))
	return getError(ret)
}

// GetXattr gets an extended attribute from the open file.
//
// Implements:
//
//	int ceph_fgetxattr(struct ceph_mount_info *cmount, int fd, const char *name,
//	                   void *value, size_t size);
func (f *File) GetXattr(name string) ([]byte, error) {
	if err := f.validate(); err != nil {
		return nil, err
	}
	if name == "" {
		return nil, errInvalid
	}
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	var (
		ret C.int
		err error
		buf []byte
	)
	// range from 1k to 64KiB
	retry.WithSizes(1024, 1<<16, func(size int) retry.Hint {
		buf = make([]byte, size)
		ret = C.ceph_fgetxattr(
			f.mount.mount,
			f.fd,
			cName,
			unsafe.Pointer(&buf[0]),
			C.size_t(size))
		err = getErrorIfNegative(ret)
		return retry.DoubleSize.If(err == errRange)
	})
	if err != nil {
		return nil, err
	}
	return buf[:ret], nil
}

// ListXattr returns a slice containing strings for the name of each xattr set
// on the file.
//
// Implements:
//
//	int ceph_flistxattr(struct ceph_mount_info *cmount, int fd, char *list, size_t size);
func (f *File) ListXattr() ([]string, error) {
	if err := f.validate(); err != nil {
		return nil, err
	}

	var (
		ret C.int
		err error
		buf []byte
	)
	// range from 1k to 64KiB
	retry.WithSizes(1024, 1<<16, func(size int) retry.Hint {
		buf = make([]byte, size)
		ret = C.ceph_flistxattr(
			f.mount.mount,
			f.fd,
			(*C.char)(unsafe.Pointer(&buf[0])),
			C.size_t(size))
		err = getErrorIfNegative(ret)
		return retry.DoubleSize.If(err == errRange)
	})
	if err != nil {
		return nil, err
	}

	names := cutil.SplitSparseBuffer(buf[:ret])
	return names, nil
}

// RemoveXattr removes the named xattr from the open file.
//
// Implements:
//
//	int ceph_fremovexattr(struct ceph_mount_info *cmount, int fd, const char *name);
func (f *File) RemoveXattr(name string) error {
	if err := f.validate(); err != nil {
		return err
	}
	if name == "" {
		return errInvalid
	}
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	ret := C.ceph_fremovexattr(
		f.mount.mount,
		f.fd,
		cName)
	return getError(ret)
}
