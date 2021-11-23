package proto_test

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"

	mgrpcc "github.com/asim/go-micro/plugins/client/grpc/v4"
	empty "github.com/golang/protobuf/ptypes/empty"
	"github.com/owncloud/ocis/accounts/pkg/config"
	"github.com/owncloud/ocis/accounts/pkg/proto/v0"
	svc "github.com/owncloud/ocis/accounts/pkg/service/v0"
	oclog "github.com/owncloud/ocis/ocis-pkg/log"
	"github.com/owncloud/ocis/ocis-pkg/service/grpc"
	settings "github.com/owncloud/ocis/settings/pkg/proto/v0"
	"github.com/stretchr/testify/assert"
	"go-micro.dev/v4/client"
	merrors "go-micro.dev/v4/errors"
	"google.golang.org/genproto/protobuf/field_mask"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var service = grpc.NewService()

var dataPath = createTmpDir()

var newCreatedAccounts = []string{}
var newCreatedGroups = []string{}

var mockedRoleAssignment = map[string]string{}

var (
	user1 = proto.Account{
		Id:                       "f9149a32-2b8e-4f04-9e8d-937d81712b9a",
		AccountEnabled:           true,
		IsResourceAccount:        true,
		CreationType:             "",
		DisplayName:              "User One",
		PreferredName:            "user1",
		OnPremisesSamAccountName: "user1",
		UidNumber:                20009,
		GidNumber:                30000,
		Mail:                     "user1@example.com",
		Identities:               []*proto.Identities{nil},
		PasswordProfile:          &proto.PasswordProfile{Password: "heysdjfsdlk"},
		MemberOf: []*proto.Group{
			{Id: "509a9dcd-bb37-4f4f-a01a-19dca27d9cfa"}, // users
		},
	}
	user2 = proto.Account{
		Id:                       "e9149a32-2b8e-4f04-9e8d-937d81712b9a",
		AccountEnabled:           true,
		IsResourceAccount:        true,
		CreationType:             "",
		DisplayName:              "User Two",
		PreferredName:            "user2",
		OnPremisesSamAccountName: "user2",
		UidNumber:                20010,
		GidNumber:                30000,
		Mail:                     "user2@example.com",
		Identities:               []*proto.Identities{nil},
		PasswordProfile:          &proto.PasswordProfile{Password: "hello123"},
		MemberOf: []*proto.Group{
			{Id: "509a9dcd-bb37-4f4f-a01a-19dca27d9cfa"}, // users
		},
	}
)

func init() {
	service = grpc.NewService(
		grpc.Namespace("com.owncloud.api"),
		grpc.Name("accounts"),
		grpc.Address("localhost:9180"),
	)

	cfg := config.New()
	cfg.Repo.Backend = "disk"
	cfg.Repo.Disk.Path = dataPath
	cfg.Server.DemoUsersAndGroups = true
	var hdlr *svc.Service
	var err error

	if hdlr, err = svc.New(svc.Logger(oclog.LoggerFromConfig("accounts", *cfg.Log)), svc.Config(cfg), svc.RoleService(buildRoleServiceMock())); err != nil {
		log.Fatalf("Could not create new service")
	}

	err = proto.RegisterAccountsServiceHandler(service.Server(), hdlr)
	if err != nil {
		log.Fatal("could not register the Accounts handler")
	}
	err = proto.RegisterGroupsServiceHandler(service.Server(), hdlr)
	if err != nil {
		log.Fatal("could not register the Groups handler")
	}

	err = service.Server().Start()
	if err != nil {
		log.Fatal(err)
	}
}

func getAccount(user string) *proto.Account {
	switch user {
	case "user1":
		return &user1
	case "user2":
		return &user2
	default:
		return &proto.Account{
			Id:                fmt.Sprintf("new-id-%s", user),
			AccountEnabled:    true,
			IsResourceAccount: true,
			CreationType:      "",
			DisplayName:       "Regular User",
			PreferredName:     user,
			UidNumber:         2,
			Mail:              fmt.Sprintf("%s@example.com", user),
			Identities:        []*proto.Identities{nil},
		}
	}
}

