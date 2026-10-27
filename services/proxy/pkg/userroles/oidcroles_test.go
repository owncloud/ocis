package userroles

import (
	"context"
	"encoding/json"
	"reflect"
	"strings"
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

func TestMapGuestRole(t *testing.T) {
	roleIDs := map[string]string{"user-light": "user-light-id", "user-light-vault": "user-light-vault-id"}
	spaceAdminMapping := []config.RoleMapping{{ClaimValue: "spaceadmin", RoleName: "user-light-vault"}}

	tests := []struct {
		name    string
		mapping []config.RoleMapping
		roleIDs map[string]string
		claims  map[string]interface{}
		want    string
		wantErr bool
	}{
		{
			name:    "spaceadmin gets the mapped role",
			mapping: spaceAdminMapping,
			roleIDs: roleIDs,
			claims:  map[string]interface{}{"roles": []interface{}{"spaceadmin"}},
			want:    "user-light-vault",
		},
		{
			name:    "spaceadmin with more claim values gets the mapped role",
			mapping: spaceAdminMapping,
			roleIDs: roleIDs,
			claims:  map[string]interface{}{"roles": []interface{}{"user", "spaceadmin"}},
			want:    "user-light-vault",
		},
		{
			name:    "user keeps the guest role",
			mapping: spaceAdminMapping,
			roleIDs: roleIDs,
			claims:  map[string]interface{}{"roles": []interface{}{"user"}},
			want:    "user-light",
		},
		{
			name:    "guest without roles claim keeps the guest role",
			mapping: spaceAdminMapping,
			roleIDs: roleIDs,
			claims:  map[string]interface{}{"sub": "someone"},
			want:    "user-light",
		},
		{
			name:    "mapping to an unknown role is an error",
			mapping: spaceAdminMapping,
			roleIDs: map[string]string{"user-light": "user-light-id"},
			claims:  map[string]interface{}{"roles": []interface{}{"spaceadmin"}},
			wantErr: true,
		},
		{
			name:    "unknown role of a non-matching mapping is no error",
			mapping: spaceAdminMapping,
			roleIDs: map[string]string{"user-light": "user-light-id"},
			claims:  map[string]interface{}{"roles": []interface{}{"user"}},
			want:    "user-light",
		},
		{
			name: "first matching mapping wins",
			mapping: []config.RoleMapping{
				{ClaimValue: "spaceadmin", RoleName: "user-light-vault"},
				{ClaimValue: ".*", RoleName: "user-light"},
			},
			roleIDs: roleIDs,
			claims:  map[string]interface{}{"roles": []interface{}{"spaceadmin"}},
			want:    "user-light-vault",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ra := oidcRoleAssigner{Options: Options{rolesClaim: "roles", guestMapping: tt.mapping}}
			got, err := ra.mapGuestRole(tt.claims, "user-light", tt.roleIDs)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected an error, got role %q", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("got role %q, want %q", got, tt.want)
			}
		})
	}
}

