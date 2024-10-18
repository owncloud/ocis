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
)

// SelectFilesystem selects a file system to be mounted. If the ceph cluster
// supports more than one cephfs this optional function selects which one to
// use. Can only be called prior to calling Mount. The name of the file system
// is not validated by this call - if the supplied file system name is not
// valid then only the subsequent mount call will fail.
//
// Implements:
//
//	int ceph_select_filesystem(struct ceph_mount_info *cmount, const char *fs_name);
func (mount *MountInfo) SelectFilesystem(name string) error {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	ret := C.ceph_select_filesystem(mount.mount, cName)
	return getError(ret)
}