func getGroup(group string) *proto.Group {
	switch group {
	case "sysusers":
		return &proto.Group{Id: "34f38767-c937-4eb6-b847-1c175829a2a0", GidNumber: 15000, OnPremisesSamAccountName: "sysusers", DisplayName: "Technical users", Description: "A group for technical users. They should not show up in sharing dialogs.", Members: []*proto.Account{
			{Id: "820ba2a1-3f54-4538-80a4-2d73007e30bf"}, // idp
			{Id: "bc596f3c-c955-4328-80a0-60d018b4ad57"}, // reva
		}}
	case "users":
		return &proto.Group{Id: "509a9dcd-bb37-4f4f-a01a-19dca27d9cfa", GidNumber: 30000, OnPremisesSamAccountName: "users", DisplayName: "Users", Description: "A group every normal user belongs to.", Members: []*proto.Account{
			{Id: "4c510ada-c86b-4815-8820-42cdf82c3d51"}, // einstein
			{Id: "f7fbf8c8-139b-4376-b307-cf0a8c2d0d9c"}, // marie
			{Id: "932b4540-8d16-481e-8ef4-588e4b6b151c"}, // feynman
		}}
	case "sailing-lovers":
		return &proto.Group{Id: "6040aa17-9c64-4fef-9bd0-77234d71bad0", GidNumber: 30001, OnPremisesSamAccountName: "sailing-lovers", DisplayName: "Sailing lovers", Members: []*proto.Account{
			{Id: "4c510ada-c86b-4815-8820-42cdf82c3d51"}, // einstein
		}}
	case "violin-haters":
		return &proto.Group{Id: "dd58e5ec-842e-498b-8800-61f2ec6f911f", GidNumber: 30002, OnPremisesSamAccountName: "violin-haters", DisplayName: "Violin haters", Members: []*proto.Account{
			{Id: "4c510ada-c86b-4815-8820-42cdf82c3d51"}, // einstein
		}}
	case "radium-lovers":
		return &proto.Group{Id: "7b87fd49-286e-4a5f-bafd-c535d5dd997a", GidNumber: 30003, OnPremisesSamAccountName: "radium-lovers", DisplayName: "Radium lovers", Members: []*proto.Account{
			{Id: "f7fbf8c8-139b-4376-b307-cf0a8c2d0d9c"}, // marie
		}}
	case "polonium-lovers":
		return &proto.Group{Id: "cedc21aa-4072-4614-8676-fa9165f598ff", GidNumber: 30004, OnPremisesSamAccountName: "polonium-lovers", DisplayName: "Polonium lovers", Members: []*proto.Account{
			{Id: "f7fbf8c8-139b-4376-b307-cf0a8c2d0d9c"}, // marie
		}}
	case "quantum-lovers":
		return &proto.Group{Id: "a1726108-01f8-4c30-88df-2b1a9d1cba1a", GidNumber: 30005, OnPremisesSamAccountName: "quantum-lovers", DisplayName: "Quantum lovers", Members: []*proto.Account{
			{Id: "932b4540-8d16-481e-8ef4-588e4b6b151c"}, // feynman
		}}
	case "philosophy-haters":
		return &proto.Group{Id: "167cbee2-0518-455a-bfb2-031fe0621e5d", GidNumber: 30006, OnPremisesSamAccountName: "philosophy-haters", DisplayName: "Philosophy haters", Members: []*proto.Account{
			{Id: "932b4540-8d16-481e-8ef4-588e4b6b151c"}, // feynman
		}}
	case "physics-lovers":
		return &proto.Group{Id: "262982c1-2362-4afa-bfdf-8cbfef64a06e", GidNumber: 30007, OnPremisesSamAccountName: "physics-lovers", DisplayName: "Physics lovers", Members: []*proto.Account{
			{Id: "4c510ada-c86b-4815-8820-42cdf82c3d51"}, // einstein
			{Id: "f7fbf8c8-139b-4376-b307-cf0a8c2d0d9c"}, // marie
			{Id: "932b4540-8d16-481e-8ef4-588e4b6b151c"}, // feynman
		}}
	}
	return nil
}

func getTestGroups(group string) *proto.Group {
	switch group {
	case "grp1":
		return &proto.Group{Id: "0779f828-4df8-41d6-82dd-40020d9fd0ef", GidNumber: 40000, OnPremisesSamAccountName: "grp1", DisplayName: "Group One", Description: "One group", Members: []*proto.Account{
			{Id: "4c510ada-c86b-4815-8820-42cdf82c3d51"}, // einstein
			{Id: "f7fbf8c8-139b-4376-b307-cf0a8c2d0d9c"}, // marie
		}}
	case "grp2":
		return &proto.Group{Id: "f446fb0a-5e57-4812-9f47-9c1ebca99b5a", GidNumber: 40001, OnPremisesSamAccountName: "grp2", DisplayName: "Group Two", Description: "Two Group", Members: []*proto.Account{
			{Id: "f7fbf8c8-139b-4376-b307-cf0a8c2d0d9c"}, // marie
			{Id: "932b4540-8d16-481e-8ef4-588e4b6b151c"}, // feynman
		}}
	case "grp3":
		return &proto.Group{Id: "cae8b5d5-d133-4d95-8595-fe33a5051017", GidNumber: 40002, OnPremisesSamAccountName: "grp3", DisplayName: "Group Three", Description: "Three Group", Members: []*proto.Account{
			{Id: "4c510ada-c86b-4815-8820-42cdf82c3d51"}, // einstein
			{Id: "f7fbf8c8-139b-4376-b307-cf0a8c2d0d9c"}, // marie
			{Id: "932b4540-8d16-481e-8ef4-588e4b6b151c"}, // feynman
		}}
	}
	return nil
}

func buildRoleServiceMock() settings.RoleService {
	return settings.MockRoleService{
		AssignRoleToUserFunc: func(ctx context.Context, req *settings.AssignRoleToUserRequest, opts ...client.CallOption) (res *settings.AssignRoleToUserResponse, err error) {
			mockedRoleAssignment[req.AccountUuid] = req.RoleId
			return &settings.AssignRoleToUserResponse{
				Assignment: &settings.UserRoleAssignment{
					AccountUuid: req.AccountUuid,
					RoleId:      req.RoleId,
				},
			}, nil
		},
	}
}

func cleanUp(t *testing.T) {
	datastore := filepath.Join(dataPath, "accounts")

	for _, id := range newCreatedAccounts {
		path := filepath.Join(datastore, id)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			continue
		}
		_, _ = deleteAccount(t, id)
	}

	datastore = filepath.Join(dataPath, "groups")

	for _, id := range newCreatedGroups {
		path := filepath.Join(datastore, id)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			continue
		}
		_, _ = deleteGroup(t, id)
	}

	newCreatedAccounts = []string{}
	newCreatedGroups = []string{}
}

func assertUserExists(t *testing.T, account *proto.Account) {
	resp, _ := listAccounts(t)
	assertResponseContainsUser(t, resp, account)
}

func assertResponseContainsUser(t *testing.T, response *proto.ListAccountsResponse, account *proto.Account) {
	var result *proto.Account
	for _, a := range response.Accounts {
		if a.Id == account.Id {
			result = account
		}
	}
	if result == nil {
		t.Fatalf("Account not found in response: %s", account.PreferredName)
	}
	assertAccountsSame(t, account, result)
}

