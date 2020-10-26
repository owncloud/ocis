package disk

import (
	"os"
	"reflect"
	"sort"
	"strconv"
)

var (
	validKinds = []reflect.Kind{
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
	}
)

// verifies an autoincrement field kind on the target struct.
func isValidKind(k reflect.Kind) bool {
	for _, v := range validKinds {
		if k == v {
			return true
		}
	}
	return false
}

func getKind(i interface{}, field string) (reflect.Kind, error) {
	r := reflect.ValueOf(i)
	return reflect.Indirect(r).FieldByName(field).Kind(), nil
}

// readDir is an implementation of os.ReadDir but with different sorting.
func readDir(dirname string) ([]os.FileInfo, error) {
	f, err := os.Open(dirname)
	if err != nil {
		return nil, err
	}
	list, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		return nil, err
	}
	sort.Slice(list, func(i, j int) bool {
		a, _ := strconv.Atoi(list[i].Name())
		b, _ := strconv.Atoi(list[j].Name())
		return a < b
	})
	return list, nil
}
