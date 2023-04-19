package cephfs

/*
#cgo LDFLAGS: -lcephfs
#cgo CPPFLAGS: -D_FILE_OFFSET_BITS=64
#include <cephfs/libcephfs.h>
#ifndef AT_STATX_DONT_SYNC
// for versions earlier than Pacific
#define AT_STATX_DONT_SYNC AT_NO_ATTR_SYNC
#endif
*/
import "C"

import (
	ts "github.com/ceph/go-ceph/internal/timespec"
)

// Timespec is a public type for the internal C 'struct timespec'
type Timespec ts.Timespec

// StatxMask values contain bit-flags indicating what data should be
// populated by a statx-type call.
type StatxMask uint32

const (
	// StatxMode requests the mode value be filled in.
	StatxMode = StatxMask(C.CEPH_STATX_MODE)
	// StatxNlink requests the nlink value be filled in.
	StatxNlink = StatxMask(C.CEPH_STATX_NLINK)
	// StatxUid requests the uid value be filled in.
	StatxUid = StatxMask(C.CEPH_STATX_UID)
	// StatxRdev requests the rdev value be filled in.
	StatxRdev = StatxMask(C.CEPH_STATX_RDEV)
	// StatxAtime requests the access-time value be filled in.
	StatxAtime = StatxMask(C.CEPH_STATX_ATIME)
	// StatxMtime requests the modified-time value be filled in.
	StatxMtime = StatxMask(C.CEPH_STATX_MTIME)
	// StatxIno requests the inode be filled in.
	StatxIno = StatxMask(C.CEPH_STATX_INO)
	// StatxSize requests the size value be filled in.
	StatxSize = StatxMask(C.CEPH_STATX_SIZE)
	// StatxBlocks requests the blocks value be filled in.
	StatxBlocks = StatxMask(C.CEPH_STATX_BLOCKS)
	// StatxBasicStats requests all the fields that are part of a
	// traditional stat call.
	StatxBasicStats = StatxMask(C.CEPH_STATX_BASIC_STATS)
	// StatxBtime requests the birth-time value be filled in.
	StatxBtime = StatxMask(C.CEPH_STATX_BTIME)
	// StatxVersion requests the version value be filled in.
	StatxVersion = StatxMask(C.CEPH_STATX_VERSION)
	// StatxAllStats requests all known stat values be filled in.
	StatxAllStats = StatxMask(C.CEPH_STATX_ALL_STATS)
)

// AtFlags represent flags to be passed to calls that control how files
// are used or referenced. For example, not following symlinks.
type AtFlags uint

const (
	// AtStatxDontSync requests that the stat call only fetch locally-cached
	// values if possible, avoiding round trips to a back-end server.
	AtStatxDontSync = AtFlags(C.AT_STATX_DONT_SYNC)
	// AtNoAttrSync requests that the stat call only fetch locally-cached
	// values if possible, avoiding round trips to a back-end server.
	//
	// Deprecated: replaced by AtStatxDontSync
	AtNoAttrSync = AtStatxDontSync
	// AtSymlinkNofollow indicates the call should not follow symlinks
	// but operate on the symlink itself.
	AtSymlinkNofollow = AtFlags(C.AT_SYMLINK_NOFOLLOW)
)

// NOTE: CephStatx fields are meant to be settable by the callers.
// This is the primary reason we use public fields and not accessors
// for the CephStatx type.

// CephStatx instances are returned by extended stat (statx) calls.
// Note that CephStatx results are similar to but not identical
// to (Linux) system statx results.
type CephStatx struct {
	// Mask is a bitmask indicating what fields have been set.
	Mask StatxMask
	// Blksize represents the file system's block size.
	Blksize uint32
	// Nlink is the number of links for the file.
	Nlink uint32
	// Uid (user id) value for the file.
	Uid uint32
	// Gid (group id) value for the file.
	Gid uint32
	// Mode is the file's type and mode value.
	Mode uint16
	// Inode value for the file.
	Inode Inode
	// Size of the file in bytes.
	Size uint64
	// Blocks indicates the number of blocks allocated to the file.
	Blocks uint64
	// Dev describes the device containing this file system.
	Dev uint64
	// Rdev describes the device of this file, if the file is a device.
	Rdev uint64
	// Atime is the access time of this file.
	Atime Timespec
	// Ctime is the status change time of this file.
	Ctime Timespec
	// Mtime is the modification time of this file.
	Mtime Timespec
	// Btime is the creation (birth) time of this file.
	Btime Timespec
	// Version value for the file.
	Version uint64
}

func cStructToCephStatx(s C.struct_ceph_statx) *CephStatx {
	return &CephStatx{
		Mask:    StatxMask(s.stx_mask),
		Blksize: uint32(s.stx_blksize),
		Nlink:   uint32(s.stx_nlink),
		Uid:     uint32(s.stx_uid),
		Gid:     uint32(s.stx_gid),
		Mode:    uint16(s.stx_mode),
		Inode:   Inode(s.stx_ino),
		Size:    uint64(s.stx_size),
		Blocks:  uint64(s.stx_blocks),
		Dev:     uint64(s.stx_dev),
		Rdev:    uint64(s.stx_rdev),
		Atime:   Timespec(ts.CStructToTimespec(ts.CTimespecPtr(&s.stx_atime))),
		Ctime:   Timespec(ts.CStructToTimespec(ts.CTimespecPtr(&s.stx_ctime))),
		Mtime:   Timespec(ts.CStructToTimespec(ts.CTimespecPtr(&s.stx_mtime))),
		Btime:   Timespec(ts.CStructToTimespec(ts.CTimespecPtr(&s.stx_btime))),
		Version: uint64(s.stx_version),
	}
}

/* TODO:
   - enable later when we can test round -trips
   - add time fields

func (c *CephStatx) toCStruct() C.struct_ceph_statx {
	var s C.struct_ceph_statx
	s.stx_mask = C.uint32_t(c.Mask)
	s.stx_blksize = C.uint32_t(c.Blksize)
	s.stx_nlink = C.uint32_t(c.Nlink)
	s.stx_uid = C.uint32_t(c.Uid)
	s.stx_gid = C.uint32_t(c.Gid)
	s.stx_mode = C.uint16_t(c.Mode)
	s.stx_ino = C.uint64_t(c.Inode)
	s.stx_size = C.uint64_t(c.Size)
	s.stx_blocks = C.uint64_t(c.Blocks)
	s.stx_dev = C.uint64_t(c.Dev)
	s.stx_rdev = C.uint64_t(c.Rdev)
	s.stx_version = C.uint64_t(c.Version)
	return s
}
*/
