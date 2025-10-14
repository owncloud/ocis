package user

import (
	"encoding/base64"
	"fmt"
	"strings"

	userpb "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
)

// FederatedID creates a federated user id by
// 1. stripping the protocol from the domain and
// 2. base64 encoding the opaque id with the domain to get a unique identifier that cannot collide with other users
func FederatedID(id *userpb.UserId, domain string) *userpb.UserId {
	opaqueId := base64.URLEncoding.EncodeToString([]byte(id.OpaqueId + "@" + id.Idp))
	return &userpb.UserId{
		Type:     userpb.UserType_USER_TYPE_FEDERATED,
		Idp:      domain,
		OpaqueId: opaqueId,
	}
}

// RemoteID creates a remote user id by
// 1. decoding the base64 encoded opaque id
// 2. splitting the opaque id at the last @ to get the opaque id and the domain
func RemoteID(id *userpb.UserId) *userpb.UserId {
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
