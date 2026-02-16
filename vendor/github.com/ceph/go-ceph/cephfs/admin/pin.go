//go:build !nautilus
// +build !nautilus

package admin

// PinSubVolume pins subvolume to ranks according to policies. A valid pin
// setting value depends on the type of pin as described in the docs from
// https://docs.ceph.com/en/latest/cephfs/multimds/#cephfs-pinning and
// https://docs.ceph.com/en/latest/cephfs/multimds/#setting-subtree-partitioning-policies
//
// Similar To:
//
//	ceph fs subvolume pin <vol_name> <sub_name> <pin_type> <pin_setting>
func (fsa *FSAdmin) PinSubVolume(volume, subvolume, pintype, pinsetting string) (string, error) {
	m := map[string]string{
		"prefix":      "fs subvolume pin",
		"format":      "json",
		"vol_name":    volume,
		"sub_name":    subvolume,
		"pin_type":    pintype,
		"pin_setting": pinsetting,
	}

	return parsePathResponse(fsa.marshalMgrCommand(m))
}

// PinSubVolumeGroup pins subvolume to ranks according to policies. A valid pin
// setting value depends on the type of pin as described in the docs from
// https://docs.ceph.com/en/latest/cephfs/multimds/#cephfs-pinning and
// https://docs.ceph.com/en/latest/cephfs/multimds/#setting-subtree-partitioning-policies
//
// Similar To:
//
//	ceph fs subvolumegroup pin <vol_name> <group_name> <pin_type> <pin_setting>
func (fsa *FSAdmin) PinSubVolumeGroup(volume, group, pintype, pinsetting string) (string, error) {
	m := map[string]string{
		"prefix":      "fs subvolumegroup pin",
		"format":      "json",
		"vol_name":    volume,
		"group_name":  group,
		"pin_type":    pintype,
		"pin_setting": pinsetting,
	}

	return parsePathResponse(fsa.marshalMgrCommand(m))
}
