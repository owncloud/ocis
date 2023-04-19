package admin

import (
	"encoding/json"
	"time"
)

// golang's date parsing approach is rather bizarre
var cephTSLayout = "2006-01-02 15:04:05"

// TimeStamp abstracts some of the details about date+time stamps
// returned by ceph via JSON.
type TimeStamp struct {
	time.Time
}

// String returns a string representing the date+time as presented
// by ceph.
func (ts TimeStamp) String() string {
	return ts.Format(cephTSLayout)
}

// UnmarshalJSON implements the json Unmarshaler interface.
func (ts *TimeStamp) UnmarshalJSON(b []byte) error {
	var raw string
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	// AFAICT, ceph always returns the time in UTC so Parse, as opposed to
	// ParseInLocation, is appropriate here.
	t, err := time.Parse(cephTSLayout, raw)
	if err != nil {
		return err
	}
	*ts = TimeStamp{t}
	return nil
}
