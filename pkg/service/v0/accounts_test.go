package service

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/micro/go-micro/v2/client"
	merrors "github.com/micro/go-micro/v2/errors"
	"github.com/micro/go-micro/v2/metadata"
	"github.com/owncloud/ocis-accounts/pkg/config"
	"github.com/owncloud/ocis-accounts/pkg/proto/v0"
	olog "github.com/owncloud/ocis-pkg/v2/log"
	"github.com/owncloud/ocis-pkg/v2/middleware"
	"github.com/owncloud/ocis-pkg/v2/roles"
	settings "github.com/owncloud/ocis-settings/pkg/proto/v0"
	ssvc "github.com/owncloud/ocis-settings/pkg/service/v0"
	"github.com/stretchr/testify/assert"
)

const dataPath = "/var/tmp/ocis-accounts-tests"

var (
	roleServiceMock settings.RoleService
	s               *Service
)

func init() {
	cfg := config.New()
	cfg.Server.Name = "accounts"
	cfg.Server.AccountsDataPath = dataPath
	logger := olog.NewLogger(olog.Color(true), olog.Pretty(true))
	roleServiceMock = buildRoleServiceMock()
	roleManager := roles.NewManager(
		roles.Logger(logger),
		roles.RoleService(roleServiceMock),
		roles.CacheTTL(time.Hour),
		roles.CacheSize(1024),
	)
	s, _ = New(
		Logger(logger),
		Config(cfg),
		RoleService(roleServiceMock),
		RoleManager(&roleManager),
	)
}

func setup() (teardown func()) {
	return func() {
		if err := os.RemoveAll(dataPath); err != nil {
			log.Printf("could not delete data root: %s", dataPath)
		} else {
			log.Println("data root deleted")
		}
	}
}

// TestPermissionsListAccounts checks permission handling on ListAccounts
func TestPermissionsListAccounts(t *testing.T) {
	var scenarios = []struct {
		name            string
		roleIDs         []string
		query           string
		permissionError error
	}{
		// TODO: remove this test when https://github.com/owncloud/ocis-accounts/pull/111 is merged
		// replace with two tests:
		// 1: "ListAccounts fails with 403 when roleIDs don't exist in context"
		// 2: "ListAccounts fails with 403 when ('no admin role in context' AND 'empty query')"
		{
			"ListAccounts succeeds when no roleIDs in context",
			nil,
			"",
			nil,
		},
		{
			"ListAccounts fails when no admin roleID in context",
			[]string{ssvc.BundleUUIDRoleUser, ssvc.BundleUUIDRoleGuest},
			"",
			merrors.Forbidden(s.id, "no permission for ListAccounts"),
		},
		{
			"ListAccounts succeeds when admin roleID in context",
			[]string{ssvc.BundleUUIDRoleAdmin},
			"",
			nil,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			teardown := setup()
			defer teardown()

			ctx := buildTestCtx(t, scenario.roleIDs)
			request := &proto.ListAccountsRequest{
				Query: scenario.query,
			}
			response := &proto.ListAccountsResponse{}
			err := s.ListAccounts(ctx, request, response)
			if scenario.permissionError != nil {
				assert.Equal(t, scenario.permissionError, err)
			} else if err != nil {
				// we are only checking permissions here, so just check that the error code is not 403
				merr := merrors.FromError(err)
				assert.NotEqual(t, http.StatusForbidden, merr.GetCode())
			}
		})
	}
}

// TestPermissionsGetAccount checks permission handling on GetAccount
// TODO: remove this test function entirely, when https://github.com/owncloud/ocis-accounts/pull/111 is merged. GetAccount will not have permission checks for the time being.
func TestPermissionsGetAccount(t *testing.T) {
	var scenarios = []struct {
		name            string
		roleIDs         []string
		permissionError error
	}{
		{
			"GetAccount succeeds when no role IDs in context",
			nil,
			nil,
		},
		{
			"GetAccount fails when no admin roleID in context",
			[]string{ssvc.BundleUUIDRoleUser, ssvc.BundleUUIDRoleGuest},
			merrors.Forbidden(s.id, "no permission for GetAccount"),
		},
		{
			"GetAccount succeeds when admin roleID in context",
			[]string{ssvc.BundleUUIDRoleAdmin},
			nil,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			teardown := setup()
			defer teardown()

			ctx := buildTestCtx(t, scenario.roleIDs)
			request := &proto.GetAccountRequest{}
			response := &proto.Account{}
			err := s.GetAccount(ctx, request, response)
			if scenario.permissionError != nil {
				assert.Equal(t, scenario.permissionError, err)
			} else if err != nil {
				// we are only checking permissions here, so just check that the error code is not 403
				merr := merrors.FromError(err)
				assert.NotEqual(t, http.StatusForbidden, merr.GetCode())
			}
		})
	}
}

