// Package mfa provides functionality for multi-factor authentication (MFA).
// mfa_test.go contains usage examples and tests.
package mfa

import (
	"net/http"
)

// MFAHeader is the header to be used across http services
// to forward the access token.
const MFAHeader = "X-Multi-Factor-Authentication"

// MFARequiredHeader is the header returned by the server if step-up authentication is required.
const MFARequiredHeader = "X-Ocis-Mfa-Required"

// SetRequiredStatus sets the MFA required header and the statuscode to 403
func SetRequiredStatus(w http.ResponseWriter) {
	w.Header().Set(MFARequiredHeader, "true")
	w.WriteHeader(http.StatusForbidden)
}

// SetHeader sets the MFA header.
func SetHeader(r *http.Request, mfa bool) {
	if mfa {
		r.Header.Set(MFAHeader, "true")
		return
	}

	r.Header.Set(MFAHeader, "false")
}

// IsMFAHeaderTrue checks if the MFA header is set to "true".
func IsMFAHeaderTrue(r *http.Request) bool {
	return r.Header.Get(MFAHeader) == "true"
}
