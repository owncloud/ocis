package backend

import "testing"

// TestGroupsClaimPresent guards the decision that protects users from losing all
// group memberships: only a present (and non-null) groups claim may trigger
// membership reconciliation. An absent or null claim must be treated as "do not
// sync", not as "no groups".
func TestGroupsClaimPresent(t *testing.T) {
	const claim = "groups"

	tests := []struct {
		name   string
		claims map[string]interface{}
		want   bool
	}{
		{"claim absent", map[string]interface{}{"other": "x"}, false},
		{"claim null", map[string]interface{}{claim: nil}, false},
		{"empty claims", map[string]interface{}{}, false},
		{"present but empty", map[string]interface{}{claim: []interface{}{}}, true},
		{"present with groups", map[string]interface{}{claim: []interface{}{"admins"}}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := groupsClaimPresent(tt.claims, claim); got != tt.want {
				t.Errorf("groupsClaimPresent(%v) = %v, want %v", tt.claims, got, tt.want)
			}
		})
	}
}
