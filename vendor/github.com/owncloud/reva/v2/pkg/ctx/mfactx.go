package ctx

import (
	"context"

	"google.golang.org/grpc/metadata"
)

// MFAOutgoingHeader is the gRPC metadata key used to propagate MFA status across
// service boundaries. The "autoprop-" prefix causes the metadata interceptor
// (internal/grpc/interceptors/metadata) to forward it automatically at every
// gRPC hop, so no manual re-forwarding is required.
// Using rgrpc.AutoPropPrefix here would cause a cyclic import.
const MFAOutgoingHeader = "autoprop-mfa-authenticated"

// The corresponding HTTP header set by the proxy is "X-Multi-Factor-Authentication".
const MFAHeader = "X-Multi-Factor-Authentication"

// AppendMFAToOutgoingContext adds the MFA status to the outgoing gRPC metadata.
func AppendMFAToOutgoingContext(ctx context.Context, hasMFA bool) context.Context {
	mfaVal := "false"
	if hasMFA {
		mfaVal = "true"
	}
	return metadata.AppendToOutgoingContext(ctx, MFAOutgoingHeader, mfaVal)
}
