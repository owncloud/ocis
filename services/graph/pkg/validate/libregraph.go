package validate

import (
	"slices"

	"github.com/go-playground/validator/v10"
	libregraph "github.com/owncloud/libre-graph-api-go"

	"github.com/owncloud/ocis/v2/services/graph/pkg/unifiedrole"
)

// initLibregraph initializes libregraph validation
func initLibregraph(v *validator.Validate) {
	driveItemInvite(v)
	permission(v)
}

// driveItemInvite validates libregraph.DriveItemInvite
func driveItemInvite(v *validator.Validate) {
	s := libregraph.DriveItemInvite{}

	v.RegisterStructValidationMapRules(map[string]string{
		"Recipients":         "min=1",
		"Roles":              "max=1",
		"ExpirationDateTime": "omitnil,gt",
	}, s)

	v.RegisterStructValidation(func(sl validator.StructLevel) {
		driveItemInvite := sl.Current().Interface().(libregraph.DriveItemInvite)

		rolesAndActions(sl, driveItemInvite.Roles, driveItemInvite.LibreGraphPermissionsActions, false)

	}, s)
}

// permission validates libregraph.Permission
func permission(v *validator.Validate) {
	s := libregraph.Permission{}

	v.RegisterStructValidationMapRules(map[string]string{
		"Roles": "max=1",
	}, s)
	v.RegisterStructValidation(func(sl validator.StructLevel) {
		permission := sl.Current().Interface().(libregraph.Permission)

		if _, ok := permission.GetIdOk(); ok {
			sl.ReportError(permission.Id, "Id", "Id", "readonly", "")
		}

		if _, ok := permission.GetHasPasswordOk(); ok {
			sl.ReportError(permission.HasPassword, "hasPassword", "HasPassword", "readonly", "")
		}

		rolesAndActions(sl, permission.Roles, permission.LibreGraphPermissionsActions, true)
	}, s)
}

func rolesAndActions(sl validator.StructLevel, roles, actions []string, allowEmpty bool) {
	totalRoles := len(roles)
	totalActions := len(actions)

	switch {
	case allowEmpty && totalRoles == 0 && totalActions == 0:
		break
	case totalRoles != 0 && totalActions != 0:
		fallthrough
	case totalRoles == totalActions:
		sl.ReportError(roles, "Roles", "Roles", "one_or_another", "")
		sl.ReportError(actions, "LibreGraphPermissionsActions", "LibreGraphPermissionsActions", "one_or_another", "")
	}

	var availableRoles []string
	var availableActions []string
	for _, definition := range append(
		unifiedrole.GetBuiltinRoleDefinitionList(true),
		unifiedrole.GetBuiltinRoleDefinitionList(false)...,
	) {
		if slices.Contains(availableRoles, definition.GetId()) {
			continue
		}

		availableRoles = append(availableRoles, definition.GetId())

		for _, permission := range definition.GetRolePermissions() {
			for _, action := range permission.GetAllowedResourceActions() {
				if slices.Contains(availableActions, action) {
					continue
				}

				availableActions = append(availableActions, action)
			}
		}
	}

	for _, role := range roles {
		if slices.Contains(availableRoles, role) {
			continue
		}

		sl.ReportError(roles, "Roles", "Roles", "available_role", "")
	}

	for _, role := range actions {
		if slices.Contains(availableActions, role) {
			continue
		}

		sl.ReportError(actions, "LibreGraphPermissionsActions", "LibreGraphPermissionsActions", "available_action", "")
	}
}
