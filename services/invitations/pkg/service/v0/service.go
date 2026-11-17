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
	// createUserTimeout bounds creating the guest in the backend.
	createUserTimeout = 10 * time.Second
	// sendMailTimeout bounds sending the invitation mail. Keycloak sends the
	// mail via SMTP before it answers, so this is longer than the others.
	sendMailTimeout = 30 * time.Second
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

	// Creating the guest, sending the mail and rolling back belong together:
	// once the guest has been created, an aborted mail call leaves a guest
	// that can not be invited again. So they do not inherit the cancellation
	// of the request (e.g. a client disconnect), and each call is bounded by
	// its own timeout instead.
	ctx = context.WithoutCancel(ctx)

	id, err := s.createUser(ctx, invitation)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrBackend, err)
	}
	if id == "" {
		s.log.Error().
			Str("email", invitation.InvitedUserEmailAddress).
			Msgf("The backend returned no user id for the guest %s. The guest may have been created but can not be addressed, no invitation mail was sent. Check the backend and remove the guest there before inviting again", invitation.InvitedUserEmailAddress)
		return nil, fmt.Errorf("%w: %s", ErrBackend, "no user id returned for the created user")
	}

	// As we only have a single backend, and that backend supports email, we don't have
	// any code to handle mailing ourself yet.
	if s.backend.CanSendMail() {
		if err := s.sendMail(ctx, id); err != nil {
			s.rollback(ctx, invitation.InvitedUserEmailAddress, id, err)
			return nil, fmt.Errorf("%w: %s", ErrBackend, err)
		}
	}

	return invitation, nil
}

func (s svc) createUser(ctx context.Context, invitation *invitations.Invitation) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, createUserTimeout)
	defer cancel()
	return s.backend.CreateUser(ctx, invitation)
}

func (s svc) sendMail(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, sendMailTimeout)
	defer cancel()
	return s.backend.SendMail(ctx, id)
}

// rollback deletes the user created for the invitation after sending the
// invitation mail failed with mailErr. The guest can not redeem the invitation
// without the mail, and a leftover user would make every retry fail with a
// conflict. It only deletes if the failure is certain: if the outcome is
// unknown (e.g. timeout), the mail may have been sent and the invitation is
// working, so the user is kept. ctx must not carry the request cancellation.
func (s svc) rollback(ctx context.Context, email, id string, mailErr error) {
	if !errors.Is(mailErr, invitations.ErrRejected) {
		// Logged as error, not warning: the invitation request failed and may
		// need an operator to act, as the guest blocks inviting the address
		// again if the mail did not arrive.
		s.log.Error().
			Str("userID", id).
			Str("email", email).
			Err(mailErr).
			Msgf("Sending the invitation mail failed with an unknown outcome. The guest was kept, the invitation mail may or may not have been sent. If it did not arrive, remove the guest %s (user id %s) in Keycloak before inviting the address again", email, id)
		return
	}

	ctx, cancel := context.WithTimeout(ctx, rollbackTimeout)
	defer cancel()
	if err := s.backend.DeleteUser(ctx, id); err != nil {
		s.log.Error().
			Str("userID", id).
			Str("email", email).
			Err(err).
			AnErr("mailError", mailErr).
			Msgf("Sending the invitation mail failed and removing the guest %s (user id %s) failed. Remove the guest manually, otherwise inviting this address again fails", email, id)
	}
}
