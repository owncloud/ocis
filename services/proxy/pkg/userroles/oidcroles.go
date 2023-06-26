package userroles

import (
	"context"
	"errors"
	"sync"
	"time"

	cs3 "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	revactx "github.com/cs3org/reva/v2/pkg/ctx"
	"github.com/cs3org/reva/v2/pkg/utils"
	"github.com/owncloud/ocis/v2/ocis-pkg/middleware"
	settingssvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/settings/v0"
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

// UpdateUserRoleAssignment assigns the role "User" to the supplied user. Unless the user
// already has a different role assigned.
func (ra oidcRoleAssigner) UpdateUserRoleAssignment(ctx context.Context, user *cs3.User, claims map[string]interface{}) (*cs3.User, error) {
	logger := ra.logger.SubloggerWithRequestID(ctx).With().Str("userid", user.GetId().GetOpaqueId()).Logger()
	roleNamesToRoleIDs, err := ra.roleNamesToRoleIDs()
	if err != nil {
		logger.Error().Err(err).Msg("Error mapping role names to role ids")
		return nil, err
	}

	claimRolesRaw, ok := claims[ra.rolesClaim].([]interface{})
	if !ok {
		logger.Error().Str("rolesClaim", ra.rolesClaim).Msg("No roles in user claims")
		return nil, errors.New("no roles in user claims")
	}

	logger.Debug().Str("rolesClaim", ra.rolesClaim).Interface("rolesInClaim", claims[ra.rolesClaim]).Msg("got roles in claim")
	claimRoles := map[string]struct{}{}
	for _, cri := range claimRolesRaw {
		cr, ok := cri.(string)
		if !ok {
			err := errors.New("invalid role in claims")
			logger.Error().Err(err).Interface("claimValue", cri).Msg("Is not a valid string.")
			return nil, err
		}

		claimRoles[cr] = struct{}{}
	}

	if len(claimRoles) == 0 {
		err := errors.New("no roles set in claim")
		logger.Error().Err(err).Msg("")
		return nil, err
	}

	// the roleMapping config is supposed to have the role mappings ordered from the highest privileged role
	// down to the lowest privileged role. Since ocis currently only can handle a single role assignment we
	// pick the highest privileged role that matches a value from the claims
	roleIDFromClaim := ""
	for _, mapping := range ra.Options.roleMapping {
		if _, ok := claimRoles[mapping.ClaimValue]; ok {
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

	assignedRoles, err := loadRolesIDs(ctx, user.GetId().GetOpaqueId(), ra.roleService)
	if err != nil {
		logger.Error().Err(err).Msg("Could not load roles")
		return nil, err
	}
	if len(assignedRoles) > 1 {
		logger.Error().Str("userID", user.GetId().GetOpaqueId()).Int("numRoles", len(assignedRoles)).Msg("The user has too many roles assigned")
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
			AccountUuid: user.GetId().GetOpaqueId(),
			RoleId:      roleIDFromClaim,
		}); err != nil {
			logger.Error().Err(err).Msg("Role assignment failed")
			return nil, err
		}
	}

	user.Opaque = utils.AppendJSONToOpaque(user.Opaque, "roles", []string{roleIDFromClaim})
	return user, nil
}

// ApplyUserRole it looks up the user's role in the settings service and adds it
// user's opaque data
func (ra oidcRoleAssigner) ApplyUserRole(ctx context.Context, user *cs3.User) (*cs3.User, error) {
	roleIDs, err := loadRolesIDs(ctx, user.Id.OpaqueId, ra.roleService)
	if err != nil {
		ra.logger.Error().Err(err).Msg("Could not load roles")
		return nil, err
	}

	user.Opaque = utils.AppendJSONToOpaque(user.Opaque, "roles", roleIDs)
	return user, nil
}

func (ra oidcRoleAssigner) prepareAdminContext() (context.Context, error) {
	newctx := context.Background()
	autoProvisionUser, err := ra.autoProvsionCreator.GetAutoProvisionAdmin()
	if err != nil {
		return nil, err
	}
	token, err := ra.autoProvsionCreator.GetAutoProvisionAdminToken(newctx)
	if err != nil {
		ra.logger.Error().Err(err).Msg("Error generating token for provisioning role assignments.")
		return nil, err
	}
	newctx = revactx.ContextSetToken(newctx, token)
	newctx = metadata.Set(newctx, middleware.AccountID, autoProvisionUser.Id.OpaqueId)
	newctx = metadata.Set(newctx, middleware.RoleIDs, string(autoProvisionUser.Opaque.Map["roles"].Value))
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
