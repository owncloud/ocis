package keycloak

import (
	"context"

	"github.com/Nerzal/gocloak/v13"
	libregraph "github.com/owncloud/libre-graph-api-go"
)

// UserAction defines a type for user actions
type UserAction int8

// An incomplete list of UserActions. The string values match Keycloak's
// built-in required-action provider IDs.
const (
	// UserActionUpdatePassword sets it that the user needs to change their password.
	UserActionUpdatePassword UserAction = iota
	// UserActionVerifyEmail sets it that the user needs to verify their email address.
	UserActionVerifyEmail
	// UserActionUpdateProfile sets it that the user needs to update their profile.
	UserActionUpdateProfile
	// UserActionConfigureTOTP sets it that the user needs to configure a one-time password.
	UserActionConfigureTOTP
	// UserActionTermsAndConditions sets it that the user needs to accept the terms and conditions.
	UserActionTermsAndConditions
	// UserActionVerifyProfile sets it that the user needs to verify their profile.
	UserActionVerifyProfile
	// UserActionDeleteAccount sets it that the user is allowed to delete their account.
	UserActionDeleteAccount
	// UserActionUpdateLocale sets it that the user needs to update their locale.
	UserActionUpdateLocale
	// UserActionWebauthnRegister sets it that the user needs to register a WebAuthn authenticator.
	UserActionWebauthnRegister
	// UserActionWebauthnRegisterPasswordless sets it that the user needs to register a passwordless WebAuthn authenticator.
	UserActionWebauthnRegisterPasswordless
	// UserActionConfigureRecoveryAuthnCodes sets it that the user needs to configure recovery authentication codes.
	UserActionConfigureRecoveryAuthnCodes
)

// A lookup table to translate user actions to their string equivalents
var userActionsToString = map[UserAction]string{
	UserActionUpdatePassword:               "UPDATE_PASSWORD",
	UserActionVerifyEmail:                  "VERIFY_EMAIL",
	UserActionUpdateProfile:                "UPDATE_PROFILE",
	UserActionConfigureTOTP:                "CONFIGURE_TOTP",
	UserActionTermsAndConditions:           "TERMS_AND_CONDITIONS",
	UserActionVerifyProfile:                "VERIFY_PROFILE",
	UserActionDeleteAccount:                "delete_account",
	UserActionUpdateLocale:                 "update_user_locale",
	UserActionWebauthnRegister:             "webauthn-register",
	UserActionWebauthnRegisterPasswordless: "webauthn-register-passwordless",
	UserActionConfigureRecoveryAuthnCodes:  "CONFIGURE_RECOVERY_AUTHN_CODES",
}

// stringToUserAction is the reverse lookup of userActionsToString. It is derived
// from userActionsToString so the two stay in sync.
var stringToUserAction = func() map[string]UserAction {
	m := make(map[string]UserAction, len(userActionsToString))
	for action, s := range userActionsToString {
		m[s] = action
	}
	return m
}()

// UserActionFromString returns the UserAction matching the given Keycloak
// required-action string (e.g. "UPDATE_PASSWORD"). The boolean return value
// reports whether the string mapped to a known action.
func UserActionFromString(s string) (UserAction, bool) {
	action, ok := stringToUserAction[s]
	return action, ok
}

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
