package http

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"strings"
	"testing"

	"github.com/asim/go-micro/v3/client"
	user "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	types "github.com/cs3org/go-cs3apis/cs3/types/v1beta1"
	"github.com/cs3org/reva/pkg/auth/scope"
	"github.com/cs3org/reva/pkg/token"
	"github.com/cs3org/reva/pkg/token/manager/jwt"
	"github.com/golang/protobuf/ptypes/empty"
	accountsCmd "github.com/owncloud/ocis/accounts/pkg/command"
	accountsCfg "github.com/owncloud/ocis/accounts/pkg/config"
	accountsProto "github.com/owncloud/ocis/accounts/pkg/proto/v0"
	accountsSvc "github.com/owncloud/ocis/accounts/pkg/service/v0"
	ocisLog "github.com/owncloud/ocis/ocis-pkg/log"
	"github.com/owncloud/ocis/ocis-pkg/service/grpc"
	"github.com/owncloud/ocis/ocs/pkg/config"
	svc "github.com/owncloud/ocis/ocs/pkg/service/v0"
	settings "github.com/owncloud/ocis/settings/pkg/proto/v0"
	ssvc "github.com/owncloud/ocis/settings/pkg/service/v0"
	"github.com/stretchr/testify/assert"
)

const (
	ocsV1 string = "v1.php"
	ocsV2 string = "v2.php"
)

const unsuccessfulResponseText string = "The response was expected to be successful but was not"

const (
	userProvisioningEndPoint  string = "/v1.php/cloud/users?format=json"
	groupProvisioningEndPoint string = "/v1.php/cloud/groups?format=json"
)

const (
	userEinstein string = "einstein"
	userMarie    string = "marie"
	userRichard  string = "richard"
	userIDP      string = "idp"
	userReva     string = "reva"
	userMoss     string = "moss"
	userAdmin    string = "admin"
)
const (
	groupPhilosophyHaters string = "philosophy-haters"
	groupPhysicsLovers    string = "physics-lovers"
	groupPoloniumLovers   string = "polonium-lovers"
	groupQuantumLovers    string = "quantum-lovers"
	groupRadiumLovers     string = "radium-lovers"
	groupSailingLovers    string = "sailing-lovers"
	groupViolinHaters     string = "violin-haters"
	groupUsers            string = "users"
	groupSysUsers         string = "sysusers"
)

var defaultMemberOf = map[string][]string{
	userEinstein: {
		groupUsers,
		groupSailingLovers,
		groupViolinHaters,
		groupPhysicsLovers,
	},
	userIDP: {
		groupSysUsers,
	},
	userRichard: {
		groupUsers,
		groupQuantumLovers,
		groupPhilosophyHaters,
		groupPhysicsLovers,
	},
	userReva: {
		groupSysUsers,
	},
	userMarie: {
		groupUsers,
		groupRadiumLovers,
		groupPoloniumLovers,
		groupPhysicsLovers,
	},
	userMoss: {
		groupUsers,
	},
	userAdmin: {
		groupUsers,
	},
}

var defaultMembers = map[string][]string{
	groupSysUsers: {
		userIDP,
		userReva,
	},
	groupUsers: {
		userEinstein,
		userMarie,
		userRichard,
	},
	groupSailingLovers: {
		userEinstein,
	},
	groupViolinHaters: {
		userEinstein,
	},
	groupPoloniumLovers: {
		userMarie,
	},
	groupQuantumLovers: {
		userRichard,
	},
	groupPhilosophyHaters: {
		userRichard,
	},
	groupPhysicsLovers: {
		userEinstein,
		userMarie,
		userRichard,
	},
}

// These account ids are only needed for cleanup
const (
	userIDEinstein string = "4c510ada-c86b-4815-8820-42cdf82c3d51"
	userIDMarie    string = "f7fbf8c8-139b-4376-b307-cf0a8c2d0d9c"
	userIDFeynman  string = "932b4540-8d16-481e-8ef4-588e4b6b151c"
	userIDIDP      string = "820ba2a1-3f54-4538-80a4-2d73007e30bf"
	userIDReva     string = "bc596f3c-c955-4328-80a0-60d018b4ad57"
	userIDMoss     string = "058bff95-6708-4fe5-91e4-9ea3d377588b"
	userIDAdmin    string = "ddc2004c-0977-11eb-9d3f-a793888cd0f8"
)

// These group ids are only needed for cleanup
const (
	groupIDPhilosophyHaters = "167cbee2-0518-455a-bfb2-031fe0621e5d"
	groupIDPhysicsLovers    = "262982c1-2362-4afa-bfdf-8cbfef64a06e"
	groupIDPoloniumLovers   = "cedc21aa-4072-4614-8676-fa9165f598ff"
	groupIDQuantumLovers    = "a1726108-01f8-4c30-88df-2b1a9d1cba1a"
	groupIDRadiumLovers     = "7b87fd49-286e-4a5f-bafd-c535d5dd997a"
	groupIDSailingLovers    = "6040aa17-9c64-4fef-9bd0-77234d71bad0"
	groupIDViolinHaters     = "dd58e5ec-842e-498b-8800-61f2ec6f911f"
	groupIDUsers            = "509a9dcd-bb37-4f4f-a01a-19dca27d9cfa"
	groupIDSysUsers         = "34f38767-c937-4eb6-b847-1c175829a2a0"
)

const jwtSecret = "HELLO-secret"

var service = grpc.Service{}
var tokenManager token.Manager

var mockedRoleAssignment = map[string]string{}

var ocsVersions = []string{ocsV1, ocsV2}

var formats = []string{"json", "xml"}

var dataPath = createTmpDir()

var defaultUsers = []string{
	userEinstein,
	userIDP,
	userRichard,
	userReva,
	userMarie,
	userMoss,
	userAdmin,
}
var defaultUserIDs = []string{
	userIDEinstein,
	userIDIDP,
	userIDFeynman,
	userIDReva,
	userIDMarie,
	userIDMoss,
	userIDAdmin,
}

var defaultGroups = []string{
	groupPhilosophyHaters,
	groupPhysicsLovers,
	groupSysUsers,
	groupUsers,
	groupSailingLovers,
	groupRadiumLovers,
	groupQuantumLovers,
	groupPoloniumLovers,
	groupViolinHaters,
}
var defaultGroupIDs = []string{
	groupIDPhilosophyHaters,
	groupIDPhysicsLovers,
	groupIDSysUsers,
	groupIDUsers,
	groupIDSailingLovers,
	groupIDRadiumLovers,
	groupIDQuantumLovers,
	groupIDPoloniumLovers,
	groupIDViolinHaters,
}

func getFormatString(format string) string {
	if format == "json" {
		return "?format=json"
	} else if format == "xml" {
		return ""
	} else {
		panic("Invalid format received")
	}
}

func createTmpDir() string {
	name, err := ioutil.TempDir("/var/tmp", "ocis-accounts-store-")
	if err != nil {
		panic(err)
	}

	return name
}

type Quota struct {
	Free       int64   `json:"free" xml:"free"`
	Used       int64   `json:"used" xml:"used"`
	Total      int64   `json:"total" xml:"total"`
	Relative   float32 `json:"relative" xml:"relative"`
	Definition string  `json:"definition" xml:"definition"`
}

type User struct {
	Enabled     string `json:"enabled" xml:"enabled"`
	ID          string `json:"id" xml:"id"`
	Email       string `json:"email" xml:"email"`
	Password    string `json:"-" xml:"-"`
	Quota       Quota  `json:"quota" xml:"quota"`
	UIDNumber   int    `json:"uidnumber" xml:"uidnumber"`
	GIDNumber   int    `json:"gidnumber" xml:"gidnumber"`
	Displayname string `json:"displayname" xml:"displayname"`
}

