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

// CephStatVFS instances are returned from the StatFS call. It reports
// file-system wide statistics.
type CephStatVFS struct {
	// Bsize reports the file system's block size.
	Bsize int64
	// Fragment reports the file system's fragment size.
	Frsize int64
	// Blocks reports the number of blocks in the file system.
	Blocks uint64
	// Bfree reports the number of free blocks.
	Bfree uint64
	// Bavail reports the number of free blocks for unprivileged users.
	Bavail uint64
	// Files reports the number of inodes in the file system.
	Files uint64
	// Ffree reports the number of free indoes.
	Ffree uint64
	// Favail reports the number of free indoes for unprivileged users.
	Favail uint64
	// Fsid reports the file system ID number.
	Fsid int64
	// Flag reports the file system mount flags.
	Flag int64
	// Namemax reports the maximum file name length.
	Namemax int64
}

// StatFS returns file system wide statistics.
// NOTE: Many of the statistics fields reported by ceph are not filled in with
// useful values.
//
// Implements:
//
//	int ceph_statfs(struct ceph_mount_info *cmount, const char *path, struct statvfs *stbuf);
func (mount *MountInfo) StatFS(path string) (*CephStatVFS, error) {
	if err := mount.validate(); err != nil {
		return nil, err
	}
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))

	var statvfs C.struct_statvfs
	ret := C.ceph_statfs(mount.mount, cPath, &statvfs)
	if ret != 0 {
		return nil, getError(ret)
	}
	csfs := &CephStatVFS{
		Bsize:   int64(statvfs.f_bsize),
		Frsize:  int64(statvfs.f_frsize),
		Blocks:  uint64(statvfs.f_blocks),
		Bfree:   uint64(statvfs.f_bfree),
		Bavail:  uint64(statvfs.f_bavail),
		Files:   uint64(statvfs.f_files),
		Ffree:   uint64(statvfs.f_ffree),
		Favail:  uint64(statvfs.f_favail),
		Fsid:    int64(statvfs.f_fsid),
		Flag:    int64(statvfs.f_flag),
		Namemax: int64(statvfs.f_namemax),
	}
	return csfs, nil
}
