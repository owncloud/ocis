package userroles

import (
	"context"
	"encoding/json"
	"errors"
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

// newRoleAssignerFixture builds an oidcRoleAssigner whose settings service knows the
// roles in knownRoles and reports the user as already holding assignedRoleID. Holding
// the role the assigner is about to pick short-circuits the re-assignment branch, which
// keeps these tests on the claim-to-role resolution they are about rather than on the
// personal-space reconciliation that follows it.
func newRoleAssignerFixture(t *testing.T, opts Options, knownRoles map[string]string, assignedRoleID string) oidcRoleAssigner {
	t.Helper()

	// the role-name cache is a package global with a 5 minute TTL; reset it so tests
	// cannot leak role maps into one another
	roleNameToID.lock.Lock()
	roleNameToID.roleNameToID = nil
	roleNameToID.lastRead = time.Time{}
	roleNameToID.lock.Unlock()

	gatewayClient := &cs3mocks.GatewayAPIClient{}
	selectorName := "GatewaySelector" + t.Name()
	gatewaySelector := pool.GetSelector[gateway.GatewayAPIClient](
		selectorName,
		"com.owncloud.api.gateway",
		func(cc grpc.ClientConnInterface) gateway.GatewayAPIClient {
			return gatewayClient
		},
	)
	t.Cleanup(func() { pool.RemoveSelector(selectorName + "com.owncloud.api.gateway") })

	gatewayClient.On("Authenticate", mock.Anything, mock.Anything).Return(&gateway.AuthenticateResponse{
		Status: &rpc.Status{Code: rpc.Code_CODE_OK},
		Token:  "service-token",
	}, nil)

	bundles := make([]*settingsmsg.Bundle, 0, len(knownRoles))
	for name, id := range knownRoles {
		bundles = append(bundles, &settingsmsg.Bundle{Id: id, Name: name})
	}
	roleService := &graphmocks.RoleService{}
	roleService.On("ListRoles", mock.Anything, mock.Anything, mock.Anything).Return(
		&settingssvc.ListBundlesResponse{Bundles: bundles}, nil)
	roleService.On("ListRoleAssignments", mock.Anything, mock.Anything, mock.Anything).Return(
		&settingssvc.ListRoleAssignmentsResponse{Assignments: []*settingsmsg.UserRoleAssignment{
			{Id: "assignment-id", AccountUuid: "user-1", RoleId: assignedRoleID},
		}}, nil)
	// Mocked so that an unexpected re-assignment shows up as a failed assertion rather
	// than as a panic in an unrelated place.
	roleService.On("AssignRoleToUser", mock.Anything, mock.Anything, mock.Anything).Return(
		&settingssvc.AssignRoleToUserResponse{Assignment: &settingsmsg.UserRoleAssignment{
			Id: "assignment-id", AccountUuid: "user-1", RoleId: assignedRoleID,
		}}, nil)

	opts.logger = log.NopLogger()
	opts.gatewaySelector = gatewaySelector
	opts.roleService = roleService
	opts.serviceAccount = config.ServiceAccount{ServiceAccountID: "service-account", ServiceAccountSecret: "secret"}
	return oidcRoleAssigner{Options: opts}
}

func assignedRoleFromOpaque(t *testing.T, user *cs3user.User) []string {
	t.Helper()
	entry := user.GetOpaque().GetMap()["roles"]
	if entry == nil {
		t.Fatal("the user's opaque data carries no roles entry")
	}
	var got []string
	if err := json.Unmarshal(entry.GetValue(), &got); err != nil {
		t.Fatalf("could not decode the roles opaque entry: %v", err)
	}
	return got
}

// TestUpdateUserRoleAssignmentFallsBackToDefaultRoleOnUnmatchedClaim covers the case the
// issue names directly: the claim is present but its value matches no role mapping.
func TestUpdateUserRoleAssignmentFallsBackToDefaultRoleOnUnmatchedClaim(t *testing.T) {
	const defaultRoleID = "user-light-id"

	ra := newRoleAssignerFixture(t,
		Options{
			rolesClaim:  "roles",
			roleMapping: []config.RoleMapping{{RoleName: "admin", ClaimValue: "ocisAdmin"}},
			defaultRole: "user-light",
		},
		map[string]string{"admin": "admin-id", "user-light": defaultRoleID},
		defaultRoleID,
	)

	user := &cs3user.User{Id: &cs3user.UserId{OpaqueId: "user-1"}}
	got, err := ra.UpdateUserRoleAssignment(context.Background(), user, map[string]interface{}{"roles": "somethingElse"}, "")
	if err != nil {
		t.Fatalf("expected the default role to be applied, got error: %v", err)
	}
	if roles := assignedRoleFromOpaque(t, got); len(roles) != 1 || roles[0] != defaultRoleID {
		t.Fatalf("expected the default role id %q, got %v", defaultRoleID, roles)
	}
}

// TestUpdateUserRoleAssignmentFallsBackToDefaultRoleWithoutAnyClaim is the case that
// actually reproduces the reported bug. Users federated into the IDP from an external
// directory have no role claim at all, so extractRoles fails outright - a fallback that
// only covered "claim present but unmatched" would not fix the report.
func TestUpdateUserRoleAssignmentFallsBackToDefaultRoleWithoutAnyClaim(t *testing.T) {
	const defaultRoleID = "user-light-id"

	ra := newRoleAssignerFixture(t,
		Options{
			rolesClaim:  "roles",
			roleMapping: []config.RoleMapping{{RoleName: "admin", ClaimValue: "ocisAdmin"}},
			defaultRole: "user-light",
		},
		map[string]string{"admin": "admin-id", "user-light": defaultRoleID},
		defaultRoleID,
	)

	// no "roles" key whatsoever
	claims := map[string]interface{}{"sub": "abcd", "email": "federated@example.org"}

	user := &cs3user.User{Id: &cs3user.UserId{OpaqueId: "user-1"}}
	got, err := ra.UpdateUserRoleAssignment(context.Background(), user, claims, "")
	if err != nil {
		t.Fatalf("expected the default role to be applied, got error: %v", err)
	}
	if roles := assignedRoleFromOpaque(t, got); len(roles) != 1 || roles[0] != defaultRoleID {
		t.Fatalf("expected the default role id %q, got %v", defaultRoleID, roles)
	}
}

// TestUpdateUserRoleAssignmentPrefersAMatchingMappingOverTheDefaultRole guards the
// fallback against swallowing normal operation.
func TestUpdateUserRoleAssignmentPrefersAMatchingMappingOverTheDefaultRole(t *testing.T) {
	const adminRoleID = "admin-id"

	ra := newRoleAssignerFixture(t,
		Options{
			rolesClaim:  "roles",
			roleMapping: []config.RoleMapping{{RoleName: "admin", ClaimValue: "ocisAdmin"}},
			defaultRole: "user-light",
		},
		map[string]string{"admin": adminRoleID, "user-light": "user-light-id"},
		adminRoleID,
	)

	user := &cs3user.User{Id: &cs3user.UserId{OpaqueId: "user-1"}}
	got, err := ra.UpdateUserRoleAssignment(context.Background(), user, map[string]interface{}{"roles": "ocisAdmin"}, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if roles := assignedRoleFromOpaque(t, got); len(roles) != 1 || roles[0] != adminRoleID {
		t.Fatalf("expected the mapped role id %q, got %v", adminRoleID, roles)
	}
}

// TestUpdateUserRoleAssignmentWithoutDefaultRoleReturnsErrNoRoleAssigned pins the
// behaviour when no default role is configured. It must stay an error - this change does
// not hand out a role to anybody who was previously refused - but it must be one the
// caller can recognise, which is what turns the opaque 500 into a 403.
func TestUpdateUserRoleAssignmentWithoutDefaultRoleReturnsErrNoRoleAssigned(t *testing.T) {
	ra := newRoleAssignerFixture(t,
		Options{
			rolesClaim:  "roles",
			roleMapping: []config.RoleMapping{{RoleName: "admin", ClaimValue: "ocisAdmin"}},
		},
		map[string]string{"admin": "admin-id"},
		"admin-id",
	)

	user := &cs3user.User{Id: &cs3user.UserId{OpaqueId: "user-1"}}
	for name, claims := range map[string]map[string]interface{}{
		"unmatched claim": {"roles": "somethingElse"},
		"no claim at all": {"sub": "abcd"},
	} {
		t.Run(name, func(t *testing.T) {
			_, err := ra.UpdateUserRoleAssignment(context.Background(), user, claims, "")
			if !errors.Is(err, ErrNoRoleAssigned) {
				t.Fatalf("expected ErrNoRoleAssigned, got %v", err)
			}
		})
	}
}

// TestUpdateUserRoleAssignmentReportsAnUnknownDefaultRole separates a deployment
// mistake from a user without a role: only one of the two is fixed by editing the
// configuration, so they must not return the same error.
func TestUpdateUserRoleAssignmentReportsAnUnknownDefaultRole(t *testing.T) {
	ra := newRoleAssignerFixture(t,
		Options{
			rolesClaim:  "roles",
			roleMapping: []config.RoleMapping{{RoleName: "admin", ClaimValue: "ocisAdmin"}},
			defaultRole: "no-such-role",
		},
		map[string]string{"admin": "admin-id"},
		"admin-id",
	)

	user := &cs3user.User{Id: &cs3user.UserId{OpaqueId: "user-1"}}
	_, err := ra.UpdateUserRoleAssignment(context.Background(), user, map[string]interface{}{"roles": "somethingElse"}, "")
	if err == nil {
		t.Fatal("expected an error for a default role that does not exist")
	}
	if errors.Is(err, ErrNoRoleAssigned) {
		t.Fatalf("a misconfigured default role must not be reported as ErrNoRoleAssigned, got %v", err)
	}
	if !strings.Contains(err.Error(), "no-such-role") {
		t.Fatalf("the error should name the offending role, got %v", err)
	}
}
