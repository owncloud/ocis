package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/services/invitations/pkg/backends/keycloak"
	"github.com/owncloud/ocis/v2/services/invitations/pkg/config"
	"github.com/owncloud/ocis/v2/services/invitations/pkg/invitations"
)

const (
	// rollbackTimeout bounds the cleanup of a user whose invitation failed.
	rollbackTimeout = 10 * time.Second

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
	// DeleteUser deletes a previously created user. It is used to roll back
	// the creation when the invitation can not be completed.
	DeleteUser(ctx context.Context, identifier string) error
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
		options.Config.Keycloak.ExecuteActions,
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
			s.rollback(ctx, invitation, id, err)
			return nil, fmt.Errorf("%w: %s", ErrBackend, err)
		}
	}

	return invitation, nil
}

// rollback deletes the user created for the invitation after sending the
// invitation mail failed with mailErr. The guest can not redeem the invitation
// without the mail, and a leftover user would make every retry fail with a
// conflict. It only deletes if the failure is certain: if the outcome is
// unknown (e.g. timeout), the mail may have been sent and the invitation is
// working, so the user is kept.
func (s svc) rollback(ctx context.Context, invitation *invitations.Invitation, id string, mailErr error) {
	email := invitation.InvitedUserEmailAddress

	if !errors.Is(mailErr, invitations.ErrRejected) {
		s.log.Warn().
			Str("userID", id).
			Str("email", email).
			Err(mailErr).
			Msg("Sending the invitation mail failed with an unknown outcome, the guest may exist in the backend and the mail may or may not have been sent. The guest was not removed")
		return
	}

	if id == "" {
		s.log.Error().
			Str("email", email).
			Err(mailErr).
			Msg("Sending the invitation mail failed, but the backend returned no user id, so the guest can not be removed automatically")
		return
	}

	// The request context may be the reason the mail failed, so the rollback
	// must not inherit its cancellation.
	delCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), rollbackTimeout)
	defer cancel()
	if err := s.backend.DeleteUser(delCtx, id); err != nil {
		s.log.Error().
			Str("userID", id).
			Str("email", email).
			Err(err).
			AnErr("mailError", mailErr).
			Msgf("Sending the invitation mail failed and removing the guest %s (user id %s) failed. Remove the guest manually, otherwise inviting this address again fails", email, id)
		return
	}
	// The invited user no longer exists.
	invitation.InvitedUser = nil
}
