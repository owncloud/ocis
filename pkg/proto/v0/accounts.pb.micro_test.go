package proto_test

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/micro/go-micro/v2/client"
	merrors "github.com/micro/go-micro/v2/errors"
	"github.com/owncloud/ocis-accounts/pkg/command"
	"github.com/owncloud/ocis-accounts/pkg/config"
	"github.com/owncloud/ocis-accounts/pkg/proto/v0"
	svc "github.com/owncloud/ocis-accounts/pkg/service/v0"
	"github.com/owncloud/ocis-pkg/v2/service/grpc"
	settings "github.com/owncloud/ocis-settings/pkg/proto/v0"
	"github.com/stretchr/testify/assert"
	"google.golang.org/genproto/protobuf/field_mask"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var service = grpc.Service{}

const dataPath = "./accounts-store"

var newCreatedAccounts = []string{}
var newCreatedGroups = []string{}

var mockedRoleAssignment = map[string]string{}

func getAccount(user string) *proto.Account {
	switch user {
	case "user1":
		return &proto.Account{
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
	case "user2":
		return &proto.Account{
			Id:                       "e9149a32-2b8e-4f04-9e8d-937d81712b9a",
			AccountEnabled:           true,
			IsResourceAccount:        true,
			CreationType:             "",
			DisplayName:              "User Two",
			PreferredName:            "user2",
			OnPremisesSamAccountName: "user2",
			UidNumber:                20009,
			GidNumber:                30000,
			Mail:                     "user2@example.com",
			Identities:               []*proto.Identities{nil},
			PasswordProfile:          &proto.PasswordProfile{Password: "hello123"},
			MemberOf: []*proto.Group{
				{Id: "509a9dcd-bb37-4f4f-a01a-19dca27d9cfa"}, // users
			},
		}
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
			{Id: "820ba2a1-3f54-4538-80a4-2d73007e30bf"}, // konnectd
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

func init() {
	service = grpc.NewService(
		grpc.Namespace("com.owncloud.api"),
		grpc.Name("accounts"),
		grpc.Address("localhost:9180"),
	)

	cfg := config.New()
	cfg.Server.AccountsDataPath = dataPath
	var hdlr *svc.Service
	var err error

	if hdlr, err = svc.New(svc.Logger(command.NewLogger(cfg)), svc.Config(cfg), svc.RoleService(buildRoleServiceMock())); err != nil {
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
		_, err := deleteAccount(t, id)
		checkNoError(t, err)
	}

	datastore = filepath.Join(dataPath, "groups")

	for _, id := range newCreatedGroups {
		path := filepath.Join(datastore, id)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			continue
		}
		_, err := deleteGroup(t, id)
		checkNoError(t, err)
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

func assertGroupHasMember(t *testing.T, grp *proto.Group, memberId string) {
	for _, m := range grp.Members {
		if m.Id == memberId {
			return
		}
	}

	t.Fatalf("Member with id %s expected to be in group '%s', but not found", memberId, grp.DisplayName)
}

func checkNoError(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("Expected Error to be nil but got %s", err)
	}
}

func createAccount(t *testing.T, user string) (*proto.Account, error) {
	client := service.Client()
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
	client := service.Client()
	cl := proto.NewGroupsService("com.owncloud.api.accounts", client)

	request := &proto.CreateGroupRequest{Group: group}
	res, err := cl.CreateGroup(context.Background(), request)
	if err == nil {
		newCreatedGroups = append(newCreatedGroups, group.Id)
	}
	return res, err
}

func updateAccount(t *testing.T, account *proto.Account, updateArray []string) (*proto.Account, error) {
	client := service.Client()
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
	client := service.Client()
	cl := proto.NewAccountsService("com.owncloud.api.accounts", client)

	response, err := cl.ListAccounts(context.Background(), request)
	return response, err
}

func listGroups(t *testing.T) *proto.ListGroupsResponse {
	request := &proto.ListGroupsRequest{}
	client := service.Client()
	cl := proto.NewGroupsService("com.owncloud.api.accounts", client)

	response, err := cl.ListGroups(context.Background(), request)
	checkNoError(t, err)
	return response
}

func deleteAccount(t *testing.T, id string) (*empty.Empty, error) {
	client := service.Client()
	cl := proto.NewAccountsService("com.owncloud.api.accounts", client)

	req := &proto.DeleteAccountRequest{Id: id}
	res, err := cl.DeleteAccount(context.Background(), req)
	return res, err
}

func deleteGroup(t *testing.T, id string) (*empty.Empty, error) {
	client := service.Client()
	cl := proto.NewGroupsService("com.owncloud.api.accounts", client)

	req := &proto.DeleteGroupRequest{Id: id}
	res, err := cl.DeleteGroup(context.Background(), req)
	return res, err
}

// https://github.com/owncloud/ocis-accounts/issues/61
func TestCreateAccount(t *testing.T) {

	resp, err := createAccount(t, "user1")
	checkNoError(t, err)
	assertUserExists(t, getAccount("user1"))
	assert.IsType(t, &proto.Account{}, resp)
	// Account is not returned in response
	// assertAccountsSame(t, getAccount("user1"), resp)

	resp, err = createAccount(t, "user2")
	checkNoError(t, err)
	assertUserExists(t, getAccount("user2"))
	assert.IsType(t, &proto.Account{}, resp)
	// Account is not returned in response
	// assertAccountsSame(t, getAccount("user2"), resp)

	cleanUp(t)
}

// https://github.com/owncloud/ocis-accounts/issues/62
func TestCreateExistingUser(t *testing.T) {
	createAccount(t, "user1")
	_, err := createAccount(t, "user1")

	// Should give error but it does not
	checkNoError(t, err)
	assertUserExists(t, getAccount("user1"))

	cleanUp(t)
}

// All tests fail after running this
// https://github.com/owncloud/ocis-accounts/issues/62
func TestCreateAccountInvalidUserName(t *testing.T) {

	resp, err := listAccounts(t)
	checkNoError(t, err)
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
	checkNoError(t, err)

	assert.Equal(t, numAccounts, len(resp.GetAccounts()))

	cleanUp(t)
}

func TestUpdateAccount(t *testing.T) {
	_, _ = createAccount(t, "user1")

	tests := []struct {
		name        string
		userAccount *proto.Account
	}{
		{
			"Update user (demonstration of updatable fields)",
			&proto.Account{
				DisplayName:                 "Alice Hansen",
				PreferredName:               "Wonderful Alice",
				OnPremisesDistinguishedName: "Alice",
				UidNumber:                   20010,
				GidNumber:                   30001,
				Mail:                        "alice@example.com",
			},
		},
		{
			"Update user with unicode data",
			&proto.Account{
				DisplayName:                 "एलिस हेन्सेन",
				PreferredName:               "अद्भुत एलिस",
				OnPremisesDistinguishedName: "एलिस",
				UidNumber:                   20010,
				GidNumber:                   30001,
				Mail:                        "एलिस@उदाहरण.com",
			},
		},
		{
			"Update user with empty data values",
			&proto.Account{
				DisplayName:                 "",
				PreferredName:               "",
				OnPremisesDistinguishedName: "",
				UidNumber:                   0,
				GidNumber:                   0,
				Mail:                        "",
			},
		},
		{
			"Update user with strange data",
			&proto.Account{
				DisplayName:                 "12345",
				PreferredName:               "12345",
				OnPremisesDistinguishedName: "54321",
				UidNumber:                   1000,
				GidNumber:                   1000,
				// No email validation
				// https://github.com/owncloud/ocis-accounts/issues/77
				Mail: "1.2@3.c_@",
			},
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
			tt.userAccount.Id = "f9149a32-2b8e-4f04-9e8d-937d81712b9a"
			tt.userAccount.AccountEnabled = false
			tt.userAccount.IsResourceAccount = false
			resp, err := updateAccount(t, tt.userAccount, updateMask)

			checkNoError(t, err)

			assert.IsType(t, &proto.Account{}, resp)
			assertAccountsSame(t, tt.userAccount, resp)
			assertUserExists(t, tt.userAccount)
		})
	}

	cleanUp(t)
}

func TestUpdateNonUpdatableFieldsInAccount(t *testing.T) {
	_, _ = createAccount(t, "user1")

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
				CreationType: "Type Test",
			},
		},
		{
			"Try to update password profile",
			[]string{
				"PasswordProfile",
			},
			&proto.Account{
				PasswordProfile: &proto.PasswordProfile{Password: "new password"},
			},
		},
		{
			"Try to update member of",
			[]string{
				"MemberOf",
			},
			&proto.Account{
				MemberOf: []*proto.Group{
					{Id: "509a9dcd-bb37-4f4f-a01a-19dca27d9cfa"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.userAccount.Id = "f9149a32-2b8e-4f04-9e8d-937d81712b9a"
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
				t.Fatal("Expected merror errors but found something else.")
			}
		})
	}
}

