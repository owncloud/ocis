package checkers

import (
	"fmt"

	"github.com/owncloud/ocis/v2/ocis-pkg/oidc"
)

// BooleanChecker checks whether the specified key has a matching boolean value
type BooleanChecker struct {
	key   string
	value bool
}

// NewBooleanChecker creates a new BooleanChecker
func NewBooleanChecker(key string, value bool) *BooleanChecker {
	return &BooleanChecker{
		key:   key,
		value: value,
	}
}

// CheckClaims checks the claims so the key held by the BooleanChecker matches
// its boolean value.
func (bcc *BooleanChecker) CheckClaims(claims map[string]interface{}) error {
	value, err := oidc.ReadBoolClaim(bcc.key, claims)
	if err != nil {
		return err
	}

	if value != bcc.value {
		return fmt.Errorf("wrong value for claim '%s' - expected '%t' actual '%t'", bcc.key, bcc.value, value)
	}
	return nil
}
