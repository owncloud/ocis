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
