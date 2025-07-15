package ext

import (
	"reflect"
)

type StreamDecoder interface {
	Code() int8
	IsType(code byte, innerType int8, dataLength int) bool
	ToValue(code byte, data []byte, k reflect.Kind) (any, error)
}