func assertResponseNotContainsUser(t *testing.T, response *proto.ListAccountsResponse, account *proto.Account) {
	for _, a := range response.Accounts {
		if a.Id == account.Id || a.PreferredName == account.PreferredName {
			t.Fatal("Account was not expected to be present in the response, but was found.")
		}
	}
}

func assertAccountsSame(t *testing.T, acc1, acc2 *proto.Account) {
	assert.Equal(t, acc1.Id, acc2.Id)
	assert.Equal(t, acc1.AccountEnabled, acc2.AccountEnabled)
	assert.Equal(t, acc1.IsResourceAccount, acc2.IsResourceAccount)
	assert.Equal(t, acc1.CreationType, acc2.CreationType)
	assert.Equal(t, acc1.DisplayName, acc2.DisplayName)
	assert.Equal(t, acc1.PreferredName, acc2.PreferredName)
	assert.Equal(t, acc1.UidNumber, acc2.UidNumber)
	assert.Equal(t, acc1.OnPremisesSamAccountName, acc2.OnPremisesSamAccountName)
	assert.Equal(t, acc1.GidNumber, acc2.GidNumber)
	assert.Equal(t, acc1.Mail, acc2.Mail)
}

func assertResponseContainsGroup(t *testing.T, response *proto.ListGroupsResponse, group *proto.Group) {
	var result *proto.Group
	for _, g := range response.Groups {
		if g.Id == group.Id {
			result = g
		}
	}
	if result == nil {
		t.Fatalf("Group not found in response: %s", group.GetDisplayName())
	}
	assertGroupsSame(t, group, result)
}

func assertResponseNotContainsGroup(t *testing.T, response *proto.ListGroupsResponse, group *proto.Group) {
	for _, g := range response.Groups {
		if g.Id == group.Id && g.DisplayName == group.DisplayName {
			t.Fatal("Group was not expected to be present in the response, but was found.")
		}
	}
}

func assertGroupsSame(t *testing.T, grp1, grp2 *proto.Group) {
	assert.Equal(t, grp1.Id, grp2.Id)
	assert.Equal(t, grp1.GidNumber, grp2.GidNumber)
	assert.Equal(t, grp1.OnPremisesSamAccountName, grp2.OnPremisesSamAccountName)
	assert.Equal(t, grp1.DisplayName, grp2.DisplayName)

	assert.Equal(t, len(grp1.Members), len(grp2.Members))
	for _, m := range grp1.Members {
		assertGroupHasMember(t, grp2, m.Id)
	}
}

func assertGroupHasMember(t *testing.T, grp *proto.Group, memberID string) {
	for _, m := range grp.Members {
		if m.Id == memberID {
			return
		}
	}

	t.Fatalf("Member with id %s expected to be in group '%s', but not found", memberID, grp.DisplayName)
}

func createAccount(t *testing.T, user string) (*proto.Account, error) {
	client := mgrpcc.NewClient()
	cl := proto.NewAccountsService("com.owncloud.api.accounts", client)

	account := getAccount(user)
	request := &proto.CreateAccountRequest{Account: account}
	res, err := cl.CreateAccount(context.Background(), request)
	if err == nil {
		newCreatedAccounts = append(newCreatedAccounts, account.Id)
	}
	return res, err
}

func createGroup(t *testing.T, group *proto.Group) (*proto.Group, error) {
	client := mgrpcc.NewClient()
	cl := proto.NewGroupsService("com.owncloud.api.accounts", client)

	request := &proto.CreateGroupRequest{Group: group}
	res, err := cl.CreateGroup(context.Background(), request)
	if err == nil {
		newCreatedGroups = append(newCreatedGroups, group.Id)
	}
	return res, err
}

func updateAccount(t *testing.T, account *proto.Account, updateArray []string) (*proto.Account, error) {
	client := mgrpcc.NewClient()
	cl := proto.NewAccountsService("com.owncloud.api.accounts", client)

	updateMask := &field_mask.FieldMask{
		Paths: updateArray,
	}
	request := &proto.UpdateAccountRequest{Account: account, UpdateMask: updateMask}
	res, err := cl.UpdateAccount(context.Background(), request)

	return res, err
}

func listAccounts(t *testing.T) (*proto.ListAccountsResponse, error) {
	request := &proto.ListAccountsRequest{}
	client := mgrpcc.NewClient()
	cl := proto.NewAccountsService("com.owncloud.api.accounts", client)

	response, err := cl.ListAccounts(context.Background(), request)
	return response, err
}

func listGroups(t *testing.T) *proto.ListGroupsResponse {
	request := &proto.ListGroupsRequest{}
	client := mgrpcc.NewClient()
	cl := proto.NewGroupsService("com.owncloud.api.accounts", client)

	response, err := cl.ListGroups(context.Background(), request)
	assert.NoError(t, err)
	return response
}

func deleteAccount(t *testing.T, id string) (*empty.Empty, error) {
	client := mgrpcc.NewClient()
	cl := proto.NewAccountsService("com.owncloud.api.accounts", client)

	req := &proto.DeleteAccountRequest{Id: id}
	res, err := cl.DeleteAccount(context.Background(), req)
	return res, err
}

func deleteGroup(t *testing.T, id string) (*empty.Empty, error) {
	client := mgrpcc.NewClient()
	cl := proto.NewGroupsService("com.owncloud.api.accounts", client)

	req := &proto.DeleteGroupRequest{Id: id}
	res, err := cl.DeleteGroup(context.Background(), req)
	return res, err
}

// createTmpDir creates a temporary dir for tests data.
func createTmpDir() string {
	name, err := ioutil.TempDir("/tmp", "ocis-accounts-store-")
	if err != nil {
		panic(err)
	}

	return name
}

