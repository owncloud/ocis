package oidc_test

import (
	"encoding/json"
	"errors"
	"reflect"
	"strings"
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
	claims map[string]interface{}

	expected interface{}

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
	var dat map[string]interface{}
	if err := json.Unmarshal(byt, &dat); err != nil {
		t.Errorf("%v", err)
	}

	tests := []walkSegmentsTest{
		{
			name:     "one segment, single value",
			segments: []string{"first"},
			claims: map[string]interface{}{
				"first": "value",
			},
			expected: "value",
			wantErr:  false,
		},
		{
			name:     "one segment, array value",
			segments: []string{"first"},
			claims: map[string]interface{}{
				"first": []string{"value1", "value2"},
			},
			expected: []string{"value1", "value2"},
			wantErr:  false,
		},
		{
			name:     "two segments, single value",
			segments: []string{"first", "second"},
			claims: map[string]interface{}{
				"first": map[string]interface{}{
					"second": "value",
				},
			},
			expected: "value",
			wantErr:  false,
		},
		{
			name:     "two segments, array value",
			segments: []string{"first", "second"},
			claims: map[string]interface{}{
				"first": map[string]interface{}{
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
			expected: []interface{}{"value1", "value2"},
			wantErr:  false,
		},
		{
			name:     "three segments, array value with interface key",
			segments: []string{"first", "second", "third"},
			claims: map[string]interface{}{
				"first": map[interface{}]interface{}{
					"second": map[interface{}]interface{}{
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

// TestWalkSegmentsDistinguishesAnAbsentSegmentFromAnUnreadableOne pins the
// distinction ErrMissingClaim exists for. A path segment that is not in the claims
// at all means the token is simply silent about whatever lives under it; a segment
// that is present but cannot be descended into means the token says something there
// that we cannot make sense of. Both used to come back as the same "unsupported
// type" error, so a caller could not tell them apart without re-walking the claims
// itself.
func TestWalkSegmentsDistinguishesAnAbsentSegmentFromAnUnreadableOne(t *testing.T) {
	t.Run("absent intermediate segment wraps ErrMissingClaim", func(t *testing.T) {
		_, err := oidc.WalkSegments(
			[]string{"resource_access", "ocis", "roles"},
			map[string]interface{}{"sub": "abcd"},
		)
		if !errors.Is(err, oidc.ErrMissingClaim) {
			t.Fatalf("expected ErrMissingClaim, got %v", err)
		}
		if !strings.Contains(err.Error(), "resource_access") {
			t.Fatalf("the error should name the missing segment, got %v", err)
		}
	})

	unreadable := map[string]map[string]interface{}{
		"intermediate segment is a string": {"resource_access": "not-an-object"},
		"intermediate segment is a number": {"resource_access": 42},
		"intermediate map has a non-string key": {
			"resource_access": map[interface{}]interface{}{7: "ocis"},
		},
	}
	for name, claims := range unreadable {
		t.Run(name+" does not wrap ErrMissingClaim", func(t *testing.T) {
			_, err := oidc.WalkSegments([]string{"resource_access", "ocis", "roles"}, claims)
			if err == nil {
				t.Fatal("expected an error for a segment that cannot be descended into")
			}
			if errors.Is(err, oidc.ErrMissingClaim) {
				t.Fatalf("a present but unreadable segment must not report as missing, got %v", err)
			}
		})
	}

	// The boundary: only intermediate segments are walked. An absent *leaf* is not an
	// error at all - the walk succeeds and yields nil, and it is the caller's job to
	// decide what an empty claim means.
	t.Run("absent leaf segment is not an error", func(t *testing.T) {
		got, err := oidc.WalkSegments(
			[]string{"resource_access", "ocis", "roles"},
			map[string]interface{}{
				"resource_access": map[string]interface{}{
					"ocis": map[string]interface{}{},
				},
			},
		)
		if err != nil {
			t.Fatalf("expected no error for an absent leaf, got %v", err)
		}
		if got != nil {
			t.Fatalf("expected a nil claim, got %v", got)
		}
	})
}
