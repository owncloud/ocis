package defaults

import (
	settingsmsg "github.com/owncloud/ocis/protogen/gen/ocis/messages/settings/v0"
)

const (
	// BundleUUIDRoleAdmin represents the admin role
	BundleUUIDRoleAdmin = "71881883-1768-46bd-a24d-a356a2afdf7f"

	// BundleUUIDRoleSpaceAdmin represents the space admin role
	BundleUUIDRoleSpaceAdmin = "2aadd357-682c-406b-8874-293091995fdd"

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

	// SetSpaceQuotaPermissionID is the hardcoded setting UUID for the set space quota permission
	SetSpaceQuotaPermissionID string = "4e6f9709-f9e7-44f1-95d4-b762d27b7896"
	// SetSpaceQuotaPermissionName is the hardcoded setting name for the set space quota permission
	SetSpaceQuotaPermissionName string = "set-space-quota"

	// ListAllSpacesPermissionID is the hardcoded setting UUID for the list all spaces permission
	ListAllSpacesPermissionID string = "016f6ddd-9501-4a0a-8ebe-64a20ee8ec82"
	// ListAllSpacesPermissionName is the hardcoded setting name for the list all spaces permission
	ListAllSpacesPermissionName string = "list-all-spaces"

	// CreateSpacePermissionID is the hardcoded setting UUID for the create space permission
	CreateSpacePermissionID string = "79e13b30-3e22-11eb-bc51-0b9f0bad9a58"
	// CreateSpacePermissionName is the hardcoded setting name for the create space permission
	CreateSpacePermissionName string = "create-space"

	settingUUIDProfileLanguage = "aa8cfbe5-95d4-4f7e-a032-c3c01f5f062f"

	// AccountManagementPermissionID is the hardcoded setting UUID for the account management permission
	AccountManagementPermissionID string = "8e587774-d929-4215-910b-a317b1e80f73"
	// AccountManagementPermissionName is the hardcoded setting name for the account management permission
	AccountManagementPermissionName string = "account-management"
	// GroupManagementPermissionID is the hardcoded setting UUID for the group management permission
	GroupManagementPermissionID string = "522adfbe-5908-45b4-b135-41979de73245"
	// GroupManagementPermissionName is the hardcoded setting name for the group management permission
	GroupManagementPermissionName string = "group-management"
	// SelfManagementPermissionID is the hardcoded setting UUID for the self management permission
	SelfManagementPermissionID string = "e03070e9-4362-4cc6-a872-1c7cb2eb2b8e"
	// SelfManagementPermissionName is the hardcoded setting name for the self management permission
	SelfManagementPermissionName string = "self-management"
)

// GenerateBundlesDefaultRoles bootstraps the default roles.
func GenerateBundlesDefaultRoles() []*settingsmsg.Bundle {
	return []*settingsmsg.Bundle{
		generateBundleAdminRole(),
		generateBundleUserRole(),
		generateBundleGuestRole(),
		generateBundleProfileRequest(),
		generateBundleSpaceAdminRole(),
	}
}

