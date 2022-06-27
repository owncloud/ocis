package admin

import (
	"github.com/ceph/go-ceph/common/admin/manager"
)

const mirroring = "mirroring"

// EnableModule will enable the specified manager module.
//
// Deprecated: use the equivalent function in cluster/admin/manager.
//
// Similar To:
//  ceph mgr module enable <module> [--force]
func (fsa *FSAdmin) EnableModule(module string, force bool) error {
	mgradmin := manager.NewFromConn(fsa.conn)
	return mgradmin.EnableModule(module, force)
}

// DisableModule will disable the specified manager module.
//
// Deprecated: use the equivalent function in cluster/admin/manager.
//
// Similar To:
//  ceph mgr module disable <module>
func (fsa *FSAdmin) DisableModule(module string) error {
	mgradmin := manager.NewFromConn(fsa.conn)
	return mgradmin.DisableModule(module)
}

// EnableMirroringModule will enable the mirroring module for cephfs.
//
// Similar To:
//  ceph mgr module enable mirroring [--force]
func (fsa *FSAdmin) EnableMirroringModule(force bool) error {
	mgradmin := manager.NewFromConn(fsa.conn)
	return mgradmin.EnableModule(mirroring, force)
}

// DisableMirroringModule will disable the mirroring module for cephfs.
//
// Similar To:
//  ceph mgr module disable mirroring
func (fsa *FSAdmin) DisableMirroringModule() error {
	mgradmin := manager.NewFromConn(fsa.conn)
	return mgradmin.DisableModule(mirroring)
}
