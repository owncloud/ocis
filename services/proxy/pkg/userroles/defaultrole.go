package userroles

import (
	"context"

	cs3 "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	"github.com/owncloud/reva/v2/pkg/utils"
)

type defaultRoleAssigner struct {
	Options
}

// NewDefaultRoleAssigner returns an implementation of the UserRoleAssigner interface
func NewDefaultRoleAssigner(opts ...Option) UserRoleAssigner {
	opt := Options{}
	for _, o := range opts {
		o(&opt)
	}

	return defaultRoleAssigner{
		Options: opt,
	}
}

// UpdateUserRoleAssignment assigns the role "User" to the supplied user. Unless the user
// already has a different role assigned.
func (d defaultRoleAssigner) UpdateUserRoleAssignment(ctx context.Context, user *cs3.User, claims map[string]interface{}, _ string) (*cs3.User, error) {
	var roleIDs []string
	if user.Id.Type != cs3.UserType_USER_TYPE_LIGHTWEIGHT {
		var err error
		roleIDs, err = loadRolesIDs(ctx, user.Id.OpaqueId, d.roleService)
		if err != nil {
			d.logger.Error().Err(err).Msg("Could not load roles")
			return nil, err
		}

		roleIDs, err = bootstrapDefaultRoleIfNeeded(ctx, user, roleIDs, d.roleService, d.logger)
		if err != nil {
			return nil, err
		}
	}

	user.Opaque = utils.AppendJSONToOpaque(user.Opaque, "roles", roleIDs)
	return user, nil
}

// ApplyUserRole it looks up the user's role in the settings service and adds it
// user's opaque data
func (d defaultRoleAssigner) ApplyUserRole(ctx context.Context, user *cs3.User) (*cs3.User, error) {
	roleIDs, err := loadRolesIDs(ctx, user.Id.OpaqueId, d.roleService)
	if err != nil {
		d.logger.Error().Err(err).Msg("Could not load roles")
		return nil, err
	}

	roleIDs, err = bootstrapDefaultRoleIfNeeded(ctx, user, roleIDs, d.roleService, d.logger)
	if err != nil {
		return nil, err
	}

	user.Opaque = utils.AppendJSONToOpaque(user.Opaque, "roles", roleIDs)
	return user, nil
}
