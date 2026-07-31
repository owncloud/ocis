package userroles

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	cs3user "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	permissions "github.com/cs3org/go-cs3apis/cs3/permissions/v1beta1"
	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	storageprovider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	settingsmsg "github.com/owncloud/ocis/v2/protogen/gen/ocis/messages/settings/v0"
	settingssvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/settings/v0"
	graphmocks "github.com/owncloud/ocis/v2/services/graph/mocks"
	"github.com/owncloud/ocis/v2/services/proxy/pkg/config"
	"github.com/owncloud/reva/v2/pkg/rgrpc/todo/pool"
	cs3mocks "github.com/owncloud/reva/v2/tests/cs3mocks/mocks"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestExtractRolesArray(t *testing.T) {
	byt := []byte(`{"roles":["a","b"]}`)

	claims := map[string]interface{}{}
	err := json.Unmarshal(byt, &claims)
	if err != nil {
		t.Fatal(err)
	}

	roles, err := extractRoles("roles", claims)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := roles["a"]; !ok {
		t.Fatal("must contain 'a'")
	}
	if _, ok := roles["b"]; !ok {
		t.Fatal("must contain 'b'")
	}
}

func TestExtractRolesString(t *testing.T) {
	byt := []byte(`{"roles":"a"}`)

	claims := map[string]interface{}{}
	err := json.Unmarshal(byt, &claims)
	if err != nil {
		t.Fatal(err)
	}

	roles, err := extractRoles("roles", claims)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := roles["a"]; !ok {
		t.Fatal("must contain 'a'")
	}
}

func TestExtractRolesPathArray(t *testing.T) {
	byt := []byte(`{"sub":{"roles":["a","b"]}}`)

	claims := map[string]interface{}{}
	err := json.Unmarshal(byt, &claims)
	if err != nil {
		t.Fatal(err)
	}

	roles, err := extractRoles("sub.roles", claims)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := roles["a"]; !ok {
		t.Fatal("must contain 'a'")
	}
	if _, ok := roles["b"]; !ok {
		t.Fatal("must contain 'b'")
	}
}

func TestExtractRolesPathString(t *testing.T) {
	byt := []byte(`{"sub":{"roles":"a"}}`)

	claims := map[string]interface{}{}
	err := json.Unmarshal(byt, &claims)
	if err != nil {
		t.Fatal(err)
	}

	roles, err := extractRoles("sub.roles", claims)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := roles["a"]; !ok {
		t.Fatal("must contain 'a'")
	}
}

func TestExtractEscapedRolesPathString(t *testing.T) {
	byt := []byte(`{"sub.roles":"a"}`)

	claims := map[string]interface{}{}
	err := json.Unmarshal(byt, &claims)
	if err != nil {
		t.Fatal(err)
	}

	roles, err := extractRoles("sub\\.roles", claims)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := roles["a"]; !ok {
		t.Fatal("must contain 'a'")
	}
}

func TestNoRoles(t *testing.T) {
	byt := []byte(`{"sub":{"foo":"a"}}`)

	claims := map[string]interface{}{}
	err := json.Unmarshal(byt, &claims)
	if err != nil {
		t.Fatal(err)
	}

	roles, err := extractRoles("sub.roles", claims)
	if err == nil {
		t.Fatal("must not find a role")
	}
	if len(roles) != 0 {
		t.Fatal("length of roles mut be 0")
	}
}

func TestMatchesClaimMappingExact(t *testing.T) {
	claimRoles := map[string]struct{}{
		"ocis-user": {},
	}
	if !matchesClaimMapping("ocis-user", claimRoles) {
		t.Fatal("expected exact match to succeed")
	}
	if matchesClaimMapping("admin", claimRoles) {
		t.Fatal("expected non-matching literal to fail")
	}
}

func TestMatchesClaimMappingRegex(t *testing.T) {
	claimRoles := map[string]struct{}{
		"ocis-user-1":   {},
		"ocis-user-42":  {},
		"ocis-user-lth": {},
		"admin":         {},
	}
	if !matchesClaimMapping("ocis-user-.*", claimRoles) {
		t.Fatal("expected regex match to succeed")
	}
	if !matchesClaimMapping("ocis-user-[a-zA-Z0-9]", claimRoles) {
		t.Fatal("expected regex match to succeed")
	}
	if matchesClaimMapping("admin-.*", claimRoles) {
		t.Fatal("expected regex match to fail for admin-.*")
	}
}

func TestMatchesClaimMappingInvalidRegexFallsBackToExact(t *testing.T) {
	claimRoles := map[string]struct{}{"ocis-user": {}}
	// invalid regex pattern
	if matchesClaimMapping("ocis-user[", claimRoles) {
		t.Fatal("invalid regex should fall back to exact and not match")
	}
}

