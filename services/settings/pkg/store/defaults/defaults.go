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
	// BundleUUIDProfile represents the user profile
	BundleUUIDProfile = "2a506de7-99bd-4f0d-994e-c38e72c28fd9"
	// SettingUUIDProfileLanguage is the hardcoded setting UUID for the user profile language
	SettingUUIDProfileLanguage = "aa8cfbe5-95d4-4f7e-a032-c3c01f5f062f"
	// SettingUUIDProfileDisableNotifications is the hardcoded setting UUID for the disable notifications setting
	SettingUUIDProfileDisableNotifications = "33ffb5d6-cd07-4dc0-afb0-84f7559ae438"
	// SettingUUIDProfileAutoAcceptShares is the hardcoded setting UUID for the disable notifications setting
	SettingUUIDProfileAutoAcceptShares = "ec3ed4a3-3946-4efc-8f9f-76d38b12d3a9"
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
			AccountManagementPermission(All),
			AutoAcceptSharesPermission(Own),
			ChangeLogoPermission(All),
			CreatePublicLinkPermission(All),
			CreateSharePermission(All),
			CreateSpacesPermission(All),
			DeletePersonalSpacesPermission(All),
			DeleteProjectSpacesPermission(All),
			DeleteReadOnlyPublicLinkPasswordPermission(All),
			DisableEmailNotificationsPermission(Own),
			GroupManagementPermission(All),
			LanguageManagementPermission(All),
			ListFavoritesPermission(Own),
			ListSpacesPermission(All),
			ManageSpacePropertiesPermission(All),
			RoleManagementPermission(All),
			SetPersonalSpaceQuotaPermission(All),
			SetProjectSpaceQuotaPermission(All),
			SettingsManagementPermission(All),
			SpaceAbilityPermission(All),
			WriteFavoritesPermission(Own),
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
			AutoAcceptSharesPermission(Own),
			CreatePublicLinkPermission(All),
			CreateSharePermission(All),
			CreateSpacesPermission(All),
			CreateSpacesPermission(Own),
			DeleteProjectSpacesPermission(All),
			DeleteReadOnlyPublicLinkPasswordPermission(All),
			DisableEmailNotificationsPermission(Own),
			LanguageManagementPermission(Own),
			ListFavoritesPermission(Own),
			ListSpacesPermission(All),
			ManageSpacePropertiesPermission(All),
			SelfManagementPermission(Own),
			SetProjectSpaceQuotaPermission(All),
			SpaceAbilityPermission(All),
			WriteFavoritesPermission(Own),
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
			AutoAcceptSharesPermission(Own),
			CreatePublicLinkPermission(All),
			CreateSharePermission(All),
			CreateSpacesPermission(Own),
			DisableEmailNotificationsPermission(Own),
			LanguageManagementPermission(Own),
			ListFavoritesPermission(Own),
			SelfManagementPermission(Own),
			WriteFavoritesPermission(Own),
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
			AutoAcceptSharesPermission(Own),
			DisableEmailNotificationsPermission(Own),
			LanguageManagementPermission(Own),
		},
	}
}

func generateBundleProfileRequest() *settingsmsg.Bundle {
	return &settingsmsg.Bundle{
		Id:        BundleUUIDProfile,
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
			{
				Id:          SettingUUIDProfileAutoAcceptShares,
				Name:        "auto-accept-shares",
				DisplayName: "Auto accept shares",
				Description: "Automatically accept shares",
				Resource: &settingsmsg.Resource{
					Type: settingsmsg.Resource_TYPE_USER,
				},
				Value: &settingsmsg.Setting_BoolValue{BoolValue: &settingsmsg.Bool{Default: true, Label: "auto accept shares"}},
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
						StringValue: "bg",
					},
				},
				DisplayValue: "български",
			},
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
			{
				Value: &settingsmsg.ListOptionValue{
					Option: &settingsmsg.ListOptionValue_StringValue{
						StringValue: "nl",
					},
				},
				DisplayValue: "Nederlands",
			},
			{
				Value: &settingsmsg.ListOptionValue{
					Option: &settingsmsg.ListOptionValue_StringValue{
						StringValue: "ko",
					},
				},
				DisplayValue: "한국어",
			},
			{
				Value: &settingsmsg.ListOptionValue{
					Option: &settingsmsg.ListOptionValue_StringValue{
						StringValue: "sq",
					},
				},
				DisplayValue: "Shqipja",
			},
			{
				Value: &settingsmsg.ListOptionValue{
					Option: &settingsmsg.ListOptionValue_StringValue{
						StringValue: "sv",
					},
				},
				DisplayValue: "Svenska",
			},
			{
				Value: &settingsmsg.ListOptionValue{
					Option: &settingsmsg.ListOptionValue_StringValue{
						StringValue: "tr",
					},
				},
				DisplayValue: "Türkçe",
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
			},
			{
				AccountUuid: "f7fbf8c8-139b-4376-b307-cf0a8c2d0d9c",
				RoleId:      BundleUUIDRoleUser,
			},
			{
				AccountUuid: "932b4540-8d16-481e-8ef4-588e4b6b151c",
				RoleId:      BundleUUIDRoleUser,
			},
			{
				// additional admin user
				AccountUuid: "058bff95-6708-4fe5-91e4-9ea3d377588b", // demo user "moss"
				RoleId:      BundleUUIDRoleAdmin,
			},
			{
				// default users with role "spaceadmin"
				AccountUuid: "534bb038-6f9d-4093-946f-133be61fa4e7",
				RoleId:      BundleUUIDRoleSpaceAdmin,
			},
			{
				// service user
				AccountUuid: "service-user-id",
				RoleId:      BundleUUIDRoleAdmin,
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

	if cfg.ServiceAccountIDAdmin != "" {
		assignments = append(assignments, &settingsmsg.UserRoleAssignment{
			AccountUuid: cfg.ServiceAccountIDAdmin,
			RoleId:      BundleUUIDRoleAdmin,
		})
	}

	return assignments
}
