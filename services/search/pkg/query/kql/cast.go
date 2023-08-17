package kql

import (
	"fmt"
)

func toIfaceSlice(in interface{}) []interface{} {
	if in == nil {
		return nil
	}
	return in.([]interface{})
}

func toString(in interface{}) (string, error) {
	switch v := in.(type) {
	case []byte:
		return string(v), nil
	case []interface{}:
		str := ""
		for _, i := range v {
			j := i.([]uint8)
			str += string(j[0])
		}
		return str, nil
	case string:
		return v, nil
	default:
		return "", fmt.Errorf("can't convert '%T' to string ", v)
	}
}
