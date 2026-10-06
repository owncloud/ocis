package svc

import (
	"reflect"
	"testing"
)

func TestTagsDiff(t *testing.T) {
	cases := []struct {
		name string
		a, b []string
		want []string
	}{
		{"added tags only", []string{"work", "play"}, []string{"work"}, []string{"play"}},
		{"nothing new", []string{"work"}, []string{"work"}, []string{}},
		{"removed tags", []string{"work", "play", "home"}, []string{"home"}, []string{"work", "play"}},
		{"none present", []string{"home"}, []string{"home"}, []string{}},
		{"empty before", []string{"a"}, nil, []string{"a"}},
	}
	for _, c := range cases {
		if got := tagsDiff(c.a, c.b); !reflect.DeepEqual(got, c.want) {
			t.Errorf("%s: tagsDiff(%v, %v) = %v, want %v", c.name, c.a, c.b, got, c.want)
		}
	}
}
