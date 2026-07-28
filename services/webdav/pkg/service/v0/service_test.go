package svc

import (
	"net/http"
	"testing"
)

func TestNewErrResponseException(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		expected   string
	}{
		{"bad request", http.StatusBadRequest, "Sabre\\DAV\\Exception\\BadRequest"},
		{"unauthorized", http.StatusUnauthorized, "Sabre\\DAV\\Exception\\NotAuthenticated"},
		{"not found", http.StatusNotFound, "Sabre\\DAV\\Exception\\NotFound"},
		{"method not allowed", http.StatusMethodNotAllowed, "Sabre\\DAV\\Exception\\MethodNotAllowed"},
		{"forbidden", http.StatusForbidden, "Sabre\\DAV\\Exception\\Forbidden"},
		{"too early", http.StatusTooEarly, "Sabre\\DAV\\Exception\\TooEarly"},
		{"too many requests", http.StatusTooManyRequests, "Sabre\\DAV\\Exception\\TooManyRequests"},
		// unmapped codes must still yield a non-empty exception
		{"internal server error", http.StatusInternalServerError, defaultException},
		{"not implemented", http.StatusNotImplemented, defaultException},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rsp := newErrResponse(tt.statusCode, "some message")
			if rsp.Exception != tt.expected {
				t.Errorf("expected exception %q, got %q", tt.expected, rsp.Exception)
			}
			if rsp.HTTPStatusCode != tt.statusCode {
				t.Errorf("expected status code %d, got %d", tt.statusCode, rsp.HTTPStatusCode)
			}
			if rsp.Message != "some message" {
				t.Errorf("expected message to be set, got %q", rsp.Message)
			}
		})
	}
}

func TestNewErrResponseEmptyMessage(t *testing.T) {
	rsp := newErrResponse(http.StatusNotFound, "")
	if rsp.Message != "" {
		t.Errorf("expected empty message, got %q", rsp.Message)
	}
	if rsp.Exception != "Sabre\\DAV\\Exception\\NotFound" {
		t.Errorf("unexpected exception %q", rsp.Exception)
	}
}
