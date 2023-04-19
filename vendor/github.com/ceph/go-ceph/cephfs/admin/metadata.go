//go:build !(nautilus || octopus) && ceph_preview && ceph_pre_quincy
// +build !nautilus,!octopus,ceph_preview,ceph_pre_quincy

package admin

// GetMetadata gets custom metadata on the subvolume in a volume belonging to
// an optional subvolume group based on provided key name.
//
// Similar To:
//  ceph fs subvolume metadata get <vol_name> <sub_name> <key_name> [--group_name <subvol_group_name>]
func (fsa *FSAdmin) GetMetadata(volume, group, subvolume, key string) (string, error) {
	m := map[string]string{
		"prefix":   "fs subvolume metadata get",
		"format":   "json",
		"vol_name": volume,
		"sub_name": subvolume,
		"key_name": key,
	}

	if group != NoGroup {
		m["group_name"] = group
	}

	return parsePathResponse(fsa.marshalMgrCommand(m))
}

// SetMetadata sets custom metadata on the subvolume in a volume belonging to
// an optional subvolume group as a key-value pair.
//
// Similar To:
//  ceph fs subvolume metadata set <vol_name> <sub_name> <key_name> <value> [--group_name <subvol_group_name>]
func (fsa *FSAdmin) SetMetadata(volume, group, subvolume, key, value string) error {
	m := map[string]string{
		"prefix":   "fs subvolume metadata set",
		"format":   "json",
		"vol_name": volume,
		"sub_name": subvolume,
		"key_name": key,
		"value":    value,
	}

	if group != NoGroup {
		m["group_name"] = group
	}

	return fsa.marshalMgrCommand(m).NoData().End()
}

// RemoveMetadata removes custom metadata set on the subvolume in a volume
// belonging to an optional subvolume group using the metadata key.
//
// Similar To:
//  ceph fs subvolume metadata rm <vol_name> <sub_name> <key_name> [--group_name <subvol_group_name>]
func (fsa *FSAdmin) RemoveMetadata(volume, group, subvolume, key string) error {
	return fsa.rmSubVolumeMetadata(volume, group, subvolume, key, commonRmFlags{})
}

// ForceRemoveMetadata attempt to forcefully remove custom metadata set on
// the subvolume in a volume belonging to an optional subvolume group using
// the metadata key.
//
// Similar To:
//  ceph fs subvolume metadata rm <vol_name> <sub_name> <key_name> [--group_name <subvol_group_name>] --force
func (fsa *FSAdmin) ForceRemoveMetadata(volume, group, subvolume, key string) error {
	return fsa.rmSubVolumeMetadata(volume, group, subvolume, key, commonRmFlags{force: true})
}

func (fsa *FSAdmin) rmSubVolumeMetadata(volume, group, subvolume, key string, o commonRmFlags) error {
	m := map[string]string{
		"prefix":   "fs subvolume metadata rm",
		"format":   "json",
		"vol_name": volume,
		"sub_name": subvolume,
		"key_name": key,
	}

	if group != NoGroup {
		m["group_name"] = group
	}

	return fsa.marshalMgrCommand(mergeFlags(m, o)).NoData().End()
}

// ListMetadata lists custom metadata (key-value pairs) set on the subvolume
// in a volume belonging to an optional subvolume group.
//
// Similar To:
//  ceph fs subvolume metadata ls <vol_name> <sub_name> [--group_name <subvol_group_name>]
func (fsa *FSAdmin) ListMetadata(volume, group, subvolume string) (map[string]string, error) {
	m := map[string]string{
		"prefix":   "fs subvolume metadata ls",
		"format":   "json",
		"vol_name": volume,
		"sub_name": subvolume,
	}

	if group != NoGroup {
		m["group_name"] = group
	}

	return parseListKeyValues(fsa.marshalMgrCommand(m))
}
