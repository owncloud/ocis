package msgpack

import (
	"github.com/shamaton/msgpack/v2/internal/encoding"
)

// MarshalAsMap encodes data as map format.
// This is the same thing that StructAsArray sets false.
func MarshalAsMap(v interface{}) ([]byte, error) {
	return encoding.Encode(v, false)
}

// MarshalAsArray encodes data as array format.
// This is the same thing that StructAsArray sets true.
func MarshalAsArray(v interface{}) ([]byte, error) {
	return encoding.Encode(v, true)
}
