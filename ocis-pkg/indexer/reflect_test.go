package indexer

import (
	"testing"
)

func Test_getTypeFQN(t *testing.T) {
	type someT struct{}

	type args struct {
		t interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "ByValue", args: args{&someT{}}, want: "github.com.owncloud.ocis.ocis-pkg.indexer.someT"},
		{name: "ByRef", args: args{someT{}}, want: "github.com.owncloud.ocis.ocis-pkg.indexer.someT"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getTypeFQN(tt.args.t); got != tt.want {
				t.Errorf("getTypeFQN() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_valueOf(t *testing.T) {
	type someT struct {
		val string
	}
	type args struct {
		v     interface{}
		field string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "ByValue", args: args{v: someT{val: "hello"}, field: "val"}, want: "hello"},
		{name: "ByRef", args: args{v: &someT{val: "hello"}, field: "val"}, want: "hello"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := valueOf(tt.args.v, tt.args.field); got != tt.want {
				t.Errorf("valueOf() = %v, want %v", got, tt.want)
			}
		})
	}
}
