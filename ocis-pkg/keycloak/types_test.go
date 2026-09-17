package keycloak_test

import (
	"testing"

	"github.com/owncloud/ocis/v2/ocis-pkg/keycloak"
	"github.com/stretchr/testify/assert"
)

func TestUserActionStringValues(t *testing.T) {
	// UserAction is a free-form string passed to Keycloak verbatim; the
	// convenience constants must carry the exact Keycloak required-action
	// provider IDs.
	assert.Equal(t, "UPDATE_PASSWORD", string(keycloak.UserActionUpdatePassword))
	assert.Equal(t, "VERIFY_EMAIL", string(keycloak.UserActionVerifyEmail))
	assert.Equal(t, "delete_account", string(keycloak.UserActionDeleteAccount))
	assert.Equal(t, "webauthn-register", string(keycloak.UserActionWebauthnRegister))

	// Arbitrary, non-builtin actions are valid - there is no allowlist.
	assert.Equal(t, "CUSTOM_ACTION", string(keycloak.UserAction("CUSTOM_ACTION")))
}
