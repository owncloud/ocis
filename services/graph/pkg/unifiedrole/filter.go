package unifiedrole

import (
	"slices"

	libregraph "github.com/owncloud/libre-graph-api-go"
)

type (
	// RoleFilter is used to filter role collections
	RoleFilter func(r *libregraph.UnifiedRoleDefinition) bool

	// RoleFilterMatch defines the match behavior of a role filter
	RoleFilterMatch int
)

const (
	// RoleFilterMatchExact is the behavior for role filters that require an exact match
	RoleFilterMatchExact RoleFilterMatch = iota

	// RoleFilterMatchSome is the behavior for role filters that require some match
	RoleFilterMatchSome
)

// RoleFilterInvert inverts the provided role filter
func RoleFilterInvert(f RoleFilter) RoleFilter {
	return func(r *libregraph.UnifiedRoleDefinition) bool {
		return !f(r)
	}
}

// RoleFilterAll returns a role filter that matches all roles
func RoleFilterAll() RoleFilter {
	return func(r *libregraph.UnifiedRoleDefinition) bool {
		return true
	}
}

// RoleFilterIDs returns a role filter that matches the provided ids
// the filter is always OR!
func RoleFilterIDs(ids ...string) RoleFilter {
	return func(r *libregraph.UnifiedRoleDefinition) bool {
		return slices.Contains(ids, r.GetId())
	}
}

// RoleFilterPermission returns a role filter that matches the provided condition and actions pair
func RoleFilterPermission(matchBehavior RoleFilterMatch, condition string, wantActions ...string) RoleFilter {
	return func(r *libregraph.UnifiedRoleDefinition) bool {
		for _, permission := range r.GetRolePermissions() {
			if permission.GetCondition() != condition {
				continue
			}

			givenActions := permission.GetAllowedResourceActions()

			switch {
			case matchBehavior == RoleFilterMatchExact && slices.Equal(givenActions, wantActions):
				return true
			case matchBehavior == RoleFilterMatchSome:
				matches := 0

				for _, action := range givenActions {
					if !slices.Contains(wantActions, action) {
						break
					}

					matches++
				}

				return len(givenActions) == matches
			}
		}

		return false
	}
}

// filterRoles filters the provided roles by the provided filter
func filterRoles(roles []*libregraph.UnifiedRoleDefinition, filter RoleFilter) []*libregraph.UnifiedRoleDefinition {
	return slices.DeleteFunc(
		slices.Clone(roles),
		func(r *libregraph.UnifiedRoleDefinition) bool {
			return !filter(r)
		},
	)
}
