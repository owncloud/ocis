package service

import (
	"context"

	mclient "github.com/micro/go-micro/v2/client"
	olog "github.com/owncloud/ocis-pkg/v2/log"
	settings "github.com/owncloud/ocis-settings/pkg/proto/v0"
)

func generateSettingsBundleProfileRequest() settings.SaveSettingsBundleRequest {
	return settings.SaveSettingsBundleRequest{
		SettingsBundle: &settings.SettingsBundle{
			Identifier: &settings.Identifier{
				Extension: "ocis-accounts",
				BundleKey: "profile",
			},
			DisplayName: "Profile",
			Settings: []*settings.Setting{
				{
					SettingKey:  "language",
					DisplayName: "Language",
					Description: "User language",
					Value: &settings.Setting_SingleChoiceValue{
						SingleChoiceValue: &settings.SingleChoiceListSetting{
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
					},
				},
			},
		},
	}
}

// RegisterSettingsBundles pushes the settings bundle definitions for this extension to the ocis-settings service.
func RegisterSettingsBundles(l *olog.Logger) {
	// TODO this won't work with a registry other than mdns. Look into Micro's client initialization.
	// https://github.com/owncloud/ocis-proxy/issues/38
	service := settings.NewBundleService("com.owncloud.api.settings", mclient.DefaultClient)

	requests := []settings.SaveSettingsBundleRequest{
		generateSettingsBundleProfileRequest(),
	}

	for i := range requests {
		res, err := service.SaveSettingsBundle(context.Background(), &requests[i])
		if err != nil {
			l.Err(err).
				Msg("Error registering settings bundle")
		} else {
			l.Info().
				Str("bundle key", res.SettingsBundle.Identifier.BundleKey).
				Msg("Successfully registered settings bundle")
		}
	}
}
