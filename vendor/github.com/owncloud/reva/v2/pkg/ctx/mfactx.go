package ctx

import (
	"context"
	"slices"

	"github.com/owncloud/reva/v2/pkg/autoprop"
)

const (
	// key to use for auto-propagation. This key is exclusive to this package
	autopropKey = "mfa-authenticated"
)

// SetMFA sets the MFA status in the context and ensures it autopropagate
// across service boundaries.
// If the MFA is already set, this method will just return the passed context,
// otherwise a new context (from the provided one) containing the MFA status
// will be returned
func SetMFA(ctx context.Context) context.Context {
	// just return the same context if already set
	if HasMFA(ctx) {
		return ctx
	}
	return autoprop.AppendMetaToContext(ctx, autopropKey, "true")
}

// RemoveMFA removes the MFA status from the context.
// The meta associated to the provided context will be copied, and the MFA
// status will be removed from the copied context. The copied context
// will be returned.
//
// WARNING: Previous MFA status will still be available in the old context.
func RemoveMFA(ctx context.Context) context.Context {
	ctx2 := autoprop.CopyMetaToContext(ctx) // ensures ctx2 to have a meta
	meta := autoprop.GetMetaFromContext(ctx2)
	meta.DeleteMeta(autopropKey)
	return ctx2
}

// HasMFA checks if the context has the MFA status. Use SetMFA to set it.
func HasMFA(ctx context.Context) bool {
	meta := autoprop.GetMetaFromContext(ctx)
	if meta == nil {
		return false
	}

	if values, exists := meta.GetMetaWithExists(autopropKey); exists {
		if slices.Contains(values, "true") {
			return true
		}
	}
	return false
}
