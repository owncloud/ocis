package user

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"

	userpb "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
)

// LocalUserFederatedID creates a federated id for local users by
// 1. stripping the protocol from the domain and
func LocalUserFederatedID(id *userpb.UserId, domain string) *userpb.UserId {
	opaqueId := id.OpaqueId + "@" + id.Idp
	return &userpb.UserId{
		Type:     userpb.UserType_USER_TYPE_FEDERATED,
		Idp:      domain,
		OpaqueId: opaqueId,
	}
}

// EncodeRemoteUserFederatedID encodes a federated id for remote users by
// 1. stripping the protocol from the domain and
// 2. base64 encoding the opaque id with the domain to get a unique identifier that cannot collide with other users
func EncodeRemoteUserFederatedID(id *userpb.UserId) *userpb.UserId {
	// strip protocol from the domain
	domain := id.Idp
	if u, err := url.Parse(domain); err == nil && u.Host != "" {
		domain = u.Host
	}
	return &userpb.UserId{
		Type:     userpb.UserType_USER_TYPE_FEDERATED,
		Idp:      domain,
		OpaqueId: base64.URLEncoding.EncodeToString([]byte(id.OpaqueId + "@" + domain)),
	}
}

// DecodeRemoteUserFederatedID decodes opaque id into remote user's federated id by
// 1. decoding the base64 encoded opaque id
// 2. splitting the opaque id at the last @ to get the opaque id and the domain
func DecodeRemoteUserFederatedID(id *userpb.UserId) *userpb.UserId {
	remoteId := &userpb.UserId{
		Type:     userpb.UserType_USER_TYPE_PRIMARY,
		Idp:      id.Idp,
		OpaqueId: id.OpaqueId,
	}
	bytes, err := base64.URLEncoding.DecodeString(id.GetOpaqueId())
	if err != nil {
		return remoteId
	}
	remote := string(bytes)
	last := strings.LastIndex(remote, "@")
	if last == -1 {
		return remoteId
	}
	remoteId.OpaqueId = remote[:last]
	remoteId.Idp = remote[last+1:]

	return remoteId
}

// FormatOCMUser formats a user id in the form of <opaque-id>@<idp> used by the OCM API in shareWith, owner and creator fields
func FormatOCMUser(u *userpb.UserId) string {
	return fmt.Sprintf("%s@%s", u.OpaqueId, u.Idp)
}
