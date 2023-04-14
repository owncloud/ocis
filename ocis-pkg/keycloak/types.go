package keycloak

import (
	"context"

	"github.com/Nerzal/gocloak/v13"
	libregraph "github.com/owncloud/libre-graph-api-go"
)

// UserAction defines a type for user actions
type UserAction int8

// An incomplete list of UserActions
const (
	// UserActionUpdatePassword sets it that the user needs to change their password.
	UserActionUpdatePassword UserAction = iota
	// UserActionVerifyEmail sets it that the user needs to verify their email address.
	UserActionVerifyEmail
)

// A lookup table to translate user actions to their string equivalents
var userActionsToString = map[UserAction]string{
	UserActionUpdatePassword: "UPDATE_PASSWORD",
	UserActionVerifyEmail:    "VERIFY_EMAIL",
}

// PIIReport is a structure of all the PersonalIdentifiableInformation contained in keycloak.
type PIIReport struct {
	UserData    *libregraph.User                    `json:"user_data,omitempty"`
	Credentials []*gocloak.CredentialRepresentation `json:"credentials,omitempty"`
}

// Client represents a keycloak client.
type Client interface {
	CreateUser(ctx context.Context, realm string, user *libregraph.User, userActions []UserAction) (string, error)
	SendActionsMail(ctx context.Context, realm, userID string, userActions []UserAction) error
	GetUserByEmail(ctx context.Context, realm, email string) (*libregraph.User, error)
	GetPIIReport(ctx context.Context, realm string, email string) (*PIIReport, error)
}
