package msgpack

import (
	"fmt"

	"github.com/shamaton/msgpack/v2/def"
	"github.com/shamaton/msgpack/v2/ext"
	"github.com/shamaton/msgpack/v2/internal/decoding"
	"github.com/shamaton/msgpack/v2/internal/encoding"
)

// StructAsArray is encoding option.
// If this option sets true, default encoding sets to array-format.
var StructAsArray = false

// Marshal returns the MessagePack-encoded byte array of v.
func Marshal(v interface{}) ([]byte, error) {
	return encoding.Encode(v, StructAsArray)
}

// Unmarshal analyzes the MessagePack-encoded data and stores
// the result into the pointer of v.
func Unmarshal(data []byte, v interface{}) error {
	return decoding.Decode(data, v, StructAsArray)
}

// AddExtCoder adds encoders for extension types.
func AddExtCoder(e ext.Encoder, d ext.Decoder) error {
	if e.Code() != d.Code() {
		return fmt.Errorf("code different %d:%d", e.Code(), d.Code())
	}
	encoding.AddExtEncoder(e)
	decoding.AddExtDecoder(d)
	return nil
}

// RemoveExtCoder removes encoders for extension types.
func RemoveExtCoder(e ext.Encoder, d ext.Decoder) error {
	if e.Code() != d.Code() {
		return fmt.Errorf("code different %d:%d", e.Code(), d.Code())
	}
	encoding.RemoveExtEncoder(e)
	decoding.RemoveExtDecoder(d)
	return nil
}

// SetComplexTypeCode sets def.complexTypeCode
func SetComplexTypeCode(code int8) {
	def.SetComplexTypeCode(code)
}
