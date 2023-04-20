package ldapserver

import (
	"errors"
	"net"

	ber "github.com/go-asn1-ber/asn1-ber"
	"github.com/go-ldap/ldap/v3"
)

type ExopHandler func(req *ber.Packet, boundDN string, server *Server, conn net.Conn) (*ber.Packet, error)

type ExtendedRequest struct {
	OID  string
	Body *ber.Packet
}

var exopRegistry = map[string]ExopHandler{}

func RegisterExtendedOperation(oid string, handler ExopHandler) {
	exopRegistry[oid] = handler
}

func HandleExtendedRequest(req *ber.Packet, boundDN string, server *Server, conn net.Conn) (*ber.Packet, error) {
	extReq, err := parseExtendedRequest(req)
	if err != nil {
		logger.V(1).Info("parsing extened request failed", "error", err)
		return nil, err
	}
	if handler, ok := exopRegistry[extReq.OID]; ok {
		innerBer, err := handler(extReq.Body, boundDN, server, conn)
		var resCode LDAPResultCode = ldap.LDAPResultSuccess
		msg := ""
		if err != nil {
			if lerr, ok := err.(*ldap.Error); ok {
				msg = lerr.Err.Error()
				resCode = LDAPResultCode(lerr.ResultCode)
			} else {
				msg = err.Error()
				resCode = ldap.LDAPResultOther
			}
		}
		return encodeExtendedResponse(resCode, msg, extReq.OID, innerBer), nil
	} else {
		return nil, ldap.NewError(ldap.LDAPResultProtocolError, errors.New("Unsupported Extented Operations"))
	}
	return nil, err
}

func parseExtendedRequest(req *ber.Packet) (*ExtendedRequest, error) {
	// RFC 4511 Extended Operation:
	// ExtendedRequest ::= [APPLICATION 23] SEQUENCE {
	//         requestName      [0] LDAPOID,
	//         requestValue     [1] OCTET STRING OPTIONAL }
	extReq := ExtendedRequest{}
	if len(req.Children) < 1 || len(req.Children) > 2 {
		return nil, ldap.NewError(ldap.LDAPResultDecodingError, errors.New("invalid extented request"))
	}

	if req.Children[0].Identifier.ClassType != ber.ClassContext || req.Children[0].Identifier.Tag != 0 {
		return nil, ldap.NewError(ldap.LDAPResultDecodingError, errors.New("error decoding extented request OID"))
	}

	extReq.OID = req.Children[0].Data.String()
	if len(req.Children) == 2 {
		if req.Children[1].Identifier.ClassType != ber.ClassContext || req.Children[1].Identifier.Tag != 1 {
			return nil, ldap.NewError(ldap.LDAPResultDecodingError, errors.New("error decoding extented request OID"))
		}
		extReq.Body = req.Children[1]
	}

	logger.V(1).Info("Extended Request", "oid", extReq.OID)
	return &extReq, nil
}

func encodeExtendedResponse(rescode LDAPResultCode, msg, oid string, responseValue *ber.Packet) *ber.Packet {
	respBer := ber.Encode(ber.ClassApplication, ber.TypeConstructed,
		ber.Tag(ldap.ApplicationExtendedResponse), nil,
		ldap.ApplicationMap[ldap.ApplicationExtendedResponse])
	respBer.AppendChild(ber.NewInteger(ber.ClassUniversal, ber.TypePrimitive, ber.TagEnumerated, uint64(rescode), "resultCode: "))
	respBer.AppendChild(ber.NewString(ber.ClassUniversal, ber.TypePrimitive, ber.TagOctetString, "", "matchedDN: "))
	respBer.AppendChild(ber.NewString(ber.ClassUniversal, ber.TypePrimitive, ber.TagOctetString, msg, "errorMessage: "))
	respBer.AppendChild(ber.NewString(ber.ClassContext, ber.TypePrimitive, 10, oid, "responseName"))
	if responseValue != nil {
		encValue := ber.Encode(ber.ClassContext, ber.TypePrimitive, 11, nil, "responseValue")
		encValue.AppendChild(responseValue)
		respBer.AppendChild(encValue)
	}
	return respBer
}
