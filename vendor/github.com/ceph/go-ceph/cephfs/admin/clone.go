package admin

import (
	"strings"
)

const notProtectedSuffix = "is not protected"

// NotProtectedError error values will be returned by CloneSubVolumeSnapshot in
// the case that the source snapshot needs to be protected but is not.  The
// requirement for a snapshot to be protected prior to cloning varies by Ceph
// version.
type NotProtectedError struct {
	response
}

// CloneOptions are used to specify optional values to be used when creating a
// new subvolume clone.
type CloneOptions struct {
	TargetGroup string
	PoolLayout  string
}

// CloneSubVolumeSnapshot clones the specified snapshot from the subvolume.
// The group, subvolume, and snapshot parameters specify the source for the
// clone, and only the source.  Additional properties of the clone, such as the
// subvolume group that the clone will be created in and the pool layout may be
// specified using the clone options parameter.
//
// Similar To:
//
//	ceph fs subvolume snapshot clone <volume> --group_name=<group> <subvolume> <snapshot> <name> [...]
func (fsa *FSAdmin) CloneSubVolumeSnapshot(volume, group, subvolume, snapshot, name string, o *CloneOptions) error {
	m := map[string]string{
		"prefix":          "fs subvolume snapshot clone",
		"vol_name":        volume,
		"sub_name":        subvolume,
		"snap_name":       snapshot,
		"target_sub_name": name,
		"format":          "json",
	}
	if group != NoGroup {
		m["group_name"] = group
	}
	if o != nil && o.TargetGroup != NoGroup {
		m["target_group_name"] = group
	}
	if o != nil && o.PoolLayout != "" {
		m["pool_layout"] = o.PoolLayout
	}
	return checkCloneResponse(fsa.marshalMgrCommand(m))
}

func checkCloneResponse(res response) error {
	if strings.HasSuffix(res.Status(), notProtectedSuffix) {
		return NotProtectedError{response: res}
	}
	return res.NoData().End()
}

// CloneState is used to define constant values used to determine the state of
// a clone.
type CloneState string

const (
	// ClonePending is the state of a pending clone.
	ClonePending = CloneState("pending")
	// CloneInProgress is the state of a clone in progress.
	CloneInProgress = CloneState("in-progress")
	// CloneComplete is the state of a complete clone.
	CloneComplete = CloneState("complete")
	// CloneFailed is the state of a failed clone.
	CloneFailed = CloneState("failed")
)

// CloneSource contains values indicating the source of a clone.
type CloneSource struct {
	Volume    string `json:"volume"`
	Group     string `json:"group"`
	SubVolume string `json:"subvolume"`
	Snapshot  string `json:"snapshot"`
}

// CloneStatus reports on the status of a subvolume clone.
type CloneStatus struct {
	State  CloneState  `json:"state"`
	Source CloneSource `json:"source"`

	// failure can be obtained through .GetFailure()
	failure *CloneFailure
}

// CloneFailure reports details of a failure after a subvolume clone failed.
type CloneFailure struct {
	Errno  string `json:"errno"`
	ErrStr string `json:"errstr"`
}

type cloneStatusWrapper struct {
	Status  CloneStatus  `json:"status"`
	Failure CloneFailure `json:"failure"`
}

func parseCloneStatus(res response) (*CloneStatus, error) {
	var status cloneStatusWrapper
	if err := res.NoStatus().Unmarshal(&status).End(); err != nil {
		return nil, err
	}
	if status.Failure.Errno != "" || status.Failure.ErrStr != "" {
		status.Status.failure = &status.Failure
	}
	return &status.Status, nil
}

// CloneStatus returns data reporting the status of a subvolume clone.
//
// Similar To:
//
//	ceph fs clone status <volume> --group_name=<group> <clone>
func (fsa *FSAdmin) CloneStatus(volume, group, clone string) (*CloneStatus, error) {
	m := map[string]string{
		"prefix":     "fs clone status",
		"vol_name":   volume,
		"clone_name": clone,
		"format":     "json",
	}
	if group != NoGroup {
		m["group_name"] = group
	}
	return parseCloneStatus(fsa.marshalMgrCommand(m))
}

// CancelClone stops the background processes that populate a clone.
// CancelClone does not delete the clone.
//
// Similar To:
//
//	ceph fs clone cancel <volume> --group_name=<group> <clone>
func (fsa *FSAdmin) CancelClone(volume, group, clone string) error {
	m := map[string]string{
		"prefix":     "fs clone cancel",
		"vol_name":   volume,
		"clone_name": clone,
		"format":     "json",
	}
	if group != NoGroup {
		m["group_name"] = group
	}
	return fsa.marshalMgrCommand(m).NoData().End()
}
