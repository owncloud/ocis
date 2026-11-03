package userroles

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	cs3user "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	"github.com/owncloud/ocis/v2/ocis-pkg/middleware"
	"github.com/owncloud/ocis/v2/ocis-pkg/oidc"
	"github.com/owncloud/ocis/v2/ocis-pkg/rolemapping"
	"github.com/owncloud/ocis/v2/ocis-pkg/shared"
	settingssvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/settings/v0"
	"github.com/owncloud/ocis/v2/services/graph/pkg/identity"
	"github.com/owncloud/ocis/v2/services/proxy/pkg/config"
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
	m := rolemapping.NewMatcher(mappingValue)
	for cr := range claimRoles {
		if m.Matches(cr) {
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
	roleNamesToRoleIDs, err := ra.roleNamesToRoleIDs(ctx, roleCacheTTL)
	if err != nil {
		logger.Error().Err(err).Msg("Error mapping role names to role ids")
		return nil, err
	}

	// the account resolver passes the guest role as overwriteRole for multi-instance guests
	if overwriteRole != "" {
		overwriteRole, roleNamesToRoleIDs, err = ra.guestRoleName(ctx, claims, overwriteRole, roleNamesToRoleIDs)
		if err != nil {
			logger.Error().Err(err).Msg("Error mapping guest role")
			return nil, err
		}
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
		var oldRole string
		if len(assignedRoles) > 0 {
			oldRole = assignedRoles[0]
		}
		newctx, err := ra.prepareAdminContext(ctx)
		if err != nil {
			logger.Error().Err(err).Msg("Error creating admin context")
			return nil, err
		}
		assignResp, err := ra.roleService.AssignRoleToUser(newctx, &settingssvc.AssignRoleToUserRequest{
			AccountUuid: userID,
			RoleId:      roleIDFromClaim,
		})
		if err != nil {
			logger.Error().Err(err).Msg("Role assignment failed")
			return nil, err
		}

		userID := user.GetId().GetOpaqueId()
		client, err := ra.gatewaySelector.Next()
		if err != nil {
			return nil, err
		}

		// TODO: check if it's ctx or newctx
		canCreateDrives, err := utils.CheckPermission(revactx.ContextSetUser(ctx, user), "Drives.Create", client)
		if err != nil {
			// The permission could not be determined. Fail closed: leave the personal
			// space untouched rather than disabling (trashing) it on an indeterminate
			// result, and revert the role assignment so the user is not left in a
			// half-applied state. The role transition is retried on the next login.
			logger.Error().Any("userID", userID).Err(err).Msg("could not determine Drives.Create permission, reverting role assignment and leaving personal space unchanged")
			ra.revertRoleAssignment(newctx, userID, oldRole, assignResp.GetAssignment().GetId())
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

// GuestRoleMapping reads the 'claimvalue:rolename' pairs of OCIS_MULTI_INSTANCE_GUEST_ROLE_MAPPING,
// see rolemapping.ParsePair.
func GuestRoleMapping(cfg *config.Config) ([]config.RoleMapping, error) {
	pairs := cfg.MultiInstance.GuestRoleMapping
	if len(pairs) > 0 && cfg.RoleAssignment.Driver != "oidc" {
		return nil, fmt.Errorf("the guest role mapping requires the 'oidc' role assignment driver, not '%s'", cfg.RoleAssignment.Driver)
	}
	// the role assigner only applies the mapping to users that get the guest role
	if len(pairs) > 0 && cfg.MultiInstance.GuestRole == "" {
		return nil, errors.New("the guest role mapping requires OCIS_MULTI_INSTANCE_GUEST_ROLE")
	}

	mappings := make([]config.RoleMapping, 0, len(pairs))
	for _, pair := range pairs {
		claimValue, roleName, err := rolemapping.ParsePair(pair)
		if err != nil {
			return nil, err
		}
		mappings = append(mappings, config.RoleMapping{ClaimValue: claimValue, RoleName: roleName})
	}
	return mappings, nil
}

// guestRoleName returns the role name of a multi-instance guest and the role name to role id map
// that knows this role. A role that the cached map does not know can be a new role, so
// guestRoleName reads the roles again one time. An unknown role is an error: do not assign an
// empty role id.
func (ra oidcRoleAssigner) guestRoleName(ctx context.Context, claims map[string]interface{}, guestRole string, roleNamesToRoleIDs map[string]string) (string, map[string]string, error) {
	resolve := func(roleNamesToRoleIDs map[string]string) (string, error) {
		roleName := guestRole
		if len(ra.guestMapping) > 0 {
			var err error
			if roleName, err = ra.mapGuestRole(claims, guestRole, roleNamesToRoleIDs); err != nil {
				return "", err
			}
		}
		if _, ok := roleNamesToRoleIDs[roleName]; !ok {
			return "", fmt.Errorf("unknown role '%s'", roleName)
		}
		return roleName, nil
	}

	roleName, err := resolve(roleNamesToRoleIDs)
	if err == nil {
		return roleName, roleNamesToRoleIDs, nil
	}
	// the minimum age limits the reads for a role that does not exist
	roleNamesToRoleIDs, rerr := ra.roleNamesToRoleIDs(ctx, roleCacheMinAge)
	if rerr != nil {
		return "", nil, rerr
	}
	roleName, err = resolve(roleNamesToRoleIDs)
	if err != nil {
		return "", nil, err
	}
	return roleName, roleNamesToRoleIDs, nil
}

// mapGuestRole returns the role of the first guest role mapping that matches a value of the
// user's roles claim, or guestRole when nothing matches. A mapping to a role that does not exist
// is an error: it can only be a misconfiguration, and falling back to guestRole could give the
// guest more permissions than the mapped role was meant to.
func (ra oidcRoleAssigner) mapGuestRole(claims map[string]interface{}, guestRole string, roleNamesToRoleIDs map[string]string) (string, error) {
	claimRoles, err := extractRoles(ra.rolesClaim, claims)
	if err != nil || len(claimRoles) == 0 {
		// a guest without roles claim is valid, log at debug level to help find a wrong claim name
		ra.logger.Debug().Err(err).Str("claim", ra.rolesClaim).Msg("no roles in claim, the guest keeps the guest role")
		return guestRole, nil
	}

	for _, mapping := range ra.guestMapping {
		if !matchesClaimMapping(mapping.ClaimValue, claimRoles) {
			continue
		}
		if _, ok := roleNamesToRoleIDs[mapping.RoleName]; !ok {
			return "", fmt.Errorf("guest role mapping '%s' points to the unknown role '%s'", mapping.ClaimValue, mapping.RoleName)
		}
		return mapping.RoleName, nil
	}
	return guestRole, nil
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

func (ra oidcRoleAssigner) prepareAdminContext(ctx context.Context) (context.Context, error) {
	gatewayClient, err := ra.gatewaySelector.Next()
	if err != nil {
		ra.logger.Error().Err(err).Msg("could not select next gateway client")
		return nil, err
	}
	newctx, err := utils.GetServiceUserContextWithContext(ctx, gatewayClient, ra.serviceAccount.ServiceAccountID, ra.serviceAccount.ServiceAccountSecret)
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

const (
	// roleCacheTTL is the maximum age of the cached role name to role id map
	roleCacheTTL = 5 * time.Minute
	// roleCacheMinAge is the age that the cached map must have before a read for an unknown role
	roleCacheMinAge = 30 * time.Second
)

// roleNamesToRoleIDs returns the cached role name to role id map. When the map is older than
// maxAge, roleNamesToRoleIDs reads the roles from the settings service.
func (ra oidcRoleAssigner) roleNamesToRoleIDs(ctx context.Context, maxAge time.Duration) (map[string]string, error) {
	roleNameToID.lock.RLock()

	if !roleNameToID.lastRead.IsZero() && time.Since(roleNameToID.lastRead) < maxAge {
		defer roleNameToID.lock.RUnlock()
		return roleNameToID.roleNameToID, nil
	}
	ra.logger.Debug().Msg("refreshing roles ids")

	// cache needs Refresh get a write lock
	roleNameToID.lock.RUnlock()
	roleNameToID.lock.Lock()
	defer roleNameToID.lock.Unlock()

	// check again, another goroutine might have updated while we "upgraded" the lock
	if !roleNameToID.lastRead.IsZero() && time.Since(roleNameToID.lastRead) < maxAge {
		return roleNameToID.roleNameToID, nil
	}

	// Get all roles to find the role IDs.
	// To list roles we need some elevated access to the settings service
	// prepare a new request context for that until we have service accounts
	newctx, err := ra.prepareAdminContext(ctx)
	if err != nil {
		ra.logger.Error().Err(err).Msg("Error creating admin context")
		return nil, err
	}

	req := &settingssvc.ListBundlesRequest{}
	res, err := ra.roleService.ListRoles(newctx, req)
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

// revertRoleAssignment best-effort restores the user's previous role after a
// failed reconciliation, so the user is not left with the new role applied but
// the personal space unreconciled. A user has exactly one role, so this means
// re-assigning the previous role, or removing the new assignment if there was
// none. Failures are only logged; the next login reconciles idempotently.
func (ra oidcRoleAssigner) revertRoleAssignment(ctx context.Context, userID, oldRoleID, newAssignmentID string) {
	if oldRoleID != "" {
		if _, err := ra.roleService.AssignRoleToUser(ctx, &settingssvc.AssignRoleToUserRequest{
			AccountUuid: userID,
			RoleId:      oldRoleID,
		}); err != nil {
			ra.logger.Error().Any("userID", userID).Str("roleID", oldRoleID).Err(err).Msg("could not revert role assignment to previous role")
		}
		return
	}
	if _, err := ra.roleService.RemoveRoleFromUser(ctx, &settingssvc.RemoveRoleFromUserRequest{
		Id: newAssignmentID,
	}); err != nil {
		ra.logger.Error().Any("userID", userID).Str("assignmentID", newAssignmentID).Err(err).Msg("could not revert role assignment by removing the new assignment")
	}
}
