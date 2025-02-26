package commands

// MgrCommander in an interface for the API needed to execute JSON formatted
// commands on the ceph mgr.
type MgrCommander interface {
	MgrCommand(buf [][]byte) ([]byte, string, error)
}

// MonCommander is an interface for the API needed to execute JSON formatted
// commands on the ceph mon(s).
type MonCommander interface {
	MonCommand(buf []byte) ([]byte, string, error)
}

// RadosCommander provides an interface for APIs needed to execute JSON
// formatted commands on the Ceph cluster.
type RadosCommander interface {
	MgrCommander
	MonCommander
}
