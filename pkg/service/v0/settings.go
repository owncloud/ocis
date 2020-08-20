package svc

import settings "github.com/owncloud/ocis-settings/pkg/proto/v0"

const (
	// BundleUUIDRoleAdmin represents the admin role
	BundleUUIDRoleAdmin = "71881883-1768-46bd-a24d-a356a2afdf7f"

	// BundleUUIDRoleUser represents the user role.
	BundleUUIDRoleUser = "d7beeea8-8ff4-406b-8fb6-ab2dd81e6b11"

	// BundleUUIDRoleGuest represents the guest role.
	BundleUUIDRoleGuest = "38071a68-456a-4553-846a-fa67bf5596cc"
)

// generateBundlesDefaultRoles bootstraps the default roles.
func generateBundlesDefaultRoles() []*settings.Bundle {
	return []*settings.Bundle{
		generateBundleAdminRole(),
		generateBundleUserRole(),
		generateBundleGuestRole(),
	}
}

func generateBundleAdminRole() *settings.Bundle {
	return &settings.Bundle{
		Id:          BundleUUIDRoleAdmin,
		Name:        "admin",
		Type:        settings.Bundle_TYPE_ROLE,
		Extension:   "ocis-roles",
		DisplayName: "Admin",
		Resource: &settings.Resource{
			Type: settings.Resource_TYPE_SYSTEM,
		},
		Settings: []*settings.Setting{},
	}
}

