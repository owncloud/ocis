package validate

import (
	"slices"

	"github.com/go-playground/validator/v10"
	libregraph "github.com/owncloud/libre-graph-api-go"

	"github.com/owncloud/ocis/v2/services/graph/pkg/unifiedrole"
)

// initLibregraph initializes libregraph validation
func initLibregraph(v *validator.Validate) {
	for _, f := range []func(*validator.Validate){
		libregraphDriveItemInvite,
		libregraphDriveRecipient,
		libregraphPermission,
	} {
		f(v)
	}
}

// libregraphDriveItemInvite validates libregraph.DriveItemInvite
func libregraphDriveItemInvite(v *validator.Validate) {
	s := libregraph.DriveItemInvite{}

	v.RegisterStructValidationMapRules(map[string]string{
		"Recipients":         "len=1,dive",
		"Roles":              "max=1",
		"ExpirationDateTime": "omitnil,gt",
	}, s)

	v.RegisterStructValidation(func(sl validator.StructLevel) {
		driveItemInvite := sl.Current().Interface().(libregraph.DriveItemInvite)

		rolesAndActions(sl, driveItemInvite.Roles, driveItemInvite.LibreGraphPermissionsActions, false)
	}, s)
}

// libregraphDriveRecipient validates libregraph.DriveRecipient
func libregraphDriveRecipient(v *validator.Validate) {
	v.RegisterStructValidationMapRules(map[string]string{
		"ObjectId":                "ne=",
		"LibreGraphRecipientType": "oneof=user group",
	}, libregraph.DriveRecipient{})
}

// libregraphPermission validates libregraph.Permission
func libregraphPermission(v *validator.Validate) {
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
		// fixMe: why twice!?
		// fixMe: should we consider all roles or only the ones that are enabled?
		unifiedrole.GetBuiltinRoleDefinitionList(),
		unifiedrole.GetBuiltinRoleDefinitionList()...,
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
