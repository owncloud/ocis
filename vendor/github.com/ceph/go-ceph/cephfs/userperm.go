package cephfs

/*
#cgo LDFLAGS: -lcephfs
#cgo CPPFLAGS: -D_FILE_OFFSET_BITS=64
#include <cephfs/libcephfs.h>
*/
import "C"

import (
	"runtime"
	"unsafe"

	"github.com/ceph/go-ceph/internal/log"
)

// UserPerm types may be used to get or change the credentials used by the
// connection or some operations.
type UserPerm struct {
	userPerm *C.UserPerm

	// cache create-time params
	managed bool // if set, the userPerm was created by go-ceph
	uid     C.uid_t
	gid     C.gid_t
	gidList []C.gid_t
}

// NewUserPerm creates a UserPerm pointer and the underlying ceph resources.
//
// Implements:
//
//	UserPerm *ceph_userperm_new(uid_t uid, gid_t gid, int ngids, gid_t *gidlist);
func NewUserPerm(uid, gid int, gidlist []int) *UserPerm {
	// the C code does not copy the content of the gid list so we keep the
	// inputs stashed in the go type. For completeness we stash everything.
	p := &UserPerm{
		managed: true,
		uid:     C.uid_t(uid),
		gid:     C.gid_t(gid),
		gidList: make([]C.gid_t, len(gidlist)),
	}
	var cgids *C.gid_t
	if len(p.gidList) > 0 {
		for i, gid := range gidlist {
			p.gidList[i] = C.gid_t(gid)
		}
		cgids = (*C.gid_t)(unsafe.Pointer(&p.gidList[0]))
	}
	p.userPerm = C.ceph_userperm_new(
		p.uid, p.gid, C.int(len(p.gidList)), cgids)
	// if the go object is unreachable, we would like to free the c-memory
	// since this has no other resources than memory associated with it.
	// This is only valid for UserPerm objects created by new, and thus have
	// the managed var set.
	runtime.SetFinalizer(p, destroyUserPerm)
	return p
}

// Destroy will explicitly free ceph resources associated with the UserPerm.
//
// Implements:
//
//	void ceph_userperm_destroy(UserPerm *perm);
func (p *UserPerm) Destroy() {
	if p.userPerm == nil || !p.managed {
		return
	}
	C.ceph_userperm_destroy(p.userPerm)
	p.userPerm = nil
	p.gidList = nil
}

func destroyUserPerm(p *UserPerm) {
	if p.userPerm != nil && p.managed {
		log.Warnf("unreachable UserPerm object has not been destroyed. Cleaning up.")
	}
	p.Destroy()
}
