package svc

import settings "github.com/owncloud/ocis/settings/pkg/proto/v0"

const (
	// BundleUUIDRoleAdmin represents the admin role
	BundleUUIDRoleAdmin = "71881883-1768-46bd-a24d-a356a2afdf7f"

	// BundleUUIDRoleUser represents the user role.
	BundleUUIDRoleUser = "d7beeea8-8ff4-406b-8fb6-ab2dd81e6b11"

	// BundleUUIDRoleGuest represents the guest role.
	BundleUUIDRoleGuest = "38071a68-456a-4553-846a-fa67bf5596cc"

	// RoleManagementPermissionID is the hardcoded setting UUID for the role management permission
	RoleManagementPermissionID string = "a53e601e-571f-4f86-8fec-d4576ef49c62"
	// RoleManagementPermissionName is the hardcoded setting name for the role management permission
	RoleManagementPermissionName string = "role-management"

	// SettingsManagementPermissionID is the hardcoded setting UUID for the settings management permission
	SettingsManagementPermissionID string = "79e13b30-3e22-11eb-bc51-0b9f0bad9a58"
	// SettingsManagementPermissionName is the hardcoded setting name for the settings management permission
	SettingsManagementPermissionName string = "settings-management"
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

func generatePermissionRequests() []*settings.AddSettingToBundleRequest {
	return []*settings.AddSettingToBundleRequest{
		{
			BundleId: BundleUUIDRoleAdmin,
			Setting: &settings.Setting{
				Id:          RoleManagementPermissionID,
				Name:        RoleManagementPermissionName,
				DisplayName: "Role Management",
				Description: "This permission gives full access to everything that is related to role management.",
				Resource: &settings.Resource{
					Type: settings.Resource_TYPE_USER,
					Id:   "all",
				},
				Value: &settings.Setting_PermissionValue{
					PermissionValue: &settings.Permission{
						Operation:  settings.Permission_OPERATION_READWRITE,
						Constraint: settings.Permission_CONSTRAINT_ALL,
					},
				},
			},
		},
		{
			BundleId: BundleUUIDRoleAdmin,
			Setting: &settings.Setting{
				Id:          SettingsManagementPermissionID,
				Name:        SettingsManagementPermissionName,
				DisplayName: "Settings Management",
				Description: "This permission gives full access to everything that is related to settings management.",
				Resource: &settings.Resource{
					Type: settings.Resource_TYPE_USER,
					Id:   "all",
				},
				Value: &settings.Setting_PermissionValue{
					PermissionValue: &settings.Permission{
						Operation:  settings.Permission_OPERATION_READWRITE,
						Constraint: settings.Permission_CONSTRAINT_ALL,
					},
				},
			},
		},
	}
}