// TestPermissionsUpdateAccount checks permission handling on UpdateAccount
func TestPermissionsUpdateAccount(t *testing.T) {
	var scenarios = []struct {
		name            string
		roleIDs         []string
		permissionError error
	}{
		// TODO: remove this test when https://github.com/owncloud/ocis-accounts/pull/111 is merged
		// replace with two tests:
		// 1: "UpdateAccount fails with 403 when roleIDs don't exist in context"
		// 2: "UpdateAccount fails with 403 when no admin role in context"
		{
			"UpdateAccount succeeds when no role IDs in context",
			nil,
			nil,
		},
		{
			"UpdateAccount fails when no admin roleID in context",
			[]string{ssvc.BundleUUIDRoleUser, ssvc.BundleUUIDRoleGuest},
			merrors.Forbidden(s.id, "no permission for UpdateAccount"),
		},
		{
			"UpdateAccount succeeds when admin roleID in context",
			[]string{ssvc.BundleUUIDRoleAdmin},
			nil,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			teardown := setup()
			defer teardown()

			ctx := buildTestCtx(t, scenario.roleIDs)
			request := &proto.UpdateAccountRequest{}
			response := &proto.Account{}
			err := s.UpdateAccount(ctx, request, response)
			if scenario.permissionError != nil {
				assert.Equal(t, scenario.permissionError, err)
			} else if err != nil {
				// we are only checking permissions here, so just check that the error code is not 403
				merr := merrors.FromError(err)
				assert.NotEqual(t, http.StatusForbidden, merr.GetCode())
			}
		})
	}
}

func buildTestCtx(t *testing.T, roleIDs []string) context.Context {
	ctx := context.Background()
	if roleIDs != nil {
		roleIDs, err := json.Marshal(roleIDs)
		assert.NoError(t, err)
		ctx = metadata.Set(ctx, middleware.RoleIDs, string(roleIDs))
	}
	return ctx
}

func buildRoleServiceMock() settings.RoleService {
	defaultRoles := map[string]*settings.Bundle{
		ssvc.BundleUUIDRoleAdmin: {
			Id:   ssvc.BundleUUIDRoleAdmin,
			Type: settings.Bundle_TYPE_ROLE,
			Resource: &settings.Resource{
				Type: settings.Resource_TYPE_SYSTEM,
			},
			Settings: []*settings.Setting{
				{
					Id: AccountManagementPermissionID,
				},
			},
		},
		ssvc.BundleUUIDRoleUser: {
			Id:   ssvc.BundleUUIDRoleUser,
			Type: settings.Bundle_TYPE_ROLE,
			Resource: &settings.Resource{
				Type: settings.Resource_TYPE_SYSTEM,
			},
			Settings: []*settings.Setting{},
		},
	}
	return settings.MockRoleService{
		ListRolesFunc: func(ctx context.Context, req *settings.ListBundlesRequest, opts ...client.CallOption) (res *settings.ListBundlesResponse, err error) {
			payload := make([]*settings.Bundle, 0)
			for _, roleID := range req.BundleIds {
				if defaultRoles[roleID] != nil {
					payload = append(payload, defaultRoles[roleID])
				}
			}
			return &settings.ListBundlesResponse{
				Bundles: payload,
			}, nil
		},
		AssignRoleToUserFunc: func(ctx context.Context, req *settings.AssignRoleToUserRequest, opts ...client.CallOption) (res *settings.AssignRoleToUserResponse, err error) {
			// mock can be empty. function is called during service start. actual role assignments not needed for the tests.
			return &settings.AssignRoleToUserResponse{
				Assignment: &settings.UserRoleAssignment{},
			}, nil
		},
	}
}
