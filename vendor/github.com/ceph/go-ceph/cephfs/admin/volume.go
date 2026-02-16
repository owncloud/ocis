package admin

import (
	"bytes"
	"encoding/json"
)

var (
	listVolumesCmd = []byte(`{"prefix":"fs volume ls"}`)
	dumpVolumesCmd = []byte(`{"prefix":"fs dump","format":"json"}`)
	listFsCmd      = []byte(`{"prefix":"fs ls","format":"json"}`)
)

// ListVolumes return a list of volumes in this Ceph cluster.
//
// Similar To:
//
//	ceph fs volume ls
func (fsa *FSAdmin) ListVolumes() ([]string, error) {
	res := fsa.rawMgrCommand(listVolumesCmd)
	return parseListNames(res)
}

// FSPoolInfo contains the name of a file system as well as the metadata and
// data pools. Pool information is available by ID or by name.
type FSPoolInfo struct {
	Name           string   `json:"name"`
	MetadataPool   string   `json:"metadata_pool"`
	MetadataPoolID int      `json:"metadata_pool_id"`
	DataPools      []string `json:"data_pools"`
	DataPoolIDs    []int    `json:"data_pool_ids"`
}

// ListFileSystems lists file systems along with the pools occupied by those
// file systems.
//
// Similar To:
//
//	ceph fs ls
func (fsa *FSAdmin) ListFileSystems() ([]FSPoolInfo, error) {
	res := fsa.rawMonCommand(listFsCmd)
	return parseFsList(res)
}

func parseFsList(res response) ([]FSPoolInfo, error) {
	var listing []FSPoolInfo
	if err := res.NoStatus().Unmarshal(&listing).End(); err != nil {
		return nil, err
	}
	return listing, nil
}

// VolumeIdent contains a pair of file system identifying values: the volume
// name and the volume ID.
type VolumeIdent struct {
	Name string
	ID   int64
}

type cephFileSystem struct {
	ID     int64 `json:"id"`
	MDSMap struct {
		FSName string `json:"fs_name"`
	} `json:"mdsmap"`
}

type fsDump struct {
	FileSystems []cephFileSystem `json:"filesystems"`
}

const (
	dumpOkPrefix = "dumped fsmap epoch"
	dumpOkLen    = len(dumpOkPrefix)

	invalidTextualResponse = "this ceph version returns a non-parsable volume status response"
)

func parseDumpToIdents(res response) ([]VolumeIdent, error) {
	if !res.Ok() {
		return nil, res.End()
	}
	var dump fsDump
	if err := res.FilterPrefix(dumpOkPrefix).NoStatus().Unmarshal(&dump).End(); err != nil {
		return nil, err
	}
	// copy the dump json into the simpler enumeration list
	idents := make([]VolumeIdent, len(dump.FileSystems))
	for i := range dump.FileSystems {
		idents[i].ID = dump.FileSystems[i].ID
		idents[i].Name = dump.FileSystems[i].MDSMap.FSName
	}
	return idents, nil
}

// EnumerateVolumes returns a list of volume-name volume-id pairs.
func (fsa *FSAdmin) EnumerateVolumes() ([]VolumeIdent, error) {
	// We base our enumeration on the ceph fs dump json.  This may not be the
	// only way to do it, but it's the only one I know of currently. Because of
	// this and to keep our initial implementation simple, we expose our own
	// simplified type only, rather do a partial implementation of dump.
	return parseDumpToIdents(fsa.rawMonCommand(dumpVolumesCmd))
}

// VolumePool reports on the pool status for a CephFS volume.
type VolumePool struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	Available uint64 `json:"avail"`
	Used      uint64 `json:"used"`
}

// VolumeStatus reports various properties of a CephFS volume.
// TODO: Fill in.
type VolumeStatus struct {
	MDSVersion string       `json:"mds_version"`
	Pools      []VolumePool `json:"pools"`
}

type mdsVersionField struct {
	Version string
	Items   []struct {
		Version string `json:"version"`
	}
}

func (m *mdsVersionField) UnmarshalJSON(data []byte) (err error) {
	if err = json.Unmarshal(data, &m.Version); err == nil {
		return
	}
	return json.Unmarshal(data, &m.Items)
}

// volumeStatusResponse deals with the changing output of the mgr
// api json
type volumeStatusResponse struct {
	Pools      []VolumePool    `json:"pools"`
	MDSVersion mdsVersionField `json:"mds_version"`
}

func (v *volumeStatusResponse) volumeStatus() *VolumeStatus {
	vstatus := &VolumeStatus{}
	vstatus.Pools = v.Pools
	if v.MDSVersion.Version != "" {
		vstatus.MDSVersion = v.MDSVersion.Version
	} else if len(v.MDSVersion.Items) > 0 {
		vstatus.MDSVersion = v.MDSVersion.Items[0].Version
	}
	return vstatus
}

func parseVolumeStatus(res response) (*volumeStatusResponse, error) {
	var vs volumeStatusResponse
	res = res.NoStatus()
	if !res.Ok() {
		return nil, res.End()
	}
	res = res.Unmarshal(&vs)
	if !res.Ok() {
		if bytes.HasPrefix(res.Body(), []byte("ceph")) {
			return nil, NotImplementedError{
				Response: newResponse(res.Body(), invalidTextualResponse, res.Unwrap()),
			}
		}
		return nil, res.End()
	}
	return &vs, nil
}

// VolumeStatus returns a VolumeStatus object for the given volume name.
//
// Similar To:
//
//	ceph fs status cephfs <name>
func (fsa *FSAdmin) VolumeStatus(name string) (*VolumeStatus, error) {
	res := fsa.marshalMgrCommand(map[string]string{
		"fs":     name,
		"prefix": "fs status",
		"format": "json",
	})
	v, err := parseVolumeStatus(res)
	if err != nil {
		return nil, err
	}
	return v.volumeStatus(), nil
}
