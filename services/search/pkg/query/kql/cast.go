package kql

import (
	"fmt"
	"time"

	"github.com/araddon/dateparse"

	"github.com/owncloud/ocis/v2/services/search/pkg/query/ast"
)

func toNode[T ast.Node](in interface{}) (T, error) {
	var t T
	out, ok := in.(T)
	if !ok {
		return t, fmt.Errorf("can't convert '%T' to '%T'", in, t)
	}

	return out, nil
}

func toNodes[T ast.Node](in interface{}) ([]T, error) {
	switch v := in.(type) {
	case []T:
		return v, nil
	default:
		return nil, fmt.Errorf("can't convert '%T' to []ast.Node", in)
	}
}

func toString(in interface{}) (string, error) {
	switch v := in.(type) {
	case []byte:
		return string(v), nil
	case []interface{}:
		var str string

		for i := range v {
			sv, err := toString(v[i])
			if err != nil {
				return "", err
			}

			str += sv
		}

		return str, nil
	case string:
		return v, nil
	default:
		return "", fmt.Errorf("can't convert '%T' to string", v)
	}
}

func toTime(in interface{}) (time.Time, error) {
	ts, err := toString(in)
	if err != nil {
		return time.Time{}, err
	}

	return dateparse.ParseLocal(ts)
}
