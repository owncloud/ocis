package keycloak

import (
	"context"

	"github.com/Nerzal/gocloak/v13"
)

// GoCloak represents the parts of gocloak.GoCloak that we use, mainly here for mockery.
type GoCloak interface {
	CreateUser(ctx context.Context, token, realm string, user gocloak.User) (string, error)
	GetUsers(ctx context.Context, token, realm string, params gocloak.GetUsersParams) ([]*gocloak.User, error)
	ExecuteActionsEmail(ctx context.Context, token, realm string, params gocloak.ExecuteActionsEmail) error
	LoginClient(ctx context.Context, clientID, clientSecret, realm string, scopes ...string) (*gocloak.JWT, error)
	RetrospectToken(ctx context.Context, accessToken, clientID, clientSecret, realm string) (*gocloak.IntroSpectTokenResult, error)
	GetCredentials(ctx context.Context, accessToken, realm, userID string) ([]*gocloak.CredentialRepresentation, error)
	GetUserSessions(ctx context.Context, token, realm, userID string) ([]*gocloak.UserSessionRepresentation, error)
}
