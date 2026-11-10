// Package keycloak offers an invitation backend for the invitation service.
package keycloak

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/Nerzal/gocloak/v13"
	"github.com/google/uuid"
	libregraph "github.com/owncloud/libre-graph-api-go"
	"github.com/owncloud/ocis/v2/ocis-pkg/keycloak"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/services/invitations/pkg/invitations"
)

const (
	userType = "Guest"
)

// defaultUserActions are the Keycloak required actions used when no valid
// actions are configured. They ensure an invited guest can always set a
// password and verify their email address.
var defaultUserActions = []keycloak.UserAction{
	keycloak.UserActionUpdatePassword,
	keycloak.UserActionVerifyEmail,
}

// Backend represents the keycloak backend.
type Backend struct {
	logger      log.Logger
	client      keycloak.Client
	userRealm   string
	userActions []keycloak.UserAction
}

// New instantiates a new keycloak.Backend with a default gocloak client.
func New(
	logger log.Logger,
	baseURL, clientID, clientSecret, clientRealm, userRealm string,
	insecureSkipVerify bool,
	executeActions []string,
) *Backend {
	logger = log.Logger{
		Logger: logger.With().
			Str("invitationBackend", "keycloak").
			Str("clientID", clientID).
			Str("clientRealm", clientRealm).
			Str("userRealm", userRealm).
			Logger(),
	}
	client := keycloak.New(baseURL, clientID, clientSecret, clientRealm, insecureSkipVerify)
	return NewWithClient(logger, client, userRealm, parseUserActions(executeActions))
}

// NewWithClient creates a new backend with the supplied keycloak client.
func NewWithClient(
	logger log.Logger,
	client keycloak.Client,
	userRealm string,
	userActions []keycloak.UserAction,
) *Backend {
	if len(userActions) == 0 {
		userActions = defaultUserActions
	}
	return &Backend{
		logger:      logger,
		client:      client,
		userRealm:   userRealm,
		userActions: userActions,
	}
}

// CreateUser creates a user in the keycloak backend.
func (b Backend) CreateUser(ctx context.Context, invitation *invitations.Invitation) (string, error) {
	u := uuid.New()

	b.logger.Info().
		Str("email", invitation.InvitedUserEmailAddress).
		Msg("Creating new user")
	user := &libregraph.User{
		Mail:                     &invitation.InvitedUserEmailAddress,
		AccountEnabled:           boolP(true),
		OnPremisesSamAccountName: invitation.InvitedUserEmailAddress,
		Id:                       stringP(u.String()),
		UserType:                 stringP(userType),
	}

	id, err := b.client.CreateUser(ctx, b.userRealm, user, b.userActions)
	if err != nil {
		b.logger.Error().
			Str("userID", u.String()).
			Str("email", invitation.InvitedUserEmailAddress).
			Err(err).
			Msg("Failed to create user")
		return "", err
	}

	// Record the created user on the invitation so the service can return it as
	// the invitedUser relation of the Graph invitation response. user.Id holds the
	// OWNCLOUD_ID, i.e. the identity oCIS uses once the guest is provisioned locally
	// (not the Keycloak-internal id returned for sending the actions mail).
	invitation.InvitedUser = user

	return id, nil
}

// CanSendMail returns true because keycloak does allow sending mail.
func (b Backend) CanSendMail() bool { return true }

// SendMail sends a mail to the user with details on how to redeem the invitation.
func (b Backend) SendMail(ctx context.Context, id string) error {
	err := b.client.SendActionsMail(ctx, b.userRealm, id, b.userActions)
	if isDefiniteRejection(err) {
		return fmt.Errorf("%w: %w", invitations.ErrRejected, err)
	}
	return err
}

// isDefiniteRejection reports whether err is an error response from Keycloak
// that certainly means the request was refused and no mail was sent. Transport
// errors and cancelled contexts are not, as Keycloak may have sent the mail
// anyway. The same holds for 502 and 504: they typically come from a reverse
// proxy in front of Keycloak, which may already have processed the request.
func isDefiniteRejection(err error) bool {
	var apiErr *gocloak.APIError
	if !errors.As(err, &apiErr) || apiErr.Code <= 0 {
		return false
	}
	switch apiErr.Code {
	case http.StatusBadGateway, http.StatusGatewayTimeout:
		return false
	}
	return true
}

// DeleteUser deletes the user with the given identifier from the keycloak backend.
func (b Backend) DeleteUser(ctx context.Context, id string) error {
	return b.client.DeleteUser(ctx, b.userRealm, id)
}

// parseUserActions converts the configured Keycloak required-action strings
// into UserActions. Each non-empty value is passed through to Keycloak as-is
// (Keycloak validates the actions), so operators can configure any action their
// Keycloak supports without this list having to be kept in sync. If none are
// configured, the defaults (UPDATE_PASSWORD, VERIFY_EMAIL) are used so an
// invited guest always has a way to set up their account.
func parseUserActions(executeActions []string) []keycloak.UserAction {
	actions := make([]keycloak.UserAction, 0, len(executeActions))
	for _, a := range executeActions {
		a = strings.TrimSpace(a)
		if a == "" {
			continue
		}
		actions = append(actions, keycloak.UserAction(a))
	}
	if len(actions) == 0 {
		return defaultUserActions
	}
	return actions
}

func boolP(b bool) *bool       { return &b }
func stringP(s string) *string { return &s }
