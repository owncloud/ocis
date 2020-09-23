package middleware

import (
	"net/http/httptest"
	"testing"
	"time"
)

func TestIsSignedRequest(t *testing.T) {
	tests := []struct {
		url      string
		expected bool
	}{
		{"https://example.com/example.jpg", false},
		{"https://example.com/example.jpg?OC-Signature=something", true},
	}

	for _, tt := range tests {
		r := httptest.NewRequest("", tt.url, nil)
		result := isSignedRequest(r)
		if result != tt.expected {
			t.Errorf("with %s expected %t got %t", tt.url, tt.expected, result)
		}
	}
}

func TestAllRequiredParametersPresent(t *testing.T) {
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
		result := allRequiredParametersArePresent(r)
		if result != tt.expected {
			t.Errorf("with %s expected %t got %t", tt.params, tt.expected, result)
		}
	}
}

func TestRequestMethodMatches(t *testing.T) {
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
		result := requestMethodMatches(r)
		if result != tt.expected {
			t.Errorf("with method %s and url %s expected %t got %t", tt.method, tt.url, tt.expected, result)
		}
	}
}

func TestRequestMethodIsAllowed(t *testing.T) {
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
		result := requestMethodIsAllowed(tt.method, tt.allowed)
		if result != tt.expected {
			t.Errorf("with method %s and allowed methods %s expected %t got %t", tt.method, tt.allowed, tt.expected, result)
		}
	}
}

func TestUrlIsExpired(t *testing.T) {
	nowFunc := func() time.Time {
		t, _ := time.Parse(time.RFC3339, "2020-08-19T15:12:43.478Z")
		return t
	}

	tests := []struct {
		url      string
		expected bool
	}{
		{"http://example.com/example.jpg?OC-Date=2020-08-19T15:02:43.478Z&OC-Expires=1200", false},
		{"http://example.com/example.jpg?OC-Date=invalid&OC-Expires=1200", true},
		{"http://example.com/example.jpg?OC-Date=2020-08-19T15:02:43.478Z&OC-Expires=invalid", true},
	}

	for _, tt := range tests {
		r := httptest.NewRequest("", tt.url, nil)
		result := urlIsExpired(r, nowFunc)
		if result != tt.expected {
			t.Errorf("with %s expected %t got %t", tt.url, tt.expected, result)
		}
	}
}

func TestCreateSignature(t *testing.T) {
	expected := "27d2ebea381384af3179235114801dcd00f91e46f99fca72575301cf3948101d"
	s := createSignature("something", []byte("somerandomkey"))

	if s != expected {
		t.Fail()
	}
}
