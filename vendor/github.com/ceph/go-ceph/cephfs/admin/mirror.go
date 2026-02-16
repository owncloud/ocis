package admin

import (
	ccom "github.com/ceph/go-ceph/common/commands"
	"github.com/ceph/go-ceph/internal/commands"
)

// SnapshotMirrorAdmin helps administer the snapshot mirroring features of
// cephfs. Snapshot mirroring is only available in ceph pacific and later.
type SnapshotMirrorAdmin struct {
	conn ccom.MgrCommander
}

// SnapshotMirror returns a new SnapshotMirrorAdmin to be used for the
// administration of snapshot mirroring features.
func (fsa *FSAdmin) SnapshotMirror() *SnapshotMirrorAdmin {
	return &SnapshotMirrorAdmin{conn: fsa.conn}
}

// Enable snapshot mirroring for the given file system.
//
// Similar To:
//
//	ceph fs snapshot mirror enable <fs_name>
func (sma *SnapshotMirrorAdmin) Enable(fsname string) error {
	m := map[string]string{
		"prefix":  "fs snapshot mirror enable",
		"fs_name": fsname,
		"format":  "json",
	}
	return commands.MarshalMgrCommand(sma.conn, m).NoStatus().EmptyBody().End()
}

// Disable snapshot mirroring for the given file system.
//
// Similar To:
//
//	ceph fs snapshot mirror disable <fs_name>
func (sma *SnapshotMirrorAdmin) Disable(fsname string) error {
	m := map[string]string{
		"prefix":  "fs snapshot mirror disable",
		"fs_name": fsname,
		"format":  "json",
	}
	return commands.MarshalMgrCommand(sma.conn, m).NoStatus().EmptyBody().End()
}

// Add a path in the file system to be mirrored.
//
// Similar To:
//
//	ceph fs snapshot mirror add <fs_name> <path>
func (sma *SnapshotMirrorAdmin) Add(fsname, path string) error {
	m := map[string]string{
		"prefix":  "fs snapshot mirror add",
		"fs_name": fsname,
		"path":    path,
		"format":  "json",
	}
	return commands.MarshalMgrCommand(sma.conn, m).NoStatus().EmptyBody().End()
}

// Remove a path in the file system from mirroring.
//
// Similar To:
//
//	ceph fs snapshot mirror remove <fs_name> <path>
func (sma *SnapshotMirrorAdmin) Remove(fsname, path string) error {
	m := map[string]string{
		"prefix":  "fs snapshot mirror remove",
		"fs_name": fsname,
		"path":    path,
		"format":  "json",
	}
	return commands.MarshalMgrCommand(sma.conn, m).NoStatus().EmptyBody().End()
}

type bootstrapTokenResponse struct {
	Token string `json:"token"`
}

// CreatePeerBootstrapToken returns a token that can be used to create
// a peering association between this site an another site.
//
// Similar To:
//
//	ceph fs snapshot mirror peer_bootstrap create <fs_name> <client_entity> <site-name>
func (sma *SnapshotMirrorAdmin) CreatePeerBootstrapToken(
	fsname, client, site string) (string, error) {
	m := map[string]string{
		"prefix":      "fs snapshot mirror peer_bootstrap create",
		"fs_name":     fsname,
		"client_name": client,
		"format":      "json",
	}
	if site != "" {
		m["site_name"] = site
	}
	var bt bootstrapTokenResponse
	err := commands.MarshalMgrCommand(sma.conn, m).NoStatus().Unmarshal(&bt).End()
	return bt.Token, err
}

// ImportPeerBoostrapToken creates an association between another site, one
// that has provided a token, with the current site.
//
// Similar To:
//
//	ceph fs snapshot mirror peer_bootstrap import <fs_name> <token>
func (sma *SnapshotMirrorAdmin) ImportPeerBoostrapToken(fsname, token string) error {
	m := map[string]string{
		"prefix":  "fs snapshot mirror peer_bootstrap import",
		"fs_name": fsname,
		"token":   token,
		"format":  "json",
	}
	return commands.MarshalMgrCommand(sma.conn, m).NoStatus().EmptyBody().End()
}

// DaemonID represents the ID of a cephfs mirroring daemon.
type DaemonID uint