func TestListAccounts(t *testing.T) {
	createAccount(t, "user1")
	createAccount(t, "user2")

	resp, err := listAccounts(t)
	checkNoError(t, err)

	assert.IsType(t, &proto.ListAccountsResponse{}, resp)
	assert.Equal(t, 8, len(resp.Accounts))

	assertResponseContainsUser(t, resp, getAccount("user1"))
	assertResponseContainsUser(t, resp, getAccount("user2"))

	cleanUp(t)
}

func TestListWithoutUserCreation(t *testing.T) {
	resp, err := listAccounts(t)

	checkNoError(t, err)

	// Only 5 default users
	assert.Equal(t, 6, len(resp.Accounts))
	cleanUp(t)
}

func TestGetAccount(t *testing.T) {
	createAccount(t, "user1")

	req := &proto.GetAccountRequest{Id: getAccount("user1").Id}

	client := service.Client()
	cl := proto.NewAccountsService("com.owncloud.api.accounts", client)

	resp, err := cl.GetAccount(context.Background(), req)

	checkNoError(t, err)
	assert.IsType(t, &proto.Account{}, resp)
	assertAccountsSame(t, getAccount("user1"), resp)

	cleanUp(t)
}

func TestDeleteAccount(t *testing.T) {
	createAccount(t, "user1")
	createAccount(t, "user2")

	req := &proto.DeleteAccountRequest{Id: getAccount("user1").Id}

	client := service.Client()
	cl := proto.NewAccountsService("com.owncloud.api.accounts", client)

	resp, err := cl.DeleteAccount(context.Background(), req)
	checkNoError(t, err)
	assert.IsType(t, resp, &empty.Empty{})

	// Check the account doesn't exists anymore
	accountList, _ := listAccounts(t)
	assertResponseContainsUser(t, accountList, getAccount("user2"))
	assertResponseNotContainsUser(t, accountList, getAccount("user1"))

	cleanUp(t)
}

