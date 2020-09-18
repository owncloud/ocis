package oidc

import "context"

// contextKey is the key for oidc claims in a context
type contextKey struct{}

// NewContext makes a new context that contains the OpenID Connect claims.
func NewContext(parent context.Context, c *StandardClaims) context.Context {
	return context.WithValue(parent, contextKey{}, c)
}

// FromContext returns the StandardClaims stored in a context, or nil if there isn't one.
func FromContext(ctx context.Context) *StandardClaims {
	s, _ := ctx.Value(contextKey{}).(*StandardClaims)
	return s
}