// FileSystemID represents the ID of a cephfs file system.
type FileSystemID uint

// PeerUUID represents the UUID of a cephfs mirroring peer.
type PeerUUID string

// DaemonStatusPeer contains fields detailing a remote peer.
type DaemonStatusPeer struct {
	ClientName  string `json:"client_name"`
	ClusterName string `json:"cluster_name"`
	FSName      string `json:"fs_name"`
}

// DaemonStatusPeerStats contains fields detailing the a remote peer's stats.
type DaemonStatusPeerStats struct {
	FailureCount  uint64 `json:"failure_count"`
	RecoveryCount uint64 `json:"recovery_count"`
}

// DaemonStatusPeerInfo contains fields representing information about a remote peer.
type DaemonStatusPeerInfo struct {
	UUID   PeerUUID              `json:"uuid"`
	Remote DaemonStatusPeer      `json:"remote"`
	Stats  DaemonStatusPeerStats `json:"stats"`
}

// DaemonStatusFileSystemInfo represents information about a mirrored file system.
type DaemonStatusFileSystemInfo struct {
	FileSystemID   FileSystemID           `json:"filesystem_id"`
	Name           string                 `json:"name"`
	DirectoryCount int64                  `json:"directory_count"`
	Peers          []DaemonStatusPeerInfo `json:"peers"`
}

// DaemonStatusInfo maps file system IDs to information about that file system.
type DaemonStatusInfo struct {
	DaemonID    DaemonID                     `json:"daemon_id"`
	FileSystems []DaemonStatusFileSystemInfo `json:"filesystems"`
}

// DaemonStatusResults maps mirroring daemon IDs to information about that
// mirroring daemon.
type DaemonStatusResults []DaemonStatusInfo

func parseDaemonStatus(res response) (DaemonStatusResults, error) {
	var dsr DaemonStatusResults
	if err := res.NoStatus().Unmarshal(&dsr).End(); err != nil {
		return nil, err
	}
	return dsr, nil
}

// DaemonStatus returns information on the status of cephfs mirroring daemons
// associated with the given file system.
//
// Similar To:
//
//	ceph fs snapshot mirror daemon status <fs_name>
func (sma *SnapshotMirrorAdmin) DaemonStatus(fsname string) (
	DaemonStatusResults, error) {
	// ---
	m := map[string]string{
		"prefix":  "fs snapshot mirror daemon status",
		"fs_name": fsname,
		"format":  "json",
	}
	return parseDaemonStatus(commands.MarshalMgrCommand(sma.conn, m))
}

// PeerInfo includes information about a cephfs mirroring peer.
type PeerInfo struct {
	ClientName string `json:"client_name"`
	SiteName   string `json:"site_name"`
	FSName     string `json:"fs_name"`
	MonHost    string `json:"mon_host"`
}

// PeerListResults maps a peer's UUID to information about that peer.
type PeerListResults map[PeerUUID]PeerInfo

func parsePeerList(res response) (PeerListResults, error) {
	var plr PeerListResults
	if err := res.NoStatus().Unmarshal(&plr).End(); err != nil {
		return nil, err
	}
	return plr, nil
}

// PeerList returns information about peers associated with the given file system.
//
// Similar To:
//
//	ceph fs snapshot mirror peer_list <fs_name>
func (sma *SnapshotMirrorAdmin) PeerList(fsname string) (
	PeerListResults, error) {
	// ---
	m := map[string]string{
		"prefix":  "fs snapshot mirror peer_list",
		"fs_name": fsname,
		"format":  "json",
	}
	return parsePeerList(commands.MarshalMgrCommand(sma.conn, m))
}

/*
 DirMap - figure out what last_shuffled is supposed to mean and, if it is a time
          like it seems to be, how best to represent in Go.

 DirMap TODO
ceph fs snapshot mirror dirmap
func (sma *SnapshotMirrorAdmin) DirMap(fsname, path string) error {
	m := map[string]string{
		"prefix":  "fs snapshot mirror dirmap",
		"fs_name": fsname,
		"path":    path,
		"format":  "json",
	}
	return commands.MarshalMgrCommand(sma.conn, m).NoStatus().EmptyBody().End()
}
*/
