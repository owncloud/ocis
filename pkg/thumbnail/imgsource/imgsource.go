package imgsource

import (
	"context"
	"image"
)

type key int

const (
	auth key = iota
)

// Source defines the interface for image sources
type Source interface {
	Get(ctx context.Context, path string) (image.Image, error)
}

// WithAuthorization puts the authorization in the context.
func WithAuthorization(parent context.Context, authorization string) context.Context {
	return context.WithValue(parent, auth, authorization)
}

func authorization(ctx context.Context) string {
	val := ctx.Value(auth)
	if val == nil {
		return ""
	}
	return val.(string)
}
