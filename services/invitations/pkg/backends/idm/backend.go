// Package idm provides an invitation backend that provisions invited guests
// directly into the oCIS identity backend (the "IDM" provisioning mode described
// in the invitations service README). Unlike the Keycloak backend, which creates
// the guest in an external IdP, this backend writes the guest into the directory
// that oCIS itself reads (via the Graph identity backend), so the invited guest
// is immediately resolvable through the Graph API and therefore an immediately
// shareable principal — closing the provisioning-delay gap that prevents
// "invite by email, then share now".
//
// PROTOTYPE SCOPE / LIMITATION: this backend only *provisions* the guest. It does
// not send a credential-setup email, because the local identity backend has no
// execute-actions-email equivalent (that is a Keycloak feature). It therefore
// fits IDM-managed deployments. In a Keycloak-backed deployment, the guest also
// needs a Keycloak account to authenticate; there the better lever is creating
// the guest in Keycloak with a WRITABLE LDAP federation, so a single write
// satisfies both authentication and directory resolution.
package idm

import (
	"context"

	"github.com/google/uuid"
	libregraph "github.com/owncloud/libre-graph-api-go"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/services/invitations/pkg/invitations"
)

const userType = "Guest"

// DirectoryProvisioner is the subset of the Graph identity backend that this
// invitation backend needs. *identity.LDAP (services/graph/pkg/identity)
// satisfies it, so in production the guest is created in the same directory oCIS
// reads for share-recipient resolution.
type DirectoryProvisioner interface {
	CreateUser(ctx context.Context, user libregraph.User) (*libregraph.User, error)
}

// Backend provisions invited guests into the oCIS identity backend.
type Backend struct {
	logger      log.Logger
	provisioner DirectoryProvisioner
}

// New instantiates an IDM invitation backend backed by the given provisioner.
func New(logger log.Logger, provisioner DirectoryProvisioner) *Backend {
	return &Backend{
		logger:      log.Logger{Logger: logger.With().Str("invitationBackend", "idm").Logger()},
		provisioner: provisioner,
	}
}

// CreateUser provisions the invited guest as a local user (userType "Guest") and
// records the created user on the invitation as invitedUser. The returned id is
// the directory id of the created user, which is immediately usable as a share
// recipient.
func (b Backend) CreateUser(ctx context.Context, invitation *invitations.Invitation) (string, error) {
	id := uuid.New().String()

	displayName := invitation.InvitedUserDisplayName
	if displayName == "" {
		displayName = invitation.InvitedUserEmailAddress
	}

	b.logger.Info().
		Str("email", invitation.InvitedUserEmailAddress).
		Msg("Provisioning new guest in the identity backend")

	user := libregraph.User{
		Id:                       &id,
		Mail:                     &invitation.InvitedUserEmailAddress,
		DisplayName:              displayName,
		OnPremisesSamAccountName: invitation.InvitedUserEmailAddress,
		AccountEnabled:           boolP(true),
		UserType:                 stringP(userType),
	}

	created, err := b.provisioner.CreateUser(ctx, user)
	if err != nil {
		b.logger.Error().
			Str("email", invitation.InvitedUserEmailAddress).
			Err(err).
			Msg("Failed to provision guest")
		return "", err
	}

	// Record the created user so the service returns it as invitedUser, and hand
	// back its directory id (immediately resolvable / shareable).
	invitation.InvitedUser = created
	return created.GetId(), nil
}

// CanSendMail reports false: the IDM backend only provisions the guest and has no
// credential-setup mail flow (see the package doc).
func (b Backend) CanSendMail() bool { return false }

// SendMail is a no-op for the IDM backend.
func (b Backend) SendMail(_ context.Context, _ string) error { return nil }

func boolP(b bool) *bool       { return &b }
func stringP(s string) *string { return &s }