// TestUpdateUserRoleAssignmentFailsClosedOnInconclusivePermission verifies the
// fix for an issue when the Drives.Create permission check returns a non-OK
// status that is not PERMISSION_DENIED (here CODE_INTERNAL), the permission is
// indeterminate. The assigner must NOT disable (trash) the personal space and
// must revert the role assignment it just persisted.
func TestUpdateUserRoleAssignmentFailsClosedOnInconclusivePermission(t *testing.T) {
	const (
		newRoleName  = "ocis-user"
		newRoleID    = "new-role-id"
		oldRoleID    = "old-role-id"
		newAssignID  = "new-assignment-id"
		userOpaqueID = "user-1"
	)

	// reset the package-global role-name cache so this test is deterministic
	roleNameToID.lock.Lock()
	roleNameToID.roleNameToID = nil
	roleNameToID.lastRead = time.Time{}
	roleNameToID.lock.Unlock()

	gatewayClient := &cs3mocks.GatewayAPIClient{}
	gatewaySelector := pool.GetSelector[gateway.GatewayAPIClient](
		"GatewaySelector",
		"com.owncloud.api.gateway",
		func(cc grpc.ClientConnInterface) gateway.GatewayAPIClient {
			return gatewayClient
		},
	)
	defer pool.RemoveSelector("GatewaySelector" + "com.owncloud.api.gateway")

	// admin context creation authenticates the service account
	gatewayClient.On("Authenticate", mock.Anything, mock.Anything).Return(&gateway.AuthenticateResponse{
		Status: &rpc.Status{Code: rpc.Code_CODE_OK},
		Token:  "service-token",
	}, nil)
	// the permission check is inconclusive: non-OK status, not PERMISSION_DENIED
	gatewayClient.On("CheckPermission", mock.Anything, mock.Anything).Return(&permissions.CheckPermissionResponse{
		Status: &rpc.Status{Code: rpc.Code_CODE_INTERNAL, Message: "settings unavailable"},
	}, nil)
	// the space lookup/delete must never be reached on the indeterminate path,
	// but mock them so a regression (calling them) is observable rather than a panic
	gatewayClient.On("ListStorageSpaces", mock.Anything, mock.Anything).Return(&storageprovider.ListStorageSpacesResponse{
		Status:        &rpc.Status{Code: rpc.Code_CODE_OK},
		StorageSpaces: []*storageprovider.StorageSpace{{Id: &storageprovider.StorageSpaceId{OpaqueId: "ps1"}}},
	}, nil)
	gatewayClient.On("DeleteStorageSpace", mock.Anything, mock.Anything).Return(&storageprovider.DeleteStorageSpaceResponse{
		Status: &rpc.Status{Code: rpc.Code_CODE_OK},
	}, nil)

	roleService := &graphmocks.RoleService{}
	roleService.On("ListRoles", mock.Anything, mock.Anything, mock.Anything).Return(
		&settingssvc.ListBundlesResponse{Bundles: []*settingsmsg.Bundle{{Id: newRoleID, Name: newRoleName}}}, nil)
	// user currently has the old role, so a re-assignment is triggered
	roleService.On("ListRoleAssignments", mock.Anything, mock.Anything, mock.Anything).Return(
		&settingssvc.ListRoleAssignmentsResponse{Assignments: []*settingsmsg.UserRoleAssignment{
			{Id: "old-assignment-id", AccountUuid: userOpaqueID, RoleId: oldRoleID},
		}}, nil)
	// the initial assignment to the new role
	roleService.On("AssignRoleToUser", mock.Anything, mock.MatchedBy(func(req *settingssvc.AssignRoleToUserRequest) bool {
		return req.GetRoleId() == newRoleID
	}), mock.Anything).Return(&settingssvc.AssignRoleToUserResponse{Assignment: &settingsmsg.UserRoleAssignment{
		Id:          newAssignID,
		AccountUuid: userOpaqueID,
		RoleId:      newRoleID,
	}}, nil)
	// the revert back to the previous role
	roleService.On("AssignRoleToUser", mock.Anything, mock.MatchedBy(func(req *settingssvc.AssignRoleToUserRequest) bool {
		return req.GetRoleId() == oldRoleID
	}), mock.Anything).Return(&settingssvc.AssignRoleToUserResponse{Assignment: &settingsmsg.UserRoleAssignment{
		Id:          "old-assignment-id",
		AccountUuid: userOpaqueID,
		RoleId:      oldRoleID,
	}}, nil)
	roleService.On("RemoveRoleFromUser", mock.Anything, mock.Anything, mock.Anything).Return(&emptypb.Empty{}, nil)

	ra := oidcRoleAssigner{Options: Options{
		logger:          log.NopLogger(),
		gatewaySelector: gatewaySelector,
		roleService:     roleService,
		serviceAccount:  config.ServiceAccount{ServiceAccountID: "service-account", ServiceAccountSecret: "secret"},
	}}

	user := &cs3user.User{Id: &cs3user.UserId{OpaqueId: userOpaqueID}}

	_, err := ra.UpdateUserRoleAssignment(context.Background(), user, nil, newRoleName)
	if err == nil {
		t.Fatal("expected an error to be returned on an inconclusive permission check")
	}

	gatewayClient.AssertNotCalled(t, "DeleteStorageSpace", mock.Anything, mock.Anything)
	// the role assignment must be reverted back to the previous role
	roleService.AssertCalled(t, "AssignRoleToUser", mock.Anything, mock.MatchedBy(func(req *settingssvc.AssignRoleToUserRequest) bool {
		return req.GetRoleId() == oldRoleID
	}), mock.Anything)
}
