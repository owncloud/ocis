// Package keycloak offers an invitation backend for the invitation service.
package keycloak

import (
	"context"
	"crypto/tls"
	"fmt"
	"strings"

	"github.com/Nerzal/gocloak/v13"
	"github.com/google/uuid"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/services/invitations/pkg/invitations"
)

const (
	idAttr       = "OWNCLOUD_ID"
	userTypeAttr = "OWNCLOUD_USER_TYPE"
	userTypeVal  = "Guest"
)

var userRequiredActions = []string{"UPDATE_PASSWORD", "VERIFY_EMAIL"}

// Backend represents the keycloak backend.
type Backend struct {
	logger       log.Logger
	client       GoCloak
	clientID     string
	clientSecret string
	clientRealm  string
	userRealm    string
}

// New instantiates a new keycloak.Backend with a default gocloak client.
func New(
	logger log.Logger,
	baseURL, clientID, clientSecret, clientRealm, userRealm string,
	insecureSkipVerify bool,
) *Backend {
	client := gocloak.NewClient(baseURL)
	restyClient := client.RestyClient()
	restyClient.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: insecureSkipVerify}) //nolint:gosec
	return NewWithClient(logger, client, clientID, clientSecret, clientRealm, userRealm)
}

// NewWithClient creates a new backend with the supplied GoCloak client.
func NewWithClient(
	logger log.Logger,
	client GoCloak,
	clientID, clientSecret, clientRealm, userRealm string,
) *Backend {
	return &Backend{
		logger: log.Logger{
			Logger: logger.With().
				Str("invitationBackend", "keycloak").
				Str("clientID", clientID).
				Str("clientRealm", clientRealm).
				Str("userRealm", userRealm).
				Logger(),
		},
		client:       client,
		clientID:     clientID,
		clientSecret: clientSecret,
		clientRealm:  clientRealm,
		userRealm:    userRealm,
	}
}

// CreateUser creates a user in the keycloak backend.
func (b Backend) CreateUser(ctx context.Context, invitation *invitations.Invitation) (string, error) {
	token, err := b.getToken(ctx)
	if err != nil {
		return "", err
	}
	u := uuid.New()

	firstName, lastName := splitDisplayName(invitation.InvitedUserDisplayName)
	b.logger.Info().
		Str(idAttr, u.String()).
		Str("email", invitation.InvitedUserEmailAddress).
		Msg("Creating new user")
	user := gocloak.User{
		FirstName: &firstName,
		LastName:  &lastName,
		Email:     &invitation.InvitedUserEmailAddress,
		Enabled:   gocloak.BoolP(true),
		Username:  &invitation.InvitedUserEmailAddress,
		Attributes: &map[string][]string{
			idAttr:       {u.String()},
			userTypeAttr: {userTypeVal},
		},
		RequiredActions: &userRequiredActions,
	}

	id, err := b.client.CreateUser(ctx, token.AccessToken, b.userRealm, user)
	if err != nil {
		b.logger.Error().
			Str(idAttr, u.String()).
			Str("email", invitation.InvitedUserEmailAddress).
			Err(err).
			Msg("Failed to create user")
		return "", err
	}

	return id, nil
}

// CanSendMail returns true because keycloak does allow to send mail.
func (b Backend) CanSendMail() bool { return true }

// SendMail sends a mail to the user with details on how to reedeem the invitation.
func (b Backend) SendMail(ctx context.Context, id string) error {
	token, err := b.getToken(ctx)
	if err != nil {
		return err
	}
	params := gocloak.ExecuteActionsEmail{
		UserID:  &id,
		Actions: &userRequiredActions,
	}
	return b.client.ExecuteActionsEmail(ctx, token.AccessToken, b.userRealm, params)
}

func (b Backend) getToken(ctx context.Context) (*gocloak.JWT, error) {
	b.logger.Debug().Msg("Logging into keycloak")
	token, err := b.client.LoginClient(ctx, b.clientID, b.clientSecret, b.clientRealm)
	if err != nil {
		b.logger.Error().Err(err).Msg("failed to get token")
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	rRes, err := b.client.RetrospectToken(ctx, token.AccessToken, b.clientID, b.clientSecret, b.clientRealm)
	if err != nil {
		b.logger.Error().Err(err).Msg("failed to introspect token")
		return nil, fmt.Errorf("failed to retrospect token: %w", err)
	}

	if !*rRes.Active {
		b.logger.Error().Msg("token not active")
		return nil, fmt.Errorf("token is not active")
	}

	return token, nil
}

// Quick and dirty way to split the last name off from the first name(s), imperfect, because
// every culture has a different conception of names.
func splitDisplayName(displayName string) (string, string) {
	parts := strings.Split(displayName, " ")
	if len(parts) <= 1 {
		return parts[0], ""
	}

	return strings.Join(parts[:len(parts)-1], " "), parts[len(parts)-1]
}
