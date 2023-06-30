package defaults

import (
	settingsmsg "github.com/owncloud/ocis/v2/protogen/gen/ocis/messages/settings/v0"
	"github.com/owncloud/ocis/v2/services/settings/pkg/config"
)

const (
	// BundleUUIDRoleAdmin represents the admin role
	BundleUUIDRoleAdmin = "71881883-1768-46bd-a24d-a356a2afdf7f"

	// BundleUUIDRoleSpaceAdmin represents the space admin role
	BundleUUIDRoleSpaceAdmin = "2aadd357-682c-406b-8874-293091995fdd"

	// BundleUUIDRoleUser represents the user role.
	BundleUUIDRoleUser = "d7beeea8-8ff4-406b-8fb6-ab2dd81e6b11"

	// BundleUUIDRoleUserLight represents the user light role.
	BundleUUIDRoleUserLight = "38071a68-456a-4553-846a-fa67bf5596cc"

	// RoleManagementPermissionID is the hardcoded setting UUID for the role management permission
	RoleManagementPermissionID string = "a53e601e-571f-4f86-8fec-d4576ef49c62"
	// RoleManagementPermissionName is the hardcoded setting name for the role management permission
	RoleManagementPermissionName string = "Roles.ReadWrite"

	// SettingsManagementPermissionID is the hardcoded setting UUID for the settings management permission
	SettingsManagementPermissionID string = "3d58f441-4a05-42f8-9411-ef5874528ae1"
	// SettingsManagementPermissionName is the hardcoded setting name for the settings management permission
	SettingsManagementPermissionName string = "Settings.ReadWrite"

	// LanguageReadWriteID is the hardcoded setting UUID for the language read write all permission
	LanguageReadWriteID string = "7d81f103-0488-4853-bce5-98dcce36d649"
	// LanguageReadWriteName is the hardcoded setting name for the language read write all permission
	LanguageReadWriteName string = "Language.ReadWrite"

	// DisableEmailNotificationsPermissionID is the hardcoded setting UUID for the disable email notifications permission
	DisableEmailNotificationsPermissionID string = "ad5bb5e5-dc13-4cd3-9304-09a424564ea8"
	// DisableEmailNotificationsPermissionName is the hardcoded setting name for the disable email notifications permission
	DisableEmailNotificationsPermissionName string = "EmailNotifications.ReadWriteDisabled"
	// DisableEmailNotificationsPermissionDisplayName is the hardcoded setting name for the disable email notifications permission
	DisableEmailNotificationsPermissionDisplayName string = "Disable Email Notifications"

	// SetPersonalSpaceQuotaPermissionID is the hardcoded setting UUID for the set personal space quota permission
	SetPersonalSpaceQuotaPermissionID string = "4e6f9709-f9e7-44f1-95d4-b762d27b7896"
	// SetPersonalSpaceQuotaPermissionName is the hardcoded setting name for the set personal space quota permission
	SetPersonalSpaceQuotaPermissionName string = "Drives.ReadWritePersonalQuota"

	// SetProjectSpaceQuotaPermissionID is the hardcoded setting UUID for the set project space quota permission
	SetProjectSpaceQuotaPermissionID string = "977f0ae6-0da2-4856-93f3-22e0a8482489"
	// SetProjectSpaceQuotaPermissionName is the hardcoded setting name for the set project space quota permission
	SetProjectSpaceQuotaPermissionName string = "Drives.ReadWriteProjectQuota"

	// ListAllSpacesPermissionID is the hardcoded setting UUID for the list all spaces permission
	ListAllSpacesPermissionID string = "016f6ddd-9501-4a0a-8ebe-64a20ee8ec82"
	// ListAllSpacesPermissionName is the hardcoded setting name for the list all spaces permission
	ListAllSpacesPermissionName string = "Drives.List"

	// CreateSpacePermissionID is the hardcoded setting UUID for the create space permission
	CreateSpacePermissionID string = "79e13b30-3e22-11eb-bc51-0b9f0bad9a58"
	// CreateSpacePermissionName is the hardcoded setting name for the create space permission
	CreateSpacePermissionName string = "Drives.Create"

	// DeleteHomeSpacesPermissionID is the hardcoded setting UUID for the delete home space permission
	DeleteHomeSpacesPermissionID string = "5de9fe0a-4bc5-4a47-b758-28f370caf169"
	// DeleteHomeSpacesPermissionName is the hardcoded setting name for the delete home space permission
	DeleteHomeSpacesPermissionName string = "Drives.DeletePersonal"

	// DeleteAllSpacesPermissionID is the hardcoded setting UUID for the delete all spaces permission
	DeleteAllSpacesPermissionID string = "fb60b004-c1fa-4f09-bf87-55ce7d46ac61"
	// DeleteAllSpacesPermissionName is the hardcoded setting name for the delete all space permission
	DeleteAllSpacesPermissionName string = "Drives.DeleteProject"

	// ManageSpacePropertiesPermissionID is the hardcoded setting UUID for the manage space properties permission
	ManageSpacePropertiesPermissionID string = "b44b4054-31a2-42b8-bb71-968b15cfbd4f"
	// ManageSpacePropertiesPermissionName is the hardcoded setting name for the manage space properties permission
	ManageSpacePropertiesPermissionName string = "Drives.ReadWrite"

	// SpaceAbilityPermissionID is the hardcoded setting UUID for the space ability permission
	SpaceAbilityPermissionID string = "cf3faa8c-50d9-4f84-9650-ff9faf21aa9d"
	// SpaceAbilityPermissionName is the hardcoded setting name for the space ability permission
	SpaceAbilityPermissionName string = "Drives.ReadWriteEnabled"

	// SettingUUIDProfileLanguage is the hardcoded setting UUID for the user profile language
	SettingUUIDProfileLanguage = "aa8cfbe5-95d4-4f7e-a032-c3c01f5f062f"
	// SettingUUIDProfileDisableNotifications is the hardcoded setting UUID for the disable notifications setting
	SettingUUIDProfileDisableNotifications = "33ffb5d6-cd07-4dc0-afb0-84f7559ae438"

	// AccountManagementPermissionID is the hardcoded setting UUID for the account management permission
	AccountManagementPermissionID string = "8e587774-d929-4215-910b-a317b1e80f73"
	// AccountManagementPermissionName is the hardcoded setting name for the account management permission
	AccountManagementPermissionName string = "Accounts.ReadWrite"
	// GroupManagementPermissionID is the hardcoded setting UUID for the group management permission
	GroupManagementPermissionID string = "522adfbe-5908-45b4-b135-41979de73245"
	// GroupManagementPermissionName is the hardcoded setting name for the group management permission
	GroupManagementPermissionName string = "Groups.ReadWrite"
	// SelfManagementPermissionID is the hardcoded setting UUID for the self management permission
	SelfManagementPermissionID string = "e03070e9-4362-4cc6-a872-1c7cb2eb2b8e"
	// SelfManagementPermissionName is the hardcoded setting name for the self management permission
	SelfManagementPermissionName string = "Self.ReadWrite"

	// ChangeLogoPermissionID is the hardcoded setting UUID for the change-logo permission
	ChangeLogoPermissionID string = "ed83fc10-1f54-4a9e-b5a7-fb517f5f3e01"
	// ChangeLogoPermissionName is the hardcoded setting name for the change-logo permission
	ChangeLogoPermissionName string = "Logo.Write"

	// WritePublicLinkPermissionID is the hardcoded setting UUID for the PublicLink.Write permission
	WritePublicLinkPermissionID string = "11516bbd-7157-49e1-b6ac-d00c820f980b"
	// WritePublicLinkPermissionName is the hardcoded setting name for the PublicLink.Write permission
	WritePublicLinkPermissionName string = "PublicLink.Write"
)

