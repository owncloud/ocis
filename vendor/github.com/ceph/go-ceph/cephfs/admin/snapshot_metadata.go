//go:build !(nautilus || octopus) && ceph_preview && ceph_pre_quincy
// +build !nautilus,!octopus,ceph_preview,ceph_pre_quincy

package admin

// GetSnapshotMetadata gets custom metadata on the subvolume snapshot in a
// volume belonging to an optional subvolume group based on provided key name.
//
// Similar To:
//  ceph fs subvolume snapshot metadata get <vol_name> <sub_name> <snap_name> <key_name> [--group_name <subvol_group_name>]
func (fsa *FSAdmin) GetSnapshotMetadata(volume, group, subvolume, snapname, key string) (string, error) {
	m := map[string]string{
		"prefix":    "fs subvolume snapshot metadata get",
		"format":    "json",
		"vol_name":  volume,
		"sub_name":  subvolume,
		"snap_name": snapname,
		"key_name":  key,
	}

	if group != NoGroup {
		m["group_name"] = group
	}

	return parsePathResponse(fsa.marshalMgrCommand(m))
}

// SetSnapshotMetadata sets custom metadata on the subvolume snapshot in a
// volume belonging to an optional subvolume group as a key-value pair.
//
// Similar To:
//  ceph fs subvolume snapshot metadata set <vol_name> <sub_name> <snap_name> <key_name> <value> [--group_name <subvol_group_name>]
func (fsa *FSAdmin) SetSnapshotMetadata(volume, group, subvolume, snapname, key, value string) error {
	m := map[string]string{
		"prefix":    "fs subvolume snapshot metadata set",
		"format":    "json",
		"vol_name":  volume,
		"sub_name":  subvolume,
		"snap_name": snapname,
		"key_name":  key,
		"value":     value,
	}

	if group != NoGroup {
		m["group_name"] = group
	}

	return fsa.marshalMgrCommand(m).NoData().End()
}

// RemoveSnapshotMetadata removes custom metadata set on the subvolume
// snapshot in a volume belonging to an optional subvolume group using the
// metadata key.
//
// Similar To:
//  ceph fs subvolume snapshot metadata rm <vol_name> <sub_name> <snap_name> <key_name> [--group_name <subvol_group_name>]
func (fsa *FSAdmin) RemoveSnapshotMetadata(volume, group, subvolume, snapname, key string) error {
	return fsa.rmSubVolumeSnapShotMetadata(volume, group, subvolume, snapname, key, commonRmFlags{})
}

// ForceRemoveSnapshotMetadata attempt to forcefully remove custom metadata
// set on the subvolume snapshot in a volume belonging to an optional
// subvolume group using the metadata key.
//
// Similar To:
//  ceph fs subvolume snapshot metadata rm <vol_name> <sub_name> <snap_name> <key_name> [--group_name <subvol_group_name>] --force
func (fsa *FSAdmin) ForceRemoveSnapshotMetadata(volume, group, subvolume, snapname, key string) error {
	return fsa.rmSubVolumeSnapShotMetadata(volume, group, subvolume, snapname, key, commonRmFlags{force: true})
}

func (fsa *FSAdmin) rmSubVolumeSnapShotMetadata(volume, group, subvolume, snapname, key string, o commonRmFlags) error {
	m := map[string]string{
		"prefix":    "fs subvolume snapshot metadata rm",
		"format":    "json",
		"vol_name":  volume,
		"sub_name":  subvolume,
		"snap_name": snapname,
		"key_name":  key,
	}

	if group != NoGroup {
		m["group_name"] = group
	}

	return fsa.marshalMgrCommand(mergeFlags(m, o)).NoData().End()
}

// ListSnapshotMetadata lists custom metadata (key-value pairs) set on the subvolume
// snapshot in a volume belonging to an optional subvolume group.
//
// Similar To:
//  ceph fs subvolume snapshot metadata ls <vol_name> <sub_name> <snap_name> [--group_name <subvol_group_name>]
func (fsa *FSAdmin) ListSnapshotMetadata(volume, group, subvolume, snapname string) (map[string]string, error) {
	m := map[string]string{
		"prefix":    "fs subvolume snapshot metadata ls",
		"format":    "json",
		"vol_name":  volume,
		"sub_name":  subvolume,
		"snap_name": snapname,
	}

	if group != NoGroup {
		m["group_name"] = group
	}

	return parseListKeyValues(fsa.marshalMgrCommand(m))
}
