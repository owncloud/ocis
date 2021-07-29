package glauth

import (
	"net"

	"github.com/glauth/glauth/pkg/config"
	"github.com/glauth/glauth/pkg/handler"
	"github.com/nmcclain/ldap"
	"github.com/owncloud/ocis/ocis-pkg/log"
)

type chainHandler struct {
	log log.Logger
	b   handler.Handler
	f   handler.Handler
}

func (h chainHandler) Bind(bindDN, bindSimplePw string, conn net.Conn) (res ldap.LDAPResultCode, err error) {
	h.log.Debug().
		Str("binddn", bindDN).
		Interface("src", conn.RemoteAddr()).
		Str("handler", "chain").
		Msg("Bind request")
	res, err = h.b.Bind(bindDN, bindSimplePw, conn)
	switch {
	case err != nil:
		h.log.Error().
			Err(err).
			Str("binddn", bindDN).
			Interface("src", conn.RemoteAddr()).
			Str("handler", "chain").
			Msg("Bind request")
		return h.f.Bind(bindDN, bindSimplePw, conn)
	case res == ldap.LDAPResultInvalidCredentials:
		return h.f.Bind(bindDN, bindSimplePw, conn)
	}
	return
}

func (h chainHandler) Search(bindDN string, searchReq ldap.SearchRequest, conn net.Conn) (res ldap.ServerSearchResult, err error) {
	h.log.Debug().
		Str("binddn", bindDN).
		Interface("src", conn.RemoteAddr()).
		Str("handler", "chain").
		Msg("Search request")
	res, err = h.b.Search(bindDN, searchReq, conn)
	switch {
	case err != nil:
		h.log.Error().
			Err(err).
			Str("binddn", bindDN).
			Interface("src", conn.RemoteAddr()).
			Str("handler", "chain").
			Msg("Search request")
		return h.f.Search(bindDN, searchReq, conn)
	case len(res.Entries) == 0:
		// yes, we only fall back if there are no results in the first backend
		// this is not supposed to work for searching lots of users, only to look up a single user
		// searching multiple users would require merging result sets. out of scope for now.
		return h.f.Search(bindDN, searchReq, conn)
	}
	return
}
func (h chainHandler) Close(boundDN string, conn net.Conn) error {
	h.log.Debug().
		Str("boundDN", boundDN).
		Interface("src", conn.RemoteAddr()).
		Str("handler", "chain").
		Msg("Close request")
	if err := h.b.Close(boundDN, conn); err != nil {
		h.log.Error().
			Err(err).
			Str("boundDN", boundDN).
			Interface("src", conn.RemoteAddr()).
			Str("handler", "chain").
			Msg("Close request")
	}
	if err := h.f.Close(boundDN, conn); err != nil {
		h.log.Error().
			Err(err).
			Str("boundDN", boundDN).
			Interface("src", conn.RemoteAddr()).
			Str("handler", "chain").
			Msg("Close request")
	}
	return nil
}

// Add is not yet supported for the chain backend
func (h chainHandler) Add(boundDN string, req ldap.AddRequest, conn net.Conn) (result ldap.LDAPResultCode, err error) {
	return ldap.LDAPResultInsufficientAccessRights, nil
}

// Modify is not yet supported for the chain backend
func (h chainHandler) Modify(boundDN string, req ldap.ModifyRequest, conn net.Conn) (result ldap.LDAPResultCode, err error) {
	return ldap.LDAPResultInsufficientAccessRights, nil
}

// Delete is not yet supported for the chain backend
func (h chainHandler) Delete(boundDN string, deleteDN string, conn net.Conn) (result ldap.LDAPResultCode, err error) {
	return ldap.LDAPResultInsufficientAccessRights, nil
}

// FindUser with the given username. Called by the ldap backend to authenticate the bind. Optional
func (h chainHandler) FindUser(userName string) (bool, config.User, error) {
	return false, config.User{}, nil
}

// NewChainHandler implements a chain backend with two backends
func NewChainHandler(log log.Logger, bh handler.Handler, fh handler.Handler) handler.Handler {
	return chainHandler{
		log: log,
		b:   bh,
		f:   fh,
	}
}