// https://github.com/owncloud/ocis/accounts/issues/61
func TestCreateAccount(t *testing.T) {
	resp, err := createAccount(t, "user1")
	assert.NoError(t, err)
	assertUserExists(t, getAccount("user1"))
	assert.IsType(t, &proto.Account{}, resp)
	// Account is not returned in response
	// assertAccountsSame(t, getAccount("user1"), resp)

	resp, err = createAccount(t, "user2")
	assert.NoError(t, err)
	assertUserExists(t, getAccount("user2"))
	assert.IsType(t, &proto.Account{}, resp)
	// Account is not returned in response
	// assertAccountsSame(t, getAccount("user2"), resp)

	cleanUp(t)
}

// https://github.com/owncloud/ocis-accounts/issues/62
func TestCreateExistingUser(t *testing.T) {
	var err error
	_, err = createAccount(t, "user1")
	assert.NoError(t, err)

	_, err = createAccount(t, "user1")
	assert.Error(t, err)
	assertUserExists(t, getAccount("user1"))

	cleanUp(t)
}

// All tests fail after running this
// https://github.com/owncloud/ocis/accounts-issues/62
func TestCreateAccountInvalidUserName(t *testing.T) {

	resp, err := listAccounts(t)
	assert.NoError(t, err)
	numAccounts := len(resp.GetAccounts())

	testData := []string{
		"",
		"0",
		"#&@#",
		".._$%203",
	}

	for _, userName := range testData {
		_, err := createAccount(t, userName)

		// Should give error
		if err == nil {
			t.Fatalf("Expected an Error when creating user '%s' but got nil", userName)
		}
	}

	// resp should have the same number of accounts
	resp, err = listAccounts(t)
	assert.NoError(t, err)

	assert.Equal(t, numAccounts, len(resp.GetAccounts()))

	cleanUp(t)
}

func TestUpdateAccount(t *testing.T) {
	tests := []struct {
		name                string
		userAccount         *proto.Account
		expectedErrOnUpdate error
	}{
		{
			"Update user (demonstration of updatable fields)",
			&proto.Account{
				DisplayName:              "Alice Hansen",
				PreferredName:            "Wonderful-Alice",
				OnPremisesSamAccountName: "Alice",
				UidNumber:                20010,
				GidNumber:                30001,
				Mail:                     "alice@example.com",
			},
			nil,
		},
		{
			"Update user with unicode data",
			&proto.Account{
				DisplayName:              "एलिस हेन्सेन",
				PreferredName:            "अद्भुत-एलिस",
				OnPremisesSamAccountName: "एलिस",
				UidNumber:                20010,
				GidNumber:                30001,
				Mail:                     "एलिस@उदाहरण.com",
			},
			merrors.BadRequest(".", "preferred_name 'अद्भुत-एलिस' must be at least the local part of an email"),
		},
		{
			"Update user with empty data values",
			&proto.Account{
				DisplayName:              "",
				PreferredName:            "",
				OnPremisesSamAccountName: "",
				UidNumber:                0,
				GidNumber:                0,
				Mail:                     "",
			},
			merrors.BadRequest(".", "preferred_name '' must be at least the local part of an email"),
		},
		{
			"Update user with strange data",
			&proto.Account{
				DisplayName:              "12345",
				PreferredName:            "a12345",
				OnPremisesSamAccountName: "a54321",
				UidNumber:                1000,
				GidNumber:                1000,
				Mail:                     "1.2@3.c_@",
			},
			merrors.BadRequest(".", "mail '1.2@3.c_@' must be a valid email"),
		},
	}

	for _, tt := range tests {
		// updatable fields for type Account
		updateMask := []string{
			"AccountEnabled",
			"IsResourceAccount",
			"DisplayName",
			"PreferredName",
			"OnPremisesSamAccountName",
			"UidNumber",
			"GidNumber",
			"Mail",
		}

		t.Run(tt.name, func(t *testing.T) {
			acc, err := createAccount(t, "user1")
			assert.NoError(t, err)

			tt.userAccount.Id = acc.Id
			tt.userAccount.AccountEnabled = false
			tt.userAccount.IsResourceAccount = false
			resp, err := updateAccount(t, tt.userAccount, updateMask)
			if tt.expectedErrOnUpdate != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErrOnUpdate.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.IsType(t, &proto.Account{}, resp)
				assertAccountsSame(t, tt.userAccount, resp)
				assertUserExists(t, tt.userAccount)
			}
			cleanUp(t)
		})
	}
}

