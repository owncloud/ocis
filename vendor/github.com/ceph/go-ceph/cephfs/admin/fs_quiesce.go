package admin

import "fmt"

// fixedPointFloat is a custom type that implements the MarshalJSON interface.
// This is used to format float64 values to two decimal places.
// By default these get converted to integers in the JSON output and
// fail the command.
type fixedPointFloat float64

// MarshalJSON provides a custom implementation for the JSON marshalling
// of fixedPointFloat. It formats the float to two decimal places.
func (fpf fixedPointFloat) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%.2f", float64(fpf))), nil
}

// fSQuiesceFields is the internal type used to create JSON for ceph.
// See FSQuiesceOptions for the type that users of the library
// interact with.
type fSQuiesceFields struct {
	Prefix     string          `json:"prefix"`
	VolName    string          `json:"vol_name"`
	GroupName  string          `json:"group_name,omitempty"`
	Members    []string        `json:"members,omitempty"`
	SetId      string          `json:"set_id,omitempty"`
	Timeout    fixedPointFloat `json:"timeout,omitempty"`
	Expiration fixedPointFloat `json:"expiration,omitempty"`
	AwaitFor   fixedPointFloat `json:"await_for,omitempty"`
	Await      bool            `json:"await,omitempty"`
	IfVersion  int             `json:"if_version,omitempty"`
	Include    bool            `json:"include,omitempty"`
	Exclude    bool            `json:"exclude,omitempty"`
	Reset      bool            `json:"reset,omitempty"`
	Release    bool            `json:"release,omitempty"`
	Query      bool            `json:"query,omitempty"`
	All        bool            `json:"all,omitempty"`
	Cancel     bool            `json:"cancel,omitempty"`
}

// FSQuiesceOptions are used to specify optional, non-identifying, values
// to be used when quiescing a cephfs volume.
type FSQuiesceOptions struct {
	Timeout    float64
	Expiration float64
	AwaitFor   float64
	Await      bool
	IfVersion  int
	Include    bool
	Exclude    bool
	Reset      bool
	Release    bool
	Query      bool
	All        bool
	Cancel     bool
}

// toFields is used to convert the FSQuiesceOptions to the internal
// fSQuiesceFields type.
func (o *FSQuiesceOptions) toFields(volume, group string, subvolumes []string, setId string) *fSQuiesceFields {
	return &fSQuiesceFields{
		Prefix:     "fs quiesce",
		VolName:    volume,
		GroupName:  group,
		Members:    subvolumes,
		SetId:      setId,
		Timeout:    fixedPointFloat(o.Timeout),
		Expiration: fixedPointFloat(o.Expiration),
		AwaitFor:   fixedPointFloat(o.AwaitFor),
		Await:      o.Await,
		IfVersion:  o.IfVersion,
		Include:    o.Include,
		Exclude:    o.Exclude,
		Reset:      o.Reset,
		Release:    o.Release,
		Query:      o.Query,
		All:        o.All,
		Cancel:     o.Cancel,
	}
}

// QuiesceState is used to report the state of a quiesced fs volume.
type QuiesceState struct {
	Name string  `json:"name"`
	Age  float64 `json:"age"`
}

// QuiesceInfoMember is used to report the state of a quiesced fs volume.
// This is part of sets members object array in the json.
type QuiesceInfoMember struct {
	Excluded bool         `json:"excluded"`
	State    QuiesceState `json:"state"`
}

// QuiesceInfo reports various informational values about a quiesced volume.
// This is returned as sets object array in the json.
type QuiesceInfo struct {
	Version    int                          `json:"version"`
	AgeRef     float64                      `json:"age_ref"`
	State      QuiesceState                 `json:"state"`
	Timeout    float64                      `json:"timeout"`
	Expiration float64                      `json:"expiration"`
	Members    map[string]QuiesceInfoMember `json:"members"`
}

// FSQuiesceInfo reports various informational values about quiesced volumes.
type FSQuiesceInfo struct {
	Epoch      int                    `json:"epoch"`
	SetVersion int                    `json:"set_version"`
	Sets       map[string]QuiesceInfo `json:"sets"`
}

// parseFSQuiesceInfo is used to parse the response from the quiesce command. It returns a FSQuiesceInfo object.
func parseFSQuiesceInfo(res response) (*FSQuiesceInfo, error) {
	var info FSQuiesceInfo
	if err := res.NoStatus().Unmarshal(&info).End(); err != nil {
		return nil, err
	}
	return &info, nil
}

// FSQuiesce will quiesce the specified subvolumes in a volume.
// Quiescing a fs will prevent new writes to the subvolumes.
// Similar To:
//
// ceph fs quiesce <volume>
func (fsa *FSAdmin) FSQuiesce(volume, group string, subvolumes []string, setId string, o *FSQuiesceOptions) (*FSQuiesceInfo, error) {
	if o == nil {
		o = &FSQuiesceOptions{}
	}
	f := o.toFields(volume, group, subvolumes, setId)

	return parseFSQuiesceInfo(fsa.marshalMgrCommand(f))
}
