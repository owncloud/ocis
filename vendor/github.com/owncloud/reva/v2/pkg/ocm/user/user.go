package user

import (
	"fmt"
	"net/url"
	"strings"

	userpb "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
)

// FederatedID creates a federated id for local users by
// 1. stripping the protocol from the domain and
// 2. concatenating the opaque id with the domain to get a unique identifier that cannot collide with other users
func FederatedID(id *userpb.UserId, domain string) *userpb.UserId {
	// strip protocol from the domain
	idp := id.Idp
	if u, err := url.Parse(id.Idp); err == nil && u.Host != "" {
		idp = u.Host
	}
	opaqueId := id.OpaqueId
	if !strings.Contains(id.OpaqueId, "@") {
		opaqueId = id.OpaqueId + "@" + idp
	}

	u := &userpb.UserId{
		Type:     userpb.UserType_USER_TYPE_FEDERATED,
		Idp:      domain,
		OpaqueId: opaqueId,
	}

	return u
}

// DecodeRemoteUserFederatedID decodes opaque id into remote user's federated id by
// splitting the opaque id at the last @ to get the opaque id and the domain
func DecodeRemoteUserFederatedID(id *userpb.UserId) *userpb.UserId {
	remoteId := &userpb.UserId{
		Type:     userpb.UserType_USER_TYPE_PRIMARY,
		Idp:      id.Idp,
		OpaqueId: id.OpaqueId,
	}
	remote := id.OpaqueId
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
	if u.Idp == "" {
		return u.OpaqueId
	}
	// strip protocol from the domain
	idp := u.Idp
	if u, err := url.Parse(u.Idp); err == nil && u.Host != "" {
		idp = u.Host
	}
	return fmt.Sprintf("%s@%s", u.OpaqueId, idp)
}
