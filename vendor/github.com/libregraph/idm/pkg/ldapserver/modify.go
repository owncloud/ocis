package ldapserver

import (
	"errors"
	"fmt"
	"log"
	"net"

	ber "github.com/go-asn1-ber/asn1-ber"
	"github.com/go-ldap/ldap/v3"
)

func HandleModifyRequest(req *ber.Packet, boundDN string, server *Server, conn net.Conn) error {
	if boundDN == "" {
		return ldap.NewError(ldap.LDAPResultInsufficientAccessRights, errors.New("anonymous Write denied"))
	}
	modReq, err := parseModifyRequest(req)
	if err != nil {
		return err
	}

	log.Printf("Parsed Modification %s", dumpModRequest(modReq))

	fnNames := []string{}
	for k := range server.ModifyFns {
		fnNames = append(fnNames, k)
	}
	fn := routeFunc(modReq.DN, fnNames)
	var modifier Modifier
	if modifier = server.ModifyFns[fn]; modifier == nil {
		if fn == "" {
			err = fmt.Errorf("no suitable handler found for dn: '%s'", modReq.DN)
		} else {
			err = fmt.Errorf("handler '%s' does not support modify", fn)
		}
		return ldap.NewError(ldap.LDAPResultUnwillingToPerform, err)
	}
	code, err := modifier.Modify(boundDN, modReq, conn)
	return ldap.NewError(uint16(code), err)
}

func dumpModRequest(mr *ldap.ModifyRequest) string {
	str := fmt.Sprintf("dn: %s\n", mr.DN)
	for _, change := range mr.Changes {
		str += fmt.Sprintf("op: %d\n attr: %s values: %v\n", change.Operation, change.Modification.Type, change.Modification.Vals)
	}
	return str
}

func parseModifyRequest(req *ber.Packet) (*ldap.ModifyRequest, error) {
	modReq := ldap.ModifyRequest{}
	// LDAP Modify requests have 2 Elements (DN and AttributeList)
	if len(req.Children) != 2 {
		return nil, ldap.NewError(ldap.LDAPResultProtocolError, errors.New("invalid modify request"))
	}

	dn, ok := req.Children[0].Value.(string)
	if !ok {
		return nil, ldap.NewError(ldap.LDAPResultProtocolError, errors.New("error decoding entry DN"))
	}

	_, err := ldap.ParseDN(dn)
	if err != nil {
		return nil, ldap.NewError(ldap.LDAPResultProtocolError, err)
	}
	modReq.DN = dn

	ml, err := parseModList(req.Children[1])
	if err != nil {
		return nil, err
	}
	modReq.Changes = ml

	return &modReq, nil
}

func parseModList(req *ber.Packet) ([]ldap.Change, error) {
	var changes []ldap.Change

	if req.ClassType != ber.ClassUniversal || req.TagType != ber.TypeConstructed || req.Tag != ber.TagSequence {
		return nil, ldap.NewError(ldap.LDAPResultProtocolError, errors.New("error decoding Changes List"))
	}

	for _, c := range req.Children {
		var change ldap.Change
		switch c.Children[0].Data.Bytes()[0] {
		default:
			return nil, ldap.NewError(ldap.LDAPResultDecodingError, errors.New("invalid change operation"))
		case ldap.AddAttribute:
			change.Operation = ldap.AddAttribute
			log.Print("op=add")
		case ldap.ReplaceAttribute:
			change.Operation = ldap.ReplaceAttribute
			log.Print("op=replace")
		case ldap.DeleteAttribute:
			change.Operation = ldap.DeleteAttribute
			log.Print("op=delete")
		}
		attr, err := parseAttribute(c.Children[1], true)
		if err != nil {
			return nil, err
		}
		change.Modification = ldap.PartialAttribute{Type: attr.Type, Vals: attr.Vals}
		changes = append(changes, change)

	}
	return changes, nil
}
