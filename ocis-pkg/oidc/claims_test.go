package oidc_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/owncloud/ocis/v2/ocis-pkg/oidc"
)

type splitWithEscapingTest struct {
	// Name of the subtest.
	name string

	// string to split
	s string

	// seperator to use
	seperator string

	// escape character to use for escaping
	escape string

	expectedParts []string
}

func (swet splitWithEscapingTest) run(t *testing.T) {
	parts := oidc.SplitWithEscaping(swet.s, swet.seperator, swet.escape)
	if len(swet.expectedParts) != len(parts) {
		t.Errorf("mismatching length")
	}
	for i, v := range swet.expectedParts {
		if parts[i] != v {
			t.Errorf("expected part %d to be '%s', got '%s'", i, v, parts[i])
		}
	}
}

func TestSplitWithEscaping(t *testing.T) {
	tests := []splitWithEscapingTest{
		{
			name:          "plain claim name",
			s:             "roles",
			seperator:     ".",
			escape:        "\\",
			expectedParts: []string{"roles"},
		},
		{
			name:          "claim with .",
			s:             "my.roles",
			seperator:     ".",
			escape:        "\\",
			expectedParts: []string{"my", "roles"},
		},
		{
			name:          "claim with escaped .",
			s:             "my\\.roles",
			seperator:     ".",
			escape:        "\\",
			expectedParts: []string{"my.roles"},
		},
		{
			name:          "claim with escaped . left",
			s:             "my\\.other.roles",
			seperator:     ".",
			escape:        "\\",
			expectedParts: []string{"my.other", "roles"},
		},
		{
			name:          "claim with escaped . right",
			s:             "my.other\\.roles",
			seperator:     ".",
			escape:        "\\",
			expectedParts: []string{"my", "other.roles"},
		},
	}
	for _, test := range tests {
		t.Run(test.name, test.run)
	}
}

type walkSegmentsTest struct {
	// Name of the subtest.
	name string

	// path segments to walk
	segments []string

	// seperator to use
	claims map[string]any

	expected any

	wantErr bool
}

func (wst walkSegmentsTest) run(t *testing.T) {
	v, err := oidc.WalkSegments(wst.segments, wst.claims)
	if err != nil && !wst.wantErr {
		t.Errorf("%v", err)
	}
	if err == nil && wst.wantErr {
		t.Errorf("expected error")
	}
	if !reflect.DeepEqual(v, wst.expected) {
		t.Errorf("expected %v got %v", wst.expected, v)
	}
}

func TestWalkSegments(t *testing.T) {
	byt := []byte(`{"first":{"second":{"third":["value1","value2"]},"foo":"bar"},"fizz":"buzz"}`)
	var dat map[string]any
	if err := json.Unmarshal(byt, &dat); err != nil {
		t.Errorf("%v", err)
	}

	tests := []walkSegmentsTest{
		{
			name:     "one segment, single value",
			segments: []string{"first"},
			claims: map[string]any{
				"first": "value",
			},
			expected: "value",
			wantErr:  false,
		},
		{
			name:     "one segment, array value",
			segments: []string{"first"},
			claims: map[string]any{
				"first": []string{"value1", "value2"},
			},
			expected: []string{"value1", "value2"},
			wantErr:  false,
		},
		{
			name:     "two segments, single value",
			segments: []string{"first", "second"},
			claims: map[string]any{
				"first": map[string]any{
					"second": "value",
				},
			},
			expected: "value",
			wantErr:  false,
		},
		{
			name:     "two segments, array value",
			segments: []string{"first", "second"},
			claims: map[string]any{
				"first": map[string]any{
					"second": []string{"value1", "value2"},
				},
			},
			expected: []string{"value1", "value2"},
			wantErr:  false,
		},
		{
			name:     "three segments, array value from json",
			segments: []string{"first", "second", "third"},
			claims:   dat,
			expected: []any{"value1", "value2"},
			wantErr:  false,
		},
		{
			name:     "three segments, array value with interface key",
			segments: []string{"first", "second", "third"},
			claims: map[string]any{
				"first": map[any]any{
					"second": map[any]any{
						"third": []string{"value1", "value2"},
					},
				},
			},
			expected: []string{"value1", "value2"},
			wantErr:  false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, test.run)
	}
}

func TestReadNumericClaim(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		claims  map[string]any
		want    int64
		wantErr bool
	}{
		{name: "float64", path: "auth_time", claims: map[string]any{"auth_time": float64(1700000000)}, want: 1700000000},
		{name: "int", path: "auth_time", claims: map[string]any{"auth_time": 1700000000}, want: 1700000000},
		{name: "int64", path: "auth_time", claims: map[string]any{"auth_time": int64(1700000000)}, want: 1700000000},
		{name: "uint32", path: "auth_time", claims: map[string]any{"auth_time": uint32(1700000000)}, want: 1700000000},
		{name: "json.Number", path: "auth_time", claims: map[string]any{"auth_time": json.Number("1700000000")}, want: 1700000000},
		{name: "nested path", path: "loa.times", claims: map[string]any{"loa": map[string]any{"times": float64(1700000000)}}, want: 1700000000},
		{name: "escaped dot", path: "loa\\.times", claims: map[string]any{"loa.times": float64(1700000000)}, want: 1700000000},
		{name: "missing", path: "auth_time", claims: map[string]any{}, wantErr: true},
		{name: "string", path: "auth_time", claims: map[string]any{"auth_time": "1700000000"}, wantErr: true},
		{name: "invalid json.Number", path: "auth_time", claims: map[string]any{"auth_time": json.Number("abc")}, wantErr: true},
		{name: "broken path", path: "loa.times", claims: map[string]any{"loa": "x"}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := oidc.ReadNumericClaim(tt.path, tt.claims)
			if (err != nil) != tt.wantErr {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("expected %d, got %d", tt.want, got)
			}
		})
	}
}
