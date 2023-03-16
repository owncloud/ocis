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

// NewOIDCRoleAssigner returns an implemenation of the UserRoleAssigner interface
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
	claimValueToRoleID, err := ra.oidcClaimvaluesToRoleIDs()
	if err != nil {
		logger.Error().Err(err).Msg("Error mapping claims to roles ids")
		return nil, err
	}

	roleIDsFromClaim := make([]string, 0, 1)
	claimRoles, ok := claims[ra.rolesClaim].([]interface{})
	if !ok {
		logger.Error().Err(err).Str("rolesClaim", ra.rolesClaim).Msg("No roles in user claims")
		return nil, err
	}
	logger.Debug().Str("rolesClaim", ra.rolesClaim).Interface("rolesInClaim", claims[ra.rolesClaim]).Msg("got roles in claim")
	for _, cri := range claimRoles {
		cr, ok := cri.(string)
		if !ok {
			err := errors.New("invalid role in claims")
			logger.Error().Err(err).Interface("claimValue", cri).Msg("Is not a valid string.")
			return nil, err
		}
		id, ok := claimValueToRoleID[cr]
		if !ok {
			logger.Debug().Str("role", cr).Msg("No mapping for claim role. Skipped.")
			continue
		}
		roleIDsFromClaim = append(roleIDsFromClaim, id)
	}
	logger.Debug().Interface("roleIDs", roleIDsFromClaim).Msg("Mapped claim roles to roleids")

	switch len(roleIDsFromClaim) {
	default:
		err := errors.New("too many roles found in claims")
		logger.Error().Err(err).Msg("Only one role per user is allowed.")
		return nil, err
	case 0:
		err := errors.New("no role in claim, maps to a ocis role")
		logger.Error().Err(err).Msg("")
		return nil, err
	case 1:
		// exactly one mapping. This is right
	}
	assignedRoles, err := loadRolesIDs(ctx, user.GetId().GetOpaqueId(), ra.roleService)
	if err != nil {
		logger.Error().Err(err).Msg("Could not load roles")
		return nil, err
	}
	if len(assignedRoles) > 1 {
		err := errors.New("too many roles assigned")
		logger.Error().Err(err).Msg("The user has too many roles assigned")
		return nil, err
	}
	logger.Debug().Interface("assignedRoleIds", assignedRoles).Msg("Currently assigned roles")
	if len(assignedRoles) == 0 || (assignedRoles[0] != roleIDsFromClaim[0]) {
		logger.Debug().Interface("assignedRoleIds", assignedRoles).Interface("newRoleIds", roleIDsFromClaim).Msg("Updating role assignment for user")
		newctx, err := ra.prepareAdminContext()
		if err != nil {
			logger.Error().Err(err).Msg("Error creating admin context")
			return nil, err
		}
		if _, err = ra.roleService.AssignRoleToUser(newctx, &settingssvc.AssignRoleToUserRequest{
			AccountUuid: user.GetId().GetOpaqueId(),
			RoleId:      roleIDsFromClaim[0],
		}); err != nil {
			logger.Error().Err(err).Msg("Role assignment failed")
			return nil, err
		}
	}

	user.Opaque = utils.AppendJSONToOpaque(user.Opaque, "roles", roleIDsFromClaim)
	return user, nil
}

// ApplyUserRole it looks up the user's role in the settings service and adds it
// user's opaque data
func (ra oidcRoleAssigner) ApplyUserRole(ctx context.Context, user *cs3.User) (*cs3.User, error) {
	roleIDs, err := loadRolesIDs(ctx, user.Id.OpaqueId, ra.roleService)
	if err != nil {
		ra.logger.Error().Err(err).Msgf("Could not load roles")
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

type roleClaimToIDCache struct {
	roleClaimToID map[string]string
	lastRead      time.Time
	lock          sync.RWMutex
}

var roleClaimToID roleClaimToIDCache

func (ra oidcRoleAssigner) oidcClaimvaluesToRoleIDs() (map[string]string, error) {
	cacheTTL := 5 * time.Minute
	roleClaimToID.lock.RLock()

	if !roleClaimToID.lastRead.IsZero() && time.Since(roleClaimToID.lastRead) < cacheTTL {
		defer roleClaimToID.lock.RUnlock()
		return roleClaimToID.roleClaimToID, nil
	}
	ra.logger.Debug().Msg("refreshing roles ids")

	// cache needs Refresh get a write lock
	roleClaimToID.lock.RUnlock()
	roleClaimToID.lock.Lock()
	defer roleClaimToID.lock.Unlock()

	// check again, another goroutine might have updated while we "upgraded" the lock
	if !roleClaimToID.lastRead.IsZero() && time.Since(roleClaimToID.lastRead) < cacheTTL {
		return roleClaimToID.roleClaimToID, nil
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
		roleClaim, ok := ra.roleMapping[role.Name]
		if !ok {
			err := errors.New("Incomplete role mapping")
			ra.logger.Error().Err(err).Str("role", role.Name).Msg("Role not mapped to a claim value")
			return map[string]string{}, err
		}
		newIDs[roleClaim] = role.Id
	}
	ra.logger.Debug().Interface("roleMap", newIDs).Msg("Claim Role to role ID map")
	roleClaimToID.roleClaimToID = newIDs
	roleClaimToID.lastRead = time.Now()
	return roleClaimToID.roleClaimToID, nil
}
