package admin

// this is the internal type used to create JSON for ceph.
// See SubVolumeGroupOptions for the type that users of the library
// interact with.
// note that the ceph json takes mode as a string.
type subVolumeGroupFields struct {
	Prefix     string `json:"prefix"`
	Format     string `json:"format"`
	VolName    string `json:"vol_name"`
	GroupName  string `json:"group_name"`
	Uid        int    `json:"uid,omitempty"`
	Gid        int    `json:"gid,omitempty"`
	Mode       string `json:"mode,omitempty"`
	PoolLayout string `json:"pool_layout,omitempty"`
}

// SubVolumeGroupOptions are used to specify optional, non-identifying, values
// to be used when creating a new subvolume group.
type SubVolumeGroupOptions struct {
	Uid        int
	Gid        int
	Mode       int
	PoolLayout string
}

func (s *SubVolumeGroupOptions) toFields(v, g string) *subVolumeGroupFields {
	return &subVolumeGroupFields{
		Prefix:     "fs subvolumegroup create",
		Format:     "json",
		VolName:    v,
		GroupName:  g,
		Uid:        s.Uid,
		Gid:        s.Gid,
		Mode:       modeString(s.Mode, false),
		PoolLayout: s.PoolLayout,
	}
}

// CreateSubVolumeGroup sends a request to create a subvolume group in a volume.
//
// Similar To:
//
//	ceph fs subvolumegroup create <volume> <group_name>  ...
func (fsa *FSAdmin) CreateSubVolumeGroup(volume, name string, o *SubVolumeGroupOptions) error {
	if o == nil {
		o = &SubVolumeGroupOptions{}
	}
	res := fsa.marshalMgrCommand(o.toFields(volume, name))
	return res.NoData().End()
}

// ListSubVolumeGroups returns a list of subvolume groups belonging to the
// specified volume.
//
// Similar To:
//
//	ceph fs subvolumegroup ls cephfs <volume>
func (fsa *FSAdmin) ListSubVolumeGroups(volume string) ([]string, error) {
	res := fsa.marshalMgrCommand(map[string]string{
		"prefix":   "fs subvolumegroup ls",
		"vol_name": volume,
		"format":   "json",
	})
	return parseListNames(res)
}

// RemoveSubVolumeGroup will delete a subvolume group in a volume.
// Similar To:
//
//	ceph fs subvolumegroup rm <volume> <group_name>
func (fsa *FSAdmin) RemoveSubVolumeGroup(volume, name string) error {
	return fsa.rmSubVolumeGroup(volume, name, commonRmFlags{})
}

// ForceRemoveSubVolumeGroup will delete a subvolume group in a volume.
// Similar To:
//
//	ceph fs subvolumegroup rm <volume> <group_name> --force
func (fsa *FSAdmin) ForceRemoveSubVolumeGroup(volume, name string) error {
	return fsa.rmSubVolumeGroup(volume, name, commonRmFlags{force: true})
}

func (fsa *FSAdmin) rmSubVolumeGroup(volume, name string, o commonRmFlags) error {
	res := fsa.marshalMgrCommand(mergeFlags(map[string]string{
		"prefix":     "fs subvolumegroup rm",
		"vol_name":   volume,
		"group_name": name,
		"format":     "json",
	}, o))
	return res.NoData().End()
}

// SubVolumeGroupPath returns the path to the subvolume from the root of the
// file system.
//
// Similar To:
//
//	ceph fs subvolumegroup getpath <volume> <group_name>
func (fsa *FSAdmin) SubVolumeGroupPath(volume, name string) (string, error) {
	m := map[string]string{
		"prefix":     "fs subvolumegroup getpath",
		"vol_name":   volume,
		"group_name": name,
		// ceph doesn't respond in json for this cmd (even if you ask)
	}
	return parsePathResponse(fsa.marshalMgrCommand(m))
}