func TestUpdateNonUpdatableFieldsInAccount(t *testing.T) {
	_, err := createAccount(t, "user1")
	assert.NoError(t, err)

	tests := []struct {
		name        string
		updateMask  []string
		userAccount *proto.Account
	}{
		{
			"Try to update creation type",
			[]string{
				"CreationType",
			},
			&proto.Account{
				Id:           user1.Id,
				CreationType: "Type Test",
			},
		},
		{
			"Try to update password profile",
			[]string{
				"PasswordProfile",
			},
			&proto.Account{
				Id:              user1.Id,
				PasswordProfile: &proto.PasswordProfile{Password: "new password"},
			},
		},
		{
			"Try to update member of",
			[]string{
				"MemberOf",
			},
			&proto.Account{
				Id: user1.Id,
				MemberOf: []*proto.Group{
					{Id: "509a9dcd-bb37-4f4f-a01a-19dca27d9cfa"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := updateAccount(t, tt.userAccount, tt.updateMask)
			if err == nil {
				t.Fatalf("Expected error while updating non updatable field, but found none.")
			}
			assert.IsType(t, &proto.Account{}, res)
			assert.Empty(t, res)
			assert.Error(t, err)

			var e *merrors.Error

			if errors.As(err, &e) {
				assert.EqualValues(t, 400, e.Code)
				assert.Equal(t, "Bad Request", e.Status)

				errMsg := fmt.Sprintf("can not update field %s, either unknown or readonly", tt.updateMask[0])
				assert.Equal(t, errMsg, e.Detail)
			} else {
				t.Fatal("Expected merrors errors but found something else.")
			}
		})
	}

	cleanUp(t)
}

func TestListAccounts(t *testing.T) {
	_, err := createAccount(t, "user1")
	assert.NoError(t, err)
	_, err = createAccount(t, "user2")
	assert.NoError(t, err)

	resp, err := listAccounts(t)
	assert.NoError(t, err)

	assert.IsType(t, &proto.ListAccountsResponse{}, resp)
	assert.Equal(t, 9, len(resp.Accounts))

	assertResponseContainsUser(t, resp, getAccount("user1"))
	assertResponseContainsUser(t, resp, getAccount("user2"))

	cleanUp(t)
}

func TestListWithoutUserCreation(t *testing.T) {
	resp, err := listAccounts(t)
	assert.NoError(t, err)

	// Only 7 default users
	assert.Equal(t, 7, len(resp.Accounts))
	cleanUp(t)
}

func TestListAccountsWithFilterQuery(t *testing.T) {
	scenarios := []struct {
		name        string
		query       string
		expectedIDs []string
	}{
		// FIXME: disabled test scenarios need to be supported when implementing OData support
		// OData implementation tracked in https://github.com/owncloud/ocis/issues/716
		//{
		//	name:        "ListAccounts with exact match on preferred_name",
		//	query:       "preferred_name eq 'user1'",
		//	expectedIDs: []string{user1.Id},
		//},
		{
			name:        "ListAccounts with exact match on on_premises_sam_account_name",
			query:       "on_premises_sam_account_name eq 'user1'",
			expectedIDs: []string{user1.Id},
		},
		{
			name:        "ListAccounts with exact match on mail",
			query:       "mail eq 'user1@example.com'",
			expectedIDs: []string{user1.Id},
		},
		//{
		//	name:        "ListAccounts with exact match on id",
		//	query:       "id eq 'f9149a32-2b8e-4f04-9e8d-937d81712b9a'",
		//	expectedIDs: []string{user1.Id},
		//},
		//{
		//	name:        "ListAccounts without match on preferred_name",
		//	query:       "preferred_name eq 'wololo'",
		//	expectedIDs: []string{},
		//},
		//{
		//	name:        "ListAccounts with exact match on preferred_name AND mail",
		//	query:       "preferred_name eq 'user1' and mail eq 'user1@example.com'",
		//	expectedIDs: []string{user1.Id},
		//},
		//{
		//	name:        "ListAccounts without match on preferred_name AND mail",
		//	query:       "preferred_name eq 'user1' and mail eq 'wololo@example.com'",
		//	expectedIDs: []string{},
		//},
		//{
		//	name:        "ListAccounts with exact match on preferred_name OR mail, preferred_name exists, mail exists",
		//	query:       "preferred_name eq 'user1' or mail eq 'user1@example.com'",
		//	expectedIDs: []string{user1.Id},
		//},
		//{
		//	name:        "ListAccounts with exact match on preferred_name OR mail, preferred_name exists, mail does not exist",
		//	query:       "preferred_name eq 'user1' or mail eq 'wololo@example.com'",
		//	expectedIDs: []string{user1.Id},
		//},
		//{
		//	name:        "ListAccounts with exact match on preferred_name OR mail, preferred_name does not exists, mail exists",
		//	query:       "preferred_name eq 'wololo' or mail eq 'user1@example.com'",
		//	expectedIDs: []string{user1.Id},
		//},
		//{
		//	name:        "ListAccounts without match on preferred_name OR mail, preferred_name and mail do not exist",
		//	query:       "preferred_name eq 'wololo' or mail eq 'wololo@example.com'",
		//	expectedIDs: []string{},
		//},
		//{
		//	name:        "ListAccounts with multiple matches on preferred_name",
		//	query:       "startswith(preferred_name,'user*')",
		//	expectedIDs: []string{user1.Id, user2.Id},
		//},
		//{
		//	name:        "ListAccounts with multiple matches on on_premises_sam_account_name",
		//	query:       "startswith(on_premises_sam_account_name,'user*')",
		//	expectedIDs: []string{user1.Id, user2.Id},
		//},
	}

	cl := proto.NewAccountsService("com.owncloud.api.accounts", service.Client())

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			_, err := createAccount(t, "user1")
			assert.NoError(t, err)
			_, err = createAccount(t, "user2")
			assert.NoError(t, err)

			req := &proto.ListAccountsRequest{Query: scenario.query}
			res, err := cl.ListAccounts(context.Background(), req)
			assert.NoError(t, err)
			ids := make([]string, 0)
			for _, acc := range res.Accounts {
				ids = append(ids, acc.Id)
			}
			assert.Equal(t, scenario.expectedIDs, ids)
			cleanUp(t)
		})
	}
}

func TestGetAccount(t *testing.T) {
	_, err := createAccount(t, "user1")
	assert.NoError(t, err)

	req := &proto.GetAccountRequest{Id: getAccount("user1").Id}

	cl := proto.NewAccountsService("com.owncloud.api.accounts", service.Client())

	resp, err := cl.GetAccount(context.Background(), req)

	assert.NoError(t, err)
	assert.IsType(t, &proto.Account{}, resp)
	assertAccountsSame(t, getAccount("user1"), resp)

	cleanUp(t)
}

//TODO: This segfaults! WIP

func TestDeleteAccount(t *testing.T) {
	_, err := createAccount(t, "user1")
	assert.NoError(t, err)
	_, err = createAccount(t, "user2")
	assert.NoError(t, err)

	req := &proto.DeleteAccountRequest{Id: getAccount("user1").Id}

	client := mgrpcc.NewClient()
	cl := proto.NewAccountsService("com.owncloud.api.accounts", client)

	resp, err := cl.DeleteAccount(context.Background(), req)
	assert.NoError(t, err)
	assert.IsType(t, resp, &empty.Empty{})

	// Check the account doesn't exists anymore
	accountList, _ := listAccounts(t)
	assertResponseContainsUser(t, accountList, getAccount("user2"))
	assertResponseNotContainsUser(t, accountList, getAccount("user1"))

	cleanUp(t)
}

