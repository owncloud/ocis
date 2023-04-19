package ldapserver

import (
	"errors"
	"fmt"
	"net"

	ber "github.com/go-asn1-ber/asn1-ber"
	"github.com/go-ldap/ldap/v3"
)

func HandleDeleteRequest(req *ber.Packet, boundDN string, server *Server, conn net.Conn) error {
	if boundDN == "" {
		return ldap.NewError(ldap.LDAPResultInsufficientAccessRights, errors.New("anonymous Write denied"))
	}
	delReq, err := parseDeleteRequest(req)
	if err != nil {
		return err
	}
	fnNames := []string{}
	for k := range server.DeleteFns {
		fnNames = append(fnNames, k)
	}
	fn := routeFunc(delReq.DN, fnNames)
	var del Deleter
	if del = server.DeleteFns[fn]; del == nil {
		if fn == "" {
			err = fmt.Errorf("no suitable handler found for dn: '%s'", delReq.DN)
		} else {
			err = fmt.Errorf("handler '%s' does not support add", fn)
		}
		return ldap.NewError(ldap.LDAPResultUnwillingToPerform, err)
	}
	code, err := del.Delete(boundDN, delReq, conn)
	return ldap.NewError(uint16(code), err)
}

func parseDeleteRequest(req *ber.Packet) (*ldap.DelRequest, error) {
	delReq := ldap.DelRequest{}
	// LDAP Delete requests contain just the DN (no Sequence, or set)
	// i.e. they have no childre
	if len(req.Children) != 0 {
		return nil, ldap.NewError(ldap.LDAPResultProtocolError, errors.New("invalid delete request"))
	}
	dn := req.Data.String()

	_, err := ldap.ParseDN(dn)
	if err != nil {
		return nil, ldap.NewError(ldap.LDAPResultProtocolError, err)
	}
	delReq.DN = dn

	return &delReq, nil
}