func (u *User) getUserRequestString() string {
	res := url.Values{}

	if u.Password != "" {
		res.Add("password", u.Password)
	}

	if u.ID != "" {
		res.Add("userid", u.ID)
	}

	if u.Email != "" {
		res.Add("email", u.Email)
	}

	if u.Displayname != "" {
		res.Add("displayname", u.Displayname)
	}

	if u.UIDNumber != 0 {
		res.Add("uidnumber", fmt.Sprint(u.UIDNumber))
	}

	if u.GIDNumber != 0 {
		res.Add("gidnumber", fmt.Sprint(u.GIDNumber))
	}

	return res.Encode()
}

type Group struct {
	ID          string `json:"id" xml:"id"`
	GIDNumber   int    `json:"gidnumber" xml:"gidnumber"`
	Displayname string `json:"displayname" xml:"displayname"`
}

func (g *Group) getGroupRequestString() string {
	res := url.Values{}

	if g.ID != "" {
		res.Add("groupid", g.ID)
	}

	if g.Displayname != "" {
		res.Add("displayname", g.Displayname)
	}
	if g.GIDNumber != 0 {
		res.Add("gidnumber", fmt.Sprint(g.GIDNumber))
	}

	return res.Encode()
}

type Meta struct {
	Status     string `json:"status" xml:"status"`
	StatusCode int    `json:"statuscode" xml:"statuscode"`
	Message    string `json:"message" xml:"message"`
}

func (m *Meta) Success(ocsVersion string) bool {
	if !(ocsVersion == ocsV1 || ocsVersion == ocsV2) {
		return false
	}
	if m.Status != "ok" {
		return false
	}
	if ocsVersion == ocsV1 && m.StatusCode != 100 {
		return false
	} else if ocsVersion == ocsV2 && m.StatusCode != 200 {
		return false
	} else {
		return true
	}
}

type SingleUserResponse struct {
	Ocs struct {
		Meta Meta `json:"meta" xml:"meta"`
		Data User `json:"data" xml:"data"`
	} `json:"ocs" xml:"ocs"`
}

type GetUsersResponse struct {
	Ocs struct {
		Meta Meta `json:"meta" xml:"meta"`
		Data struct {
			Users []string `json:"users" xml:"users>element"`
		} `json:"data" xml:"data"`
	} `json:"ocs" xml:"ocs"`
}

type EmptyResponse struct {
	Ocs struct {
		Meta Meta `json:"meta" xml:"meta"`
		Data struct {
		} `json:"data" xml:"data"`
	} `json:"ocs" xml:"ocs"`
}

func assertEmptyResponse(t *testing.T, format string, res *httptest.ResponseRecorder) *EmptyResponse {
	var response EmptyResponse
	if format == "json" {
		if err := json.Unmarshal(res.Body.Bytes(), &response); err != nil {
			t.Log(res.Body.String())
			t.Fatal(err)
		}
	} else {
		if err := xml.Unmarshal(res.Body.Bytes(), &response.Ocs); err != nil {
			t.Log(res.Body.String())
			t.Fatal(err)
		}
	}
	return &response
}

type GetUsersGroupsResponse struct {
	Ocs struct {
		Meta Meta `json:"meta" xml:"meta"`
		Data struct {
			Groups []string `json:"groups" xml:"groups>element"`
		} `json:"data" xml:"data"`
	} `json:"ocs" xml:"ocs"`
}

type OcsConfig struct {
	Version string `json:"version" xml:"version"`
	Website string `json:"website" xml:"website"`
	Host    string `json:"host" xml:"host"`
	Contact string `json:"contact" xml:"contact"`
	Ssl     string `json:"ssl" xml:"ssl"`
}

type GetConfigResponse struct {
	Ocs struct {
		Meta Meta      `json:"meta" xml:"meta"`
		Data OcsConfig `json:"data" xml:"data"`
	} `json:"ocs" xml:"ocs"`
}

func assertStatusCode(t *testing.T, statusCode int, res *httptest.ResponseRecorder, ocsVersion string) {
	if ocsVersion == ocsV1 {
		assert.Equal(t, 200, res.Code)
	} else {
		assert.Equal(t, statusCode, res.Code)
	}
}

type GetGroupsResponse struct {
	Ocs struct {
		Meta Meta `json:"meta" xml:"meta"`
		Data struct {
			Groups []string `json:"groups" xml:"groups>element"`
		} `json:"data" xml:"data"`
	} `json:"ocs" xml:"ocs"`
}

type GetGroupMembersResponse struct {
	Ocs struct {
		Meta Meta `json:"meta" xml:"meta"`
		Data struct {
			Users []string `json:"users" xml:"users>element"`
		} `json:"data" xml:"data"`
	} `json:"ocs" xml:"ocs"`
}

func assertResponseMeta(t *testing.T, expected, actual Meta) {
	assert.Equal(t, expected.Status, actual.Status, "The status of response doesn't match")
	assert.Equal(t, expected.StatusCode, actual.StatusCode, "The Status code of response doesn't match")
	assert.Equal(t, expected.Message, actual.Message, "The Message of response doesn't match")
}

// compares users at tha /user endpoint
func assertUserSame(t *testing.T, expected, actual User) {
	if expected.ID == "" {
		// Check the auto generated userId
		assert.Regexp(
			t,
			"[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}",
			actual.ID, "the userid is not a valid uuid",
		)
	} else {
		assert.Equal(t, expected.ID, actual.ID, "UserId doesn't match for user %v", expected.ID)
	}
	assert.Equal(t, expected.Email, actual.Email, "email doesn't match for user %v", expected.ID)
	// /user has no enabled flag
	if expected.Displayname == "" {
		assert.Equal(t, expected.ID, actual.Displayname, "displayname doesn't match for user %v", expected.ID)
	} else {
		assert.Equal(t, expected.Displayname, actual.Displayname, "displayname doesn't match for user %v", expected.ID)
	}
	// /user has no quota, uid or gid
}

// compares users at the /users/<userid> endpoint
func assertUsersSame(t *testing.T, expected, actual User, quotaAvailable bool) {
	if expected.ID == "" {
		// Check the auto generated userId
		assert.Regexp(
			t,
			"[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}",
			actual.ID, "the userid is not a valid uuid",
		)
	} else {
		assert.Equal(t, expected.ID, actual.ID, "UserId doesn't match for user %v", expected.ID)
	}
	assert.Equal(t, expected.Email, actual.Email, "email doesn't match for user %v", expected.ID)
	assert.Equal(t, expected.Enabled, actual.Enabled, "enabled doesn't match for user %v", expected.ID)
	if expected.Displayname == "" {
		assert.Equal(t, expected.ID, actual.Displayname, "displayname doesn't match for user %v", expected.ID)
	} else {
		assert.Equal(t, expected.Displayname, actual.Displayname, "displayname doesn't match for user %v", expected.ID)
	}
	if quotaAvailable {
		assert.NotZero(t, actual.Quota.Free)
		assert.NotZero(t, actual.Quota.Used)
		assert.NotZero(t, actual.Quota.Total)
		assert.Equal(t, "default", actual.Quota.Definition)
	} else {
		assert.Equal(t, expected.Quota, actual.Quota, "Quota match for user %v", expected.ID)
	}

	if expected.UIDNumber != 0 {
		assert.Equal(t, expected.UIDNumber, actual.UIDNumber, "UidNumber doesn't match for user %s", expected.ID)
	}
	if expected.GIDNumber != 0 {
		assert.Equal(t, expected.GIDNumber, actual.GIDNumber, "GidNumber doesn't match for user %s", expected.ID)
	}
}

func deleteAccount(t *testing.T, id string) (*empty.Empty, error) {
	cl := accountsProto.NewAccountsService("com.owncloud.api.accounts", service.Client())

	req := &accountsProto.DeleteAccountRequest{Id: id}
	res, err := cl.DeleteAccount(context.Background(), req)
	return res, err
}

