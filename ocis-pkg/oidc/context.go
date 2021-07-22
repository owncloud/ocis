package oidc

import "context"

// contextKey is the key for oidc claims in a context
type contextKey struct{}

// NewContext makes a new context that contains the OpenID connect claims in a map.
func NewContext(parent context.Context, c map[string]interface{}) context.Context {
	return context.WithValue(parent, contextKey{}, c)
}

// FromContext returns the claims map stored in a context, or nil if there isn't one.
func FromContext(ctx context.Context) map[string]interface{} {
	s, _ := ctx.Value(contextKey{}).(map[string]interface{})
	return s
}
