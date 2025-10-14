package service

import (
	"context"
	"fmt"

	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/services/invitations/pkg/backends/keycloak"
	"github.com/owncloud/ocis/v2/services/invitations/pkg/config"
	"github.com/owncloud/ocis/v2/services/invitations/pkg/invitations"
)

const (
	OwnCloudInstanceRel = "http://invitations.owncloud/rel/server-instance"
	OpenIDConnectRel    = "http://openid.net/specs/connect/1.0/issuer"
)

// Service defines the extension handlers.
type Service interface {
	// Invite creates a new invitation. Invitation adds an external user to the organization.
	//
	// When creating a new invitation you have several options available:
	// 1. On invitation creation, Microsoft Graph can automatically send an
	//    invitation email directly to the invited user, or your app can use
	//    the inviteRedeemUrl returned in the creation response to craft your
	//    own invitation (through your communication mechanism of choice) to
	//    the invited user. If you decide to have Microsoft Graph send an
	//    invitation email automatically, you can control the content and
	//    language of the email using invitedUserMessageInfo.
	// 2. When the user is invited, a user entity (of userType Guest) is
	//    created and can now be used to control access to resources. The
	//    invited user has to go through the redemption process to access any
	//    resources they have been invited to.
	Invite(ctx context.Context, invitation *invitations.Invitation) (*invitations.Invitation, error)
}

// Backend defines the behaviour of a user backend.
type Backend interface {
	// CreateUser creates a user in the backend and returns an identifier string.
	CreateUser(ctx context.Context, invitation *invitations.Invitation) (string, error)
	// CanSendMail should return true if the backend can send mail
	CanSendMail() bool
	// SendMail sends a mail to the user with details on how to reedeem the invitation.
	SendMail(ctx context.Context, identifier string) error
}

// New returns a new instance of Service
func New(opts ...Option) (Service, error) {
	options := newOptions(opts...)

	// Harcode keycloak backend for now, but this should be configurable in the future.
	backend := keycloak.New(
		options.Logger,
		options.Config.Keycloak.BasePath,
		options.Config.Keycloak.ClientID,
		options.Config.Keycloak.ClientSecret,
		options.Config.Keycloak.ClientRealm,
		options.Config.Keycloak.UserRealm,
		options.Config.Keycloak.InsecureSkipVerify,
	)

	return svc{
		log:     options.Logger,
		config:  options.Config,
		backend: backend,
	}, nil
}

type svc struct {
	config  *config.Config
	log     log.Logger
	backend Backend
}

// Invite implements the service interface
func (s svc) Invite(ctx context.Context, invitation *invitations.Invitation) (*invitations.Invitation, error) {
	if invitation == nil {
		return nil, ErrBadRequest
	}

	if invitation.InvitedUserEmailAddress == "" {
		return nil, ErrMissingEmail
	}

	id, err := s.backend.CreateUser(ctx, invitation)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrBackend, err)
	}

	// As we only have a single backend, and that backend supports email, we don't have
	// any code to handle mailing ourself yet.
	if s.backend.CanSendMail() {
		err := s.backend.SendMail(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("%w: %s", ErrBackend, err)
		}
	}

	return invitation, nil
}
