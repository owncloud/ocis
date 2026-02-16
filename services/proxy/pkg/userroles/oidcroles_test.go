package userroles

import (
	"encoding/json"
	"testing"
)

func TestExtractRolesArray(t *testing.T) {
	byt := []byte(`{"roles":["a","b"]}`)

	claims := map[string]interface{}{}
	err := json.Unmarshal(byt, &claims)
	if err != nil {
		t.Fatal(err)
	}

	roles, err := extractRoles("roles", claims)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := roles["a"]; !ok {
		t.Fatal("must contain 'a'")
	}
	if _, ok := roles["b"]; !ok {
		t.Fatal("must contain 'b'")
	}
}

func TestExtractRolesString(t *testing.T) {
	byt := []byte(`{"roles":"a"}`)

	claims := map[string]interface{}{}
	err := json.Unmarshal(byt, &claims)
	if err != nil {
		t.Fatal(err)
	}

	roles, err := extractRoles("roles", claims)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := roles["a"]; !ok {
		t.Fatal("must contain 'a'")
	}
}

func TestExtractRolesPathArray(t *testing.T) {
	byt := []byte(`{"sub":{"roles":["a","b"]}}`)

	claims := map[string]interface{}{}
	err := json.Unmarshal(byt, &claims)
	if err != nil {
		t.Fatal(err)
	}

	roles, err := extractRoles("sub.roles", claims)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := roles["a"]; !ok {
		t.Fatal("must contain 'a'")
	}
	if _, ok := roles["b"]; !ok {
		t.Fatal("must contain 'b'")
	}
}

func TestExtractRolesPathString(t *testing.T) {
	byt := []byte(`{"sub":{"roles":"a"}}`)

	claims := map[string]interface{}{}
	err := json.Unmarshal(byt, &claims)
	if err != nil {
		t.Fatal(err)
	}

	roles, err := extractRoles("sub.roles", claims)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := roles["a"]; !ok {
		t.Fatal("must contain 'a'")
	}
}

func TestExtractEscapedRolesPathString(t *testing.T) {
	byt := []byte(`{"sub.roles":"a"}`)

	claims := map[string]interface{}{}
	err := json.Unmarshal(byt, &claims)
	if err != nil {
		t.Fatal(err)
	}

	roles, err := extractRoles("sub\\.roles", claims)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := roles["a"]; !ok {
		t.Fatal("must contain 'a'")
	}
}

func TestNoRoles(t *testing.T) {
	byt := []byte(`{"sub":{"foo":"a"}}`)

	claims := map[string]interface{}{}
	err := json.Unmarshal(byt, &claims)
	if err != nil {
		t.Fatal(err)
	}

	roles, err := extractRoles("sub.roles", claims)
	if err == nil {
		t.Fatal("must not find a role")
	}
	if len(roles) != 0 {
		t.Fatal("length of roles mut be 0")
	}
}

func TestMatchesClaimMappingExact(t *testing.T) {
	claimRoles := map[string]struct{}{
		"ocis-user": {},
	}
	if !matchesClaimMapping("ocis-user", claimRoles) {
		t.Fatal("expected exact match to succeed")
	}
	if matchesClaimMapping("admin", claimRoles) {
		t.Fatal("expected non-matching literal to fail")
	}
}

func TestMatchesClaimMappingRegex(t *testing.T) {
	claimRoles := map[string]struct{}{
		"ocis-user-1":   {},
		"ocis-user-42":  {},
		"ocis-user-lth": {},
		"admin":         {},
	}
	if !matchesClaimMapping("ocis-user-.*", claimRoles) {
		t.Fatal("expected regex match to succeed")
	}
	if !matchesClaimMapping("ocis-user-[a-zA-Z0-9]", claimRoles) {
		t.Fatal("expected regex match to succeed")
	}
	if matchesClaimMapping("admin-.*", claimRoles) {
		t.Fatal("expected regex match to fail for admin-.*")
	}
}

func TestMatchesClaimMappingInvalidRegexFallsBackToExact(t *testing.T) {
	claimRoles := map[string]struct{}{"ocis-user": {}}
	// invalid regex pattern
	if matchesClaimMapping("ocis-user[", claimRoles) {
		t.Fatal("invalid regex should fall back to exact and not match")
	}
}
