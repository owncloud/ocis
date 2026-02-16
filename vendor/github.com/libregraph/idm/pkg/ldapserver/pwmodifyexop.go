package ldapserver

import (
	"errors"
	"fmt"
	"net"

	ber "github.com/go-asn1-ber/asn1-ber"
	"github.com/go-ldap/ldap/v3"
	"github.com/libregraph/idm/pkg/ldappassword"
)

const pwmodOID = "1.3.6.1.4.1.4203.1.11.1"

const (
	TagReqIdentity = 0
	TagReqOldPW    = 1
	TagReqNewPW    = 2
	TagRespGenPW   = 0
)

func init() {
	RegisterExtendedOperation(pwmodOID, HandlePasswordModifyExOp)
}

func HandlePasswordModifyExOp(req *ber.Packet, boundDN string, server *Server, conn net.Conn) (*ber.Packet, error) {
	var passwordGenerated bool
	logger.V(1).Info("HandlePasswordModifyExOp")
	if boundDN == "" {
		return nil, ldap.NewError(ldap.LDAPResultUnwillingToPerform, errors.New("authentication required"))
	}

	pwReq, err := parsePasswordModifyExop(req)
	if err != nil {
		return nil, err
	}

	// If `UserIdentity` is empty, this is a request to update the bound user's own password
	if pwReq.UserIdentity == "" {
		pwReq.UserIdentity = boundDN
	}

	if pwReq.NewPassword == "" {
		// New password empty means, we're requested to generate a new password
		var err error
		if pwReq.NewPassword, err = ldappassword.GenerateRandomPassword(server.GeneratedPasswordLength); err != nil {
			logger.Error(err, "Failed to generate new password")
			return nil, ldap.NewError(ldap.LDAPResultOperationsError, errors.New("Failed to generate new Password"))
		}
		passwordGenerated = true
	}
	pwReq.NewPassword, err = ldappassword.Hash(pwReq.NewPassword, "{ARGON2}")
	if err != nil {
		return nil, ldap.NewError(ldap.LDAPResultOperationsError, err)
	}

	logger.V(1).Info("Modify password extended operation", "dn", pwReq.UserIdentity)

	fnNames := []string{}
	for k := range server.PasswordExOpFns {
		fnNames = append(fnNames, k)
	}
	fn := routeFunc(pwReq.UserIdentity, fnNames)
	var pwUpdatefn PasswordUpdater
	if pwUpdatefn = server.PasswordExOpFns[fn]; pwUpdatefn == nil {
		if fn == "" {
			err = fmt.Errorf("no suitable handler found for dn: '%s'", pwReq.UserIdentity)
		} else {
			err = fmt.Errorf("handler '%s' does not support add", fn)
		}
		return nil, ldap.NewError(ldap.LDAPResultUnwillingToPerform, err)
	}
	code, err := pwUpdatefn.ModifyPasswordExop(boundDN, pwReq, conn)
	if code != ldap.LDAPResultSuccess {
		return nil, ldap.NewError(uint16(code), err)
	}
	var response *ber.Packet
	if passwordGenerated {
		response = ber.NewSequence("PasswdModifyResponseValue")
		response.AppendChild(ber.NewString(ber.ClassContext, ber.TypePrimitive, 0, pwReq.NewPassword, "genPasswd"))
	}
	return response, nil
}

func parsePasswordModifyExop(req *ber.Packet) (*ldap.PasswordModifyRequest, error) {
	pwReq := ldap.PasswordModifyRequest{}

	// An absent (or empty) body of the request is valid. Translates into: "generate a new password for
	// for the current user"
	if req == nil {
		return &pwReq, nil
	}

	inner := ber.DecodePacket(req.Data.Bytes())
	if inner == nil {
		return &pwReq, nil
	}

	if len(inner.Children) > 3 {
		return nil, ldap.NewError(ldap.LDAPResultDecodingError, errors.New("invalid request"))
	}

	for _, kid := range inner.Children {
		if kid.ClassType != ber.ClassContext {
			return nil, ldap.NewError(ldap.LDAPResultDecodingError, errors.New("invalid request"))
		}
		switch kid.Tag {
		default:
			return nil, ldap.NewError(ldap.LDAPResultDecodingError, errors.New("invalid request"))
		case TagReqIdentity:
			pwReq.UserIdentity = kid.Data.String()
		case TagReqOldPW:
			pwReq.OldPassword = kid.Data.String()
		case TagReqNewPW:
			pwReq.NewPassword = kid.Data.String()
		}
	}
	return &pwReq, nil
}