func TestListGroups(t *testing.T) {
	req := &proto.ListGroupsRequest{}

	client := service.Client()
	cl := proto.NewGroupsService("com.owncloud.api.accounts", client)

	resp, err := cl.ListGroups(context.Background(), req)
	checkNoError(t, err)
	assert.IsType(t, &proto.ListGroupsResponse{}, resp)
	assert.Equal(t, len(resp.Groups), 9)

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
	client := service.Client()
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

		checkNoError(t, err)
		assert.IsType(t, &proto.Group{}, resp)
		assertGroupsSame(t, group, resp)
	}
	cleanUp(t)
}

// https://github.com/owncloud/ocis-accounts/issues/61
func TestCreateGroup(t *testing.T) {
	group := &proto.Group{Id: "2d58e5ec-842e-498b-8800-61f2ec6f911f", GidNumber: 30042, OnPremisesSamAccountName: "quantum-group", DisplayName: "Quantum Group", Members: []*proto.Account{
		{Id: "4c510ada-c86b-4815-8820-42cdf82c3d51"}, // einstein
		{Id: "932b4540-8d16-481e-8ef4-588e4b6b151c"}, // feynman
	}}

	res, err := createGroup(t, group)
	checkNoError(t, err)

	assert.IsType(t, &proto.Group{}, res)

	// Should return the group but does not
	// assertGroupsSame(t, res, group)

	groupsResponse := listGroups(t)
	assertResponseContainsGroup(t, groupsResponse, group)
	cleanUp(t)
}

func TestGetGroupInvalidID(t *testing.T) {
	client := service.Client()
	cl := proto.NewGroupsService("com.owncloud.api.accounts", client)

	req := &proto.GetGroupRequest{Id: "42"}
	resp, err := cl.GetGroup(context.Background(), req)

	assert.IsType(t, &proto.Group{}, resp)
	assert.Empty(t, resp)
	assert.Error(t, err)
	assert.Equal(t, "{\"id\":\".\",\"code\":404,\"detail\":\"could not read group: open accounts-store/groups/42: no such file or directory\",\"status\":\"Not Found\"}", err.Error())
	cleanUp(t)
}