func generateDefaultPermissionsRequests() []settings.AddSettingToBundleRequest {
	return []settings.AddSettingToBundleRequest{
		// ADMIN permissions
		{
			BundleId: BundleUUIDRoleAdmin,
			Setting: &settings.Setting{
				Id:          "560c6270-b29a-49de-8a1b-b655aa8b9c84",
				Name:        "read-settings-all",
				DisplayName: "Permission to read values for all settings",
				Resource: &settings.Resource{
					Type: settings.Resource_TYPE_SETTING,
					Id:   "all",
				},
				Value: &settings.Setting_PermissionValue{
					PermissionValue: &settings.Permission{
						Operation:  settings.Permission_OPERATION_READ,
						Constraint: settings.Permission_CONSTRAINT_ALL,
					},
				},
			},
		},
		{
			BundleId: BundleUUIDRoleAdmin,
			Setting: &settings.Setting{
				Id:          "4503cf00-8598-453d-8bd4-81ba552fd1fc",
				Name:        "create-settings-all",
				DisplayName: "Permission to create values for all settings",
				Resource: &settings.Resource{
					Type: settings.Resource_TYPE_SETTING,
					Id:   "all",
				},
				Value: &settings.Setting_PermissionValue{
					PermissionValue: &settings.Permission{
						Operation:  settings.Permission_OPERATION_CREATE,
						Constraint: settings.Permission_CONSTRAINT_ALL,
					},
				},
			},
		},
		{
			BundleId: BundleUUIDRoleAdmin,
			Setting: &settings.Setting{
				Id:          "08dda0ab-f087-4d9f-92f2-64f8f6c5a463",
				Name:        "update-settings-all",
				DisplayName: "Permission to update values for all settings",
				Resource: &settings.Resource{
					Type: settings.Resource_TYPE_SETTING,
					Id:   "all",
				},
				Value: &settings.Setting_PermissionValue{
					PermissionValue: &settings.Permission{
						Operation:  settings.Permission_OPERATION_UPDATE,
						Constraint: settings.Permission_CONSTRAINT_ALL,
					},
				},
			},
		},
		{
			BundleId: BundleUUIDRoleAdmin,
			Setting: &settings.Setting{
				Id:          "b996707a-e122-4490-b3ed-a3d22713692e",
				Name:        "delete-settings-all",
				DisplayName: "Permission to delete values for all settings",
				Resource: &settings.Resource{
					Type: settings.Resource_TYPE_SETTING,
					Id:   "all",
				},
				Value: &settings.Setting_PermissionValue{
					PermissionValue: &settings.Permission{
						Operation:  settings.Permission_OPERATION_DELETE,
						Constraint: settings.Permission_CONSTRAINT_ALL,
					},
				},
			},
		},
		// USER permissions
		{
			BundleId: BundleUUIDRoleUser,
			Setting: &settings.Setting{
				Id:          "fb1036a5-6356-4dd0-b4c6-90dc6f6e86b0",
				Name:        "read-settings-all",
				DisplayName: "Permission to read values for all settings",
				Resource: &settings.Resource{
					Type: settings.Resource_TYPE_SETTING,
					Id:   "all",
				},
				Value: &settings.Setting_PermissionValue{
					PermissionValue: &settings.Permission{
						Operation:  settings.Permission_OPERATION_READ,
						Constraint: settings.Permission_CONSTRAINT_OWN,
					},
				},
			},
		},
		{
			BundleId: BundleUUIDRoleUser,
			Setting: &settings.Setting{
				Id:          "8eae5d66-cc72-4b15-a7db-33c84dbaa305",
				Name:        "create-settings-all",
				DisplayName: "Permission to create values for all settings",
				Resource: &settings.Resource{
					Type: settings.Resource_TYPE_SETTING,
					Id:   "all",
				},
				Value: &settings.Setting_PermissionValue{
					PermissionValue: &settings.Permission{
						Operation:  settings.Permission_OPERATION_CREATE,
						Constraint: settings.Permission_CONSTRAINT_OWN,
					},
				},
			},
		},
		{
			BundleId: BundleUUIDRoleUser,
			Setting: &settings.Setting{
				Id:          "d818ba99-1c81-4773-a2f3-89cecdb19b92",
				Name:        "update-settings-all",
				DisplayName: "Permission to update values for all settings",
				Resource: &settings.Resource{
					Type: settings.Resource_TYPE_SETTING,
					Id:   "all",
				},
				Value: &settings.Setting_PermissionValue{
					PermissionValue: &settings.Permission{
						Operation:  settings.Permission_OPERATION_UPDATE,
						Constraint: settings.Permission_CONSTRAINT_OWN,
					},
				},
			},
		},
		{
			BundleId: BundleUUIDRoleUser,
			Setting: &settings.Setting{
				Id:          "40273d13-9bdc-4234-8b76-56a6572d2619",
				Name:        "delete-settings-all",
				DisplayName: "Permission to create values for all settings",
				Resource: &settings.Resource{
					Type: settings.Resource_TYPE_SETTING,
					Id:   "all",
				},
				Value: &settings.Setting_PermissionValue{
					PermissionValue: &settings.Permission{
						Operation:  settings.Permission_OPERATION_DELETE,
						Constraint: settings.Permission_CONSTRAINT_OWN,
					},
				},
			},
		},
		// GUEST user
		{
			BundleId: BundleUUIDRoleGuest,
			Setting: &settings.Setting{
				Id:          "5fb4ea7f-f351-4dd7-a9af-4550c44e2362",
				Name:        "read-settings-all",
				DisplayName: "Permission to read values for all settings",
				Resource: &settings.Resource{
					Type: settings.Resource_TYPE_SETTING,
					Id:   "all",
				},
				Value: &settings.Setting_PermissionValue{
					PermissionValue: &settings.Permission{
						Operation:  settings.Permission_OPERATION_READ,
						Constraint: settings.Permission_CONSTRAINT_OWN,
					},
				},
			},
		},
	}
}

func generateBundleUserRole() *settings.Bundle {
	return &settings.Bundle{
		Id:          BundleUUIDRoleUser,
		Name:        "user",
		Type:        settings.Bundle_TYPE_ROLE,
		Extension:   "ocis-roles",
		DisplayName: "User",
		Resource: &settings.Resource{
			Type: settings.Resource_TYPE_SYSTEM,
		},
		Settings: []*settings.Setting{},
	}
}

func generateBundleGuestRole() *settings.Bundle {
	return &settings.Bundle{
		Id:          BundleUUIDRoleGuest,
		Name:        "guest",
		Type:        settings.Bundle_TYPE_ROLE,
		Extension:   "ocis-roles",
		DisplayName: "Guest",
		Resource: &settings.Resource{
			Type: settings.Resource_TYPE_SYSTEM,
		},
		Settings: []*settings.Setting{},
	}
}
