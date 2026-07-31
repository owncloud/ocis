package service

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"

	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	graphdefaults "github.com/owncloud/ocis/v2/services/graph/pkg/config/defaults"
	"github.com/owncloud/ocis/v2/services/graph/pkg/identity"
	graphldap "github.com/owncloud/ocis/v2/services/graph/pkg/identity/ldap"
	"github.com/owncloud/ocis/v2/services/invitations/pkg/backends/idm"
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

	backend, err := newBackend(options)
	if err != nil {
		return nil, err
	}

	return svc{
		log:     options.Logger,
		config:  options.Config,
		backend: backend,
	}, nil
}

// newBackend selects the invitation backend based on configuration.
func newBackend(options Options) (Backend, error) {
	switch options.Config.Backend {
	case "ldap", "cs3", "idm":
		return newLDAPBackend(options)
	default:
		return keycloak.New(
			options.Logger,
			options.Config.Keycloak.BasePath,
			options.Config.Keycloak.ClientID,
			options.Config.Keycloak.ClientSecret,
			options.Config.Keycloak.ClientRealm,
			options.Config.Keycloak.UserRealm,
			options.Config.Keycloak.InsecureSkipVerify,
		), nil
	}
}

// newLDAPBackend builds the backend that provisions guests directly into the
// oCIS identity backend. The LDAP schema defaults are reused from the graph
// service so invitations and graph resolve the same directory entries; only the
// connection and write settings come from the invitations configuration.
func newLDAPBackend(options Options) (Backend, error) {
	logger := options.Logger
	cfg := options.Config.LDAP

	lc := graphdefaults.DefaultConfig().Identity.LDAP
	if cfg.URI != "" {
		lc.URI = cfg.URI
	}
	if cfg.BindDN != "" {
		lc.BindDN = cfg.BindDN
	}
	if cfg.BindPassword != "" {
		lc.BindPassword = cfg.BindPassword
	}
	if cfg.CACert != "" {
		lc.CACert = cfg.CACert
	}
	if cfg.UserBaseDN != "" {
		lc.UserBaseDN = cfg.UserBaseDN
	}
	lc.Insecure = cfg.Insecure
	lc.WriteEnabled = cfg.WriteEnabled

	tlsConf := &tls.Config{InsecureSkipVerify: lc.Insecure} //nolint:gosec
	if !lc.Insecure && lc.CACert != "" {
		certs := x509.NewCertPool()
		pemData, err := os.ReadFile(lc.CACert)
		if err != nil {
			return nil, fmt.Errorf("invitations: reading LDAP CA cert %q: %w", lc.CACert, err)
		}
		if !certs.AppendCertsFromPEM(pemData) {
			return nil, fmt.Errorf("invitations: adding LDAP CA cert %q failed", lc.CACert)
		}
		tlsConf.RootCAs = certs
	}

	conn := graphldap.NewLDAPWithReconnect(&logger, graphldap.Config{
		URI:          lc.URI,
		BindDN:       lc.BindDN,
		BindPassword: lc.BindPassword,
		TLSConfig:    tlsConf,
	})
	lb, err := identity.NewLDAPBackend(conn, lc, &logger, "", "")
	if err != nil {
		return nil, fmt.Errorf("invitations: initializing LDAP backend: %w", err)
	}
	return idm.New(logger, lb), nil
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
