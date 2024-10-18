//go:build !(nautilus || octopus)
// +build !nautilus,!octopus

package admin

// PoolInfo reports various properties of a pool.
type PoolInfo struct {
	Available int    `json:"avail"`
	Name      string `json:"name"`
	Used      int    `json:"used"`
}

// PoolType indicates the type of pool related to a volume.
type PoolType struct {
	DataPool     []PoolInfo `json:"data"`
	MetadataPool []PoolInfo `json:"metadata"`
}

// VolInfo holds various informational values about a volume.
type VolInfo struct {
	MonAddrs          []string `json:"mon_addrs"`
	PendingSubvolDels int      `json:"pending_subvolume_deletions"`
	Pools             PoolType `json:"pools"`
	UsedSize          int      `json:"used_size"`
}

func parseVolumeInfo(res response) (*VolInfo, error) {
	var info VolInfo
	if err := res.NoStatus().Unmarshal(&info).End(); err != nil {
		return nil, err
	}
	return &info, nil
}

// FetchVolumeInfo fetches the information of a CephFS volume.
//
// Similar To:
//
//	ceph fs volume info <vol_name>
func (fsa *FSAdmin) FetchVolumeInfo(volume string) (*VolInfo, error) {
	m := map[string]string{
		"prefix":   "fs volume info",
		"vol_name": volume,
		"format":   "json",
	}

	return parseVolumeInfo(fsa.marshalMgrCommand(m))
}