func generateBundleAdminRole() *settingsmsg.Bundle {
	return &settingsmsg.Bundle{
		Id:          BundleUUIDRoleAdmin,
		Name:        "admin",
		Type:        settingsmsg.Bundle_TYPE_ROLE,
		Extension:   "ocis-roles",
		DisplayName: "Admin",
		Resource: &settingsmsg.Resource{
			Type: settingsmsg.Resource_TYPE_SYSTEM,
		},
		Settings: []*settingsmsg.Setting{
			{
				Id:          RoleManagementPermissionID,
				Name:        RoleManagementPermissionName,
				DisplayName: "Role Management",
				Description: "This permission gives full access to everything that is related to role management.",
				Resource: &settingsmsg.Resource{
					Type: settingsmsg.Resource_TYPE_USER,
					Id:   "all",
				},
				Value: &settingsmsg.Setting_PermissionValue{
					PermissionValue: &settingsmsg.Permission{
						Operation:  settingsmsg.Permission_OPERATION_READWRITE,
						Constraint: settingsmsg.Permission_CONSTRAINT_ALL,
					},
				},
			},
			{
				Id:          SettingsManagementPermissionID,
				Name:        SettingsManagementPermissionName,
				DisplayName: "Settings Management",
				Description: "This permission gives full access to everything that is related to settings management.",
				Resource: &settingsmsg.Resource{
					Type: settingsmsg.Resource_TYPE_USER,
					Id:   "all",
				},
				Value: &settingsmsg.Setting_PermissionValue{
					PermissionValue: &settingsmsg.Permission{
						Operation:  settingsmsg.Permission_OPERATION_READWRITE,
						Constraint: settingsmsg.Permission_CONSTRAINT_ALL,
					},
				},
			},
			{
				Id:          "7d81f103-0488-4853-bce5-98dcce36d649",
				Name:        "language-readwrite",
				DisplayName: "Permission to read and set the language (anyone)",
				Resource: &settingsmsg.Resource{
					Type: settingsmsg.Resource_TYPE_SETTING,
					Id:   settingUUIDProfileLanguage,
				},
				Value: &settingsmsg.Setting_PermissionValue{
					PermissionValue: &settingsmsg.Permission{
						Operation:  settingsmsg.Permission_OPERATION_READWRITE,
						Constraint: settingsmsg.Permission_CONSTRAINT_ALL,
					},
				},
			},
			{
				Id:          AccountManagementPermissionID,
				Name:        AccountManagementPermissionName,
				DisplayName: "Account Management",
				Description: "This permission gives full access to everything that is related to account management.",
				Resource: &settingsmsg.Resource{
					Type: settingsmsg.Resource_TYPE_USER,
					Id:   "all",
				},
				Value: &settingsmsg.Setting_PermissionValue{
					PermissionValue: &settingsmsg.Permission{
						Operation:  settingsmsg.Permission_OPERATION_READWRITE,
						Constraint: settingsmsg.Permission_CONSTRAINT_ALL,
					},
				},
			},
			{
				Id:          GroupManagementPermissionID,
				Name:        GroupManagementPermissionName,
				DisplayName: "Group Management",
				Description: "This permission gives full access to everything that is related to group management.",
				Resource: &settingsmsg.Resource{
					Type: settingsmsg.Resource_TYPE_GROUP,
					Id:   "all",
				},
				Value: &settingsmsg.Setting_PermissionValue{
					PermissionValue: &settingsmsg.Permission{
						Operation:  settingsmsg.Permission_OPERATION_READWRITE,
						Constraint: settingsmsg.Permission_CONSTRAINT_ALL,
					},
				},
			},
			{
				Id:          SetSpaceQuotaPermissionID,
				Name:        SetSpaceQuotaPermissionName,
				DisplayName: "Set Space Quota",
				Description: "This permission allows to manage space quotas.",
				Resource: &settingsmsg.Resource{
					Type: settingsmsg.Resource_TYPE_SYSTEM,
				},
				Value: &settingsmsg.Setting_PermissionValue{
					PermissionValue: &settingsmsg.Permission{
						Operation:  settingsmsg.Permission_OPERATION_READWRITE,
						Constraint: settingsmsg.Permission_CONSTRAINT_ALL,
					},
				},
			},
			{
				Id:          CreateSpacePermissionID,
				Name:        CreateSpacePermissionName,
				DisplayName: "Create Space",
				Description: "This permission allows to create new spaces.",
				Resource: &settingsmsg.Resource{
					Type: settingsmsg.Resource_TYPE_SYSTEM,
				},
				Value: &settingsmsg.Setting_PermissionValue{
					PermissionValue: &settingsmsg.Permission{
						Operation:  settingsmsg.Permission_OPERATION_READWRITE,
						Constraint: settingsmsg.Permission_CONSTRAINT_ALL,
					},
				},
			},
			{
				Id:          ListAllSpacesPermissionID,
				Name:        ListAllSpacesPermissionName,
				DisplayName: "List All Spaces",
				Description: "This permission allows list all spaces.",
				Resource: &settingsmsg.Resource{
					Type: settingsmsg.Resource_TYPE_SYSTEM,
				},
				Value: &settingsmsg.Setting_PermissionValue{
					PermissionValue: &settingsmsg.Permission{
						Operation:  settingsmsg.Permission_OPERATION_READ,
						Constraint: settingsmsg.Permission_CONSTRAINT_ALL,
					},
				},
			},
		},
	}
}

