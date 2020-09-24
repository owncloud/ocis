package glauth

import (
	"errors"
	"fmt"

	"github.com/GeertJohan/yubigo"
	"github.com/glauth/glauth/pkg/config"
	"github.com/glauth/glauth/pkg/handler"
	"github.com/go-logr/logr"
	"github.com/nmcclain/ldap"
	"github.com/owncloud/ocis/glauth/pkg/mlogr"
)

// LdapSvc holds the ldap server struct
type LdapSvc struct {
	log      logr.Logger
	c        *config.Config
	yubiAuth *yubigo.YubiAuth
	l        *ldap.Server
}

// Server initializes the ldap server.
// It is a fork github.com/glauth/pkg/server because it would introduce a go-micro dependency upstream.
func Server(opts ...Option) (*LdapSvc, error) {
	options := newOptions(opts...)

	s := LdapSvc{
		log: mlogr.New(&options.Logger),
		c:   options.Config,
	}

	var err error

	if len(s.c.YubikeyClientID) > 0 && len(s.c.YubikeySecret) > 0 {
		s.yubiAuth, err = yubigo.NewYubiAuth(s.c.YubikeyClientID, s.c.YubikeySecret)

		if err != nil {
			return nil, errors.New("yubikey auth failed")
		}
	}

	// configure the backend
	s.l = ldap.NewServer()
	s.l.EnforceLDAP = true
	var h handler.Handler
	switch s.c.Backend.Datastore {
	/* TODO bring back file config
	case "config":
		h = handler.NewConfigHandler(
			handler.Logger(s.log),
			handler.Config(s.c),
			handler.YubiAuth(s.yubiAuth),
		)
	*/
	case "ldap":
		h = handler.NewLdapHandler(
			handler.Logger(s.log),
			handler.Config(s.c),
		)
	case "owncloud":
		h = handler.NewOwnCloudHandler(
			handler.Logger(s.log),
			handler.Config(s.c),
		)
	case "accounts":
		h = NewOCISHandler(
			AccountsService(options.AccountsService),
			GroupsService(options.GroupsService),
			Logger(options.Logger),
			Config(s.c),
		)
	default:
		return nil, fmt.Errorf("unsupported backend %s - must be 'ldap', 'owncloud' or 'accounts'", s.c.Backend.Datastore)
		//return nil, fmt.Errorf("unsupported backend %s - must be 'config', 'homed', 'ldap', 'owncloud' or 'accounts'", s.c.Backend.Datastore)
	}
	s.log.V(3).Info("Using backend", "datastore", s.c.Backend.Datastore)
	s.l.BindFunc(s.c.Backend.BaseDN, h)
	s.l.SearchFunc(s.c.Backend.BaseDN, h)
	s.l.CloseFunc(s.c.Backend.BaseDN, h)

	return &s, nil
}

// ListenAndServe listens on the TCP network address s.c.LDAP.Listen
func (s *LdapSvc) ListenAndServe() error {
	s.log.V(3).Info("LDAP server listening", "address", s.c.LDAP.Listen)
	return s.l.ListenAndServe(s.c.LDAP.Listen)
}

// ListenAndServeTLS listens on the TCP network address s.c.LDAPS.Listen
func (s *LdapSvc) ListenAndServeTLS() error {
	s.log.V(3).Info("LDAPS server listening", "address", s.c.LDAPS.Listen)
	return s.l.ListenAndServeTLS(
		s.c.LDAPS.Listen,
		s.c.LDAPS.Cert,
		s.c.LDAPS.Key,
	)
}

// Shutdown ends listeners by sending true to the ldap serves quit channel
func (s *LdapSvc) Shutdown() {
	s.l.Quit <- true
}
