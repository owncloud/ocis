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

// CurrentDir gets the current working directory.
func (mount *MountInfo) CurrentDir() string {
	if err := mount.validate(); err != nil {
		return ""
	}
	cDir := C.ceph_getcwd(mount.mount)
	return C.GoString(cDir)
}

// ChangeDir changes the current working directory.
func (mount *MountInfo) ChangeDir(path string) error {
	if err := mount.validate(); err != nil {
		return err
	}
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))

	ret := C.ceph_chdir(mount.mount, cPath)
	return getError(ret)
}

// MakeDir creates a directory.
func (mount *MountInfo) MakeDir(path string, mode uint32) error {
	if err := mount.validate(); err != nil {
		return err
	}
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))

	ret := C.ceph_mkdir(mount.mount, cPath, C.mode_t(mode))
	return getError(ret)
}

// RemoveDir removes a directory.
func (mount *MountInfo) RemoveDir(path string) error {
	if err := mount.validate(); err != nil {
		return err
	}
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))

	ret := C.ceph_rmdir(mount.mount, cPath)
	return getError(ret)
}

// Unlink removes a file.
//
// Implements:
//
//	int ceph_unlink(struct ceph_mount_info *cmount, const char *path);
func (mount *MountInfo) Unlink(path string) error {
	if err := mount.validate(); err != nil {
		return err
	}
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))

	ret := C.ceph_unlink(mount.mount, cPath)
	return getError(ret)
}

// Link creates a new link to an existing file.
//
// Implements:
//
//	int ceph_link (struct ceph_mount_info *cmount, const char *existing, const char *newname);
func (mount *MountInfo) Link(oldname, newname string) error {
	if err := mount.validate(); err != nil {
		return err
	}
	cOldname := C.CString(oldname)
	defer C.free(unsafe.Pointer(cOldname))

	cNewname := C.CString(newname)
	defer C.free(unsafe.Pointer(cNewname))

	ret := C.ceph_link(mount.mount, cOldname, cNewname)
	return getError(ret)
}

// Symlink creates a symbolic link to an existing path.
//
// Implements:
//
//	int ceph_symlink(struct ceph_mount_info *cmount, const char *existing, const char *newname);
func (mount *MountInfo) Symlink(existing, newname string) error {
	if err := mount.validate(); err != nil {
		return err
	}
	cExisting := C.CString(existing)
	defer C.free(unsafe.Pointer(cExisting))

	cNewname := C.CString(newname)
	defer C.free(unsafe.Pointer(cNewname))

	ret := C.ceph_symlink(mount.mount, cExisting, cNewname)
	return getError(ret)
}

// Readlink returns the value of a symbolic link.
//
// Implements:
//
//	int ceph_readlink(struct ceph_mount_info *cmount, const char *path, char *buf, int64_t size);
func (mount *MountInfo) Readlink(path string) (string, error) {
	if err := mount.validate(); err != nil {
		return "", err
	}
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))

	buf := make([]byte, 4096)
	ret := C.ceph_readlink(mount.mount,
		cPath,
		(*C.char)(unsafe.Pointer(&buf[0])),
		C.int64_t(len(buf)))
	if ret < 0 {
		return "", getError(ret)
	}

	return string(buf[:ret]), nil
}

// Statx returns information about a file/directory.
//
// Implements:
//
//	int ceph_statx(struct ceph_mount_info *cmount, const char *path, struct ceph_statx *stx,
//	               unsigned int want, unsigned int flags);
func (mount *MountInfo) Statx(path string, want StatxMask, flags AtFlags) (*CephStatx, error) {
	if err := mount.validate(); err != nil {
		return nil, err
	}
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))

	var stx C.struct_ceph_statx
	ret := C.ceph_statx(
		mount.mount,
		cPath,
		&stx,
		C.uint(want),
		C.uint(flags),
	)
	if err := getError(ret); err != nil {
		return nil, err
	}
	return cStructToCephStatx(stx), nil
}

// Rename a file or directory.
//
// Implements:
//
//	int ceph_rename(struct ceph_mount_info *cmount, const char *from, const char *to);
func (mount *MountInfo) Rename(from, to string) error {
	if err := mount.validate(); err != nil {
		return err
	}
	cFrom := C.CString(from)
	defer C.free(unsafe.Pointer(cFrom))
	cTo := C.CString(to)
	defer C.free(unsafe.Pointer(cTo))

	ret := C.ceph_rename(mount.mount, cFrom, cTo)
	return getError(ret)
}

// Truncate sets the size of the specified file.
//
// Implements:
//
//	int ceph_truncate(struct ceph_mount_info *cmount, const char *path, int64_t size);
func (mount *MountInfo) Truncate(path string, size int64) error {
	if err := mount.validate(); err != nil {
		return err
	}
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))

	ret := C.ceph_truncate(
		mount.mount,
		cPath,
		C.int64_t(size),
	)
	return getError(ret)
}
