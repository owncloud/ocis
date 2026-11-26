package decodingutil

import (
	"errors"
	"fmt"
	"math"
	"reflect"

	"github.com/shamaton/msgpack/v2/def"
)

// ErrDeclaredLengthTooLarge is returned when a declared 32-bit length cannot
// be represented as a non-negative int (32-bit platforms).
var ErrDeclaredLengthTooLarge = errors.New("declared length is too large")

// LengthFromUint32 converts a MessagePack-declared 32-bit length to int,
// rejecting values that are not representable as a non-negative int.
func LengthFromUint32(u uint32) (int, error) {
	if int64(u) > int64(math.MaxInt) {
		return 0, fmt.Errorf("%w: %d", ErrDeclaredLengthTooLarge, u)
	}
	return int(u), nil // #nosec G115 -- checked above
}

// MaxDepth limits the nesting of arrays, maps and structs while decoding.
// Without it, deeply nested input exhausts the goroutine stack and the runtime
// aborts with a fatal error that cannot be recovered.
const MaxDepth = 10000

// DepthError reports that the nesting limit was exceeded.
func DepthError() error {
	return fmt.Errorf("%w (max %d)", def.ErrExceededMaxDepth, MaxDepth)
}

// IsExtCode reports whether code starts an ext value. Map keys decoded from
// ext values come from user-provided decoders and may be unhashable, while
// other non-container codes always yield hashable values.
func IsExtCode(code byte) bool {
	switch code {
	case def.Fixext1, def.Fixext2, def.Fixext4, def.Fixext8, def.Fixext16,
		def.Ext8, def.Ext16, def.Ext32:
		return true
	}
	return false
}

// CheckMapKey reports an error when key would panic as a map key.
func CheckMapKey(key interface{}) error {
	if key == nil {
		return nil
	}
	if !reflect.ValueOf(key).Comparable() {
		return fmt.Errorf("%w. type: %T", def.ErrUnhashableMapKey, key)
	}
	return nil
}

// MapKeyNeedsCheck reports whether a concrete map key type t may still hold
// unhashable values at runtime (e.g. a struct with an interface field).
// Interface key types are handled separately by the decoders.
func MapKeyNeedsCheck(t reflect.Type) bool {
	switch t.Kind() {
	case reflect.Struct, reflect.Array:
		return true
	}
	return false
}

// CheckMapKeyValue is the reflect.Value variant of CheckMapKey.
func CheckMapKeyValue(key reflect.Value) error {
	if !key.Comparable() {
		return fmt.Errorf("%w. type: %v", def.ErrUnhashableMapKey, key.Type())
	}
	return nil
}

// EmptyInterfaceType is the type of interface{}. Decoders compare against
// it to skip SetInterface on the common path.
var EmptyInterfaceType = reflect.TypeOf((*interface{})(nil)).Elem()

// SetInterface stores a decoded value into the interface destination rv,
// leaving rv untouched for nil. Non-empty interfaces only accept values that
// implement them; reflect.Set would panic otherwise.
func SetInterface(rv reflect.Value, v interface{}) error {
	if v == nil {
		return nil
	}
	vv := reflect.ValueOf(v)
	if t := rv.Type(); !vv.Type().Implements(t) {
		return fmt.Errorf("%w %v decoding as %v", def.ErrCanNotDecode, vv.Type(), t)
	}
	rv.Set(vv)
	return nil
}
