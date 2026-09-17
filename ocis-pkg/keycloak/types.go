package keycloak

import (
	"context"

	"github.com/Nerzal/gocloak/v13"
	libregraph "github.com/owncloud/libre-graph-api-go"
)

// UserAction is a Keycloak required-action provider ID. It is a free-form
// string so operators can configure any action Keycloak supports; it is passed
// to Keycloak verbatim. The constants below are the common built-ins, provided
// for convenience only - the list does not need to be exhaustive or kept in
// sync with Keycloak.
type UserAction string

const (
	// UserActionUpdatePassword sets it that the user needs to change their password.
	UserActionUpdatePassword UserAction = "UPDATE_PASSWORD"
	// UserActionVerifyEmail sets it that the user needs to verify their email address.
	UserActionVerifyEmail UserAction = "VERIFY_EMAIL"
	// UserActionUpdateProfile sets it that the user needs to update their profile.
	UserActionUpdateProfile UserAction = "UPDATE_PROFILE"
	// UserActionConfigureTOTP sets it that the user needs to configure a one-time password.
	UserActionConfigureTOTP UserAction = "CONFIGURE_TOTP"
	// UserActionTermsAndConditions sets it that the user needs to accept the terms and conditions.
	UserActionTermsAndConditions UserAction = "TERMS_AND_CONDITIONS"
	// UserActionVerifyProfile sets it that the user needs to verify their profile.
	UserActionVerifyProfile UserAction = "VERIFY_PROFILE"
	// UserActionDeleteAccount sets it that the user is allowed to delete their account.
	UserActionDeleteAccount UserAction = "delete_account"
	// UserActionUpdateLocale sets it that the user needs to update their locale.
	UserActionUpdateLocale UserAction = "update_user_locale"
	// UserActionWebauthnRegister sets it that the user needs to register a WebAuthn authenticator.
	UserActionWebauthnRegister UserAction = "webauthn-register"
	// UserActionWebauthnRegisterPasswordless sets it that the user needs to register a passwordless WebAuthn authenticator.
	UserActionWebauthnRegisterPasswordless UserAction = "webauthn-register-passwordless"
	// UserActionConfigureRecoveryAuthnCodes sets it that the user needs to configure recovery authentication codes.
	UserActionConfigureRecoveryAuthnCodes UserAction = "CONFIGURE_RECOVERY_AUTHN_CODES"
)

// PIIReport is a structure of all the PersonalIdentifiableInformation contained in keycloak.
type PIIReport struct {
	UserData *libregraph.User
	Sessions []*gocloak.UserSessionRepresentation
}

// Client represents a keycloak client.
type Client interface {
	CreateUser(ctx context.Context, realm string, user *libregraph.User, userActions []UserAction) (string, error)
	SendActionsMail(ctx context.Context, realm, userID string, userActions []UserAction) error
	GetUserByUsername(ctx context.Context, realm, username string) (*libregraph.User, error)
	GetPIIReport(ctx context.Context, realm, username string) (*PIIReport, error)
}
