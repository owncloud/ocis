package cephfs

/*
#cgo LDFLAGS: -lcephfs
#cgo CPPFLAGS: -D_FILE_OFFSET_BITS=64
#include <stdlib.h>
#include <cephfs/libcephfs.h>
*/
import "C"

import (
	"unsafe"
)

// MakeDirs creates multiple directories at once.
//
// Implements:
//
//	int ceph_mkdirs(struct ceph_mount_info *cmount, const char *path, mode_t mode);
func (mount *MountInfo) MakeDirs(path string, mode uint32) error {
	if err := mount.validate(); err != nil {
		return err
	}
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))

	ret := C.ceph_mkdirs(mount.mount, cPath, C.mode_t(mode))
	return getError(ret)
}
