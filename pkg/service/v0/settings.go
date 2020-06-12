package service

import settings "github.com/owncloud/ocis-settings/pkg/proto/v0"

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
