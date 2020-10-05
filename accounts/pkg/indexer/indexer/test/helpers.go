package test

import (
	"io/ioutil"
	"reflect"
	"testing"
)

func CreateTmpDir(t *testing.T) string {
	name, err := ioutil.TempDir("/var/tmp", "testfiles-*")
	if err != nil {
		t.Fatal(err)
	}

	return name
}

func ValueOf(v interface{}, field string) string {
	r := reflect.ValueOf(v)
	f := reflect.Indirect(r).FieldByName(field)

	return f.String()
}
