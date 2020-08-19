package service

import (
	"context"

	mclient "github.com/micro/go-micro/v2/client"
	olog "github.com/owncloud/ocis-pkg/v2/log"
	settings "github.com/owncloud/ocis-settings/pkg/proto/v0"
)

const (
	settingUuidProfileLanguage = "aa8cfbe5-95d4-4f7e-a032-c3c01f5f062f"
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
					Id:          settingUuidProfileLanguage,
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
