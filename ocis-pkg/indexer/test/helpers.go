package test

import (
	"errors"
	"io/ioutil"
	"path"
	"reflect"
	"strings"
)

// CreateTmpDir creates a temporary dir for tests data.
func CreateTmpDir() (string, error) {
	name, err := ioutil.TempDir("/tmp", "testfiles-")
	if err != nil {
		return "", err
	}

	return name, nil
}

// ValueOf gets the value of a type v on a given field <field>.
func ValueOf(v interface{}, field string) string {
	r := reflect.ValueOf(v)
	f := reflect.Indirect(r).FieldByName(field)

	return f.String()
}

func getType(v interface{}) (reflect.Value, error) {
	rv := reflect.ValueOf(v)
	for rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Interface {
		rv = rv.Elem()
	}
	if !rv.IsValid() {
		return reflect.Value{}, errors.New("failed to read value via reflection")
	}

	return rv, nil
}

// GetTypeFQN formats a valid name from a type <t>. This is a duplication of the already existing function in the
// indexer package, but since there is a circular dependency we chose to duplicate it.
func GetTypeFQN(t interface{}) string {
	typ, _ := getType(t)
	typeName := path.Join(typ.Type().PkgPath(), typ.Type().Name())
	typeName = strings.ReplaceAll(typeName, "/", ".")
	return typeName
}
