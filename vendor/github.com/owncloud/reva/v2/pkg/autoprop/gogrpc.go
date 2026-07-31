package autoprop

import (
	"context"

	"google.golang.org/grpc"
)

// GetGoGRPCUnaryClientInterceptor will create a new UnaryClientInterceptor
// that will include the oCIS metadata present in the context. The keys from
// the metadata will have the AutoPropPrefix prepended so they can be easily
// identified and propagated.
func GetGoGRPCUnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		outCtx := moveOcisMetaToOutgoingContext(ctx)
		return invoker(outCtx, method, req, reply, cc, opts...)
	}
}

// GetGoGRPCStreamClientInterceptor will create a new StreamClientInterceptor
// that will include the oCIS metadata present in the context. The keys from
// the metadata will have the AutoPropPrefix prepended so they can be easily
// identified and propagated.
func GetGoGRPCStreamClientInterceptor() grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		outCtx := moveOcisMetaToOutgoingContext(ctx)
		return streamer(outCtx, desc, cc, method, opts...)
	}
}

// GetGoGRPCUnaryServerInterceptor will create a new UnaryServerInterceptor
// that will copy the incoming context values to the oCIS metadata. Only the
// keys prepended with the AutoPropPrefix will be copied.
// This method also set the oCIS metadata in the context.
// NOTE: keys without the AutoPropPrefix will still be available (in the
// incoming context, not in the oCIS metadata).
func GetGoGRPCUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		newCtx := moveIncomingContextToOcisMeta(ctx)
		return handler(newCtx, req)
	}
}

// GetGoGRPCStreamServerInterceptor will create a new StreamServerInterceptor
// that will copy the incoming context values to the oCIS metadata. Only the
// keys prepended with the AutoPropPrefix will be copied.
// This method also set the oCIS metadata in the context.
// NOTE: keys without the AutoPropPrefix will still be available (in the
// incoming context, not in the oCIS metadata).
func GetGoGRPCStreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx := ss.Context()
		newCtx := moveIncomingContextToOcisMeta(ctx)

		wrapped := newWrappedServerStream(newCtx, ss)
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
