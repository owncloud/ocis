package checkers

import (
	"fmt"
	"regexp"

	"github.com/owncloud/ocis/v2/ocis-pkg/oidc"
)

// RegexpChecker checks whether the specific key matches the provided regular
// expresion.
type RegexpChecker struct {
	key string
	exp *regexp.Regexp
}

// NewRegexpChecker creates a new RegexpChecker
func NewRegexpChecker(key, pattern string) *RegexpChecker {
	return &RegexpChecker{
		key: key,
		exp: regexp.MustCompile(pattern),
	}
}

// CheckClaims checks in the claims if the claims key's value matches the
// provided regular expresion. If it doesn't match, an error is returned.
func (rc *RegexpChecker) CheckClaims(claims map[string]interface{}) error {
	value, err := oidc.ReadStringClaim(rc.key, claims)
	if err != nil {
		return err
	}

	if !rc.exp.MatchString(value) {
		return fmt.Errorf("wrong value for claim '%s' - value '%s' doesn't match regular expresion '%s'", rc.key, value, rc.exp.String())
	}
	return nil
}
