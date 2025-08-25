package checkers

// Checker will allow different checks based on the OIDC claims.
// Each implementation will perform different checks.
type Checker interface {
	// CheckClaims will check whether the claims match specific criteria.
	// If the check passes, nil will be returned, otherwise a proper error
	// will be returned instead.
	CheckClaims(claims map[string]interface{}) error
	// RequireMap returns a map with the expected headers in the failed
	// response. The headers should have enough information for the client
	// to know what's going on.
	// Expected keys to be returned are "Type" and "Data".
	// All the returned keys will be prepended with "X-OCIS-<auth>-Requires-*"
	// For example:
	// {"X-OCIS-OIDC-Requires-Type": "Bool", "X-OCIS-OIDC-Requires-Data": "email_verified=true"}
	// or
	// {"X-OCIS-OIDC-Requires-Type": "Acr", "X-OCIS-OIDC-Requires-Data": "acr=advanced"}
	//
	// It's up to the client to decide what to do with those headers.
	// Usually, if any "X-OCIS-<auth>-Requires-*" header is received, the client
	// might just show a popup with a message such as "not enough permissions
	// to access to this resource", because the client might not be able
	// to do anything to fix the problem.
	// For the step up auth scenario, the client should be able to detect it
	// with the headers ("X-OCIS-OIDC-Requires-Type": "Acr") and
	// ("X-OCIS-OIDC-Requires-Data": "acr=advanced" - which means the acr should
	// be advanced), and then act accordingly.
	RequireMap() map[string]string
}
