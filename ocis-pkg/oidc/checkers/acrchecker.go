package checkers

import (
	"fmt"

	"github.com/owncloud/ocis/v2/ocis-pkg/oidc"
)

// AcrChecker check if the acr in the claims has the exact same value
// as the provided one
type AcrChecker struct {
	value string
}

// NewAcrChecker creates a new AcrChecker
func NewAcrChecker(value string) *AcrChecker {
	return &AcrChecker{
		value: value,
	}
}

// CheckClaims checks if the provided value matches the acr value in the
// claims. It's an exact match.
func (ac *AcrChecker) CheckClaims(claims map[string]interface{}) error {
	value, err := oidc.ReadStringClaim("acr", claims)
	if err != nil {
		return err
	}

	if ac.value != value {
		return fmt.Errorf("wrong value for 'acr' - expected '%s' actual '%s'", ac.value, value)
	}
	return nil
}

func (ac *AcrChecker) RequireMap() map[string]string {
	return map[string]string{
		"Type": "Acr",
		"Data": "acr=" + ac.value,
	}
}
