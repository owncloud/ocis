//go:build !nautilus
// +build !nautilus

package cephfs

/*
#cgo LDFLAGS: -lcephfs
#cgo CPPFLAGS: -D_FILE_OFFSET_BITS=64
#include <errno.h>
#include <stdlib.h>
#include <cephfs/libcephfs.h>
*/
import "C"

import (
	ts "github.com/ceph/go-ceph/internal/timespec"
	"unsafe"
)

// Mknod creates a regular, block or character special file.
//
// Implements:
//
//	int ceph_mknod(struct ceph_mount_info *cmount, const char *path, mode_t mode,
//				   dev_t rdev);
func (mount *MountInfo) Mknod(path string, mode uint16, dev uint16) error {
	if err := mount.validate(); err != nil {
		return err
	}

	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))

	ret := C.ceph_mknod(mount.mount, cPath, C.mode_t(mode), C.dev_t(dev))
	return getError(ret)
}

// Utime struct is the equivalent of C.struct_utimbuf
type Utime struct {
	// AcTime  represents the file's access time in seconds since the Unix epoch.
	AcTime int64
	// ModTime represents the file's modification time in seconds since the Unix epoch.
	ModTime int64
}

// Futime changes file/directory last access and modification times.
//
// Implements:
//
//	int ceph_futime(struct ceph_mount_info *cmount, int fd, struct utimbuf *buf);
func (mount *MountInfo) Futime(fd int, times *Utime) error {
	if err := mount.validate(); err != nil {
		return err
	}

	cFd := C.int(fd)
	uTimeBuf := &C.struct_utimbuf{
		actime:  C.time_t(times.AcTime),
		modtime: C.time_t(times.ModTime),
	}

	ret := C.ceph_futime(mount.mount, cFd, uTimeBuf)
	return getError(ret)
}

// Timeval struct is the go equivalent of C.struct_timeval type
type Timeval struct {
	// Sec represents seconds
	Sec int64
	// USec represents microseconds
	USec int64
}

// Futimens changes file/directory last access and modification times, here times param
// is an array of Timespec struct having length 2, where times[0] represents the access time
// and times[1] represents the modification time.
//
// Implements:
//
//	int ceph_futimens(struct ceph_mount_info *cmount, int fd, struct timespec times[2]);
func (mount *MountInfo) Futimens(fd int, times []Timespec) error {
	if err := mount.validate(); err != nil {
		return err
	}

	if len(times) != 2 {
		return getError(-C.EINVAL)
	}

	cFd := C.int(fd)
	cTimes := []C.struct_timespec{}
	for _, val := range times {
		cTs := &C.struct_timespec{}
		ts.CopyToCStruct(
			ts.Timespec(val),
			ts.CTimespecPtr(cTs),
		)
		cTimes = append(cTimes, *cTs)
	}

	ret := C.ceph_futimens(mount.mount, cFd, &cTimes[0])
	return getError(ret)
}

// Futimes changes file/directory last access and modification times, here times param
// is an array of Timeval struct type having length 2, where times[0] represents the access time
// and times[1] represents the modification time.
//
// Implements:
//
//	int ceph_futimes(struct ceph_mount_info *cmount, int fd, struct timeval times[2]);
func (mount *MountInfo) Futimes(fd int, times []Timeval) error {
	if err := mount.validate(); err != nil {
		return err
	}

	if len(times) != 2 {
		return getError(-C.EINVAL)
	}

	cFd := C.int(fd)
	cTimes := []C.struct_timeval{}
	for _, val := range times {
		cTimes = append(cTimes, C.struct_timeval{
			tv_sec:  C.time_t(val.Sec),
			tv_usec: C.suseconds_t(val.USec),
		})
	}

	ret := C.ceph_futimes(mount.mount, cFd, &cTimes[0])
	return getError(ret)
}
