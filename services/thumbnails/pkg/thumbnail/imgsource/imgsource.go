package imgsource

import (
	"context"
	"io"
)

type key int

const (
	auth key = iota
)

// Source defines the interface for image sources
type Source interface {
	Get(ctx context.Context, path string) (io.ReadCloser, error)
}

// ContextSetAuthorization puts the authorization in the context.
func ContextSetAuthorization(parent context.Context, authorization string) context.Context {
	return context.WithValue(parent, auth, authorization)
}

// ContextGetAuthorization gets the authorization from the context.
func ContextGetAuthorization(ctx context.Context) (string, bool) {
	val := ctx.Value(auth)
	if val == nil {
		return "", false
	}
	return val.(string), true
}
