package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	userpb "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	revactx "github.com/owncloud/reva/v2/pkg/ctx"
	"go-micro.dev/v4/client"

	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/ocis-pkg/roles"
	settingsmsg "github.com/owncloud/ocis/v2/protogen/gen/ocis/messages/settings/v0"
	settingssvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/settings/v0"
	settings "github.com/owncloud/ocis/v2/services/settings/pkg/service/v0"
	"google.golang.org/protobuf/types/known/emptypb"
)

// fakeRoleService returns a single role that carries the vault permission only
// when hasVaultPermission is true.
type fakeRoleService struct {
	roleID             string
	hasVaultPermission bool
}

func (f fakeRoleService) ListRoles(_ context.Context, _ *settingssvc.ListBundlesRequest, _ ...client.CallOption) (*settingssvc.ListBundlesResponse, error) {
	bundle := &settingsmsg.Bundle{Id: f.roleID}
	if f.hasVaultPermission {
		bundle.Settings = []*settingsmsg.Setting{{Id: settings.VaultModePermissionID}}
	} else {
		bundle.Settings = []*settingsmsg.Setting{{Id: "some-other-permission"}}
	}
	return &settingssvc.ListBundlesResponse{Bundles: []*settingsmsg.Bundle{bundle}}, nil
}

func (f fakeRoleService) ListRoleAssignments(_ context.Context, _ *settingssvc.ListRoleAssignmentsRequest, _ ...client.CallOption) (*settingssvc.ListRoleAssignmentsResponse, error) {
	return &settingssvc.ListRoleAssignmentsResponse{
		Assignments: []*settingsmsg.UserRoleAssignment{{RoleId: f.roleID}},
	}, nil
}

func (f fakeRoleService) ListRoleAssignmentsFiltered(_ context.Context, _ *settingssvc.ListRoleAssignmentsFilteredRequest, _ ...client.CallOption) (*settingssvc.ListRoleAssignmentsResponse, error) {
	return &settingssvc.ListRoleAssignmentsResponse{}, nil
}

func (f fakeRoleService) AssignRoleToUser(_ context.Context, _ *settingssvc.AssignRoleToUserRequest, _ ...client.CallOption) (*settingssvc.AssignRoleToUserResponse, error) {
	return &settingssvc.AssignRoleToUserResponse{}, nil
}

func (f fakeRoleService) RemoveRoleFromUser(_ context.Context, _ *settingssvc.RemoveRoleFromUserRequest, _ ...client.CallOption) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func newManager(rs settingssvc.RoleService) *roles.Manager {
	m := roles.NewManager(
		roles.Logger(log.NopLogger()),
		roles.RoleService(rs),
	)
	return &m
}

func requestWithUser() *http.Request {
	req := httptest.NewRequest("GET", "/vault/graph/v1.0/me/drives", nil)
	ctx := revactx.ContextSetUser(req.Context(), &userpb.User{
		Id: &userpb.UserId{OpaqueId: "student-opaque-id"},
	})
	return req.WithContext(ctx)
}

func TestRequireVaultPermission_Denied(t *testing.T) {
	rm := newManager(fakeRoleService{roleID: "user-light", hasVaultPermission: false})
	handler := RequireVaultPermission(rm, log.NopLogger())(dummyHandler{})

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, requestWithUser())

	if rr.Code != http.StatusForbidden {
		t.Errorf("got status %d, want %d", rr.Code, http.StatusForbidden)
	}
	// a role denial must not carry the MFA step-up header
	if got := rr.Header().Get("X-Ocis-Mfa-Required"); got != "" {
		t.Errorf("X-Ocis-Mfa-Required header must be absent for a role denial, got %q", got)
	}
}

func TestRequireVaultPermission_Allowed(t *testing.T) {
	rm := newManager(fakeRoleService{roleID: "user", hasVaultPermission: true})
	handler := RequireVaultPermission(rm, log.NopLogger())(dummyHandler{})

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, requestWithUser())

	if rr.Code != http.StatusOK {
		t.Errorf("got status %d, want %d", rr.Code, http.StatusOK)
	}
}

func TestRequireVaultPermission_NoUser(t *testing.T) {
	rm := newManager(fakeRoleService{roleID: "user", hasVaultPermission: true})
	handler := RequireVaultPermission(rm, log.NopLogger())(dummyHandler{})

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, httptest.NewRequest("GET", "/vault/graph/v1.0/me/drives", nil))

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("got status %d, want %d", rr.Code, http.StatusUnauthorized)
	}
}
