package metadata

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// NewUnary returns a server interceptor that will propagate GRPC metadata
func NewUnary(autoPropPrefix string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			md = metadata.MD{}
		}

		pairs := make([]string, 0, md.Len()*2)
		for key, values := range md {
			if strings.HasPrefix(key, autoPropPrefix) {
				for _, value := range values {
					pairs = append(pairs, key, value)
				}
			}
		}

		newctx := metadata.AppendToOutgoingContext(ctx, pairs...)
		return handler(newctx, req)
	}
}

// NewStream returns a server interceptor that will propagate GRPC metadata
func NewStream(autoPropPrefix string) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx := ss.Context()

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			md = metadata.MD{}
		}

		pairs := make([]string, 0, md.Len()*2)
		for key, values := range md {
			if strings.HasPrefix(key, autoPropPrefix) {
				for _, value := range values {
					pairs = append(pairs, key, value)
				}
			}
		}

		newctx := metadata.AppendToOutgoingContext(ctx, pairs...)

		wrapped := newWrappedServerStream(newctx, ss)
		return handler(srv, wrapped)
	}
}

// newWrappedServerStream returns a wrapped server stream in order to customize
// the GRPC server stream's context. It will use the one provided.
func newWrappedServerStream(ctx context.Context, ss grpc.ServerStream) *wrappedServerStream {
	return &wrappedServerStream{ServerStream: ss, newCtx: ctx}
}

type wrappedServerStream struct {
	grpc.ServerStream
	newCtx context.Context
}

// Context returns the context of the server stream
func (ss *wrappedServerStream) Context() context.Context {
	return ss.newCtx
}
