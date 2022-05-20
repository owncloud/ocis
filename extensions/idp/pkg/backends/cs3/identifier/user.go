package cs3

import (
	cs3user "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	konnect "github.com/libregraph/lico"
)

type cs3User struct {
	u *cs3user.User
}

func newCS3User(u *cs3user.User) (*cs3User, error) {
	return &cs3User{
		u: u,
	}, nil
}

func (u *cs3User) Subject() string {
	return u.u.GetId().GetOpaqueId()
}

func (u *cs3User) Email() string {
	return u.u.GetMail()
}

func (u *cs3User) EmailVerified() bool {
	return u.u.GetMailVerified()
}

func (u *cs3User) Name() string {
	return u.u.GetDisplayName()
}

func (u *cs3User) FamilyName() string {
	return ""
}

func (u *cs3User) GivenName() string {
	return ""
}

func (u *cs3User) Username() string {
	return u.u.GetUsername()
}

func (u *cs3User) UniqueID() string {
	return u.u.GetId().GetOpaqueId()
}

func (u *cs3User) BackendClaims() map[string]interface{} {
	claims := make(map[string]interface{})
	claims[konnect.IdentifiedUserIDClaim] = u.u.GetId().GetOpaqueId()

	return claims
}

func (u *cs3User) BackendScopes() []string {
	return nil
}

func (u *cs3User) RequiredScopes() []string {
	return nil
}
