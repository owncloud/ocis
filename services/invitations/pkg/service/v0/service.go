package service

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"net/http"

	libregraph "github.com/owncloud/libre-graph-api-go"
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

	if invitation == nil {
		return nil, ErrBadRequest
	}

	if invitation.InvitedUserEmailAddress == "" {
		return nil, ErrMissingEmail
	}

	user := &libregraph.User{
		Mail: &invitation.InvitedUserEmailAddress,
		// TODO we cannot set the user type here
	}

	if invitation.InvitedUserDisplayName != "" {
		user.DisplayName = &invitation.InvitedUserDisplayName
	}
	// we don't really need a username as guests have to log in with their email address anyway
	// what if later a user is provisioned with a guest accounts email address?

	data, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}

	// send a request to the provisioning endpoint
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true /*TODO make configurable*/},
	}
	client := &http.Client{Transport: tr}

	req, err := http.NewRequest("POST", "/graph/v1.0/users", bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	// TODO either forward current user token or use bearer token?
	req.Header.Set("Authorization", "Bearer some-token")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	invitedUser := &libregraph.User{}
	err = json.NewDecoder(res.Body).Decode(invitedUser)
	if err != nil {
		return nil, err
	}

	response := &invitations.Invitation{
		InvitedUser: invitedUser,
	}
	if res.StatusCode == http.StatusCreated {
		response.Status = "Completed"
	}

	// optionally send an email

	return response, nil
}
