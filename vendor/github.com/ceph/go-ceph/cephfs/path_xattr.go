package cephfs

/*
#cgo LDFLAGS: -lcephfs
#cgo CPPFLAGS: -D_FILE_OFFSET_BITS=64
#define _GNU_SOURCE
#include <stdlib.h>
#include <cephfs/libcephfs.h>
*/
import "C"

import (
	"unsafe"

	"github.com/ceph/go-ceph/internal/cutil"
	"github.com/ceph/go-ceph/internal/retry"
)

// SetXattr sets an extended attribute on the file at the supplied path.
//
// NOTE: Attempting to set an xattr value with an empty value may cause
// the xattr to be unset. Please refer to https://tracker.ceph.com/issues/46084
//
// Implements:
//
//	int ceph_setxattr(struct ceph_mount_info *cmount, const char *path, const char *name,
//	                  const void *value, size_t size, int flags);
func (mount *MountInfo) SetXattr(path, name string, value []byte, flags XattrFlags) error {
	if err := mount.validate(); err != nil {
		return err
	}
	if name == "" {
		return errInvalid
	}
	var vptr unsafe.Pointer
	if len(value) > 0 {
		vptr = unsafe.Pointer(&value[0])
	}
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	ret := C.ceph_setxattr(
		mount.mount,
		cPath,
		cName,
		vptr,
		C.size_t(len(value)),
		C.int(flags))
	return getError(ret)
}

// GetXattr gets an extended attribute from the file at the supplied path.
//
// Implements:
//
//	int ceph_getxattr(struct ceph_mount_info *cmount, const char *path, const char *name,
//	                  void *value, size_t size);
func (mount *MountInfo) GetXattr(path, name string) ([]byte, error) {
	if err := mount.validate(); err != nil {
		return nil, err
	}
	if name == "" {
		return nil, errInvalid
	}
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))
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
		ret = C.ceph_getxattr(
			mount.mount,
			cPath,
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
// on the file at the supplied path.
//
// Implements:
//
//	int ceph_listxattr(struct ceph_mount_info *cmount, const char *path, char *list, size_t size);
func (mount *MountInfo) ListXattr(path string) ([]string, error) {
	if err := mount.validate(); err != nil {
		return nil, err
	}
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))

	var (
		ret C.int
		err error
		buf []byte
	)
	// range from 1k to 64KiB
	retry.WithSizes(1024, 1<<16, func(size int) retry.Hint {
		buf = make([]byte, size)
		ret = C.ceph_listxattr(
			mount.mount,
			cPath,
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
//	int ceph_removexattr(struct ceph_mount_info *cmount, const char *path, const char *name);
func (mount *MountInfo) RemoveXattr(path, name string) error {
	if err := mount.validate(); err != nil {
		return err
	}
	if name == "" {
		return errInvalid
	}
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	ret := C.ceph_removexattr(
		mount.mount,
		cPath,
		cName)
	return getError(ret)
}

// LsetXattr sets an extended attribute on the file at the supplied path.
//
// NOTE: Attempting to set an xattr value with an empty value may cause
// the xattr to be unset. Please refer to https://tracker.ceph.com/issues/46084
//
// Implements:
//
//	int ceph_lsetxattr(struct ceph_mount_info *cmount, const char *path, const char *name,
//	                  const void *value, size_t size, int flags);
func (mount *MountInfo) LsetXattr(path, name string, value []byte, flags XattrFlags) error {
	if err := mount.validate(); err != nil {
		return err
	}
	if name == "" {
		return errInvalid
	}
	var vptr unsafe.Pointer
	if len(value) > 0 {
		vptr = unsafe.Pointer(&value[0])
	}
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	ret := C.ceph_lsetxattr(
		mount.mount,
		cPath,
		cName,
		vptr,
		C.size_t(len(value)),
		C.int(flags))
	return getError(ret)
}

// LgetXattr gets an extended attribute from the file at the supplied path.
//
// Implements:
//
//	int ceph_lgetxattr(struct ceph_mount_info *cmount, const char *path, const char *name,
//	                  void *value, size_t size);
func (mount *MountInfo) LgetXattr(path, name string) ([]byte, error) {
	if err := mount.validate(); err != nil {
		return nil, err
	}
	if name == "" {
		return nil, errInvalid
	}
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))
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
		ret = C.ceph_lgetxattr(
			mount.mount,
			cPath,
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

// LlistXattr returns a slice containing strings for the name of each xattr set
// on the file at the supplied path.
//
// Implements:
//
//	int ceph_llistxattr(struct ceph_mount_info *cmount, const char *path, char *list, size_t size);
func (mount *MountInfo) LlistXattr(path string) ([]string, error) {
	if err := mount.validate(); err != nil {
		return nil, err
	}
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))

	var (
		ret C.int
		err error
		buf []byte
	)
	// range from 1k to 64KiB
	retry.WithSizes(1024, 1<<16, func(size int) retry.Hint {
		buf = make([]byte, size)
		ret = C.ceph_llistxattr(
			mount.mount,
			cPath,
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

// LremoveXattr removes the named xattr from the open file.
//
// Implements:
//
//	int ceph_lremovexattr(struct ceph_mount_info *cmount, const char *path, const char *name);
func (mount *MountInfo) LremoveXattr(path, name string) error {
	if err := mount.validate(); err != nil {
		return err
	}
	if name == "" {
		return errInvalid
	}
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	ret := C.ceph_lremovexattr(
		mount.mount,
		cPath,
		cName)
	return getError(ret)
}
