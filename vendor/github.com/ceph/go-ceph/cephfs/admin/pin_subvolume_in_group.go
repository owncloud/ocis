//go:build !nautilus && ceph_preview

package admin

// PinSubVolumeInGroup pins a subvolume in a volume and optional subvolume group
// to ranks according to policies. A valid pin setting value depends on the type
// of pin as described in the docs from
// https://docs.ceph.com/en/latest/cephfs/multimds/#cephfs-pinning and
// https://docs.ceph.com/en/latest/cephfs/multimds/#setting-subtree-partitioning-policies
//
// Similar To:
//
//	ceph fs subvolume pin <vol_name> <sub_name> <pin_type> <pin_setting> [--group_name=<group>]
func (fsa *FSAdmin) PinSubVolumeInGroup(volume, group, subvolume, pintype, pinsetting string) (string, error) {
	m := map[string]string{
		"prefix":      "fs subvolume pin",
		"format":      "json",
		"vol_name":    volume,
		"sub_name":    subvolume,
		"pin_type":    pintype,
		"pin_setting": pinsetting,
	}
	if group != NoGroup {
		m["group_name"] = group
	}

	return parsePathResponse(fsa.marshalMgrCommand(m))
}
