package manager

import (
	ccom "github.com/ceph/go-ceph/common/commands"
)

// MgrAdmin is used to administrate ceph's manager (mgr).
type MgrAdmin struct {
	conn ccom.RadosCommander
}

// NewFromConn creates an new management object from a preexisting
// rados connection. The existing connection can be rados.Conn or any
// type implementing the RadosCommander interface.
func NewFromConn(conn ccom.RadosCommander) *MgrAdmin {
	return &MgrAdmin{conn}
}
