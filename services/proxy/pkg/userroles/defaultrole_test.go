package userroles

import (
	"context"
	"testing"

	cs3user "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	settingsmsg "github.com/owncloud/ocis/v2/protogen/gen/ocis/messages/settings/v0"
	settingssvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/settings/v0"
	graphmocks "github.com/owncloud/ocis/v2/services/graph/mocks"
	settingsService "github.com/owncloud/ocis/v2/services/settings/pkg/service/v0"
	"github.com/owncloud/reva/v2/pkg/utils"
	"github.com/stretchr/testify/mock"
)

// TestApplyUserRoleBootstrapsDefaultRoleForNewUser verifies that ApplyUserRole,
// just like UpdateUserRoleAssignment, assigns a default role to a user who
// doesn't have one yet. This matters for authentication paths (such as app
// passwords) that only call ApplyUserRole and never go through claims-based
// UpdateUserRoleAssignment, so a user who never logged in interactively must
// still end up with a role.
func TestApplyUserRoleBootstrapsDefaultRoleForNewUser(t *testing.T) {
	const userOpaqueID = "user-1"

	roleService := &graphmocks.RoleService{}
	roleService.On("ListRoleAssignments", mock.Anything, mock.Anything, mock.Anything).Return(
		&settingssvc.ListRoleAssignmentsResponse{}, nil)
	roleService.On("AssignRoleToUser", mock.Anything, mock.MatchedBy(func(req *settingssvc.AssignRoleToUserRequest) bool {
		return req.GetAccountUuid() == userOpaqueID && req.GetRoleId() == settingsService.BundleUUIDRoleUser
	}), mock.Anything).Return(&settingssvc.AssignRoleToUserResponse{Assignment: &settingsmsg.UserRoleAssignment{
		Id:          "new-assignment-id",
		AccountUuid: userOpaqueID,
		RoleId:      settingsService.BundleUUIDRoleUser,
	}}, nil)

	ra := NewDefaultRoleAssigner(WithRoleService(roleService), WithLogger(log.NopLogger()))

	user := &cs3user.User{Id: &cs3user.UserId{OpaqueId: userOpaqueID, Type: cs3user.UserType_USER_TYPE_PRIMARY}}

	updated, err := ra.ApplyUserRole(context.Background(), user)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var roleIDs []string
	if err := utils.ReadJSONFromOpaque(updated.GetOpaque(), "roles", &roleIDs); err != nil {
		t.Fatalf("could not read roles from opaque: %v", err)
	}
	if len(roleIDs) != 1 || roleIDs[0] != settingsService.BundleUUIDRoleUser {
		t.Fatalf("expected the default user role to be bootstrapped, got %v", roleIDs)
	}

	roleService.AssertCalled(t, "AssignRoleToUser", mock.Anything, mock.MatchedBy(func(req *settingssvc.AssignRoleToUserRequest) bool {
		return req.GetAccountUuid() == userOpaqueID && req.GetRoleId() == settingsService.BundleUUIDRoleUser
	}), mock.Anything)
}

// TestApplyUserRoleDoesNotBootstrapForLightweightUser verifies that
// ApplyUserRole never assigns a default role to a lightweight (federated)
// user, mirroring UpdateUserRoleAssignment's existing behavior.
func TestApplyUserRoleDoesNotBootstrapForLightweightUser(t *testing.T) {
	const userOpaqueID = "lightweight-user-1"

	roleService := &graphmocks.RoleService{}
	roleService.On("ListRoleAssignments", mock.Anything, mock.Anything, mock.Anything).Return(
		&settingssvc.ListRoleAssignmentsResponse{}, nil)

	ra := NewDefaultRoleAssigner(WithRoleService(roleService), WithLogger(log.NopLogger()))

	user := &cs3user.User{Id: &cs3user.UserId{OpaqueId: userOpaqueID, Type: cs3user.UserType_USER_TYPE_LIGHTWEIGHT}}

	_, err := ra.ApplyUserRole(context.Background(), user)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	roleService.AssertNotCalled(t, "AssignRoleToUser", mock.Anything, mock.Anything, mock.Anything)
}
