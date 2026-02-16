package admin

import (
	"encoding/json"
	"fmt"
)

// ByteCount represents the size of a volume in bytes.
type ByteCount uint64

// SI byte size constants. keep these private for now.
const (
	kibiByte ByteCount = 1024
	mebiByte           = 1024 * kibiByte
	gibiByte           = 1024 * mebiByte
	tebiByte           = 1024 * gibiByte
)

// resizeValue returns a size value as a string, as needed by the subvolume
// resize command json.
func (bc ByteCount) resizeValue() string {
	return uint64String(uint64(bc))
}

// QuotaSize interface values can be used to change the size of a volume.
type QuotaSize interface {
	resizeValue() string
}

// specialSize is a custom non-numeric quota size value.
type specialSize string

// resizeValue for a specialSize returns the original string value.
func (s specialSize) resizeValue() string {
	return string(s)
}

// Infinite is a special QuotaSize value that can be used to clear size limits
// on a subvolume.
const Infinite = specialSize("infinite")

// quotaSizePlaceholder types are helpful to extract QuotaSize typed values
// from JSON responses.
type quotaSizePlaceholder struct {
	Value QuotaSize
}

func (p *quotaSizePlaceholder) UnmarshalJSON(b []byte) error {
	var val interface{}
	if err := json.Unmarshal(b, &val); err != nil {
		return err
	}
	switch v := val.(type) {
	case string:
		if v == string(Infinite) {
			p.Value = Infinite
		} else {
			return fmt.Errorf("quota size: invalid string value: %q", v)
		}
	case float64:
		p.Value = ByteCount(v)
	default:
		return fmt.Errorf("quota size: invalid type, string or number required: %v (%T)", val, val)
	}
	return nil
}
