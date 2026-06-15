package userroles

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"sync"
	"time"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	cs3user "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	permissions "github.com/cs3org/go-cs3apis/cs3/permissions/v1beta1"
	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	"github.com/owncloud/ocis/v2/ocis-pkg/middleware"
	"github.com/owncloud/ocis/v2/ocis-pkg/oidc"
	"github.com/owncloud/ocis/v2/ocis-pkg/shared"
	settingssvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/settings/v0"
	"github.com/owncloud/ocis/v2/services/graph/pkg/identity"
	revactx "github.com/owncloud/reva/v2/pkg/ctx"
	"github.com/owncloud/reva/v2/pkg/utils"
	"go-micro.dev/v4/metadata"
)

type oidcRoleAssigner struct {
	Options
}

// NewOIDCRoleAssigner returns an implementation of the UserRoleAssigner interface
func NewOIDCRoleAssigner(opts ...Option) UserRoleAssigner {
	opt := Options{}
	for _, o := range opts {
		o(&opt)
	}

	return oidcRoleAssigner{
		Options: opt,
	}
}

func extractRoles(rolesClaim string, claims map[string]interface{}) (map[string]struct{}, error) {

	claimRoles := map[string]struct{}{}
	// happy path
	value, _ := claims[rolesClaim].(string)
	if value != "" {
		claimRoles[value] = struct{}{}
		return claimRoles, nil
	}

	claim, err := oidc.WalkSegments(oidc.SplitWithEscaping(rolesClaim, ".", "\\"), claims)
	if err != nil {
		return nil, err
	}

	switch v := claim.(type) {
	case []string:
		for _, cr := range v {
			claimRoles[cr] = struct{}{}
		}
	case []interface{}:
		for _, cri := range v {
			cr, ok := cri.(string)
			if !ok {
				err := errors.New("invalid role in claims")
				return nil, err
			}

			claimRoles[cr] = struct{}{}
		}
	case string:
		claimRoles[v] = struct{}{}
	default:
		return nil, errors.New("no roles in user claims")
	}

	return claimRoles, nil
}

// matchesClaimMapping returns true if the provided mapping pattern matches at least
// one of the values present in claimRoles. It supports:
// - exact match when ClaimValue is a literal equal to a claim value
// - regex match when ClaimValue is a regex pattern (e.g. "ocis-user-.*")
// The regex is matched against the entire claim value, not a substring.
func matchesClaimMapping(mappingValue string, claimRoles map[string]struct{}) bool {
	if _, ok := claimRoles[mappingValue]; ok {
		return true
	}

	rx, err := regexp.Compile("^(?:" + mappingValue + ")$")
	if err != nil {
		return false
	}
	for cr := range claimRoles {
		if rx.MatchString(cr) {
			return true
		}
	}
	return false
}

