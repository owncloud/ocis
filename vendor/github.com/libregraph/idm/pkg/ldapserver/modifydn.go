package ldapserver

import (
	"errors"
	"fmt"
	"net"

	ber "github.com/go-asn1-ber/asn1-ber"
	"github.com/go-ldap/ldap/v3"
)

func HandleModifyDNRequest(req *ber.Packet, boundDN string, server *Server, conn net.Conn) error {
	if boundDN == "" {
		return ldap.NewError(ldap.LDAPResultInsufficientAccessRights, errors.New("anonymous Write denied"))
	}

	modDNReq, err := parseModifyDNRequest(req)
	if err != nil {
		return err
	}

	fnNames := []string{}
	for k := range server.ModifyDNFns {
		fnNames = append(fnNames, k)
	}
	fn := routeFunc(modDNReq.DN, fnNames)
	var rename Renamer
	if rename = server.ModifyDNFns[fn]; rename == nil {
		if fn == "" {
			err = fmt.Errorf("no suitable handler found for dn: '%s'", modDNReq.DN)
		} else {
			err = fmt.Errorf("handler '%s' does not support rename", fn)
		}
		return ldap.NewError(ldap.LDAPResultUnwillingToPerform, err)
	}
	code, err := rename.ModifyDN(boundDN, modDNReq, conn)
	return ldap.NewError(uint16(code), err)
	return ldap.NewError(ldap.LDAPResultProtocolError, errors.New("invalid ModifyDN request"))
}

func parseModifyDNRequest(req *ber.Packet) (*ldap.ModifyDNRequest, error) {
	modDNReq := ldap.ModifyDNRequest{}
	// LDAP ModifyDN requests have up to 4 Elements (DN, newRDN, deleteOld flag and new superior)
	if len(req.Children) < 3 || len(req.Children) > 4 {
		return nil, ldap.NewError(ldap.LDAPResultProtocolError, errors.New("invalid ModifyDN request"))
	}

	dn, ok := req.Children[0].Value.(string)
	if !ok {
		return nil, ldap.NewError(ldap.LDAPResultProtocolError, errors.New("error decoding entry DN"))
	}

	_, err := ldap.ParseDN(dn)
	if err != nil {
		return nil, ldap.NewError(ldap.LDAPResultProtocolError, err)
	}
	modDNReq.DN = dn

	newRDN, ok := req.Children[1].Value.(string)
	if !ok {
		return nil, ldap.NewError(ldap.LDAPResultProtocolError, errors.New("error decoding entry new RDN"))
	}

	_, err = ldap.ParseDN(newRDN)
	if err != nil {
		return nil, ldap.NewError(ldap.LDAPResultProtocolError, err)
	}
	modDNReq.NewRDN = newRDN

	removeOld, ok := req.Children[2].Value.(bool)
	if !ok {
		return nil, ldap.NewError(ldap.LDAPResultProtocolError, errors.New("error decoding 'deleteOldRDN` flag"))
	}

	modDNReq.DeleteOldRDN = removeOld

	// moving to a new subtree is not yet supported
	if len(req.Children) == 4 {
		return nil, ldap.NewError(ldap.LDAPResultUnwillingToPerform, errors.New("moving to 'newSuperior' is not implemented"))
	}

	return &modDNReq, nil
}