func generateBundleSpaceAdminRole() *settingsmsg.Bundle {
	return &settingsmsg.Bundle{
		Id:          BundleUUIDRoleSpaceAdmin,
		Name:        "spaceadmin",
		Type:        settingsmsg.Bundle_TYPE_ROLE,
		Extension:   "ocis-roles",
		DisplayName: "Space Admin",
		Resource: &settingsmsg.Resource{
			Type: settingsmsg.Resource_TYPE_SYSTEM,
		},
		Settings: []*settingsmsg.Setting{
			{
				Id:          SetSpaceQuotaPermissionID,
				Name:        SetSpaceQuotaPermissionName,
				DisplayName: "Set Space Quota",
				Description: "This permission allows to manage space quotas.",
				Resource: &settingsmsg.Resource{
					Type: settingsmsg.Resource_TYPE_SYSTEM,
				},
				Value: &settingsmsg.Setting_PermissionValue{
					PermissionValue: &settingsmsg.Permission{
						Operation:  settingsmsg.Permission_OPERATION_READWRITE,
						Constraint: settingsmsg.Permission_CONSTRAINT_ALL,
					},
				},
			},
			{
				Id:          CreateSpacePermissionID,
				Name:        CreateSpacePermissionName,
				DisplayName: "Create Space",
				Description: "This permission allows to create new spaces.",
				Resource: &settingsmsg.Resource{
					Type: settingsmsg.Resource_TYPE_SYSTEM,
				},
				Value: &settingsmsg.Setting_PermissionValue{
					PermissionValue: &settingsmsg.Permission{
						Operation:  settingsmsg.Permission_OPERATION_READWRITE,
						Constraint: settingsmsg.Permission_CONSTRAINT_ALL,
					},
				},
			},
			{
				Id:          ListAllSpacesPermissionID,
				Name:        ListAllSpacesPermissionName,
				DisplayName: "List All Spaces",
				Description: "This permission allows list all spaces.",
				Resource: &settingsmsg.Resource{
					Type: settingsmsg.Resource_TYPE_SYSTEM,
				},
				Value: &settingsmsg.Setting_PermissionValue{
					PermissionValue: &settingsmsg.Permission{
						Operation:  settingsmsg.Permission_OPERATION_READ,
						Constraint: settingsmsg.Permission_CONSTRAINT_ALL,
					},
				},
			},
			{
				Id:          "640e00d2-4df8-41bd-b1c2-9f30a01e0e99",
				Name:        "language-readwrite",
				DisplayName: "Permission to read and set the language (self)",
				Resource: &settingsmsg.Resource{
					Type: settingsmsg.Resource_TYPE_SETTING,
					Id:   settingUUIDProfileLanguage,
				},
				Value: &settingsmsg.Setting_PermissionValue{
					PermissionValue: &settingsmsg.Permission{
						Operation:  settingsmsg.Permission_OPERATION_READWRITE,
						Constraint: settingsmsg.Permission_CONSTRAINT_OWN,
					},
				},
			},
			{
				Id:          SelfManagementPermissionID,
				Name:        SelfManagementPermissionName,
				DisplayName: "Self Management",
				Description: "This permission gives access to self management.",
				Resource: &settingsmsg.Resource{
					Type: settingsmsg.Resource_TYPE_USER,
					Id:   "me",
				},
				Value: &settingsmsg.Setting_PermissionValue{
					PermissionValue: &settingsmsg.Permission{
						Operation:  settingsmsg.Permission_OPERATION_READWRITE,
						Constraint: settingsmsg.Permission_CONSTRAINT_OWN,
					},
				},
			},
			{
				Id:          CreateSpacePermissionID,
				Name:        CreateSpacePermissionName,
				DisplayName: "Create own Space",
				Description: "This permission allows to create a space owned by the current user.",
				Resource: &settingsmsg.Resource{
					Type: settingsmsg.Resource_TYPE_SYSTEM, // TODO resource type space? self? me? own?
				},
				Value: &settingsmsg.Setting_PermissionValue{
					PermissionValue: &settingsmsg.Permission{
						Operation:  settingsmsg.Permission_OPERATION_CREATE,
						Constraint: settingsmsg.Permission_CONSTRAINT_OWN,
					},
				},
			},
		},
	}
}

