package cephfs

/*
#cgo LDFLAGS: -lcephfs
#cgo CPPFLAGS: -D_FILE_OFFSET_BITS=64
#include <cephfs/libcephfs.h>
*/
import "C"

// Some general connectivity and mounting functions are new in
// Ceph Nautilus.

// GetFsCid returns the cluster ID for a mounted ceph file system.
// If the object does not refer to a mounted file system, an error
// will be returned.
//
// Note:
//
//	Only supported in Ceph Nautilus and newer.
//
// Implements:
//
//	int64_t ceph_get_fs_cid(struct ceph_mount_info *cmount);
func (mount *MountInfo) GetFsCid() (int64, error) {
	ret := C.ceph_get_fs_cid(mount.mount)
	if ret < 0 {
		return 0, getError(C.int(ret))
	}
	return int64(ret), nil
}