func TestListGroups(t *testing.T) {
	req := &proto.ListGroupsRequest{}

	client := mgrpcc.NewClient()
	cl := proto.NewGroupsService("com.owncloud.api.accounts", client)

	resp, err := cl.ListGroups(context.Background(), req)
	assert.NoError(t, err)
	assert.IsType(t, &proto.ListGroupsResponse{}, resp)
	assert.Equal(t, 9, len(resp.Groups))

	groups := []string{
		"sysusers",
		"users",
		"sailing-lovers",
		"violin-haters",
		"radium-lovers",
		"polonium-lovers",
		"quantum-lovers",
		"philosophy-haters",
		"physics-lovers",
	}

	for _, g := range groups {
		assertResponseContainsGroup(t, resp, getGroup(g))
	}
	cleanUp(t)
}

func TestGetGroups(t *testing.T) {
	client := mgrpcc.NewClient()
	cl := proto.NewGroupsService("com.owncloud.api.accounts", client)

	groups := []string{
		"sysusers",
		"users",
		"sailing-lovers",
		"violin-haters",
		"radium-lovers",
		"polonium-lovers",
		"quantum-lovers",
		"philosophy-haters",
		"physics-lovers",
	}

	for _, g := range groups {
		group := getGroup(g)
		req := &proto.GetGroupRequest{Id: group.Id}
		resp, err := cl.GetGroup(context.Background(), req)

		assert.NoError(t, err)
		assert.IsType(t, &proto.Group{}, resp)
		assertGroupsSame(t, group, resp)
	}
	cleanUp(t)
}

// https://github.com/owncloud/ocis/accounts/issues/61
func TestCreateGroup(t *testing.T) {
	group := &proto.Group{Id: "2d58e5ec-842e-498b-8800-61f2ec6f911f", GidNumber: 30042, OnPremisesSamAccountName: "quantum-group", DisplayName: "Quantum Group", Members: []*proto.Account{
		{Id: "4c510ada-c86b-4815-8820-42cdf82c3d51"}, // einstein
		{Id: "932b4540-8d16-481e-8ef4-588e4b6b151c"}, // feynman
	}}

	res, err := createGroup(t, group)
	assert.NoError(t, err)

	assert.IsType(t, &proto.Group{}, res)

	// Should return the group but does not
	// assertGroupsSame(t, res, group)

	groupsResponse := listGroups(t)
	assertResponseContainsGroup(t, groupsResponse, group)
	cleanUp(t)
}

func TestGetGroupInvalidID(t *testing.T) {
	client := mgrpcc.NewClient()
	cl := proto.NewGroupsService("com.owncloud.api.accounts", client)

	req := &proto.GetGroupRequest{Id: "42"}
	resp, err := cl.GetGroup(context.Background(), req)

	assert.IsType(t, &proto.Group{}, resp)
	assert.Empty(t, resp)
	assert.Error(t, err)
	cleanUp(t)
}

func TestDeleteGroup(t *testing.T) {
	grp1 := getTestGroups("grp1")
	grp2 := getTestGroups("grp2")
	grp3 := getTestGroups("grp3")
	_, err := createGroup(t, grp1)
	assert.NoError(t, err)
	_, err = createGroup(t, grp2)
	assert.NoError(t, err)
	_, err = createGroup(t, grp3)
	assert.NoError(t, err)

	client := mgrpcc.NewClient()
	cl := proto.NewGroupsService("com.owncloud.api.accounts", client)

	req := &proto.DeleteGroupRequest{Id: grp1.Id}
	res, err := cl.DeleteGroup(context.Background(), req)
	assert.IsType(t, res, &empty.Empty{})
	assert.NoError(t, err)

	req = &proto.DeleteGroupRequest{Id: grp2.Id}
	res, err = cl.DeleteGroup(context.Background(), req)
	assert.IsType(t, res, &empty.Empty{})
	assert.NoError(t, err)

	groupsResponse := listGroups(t)
	assertResponseNotContainsGroup(t, groupsResponse, grp1)
	assertResponseNotContainsGroup(t, groupsResponse, grp2)
	assertResponseContainsGroup(t, groupsResponse, grp3)
	cleanUp(t)
}

func TestDeleteGroupNotExisting(t *testing.T) {
	invalidIds := []string{
		"$@dsfd",
		"42",
		"happyString",
		"0ed84f08-aa0a-46e4-8e42-f0a5d6e1b059",
		"  ",
	}

	client := mgrpcc.NewClient()
	cl := proto.NewGroupsService("com.owncloud.api.accounts", client)

	for _, id := range invalidIds {
		req := &proto.DeleteGroupRequest{Id: id}
		res, err := cl.DeleteGroup(context.Background(), req)
		assert.IsType(t, &empty.Empty{}, res)
		assert.Empty(t, res)
		assert.Error(t, err)
	}
	cleanUp(t)
}

func TestDeleteGroupInvalidId(t *testing.T) {
	invalidIds := map[string]string{
		".":                                     ".",
		"hello/world":                           "hello/world",
		"/new-id":                               "/new-id",
		"/new-id/":                              "/new-id",
		"/0ed84f08-aa0a-46e4-8e42-f0a5d6e1b059": "/0ed84f08-aa0a-46e4-8e42-f0a5d6e1b059",
		"":                                      ".",
	}

	client := mgrpcc.NewClient()
	cl := proto.NewGroupsService("com.owncloud.api.accounts", client)

	for id := range invalidIds {
		req := &proto.DeleteGroupRequest{Id: id}
		res, err := cl.DeleteGroup(context.Background(), req)
		assert.IsType(t, &empty.Empty{}, res)
		assert.Empty(t, res)
		assert.Error(t, err)
	}
	cleanUp(t)
}

