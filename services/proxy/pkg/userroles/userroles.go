package userroles

import (
	"context"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	cs3 "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	ocismiddleware "github.com/owncloud/ocis/v2/ocis-pkg/middleware"
	settingssvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/settings/v0"
	"github.com/owncloud/ocis/v2/services/proxy/pkg/config"
	settingsService "github.com/owncloud/ocis/v2/services/settings/pkg/service/v0"
	"github.com/owncloud/reva/v2/pkg/rgrpc/todo/pool"
	"go-micro.dev/v4/metadata"
)

// UserRoleAssigner allows providing different implementations for how users get their default roles
// assigned by the proxy during authentication
type UserRoleAssigner interface {
	// UpdateUserRoleAssignment is called by the account resolver middleware. It updates the user's role assignment
	// based on the user's (OIDC) claims. It adds the user's roles to the opaque data of the cs3.User struct
	UpdateUserRoleAssignment(ctx context.Context, user *cs3.User, claims map[string]interface{}, forceRole string) (*cs3.User, error)
	// ApplyUserRole can be called by proxy middlewares, it looks up the user's roles and adds them
	// the users "roles" key in the user's opaque data
	ApplyUserRole(ctx context.Context, user *cs3.User) (*cs3.User, error)
}

// Options defines the available options for this package.
type Options struct {
	gatewaySelector pool.Selectable[gateway.GatewayAPIClient]
	roleService     settingssvc.RoleService
	rolesClaim      string
	roleMapping     []config.RoleMapping
	serviceAccount  config.ServiceAccount
	logger          log.Logger
}

// Option defines a single option function.
type Option func(o *Options)

// WithLogger configure the logger
func WithLogger(l log.Logger) Option {
	return func(o *Options) {
		o.logger = l
	}
}

// WithRoleService sets the roleservice instance to use
func WithRoleService(rs settingssvc.RoleService) Option {
	return func(o *Options) {
		o.roleService = rs
	}
}

// WithRolesClaim sets the OIDC claim for looking up role names
func WithRolesClaim(claim string) Option {
	return func(o *Options) {
		o.rolesClaim = claim
	}
}

// WithRoleMapping configures the map of ocis role names to claims values
func WithRoleMapping(roleMap []config.RoleMapping) Option {
	return func(o *Options) {
		o.roleMapping = roleMap
	}
}

// WithRevaGatewaySelector set the gatewaySelector option
func WithRevaGatewaySelector(selectable pool.Selectable[gateway.GatewayAPIClient]) Option {
	return func(o *Options) {
		o.gatewaySelector = selectable
	}
}

// WithServiceAccount configures the service account creator to use
func WithServiceAccount(c config.ServiceAccount) Option {
	return func(o *Options) {
		o.serviceAccount = c
	}
}

// bootstrapDefaultRoleIfNeeded assigns the default "User" role (or "User Light"
// for guests) when the user has no role assignment yet, so that users who only
// ever authenticate without going through UpdateUserRoleAssignment (e.g. via an
// app password) still end up with a role. See
// https://github.com/owncloud/ocis/issues/1825 for context.
func bootstrapDefaultRoleIfNeeded(ctx context.Context, user *cs3.User, roleIDs []string, roleService settingssvc.RoleService, logger log.Logger) ([]string, error) {
	if len(roleIDs) != 0 {
		return roleIDs, nil
	}
	if user.Id.Type != cs3.UserType_USER_TYPE_PRIMARY && user.Id.Type != cs3.UserType_USER_TYPE_GUEST {
		return roleIDs, nil
	}

	roleID := settingsService.BundleUUIDRoleUser
	if user.Id.Type == cs3.UserType_USER_TYPE_GUEST {
		roleID = settingsService.BundleUUIDRoleUserLight
	}
	logger.Info().Str("userid", user.Id.OpaqueId).Msg("user has no role assigned, assigning default user role")
	ctx = metadata.Set(ctx, ocismiddleware.AccountID, user.Id.OpaqueId)
	if _, err := roleService.AssignRoleToUser(ctx, &settingssvc.AssignRoleToUserRequest{
		AccountUuid: user.Id.OpaqueId,
		RoleId:      roleID,
	}); err != nil {
		logger.Error().Err(err).Msg("Could not add default role")
		return nil, err
	}
	return append(roleIDs, roleID), nil
}

// loadRolesIDs returns the role-ids assigned to an user
func loadRolesIDs(ctx context.Context, opaqueUserID string, rs settingssvc.RoleService) ([]string, error) {
	req := &settingssvc.ListRoleAssignmentsRequest{AccountUuid: opaqueUserID}
	assignmentResponse, err := rs.ListRoleAssignments(ctx, req)

	if err != nil {
		return nil, err
	}

	roleIDs := make([]string, 0)

	for _, assignment := range assignmentResponse.Assignments {
		roleIDs = append(roleIDs, assignment.RoleId)
	}

	return roleIDs, nil
}
