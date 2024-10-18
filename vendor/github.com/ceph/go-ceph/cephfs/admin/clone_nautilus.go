package admin

// GetFailure returns details about the CloneStatus when in CloneFailed state.
//
// Similar To:
//
//	Reading the .failure object from the JSON returned by "ceph fs subvolume
//	snapshot clone"
func (cs *CloneStatus) GetFailure() *CloneFailure {
	return cs.failure
}
