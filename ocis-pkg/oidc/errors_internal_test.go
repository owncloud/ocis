package oidc

import (
	"context"
	"errors"
	"net"
	"net/http"
	"testing"
)

type timeoutErr struct{}

func (timeoutErr) Error() string   { return "i/o timeout" }
func (timeoutErr) Timeout() bool   { return true }
func (timeoutErr) Temporary() bool { return true }

var _ net.Error = timeoutErr{}

func TestClassifyTransport(t *testing.T) {
	tests := []struct {
		name      string
		err       error
		transient bool
	}{
		{"nil", nil, false},
		{"deadline exceeded", context.DeadlineExceeded, true},
		{"canceled", context.Canceled, true},
		{"wrapped deadline", errors.Join(errors.New("get userinfo"), context.DeadlineExceeded), true},
		{"net error", timeoutErr{}, true},
		{"plain error", errors.New("token revoked"), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := classifyTransport(tt.err)
			if tt.err == nil {
				if got != nil {
					t.Fatalf("expected nil, got %v", got)
				}
				return
			}
			if errors.Is(got, ErrTemporarilyUnavailable) != tt.transient {
				t.Fatalf("transient = %v, want %v (err=%v)", !tt.transient, tt.transient, got)
			}
			if !errors.Is(got, tt.err) {
				t.Fatalf("classifyTransport dropped the underlying error: %v", got)
			}
		})
	}
}

func TestStatusIsTransient(t *testing.T) {
	tests := map[int]bool{
		http.StatusUnauthorized:            false,
		http.StatusForbidden:               false,
		http.StatusNotFound:                false,
		http.StatusTooManyRequests:         true,
		http.StatusInternalServerError:     true,
		http.StatusBadGateway:              true,
		http.StatusServiceUnavailable:      true,
		http.StatusGatewayTimeout:          true,
		http.StatusNotImplemented:          false, // permanent 5xx must fail fast
		http.StatusHTTPVersionNotSupported: false,
	}
	for status, want := range tests {
		if got := statusIsTransient(status); got != want {
			t.Errorf("statusIsTransient(%d) = %v, want %v", status, got, want)
		}
	}
}
