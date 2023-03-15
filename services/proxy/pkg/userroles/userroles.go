package userroles

import (
	"context"

	cs3 "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	settingssvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/settings/v0"
)

//go:generate mockery --name=UserRoleAssigner

// UserRoleAssigner allows to provide different implemenation for how users get their default roles
// assigned by the proxy during authentication
type UserRoleAssigner interface {
	// UpdateUserRoleAssignment is called by the account resolver middleware. It updates the user's role assignment
	// based on the user's (OIDC) claims. It adds the user's roles to the opaque data of the cs3.User struct
	UpdateUserRoleAssignment(ctx context.Context, user *cs3.User, claims map[string]interface{}) (*cs3.User, error)
	// ApplyUserRole can be called by proxy middlewares, it looks up the user's roles and adds them
	// the users "roles" key in the user's opaque data
	ApplyUserRole(ctx context.Context, user *cs3.User) (*cs3.User, error)
}

type Options struct {
	roleService settingssvc.RoleService
	rolesClaim  string
	logger      log.Logger
}

type Option func(o *Options)

func WithLogger(l log.Logger) Option {
	return func(o *Options) {
		o.logger = l
	}
}

func WithRoleService(rs settingssvc.RoleService) Option {
	return func(o *Options) {
		o.roleService = rs
	}
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
