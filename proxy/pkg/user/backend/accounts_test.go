package backend

import (
	"context"
	"testing"

	userv1beta1 "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	"github.com/cs3org/reva/pkg/token/manager/jwt"
	"github.com/owncloud/ocis/ocis-pkg/log"
	"github.com/owncloud/ocis/ocis-pkg/oidc"
	accountsmsg "github.com/owncloud/ocis/protogen/gen/ocis/messages/accounts/v0"
	settingsmsg "github.com/owncloud/ocis/protogen/gen/ocis/messages/settings/v0"
	accountssvc "github.com/owncloud/ocis/protogen/gen/ocis/services/accounts/v0"
	settingssvc "github.com/owncloud/ocis/protogen/gen/ocis/services/settings/v0"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"go-micro.dev/v4/client"
)

var mockAccResp = []*accountsmsg.Account{
	{
		Id:                       "1234",
		AccountEnabled:           true,
		DisplayName:              "foo",
		PreferredName:            "prefname",
		UidNumber:                1,
		GidNumber:                2,
		Mail:                     "foo@example.org",
		OnPremisesSamAccountName: "samaccount",
		MemberOf: []*accountsmsg.Group{
			{OnPremisesSamAccountName: "g1"},
			{OnPremisesSamAccountName: "g2"},
		},
	},
}

var expectedRoles = []*settingsmsg.UserRoleAssignment{
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
			u, _, err := accBackend.GetUserByClaims(context.Background(), tests[k].claim, tests[k].value, true)

			assert.NoError(t, err)
			assert.NotNil(t, u)
			assertUserMatchesAccount(t, mockAccResp[0], u)

		})
	}
}

func TestGetUserByClaimsNotFound(t *testing.T) {
	accBackend := newAccountsBackend([]*accountsmsg.Account{}, expectedRoles)

	u, _, err := accBackend.GetUserByClaims(context.Background(), "mail", "foo@example.com", true)

	assert.Error(t, err)
	assert.Nil(t, u)
	assert.Equal(t, ErrAccountNotFound, err)
}

func TestGetUserByClaimsInvalidClaim(t *testing.T) {
	accBackend := newAccountsBackend([]*accountsmsg.Account{}, expectedRoles)
	u, _, err := accBackend.GetUserByClaims(context.Background(), "invalidClaimName", "efwfwfwfe", true)

	assert.Nil(t, u)
	assert.Error(t, err)
}

func TestGetUserByClaimsDisabledAccount(t *testing.T) {
	accBackend := newAccountsBackend([]*accountsmsg.Account{{AccountEnabled: false}}, expectedRoles)
	u, _, err := accBackend.GetUserByClaims(context.Background(), "mail", "foo@example.com", true)

	assert.Nil(t, u)
	assert.Error(t, err)
	assert.Equal(t, ErrAccountDisabled, err)
}

func TestAuthenticate(t *testing.T) {
	accBackend := newAccountsBackend(mockAccResp, expectedRoles)
	u, _, err := accBackend.Authenticate(context.Background(), "foo", "secret")

	assert.NoError(t, err)
	assert.NotNil(t, u)
	assertUserMatchesAccount(t, mockAccResp[0], u)
}

func TestAuthenticateFailed(t *testing.T) {
	accBackend := newAccountsBackend([]*accountsmsg.Account{}, expectedRoles)
	u, _, err := accBackend.Authenticate(context.Background(), "foo", "secret")

	assert.Nil(t, u)
	assert.Error(t, err)
}

func TestCreateUserFromClaims(t *testing.T) {
	exp := mockAccResp[0]
	accBackend := newAccountsBackend([]*accountsmsg.Account{}, expectedRoles)
	act, _ := accBackend.CreateUserFromClaims(context.Background(), map[string]interface{}{
		oidc.Name:              mockAccResp[0].DisplayName,
		oidc.PreferredUsername: mockAccResp[0].OnPremisesSamAccountName,
		oidc.Email:             mockAccResp[0].Mail,
		oidc.UIDNumber:         "1",
		oidc.GIDNumber:         "2",
		oidc.Groups:            []string{"g1", "g2"},
	})

	assert.NotNil(t, act.Id)
	assert.Equal(t, exp.Id, act.Id.OpaqueId)
	assert.Equal(t, exp.Mail, act.Mail)
	assert.Equal(t, exp.DisplayName, act.DisplayName)
	assert.Equal(t, exp.OnPremisesSamAccountName, act.Username)
}

func TestGetUserGroupsUnimplemented(t *testing.T) {
	accBackend := newAccountsBackend([]*accountsmsg.Account{}, expectedRoles)
	assert.Panics(t, func() { accBackend.GetUserGroups(context.Background(), "foo") })
}

func assertUserMatchesAccount(t *testing.T, exp *accountsmsg.Account, act *userv1beta1.User) {
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
	assert.Equal(t, int64(1), act.UidNumber)
	assert.Equal(t, int64(2), act.GidNumber)
}

func newAccountsBackend(mockAccounts []*accountsmsg.Account, mockRoles []*settingsmsg.UserRoleAssignment) UserBackend {
	accSvc, roleSvc := getAccountService(mockAccounts, nil), getRoleService(mockRoles, nil)
	tokenManager, _ := jwt.New(map[string]interface{}{
		"secret":  "change-me",
		"expires": int64(24 * 60 * 60),
	})
	accBackend := NewAccountsServiceUserBackend(accSvc, roleSvc, "https://idp.example.org", tokenManager, log.NewLogger())
	zerolog.SetGlobalLevel(zerolog.Disabled)
	return accBackend
}

func getAccountService(expectedResponse []*accountsmsg.Account, err error) *accountssvc.MockAccountsService {
	return &accountssvc.MockAccountsService{
		ListFunc: func(ctx context.Context, in *accountssvc.ListAccountsRequest, opts ...client.CallOption) (*accountssvc.ListAccountsResponse, error) {
			return &accountssvc.ListAccountsResponse{Accounts: expectedResponse}, err
		},
		CreateFunc: func(ctx context.Context, in *accountssvc.CreateAccountRequest, opts ...client.CallOption) (*accountsmsg.Account, error) {
			a := in.Account
			a.Id = "1234"
			return a, nil
		},
	}
}

func getRoleService(expectedResponse []*settingsmsg.UserRoleAssignment, err error) *settingssvc.MockRoleService {
	return &settingssvc.MockRoleService{
		ListRoleAssignmentsFunc: func(ctx context.Context, req *settingssvc.ListRoleAssignmentsRequest, opts ...client.CallOption) (*settingssvc.ListRoleAssignmentsResponse, error) {
			return &settingssvc.ListRoleAssignmentsResponse{Assignments: expectedResponse}, err
		},
	}

}