func deleteGroup(t *testing.T, id string) (*empty.Empty, error) {
	cl := accountsProto.NewGroupsService("com.owncloud.api.accounts", service.Client())

	req := &accountsProto.DeleteGroupRequest{Id: id}
	res, err := cl.DeleteGroup(context.Background(), req)
	return res, err
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
		ListRolesFunc: func(ctx context.Context, req *settings.ListBundlesRequest, opts ...client.CallOption) (*settings.ListBundlesResponse, error) {
			return &settings.ListBundlesResponse{
				Bundles: []*settings.Bundle{
					{
						Id: ssvc.BundleUUIDRoleAdmin,
						Settings: []*settings.Setting{
							{
								Id: accountsSvc.AccountManagementPermissionID,
							},
						},
					},
					{
						Id: ssvc.BundleUUIDRoleUser,
						Settings: []*settings.Setting{
							{
								Id: accountsSvc.SelfManagementPermissionID,
							},
						},
					},
				},
			}, nil
		},
	}
}

func init() {
	service = grpc.NewService(
		grpc.Namespace("com.owncloud.api"),
		grpc.Name("accounts"),
		grpc.Address("localhost:9180"),
	)

	c := &accountsCfg.Config{
		Server: accountsCfg.Server{},
		Repo: accountsCfg.Repo{
			Disk: accountsCfg.Disk{
				Path: dataPath,
			},
		},
		Log: accountsCfg.Log{
			Level:  "info",
			Pretty: true,
			Color:  true,
		},
	}

	var hdlr *accountsSvc.Service
	var err error

	if hdlr, err = accountsSvc.New(
		accountsSvc.Logger(accountsCmd.NewLogger(c)),
		accountsSvc.Config(c),
		accountsSvc.RoleService(buildRoleServiceMock()),
	); err != nil {
		log.Fatalf("Could not create new service")
	}

	err = accountsProto.RegisterAccountsServiceHandler(service.Server(), hdlr)
	if err != nil {
		log.Fatal("could not register the Accounts handler")
	}
	err = accountsProto.RegisterGroupsServiceHandler(service.Server(), hdlr)
	if err != nil {
		log.Fatal("could not register the Groups handler")
	}

	err = service.Server().Start()
	if err != nil {
		log.Fatalf("could not start server: %v", err)
	}

	// a token manager to mint tokens
	tokenManager, err = jwt.New(map[string]interface{}{
		"secret": jwtSecret,
	})
	if err != nil {
		log.Fatalf("could not create token manager: %v", err)
	}
}

func cleanUp(t *testing.T) {
	datastore := filepath.Join(dataPath, "accounts")

	files, err := ioutil.ReadDir(datastore)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		found := false
		for _, defUser := range defaultUserIDs {
			if f.Name() == defUser {
				found = true
				break
			}
		}

		if !found {
			if _, err := deleteAccount(t, f.Name()); err != nil {
				panic(err)
			}
		}
	}

	datastoreGroups := filepath.Join(dataPath, "groups")

	files, err = ioutil.ReadDir(datastoreGroups)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		found := false
		for _, defGrp := range defaultGroupIDs {
			if f.Name() == defGrp {
				found = true
				break
			}
		}

		if !found {
			if _, err := deleteGroup(t, f.Name()); err != nil {
				panic(err)
			}
		}
	}
}

func mintToken(ctx context.Context, su *User, roleIds []string) (token string, err error) {
	roleIDsJSON, err := json.Marshal(roleIds)
	if err != nil {
		return "", err
	}
	u := &user.User{
		Id: &user.UserId{
			OpaqueId: su.ID,
		},
		Opaque: &types.Opaque{
			Map: map[string]*types.OpaqueEntry{
				"roles": {
					Decoder: "json",
					Value:   roleIDsJSON,
				},
			},
		},
		Groups:    []string{},
		UidNumber: int64(su.UIDNumber),
		GidNumber: int64(su.GIDNumber),
	}
	s, _ := scope.GetOwnerScope()
	return tokenManager.MintToken(ctx, u, s)
}

