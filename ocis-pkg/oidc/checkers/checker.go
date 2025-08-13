package checkers

// Checker will allow different checks based on the OIDC claims.
// Each implementation will perform different checks.
type Checker interface {
	// CheckClaims will check whether the claims match specific criteria.
	// If the check passes, nil will be returned, otherwise a proper error
	// will be returned instead.
	CheckClaims(claims map[string]interface{}) error
}
