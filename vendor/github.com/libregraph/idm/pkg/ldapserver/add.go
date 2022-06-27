package ldapserver

import (
	"errors"
	"fmt"
	"net"

	ber "github.com/go-asn1-ber/asn1-ber"
	"github.com/go-ldap/ldap/v3"
)

func HandleAddRequest(req *ber.Packet, boundDN string, server *Server, conn net.Conn) error {
	if boundDN == "" {
		return ldap.NewError(ldap.LDAPResultInsufficientAccessRights, errors.New("anonymous Write denied"))
	}
	addReq, err := parseAddRequest(req)
	if err != nil {
		return err
	}
	fnNames := []string{}
	for k := range server.AddFns {
		fnNames = append(fnNames, k)
	}
	fn := routeFunc(addReq.DN, fnNames)
	var adder Adder
	if adder = server.AddFns[fn]; adder == nil {
		if fn == "" {
			err = fmt.Errorf("no suitable handler found for dn: '%s'", addReq.DN)
		} else {
			err = fmt.Errorf("handler '%s' does not support add", fn)
		}
		return ldap.NewError(ldap.LDAPResultUnwillingToPerform, err)
	}
	code, err := adder.Add(boundDN, addReq, conn)
	return ldap.NewError(uint16(code), err)
}

func parseAddRequest(req *ber.Packet) (*ldap.AddRequest, error) {
	addReq := ldap.AddRequest{}
	// LDAP Add request have 2 Elements (DN and AttributeList)
	if len(req.Children) != 2 {
		return nil, ldap.NewError(ldap.LDAPResultProtocolError, errors.New("invalid add request"))
	}

	dn, ok := req.Children[0].Value.(string)
	if !ok {
		return nil, ldap.NewError(ldap.LDAPResultProtocolError, errors.New("error decoding entry DN"))
	}

	_, err := ldap.ParseDN(dn)
	if err != nil {
		return nil, ldap.NewError(ldap.LDAPResultProtocolError, err)
	}
	addReq.DN = dn

	al, err := parseAttributeList(req.Children[1])
	if err != nil {
		return nil, err
	}
	addReq.Attributes = al

	return &addReq, nil
}

func parseAttributeList(req *ber.Packet) ([]ldap.Attribute, error) {
	ldapAttrs := []ldap.Attribute{}

	if req.ClassType != ber.ClassUniversal || req.TagType != ber.TypeConstructed || req.Tag != ber.TagSequence {
		return nil, ldap.NewError(ldap.LDAPResultProtocolError, errors.New("error decoding Attribute List"))
	}

	for _, a := range req.Children {
		attr, err := parseAttribute(a, false)
		if err != nil {
			return nil, err
		}
		ldapAttrs = append(ldapAttrs, *attr)

	}
	return ldapAttrs, nil
}

func parseAttribute(attr *ber.Packet, partial bool) (*ldap.Attribute, error) {
	var la ldap.Attribute
	var ok bool
	var err error

	// Partial attributes, might just contain a type and allow the Value to be absent
	if partial && (len(attr.Children) < 1 || len(attr.Children) > 2) {
		return nil, ldap.NewError(ldap.LDAPResultProtocolError, errors.New("error decoding partial Attribute"))
	} else if !partial && len(attr.Children) != 2 {
		return nil, ldap.NewError(ldap.LDAPResultProtocolError, errors.New("error decoding Attribute"))
	}

	ad := attr.Children[0]
	if ad.ClassType != ber.ClassUniversal || ad.TagType != ber.TypePrimitive || ad.Tag != ber.TagOctetString {
		return nil, ldap.NewError(ldap.LDAPResultProtocolError, errors.New("error decoding Attribute Description"))
	}
	la.Type, ok = ad.Value.(string)
	if !ok {
		return nil, ldap.NewError(ldap.LDAPResultProtocolError, errors.New("error decoding Attribute Description"))
	}

	// We can return here if this is a Partial Attribute without values
	if partial && len(attr.Children) == 1 {
		return &la, nil
	}
	if la.Vals, err = parseAttributeValues(attr.Children[1]); err != nil {
		return nil, err
	}
	return &la, nil
}

func parseAttributeValues(values *ber.Packet) ([]string, error) {
	var strVals []string
	if values.ClassType != ber.ClassUniversal || values.TagType != ber.TypeConstructed || values.Tag != ber.TagSet {
		return nil, ldap.NewError(ldap.LDAPResultProtocolError, errors.New("error decoding Attribute Values"))
	}
	for _, value := range values.Children {
		strVal, ok := value.Value.(string)
		if !ok {
			return nil, ldap.NewError(ldap.LDAPResultProtocolError, errors.New("error decoding Attribute Value"))
		}
		strVals = append(strVals, strVal)
	}
	return strVals, nil
}