func TestUpdateGroup(t *testing.T) {
	grp1 := getTestGroups("grp1")
	_, err := createGroup(t, grp1)
	assert.NoError(t, err)

	client := mgrpcc.NewClient()
	cl := proto.NewGroupsService("com.owncloud.api.accounts", client)

	updateGrp := &proto.Group{
		Id: grp1.Id,
	}

	req := &proto.UpdateGroupRequest{Group: updateGrp}

	res, err := cl.UpdateGroup(context.Background(), req)

	assert.IsType(t, &proto.Group{}, res)
	assert.Empty(t, res)
	assert.Error(t, err)

	cleanUp(t)
}

// https://github.com/owncloud/ocis/accounts/issues/61
func TestAddMember(t *testing.T) {
	grp1 := getTestGroups("grp1")
	account := getAccount("user1")

	updatedGroup := grp1
	updatedGroup.Members = append(updatedGroup.Members, &proto.Account{Id: account.Id})

	_, err := createGroup(t, grp1)
	assert.NoError(t, err)
	_, err = createAccount(t, account.PreferredName)
	assert.NoError(t, err)

	client := mgrpcc.NewClient()
	cl := proto.NewGroupsService("com.owncloud.api.accounts", client)

	req := &proto.AddMemberRequest{GroupId: grp1.Id, AccountId: account.Id}

	res, err := cl.AddMember(context.Background(), req)
	assert.NoError(t, err)

	assert.IsType(t, &proto.Group{}, res)

	// Should return the group but returns empty
	// assertGroupsSame(t, updatedGroup, res)

	resp := listGroups(t)
	assertResponseContainsGroup(t, resp, updatedGroup)

	cleanUp(t)
}

// https://github.com/owncloud/ocis/accounts/issues/62
func TestAddMemberAlreadyInGroup(t *testing.T) {
	grp1 := getTestGroups("grp1")
	account := getAccount("user1")

	updatedGroup := grp1
	updatedGroup.Members = append(updatedGroup.Members, &proto.Account{Id: account.Id})

	_, err := createGroup(t, grp1)
	assert.NoError(t, err)
	_, err = createAccount(t, account.PreferredName)
	assert.NoError(t, err)

	_, err = addMemberToGroup(t, grp1.Id, account.Id)
	assert.NoError(t, err)

	client := mgrpcc.NewClient()
	cl := proto.NewGroupsService("com.owncloud.api.accounts", client)

	req := &proto.AddMemberRequest{GroupId: grp1.Id, AccountId: account.Id}

	res, err := cl.AddMember(context.Background(), req)

	// Should Give Error
	assert.NoError(t, err)
	assert.IsType(t, &proto.Group{}, res)
	//assert.Equal(t, proto.Group{}, *res)
	//assertGroupsSame(t, updatedGroup, res)

	// Check the group is truly updated
	resp := listGroups(t)
	assertResponseContainsGroup(t, resp, updatedGroup)

	cleanUp(t)
}

func TestAddMemberNonExisting(t *testing.T) {
	grp1 := getTestGroups("grp1")

	_, err := createGroup(t, grp1)
	assert.NoError(t, err)

	client := mgrpcc.NewClient()
	cl := proto.NewGroupsService("com.owncloud.api.accounts", client)

	invalidIds := []string{
		"$@dsfd",
		"42",
		"happyString",
		"0ed84f08-aa0a-46e4-8e42-f0a5d6e1b059",
		"  ",
	}

	for _, id := range invalidIds {
		req := &proto.AddMemberRequest{GroupId: grp1.Id, AccountId: id}

		res, err := cl.AddMember(context.Background(), req)
		assert.IsType(t, &proto.Group{}, res)
		assert.Empty(t, res)
		assert.Error(t, err)
	}

	// Check group is not changed
	resp := listGroups(t)
	assertResponseContainsGroup(t, resp, grp1)

	cleanUp(t)
}

func addMemberToGroup(t *testing.T, groupID, memberID string) (*proto.Group, error) {
	client := mgrpcc.NewClient()
	cl := proto.NewGroupsService("com.owncloud.api.accounts", client)

	req := &proto.AddMemberRequest{GroupId: groupID, AccountId: memberID}

	res, err := cl.AddMember(context.Background(), req)

	return res, err
}

// https://github.com/owncloud/ocis/accounts/issues/61
func TestRemoveMember(t *testing.T) {
	grp1 := getTestGroups("grp1")
	account := getAccount("user1")

	_, err := createGroup(t, grp1)
	assert.NoError(t, err)
	_, err = createAccount(t, account.PreferredName)
	assert.NoError(t, err)

	_, err = addMemberToGroup(t, grp1.Id, account.Id)
	assert.NoError(t, err)

	client := mgrpcc.NewClient()
	cl := proto.NewGroupsService("com.owncloud.api.accounts", client)

	req := &proto.RemoveMemberRequest{GroupId: grp1.Id, AccountId: account.Id}

	res, err := cl.RemoveMember(context.Background(), req)
	assert.NoError(t, err)

	assert.IsType(t, &proto.Group{}, res)
	//assert.Equal(t, proto.Group{}, *res)
	// assertGroupsSame(t, grp1, res)

	resp := listGroups(t)
	assertResponseContainsGroup(t, resp, grp1)

	cleanUp(t)
}