// UpdateUserRoleAssignment assigns the role "User" to the supplied user. Unless the user
// already has a different role assigned.
func (ra oidcRoleAssigner) UpdateUserRoleAssignment(ctx context.Context, user *cs3user.User, claims map[string]interface{}, overwriteRole string) (*cs3user.User, error) {
	userID := user.GetId().GetOpaqueId()
	logger := ra.logger.SubloggerWithRequestID(ctx).With().Str("userid", userID).Logger()
	roleNamesToRoleIDs, err := ra.roleNamesToRoleIDs()
	if err != nil {
		logger.Error().Err(err).Msg("Error mapping role names to role ids")
		return nil, err
	}

	roleIDFromClaim := roleNamesToRoleIDs[overwriteRole]
	if overwriteRole == "" {
		claimRoles, err := extractRoles(ra.rolesClaim, claims)
		if err != nil {
			logger.Error().Err(err).Str("Claim", ra.rolesClaim).Interface("claims", claims).Msg("Error mapping role names to role ids")
			return nil, err
		}

		if len(claimRoles) == 0 {
			err := errors.New("no roles set in claim")
			logger.Error().Err(err).Msg("")
			return nil, err
		}

		// the roleMapping config is supposed to have the role mappings ordered from the highest privileged role
		// down to the lowest privileged role. Since ocis currently only can handle a single role assignment we
		// pick the highest privileged role that matches a value from the claims
		for _, mapping := range ra.Options.roleMapping {
			if matchesClaimMapping(mapping.ClaimValue, claimRoles) {
				logger.Debug().Str("ocisRole", mapping.RoleName).Str("role id", roleNamesToRoleIDs[mapping.RoleName]).Msg("first matching role")
				roleIDFromClaim = roleNamesToRoleIDs[mapping.RoleName]
				break
			}
		}

		if roleIDFromClaim == "" {
			err := errors.New("no role in claim maps to an ocis role")
			logger.Error().Err(err).Msg("")
			return nil, err
		}
	}

	assignedRoles, err := loadRolesIDs(ctx, userID, ra.roleService)
	if err != nil {
		logger.Error().Err(err).Msg("Could not load roles")
		return nil, err
	}
	if len(assignedRoles) > 1 {
		logger.Error().Str("userID", userID).Int("numRoles", len(assignedRoles)).Msg("The user has too many roles assigned")
	}
	logger.Debug().Interface("assignedRoleIds", assignedRoles).Msg("Currently assigned roles")

	if len(assignedRoles) != 1 || (assignedRoles[0] != roleIDFromClaim) {
		logger.Debug().Interface("assignedRoleIds", assignedRoles).Interface("newRoleId", roleIDFromClaim).Msg("Updating role assignment for user")
		newctx, err := ra.prepareAdminContext()
		if err != nil {
			logger.Error().Err(err).Msg("Error creating admin context")
			return nil, err
		}
		if _, err = ra.roleService.AssignRoleToUser(newctx, &settingssvc.AssignRoleToUserRequest{
			AccountUuid: userID,
			RoleId:      roleIDFromClaim,
		}); err != nil {
			logger.Error().Err(err).Msg("Role assignment failed")
			return nil, err
		}

		userID := user.GetId().GetOpaqueId()
		client, err := ra.gatewaySelector.Next()
		if err != nil {
			return nil, err
		}

		canCreateDrives, err := ra.checkPermission("Drives.Create", user, client)
		if err != nil {
			// The permission could not be determined. Fail closed and leave the
			// personal space untouched rather than disabling (trashing) it on an
			// indeterminate result. The role transition is retried on the next login.
			logger.Error().Any("userID", userID).Err(err).Msg("could not determine Drives.Create permission, leaving personal space unchanged")
			return nil, err
		}
		if canCreateDrives {
			libreUser := identity.CreateUserModelFromCS3(user)
			err = shared.RestorePersonalSpace(newctx, client, libreUser.GetId())
			if err != nil {
				logger.Error().Any("userID", userID).Err(err).Msg("can't ensure the personal space")
				return nil, err
			}
		} else {
			err := shared.DisablePersonalSpace(newctx, client, userID)
			if err != nil {
				logger.Error().Any("userID", userID).Err(err).Msg("can't disable the personal space")
				return nil, err
			}
		}
	}

	user.Opaque = utils.AppendJSONToOpaque(user.Opaque, "roles", []string{roleIDFromClaim})
	return user, nil
}

// ApplyUserRole it looks up the user's role in the settings service and adds it
// user's opaque data
func (ra oidcRoleAssigner) ApplyUserRole(ctx context.Context, user *cs3user.User) (*cs3user.User, error) {
	roleIDs, err := loadRolesIDs(ctx, user.Id.OpaqueId, ra.roleService)
	if err != nil {
		ra.logger.Error().Err(err).Msg("Could not load roles")
		return nil, err
	}

	user.Opaque = utils.AppendJSONToOpaque(user.Opaque, "roles", roleIDs)
	return user, nil
}

