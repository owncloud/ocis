package glauth

import (
	"errors"
	"fmt"

	"github.com/GeertJohan/yubigo"
	"github.com/glauth/glauth/v2/pkg/config"
	"github.com/glauth/glauth/v2/pkg/handler"
	"github.com/go-logr/logr"
	"github.com/nmcclain/ldap"
	"github.com/owncloud/ocis/glauth/pkg/mlogr"
)

// LdapSvc holds the ldap server struct
type LdapSvc struct {
	log      logr.Logger
	ldap     *config.LDAP
	ldaps    *config.LDAPS
	backend  *config.Config
	fallback *config.Config
	yubiAuth *yubigo.YubiAuth
	l        *ldap.Server
}

// Server initializes the ldap server.
// It is a fork github.com/glauth/pkg/server because it would introduce a go-micro dependency upstream.
func Server(opts ...Option) (*LdapSvc, error) {
	options := newOptions(opts...)

	s := LdapSvc{
		log:      mlogr.New(&options.Logger),
		backend:  options.Backend,
		fallback: options.Fallback,
		ldap:     options.LDAP,
		ldaps:    options.LDAPS,
	}

	var err error

	if len(s.backend.YubikeyClientID) > 0 && len(s.backend.YubikeySecret) > 0 {
		s.yubiAuth, err = yubigo.NewYubiAuth(s.backend.YubikeyClientID, s.backend.YubikeySecret)

		if err != nil {
			return nil, errors.New("yubikey auth failed")
		}
	}

	// configure the backend
	s.l = ldap.NewServer()
	s.l.EnforceLDAP = true
	var bh handler.Handler

	switch s.backend.Backend.Datastore {
	/* TODO bring back file config
	case "config":
		bh = handler.NewConfigHandler(
			handler.Logger(s.log),
			handler.Config(s.c),
			handler.YubiAuth(s.yubiAuth),
		)
	*/
	case "ldap":
		bh = handler.NewLdapHandler(
			handler.Logger(s.log),
			handler.Backend(s.backend.Backend),
		)
	case "owncloud":
		bh = handler.NewOwnCloudHandler(
			handler.Logger(s.log),
			handler.Backend(s.backend.Backend),
		)
	case "accounts":
		bh = NewOCISHandler(
			AccountsService(options.AccountsService),
			GroupsService(options.GroupsService),
			Logger(options.Logger),
			BaseDN(s.backend.Backend.BaseDN),
			NameFormat(s.backend.Backend.NameFormat),
			GroupFormat(s.backend.Backend.GroupFormat),
			RoleBundleUUID(options.RoleBundleUUID),
		)
	default:
		return nil, fmt.Errorf("unsupported backend %s - must be 'ldap', 'owncloud' or 'accounts'", s.backend.Backend.Datastore)
	}
	s.log.V(3).Info("Using backend", "backend", s.backend.Backend)

	if s.fallback != nil && s.fallback.Backend.Datastore != "" {

		var fh handler.Handler

		switch s.fallback.Backend.Datastore {
		/* TODO bring back file config
		case "config":
			fh = handler.NewConfigHandler(
				handler.Logger(s.log),
				handler.Config(s.c),
				handler.YubiAuth(s.yubiAuth),
			)
		*/
		case "ldap":
			fh = handler.NewLdapHandler(
				handler.Logger(s.log),
				handler.Backend(s.fallback.Backend),
			)
		case "owncloud":
			fh = handler.NewOwnCloudHandler(
				handler.Logger(s.log),
				handler.Backend(s.fallback.Backend),
			)
		case "accounts":
			fh = NewOCISHandler(
				AccountsService(options.AccountsService),
				GroupsService(options.GroupsService),
				Logger(options.Logger),
				BaseDN(s.fallback.Backend.BaseDN),
				NameFormat(s.fallback.Backend.NameFormat),
				GroupFormat(s.fallback.Backend.GroupFormat),
				RoleBundleUUID(options.RoleBundleUUID),
			)
		default:
			return nil, fmt.Errorf("unsupported fallback %s - must be 'ldap', 'owncloud' or 'accounts'", s.fallback.Backend.Datastore)
		}
		s.log.V(3).Info("Using fallback", "backend", s.fallback.Backend)

		bh = NewChainHandler(options.Logger, bh, fh)
	}

	s.l.BindFunc(s.backend.Backend.BaseDN, bh)
	s.l.SearchFunc(s.backend.Backend.BaseDN, bh)
	s.l.CloseFunc(s.backend.Backend.BaseDN, bh)

	return &s, nil
}

// ListenAndServe listens on the TCP network address s.c.LDAP.Listen
func (s *LdapSvc) ListenAndServe() error {
	s.log.V(3).Info("ldap server listening", "address", s.ldap.Listen)
	return s.l.ListenAndServe(s.ldap.Listen)
}

// ListenAndServeTLS listens on the TCP network address s.c.LDAPS.Listen
func (s *LdapSvc) ListenAndServeTLS() error {
	s.log.V(3).Info("ldaps server listening", "address", s.ldaps.Listen)
	return s.l.ListenAndServeTLS(
		s.ldaps.Listen,
		s.ldaps.Cert,
		s.ldaps.Key,
	)
}

// Shutdown ends listeners by sending true to the ldap serves quit channel
func (s *LdapSvc) Shutdown() {
	s.l.Quit <- true
}
