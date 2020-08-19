package service

import (
	"context"

	mclient "github.com/micro/go-micro/v2/client"
	olog "github.com/owncloud/ocis-pkg/v2/log"
	settings "github.com/owncloud/ocis-settings/pkg/proto/v0"
	ssvc "github.com/owncloud/ocis-settings/pkg/service/v0"
)

const (
	settingUUIDProfileLanguage = "aa8cfbe5-95d4-4f7e-a032-c3c01f5f062f"
)

// RegisterSettingsBundles pushes the settings bundle definitions for this extension to the ocis-settings service.
func RegisterSettingsBundles(l *olog.Logger) {
	// TODO this won't work with a registry other than mdns. Look into Micro's client initialization.
	// https://github.com/owncloud/ocis-proxy/issues/38
	service := settings.NewBundleService("com.owncloud.api.settings", mclient.DefaultClient)

	bundleRequests := []settings.SaveBundleRequest{
		generateBundleProfileRequest(),
	}

	for i := range bundleRequests {
		res, err := service.SaveBundle(context.Background(), &bundleRequests[i])
		if err != nil {
			l.Err(err).Str("bundle", bundleRequests[i].Bundle.Id).Msg("Error registering bundle")
		} else {
			l.Info().Str("bundle", res.Bundle.Id).Msg("Successfully registered bundle")
		}
	}

	permissionRequests := generateProfilePermissionsRequests()
	for i := range permissionRequests {
		res, err := service.AddSettingToBundle(context.Background(), &permissionRequests[i])
		bundleId := permissionRequests[i].BundleId
		if err != nil {
			l.Err(err).Str("bundle", bundleId).Str("setting", permissionRequests[i].Setting.Id).Msg("Error adding setting to bundle")
		} else {
			l.Info().Str("bundle", bundleId).Str("setting", res.Setting.Id).Msg("Successfully added setting to bundle")
		}
	}
}

var languageSetting = settings.Setting_SingleChoiceValue{
	SingleChoiceValue: &settings.SingleChoiceList{
		Options: []*settings.ListOption{
			{
				Value: &settings.ListOptionValue{
					Option: &settings.ListOptionValue_StringValue{
						StringValue: "cs",
					},
				},
				DisplayValue: "Czech",
			},
			{
				Value: &settings.ListOptionValue{
					Option: &settings.ListOptionValue_StringValue{
						StringValue: "de",
					},
				},
				DisplayValue: "Deutsch",
			},
			{
				Value: &settings.ListOptionValue{
					Option: &settings.ListOptionValue_StringValue{
						StringValue: "en",
					},
				},
				DisplayValue: "English",
			},
			{
				Value: &settings.ListOptionValue{
					Option: &settings.ListOptionValue_StringValue{
						StringValue: "es",
					},
				},
				DisplayValue: "Español",
			},
			{
				Value: &settings.ListOptionValue{
					Option: &settings.ListOptionValue_StringValue{
						StringValue: "fr",
					},
				},
				DisplayValue: "Français",
			},
			{
				Value: &settings.ListOptionValue{
					Option: &settings.ListOptionValue_StringValue{
						StringValue: "gl",
					},
				},
				DisplayValue: "Galego",
			},
			{
				Value: &settings.ListOptionValue{
					Option: &settings.ListOptionValue_StringValue{
						StringValue: "it",
					},
				},
				DisplayValue: "Italiano",
			},
		},
	},
}

func generateBundleProfileRequest() settings.SaveBundleRequest {
	return settings.SaveBundleRequest{
		Bundle: &settings.Bundle{
			Id:        "2a506de7-99bd-4f0d-994e-c38e72c28fd9",
			Name:      "profile",
			Extension: "ocis-accounts",
			Type:      settings.Bundle_TYPE_DEFAULT,
			Resource: &settings.Resource{
				Type: settings.Resource_TYPE_SYSTEM,
			},
			DisplayName: "Profile",
			Settings: []*settings.Setting{
				{
					Id:          settingUUIDProfileLanguage,
					Name:        "language",
					DisplayName: "Language",
					Description: "User language",
					Resource: &settings.Resource{
						Type: settings.Resource_TYPE_USER,
					},
					Value: &languageSetting,
				},
			},
		},
	}
}

func generateProfilePermissionsRequests() []settings.AddSettingToBundleRequest {
	// TODO: we don't want to set up permissions for settings manually in the future. Instead each setting should come with
	// a set of default permissions for the default roles (guest, user, admin).
	return []settings.AddSettingToBundleRequest{
		{
			BundleId: ssvc.BundleUUIDRoleAdmin,
			Setting: &settings.Setting{
				Id:          "7d81f103-0488-4853-bce5-98dcce36d649",
				Name:        "language-create",
				DisplayName: "Permission to set the language",
				Resource: &settings.Resource{
					Type: settings.Resource_TYPE_SETTING,
					Id:   settingUUIDProfileLanguage,
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
			BundleId: ssvc.BundleUUIDRoleAdmin,
			Setting: &settings.Setting{
				Id:          "04ef2fd3-e724-48f6-a411-129dd461c820",
				Name:        "language-read",
				DisplayName: "Permission to read the language",
				Resource: &settings.Resource{
					Type: settings.Resource_TYPE_SETTING,
					Id:   settingUUIDProfileLanguage,
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
			BundleId: ssvc.BundleUUIDRoleAdmin,
			Setting: &settings.Setting{
				Id:          "30ac1e63-10e2-4ef8-bf0a-941cd5b56c5c",
				Name:        "language-update",
				DisplayName: "Permission to update the language",
				Resource: &settings.Resource{
					Type: settings.Resource_TYPE_SETTING,
					Id:   settingUUIDProfileLanguage,
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
			BundleId: ssvc.BundleUUIDRoleUser,
			Setting: &settings.Setting{
				Id:          "640e00d2-4df8-41bd-b1c2-9f30a01e0e99",
				Name:        "language-create",
				DisplayName: "Permission to set the language",
				Resource: &settings.Resource{
					Type: settings.Resource_TYPE_SETTING,
					Id:   settingUUIDProfileLanguage,
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
			BundleId: ssvc.BundleUUIDRoleUser,
			Setting: &settings.Setting{
				Id:          "dcaeb961-da25-46f2-9892-731603a20d3b",
				Name:        "language-read",
				DisplayName: "Permission to read the language",
				Resource: &settings.Resource{
					Type: settings.Resource_TYPE_SETTING,
					Id:   settingUUIDProfileLanguage,
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
			BundleId: ssvc.BundleUUIDRoleUser,
			Setting: &settings.Setting{
				Id:          "e43f364c-ffa5-4080-9621-0d186632a169",
				Name:        "language-update",
				DisplayName: "Permission to update the language",
				Resource: &settings.Resource{
					Type: settings.Resource_TYPE_SETTING,
					Id:   settingUUIDProfileLanguage,
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
			BundleId: ssvc.BundleUUIDRoleGuest,
			Setting: &settings.Setting{
				Id:          "ca878636-8b1a-4fae-8282-8617a4c13597",
				Name:        "language-read",
				DisplayName: "Permission to read the language",
				Resource: &settings.Resource{
					Type: settings.Resource_TYPE_SETTING,
					Id:   settingUUIDProfileLanguage,
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
