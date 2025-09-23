// Package mfa provides functionality for multi-factor authentication (MFA).
// mfa_test.go contains usage examples and tests.
package mfa

import (
	"context"
	"net/http"
)

// MFAHeader is the header to be used across grpc and http services
// to forward the access token.
const MFAHeader = "X-Multi-Factor-Authentication"

// MFARequiredHeader is the header returned by the server if step-up authentication is required.
const MFARequiredHeader = "X-Ocis-Mfa-Required"

type mfaKeyType struct{}

var mfaKey = mfaKeyType{}

// EnhanceRequest enhances the request context with the MFA status from the header.
// This operation does not overwrite existing context values.
func EnhanceRequest(req *http.Request) *http.Request {
	ctx := req.Context()
	if Has(ctx) {
		return req
	}
	return req.WithContext(Set(ctx, req.Header.Get(MFAHeader) == "true"))
}

// SetRequiredStatus sets the MFA required header and the statuscode to 403
func SetRequiredStatus(w http.ResponseWriter) {
	w.Header().Set(MFARequiredHeader, "true")
	w.WriteHeader(http.StatusForbidden)
}

// Has returns the mfa status from the context.
func Has(ctx context.Context) bool {
	mfa, ok := ctx.Value(mfaKey).(bool)
	if !ok {
		return false
	}
	return mfa
}

// Set stores the mfa status in the context.
func Set(ctx context.Context, mfa bool) context.Context {
	return context.WithValue(ctx, mfaKey, mfa)
}

// SetHeader sets the MFA header.
func SetHeader(r *http.Request, mfa bool) {
	if mfa {
		r.Header.Set(MFAHeader, "true")
		return
	}

	r.Header.Set(MFAHeader, "false")
}
