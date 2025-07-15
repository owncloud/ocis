package admin

import (
	"github.com/ceph/go-ceph/common/admin/manager"
)

const mirroring = "mirroring"

// EnableMirroringModule will enable the mirroring module for cephfs.
//
// Similar To:
//
//	ceph mgr module enable mirroring [--force]
func (fsa *FSAdmin) EnableMirroringModule(force bool) error {
	mgradmin := manager.NewFromConn(fsa.conn)
	return mgradmin.EnableModule(mirroring, force)
}

// DisableMirroringModule will disable the mirroring module for cephfs.
//
// Similar To:
//
//	ceph mgr module disable mirroring
func (fsa *FSAdmin) DisableMirroringModule() error {
	mgradmin := manager.NewFromConn(fsa.conn)
	return mgradmin.DisableModule(mirroring)
}
