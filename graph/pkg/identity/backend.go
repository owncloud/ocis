package identity

import (
	"context"
	"net/url"

	cs3 "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	msgraph "github.com/yaegashi/msgraph.go/beta"
)

type Backend interface {
	GetUser(ctx context.Context, nameOrId string) (*msgraph.User, error)
	GetUsers(ctx context.Context, queryParam url.Values) ([]*msgraph.User, error)

	GetGroup(ctx context.Context, nameOrId string) (*msgraph.Group, error)
	GetGroups(ctx context.Context, queryParam url.Values) ([]*msgraph.Group, error)
}

func CreateUserModelFromCS3(u *cs3.User) *msgraph.User {
	if u.Id == nil {
		u.Id = &cs3.UserId{}
	}
	return &msgraph.User{
		DisplayName: &u.DisplayName,
		Mail:        &u.Mail,
		// TODO u.Groups are those ids or group names?
		OnPremisesSamAccountName: &u.Username,
		DirectoryObject: msgraph.DirectoryObject{
			Entity: msgraph.Entity{
				ID: &u.Id.OpaqueId,
				Object: msgraph.Object{
					AdditionalData: map[string]interface{}{
						"uidnumber": u.UidNumber,
						"gidnumber": u.GidNumber,
					},
				},
			},
		},
	}
}
