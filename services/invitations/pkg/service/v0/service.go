package service

import (
	"context"

	"github.com/owncloud/ocis/v2/ocis-pkg/log"
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

// New returns a new instance of Service
func New(opts ...Option) (Service, error) {
	options := newOptions(opts...)

	return svc{
		log:    options.Logger,
		config: options.Config,
	}, nil
}

type svc struct {
	config *config.Config
	log    log.Logger
}

// Invite implements the service interface
func (s svc) Invite(ctx context.Context, invitation *invitations.Invitation) (*invitations.Invitation, error) {
	return &invitations.Invitation{
		InvitedUserDisplayName: "Yay",
	}, nil
}