func TestDeleteGroup(t *testing.T) {
	grp1 := getTestGroups("grp1")
	grp2 := getTestGroups("grp2")
	grp3 := getTestGroups("grp3")
	createGroup(t, grp1)
	createGroup(t, grp2)
	createGroup(t, grp3)

	client := service.Client()
	cl := proto.NewGroupsService("com.owncloud.api.accounts", client)

	req := &proto.DeleteGroupRequest{Id: grp1.Id}
	res, err := cl.DeleteGroup(context.Background(), req)
	assert.IsType(t, res, &empty.Empty{})
	checkNoError(t, err)

	req = &proto.DeleteGroupRequest{Id: grp2.Id}
	res, err = cl.DeleteGroup(context.Background(), req)
	assert.IsType(t, res, &empty.Empty{})
	checkNoError(t, err)

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

	client := service.Client()
	cl := proto.NewGroupsService("com.owncloud.api.accounts", client)

	for _, id := range invalidIds {
		req := &proto.DeleteGroupRequest{Id: id}
		res, err := cl.DeleteGroup(context.Background(), req)
		assert.IsType(t, &empty.Empty{}, res)
		assert.Empty(t, res)
		assert.Error(t, err)
		assert.Equal(
			t,
			fmt.Sprintf("{\"id\":\".\",\"code\":404,\"detail\":\"could not read group: open accounts-store/groups/%v: no such file or directory\",\"status\":\"Not Found\"}", id),
			err.Error(),
		)
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

	client := service.Client()
	cl := proto.NewGroupsService("com.owncloud.api.accounts", client)

	for id, val := range invalidIds {
		req := &proto.DeleteGroupRequest{Id: id}
		res, err := cl.DeleteGroup(context.Background(), req)
		assert.IsType(t, &empty.Empty{}, res)
		assert.Empty(t, res)
		assert.Error(t, err)
		assert.Equal(
			t,
			fmt.Sprintf("{\"id\":\".\",\"code\":500,\"detail\":\"could not clean up group id: invalid id %v\",\"status\":\"Internal Server Error\"}", val),
			err.Error(),
		)
	}
	cleanUp(t)
}

func TestUpdateGroup(t *testing.T) {
	grp1 := getTestGroups("grp1")
	createGroup(t, grp1)

	client := service.Client()
	cl := proto.NewGroupsService("com.owncloud.api.accounts", client)

	updateGrp := &proto.Group{
		Id: grp1.Id,
	}

	req := &proto.UpdateGroupRequest{Group: updateGrp}

	res, err := cl.UpdateGroup(context.Background(), req)

	assert.IsType(t, &proto.Group{}, res)
	assert.Empty(t, res)
	assert.Error(t, err)
	assert.Equal(
		t,
		"{\"id\":\".\",\"code\":500,\"detail\":\"not implemented\",\"status\":\"Internal Server Error\"}",
		err.Error(),
	)
	cleanUp(t)
}

// https://github.com/owncloud/ocis-accounts/issues/61
func TestAddMember(t *testing.T) {
	grp1 := getTestGroups("grp1")
	account := getAccount("user1")

	updatedGroup := grp1
	updatedGroup.Members = append(updatedGroup.Members, &proto.Account{Id: account.Id})

	createGroup(t, grp1)
	createAccount(t, account.PreferredName)

	client := service.Client()
	cl := proto.NewGroupsService("com.owncloud.api.accounts", client)

	req := &proto.AddMemberRequest{GroupId: grp1.Id, AccountId: account.Id}

	res, err := cl.AddMember(context.Background(), req)
	checkNoError(t, err)

	assert.IsType(t, &proto.Group{}, res)

	// Should return the group but returns empty
	// assertGroupsSame(t, updatedGroup, res)

	resp := listGroups(t)
	assertResponseContainsGroup(t, resp, updatedGroup)

	cleanUp(t)
}

// https://github.com/owncloud/ocis-accounts/issues/62
func TestAddMemberAlreadyInGroup(t *testing.T) {
	grp1 := getTestGroups("grp1")
	account := getAccount("user1")

	updatedGroup := grp1
	updatedGroup.Members = append(updatedGroup.Members, &proto.Account{Id: account.Id})

	createGroup(t, grp1)
	createAccount(t, account.PreferredName)

	addMemberToGroup(t, grp1.Id, account.Id)

	client := service.Client()
	cl := proto.NewGroupsService("com.owncloud.api.accounts", client)

	req := &proto.AddMemberRequest{GroupId: grp1.Id, AccountId: account.Id}

	res, err := cl.AddMember(context.Background(), req)

	// Should Give Error
	checkNoError(t, err)
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

	createGroup(t, grp1)

	client := service.Client()
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
		assert.Equal(
			t,
			fmt.Sprintf("{\"id\":\".\",\"code\":404,\"detail\":\"could not read account: open accounts-store/accounts/%v: no such file or directory\",\"status\":\"Not Found\"}", id),
			err.Error(),
		)
	}

	// Check group is not changed
	resp := listGroups(t)
	assertResponseContainsGroup(t, resp, grp1)

	cleanUp(t)
}

