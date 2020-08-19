package svc

import "github.com/owncloud/ocis-settings/pkg/proto/v0"

func (g Service) hasPermission(
	assignments []*proto.UserRoleAssignment,
	resource *proto.Resource,
	operation proto.Permission_Operation,
	constraint proto.Permission_Constraint,
) bool {
	for index := range assignments {
		if g.isAllowedByRole(assignments[index], resource, operation, constraint) {
			return true
		}
	}
	return false
}

func (g Service) isAllowedByRole(
	assignment *proto.UserRoleAssignment,
	resource *proto.Resource,
	operation proto.Permission_Operation,
	constraint proto.Permission_Constraint,
) bool {
	role, err := g.manager.ReadBundle(assignment.RoleId)
	if err != nil {
		g.logger.Err(err).Str("bundle", assignment.RoleId).Msg("Failed to fetch role")
		return false
	}
	for _, setting := range role.Settings {
		if _, ok := setting.Value.(*proto.Setting_PermissionValue); ok {
			value := setting.Value.(*proto.Setting_PermissionValue).PermissionValue
			if resource.Type == setting.Resource.Type &&
				resource.Id == setting.Resource.Id &&
				operation == value.Operation &&
				isConstraintMatch(constraint, value.Constraint) {
				return true
			}
		}
	}
	return false
}

// isConstraintMatch checks if the `given` constraint is the same or a superset of the `required` constraint.
// this is only a comparison on ENUM level. this is not a check about the appropriate constraint for a resource.
func isConstraintMatch(given, required proto.Permission_Constraint) bool {
	// comparing enum by order is not a feasible solution, because `SHARED` is not a superset of `OWN`.
	if given == proto.Permission_CONSTRAINT_ALL {
		return true
	}
	return given != proto.Permission_CONSTRAINT_UNKNOWN && given == required
}
