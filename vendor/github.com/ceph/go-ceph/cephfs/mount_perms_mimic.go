//
// ceph_mount_perms_set available in mimic & later

package cephfs

/*
#cgo LDFLAGS: -lcephfs
#cgo CPPFLAGS: -D_FILE_OFFSET_BITS=64
#include <cephfs/libcephfs.h>
*/
import "C"

// SetMountPerms applies the given UserPerm to the mount object, which it will
// then use to define the connection's ownership credentials.
// This function must be called after Init but before Mount.
//
// Implements:
//
//	int ceph_mount_perms_set(struct ceph_mount_info *cmount, UserPerm *perm);
func (mount *MountInfo) SetMountPerms(perm *UserPerm) error {
	return getError(C.ceph_mount_perms_set(mount.mount, perm.userPerm))
}
