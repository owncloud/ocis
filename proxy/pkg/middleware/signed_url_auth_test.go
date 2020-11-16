package middleware

import (
	"net/http/httptest"
	"testing"
	"time"
)

func TestSignedURLAuth_shouldServe(t *testing.T) {
	pua := signedURLAuth{}
	tests := []struct {
		url      string
		expected bool
	}{
		{"https://example.com/example.jpg", false},
		{"https://example.com/example.jpg?OC-Signature=something", true},
	}

	for _, tt := range tests {
		r := httptest.NewRequest("", tt.url, nil)
		result := pua.shouldServe(r)

		if result != tt.expected {
			t.Errorf("with %s expected %t got %t", tt.url, tt.expected, result)
		}
	}
}

func TestSignedURLAuth_allRequiredParametersPresent(t *testing.T) {
	pua := signedURLAuth{}
	baseURL := "https://example.com/example.jpg?"
	tests := []struct {
		params   string
		expected bool
	}{
		{"OC-Signature=something&OC-Credential=something&OC-Date=something&OC-Expires=something&OC-Verb=something", true},
		{"OC-Credential=something&OC-Date=something&OC-Expires=something&OC-Verb=something", false},
		{"OC-Signature=something&OC-Date=something&OC-Expires=something&OC-Verb=something", false},
		{"OC-Signature=something&OC-Credential=something&OC-Expires=something&OC-Verb=something", false},
		{"OC-Signature=something&OC-Credential=something&OC-Date=something&OC-Verb=something", false},
		{"OC-Signature=something&OC-Credential=something&OC-Date=something&OC-Expires=something", false},
	}
	for _, tt := range tests {
		r := httptest.NewRequest("", baseURL+tt.params, nil)
		ok, _ := pua.allRequiredParametersArePresent(r.URL.Query())
		if ok != tt.expected {
			t.Errorf("with %s expected %t got %t", tt.params, tt.expected, ok)
		}
	}
}

func TestSignedURLAuth_requestMethodMatches(t *testing.T) {
	pua := signedURLAuth{}
	tests := []struct {
		method   string
		url      string
		expected bool
	}{
		{"GET", "https://example.com/example.jpg?OC-Verb=GET", true},
		{"GET", "https://example.com/example.jpg?OC-Verb=get", true},
		{"POST", "https://example.com/example.jpg?OC-Verb=GET", false},
	}

	for _, tt := range tests {
		r := httptest.NewRequest(tt.method, tt.url, nil)
		ok, _ := pua.requestMethodMatches(r.Method, r.URL.Query())
		if ok != tt.expected {
			t.Errorf("with method %s and url %s expected %t got %t", tt.method, tt.url, tt.expected, ok)
		}
	}
}

func TestSignedURLAuth_requestMethodIsAllowed(t *testing.T) {
	pua := signedURLAuth{}
	tests := []struct {
		method   string
		allowed  []string
		expected bool
	}{
		{"GET", []string{}, false},
		{"GET", []string{"POST"}, false},
		{"GET", []string{"GET"}, true},
		{"GET", []string{"get"}, true},
		{"GET", []string{"POST", "GET"}, true},
	}

	for _, tt := range tests {
		pua.preSignedURLConfig.AllowedHTTPMethods = tt.allowed
		ok, _ := pua.requestMethodIsAllowed(tt.method)

		if ok != tt.expected {
			t.Errorf("with method %s and allowed methods %s expected %t got %t", tt.method, tt.allowed, tt.expected, ok)
		}
	}
}

func TestSignedURLAuth_urlIsExpired(t *testing.T) {
	pua := signedURLAuth{}
	nowFunc := func() time.Time {
		t, _ := time.Parse(time.RFC3339, "2020-02-02T12:30:00.000Z")
		return t
	}

	tests := []struct {
		url      string
		isExpired bool
	}{
		{"http://example.com/example.jpg?OC-Date=2020-02-02T12:29:00.000Z&OC-Expires=61", false},
		{"http://example.com/example.jpg?OC-Date=2020-02-02T12:29:00.000Z&OC-Expires=invalid", true},
		{"http://example.com/example.jpg?OC-Date=2020-02-02T12:29:00.000Z&OC-Expires=59", true},
		{"http://example.com/example.jpg?OC-Date=2020-02-03T12:29:00.000Z&OC-Expires=59", true},
		{"http://example.com/example.jpg?OC-Date=2020-02-01T12:29:00.000Z&OC-Expires=59", true},
	}

	for _, tt := range tests {
		r := httptest.NewRequest("", tt.url, nil)
		expired, _ := pua.urlIsExpired(r.URL.Query(), nowFunc)
		if expired != tt.isExpired {
			t.Errorf("with %s expected %t got %t", tt.url, tt.isExpired, expired)
		}
	}
}

func TestSignedURLAuth_createSignature(t *testing.T) {
	pua := signedURLAuth{}
	expected := "27d2ebea381384af3179235114801dcd00f91e46f99fca72575301cf3948101d"
	s := pua.createSignature("something", []byte("somerandomkey"))

	if s != expected {
		t.Fail()
	}
}