func generateBundleUserRole() *settingsmsg.Bundle {
	return &settingsmsg.Bundle{
		Id:          BundleUUIDRoleUser,
		Name:        "user",
		Type:        settingsmsg.Bundle_TYPE_ROLE,
		Extension:   "ocis-roles",
		DisplayName: "User",
		Resource: &settingsmsg.Resource{
			Type: settingsmsg.Resource_TYPE_SYSTEM,
		},
		Settings: []*settingsmsg.Setting{
			{
				Id:          "640e00d2-4df8-41bd-b1c2-9f30a01e0e99",
				Name:        "language-readwrite",
				DisplayName: "Permission to read and set the language (self)",
				Resource: &settingsmsg.Resource{
					Type: settingsmsg.Resource_TYPE_SETTING,
					Id:   settingUUIDProfileLanguage,
				},
				Value: &settingsmsg.Setting_PermissionValue{
					PermissionValue: &settingsmsg.Permission{
						Operation:  settingsmsg.Permission_OPERATION_READWRITE,
						Constraint: settingsmsg.Permission_CONSTRAINT_OWN,
					},
				},
			},
			{
				Id:          SelfManagementPermissionID,
				Name:        SelfManagementPermissionName,
				DisplayName: "Self Management",
				Description: "This permission gives access to self management.",
				Resource: &settingsmsg.Resource{
					Type: settingsmsg.Resource_TYPE_USER,
					Id:   "me",
				},
				Value: &settingsmsg.Setting_PermissionValue{
					PermissionValue: &settingsmsg.Permission{
						Operation:  settingsmsg.Permission_OPERATION_READWRITE,
						Constraint: settingsmsg.Permission_CONSTRAINT_OWN,
					},
				},
			},
			{
				Id:          CreateSpacePermissionID,
				Name:        CreateSpacePermissionName,
				DisplayName: "Create own Space",
				Description: "This permission allows to create a space owned by the current user.",
				Resource: &settingsmsg.Resource{
					Type: settingsmsg.Resource_TYPE_SYSTEM, // TODO resource type space? self? me? own?
				},
				Value: &settingsmsg.Setting_PermissionValue{
					PermissionValue: &settingsmsg.Permission{
						Operation:  settingsmsg.Permission_OPERATION_CREATE,
						Constraint: settingsmsg.Permission_CONSTRAINT_OWN,
					},
				},
			},
		},
	}
}

func generateBundleGuestRole() *settingsmsg.Bundle {
	return &settingsmsg.Bundle{
		Id:          BundleUUIDRoleGuest,
		Name:        "guest",
		Type:        settingsmsg.Bundle_TYPE_ROLE,
		Extension:   "ocis-roles",
		DisplayName: "Guest",
		Resource: &settingsmsg.Resource{
			Type: settingsmsg.Resource_TYPE_SYSTEM,
		},
		Settings: []*settingsmsg.Setting{
			{
				Id:          "ca878636-8b1a-4fae-8282-8617a4c13597",
				Name:        "language-readwrite",
				DisplayName: "Permission to read and set the language (self)",
				Resource: &settingsmsg.Resource{
					Type: settingsmsg.Resource_TYPE_SETTING,
					Id:   settingUUIDProfileLanguage,
				},
				Value: &settingsmsg.Setting_PermissionValue{
					PermissionValue: &settingsmsg.Permission{
						Operation:  settingsmsg.Permission_OPERATION_READWRITE,
						Constraint: settingsmsg.Permission_CONSTRAINT_OWN,
					},
				},
			},
		},
	}
}

