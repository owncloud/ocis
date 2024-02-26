package store

import (
	"context"

	"go-micro.dev/v4/client"
	"go-micro.dev/v4/store"
)

type grpcClientContextKey struct{}

// WithGRPCClient sets the grpc client
func WithGRPCClient(c client.Client) store.Option {
	return func(o *store.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, grpcClientContextKey{}, c)
	}
}
