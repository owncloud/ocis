package svc

import "github.com/owncloud/ocis-settings/pkg/proto/v0"

func (g Service) hasPermission(
	roleIDs []string,
	resource *proto.Resource,
	operations []proto.Permission_Operation,
	constraint proto.Permission_Constraint,
) bool {
	permissions, err := g.manager.ListPermissionsByResource(resource, roleIDs)
	if err != nil {
		g.logger.Debug().Err(err).
			Str("resource-type", resource.Type.String()).
			Str("resource-id", resource.Id).
			Msg("permissions could not be loaded for resource")
		return false
	}
	permissions = getFilteredPermissionsByOperations(permissions, operations)
	return isConstraintFulfilled(permissions, constraint)
}

// filterPermissionsByOperations returns the subset of the given permissions, where at least one of the given operations is fulfilled.
func getFilteredPermissionsByOperations(permissions []*proto.Permission, operations []proto.Permission_Operation) []*proto.Permission {
	var filteredPermissions []*proto.Permission
	for _, permission := range permissions {
		if isAnyOperationFulfilled(permission, operations) {
			filteredPermissions = append(filteredPermissions, permission)
		}
	}
	return filteredPermissions
}

// isAnyOperationFulfilled checks if the permissions is about any of the operations
func isAnyOperationFulfilled(permission *proto.Permission, operations []proto.Permission_Operation) bool {
	for _, operation := range operations {
		if operation == permission.Operation {
			return true
		}
	}
	return false
}

// isConstraintFulfilled checks if one of the permissions has the same or a parent of the constraint.
// this is only a comparison on ENUM level. More sophisticated checks cannot happen here...
func isConstraintFulfilled(permissions []*proto.Permission, constraint proto.Permission_Constraint) bool {
	for _, permission := range permissions {
		// comparing enum by order is not a feasible solution, because `SHARED` is not a superset of `OWN`.
		if permission.Constraint == proto.Permission_CONSTRAINT_ALL {
			return true
		}
		return permission.Constraint != proto.Permission_CONSTRAINT_UNKNOWN && permission.Constraint == constraint
	}
	return false
}
