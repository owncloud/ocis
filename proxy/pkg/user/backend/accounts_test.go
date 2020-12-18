package backend

import (
	"context"
	userv1beta1 "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	"github.com/micro/go-micro/v2/client"
	accounts "github.com/owncloud/ocis/accounts/pkg/proto/v0"
	"github.com/owncloud/ocis/ocis-pkg/log"
	"github.com/owncloud/ocis/ocis-pkg/oidc"
	settings "github.com/owncloud/ocis/settings/pkg/proto/v0"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"testing"
)

var mockAccResp = []*accounts.Account{
	{
		Id:                       "1234",
		AccountEnabled:           true,
		DisplayName:              "foo",
		PreferredName:            "prefname",
		UidNumber:                1,
		GidNumber:                2,
		Mail:                     "foo@example.org",
		OnPremisesSamAccountName: "samaccount",
		MemberOf: []*accounts.Group{
			{OnPremisesSamAccountName: "g1"},
			{OnPremisesSamAccountName: "g2"},
		},
	},
}

var expectedRoles = []*settings.UserRoleAssignment{
	{Id: "abc", AccountUuid: "1234", RoleId: "a"},
	{Id: "def", AccountUuid: "1234", RoleId: "b"},
}

func TestGetUserByClaimsFound(t *testing.T) {
	type testCase struct {
		id, claim, value string
	}

	var tests = []testCase{
		{id: "Mail", claim: "mail", value: "foo@example.org"},
		{id: "Username", claim: "username", value: "prefname"},
		{id: "ID", claim: "id", value: "1234"},
	}

	accBackend := newAccountsBackend(mockAccResp, expectedRoles)

	for k := range tests {
		t.Run(tests[k].id, func(t *testing.T) {
			u, err := accBackend.GetUserByClaims(context.Background(), tests[k].claim, tests[k].value, true)

			assert.NoError(t, err)
			assert.NotNil(t, u)
			assertUserMatchesAccount(t, mockAccResp[0], u)

		})
	}
}

func TestGetUserByClaimsNotFound(t *testing.T) {
	accBackend := newAccountsBackend([]*accounts.Account{}, expectedRoles)

	u, err := accBackend.GetUserByClaims(context.Background(), "mail", "foo@example.com", true)

	assert.Error(t, err)
	assert.Nil(t, u)
	assert.Equal(t, ErrAccountNotFound, err)
}

func TestGetUserByClaimsInvalidClaim(t *testing.T) {
	accBackend := newAccountsBackend([]*accounts.Account{}, expectedRoles)
	u, err := accBackend.GetUserByClaims(context.Background(), "invalidClaimName", "efwfwfwfe", true)

	assert.Nil(t, u)
	assert.Error(t, err)
}

func TestGetUserByClaimsDisabledAccount(t *testing.T) {
	accBackend := newAccountsBackend([]*accounts.Account{{AccountEnabled: false}}, expectedRoles)
	u, err := accBackend.GetUserByClaims(context.Background(), "mail", "foo@example.com", true)

	assert.Nil(t, u)
	assert.Error(t, err)
	assert.Equal(t, ErrAccountDisabled, err)
}

func TestAuthenticate(t *testing.T) {
	accBackend := newAccountsBackend(mockAccResp, expectedRoles)
	u, err := accBackend.Authenticate(context.Background(), "foo", "secret")

	assert.NoError(t, err)
	assert.NotNil(t, u)
	assertUserMatchesAccount(t, mockAccResp[0], u)
}

func TestAuthenticateFailed(t *testing.T) {
	accBackend := newAccountsBackend([]*accounts.Account{}, expectedRoles)
	u, err := accBackend.Authenticate(context.Background(), "foo", "secret")

	assert.Nil(t, u)
	assert.Error(t, err)
}

func TestCreateUserFromClaims(t *testing.T) {
	exp := mockAccResp[0]
	accBackend := newAccountsBackend([]*accounts.Account{}, expectedRoles)
	act, _ := accBackend.CreateUserFromClaims(context.Background(), &oidc.StandardClaims{
		DisplayName:       mockAccResp[0].DisplayName,
		PreferredUsername: mockAccResp[0].OnPremisesSamAccountName,
		Email:             mockAccResp[0].Mail,
		UIDNumber:         "1",
		GIDNumber:         "2",
		Groups:            []string{"g1", "g2"},
	})

	assert.NotNil(t, act.Id)
	assert.Equal(t, exp.Id, act.Id.OpaqueId)
	assert.Equal(t, exp.Mail, act.Mail)
	assert.Equal(t, exp.DisplayName, act.DisplayName)
	assert.Equal(t, exp.OnPremisesSamAccountName, act.Username)
}

func TestGetUserGroupsUnimplemented(t *testing.T) {
	accBackend := newAccountsBackend([]*accounts.Account{}, expectedRoles)
	assert.Panics(t, func() { accBackend.GetUserGroups(context.Background(), "foo") })
}

func assertUserMatchesAccount(t *testing.T, exp *accounts.Account, act *userv1beta1.User) {
	// User
	assert.NotNil(t, act.Id)
	assert.Equal(t, exp.Id, act.Id.OpaqueId)
	assert.Equal(t, exp.Mail, act.Mail)
	assert.Equal(t, exp.DisplayName, act.DisplayName)
	assert.Equal(t, exp.OnPremisesSamAccountName, act.Username)

	// Groups
	assert.ElementsMatch(t, []string{"g1", "g2"}, act.Groups)

	// Roles
	assert.NotNil(t, act.Opaque.Map["roles"])
	assert.Equal(t, `["a","b"]`, string(act.Opaque.Map["roles"].GetValue()))

	// UID/GID
	assert.NotNil(t, act.Opaque.Map["uid"])
	assert.Equal(t, "1", string(act.Opaque.Map["uid"].GetValue()))

	assert.NotNil(t, act.Opaque.Map["gid"])
	assert.Equal(t, "2", string(act.Opaque.Map["gid"].GetValue()))
}

func newAccountsBackend(mockAccounts []*accounts.Account, mockRoles []*settings.UserRoleAssignment) UserBackend {
	accSvc, roleSvc := getAccountService(mockAccounts, nil), getRoleService(mockRoles, nil)
	accBackend := NewAccountsServiceUserBackend(accSvc, roleSvc, "https://idp.example.org", log.NewLogger())
	zerolog.SetGlobalLevel(zerolog.Disabled)
	return accBackend
}

func getAccountService(expectedResponse []*accounts.Account, err error) *accounts.MockAccountsService {
	return &accounts.MockAccountsService{
		ListFunc: func(ctx context.Context, in *accounts.ListAccountsRequest, opts ...client.CallOption) (*accounts.ListAccountsResponse, error) {
			return &accounts.ListAccountsResponse{Accounts: expectedResponse}, err
		},
		CreateFunc: func(ctx context.Context, in *accounts.CreateAccountRequest, opts ...client.CallOption) (*accounts.Account, error) {
			a := in.Account
			a.Id = "1234"
			return a, nil
		},
	}
}

func getRoleService(expectedRespone []*settings.UserRoleAssignment, err error) *settings.MockRoleService {
	return &settings.MockRoleService{
		ListRoleAssignmentsFunc: func(ctx context.Context, req *settings.ListRoleAssignmentsRequest, opts ...client.CallOption) (*settings.ListRoleAssignmentsResponse, error) {
			return &settings.ListRoleAssignmentsResponse{Assignments: expectedRespone}, err
		},
	}

}
