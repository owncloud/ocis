// Copyright 2011 The Go Authors. All rights reserved.
// Copyright 2021 The LibreGraph Authors.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ldapserver

import (
	"errors"
	"net"
	"strings"

	ber "github.com/go-asn1-ber/asn1-ber"
	"github.com/go-ldap/ldap/v3"
)

func HandleSearchRequest(req *ber.Packet, controls *[]ldap.Control, messageID int64, boundDN string, server *Server, conn net.Conn) (doneControls *[]ldap.Control, resultErr error) {
	searchReq, err := parseSearchRequest(boundDN, req, controls)
	if err != nil {
		return nil, ldap.NewError(ldap.LDAPResultOperationsError, err)
	}

	var filterPacket *ber.Packet
	if server.EnforceLDAP {
		filterPacket, err = ldap.CompileFilter(searchReq.Filter)
		if err != nil {
			return nil, ldap.NewError(ldap.LDAPResultOperationsError, err)
		}
	}

	fnNames := []string{}
	for k := range server.SearchFns {
		fnNames = append(fnNames, k)
	}
	fn := routeFunc(searchReq.BaseDN, fnNames)
	searchResp, err := server.SearchFns[fn].Search(boundDN, searchReq, conn)
	if err != nil {
		return &searchResp.Controls, ldap.NewError(uint16(searchResp.ResultCode), err)
	}

	if server.EnforceLDAP {
		if searchReq.DerefAliases != ldap.NeverDerefAliases { // [-a {never|always|search|find}
			// TODO: Server DerefAliases not supported: RFC4511 4.5.1.3
		}
		if searchReq.TimeLimit > 0 {
			// TODO: Server TimeLimit not implemented
		}
	}

	i := 0
	for _, entry := range searchResp.Entries {
		if server.EnforceLDAP {
			// filter
			keep, resultCode := ServerApplyFilter(filterPacket, entry)
			if resultCode != ldap.LDAPResultSuccess {
				return &searchResp.Controls, ldap.NewError(uint16(resultCode), errors.New("ServerApplyFilter error"))
			}
			if !keep {
				continue
			}

			keep, resultCode = ServerFilterScope(searchReq.BaseDN, searchReq.Scope, entry)
			if resultCode != ldap.LDAPResultSuccess {
				return &searchResp.Controls, ldap.NewError(uint16(resultCode), errors.New("ServerApplyScope error"))
			}
			if !keep {
				continue
			}

			resultCode, err = ServerFilterAttributes(searchReq.Attributes, entry)
			if err != nil {
				return &searchResp.Controls, ldap.NewError(uint16(resultCode), err)
			}

			// size limit
			if searchReq.SizeLimit > 0 && i >= searchReq.SizeLimit {
				break
			}
			i++
		}

		// respond
		responsePacket := encodeSearchResponse(messageID, searchReq, entry)
		if err = sendPacket(conn, responsePacket); err != nil {
			return &searchResp.Controls, ldap.NewError(ldap.LDAPResultOperationsError, err)
		}
	}
	return &searchResp.Controls, nil
}

func parseSearchRequest(boundDN string, req *ber.Packet, controls *[]ldap.Control) (*ldap.SearchRequest, error) {
	if len(req.Children) != 8 {
		return &ldap.SearchRequest{}, ldap.NewError(ldap.LDAPResultOperationsError, errors.New("Bad search request"))
	}

	// Parse the request.
	baseObject, ok := req.Children[0].Value.(string)
	if !ok {
		return &ldap.SearchRequest{}, ldap.NewError(ldap.LDAPResultProtocolError, errors.New("Bad search request"))
	}
	s, ok := req.Children[1].Value.(int64)
	if !ok {
		return &ldap.SearchRequest{}, ldap.NewError(ldap.LDAPResultProtocolError, errors.New("Bad search request"))
	}
	scope := int(s)
	d, ok := req.Children[2].Value.(int64)
	if !ok {
		return &ldap.SearchRequest{}, ldap.NewError(ldap.LDAPResultProtocolError, errors.New("Bad search request"))
	}
	derefAliases := int(d)
	s, ok = req.Children[3].Value.(int64)
	if !ok {
		return &ldap.SearchRequest{}, ldap.NewError(ldap.LDAPResultProtocolError, errors.New("Bad search request"))
	}
	sizeLimit := int(s)
	t, ok := req.Children[4].Value.(int64)
	if !ok {
		return &ldap.SearchRequest{}, ldap.NewError(ldap.LDAPResultProtocolError, errors.New("Bad search request"))
	}
	timeLimit := int(t)
	typesOnly := false
	if req.Children[5].Value != nil {
		typesOnly, ok = req.Children[5].Value.(bool)
		if !ok {
			return &ldap.SearchRequest{}, ldap.NewError(ldap.LDAPResultProtocolError, errors.New("Bad search request"))
		}
	}
	filter, err := DecompileFilter(req.Children[6])
	if err != nil {
		return &ldap.SearchRequest{}, err
	}
	attributes := []string{}
	for _, attr := range req.Children[7].Children {
		a, ok := attr.Value.(string)
		if !ok {
			return &ldap.SearchRequest{}, ldap.NewError(ldap.LDAPResultProtocolError, errors.New("Bad search request"))
		}
		attributes = append(attributes, a)
	}
	searchReq := &ldap.SearchRequest{baseObject, scope,
		derefAliases, sizeLimit, timeLimit,
		typesOnly, filter, attributes, *controls}

	return searchReq, nil
}

