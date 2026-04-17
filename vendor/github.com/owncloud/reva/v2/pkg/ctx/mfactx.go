package ctx

// MFAOutgoingHeader is the gRPC metadata key used to propagate MFA status across
// service boundaries. The "autoprop-" prefix causes the metadata interceptor
// (internal/grpc/interceptors/metadata) to forward it automatically at every
// gRPC hop, so no manual re-forwarding is required.
// The const rgrpc.AutoPropPrefix causes the cycle import
const MFAOutgoingHeader = "autoprop-mfa-authenticated"

// The corresponding HTTP header set by the proxy is "X-Multi-Factor-Authentication".
const MFAHeader = "X-Multi-Factor-Authentication"