// TestUpdateUserRoleAssignmentMapsGuestRole verifies that a guest who already holds the
// guest role is moved to the role the guest role mapping picks for their claims, which is
// how guests invited before the mapping was configured get the new role on their next login.
func TestUpdateUserRoleAssignmentMapsGuestRole(t *testing.T) {
	const (
		guestRoleID  = "user-light-id"
		vaultRoleID  = "user-light-vault-id"
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

	gatewayClient.On("Authenticate", mock.Anything, mock.Anything).Return(&gateway.AuthenticateResponse{
		Status: &rpc.Status{Code: rpc.Code_CODE_OK},
		Token:  "service-token",
	}, nil)
	// like the guest role, the mapped role can't create drives, so the assigner disables
	// the (here non-existent) personal space
	gatewayClient.On("CheckPermission", mock.Anything, mock.Anything).Return(&permissions.CheckPermissionResponse{
		Status: &rpc.Status{Code: rpc.Code_CODE_PERMISSION_DENIED},
	}, nil)
	gatewayClient.On("ListStorageSpaces", mock.Anything, mock.Anything).Return(&storageprovider.ListStorageSpacesResponse{
		Status: &rpc.Status{Code: rpc.Code_CODE_OK},
	}, nil)

	roleService := &graphmocks.RoleService{}
	roleService.On("ListRoles", mock.Anything, mock.Anything, mock.Anything).Return(
		&settingssvc.ListBundlesResponse{Bundles: []*settingsmsg.Bundle{
			{Id: guestRoleID, Name: "user-light"},
			{Id: vaultRoleID, Name: "user-light-vault"},
		}}, nil)
	roleService.On("ListRoleAssignments", mock.Anything, mock.Anything, mock.Anything).Return(
		&settingssvc.ListRoleAssignmentsResponse{Assignments: []*settingsmsg.UserRoleAssignment{
			{Id: "guest-assignment-id", AccountUuid: userOpaqueID, RoleId: guestRoleID},
		}}, nil)
	roleService.On("AssignRoleToUser", mock.Anything, mock.Anything, mock.Anything).Return(
		&settingssvc.AssignRoleToUserResponse{Assignment: &settingsmsg.UserRoleAssignment{
			Id:          "vault-assignment-id",
			AccountUuid: userOpaqueID,
			RoleId:      vaultRoleID,
		}}, nil)

	ra := oidcRoleAssigner{Options: Options{
		logger:          log.NopLogger(),
		gatewaySelector: gatewaySelector,
		roleService:     roleService,
		rolesClaim:      "roles",
		guestMapping:    []config.RoleMapping{{ClaimValue: "spaceadmin", RoleName: "user-light-vault"}},
		serviceAccount:  config.ServiceAccount{ServiceAccountID: "service-account", ServiceAccountSecret: "secret"},
	}}

	user := &cs3user.User{Id: &cs3user.UserId{OpaqueId: userOpaqueID}}
	claims := map[string]interface{}{"roles": []interface{}{"spaceadmin"}}

	if _, err := ra.UpdateUserRoleAssignment(context.Background(), user, claims, "user-light"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	roleService.AssertCalled(t, "AssignRoleToUser", mock.Anything, mock.MatchedBy(func(req *settingssvc.AssignRoleToUserRequest) bool {
		return req.GetAccountUuid() == userOpaqueID && req.GetRoleId() == vaultRoleID
	}), mock.Anything)
	roleService.AssertNotCalled(t, "AssignRoleToUser", mock.Anything, mock.MatchedBy(func(req *settingssvc.AssignRoleToUserRequest) bool {
		return req.GetRoleId() == guestRoleID
	}), mock.Anything)
}

// TestUpdateUserRoleAssignmentFailsOnUnknownGuestRole verifies that a guest role mapping that
// points to a role which does not exist fails the role assignment instead of silently assigning
// the guest role.
func TestUpdateUserRoleAssignmentFailsOnUnknownGuestRole(t *testing.T) {
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

	// listing the roles needs a service account context
	gatewayClient.On("Authenticate", mock.Anything, mock.Anything).Return(&gateway.AuthenticateResponse{
		Status: &rpc.Status{Code: rpc.Code_CODE_OK},
		Token:  "service-token",
	}, nil)

	roleService := &graphmocks.RoleService{}
	roleService.On("ListRoles", mock.Anything, mock.Anything, mock.Anything).Return(
		&settingssvc.ListBundlesResponse{Bundles: []*settingsmsg.Bundle{{Id: "user-light-id", Name: "user-light"}}}, nil)

	ra := oidcRoleAssigner{Options: Options{
		logger:          log.NopLogger(),
		gatewaySelector: gatewaySelector,
		roleService:     roleService,
		rolesClaim:      "roles",
		guestMapping:    []config.RoleMapping{{ClaimValue: "spaceadmin", RoleName: "user-light-vault"}},
		serviceAccount:  config.ServiceAccount{ServiceAccountID: "service-account", ServiceAccountSecret: "secret"},
	}}

	user := &cs3user.User{Id: &cs3user.UserId{OpaqueId: "user-1"}}
	claims := map[string]interface{}{"roles": []interface{}{"spaceadmin"}}

	_, err := ra.UpdateUserRoleAssignment(context.Background(), user, claims, "user-light")
	if err == nil || !strings.Contains(err.Error(), "unknown role 'user-light-vault'") {
		t.Fatalf("expected an error for the guest role mapping to an unknown role, got %v", err)
	}
	roleService.AssertNotCalled(t, "AssignRoleToUser", mock.Anything, mock.Anything, mock.Anything)

	// a guest that no mapping matches gets the guest role, which must exist too
	claims = map[string]interface{}{"roles": []interface{}{"user"}}
	_, err = ra.UpdateUserRoleAssignment(context.Background(), user, claims, "no-such-guest-role")
	if err == nil || !strings.Contains(err.Error(), "unknown role 'no-such-guest-role'") {
		t.Fatalf("expected an error for an unknown guest role, got %v", err)
	}
	roleService.AssertNotCalled(t, "AssignRoleToUser", mock.Anything, mock.Anything, mock.Anything)
}

func TestGuestRoleMapping(t *testing.T) {
	tests := []struct {
		name    string
		driver  string
		pairs   []string
		want    []config.RoleMapping
		wantErr bool
	}{
		{
			name:   "empty",
			driver: "default",
			pairs:  nil,
			want:   []config.RoleMapping{},
		},
		{
			name:   "pairs keep their order",
			driver: "oidc",
			pairs:  []string{"spaceadmin:user-light-vault", "admin:guest-admin"},
			want: []config.RoleMapping{
				{ClaimValue: "spaceadmin", RoleName: "user-light-vault"},
				{ClaimValue: "admin", RoleName: "guest-admin"},
			},
		},
		{
			name:   "claim value may contain colons",
			driver: "oidc",
			pairs:  []string{"urn:ocis:spaceadmin:user-light-vault"},
			want:   []config.RoleMapping{{ClaimValue: "urn:ocis:spaceadmin", RoleName: "user-light-vault"}},
		},
		{name: "requires the oidc driver", driver: "default", pairs: []string{"spaceadmin:user-light-vault"}, wantErr: true},
		{name: "missing colon", driver: "oidc", pairs: []string{"spaceadmin"}, wantErr: true},
		{name: "missing claim value", driver: "oidc", pairs: []string{":user-light-vault"}, wantErr: true},
		{name: "missing role name", driver: "oidc", pairs: []string{"spaceadmin:"}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{
				RoleAssignment: config.RoleAssignment{Driver: tt.driver},
				MultiInstance:  config.MultiInstanceConfig{GuestRoleMapping: tt.pairs},
			}
			got, err := GuestRoleMapping(cfg)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected an error for %v", tt.pairs)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
		})
	}
}
