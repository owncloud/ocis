package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	userv1beta1 "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	ocismw "github.com/owncloud/ocis/v2/ocis-pkg/middleware"
	"github.com/owncloud/ocis/v2/ocis-pkg/roles"
	settingsmsg "github.com/owncloud/ocis/v2/protogen/gen/ocis/messages/settings/v0"
	settingssvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/settings/v0"
	settingsdefaults "github.com/owncloud/ocis/v2/services/settings/pkg/store/defaults"
	revactx "github.com/owncloud/reva/v2/pkg/ctx"
	"go-micro.dev/v4/client"
	"go-micro.dev/v4/metadata"
)

// roleService adds the method the generated MockRoleService lacks, so it satisfies settingssvc.RoleService.
type roleService struct {
	settingssvc.MockRoleService
}

func (roleService) ListRoleAssignmentsFiltered(context.Context, *settingssvc.ListRoleAssignmentsFilteredRequest, ...client.CallOption) (*settingssvc.ListRoleAssignmentsResponse, error) {
	panic("ListRoleAssignmentsFiltered was called in test but not mocked")
}

func TestRequireVaultPermission(t *testing.T) {
	const roleID = "role-1"
	vaultPermissionID := settingsdefaults.VaultModePermission(settingsdefaults.Own).GetId()

	bundleWith := func(permissionIDs ...string) *settingsmsg.Bundle {
		settings := make([]*settingsmsg.Setting, 0, len(permissionIDs))
		for _, id := range permissionIDs {
			settings = append(settings, &settingsmsg.Setting{Id: id})
		}
		return &settingsmsg.Bundle{Id: roleID, Settings: settings}
	}

	newManager := func(bundle *settingsmsg.Bundle) *roles.Manager {
		rs := roleService{MockRoleService: settingssvc.MockRoleService{
			ListRolesFunc: func(ctx context.Context, req *settingssvc.ListBundlesRequest, opts ...client.CallOption) (*settingssvc.ListBundlesResponse, error) {
				return &settingssvc.ListBundlesResponse{Bundles: []*settingsmsg.Bundle{bundle}}, nil
			},
		}}
		m := roles.NewManager(roles.Logger(log.NopLogger()), roles.RoleService(rs))
		return &m
	}

	withUserAndRoles := func(ctx context.Context) context.Context {
		ctx = revactx.ContextSetUser(ctx, &userv1beta1.User{Id: &userv1beta1.UserId{OpaqueId: "user-1"}})
		roleIDsJSON, _ := json.Marshal([]string{roleID})
		return metadata.Set(ctx, ocismw.RoleIDs, string(roleIDsJSON))
	}

	tests := []struct {
		name       string
		ctx        context.Context
		bundle     *settingsmsg.Bundle
		wantStatus int
	}{
		{
			name:       "role holds the vault permission -> allowed",
			ctx:        withUserAndRoles(context.Background()),
			bundle:     bundleWith(vaultPermissionID),
			wantStatus: http.StatusOK,
		},
		{
			name:       "role lacks the vault permission -> forbidden",
			ctx:        withUserAndRoles(context.Background()),
			bundle:     bundleWith("some-other-permission"),
			wantStatus: http.StatusForbidden,
		},
		{
			name:       "no user in context -> unauthorized",
			ctx:        context.Background(),
			bundle:     bundleWith(vaultPermissionID),
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := RequireVaultPermission(newManager(tt.bundle), log.NopLogger())(dummyHandler{})

			req := httptest.NewRequest(http.MethodGet, "/vault/graph", nil).WithContext(tt.ctx)
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			if rr.Code != tt.wantStatus {
				t.Errorf("got status %d, want %d", rr.Code, tt.wantStatus)
			}
		})
	}
}
