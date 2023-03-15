package userroles

import (
	"context"
	"errors"

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

	// To list roles and update assignment we need some elevated access to the settings service
	// prepare a new request context for that until we have service accounts
	newctx, err := ra.prepareAdminContext()
	if err != nil {
		ra.logger.Error().Err(err).Msg("Error creating admin context")
		return nil, err
	}

	claimValueToRoleID, err := ra.oidcClaimvaluesToRoleIDs(newctx)
	if err != nil {
		ra.logger.Error().Err(err).Msg("Error mapping claims to roles ids")
		return nil, err
	}

	roleIDsFromClaim := make([]string, 0, 1)
	ra.logger.Error().Interface("rolesclaim", claims[ra.rolesClaim]).Msg("Got ClaimRoles")
	claimRoles, ok := claims[ra.rolesClaim].([]interface{})
	if !ok {
		ra.logger.Error().Err(err).Msg("No roles in user claims.")
		return nil, err
	}
	for _, cri := range claimRoles {
		cr, ok := cri.(string)
		if !ok {
			err := errors.New("invalid role in claims")
			ra.logger.Error().Err(err).Interface("claim value", cri).Msg("Is not a valid string.")
			return nil, err
		}
		id, ok := claimValueToRoleID[cr]
		if !ok {
			ra.logger.Error().Str("role", cr).Msg("Skipping unmaped role from claims.")
			continue
		}
		roleIDsFromClaim = append(roleIDsFromClaim, id)
	}
	ra.logger.Error().Interface("roleIDs", roleIDsFromClaim).Msg("Mapped roles from claim")

	switch len(roleIDsFromClaim) {
	default:
		err := errors.New("too many roles found in claims")
		ra.logger.Error().Err(err).Msg("Only one role per user is allowed.")
		return nil, err
	case 0:
		err := errors.New("no role in claim, maps to a ocis role")
		ra.logger.Error().Err(err).Msg("")
		return nil, err
	case 1:
		// exactly one mapping. This is right
	}

	assignedRoles, err := loadRolesIDs(newctx, user.GetId().GetOpaqueId(), ra.roleService)
	if err != nil {
		ra.logger.Error().Err(err).Msgf("Could not load roles")
		return nil, err
	}
	if len(assignedRoles) > 1 {
		err := errors.New("too many roles assigned")
		ra.logger.Error().Err(err).Msg("The user has too many roles assigned")
		return nil, err
	}
	ra.logger.Error().Interface("assignedRoleIds", assignedRoles).Msg("Currently assigned roles")
	if len(assignedRoles) == 0 || (assignedRoles[0] != roleIDsFromClaim[0]) {
		if _, err = ra.roleService.AssignRoleToUser(newctx, &settingssvc.AssignRoleToUserRequest{
			AccountUuid: user.GetId().GetOpaqueId(),
			RoleId:      roleIDsFromClaim[0],
		}); err != nil {
			ra.logger.Error().Err(err).Msg("Role assignment failed")
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

func (ra oidcRoleAssigner) oidcClaimvaluesToRoleIDs(ctx context.Context) (map[string]string, error) {
	roleClaimToID := map[string]string{}
	// Get all roles to find the role IDs.
	// TODO: we need to cache this. Roles IDs change rarely and this is a pretty expensiveV call
	req := &settingssvc.ListBundlesRequest{}
	res, err := ra.roleService.ListRoles(ctx, req)
	if err != nil {
		ra.logger.Error().Err(err).Msg("Failed to list all roles")
		return roleClaimToID, err
	}

	for _, role := range res.Bundles {
		ra.logger.Error().Str("role", role.Name).Str("id", role.Id).Msg("Got Role")
		roleClaim, ok := ra.roleMapping[role.Name]
		if !ok {
			err := errors.New("Incomplete role mapping")
			ra.logger.Error().Err(err).Str("role", role.Name).Msg("Role not mapped to a claim value")
			return roleClaimToID, err
		}
		roleClaimToID[roleClaim] = role.Id
	}
	return roleClaimToID, nil
}