func generateBundleProfileRequest() *settingsmsg.Bundle {
	return &settingsmsg.Bundle{
		Id:        "2a506de7-99bd-4f0d-994e-c38e72c28fd9",
		Name:      "profile",
		Extension: "ocis-accounts",
		Type:      settingsmsg.Bundle_TYPE_DEFAULT,
		Resource: &settingsmsg.Resource{
			Type: settingsmsg.Resource_TYPE_SYSTEM,
		},
		DisplayName: "Profile",
		Settings: []*settingsmsg.Setting{
			{
				Id:          settingUUIDProfileLanguage,
				Name:        "language",
				DisplayName: "Language",
				Description: "User language",
				Resource: &settingsmsg.Resource{
					Type: settingsmsg.Resource_TYPE_USER,
				},
				Value: &languageSetting,
			},
		},
	}
}

// TODO: languageSetting needed?
var languageSetting = settingsmsg.Setting_SingleChoiceValue{
	SingleChoiceValue: &settingsmsg.SingleChoiceList{
		Options: []*settingsmsg.ListOption{
			{
				Value: &settingsmsg.ListOptionValue{
					Option: &settingsmsg.ListOptionValue_StringValue{
						StringValue: "cs",
					},
				},
				DisplayValue: "Czech",
			},
			{
				Value: &settingsmsg.ListOptionValue{
					Option: &settingsmsg.ListOptionValue_StringValue{
						StringValue: "de",
					},
				},
				DisplayValue: "Deutsch",
			},
			{
				Value: &settingsmsg.ListOptionValue{
					Option: &settingsmsg.ListOptionValue_StringValue{
						StringValue: "en",
					},
				},
				DisplayValue: "English",
				Default:      true,
			},
			{
				Value: &settingsmsg.ListOptionValue{
					Option: &settingsmsg.ListOptionValue_StringValue{
						StringValue: "es",
					},
				},
				DisplayValue: "Español",
			},
			{
				Value: &settingsmsg.ListOptionValue{
					Option: &settingsmsg.ListOptionValue_StringValue{
						StringValue: "fr",
					},
				},
				DisplayValue: "Français",
			},
			{
				Value: &settingsmsg.ListOptionValue{
					Option: &settingsmsg.ListOptionValue_StringValue{
						StringValue: "gl",
					},
				},
				DisplayValue: "Galego",
			},
			{
				Value: &settingsmsg.ListOptionValue{
					Option: &settingsmsg.ListOptionValue_StringValue{
						StringValue: "it",
					},
				},
				DisplayValue: "Italiano",
			},
		},
	},
}

// DefaultRoleAssignments returns (as one might guess) the default role assignments
func DefaultRoleAssignments() []*settingsmsg.UserRoleAssignment {
	return []*settingsmsg.UserRoleAssignment{
		// default admin users
		{
			AccountUuid: "058bff95-6708-4fe5-91e4-9ea3d377588b", // demo user "moss"
			RoleId:      BundleUUIDRoleAdmin,
		}, {
			AccountUuid: "ddc2004c-0977-11eb-9d3f-a793888cd0f8",
			RoleId:      BundleUUIDRoleAdmin,
		},
		// default users with role "user"
		{
			AccountUuid: "4c510ada-c86b-4815-8820-42cdf82c3d51",
			RoleId:      BundleUUIDRoleUser,
		}, {
			AccountUuid: "f7fbf8c8-139b-4376-b307-cf0a8c2d0d9c",
			RoleId:      BundleUUIDRoleUser,
		}, {
			AccountUuid: "932b4540-8d16-481e-8ef4-588e4b6b151c",
			RoleId:      BundleUUIDRoleUser,
		},
		// default users with role "spaceadmin"
		{
			AccountUuid: "534bb038-6f9d-4093-946f-133be61fa4e7",
			RoleId:      BundleUUIDRoleSpaceAdmin,
		},
	}
}