func addMemberToGroup(t *testing.T, groupId, memberId string) (*proto.Group, error) {
	client := service.Client()
	cl := proto.NewGroupsService("com.owncloud.api.accounts", client)

	req := &proto.AddMemberRequest{GroupId: groupId, AccountId: memberId}

	res, err := cl.AddMember(context.Background(), req)

	return res, err
}

// https://github.com/owncloud/ocis-accounts/issues/61
func TestRemoveMember(t *testing.T) {
	grp1 := getTestGroups("grp1")
	account := getAccount("user1")

	createGroup(t, grp1)
	createAccount(t, account.PreferredName)

	addMemberToGroup(t, grp1.Id, account.Id)

	client := service.Client()
	cl := proto.NewGroupsService("com.owncloud.api.accounts", client)

	req := &proto.RemoveMemberRequest{GroupId: grp1.Id, AccountId: account.Id}

	res, err := cl.RemoveMember(context.Background(), req)
	checkNoError(t, err)

	assert.IsType(t, &proto.Group{}, res)
	//assert.Equal(t, proto.Group{}, *res)
	// assertGroupsSame(t, grp1, res)

	resp := listGroups(t)
	assertResponseContainsGroup(t, resp, grp1)

	cleanUp(t)
}

func TestRemoveMemberNonExistingUser(t *testing.T) {
	grp1 := getTestGroups("grp1")

	createGroup(t, grp1)

	client := service.Client()
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
		assert.Equal(
			t,
			fmt.Sprintf("{\"id\":\".\",\"code\":404,\"detail\":\"could not read account: open accounts-store/accounts/%v: no such file or directory\",\"status\":\"Not Found\"}", id),
			err.Error(),
		)
	}

	// Check group is not changed
	resp := listGroups(t)
	assertResponseContainsGroup(t, resp, grp1)
	cleanUp(t)
}

// https://github.com/owncloud/ocis-accounts/issues/62
func TestRemoveMemberNotInGroup(t *testing.T) {
	grp1 := getTestGroups("grp1")
	account := getAccount("user1")

	createGroup(t, grp1)
	createAccount(t, account.PreferredName)

	client := service.Client()
	cl := proto.NewGroupsService("com.owncloud.api.accounts", client)

	req := &proto.RemoveMemberRequest{GroupId: grp1.Id, AccountId: account.Id}

	res, err := cl.RemoveMember(context.Background(), req)

	// Should give an error
	checkNoError(t, err)
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

	client := service.Client()
	cl := proto.NewGroupsService("com.owncloud.api.accounts", client)

	for _, group := range groups {
		expectedGroup := getGroup(group)
		req := &proto.ListMembersRequest{Id: expectedGroup.Id}

		res, err := cl.ListMembers(context.Background(), req)
		checkNoError(t, err)

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
	group := &proto.Group{Id: "5d58e5ec-842e-498b-8800-61f2ec6f911c", GidNumber: 30002, OnPremisesSamAccountName: "quantum-group", DisplayName: "Quantum Group", Members: []*proto.Account{}}

	createGroup(t, group)

	client := service.Client()
	cl := proto.NewGroupsService("com.owncloud.api.accounts", client)

	req := &proto.ListMembersRequest{Id: group.Id}

	res, err := cl.ListMembers(context.Background(), req)

	checkNoError(t, err)
	assert.Empty(t, res.Members)

	cleanUp(t)
}

func TestAccountUpdateMask(t *testing.T) {
	createAccount(t, "user1")
	user1 := getAccount("user1")
	client := service.Client()
	req := &proto.UpdateAccountRequest{
		// We only want to update the display-name, rest should be ignored
		UpdateMask: &field_mask.FieldMask{Paths: []string{"DisplayName"}},
		Account: &proto.Account{
			Id:            user1.Id,
			DisplayName:   "ShouldBeUpdated",
			PreferredName: "ShouldStaySame",
		}}

	cl := proto.NewAccountsService("com.owncloud.api.accounts", client)
	res, err := cl.UpdateAccount(context.Background(), req)
	checkNoError(t, err)

	assert.Equal(t, "ShouldBeUpdated", res.DisplayName)
	assert.Equal(t, user1.PreferredName, res.PreferredName)

	cleanUp(t)
}

func TestAccountUpdateReadOnlyField(t *testing.T) {
	createAccount(t, "user1")
	user1 := getAccount("user1")
	client := service.Client()
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
