// Copyright 2012 The Go Authors. All rights reserved.
// Copyright 2021 The LibreGraph Authors.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ldapserver

import (
	"net"

	ber "github.com/go-asn1-ber/asn1-ber"
	"github.com/go-ldap/ldap/v3"
)

func HandleBindRequest(req *ber.Packet, fns map[string]Binder, conn net.Conn) (resultCode LDAPResultCode) {
	defer func() {
		if r := recover(); r != nil {
			resultCode = ldap.LDAPResultOperationsError
		}
	}()

	// we only support ldapv3
	ldapVersion, ok := req.Children[0].Value.(int64)
	if !ok {
		return ldap.LDAPResultProtocolError
	}
	if ldapVersion != 3 {
		logger.V(1).Info("Unsupported LDAP version", "version", ldapVersion)
		return ldap.LDAPResultInappropriateAuthentication
	}

	// auth types
	bindDN, ok := req.Children[1].Value.(string)
	if !ok {
		return ldap.LDAPResultProtocolError
	}
	bindAuth := req.Children[2]
	switch bindAuth.Tag {
	default:
		logger.V(1).Info("Unknown LDAP authentication method", "tag", bindAuth.Tag)
		return ldap.LDAPResultInappropriateAuthentication
	case LDAPBindAuthSimple:
		if len(req.Children) == 3 {
			fnNames := []string{}
			for k := range fns {
				fnNames = append(fnNames, k)
			}
			fn := routeFunc(bindDN, fnNames)
			resultCode, err := fns[fn].Bind(bindDN, bindAuth.Data.String(), conn)
			if err != nil {
				logger.Error(err, "BindFn Error")
				return ldap.LDAPResultOperationsError
			}
			return resultCode
		} else {
			logger.V(1).Info("Simple bind request has wrong # children.  len(req.Children) != 3")
			return ldap.LDAPResultInappropriateAuthentication
		}
	case LDAPBindAuthSASL:
		logger.V(1).Info("SASL authentication is not supported")
		return ldap.LDAPResultInappropriateAuthentication
	}
	return ldap.LDAPResultOperationsError
}

func encodeBindResponse(messageID int64, ldapResultCode LDAPResultCode) *ber.Packet {
	responsePacket := ber.Encode(ber.ClassUniversal, ber.TypeConstructed, ber.TagSequence, nil, "LDAP Response")
	responsePacket.AppendChild(ber.NewInteger(ber.ClassUniversal, ber.TypePrimitive, ber.TagInteger, messageID, "Message ID"))

	bindReponse := ber.Encode(ber.ClassApplication, ber.TypeConstructed, ldap.ApplicationBindResponse, nil, "Bind Response")
	bindReponse.AppendChild(ber.NewInteger(ber.ClassUniversal, ber.TypePrimitive, ber.TagEnumerated, uint64(ldapResultCode), "resultCode: "))
	bindReponse.AppendChild(ber.NewString(ber.ClassUniversal, ber.TypePrimitive, ber.TagOctetString, "", "matchedDN: "))
	bindReponse.AppendChild(ber.NewString(ber.ClassUniversal, ber.TypePrimitive, ber.TagOctetString, "", "errorMessage: "))

	responsePacket.AppendChild(bindReponse)

	// ber.PrintPacket(responsePacket)
	return responsePacket
}
