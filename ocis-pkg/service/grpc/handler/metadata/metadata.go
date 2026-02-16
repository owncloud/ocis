package metadata

import (
	"context"
	"strings"

	"github.com/owncloud/reva/v2/pkg/rgrpc"
	"go-micro.dev/v4/server"
	"google.golang.org/grpc/metadata"
)

const (
	// Prefix used to auto propagate GRPC metadata keys across servers.
	// This is used in the NewHandlerWrapper and NewSubscriberWrapper, and
	// it's expected to work in go-micro services.
	// It needs to match the prefix used in reva for the same purpose.
	AutoPropPrefix = rgrpc.AutoPropPrefix
)

// NewHandlerWrapper propagates the grpc metadata.
func NewHandlerWrapper() server.HandlerWrapper {
	return func(h server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp interface{}) error {
			md, ok := metadata.FromIncomingContext(ctx)
			if !ok {
				md = metadata.MD{}
			}

			pairs := make([]string, 0, md.Len()*2)
			for key, values := range md {
				if strings.HasPrefix(key, AutoPropPrefix) {
					for _, value := range values {
						pairs = append(pairs, key, value)
					}
				}
			}

			newctx := metadata.AppendToOutgoingContext(ctx, pairs...)
			return h(newctx, req, rsp)
		}
	}
}

// NewSubscriberWrapper propagates the grpc metadata
func NewSubscriberWrapper() server.SubscriberWrapper {
	return func(next server.SubscriberFunc) server.SubscriberFunc {
		return func(ctx context.Context, msg server.Message) error {
			md, ok := metadata.FromIncomingContext(ctx)
			if !ok {
				md = metadata.MD{}
			}

			pairs := make([]string, 0, md.Len()*2)
			for key, values := range md {
				if strings.HasPrefix(key, AutoPropPrefix) {
					for _, value := range values {
						pairs = append(pairs, key, value)
					}
				}
			}

			newctx := metadata.AppendToOutgoingContext(ctx, pairs...)
			return next(newctx, msg)
		}
	}
}