func (ra oidcRoleAssigner) prepareAdminContext() (context.Context, error) {
	gatewayClient, err := ra.gatewaySelector.Next()
	if err != nil {
		ra.logger.Error().Err(err).Msg("could not select next gateway client")
		return nil, err
	}
	newctx, err := utils.GetServiceUserContext(ra.serviceAccount.ServiceAccountID, gatewayClient, ra.serviceAccount.ServiceAccountSecret)
	if err != nil {
		ra.logger.Error().Err(err).Msg("Error preparing request context for provisioning role assignments.")
		return nil, err
	}

	newctx = metadata.Set(newctx, middleware.AccountID, ra.serviceAccount.ServiceAccountID)
	return newctx, nil
}

type roleNameToIDCache struct {
	roleNameToID map[string]string
	lastRead     time.Time
	lock         sync.RWMutex
}

var roleNameToID roleNameToIDCache

func (ra oidcRoleAssigner) roleNamesToRoleIDs() (map[string]string, error) {
	cacheTTL := 5 * time.Minute
	roleNameToID.lock.RLock()

	if !roleNameToID.lastRead.IsZero() && time.Since(roleNameToID.lastRead) < cacheTTL {
		defer roleNameToID.lock.RUnlock()
		return roleNameToID.roleNameToID, nil
	}
	ra.logger.Debug().Msg("refreshing roles ids")

	// cache needs Refresh get a write lock
	roleNameToID.lock.RUnlock()
	roleNameToID.lock.Lock()
	defer roleNameToID.lock.Unlock()

	// check again, another goroutine might have updated while we "upgraded" the lock
	if !roleNameToID.lastRead.IsZero() && time.Since(roleNameToID.lastRead) < cacheTTL {
		return roleNameToID.roleNameToID, nil
	}

	// Get all roles to find the role IDs.
	// To list roles we need some elevated access to the settings service
	// prepare a new request context for that until we have service accounts
	ctx, err := ra.prepareAdminContext()
	if err != nil {
		ra.logger.Error().Err(err).Msg("Error creating admin context")
		return nil, err
	}

	req := &settingssvc.ListBundlesRequest{}
	res, err := ra.roleService.ListRoles(ctx, req)
	if err != nil {
		ra.logger.Error().Err(err).Msg("Failed to list all roles")
		return map[string]string{}, err
	}

	newIDs := map[string]string{}
	for _, role := range res.Bundles {
		ra.logger.Debug().Str("role", role.Name).Str("id", role.Id).Msg("Got Role")
		newIDs[role.Name] = role.Id
	}
	ra.logger.Debug().Interface("roleMap", newIDs).Msg("Role Name to role ID map")
	roleNameToID.roleNameToID = newIDs
	roleNameToID.lastRead = time.Now()
	return roleNameToID.roleNameToID, nil
}

// checkPermission reports whether the user holds the given permission.
//
// It distinguishes three outcomes that must not be collapsed into a single
// boolean: the permission is granted (true, nil), the permission is
// authoritatively denied (false, nil), or the permission could not be
// determined (false, err). The last case happens when the gateway/settings
// call fails at the transport level or returns a non-OK status other than
// PERMISSION_DENIED. Callers must treat a non-nil error as "unknown" and fail
// closed: an indeterminate result must never be mistaken for a deliberate
// denial, because the caller acts destructively (disabling the personal space)
// on denial.
func (ra oidcRoleAssigner) checkPermission(perm string, user *cs3user.User, gwc gateway.GatewayAPIClient) (bool, error) {
	resp, err := gwc.CheckPermission(revactx.ContextSetUser(context.Background(), user), &permissions.CheckPermissionRequest{
		SubjectRef: &permissions.SubjectReference{
			Spec: &permissions.SubjectReference_UserId{
				UserId: user.GetId(),
			},
		},
		Permission: perm,
	})
	if err != nil {
		return false, err
	}
	switch code := resp.GetStatus().GetCode(); code {
	case rpc.Code_CODE_OK:
		return true, nil
	case rpc.Code_CODE_PERMISSION_DENIED:
		return false, nil
	default:
		return false, fmt.Errorf("permission check for %q returned non-authoritative status %q: %s", perm, code, resp.GetStatus().GetMessage())
	}
}