func filterAttributes(entry *ldap.Entry, attributes []string) (*ldap.Entry, error) {
	// Only return requested attributes.
	newAttributes := []*ldap.EntryAttribute{}

	for _, attr := range entry.Attributes {
		for _, requested := range attributes {
			if requested == "*" || strings.EqualFold(attr.Name, requested) {
				newAttributes = append(newAttributes, attr)
			}
		}
	}
	entry.Attributes = newAttributes

	return entry, nil
}

func encodeSearchResponse(messageID int64, req *ldap.SearchRequest, res *ldap.Entry) *ber.Packet {
	responsePacket := ber.Encode(ber.ClassUniversal, ber.TypeConstructed, ber.TagSequence, nil, "LDAP Response")
	responsePacket.AppendChild(ber.NewInteger(ber.ClassUniversal, ber.TypePrimitive, ber.TagInteger, messageID, "Message ID"))

	searchEntry := ber.Encode(ber.ClassApplication, ber.TypeConstructed, ldap.ApplicationSearchResultEntry, nil, "Search Result Entry")
	searchEntry.AppendChild(ber.NewString(ber.ClassUniversal, ber.TypePrimitive, ber.TagOctetString, res.DN, "Object Name"))

	attrs := ber.Encode(ber.ClassUniversal, ber.TypeConstructed, ber.TagSequence, nil, "Attributes: ")
	for _, attribute := range res.Attributes {
		attrs.AppendChild(encodeSearchAttribute(attribute.Name, attribute.Values))
	}

	searchEntry.AppendChild(attrs)
	responsePacket.AppendChild(searchEntry)

	return responsePacket
}

func encodeSearchAttribute(name string, values []string) *ber.Packet {
	packet := ber.Encode(ber.ClassUniversal, ber.TypeConstructed, ber.TagSequence, nil, "Attribute")
	packet.AppendChild(ber.NewString(ber.ClassUniversal, ber.TypePrimitive, ber.TagOctetString, name, "Attribute Name"))

	valuesPacket := ber.Encode(ber.ClassUniversal, ber.TypeConstructed, ber.TagSet, nil, "Attribute Values: ")
	for _, value := range values {
		valuesPacket.AppendChild(ber.NewString(ber.ClassUniversal, ber.TypePrimitive, ber.TagOctetString, value, "Attribute Value"))
	}

	packet.AppendChild(valuesPacket)

	return packet
}

func encodeSearchDone(messageID int64, ldapResultCode LDAPResultCode, doneControls *[]ldap.Control) *ber.Packet {
	responsePacket := ber.Encode(ber.ClassUniversal, ber.TypeConstructed, ber.TagSequence, nil, "LDAP Response")
	responsePacket.AppendChild(ber.NewInteger(ber.ClassUniversal, ber.TypePrimitive, ber.TagInteger, messageID, "Message ID"))
	donePacket := ber.Encode(ber.ClassApplication, ber.TypeConstructed, ldap.ApplicationSearchResultDone, nil, "Search result done")
	donePacket.AppendChild(ber.NewInteger(ber.ClassUniversal, ber.TypePrimitive, ber.TagEnumerated, uint64(ldapResultCode), "resultCode: "))
	donePacket.AppendChild(ber.NewString(ber.ClassUniversal, ber.TypePrimitive, ber.TagOctetString, "", "matchedDN: "))
	donePacket.AppendChild(ber.NewString(ber.ClassUniversal, ber.TypePrimitive, ber.TagOctetString, "", "errorMessage: "))
	responsePacket.AppendChild(donePacket)

	if doneControls != nil {
		contextPacket := ber.Encode(ber.ClassContext, ber.TypeConstructed, ber.TagEOC, nil, "Controls: ")
		for _, control := range *doneControls {
			contextPacket.AppendChild(control.Encode())
		}
		responsePacket.AppendChild(contextPacket)
	}
	// ber.PrintPacket(responsePacket)
	return responsePacket
}
