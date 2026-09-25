package oidc

import (
	"context"
	"errors"
	"net"
	"net/http"
)

// ErrTemporarilyUnavailable marks a transient IdP failure (timeout, network error
// or 5xx/429), which callers should treat as retryable (503) rather than as an
// authentication failure (401).
var ErrTemporarilyUnavailable = errors.New("oidc: identity provider temporarily unavailable")

func classifyTransport(err error) error {
	if err == nil {
		return nil
	}
	var netErr net.Error
	if errors.Is(err, context.DeadlineExceeded) ||
		errors.Is(err, context.Canceled) ||
		errors.As(err, &netErr) {
		return errors.Join(ErrTemporarilyUnavailable, err)
	}
	return err
}

// statusIsTransient lists the retryable statuses explicitly so a permanent 5xx
// (501, 505) fails fast instead of being retried forever.
func statusIsTransient(status int) bool {
	switch status {
	case http.StatusTooManyRequests,
		http.StatusInternalServerError,
		http.StatusBadGateway,
		http.StatusServiceUnavailable,
		http.StatusGatewayTimeout:
		return true
	default:
		return false
	}
}