func sendRequest(method, endpoint, body string, u *User, roleIds []string) (*httptest.ResponseRecorder, error) {
	var reader = strings.NewReader(body)
	req, err := http.NewRequest(method, endpoint, reader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	if u != nil {
		token, err := mintToken(context.Background(), u, roleIds)
		if err != nil {
			return nil, err
		}
		req.Header.Set("x-access-token", token)
	}

	rr := httptest.NewRecorder()

	service := getService()
	service.ServeHTTP(rr, req)

	return rr, nil
}

func getService() svc.Service {
	c := &config.Config{
		HTTP: config.HTTP{
			Root: "/",
			Addr: "localhost:9110",
		},
		TokenManager: config.TokenManager{
			JWTSecret: jwtSecret,
		},
		Log: config.Log{
			Level: "debug",
		},
	}

	var logger ocisLog.Logger

	return svc.NewService(
		svc.Logger(logger),
		svc.Config(c),
		svc.RoleService(buildRoleServiceMock()),
	)
}

func createUser(u User) error {
	_, err := sendRequest(
		"POST",
		userProvisioningEndPoint,
		u.getUserRequestString(),
		&User{ID: userIDAdmin},
		[]string{ssvc.BundleUUIDRoleAdmin},
	)

	if err != nil {
		return err
	}
	return nil
}

func createGroup(g Group) error { //lint:file-ignore U1000 not implemented
	_, err := sendRequest(
		"POST",
		groupProvisioningEndPoint,
		g.getGroupRequestString(),
		&User{ID: userIDAdmin},
		[]string{ssvc.BundleUUIDRoleAdmin},
	)

	if err != nil {
		return err
	}
	return nil
}

func TestCreateUser(t *testing.T) {
	scenarios := []struct {
		name string
		user User
		err  *Meta
	}{
		{
			"A simple user",
			User{
				Enabled:     "true",
				ID:          "rutherford",
				Email:       "rutherford@example.com",
				Displayname: "ErnestRutherFord",
				Password:    "newPassword",
			},
			nil,
		},
		{
			"User with Uid and Gid defined",
			User{
				Enabled:     "true",
				ID:          "thomson",
				Email:       "thomson@example.com",
				Displayname: "J. J. Thomson",
				UIDNumber:   20027,
				GIDNumber:   30000,
				Password:    "newPassword",
			},
			nil,
		},
		// https://github.com/owncloud/ocis-ocs/issues/50
		{
			"User without password",
			User{
				Enabled:     "true",
				ID:          "john",
				Email:       "john@example.com",
				Displayname: "John Dalton",
			},
			nil,
		},
		// https://github.com/owncloud/ocis-ocs/issues/49
		{
			"User with special character in userid",
			User{
				Enabled:     "true",
				ID:          "schrödinger",
				Email:       "schrödinger@example.com",
				Displayname: "Erwin Schrödinger",
				Password:    "newPassword",
			},
			&Meta{
				Status:     "error",
				StatusCode: 400,
				Message:    "preferred_name 'schrödinger' must be at least the local part of an email",
			},
		},
		{
			"User with different userid and email",
			User{
				Enabled:     "true",
				ID:          "planck",
				Email:       "max@example.com",
				Displayname: "Max Planck",
				Password:    "newPassword",
			},
			nil,
		},
		{
			"User without displayname",
			User{
				Enabled:  "true",
				ID:       "oppenheimer",
				Email:    "robert@example.com",
				Password: "newPassword",
			},
			nil,
		},
		{
			"User wit invalid email",
			User{
				Enabled:  "true",
				ID:       "chadwick",
				Email:    "not_a_email",
				Password: "newPassword",
			},
			&Meta{
				Status:     "error",
				StatusCode: 400,
				Message:    "mail 'not_a_email' must be a valid email",
			},
		},
		{
			"User without email",
			User{
				Enabled:  "true",
				ID:       "chadwick",
				Password: "newPassword",
			},
			&Meta{
				Status:     "error",
				StatusCode: 400,
				Message:    "mail '' must be a valid email",
			},
		},
		{
			"User without userid",
			User{
				Enabled:  "true",
				Email:    "james@example.com",
				Password: "newPassword",
			},
			&Meta{
				Status:     "error",
				StatusCode: 400,
				Message:    "preferred_name '' must be at least the local part of an email",
			},
		},
	}

	for _, ocsVersion := range ocsVersions {
		for _, format := range formats {
			for _, scenario := range scenarios {
				t.Run(fmt.Sprintf("%s (ocs=%s, format=%s)", scenario.name, ocsVersion, format), func(t *testing.T) {
					formatpart := getFormatString(format)
					res, err := sendRequest(
						"POST",
						fmt.Sprintf("/%v/cloud/users%v", ocsVersion, formatpart),
						scenario.user.getUserRequestString(),
						&User{ID: userIDAdmin},
						[]string{ssvc.BundleUUIDRoleAdmin},
					)
					assert.NoError(t, err)

					var response SingleUserResponse
					if format == "json" {
						err = json.Unmarshal(res.Body.Bytes(), &response)
						assert.NoError(t, err)
					} else {
						err = xml.Unmarshal(res.Body.Bytes(), &response.Ocs)
						assert.NoError(t, err)
					}

					if scenario.err == nil {
						assert.True(t, response.Ocs.Meta.Success(ocsVersion), unsuccessfulResponseText)
						assertStatusCode(t, 200, res, ocsVersion)
						assertUsersSame(t, scenario.user, response.Ocs.Data, false)
					} else {
						assertStatusCode(t, 400, res, ocsVersion)
						assertResponseMeta(t, *scenario.err, response.Ocs.Meta)
					}

					var id string
					if scenario.user.ID != "" {
						id = scenario.user.ID
					} else {
						id = response.Ocs.Data.ID
					}

					res, err = sendRequest(
						"GET",
						userProvisioningEndPoint,
						"",
						&User{ID: userIDAdmin},
						[]string{ssvc.BundleUUIDRoleAdmin},
					)
					assert.NoError(t, err)

					var usersResponse GetUsersResponse
					err = json.Unmarshal(res.Body.Bytes(), &usersResponse)
					assert.NoError(t, err)

					assert.True(t, usersResponse.Ocs.Meta.Success(ocsV1), unsuccessfulResponseText)

					if scenario.err == nil {
						assert.Contains(t, usersResponse.Ocs.Data.Users, id)
					} else {
						assert.NotContains(t, usersResponse.Ocs.Data.Users, scenario.user.ID)
					}
					cleanUp(t)
				})
			}
		}
	}
}

func TestGetUsers(t *testing.T) {
	users := []User{
		{
			Enabled:     "true",
			ID:          "rutherford",
			Email:       "rutherford@example.com",
			Displayname: "Ernest RutherFord",
		},
		{
			Enabled:     "true",
			ID:          "thomson",
			Email:       "thomson@example.com",
			Displayname: "J. J. Thomson",
		},
	}

	for _, ocsVersion := range ocsVersions {
		for _, format := range formats {
			for _, user := range users {
				err := createUser(user)
				if err != nil {
					t.Fatal(err)
				}
			}

			formatpart := getFormatString(format)
			res, err := sendRequest(
				"GET",
				fmt.Sprintf("/%v/cloud/users%v", ocsVersion, formatpart),
				"",
				&User{ID: userIDAdmin},
				[]string{ssvc.BundleUUIDRoleAdmin},
			)

			if err != nil {
				t.Fatal(err)
			}

			var response GetUsersResponse

			if format == "json" {
				if err := json.Unmarshal(res.Body.Bytes(), &response); err != nil {
					t.Fatal(err)
				}
			} else {
				if err := xml.Unmarshal(res.Body.Bytes(), &response.Ocs); err != nil {
					t.Fatal(err)
				}
			}

			assertStatusCode(t, 200, res, ocsVersion)
			assert.True(t, response.Ocs.Meta.Success(ocsVersion), unsuccessfulResponseText)
			for _, user := range users {
				assert.Contains(t, response.Ocs.Data.Users, user.ID)
			}
			cleanUp(t)
		}
	}
}

func TestGetUsersDefaultUsers(t *testing.T) {
	for _, ocsVersion := range ocsVersions {
		for _, format := range formats {
			formatpart := getFormatString(format)
			res, err := sendRequest(
				"GET",
				fmt.Sprintf("/%v/cloud/users%v", ocsVersion, formatpart),
				"",
				&User{ID: userIDAdmin},
				[]string{ssvc.BundleUUIDRoleAdmin},
			)

			if err != nil {
				t.Fatal(err)
			}

			var response GetUsersResponse

			if format == "json" {
				if err := json.Unmarshal(res.Body.Bytes(), &response); err != nil {
					t.Fatal(err)
				}
			} else {
				if err := xml.Unmarshal(res.Body.Bytes(), &response.Ocs); err != nil {
					t.Fatal(err)
				}
			}

			assertStatusCode(t, 200, res, ocsVersion)
			assert.True(t, response.Ocs.Meta.Success(ocsVersion), unsuccessfulResponseText)
			for _, user := range defaultUsers {
				assert.Contains(t, response.Ocs.Data.Users, user)
			}
			cleanUp(t)
		}
	}
}

func TestGetUser(t *testing.T) {
	users := []User{
		{
			Enabled:     "true",
			ID:          "rutherford",
			Email:       "rutherford@example.com",
			Displayname: "Ernest RutherFord",
		},
		{
			Enabled:     "true",
			ID:          "thomson",
			Email:       "thomson@example.com",
			Displayname: "J. J. Thomson",
		},
	}

	for _, ocsVersion := range ocsVersions {
		for _, format := range formats {

			for _, user := range users {
				err := createUser(user)
				if err != nil {
					t.Fatal(err)
				}
			}
			formatpart := getFormatString(format)
			for _, user := range users {
				res, err := sendRequest(
					"GET",
					fmt.Sprintf("/%s/cloud/users/%s%s", ocsVersion, user.ID, formatpart),
					"",
					&User{ID: userIDAdmin},
					[]string{ssvc.BundleUUIDRoleAdmin},
				)

				if err != nil {
					t.Fatal(err)
				}

				var response SingleUserResponse
				if format == "json" {
					if err := json.Unmarshal(res.Body.Bytes(), &response); err != nil {
						t.Fatal(err)
					}
				} else {
					if err := xml.Unmarshal(res.Body.Bytes(), &response.Ocs); err != nil {
						t.Fatal(err)
					}
				}

				assertStatusCode(t, 200, res, ocsVersion)
				assert.True(t, response.Ocs.Meta.Success(ocsVersion), "The response was expected to pass but it failed")
				assertUsersSame(t, user, response.Ocs.Data, true)
			}
			cleanUp(t)
		}
	}
}

func TestGetUserInvalidId(t *testing.T) {
	invalidUsers := []string{
		"1",
		"invalid",
		"3434234233",
		"1am41validUs3r",
		"_-@@--$$__",
	}
	for _, ocsVersion := range ocsVersions {
		for _, format := range formats {
			formatpart := getFormatString(format)
			for _, user := range invalidUsers {
				res, err := sendRequest(
					"GET",
					fmt.Sprintf("/%s/cloud/user/%s%s", ocsVersion, user, formatpart),
					"",
					&User{ID: userIDAdmin},
					[]string{ssvc.BundleUUIDRoleAdmin},
				)

				if err != nil {
					t.Fatal(err)
				}

				var response SingleUserResponse
				if format == "json" {
					if err := json.Unmarshal(res.Body.Bytes(), &response); err != nil {
						t.Fatal(err)
					}
				} else {
					if err := xml.Unmarshal(res.Body.Bytes(), &response.Ocs); err != nil {
						t.Fatal(err)
					}
				}

				assertStatusCode(t, 404, res, ocsVersion)
				assert.False(t, response.Ocs.Meta.Success(ocsVersion), "the response was expected to fail but passed")
				assertResponseMeta(t, Meta{
					Status:     "error",
					StatusCode: 998,
					Message:    "not found",
				}, response.Ocs.Meta)
				cleanUp(t)
			}
		}
	}
}
func TestDeleteUser(t *testing.T) {
	users := []User{
		{
			Enabled:     "true",
			ID:          "rutherford",
			Email:       "rutherford@example.com",
			Displayname: "Ernest RutherFord",
		},
		{
			Enabled:     "true",
			ID:          "thomson",
			Email:       "thomson@example.com",
			Displayname: "J. J. Thomson",
		},
	}

	for _, ocsVersion := range ocsVersions {
		for _, format := range formats {
			for _, user := range users {
				err := createUser(user)
				if err != nil {
					t.Fatal(err)
				}
			}

			formatpart := getFormatString(format)
			res, err := sendRequest(
				"DELETE",
				fmt.Sprintf("/%s/cloud/users/rutherford%s", ocsVersion, formatpart),
				"",
				&User{ID: userIDAdmin},
				[]string{ssvc.BundleUUIDRoleAdmin},
			)

			if err != nil {
				t.Fatal(err)
			}

			response := assertEmptyResponse(t, format, res)

			assertStatusCode(t, 200, res, ocsVersion)
			assert.True(t, response.Ocs.Meta.Success(ocsVersion), unsuccessfulResponseText)
			assert.Empty(t, response.Ocs.Data)

			// Check deleted user doesn't exist and the other user does
			res, err = sendRequest(
				"GET",
				userProvisioningEndPoint,
				"",
				&User{ID: userIDAdmin},
				[]string{ssvc.BundleUUIDRoleAdmin},
			)

			if err != nil {
				t.Fatal(err)
			}

			var usersResponse GetUsersResponse
			if err := json.Unmarshal(res.Body.Bytes(), &usersResponse); err != nil {
				t.Fatal(err)
			}

			assert.True(t, usersResponse.Ocs.Meta.Success(ocsV1), unsuccessfulResponseText)
			assert.Contains(t, usersResponse.Ocs.Data.Users, "thomson")
			assert.NotContains(t, usersResponse.Ocs.Data.Users, "rutherford")

			cleanUp(t)
		}
	}
}

func TestDeleteUserInvalidId(t *testing.T) {

	invalidUsers := []string{
		"1",
		"invalid",
		"3434234233",
		"1am41validUs3r",
		"_-@@--$$__",
	}
	for _, ocsVersion := range ocsVersions {
		for _, format := range formats {
			for _, user := range invalidUsers {
				t.Run(fmt.Sprintf("%s (ocs=%s, format=%s)", user, ocsVersion, format), func(t *testing.T) {
					formatpart := getFormatString(format)
					res, err := sendRequest(
						"DELETE",
						fmt.Sprintf("/%s/cloud/users/%s%s", ocsVersion, user, formatpart),
						"",
						&User{ID: userIDAdmin},
						[]string{ssvc.BundleUUIDRoleAdmin},
					)
					assert.NoError(t, err)

					response := assertEmptyResponse(t, format, res)

					assertStatusCode(t, 404, res, ocsVersion)
					assert.False(t, response.Ocs.Meta.Success(ocsVersion), "The response was not expected to be successful but was")
					assert.Empty(t, response.Ocs.Data)

					assertResponseMeta(t, Meta{
						Status:     "error",
						StatusCode: 998,
						Message:    "The requested user could not be found",
					}, response.Ocs.Meta)
				})
			}
		}
	}
}

func TestUpdateUser(t *testing.T) {
	user := User{
		Enabled:     "true",
		ID:          "rutherford",
		Email:       "rutherford@example.com",
		Displayname: "Ernest RutherFord",
	}

	testData := []struct {
		UpdateKey   string
		UpdateValue string
		Error       *Meta
	}{
		{
			"displayname",
			"James Chadwick",
			nil,
		},
		{
			"display",
			"Neils Bohr",
			nil,
		},
		{
			"email",
			"ford@user.org",
			nil,
		},
		{
			"email",
			"not_a_valid_email",
			&Meta{
				Status:     "error",
				StatusCode: 400,
				Message:    "mail 'not_a_valid_email' must be a valid email",
			},
		},
		{
			"password",
			"strongpass1234",
			nil,
		},
		{
			"email",
			"",
			nil,
		},
		{
			"password",
			"",
			nil,
		},
		// Invalid Keys
		{
			"invalid_key",
			"validvalue",
			&Meta{
				Status:     "error",
				StatusCode: 103,
				Message:    "unknown key 'invalid_key'",
			},
		},
		{
			"12345",
			"validvalue",
			&Meta{
				Status:     "error",
				StatusCode: 103,
				Message:    "unknown key '12345'",
			},
		},
		{
			"",
			"validvalue",
			&Meta{
				Status:     "error",
				StatusCode: 103,
				Message:    "unknown key ''",
			},
		},
		{
			"",
			"",
			&Meta{
				Status:     "error",
				StatusCode: 103,
				Message:    "unknown key ''",
			},
		},
	}

	for _, ocsVersion := range ocsVersions {
		for _, format := range formats {
			formatpart := getFormatString(format)
			for _, data := range testData {
				err := createUser(user)
				if err != nil {
					t.Fatalf("Failed while creating user: %v", err)
				}

				params := url.Values{}

				params.Add("key", data.UpdateKey)
				params.Add("value", data.UpdateValue)

				res, err := sendRequest(
					"PUT",
					fmt.Sprintf("/%s/cloud/users/rutherford%s", ocsVersion, formatpart),
					params.Encode(),
					&User{ID: userIDAdmin},
					[]string{ssvc.BundleUUIDRoleAdmin},
				)

				updatedUser := user
				switch data.UpdateKey {
				case "email":
					updatedUser.Email = data.UpdateValue
				case "displayname":
					updatedUser.Displayname = data.UpdateValue
				case "display":
					updatedUser.Displayname = data.UpdateValue
				}

				if err != nil {
					t.Fatal(err)
				}

				var response struct {
					Ocs struct {
						Meta Meta `json:"meta" xml:"meta"`
					} `json:"ocs" xml:"ocs"`
				}

				if format == "json" {
					if err := json.Unmarshal(res.Body.Bytes(), &response); err != nil {
						t.Fatal(err)
					}
				} else {
					if err := xml.Unmarshal(res.Body.Bytes(), &response.Ocs); err != nil {
						t.Fatal(err)
					}
				}

				if data.Error != nil {
					assertResponseMeta(t, *data.Error, response.Ocs.Meta)
					assertStatusCode(t, 400, res, ocsVersion)
				} else {
					assert.True(t, response.Ocs.Meta.Success(ocsVersion), unsuccessfulResponseText)
					assertStatusCode(t, 200, res, ocsVersion)
				}

				// Check deleted user doesn't exist and the other user does
				res, err = sendRequest(
					"GET",
					"/v1.php/cloud/users/rutherford?format=json",
					"",
					&User{ID: userIDAdmin},
					[]string{ssvc.BundleUUIDRoleAdmin},
				)

				if err != nil {
					t.Fatal(err)
				}

				var usersResponse SingleUserResponse
				if err := json.Unmarshal(res.Body.Bytes(), &usersResponse); err != nil {
					t.Fatal(err)
				}

				assert.True(t, usersResponse.Ocs.Meta.Success(ocsV1), unsuccessfulResponseText)
				if data.Error == nil {
					assertUsersSame(t, updatedUser, usersResponse.Ocs.Data, true)
				} else {
					assertUsersSame(t, user, usersResponse.Ocs.Data, true)
				}
				cleanUp(t)
			}
		}
	}
}

// This is a bug verification test for endpoint '/cloud/user'
// Link to the fixed issue: https://github.com/owncloud/ocis-ocs/issues/52
func TestGetSingleUser(t *testing.T) {
	user := User{
		Enabled:     "true",
		ID:          "rutherford",
		Email:       "rutherford@example.com",
		Displayname: "Ernest RutherFord",
		Password:    "password",
	}

	for _, ocsVersion := range ocsVersions {
		for _, format := range formats {
			err := createUser(user)
			if err != nil {
				t.Fatal(err)
			}

			formatpart := getFormatString(format)
			res, err := sendRequest(
				"GET",
				fmt.Sprintf("/%v/cloud/user%v", ocsVersion, formatpart),
				"",
				&User{ID: user.ID},
				[]string{ssvc.BundleUUIDRoleUser},
			)

			if err != nil {
				t.Fatal(err)
			}

			var userResponse SingleUserResponse
			if format == "json" {
				if err := json.Unmarshal(res.Body.Bytes(), &userResponse); err != nil {
					t.Fatal(err)
				}
			} else {
				if err := xml.Unmarshal(res.Body.Bytes(), &userResponse.Ocs); err != nil {
					t.Fatal(err)
				}
			}

			assertStatusCode(t, 200, res, ocsVersion)
			assert.True(t, userResponse.Ocs.Meta.Success(ocsVersion), "The response was expected to pass but it failed")
			assertUserSame(t, user, userResponse.Ocs.Data)

			cleanUp(t)
		}
	}
}

// This is a bug demonstration test for endpoint '/cloud/user'
// Link to the issue: https://github.com/owncloud/ocis/ocs/issues/53
func TestGetUserSigningKey(t *testing.T) {
	user := User{
		Enabled:     "true",
		ID:          "rutherford",
		Email:       "rutherford@example.com",
		Displayname: "Ernest RutherFord",
		Password:    "password",
	}

	for _, ocsVersion := range ocsVersions {
		for _, format := range formats {
			err := createUser(user)
			if err != nil {
				t.Fatal(err)
			}

			formatpart := getFormatString(format)
			res, err := sendRequest(
				"GET",
				fmt.Sprintf("/%v/cloud/user/signing-key%v", ocsVersion, formatpart),
				"",
				&User{ID: user.ID},
				[]string{ssvc.BundleUUIDRoleUser},
			)

			if err != nil {
				t.Fatal(err)
			}

			response := assertEmptyResponse(t, format, res)

			assertStatusCode(t, 500, res, ocsVersion)
			assert.False(t, response.Ocs.Meta.Success(ocsVersion), "The response was expected to be a failure but was not")
			assertResponseMeta(t, Meta{
				Status:     "error",
				StatusCode: 996,
				Message:    "error reading from store", // because the store service is not started
			}, response.Ocs.Meta)
			assert.Empty(t, response.Ocs.Data)
			cleanUp(t)
		}
	}
}

func AddUserToGroup(userid, groupid string) error {
	res, err := sendRequest(
		"POST",
		fmt.Sprintf("/v2.php/cloud/users/%s/groups", userid),
		fmt.Sprintf("groupid=%v", groupid),
		&User{ID: userIDAdmin},
		[]string{ssvc.BundleUUIDRoleAdmin},
	)
	if err != nil {
		return err
	}
	if res.Code != 200 {
		return fmt.Errorf("Failed while adding the user to group")
	}
	return nil
}

func TestListUsersGroupNewUsers(t *testing.T) {
	users := []User{
		{
			Enabled:     "true",
			ID:          "rutherford",
			Email:       "rutherford@example.com",
			Displayname: "Ernest RutherFord",
		},
		{
			Enabled:     "true",
			ID:          "thomson",
			Email:       "thomson@example.com",
			Displayname: "J. J. Thomson",
		},
	}

	for _, ocsVersion := range ocsVersions {
		for _, format := range formats {
			formatpart := getFormatString(format)
			for _, user := range users {
				err := createUser(user)
				if err != nil {
					t.Fatal(err)
				}

				res, err := sendRequest(
					"GET",
					fmt.Sprintf("/%s/cloud/users/%s/groups%s", ocsVersion, user.ID, formatpart),
					"",
					&User{ID: userIDAdmin},
					[]string{ssvc.BundleUUIDRoleAdmin},
				)

				if err != nil {
					t.Fatal(err)
				}

				var response GetUsersGroupsResponse
				if format == "json" {
					if err := json.Unmarshal(res.Body.Bytes(), &response); err != nil {
						t.Fatal(err)
					}
				} else {
					if err := xml.Unmarshal(res.Body.Bytes(), &response.Ocs); err != nil {
						t.Fatal(err)
					}
				}

				assertStatusCode(t, 200, res, ocsVersion)
				assert.True(t, response.Ocs.Meta.Success(ocsVersion), unsuccessfulResponseText)
				// TODO why should new users be in the users group?
				assert.Equal(t, []string{groupUsers}, response.Ocs.Data.Groups)

				cleanUp(t)
			}
		}
	}
}

func TestListUsersGroupDefaultUsers(t *testing.T) {

	for _, ocsVersion := range ocsVersions {
		for _, format := range formats {
			formatpart := getFormatString(format)
			for _, user := range defaultUsers {
				res, err := sendRequest(
					"GET",
					fmt.Sprintf("/%s/cloud/users/%s/groups%s", ocsVersion, user, formatpart),
					"",
					&User{ID: userIDAdmin},
					[]string{ssvc.BundleUUIDRoleAdmin},
				)

				if err != nil {
					t.Fatal(err)
				}

				var response GetUsersGroupsResponse
				if format == "json" {
					if err := json.Unmarshal(res.Body.Bytes(), &response); err != nil {
						t.Fatal(err)
					}
				} else {
					if err := xml.Unmarshal(res.Body.Bytes(), &response.Ocs); err != nil {
						t.Fatal(err)
					}
				}

				assertStatusCode(t, 200, res, ocsVersion)
				assert.True(t, response.Ocs.Meta.Success(ocsVersion), unsuccessfulResponseText)

				assert.Equal(t, defaultMemberOf[user], response.Ocs.Data.Groups)
			}
		}
	}
	cleanUp(t)
}

func TestGetGroupForUserInvalidUserId(t *testing.T) {

	invalidUsers := []string{
		"1",
		"invalid",
		"3434234233",
		"1am41validUs3r",
		"_-@@--$$__",
	}
	for _, ocsVersion := range ocsVersions {
		for _, format := range formats {
			formatpart := getFormatString(format)
			for _, user := range invalidUsers {
				t.Run(fmt.Sprintf("%s (ocs=%s, format=%s)", user, ocsVersion, format), func(t *testing.T) {
					res, err := sendRequest(
						"GET",
						fmt.Sprintf("/%s/cloud/users/%s/groups%s", ocsVersion, user, formatpart),
						"",
						&User{ID: userIDAdmin},
						[]string{ssvc.BundleUUIDRoleAdmin},
					)

					if err != nil {
						t.Fatal(err)
					}

					response := assertEmptyResponse(t, format, res)

					assertStatusCode(t, 404, res, ocsVersion)
					assert.False(t, response.Ocs.Meta.Success(ocsVersion), unsuccessfulResponseText)
					assertResponseMeta(t, Meta{
						Status:     "error",
						StatusCode: 998,
						Message:    "The requested user could not be found",
					}, response.Ocs.Meta)

					assert.Empty(t, response.Ocs.Data)
				})
			}
		}
	}
}

func TestAddUsersToGroupsNewUsers(t *testing.T) {
	users := []User{
		{
			Enabled:     "true",
			ID:          "rutherford",
			Email:       "rutherford@example.com",
			Displayname: "Ernest RutherFord",
		},
		{
			Enabled:     "true",
			ID:          "thomson",
			Email:       "thomson@example.com",
			Displayname: "J. J. Thomson",
		},
	}

	for _, ocsVersion := range ocsVersions {
		for _, format := range formats {
			formatpart := getFormatString(format)
			for _, user := range users {
				t.Run(fmt.Sprintf("%s (ocs=%s, format=%s)", user.ID, ocsVersion, format), func(t *testing.T) {
					err := createUser(user)
					if err != nil {
						t.Fatal(err)
					}

					// group id for Physics lover
					groupid := groupPhysicsLovers

					res, err := sendRequest(
						"POST",
						fmt.Sprintf("/%s/cloud/users/%s/groups%s", ocsVersion, user.ID, formatpart),
						"groupid="+groupid,
						&User{ID: userIDAdmin},
						[]string{ssvc.BundleUUIDRoleAdmin},
					)

					if err != nil {
						t.Fatal(err)
					}

					response := assertEmptyResponse(t, format, res)

					assertStatusCode(t, 200, res, ocsVersion)
					assert.True(t, response.Ocs.Meta.Success(ocsVersion), unsuccessfulResponseText)
					assert.Empty(t, response.Ocs.Data)

					// Check the user is in the group
					res, err = sendRequest(
						"GET",
						fmt.Sprintf("/%s/cloud/users/%s/groups?format=json", ocsVersion, user.ID),
						"",
						&User{ID: userIDAdmin},
						[]string{ssvc.BundleUUIDRoleAdmin},
					)
					if err != nil {
						t.Fatal(err)
					}
					var grpResponse GetUsersGroupsResponse
					if err := json.Unmarshal(res.Body.Bytes(), &grpResponse); err != nil {
						t.Fatal(err)
					}
					assert.Contains(t, grpResponse.Ocs.Data.Groups, groupid)

					cleanUp(t)
				})
			}
		}
	}
}

func TestAddUsersToGroupInvalidGroup(t *testing.T) {
	user := User{
		Enabled:     "true",
		ID:          "rutherford",
		Email:       "rutherford@example.com",
		Displayname: "Ernest RutherFord",
	}
	err := createUser(user)
	if err != nil {
		t.Fatal(err)
	}

	invalidGroups := []string{
		"1",
		"invalid",
		"3434234233",
		"1am41validUs3r",
		"_-@@--$$__",
		"c7fbe8c4-139b-4376-b307-cf0a8c2d0d9c",
	}

	for _, ocsVersion := range ocsVersions {
		for _, format := range formats {
			formatpart := getFormatString(format)
			for _, groupid := range invalidGroups {
				res, err := sendRequest(
					"POST",
					fmt.Sprintf("/%s/cloud/users/rutherford/groups%s", ocsVersion, formatpart),
					"groupid="+groupid,
					&User{ID: userIDAdmin},
					[]string{ssvc.BundleUUIDRoleAdmin},
				)

				if err != nil {
					t.Fatal(err)
				}

				response := assertEmptyResponse(t, format, res)

				assertStatusCode(t, 404, res, ocsVersion)
				assert.False(t, response.Ocs.Meta.Success(ocsVersion), "The response was expected to be fail but was successful")
				assertResponseMeta(t, Meta{
					"error",
					998,
					"The requested group could not be found",
				}, response.Ocs.Meta)
				assert.Empty(t, response.Ocs.Data)
			}
		}
	}
	cleanUp(t)
}

// Issue: https://github.com/owncloud/ocis/ocs/issues/57 - cannot remove user from group
func TestRemoveUserFromGroup(t *testing.T) {
	user := User{
		Enabled:     "true",
		ID:          "rutherford",
		Email:       "rutherford@example.com",
		Displayname: "Ernest RutherFord",
	}

	groups := []string{
		groupRadiumLovers,
		groupPoloniumLovers,
		groupPhysicsLovers,
	}

	var err error
	for _, ocsVersion := range ocsVersions {
		for _, format := range formats {
			formatpart := getFormatString(format)

			err = createUser(user)
			if err != nil {
				t.Fatalf("Failed while creating new user: %v", err)
			}
			for _, group := range groups {
				err := AddUserToGroup(user.ID, group)
				if err != nil {
					t.Fatalf("Failed while creating new user: %v", err)
				}
			}

			// Remove user from one group
			res, err := sendRequest(
				"DELETE",
				fmt.Sprintf("/%s/cloud/users/%s/groups%s", ocsVersion, user.ID, formatpart),
				"groupid="+groups[0],
				&User{ID: userIDAdmin},
				[]string{ssvc.BundleUUIDRoleAdmin},
			)

			if err != nil {
				t.Fatal(err)
			}

			response := assertEmptyResponse(t, format, res)
			assert.True(t, response.Ocs.Meta.Success(ocsVersion), unsuccessfulResponseText)
			assert.Empty(t, response.Ocs.Data)

			// Check the users are correctly added to group
			res, err = sendRequest(
				"GET",
				fmt.Sprintf("/%s/cloud/users/%s/groups?format=json", ocsVersion, user.ID),
				"",
				&User{ID: userIDAdmin},
				[]string{ssvc.BundleUUIDRoleAdmin},
			)
			if err != nil {
				t.Fatal(err)
			}
			var grpResponse GetUsersGroupsResponse
			if err := json.Unmarshal(res.Body.Bytes(), &grpResponse); err != nil {
				t.Fatal(err)
			}

			assert.NotContains(t, grpResponse.Ocs.Data.Groups, groups[0])
			assert.Contains(t, grpResponse.Ocs.Data.Groups, groups[1])
			assert.Contains(t, grpResponse.Ocs.Data.Groups, groups[2])
			cleanUp(t)
		}
	}
}

// Issue: https://github.com/owncloud/ocis-ocs/issues/59 - cloud/capabilities endpoint not implemented
func TestCapabilities(t *testing.T) {
	for _, ocsVersion := range ocsVersions {
		for _, format := range formats {
			formatpart := getFormatString(format)
			res, err := sendRequest(
				"GET",
				fmt.Sprintf("/%s/cloud/capabilities%s", ocsVersion, formatpart),
				"",
				&User{ID: userIDAdmin},
				[]string{ssvc.BundleUUIDRoleAdmin},
			)

			if err != nil {
				t.Fatal(err)
			}

			response := assertEmptyResponse(t, format, res)

			assertStatusCode(t, 404, res, ocsVersion)
			assertResponseMeta(t, Meta{
				"error",
				998,
				"not found",
			}, response.Ocs.Meta)
			assert.Empty(t, response.Ocs.Data)
		}
	}
}

func TestGetConfig(t *testing.T) {
	for _, ocsVersion := range ocsVersions {
		for _, format := range formats {
			formatpart := getFormatString(format)
			res, err := sendRequest(
				"GET",
				fmt.Sprintf("/%s/config%s", ocsVersion, formatpart),
				"",
				&User{ID: userIDAdmin},
				[]string{ssvc.BundleUUIDRoleAdmin},
			)

			if err != nil {
				t.Fatal(err)
			}

			var response GetConfigResponse
			if format == "json" {
				if err := json.Unmarshal(res.Body.Bytes(), &response); err != nil {
					t.Fatal(err)
				}
			} else {
				if err := xml.Unmarshal(res.Body.Bytes(), &response.Ocs); err != nil {
					t.Fatal(err)
				}
			}

			assertStatusCode(t, 200, res, ocsVersion)
			assert.True(t, response.Ocs.Meta.Success(ocsVersion), unsuccessfulResponseText)
			assert.Equal(t, OcsConfig{
				"1.7", "ocis", "", "", "true",
			}, response.Ocs.Data)
		}
	}
}
func TestGetGroupsDefaultGroups(t *testing.T) {
	for _, ocsVersion := range ocsVersions {
		for _, format := range formats {
			formatpart := getFormatString(format)

			res, err := sendRequest(
				"GET",
				fmt.Sprintf("/%s/cloud/groups%s", ocsVersion, formatpart),
				"",
				&User{ID: userIDAdmin},
				[]string{ssvc.BundleUUIDRoleAdmin},
			)

			if err != nil {
				t.Fatal(err)
			}

			var response GetGroupsResponse
			if format == "json" {
				if err := json.Unmarshal(res.Body.Bytes(), &response); err != nil {
					t.Fatal(err)
				}
			} else {
				if err := xml.Unmarshal(res.Body.Bytes(), &response.Ocs); err != nil {
					t.Fatal(err)
				}
			}

			assertStatusCode(t, 200, res, ocsVersion)
			assert.True(t, response.Ocs.Meta.Success(ocsVersion), unsuccessfulResponseText)
			assert.Subset(t, defaultGroups, response.Ocs.Data.Groups)
		}
	}
}

func TestCreateGroup(t *testing.T) {
	testData := []struct {
		group Group
		err   *Meta
	}{
		// A simple group
		{
			Group{
				ID:          "grp1",
				GIDNumber:   32222,
				Displayname: "Group Name",
			},
			nil,
		},
	}
	for _, ocsVersion := range ocsVersions {
		for _, format := range formats {
			for _, data := range testData {
				formatpart := getFormatString(format)
				res, err := sendRequest(
					"POST",
					fmt.Sprintf("/%v/cloud/groups%v", ocsVersion, formatpart),
					data.group.getGroupRequestString(),
					&User{ID: userIDAdmin},
					[]string{ssvc.BundleUUIDRoleAdmin},
				)

				if err != nil {
					t.Fatal(err)
				}

				response := assertEmptyResponse(t, format, res)

				if data.err == nil {
					assert.True(t, response.Ocs.Meta.Success(ocsVersion), unsuccessfulResponseText)
				} else {
					assertResponseMeta(t, *data.err, response.Ocs.Meta)
				}

				// Check the group exists of not
				res, err = sendRequest(
					"GET",
					"/v2.php/cloud/groups?format=json",
					"",
					&User{ID: userIDAdmin},
					[]string{ssvc.BundleUUIDRoleAdmin},
				)
				if err != nil {
					t.Fatal(err)
				}
				var groupResponse GetGroupsResponse
				if err := json.Unmarshal(res.Body.Bytes(), &groupResponse); err != nil {
					t.Fatal(err)
				}
				if data.err == nil {
					assert.Contains(t, groupResponse.Ocs.Data.Groups, data.group.ID)
				} else {
					assert.NotContains(t, groupResponse.Ocs.Data.Groups, data.group.ID)
				}
				cleanUp(t)
			}
		}
	}
}

func TestDeleteGroup(t *testing.T) {
	testData := []Group{
		{
			ID:          "grp1",
			GIDNumber:   32222,
			Displayname: "Group Name",
		},
	}
	for _, ocsVersion := range ocsVersions {
		for _, format := range formats {
			formatpart := getFormatString(format)
			for _, data := range testData {
				err := createGroup(data)
				if err != nil {
					t.Fatal(err)
				}
				res, err := sendRequest(
					"DELETE",
					fmt.Sprintf("/%v/cloud/groups/%v%v", ocsVersion, data.ID, formatpart),
					"groupid="+data.ID,
					&User{ID: userIDAdmin},
					[]string{ssvc.BundleUUIDRoleAdmin},
				)
				if err != nil {
					t.Fatal(err)
				}

				response := assertEmptyResponse(t, format, res)
				assert.True(t, response.Ocs.Meta.Success(ocsVersion), unsuccessfulResponseText)

				// Check the group does not exists
				res, err = sendRequest(
					"GET",
					"/v2.php/cloud/groups?format=json",
					"",
					&User{ID: userIDAdmin},
					[]string{ssvc.BundleUUIDRoleAdmin},
				)
				if err != nil {
					t.Fatal(err)
				}
				var groupResponse GetGroupsResponse
				if err := json.Unmarshal(res.Body.Bytes(), &groupResponse); err != nil {
					t.Fatal(err)
				}

				assert.NotContains(t, groupResponse.Ocs.Data.Groups, data.ID)
				cleanUp(t)
			}
		}
	}
}

func TestDeleteGroupInvalidGroups(t *testing.T) {
	testData := []string{
		"1",
		"invalid",
		"3434234233",
		"1am41validUs3r",
		"_-@@--$$__",
	}
	for _, ocsVersion := range ocsVersions {
		for _, format := range formats {
			formatpart := getFormatString(format)
			for _, data := range testData {
				res, err := sendRequest(
					"DELETE",
					fmt.Sprintf("/%v/cloud/groups/%v%v", ocsVersion, data, formatpart),
					"groupid="+data,
					&User{ID: userIDAdmin},
					[]string{ssvc.BundleUUIDRoleAdmin},
				)

				if err != nil {
					t.Fatal(err)
				}

				response := assertEmptyResponse(t, format, res)

				assertStatusCode(t, 404, res, ocsVersion)
				assert.False(t, response.Ocs.Meta.Success(ocsVersion), "The response was expected to fail but was successful")
				assertResponseMeta(t, Meta{
					"error",
					998,
					"The requested group could not be found",
				}, response.Ocs.Meta)
				cleanUp(t)
			}
		}
	}
}

func TestGetGroupMembersDefaultGroups(t *testing.T) {
	for _, ocsVersion := range ocsVersions {
		for _, format := range formats {
			for group, members := range defaultMembers {
				formatpart := getFormatString(format)
				res, err := sendRequest(
					"GET",
					fmt.Sprintf("/%v/cloud/groups/%v%v", ocsVersion, group, formatpart),
					"",
					&User{ID: userIDAdmin},
					[]string{ssvc.BundleUUIDRoleAdmin},
				)

				if err != nil {
					t.Fatal(err)
				}

				var response GetGroupMembersResponse

				if format == "json" {
					if err := json.Unmarshal(res.Body.Bytes(), &response); err != nil {
						t.Fatal(err)
					}
				} else {
					if err := xml.Unmarshal(res.Body.Bytes(), &response.Ocs); err != nil {
						t.Fatal(err)
					}
				}

				assertStatusCode(t, 200, res, ocsVersion)
				assert.True(t, response.Ocs.Meta.Success(ocsVersion), unsuccessfulResponseText+" for group "+group)
				assert.Equal(t, members, response.Ocs.Data.Users)

				cleanUp(t)
			}
		}
	}
}

func TestListMembersInvalidGroups(t *testing.T) {
	testData := []string{
		"1",
		"invalid",
		"3434234233",
		"1am41validUs3r",
		"_-@@--$$__",
	}
	for _, ocsVersion := range ocsVersions {
		for _, format := range formats {
			formatpart := getFormatString(format)
			for _, group := range testData {
				res, err := sendRequest(
					"GET",
					fmt.Sprintf("/%v/cloud/groups/%v%v", ocsVersion, group, formatpart),
					"",
					&User{ID: userIDAdmin},
					[]string{ssvc.BundleUUIDRoleAdmin},
				)

				if err != nil {
					t.Fatal(err)
				}

				response := assertEmptyResponse(t, format, res)

				assertStatusCode(t, 404, res, ocsVersion)
				assert.False(t, response.Ocs.Meta.Success(ocsVersion), "The response was expected to fail but was successful")
				assertResponseMeta(t, Meta{
					"error",
					998,
					"The requested group could not be found",
				}, response.Ocs.Meta)
				cleanUp(t)
			}
		}
	}
}