func TestRemoveMemberNonExistingUser(t *testing.T) {
	grp1 := getTestGroups("grp1")

	_, err := createGroup(t, grp1)
	assert.NoError(t, err)

	client := mgrpcc.NewClient()
	cl := proto.NewGroupsService("com.owncloud.api.accounts", client)

	invalidIds := []string{
		"$@dsfd",
		"42",
		"happyString",
		"0ed84f08-aa0a-46e4-8e42-f0a5d6e1b059",
		"  ",
	}

	for _, id := range invalidIds {
		req := &proto.RemoveMemberRequest{GroupId: grp1.Id, AccountId: id}

		res, err := cl.RemoveMember(context.Background(), req)
		assert.IsType(t, &proto.Group{}, res)
		assert.Empty(t, res)
		assert.Error(t, err)
	}

	// Check group is not changed
	resp := listGroups(t)
	assertResponseContainsGroup(t, resp, grp1)
	cleanUp(t)
}

// https://github.com/owncloud/ocis/accounts/issues/62
func TestRemoveMemberNotInGroup(t *testing.T) {
	grp1 := getTestGroups("grp1")
	account := getAccount("user1")

	_, err := createGroup(t, grp1)
	assert.NoError(t, err)
	_, err = createAccount(t, account.PreferredName)
	assert.NoError(t, err)

	client := mgrpcc.NewClient()
	cl := proto.NewGroupsService("com.owncloud.api.accounts", client)

	req := &proto.RemoveMemberRequest{GroupId: grp1.Id, AccountId: account.Id}

	res, err := cl.RemoveMember(context.Background(), req)

	// Should give an error
	assert.NoError(t, err)
	assert.IsType(t, &proto.Group{}, res)

	//assert.Error(t, err)
	//assert.Equal(
	//	t,
	//	fmt.Sprintf("{\"id\":\".\",\"code\":404,\"detail\":\"User not found in the group\",\"status\":\"Not Found\"}", account.Id),
	//	err.Error(),
	//)

	// Check group is not changed
	resp := listGroups(t)
	assertResponseContainsGroup(t, resp, grp1)
	cleanUp(t)
}

func TestListMembers(t *testing.T) {
	groups := []string{
		"sysusers",
		"users",
		"sailing-lovers",
		"violin-haters",
		"radium-lovers",
		"polonium-lovers",
		"quantum-lovers",
		"philosophy-haters",
		"physics-lovers",
	}

	client := mgrpcc.NewClient()
	cl := proto.NewGroupsService("com.owncloud.api.accounts", client)

	for _, group := range groups {
		expectedGroup := getGroup(group)
		req := &proto.ListMembersRequest{Id: expectedGroup.Id}

		res, err := cl.ListMembers(context.Background(), req)
		assert.NoError(t, err)

		assert.Equal(t, len(res.Members), len(expectedGroup.Members))

		for _, member := range expectedGroup.Members {
			found := false
			for _, resMember := range res.Members {
				if resMember.Id == member.Id {
					found = true
					break
				}
			}

			if !found {
				t.Fatalf("Group with Id %v Expected to be in response but not found", member.Id)
			}
		}
	}
	cleanUp(t)
}

func TestListMembersEmptyGroup(t *testing.T) {
	group := &proto.Group{Id: "5d58e5ec-842e-498b-8800-61f2ec6f911c", GidNumber: 60000, OnPremisesSamAccountName: "quantum-group", DisplayName: "Quantum Group", Members: []*proto.Account{}}

	client := mgrpcc.NewClient()
	cl := proto.NewGroupsService("com.owncloud.api.accounts", client)

	request := &proto.CreateGroupRequest{Group: group}
	_, err := cl.CreateGroup(context.Background(), request)
	if err == nil {
		newCreatedGroups = append(newCreatedGroups, group.Id)
	}

	req := &proto.ListMembersRequest{Id: group.Id}

	listRes, err := cl.ListMembers(context.Background(), req)

	assert.NoError(t, err)
	assert.Empty(t, listRes.Members)

	cleanUp(t)
}

func TestAccountUpdateMask(t *testing.T) {
	_, err := createAccount(t, "user1")
	assert.NoError(t, err)

	user1 := getAccount("user1")
	client := mgrpcc.NewClient()
	req := &proto.UpdateAccountRequest{
		// We only want to update the display-name, rest should be ignored
		UpdateMask: &field_mask.FieldMask{Paths: []string{"DisplayName"}},
		Account: &proto.Account{
			Id:            user1.Id,
			DisplayName:   "ShouldBeUpdated",
			PreferredName: "ShouldStaySame And Is Invalid Anyway",
		}}

	cl := proto.NewAccountsService("com.owncloud.api.accounts", client)
	res, err := cl.UpdateAccount(context.Background(), req)
	assert.NoError(t, err)

	assert.Equal(t, "ShouldBeUpdated", res.DisplayName)
	assert.Equal(t, user1.PreferredName, res.PreferredName)

	cleanUp(t)
}

func TestAccountUpdateReadOnlyField(t *testing.T) {
	_, err := createAccount(t, "user1")
	assert.NoError(t, err)

	user1 := getAccount("user1")
	client := mgrpcc.NewClient()
	req := &proto.UpdateAccountRequest{
		// We only want to update the display-name, rest should be ignored
		UpdateMask: &field_mask.FieldMask{Paths: []string{"CreatedDateTime"}},
		Account: &proto.Account{
			Id:              user1.Id,
			CreatedDateTime: timestamppb.Now(),
		}}

	cl := proto.NewAccountsService("com.owncloud.api.accounts", client)
	res, err := cl.UpdateAccount(context.Background(), req)
	assert.Nil(t, res)
	assert.Error(t, err)

	var e *merrors.Error

	if errors.As(err, &e) {
		assert.EqualValues(t, 400, e.Code)
		assert.Equal(t, "Bad Request", e.Status)
	} else {
		t.Fatal("Unexpected error type")
	}

	cleanUp(t)
}
