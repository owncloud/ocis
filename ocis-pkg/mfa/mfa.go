// package mfa provides functionality for multi-factor authentication (MFA).
package mfa

import (
	"context"
	"net/http"
)

// MFAHeader is the header to be used across grpc and http services
// to forward the access token.
const MFAHeader = "x-multi-factor-authentication"

// MFARequiredHeader is the header returned by the server if step-up authentication is required.
const MFARequiredHeader = "X-OCIS-MFA-Required"

type mfaKeyType struct{}

var mfaKey = mfaKeyType{}

// FromRequest extracts the mfa status from the request headers
func FromRequest(ctx context.Context, req *http.Request) context.Context {
	return Set(ctx, req.Header.Get(MFAHeader) == "true")
}

// Accepted checks if the context has MFA authentication. If not, it sets the required header and status code.
func Accepted(ctx context.Context, w http.ResponseWriter) bool {
	hasMFA := Has(ctx)
	if !hasMFA {
		SetRequiredHeader(w)
	}
	return hasMFA
}

// Has returns true if the user is MFA authenticated.
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

// SetRequiredHeader sets the MFA required header and the statuscode to 403
func SetRequiredHeader(w http.ResponseWriter) {
	w.Header().Set(MFARequiredHeader, "true")
	w.WriteHeader(http.StatusForbidden)
}