// GenerateBundlesDefaultRoles bootstraps the default roles.
func GenerateBundlesDefaultRoles() []*settingsmsg.Bundle {
	return []*settingsmsg.Bundle{
		generateBundleAdminRole(),
		generateBundleUserRole(),
		generateBundleUserLightRole(),
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
				Id:          LanguageReadWriteID,
				Name:        LanguageReadWriteName,
				DisplayName: "Permission to read and set the language (anyone)",
				Resource: &settingsmsg.Resource{
					Type: settingsmsg.Resource_TYPE_SETTING,
					Id:   SettingUUIDProfileLanguage,
				},
				Value: &settingsmsg.Setting_PermissionValue{
					PermissionValue: &settingsmsg.Permission{
						Operation:  settingsmsg.Permission_OPERATION_READWRITE,
						Constraint: settingsmsg.Permission_CONSTRAINT_ALL,
					},
				},
			},
			{
				Id:          DisableEmailNotificationsPermissionID,
				Name:        DisableEmailNotificationsPermissionName,
				DisplayName: DisableEmailNotificationsPermissionDisplayName,
				Resource: &settingsmsg.Resource{
					Type: settingsmsg.Resource_TYPE_SETTING,
					Id:   SettingUUIDProfileDisableNotifications,
				},
				Value: &settingsmsg.Setting_PermissionValue{
					PermissionValue: &settingsmsg.Permission{
						Operation:  settingsmsg.Permission_OPERATION_READWRITE,
						Constraint: settingsmsg.Permission_CONSTRAINT_OWN,
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
				Id:          SetPersonalSpaceQuotaPermissionID,
				Name:        SetPersonalSpaceQuotaPermissionName,
				DisplayName: "Set Personal Space Quota",
				Description: "This permission allows managing personal space quotas.",
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
				Id:          SetProjectSpaceQuotaPermissionID,
				Name:        SetProjectSpaceQuotaPermissionName,
				DisplayName: "Set Project Space Quota",
				Description: "This permission allows managing project space quotas.",
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
				Description: "This permission allows creating new spaces.",
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
				Description: "This permission allows listing all spaces.",
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
				Id:          DeleteHomeSpacesPermissionID,
				Name:        DeleteHomeSpacesPermissionName,
				DisplayName: "Delete All Home Spaces",
				Description: "This permission allows deleting home spaces.",
				Resource: &settingsmsg.Resource{
					Type: settingsmsg.Resource_TYPE_SYSTEM,
				},
				Value: &settingsmsg.Setting_PermissionValue{
					PermissionValue: &settingsmsg.Permission{
						Operation:  settingsmsg.Permission_OPERATION_DELETE,
						Constraint: settingsmsg.Permission_CONSTRAINT_ALL,
					},
				},
			},
			{
				Id:          DeleteAllSpacesPermissionID,
				Name:        DeleteAllSpacesPermissionName,
				DisplayName: "Delete AllSpaces",
				Description: "This permission allows deleting all spaces.",
				Resource: &settingsmsg.Resource{
					Type: settingsmsg.Resource_TYPE_SYSTEM,
				},
				Value: &settingsmsg.Setting_PermissionValue{
					PermissionValue: &settingsmsg.Permission{
						Operation:  settingsmsg.Permission_OPERATION_DELETE,
						Constraint: settingsmsg.Permission_CONSTRAINT_ALL,
					},
				},
			},
			{
				Id:          ChangeLogoPermissionID,
				Name:        ChangeLogoPermissionName,
				DisplayName: "Change logo",
				Description: "This permission permits to change the system logo.",
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
				Id:          WritePublicLinkPermissionID,
				Name:        WritePublicLinkPermissionName,
				DisplayName: "Write publiclink",
				Description: "This permission allows creating public links.",
				Resource: &settingsmsg.Resource{
					Type: settingsmsg.Resource_TYPE_SHARE,
				},
				Value: &settingsmsg.Setting_PermissionValue{
					PermissionValue: &settingsmsg.Permission{
						Operation:  settingsmsg.Permission_OPERATION_WRITE,
						Constraint: settingsmsg.Permission_CONSTRAINT_ALL,
					},
				},
			},
			{
				Id:          ManageSpacePropertiesPermissionID,
				Name:        ManageSpacePropertiesPermissionName,
				DisplayName: "Manage space properties",
				Description: "This permission allows managing space properties such as name and description.",
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
				Id:          SpaceAbilityPermissionID,
				Name:        SpaceAbilityPermissionName,
				DisplayName: "Space ability",
				Description: "This permission allows enabling and disabling spaces.",
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
				Id:          ManageSpacePropertiesPermissionID,
				Name:        ManageSpacePropertiesPermissionName,
				DisplayName: "Manage space properties",
				Description: "This permission allows managing space properties such as name and description.",
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
				Id:          SpaceAbilityPermissionID,
				Name:        SpaceAbilityPermissionName,
				DisplayName: "Space ability",
				Description: "This permission allows enabling and disabling spaces.",
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
				Id:          DeleteAllSpacesPermissionID,
				Name:        DeleteAllSpacesPermissionName,
				DisplayName: "Delete AllSpaces",
				Description: "This permission allows to delete all spaces.",
				Resource: &settingsmsg.Resource{
					Type: settingsmsg.Resource_TYPE_SYSTEM,
				},
				Value: &settingsmsg.Setting_PermissionValue{
					PermissionValue: &settingsmsg.Permission{
						Operation:  settingsmsg.Permission_OPERATION_DELETE,
						Constraint: settingsmsg.Permission_CONSTRAINT_ALL,
					},
				},
			},
			{
				Id:          SetProjectSpaceQuotaPermissionID,
				Name:        SetProjectSpaceQuotaPermissionName,
				DisplayName: "Set Project Space Quota",
				Description: "This permission allows managing project space quotas.",
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
				Description: "This permission allows creating new spaces.",
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
				Id:          LanguageReadWriteID,
				Name:        LanguageReadWriteName,
				DisplayName: "Permission to read and set the language (self)",
				Resource: &settingsmsg.Resource{
					Type: settingsmsg.Resource_TYPE_SETTING,
					Id:   SettingUUIDProfileLanguage,
				},
				Value: &settingsmsg.Setting_PermissionValue{
					PermissionValue: &settingsmsg.Permission{
						Operation:  settingsmsg.Permission_OPERATION_READWRITE,
						Constraint: settingsmsg.Permission_CONSTRAINT_OWN,
					},
				},
			},
			{
				Id:          DisableEmailNotificationsPermissionID,
				Name:        DisableEmailNotificationsPermissionName,
				DisplayName: DisableEmailNotificationsPermissionDisplayName,
				Resource: &settingsmsg.Resource{
					Type: settingsmsg.Resource_TYPE_SETTING,
					Id:   SettingUUIDProfileDisableNotifications,
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
				Description: "This permission allows creating a space owned by the current user.",
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
			{
				Id:          WritePublicLinkPermissionID,
				Name:        WritePublicLinkPermissionName,
				DisplayName: "Write publiclink",
				Description: "This permission permits to write a public link.",
				Resource: &settingsmsg.Resource{
					Type: settingsmsg.Resource_TYPE_SHARE,
				},
				Value: &settingsmsg.Setting_PermissionValue{
					PermissionValue: &settingsmsg.Permission{
						Operation:  settingsmsg.Permission_OPERATION_WRITE,
						Constraint: settingsmsg.Permission_CONSTRAINT_ALL,
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
				Id:          LanguageReadWriteID,
				Name:        LanguageReadWriteName,
				DisplayName: "Permission to read and set the language (self)",
				Resource: &settingsmsg.Resource{
					Type: settingsmsg.Resource_TYPE_SETTING,
					Id:   SettingUUIDProfileLanguage,
				},
				Value: &settingsmsg.Setting_PermissionValue{
					PermissionValue: &settingsmsg.Permission{
						Operation:  settingsmsg.Permission_OPERATION_READWRITE,
						Constraint: settingsmsg.Permission_CONSTRAINT_OWN,
					},
				},
			},
			{
				Id:          DisableEmailNotificationsPermissionID,
				Name:        DisableEmailNotificationsPermissionName,
				DisplayName: DisableEmailNotificationsPermissionDisplayName,
				Resource: &settingsmsg.Resource{
					Type: settingsmsg.Resource_TYPE_SETTING,
					Id:   SettingUUIDProfileDisableNotifications,
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
				Description: "This permission allows creating a space owned by the current user.",
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
			{
				Id:          WritePublicLinkPermissionID,
				Name:        WritePublicLinkPermissionName,
				DisplayName: "Write publiclink",
				Description: "This permission permits to write a public link.",
				Resource: &settingsmsg.Resource{
					Type: settingsmsg.Resource_TYPE_SHARE,
				},
				Value: &settingsmsg.Setting_PermissionValue{
					PermissionValue: &settingsmsg.Permission{
						Operation:  settingsmsg.Permission_OPERATION_WRITE,
						Constraint: settingsmsg.Permission_CONSTRAINT_ALL,
					},
				},
			},
		},
	}
}

func generateBundleUserLightRole() *settingsmsg.Bundle {
	return &settingsmsg.Bundle{
		Id:          BundleUUIDRoleUserLight,
		Name:        "user-light",
		Type:        settingsmsg.Bundle_TYPE_ROLE,
		Extension:   "ocis-roles",
		DisplayName: "User Light",
		Resource: &settingsmsg.Resource{
			Type: settingsmsg.Resource_TYPE_SYSTEM,
		},
		Settings: []*settingsmsg.Setting{
			{
				Id:          LanguageReadWriteID,
				Name:        LanguageReadWriteName,
				DisplayName: "Permission to read and set the language (self)",
				Resource: &settingsmsg.Resource{
					Type: settingsmsg.Resource_TYPE_SETTING,
					Id:   SettingUUIDProfileLanguage,
				},
				Value: &settingsmsg.Setting_PermissionValue{
					PermissionValue: &settingsmsg.Permission{
						Operation:  settingsmsg.Permission_OPERATION_READWRITE,
						Constraint: settingsmsg.Permission_CONSTRAINT_OWN,
					},
				},
			},
			{
				Id:          DisableEmailNotificationsPermissionID,
				Name:        DisableEmailNotificationsPermissionName,
				DisplayName: DisableEmailNotificationsPermissionDisplayName,
				Resource: &settingsmsg.Resource{
					Type: settingsmsg.Resource_TYPE_SETTING,
					Id:   SettingUUIDProfileDisableNotifications,
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
				Id:          SettingUUIDProfileLanguage,
				Name:        "language",
				DisplayName: "Language",
				Description: "User language",
				Resource: &settingsmsg.Resource{
					Type: settingsmsg.Resource_TYPE_USER,
				},
				Value: &languageSetting,
			},
			{
				Id:          SettingUUIDProfileDisableNotifications,
				Name:        "disable-email-notifications",
				DisplayName: "Disable Email Notifications",
				Description: "Disable email notifications",
				Resource: &settingsmsg.Resource{
					Type: settingsmsg.Resource_TYPE_USER,
				},
				Value: &settingsmsg.Setting_BoolValue{BoolValue: &settingsmsg.Bool{Default: false, Label: "disable notifications"}},
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
func DefaultRoleAssignments(cfg *config.Config) []*settingsmsg.UserRoleAssignment {
	assignments := []*settingsmsg.UserRoleAssignment{}

	if cfg.SetupDefaultAssignments {
		assignments = []*settingsmsg.UserRoleAssignment{
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
			{
				// additional admin user
				AccountUuid: "058bff95-6708-4fe5-91e4-9ea3d377588b", // demo user "moss"
				RoleId:      BundleUUIDRoleAdmin,
			}, {
				// default users with role "spaceadmin"
				AccountUuid: "534bb038-6f9d-4093-946f-133be61fa4e7",
				RoleId:      BundleUUIDRoleSpaceAdmin,
			},
		}
	}

	if cfg.AdminUserID != "" {
		// default admin user
		assignments = append(assignments, &settingsmsg.UserRoleAssignment{
			AccountUuid: cfg.AdminUserID,
			RoleId:      BundleUUIDRoleAdmin,
		})
	}

	return assignments
}
