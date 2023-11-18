package validate

import (
	"github.com/go-playground/validator/v10"
	libregraph "github.com/owncloud/libre-graph-api-go"

	"golang.org/x/exp/slices"

	"github.com/owncloud/ocis/v2/services/graph/pkg/unifiedrole"
)

// initLibregraph initializes libregraph validation
func initLibregraph(v *validator.Validate) {
	driveItemInvite(v)
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

		totalRoles := len(driveItemInvite.Roles)
		totalActions := len(driveItemInvite.LibreGraphPermissionsActions)

		switch {
		case totalRoles != 0 && totalActions != 0:
			fallthrough
		case totalRoles == totalActions:
			sl.ReportError(driveItemInvite.Roles, "Roles", "Roles", "one_or_another", "")
			sl.ReportError(driveItemInvite.LibreGraphPermissionsActions, "LibreGraphPermissionsActions", "LibreGraphPermissionsActions", "one_or_another", "")
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

		for _, role := range driveItemInvite.Roles {
			if slices.Contains(availableRoles, role) {
				continue
			}

			sl.ReportError(driveItemInvite.Roles, "Roles", "Roles", "available_role", "")
		}

		for _, role := range driveItemInvite.LibreGraphPermissionsActions {
			if slices.Contains(availableActions, role) {
				continue
			}

			sl.ReportError(driveItemInvite.LibreGraphPermissionsActions, "LibreGraphPermissionsActions", "LibreGraphPermissionsActions", "available_action", "")
		}

	}, s)
}
