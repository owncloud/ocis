package ctx

import "context"

// MFAHeader is the gRPC metadata key used to propagate MFA status across
// service boundaries. Lowercased to satisfy gRPC metadata requirements.
// The corresponding HTTP header set by the proxy is "X-Multi-Factor-Authentication".
const MFAHeader = "x-mfa-authenticated"

// ContextGetMFA returns the MFA status stored in the context, and whether it
// was set at all. A missing value (second return = false) should be treated as
// MFA not satisfied.
func ContextGetMFA(ctx context.Context) (bool, bool) {
	v, ok := ctx.Value(mfaKey).(bool)
	return v, ok
}

// ContextSetMFA stores the MFA status in the context.
func ContextSetMFA(ctx context.Context, mfa bool) context.Context {
	return context.WithValue(ctx, mfaKey, mfa)
}
