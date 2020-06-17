package glauth

import (
	"errors"

	"github.com/GeertJohan/yubigo"
	"github.com/glauth/glauth/pkg/config"
	"github.com/glauth/glauth/pkg/handler"
	"github.com/go-logr/logr"
	"github.com/nmcclain/ldap"
	"github.com/owncloud/ocis-glauth/pkg/mlogr"
)

// LdapSvc holds the ldap server struct
type LdapSvc struct {
	log      logr.Logger
	c        *config.Config
	yubiAuth *yubigo.YubiAuth
	l        *ldap.Server
}

// Server initializes the debug service and server.
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
			return nil, errors.New("Yubikey Auth failed")
		}
	}

	// configure the backend
	s.l = ldap.NewServer()
	s.l.EnforceLDAP = true
	var h handler.Handler
	h = NewOCISHandler(
		AccountsService(options.AccountsService),
		Logger(options.Logger),
		Config(s.c),
	)
	s.l.BindFunc("", h)
	s.l.SearchFunc("", h)
	s.l.CloseFunc("", h)

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
