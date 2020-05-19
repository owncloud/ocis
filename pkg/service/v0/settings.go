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
					SettingKey:  "firstname",
					DisplayName: "Firstname",
					Description: "Input for firstname",
					Value: &settings.Setting_StringValue{
						StringValue: &settings.StringSetting{
							Placeholder: "Set firstname",
						},
					},
				},
				{
					SettingKey:  "lastname",
					DisplayName: "Lastname",
					Description: "Input for lastname",
					Value: &settings.Setting_StringValue{
						StringValue: &settings.StringSetting{
							Placeholder: "Set lastname",
						},
					},
				},
				{
					SettingKey:  "age",
					DisplayName: "Age",
					Description: "Input for age",
					Value: &settings.Setting_IntValue{
						IntValue: &settings.IntSetting{
							Placeholder: "Set age",
							Min:         16,
							Max:         200,
							Step:        2,
						},
					},
				},
				{
					SettingKey:  "timezone",
					DisplayName: "Timezone",
					Description: "User timezone",
					Value: &settings.Setting_SingleChoiceValue{
						SingleChoiceValue: &settings.SingleChoiceListSetting{
							Options: []*settings.ListOption{
								{
									Value: &settings.ListOptionValue{
										Option: &settings.ListOptionValue_StringValue{
											StringValue: "Europe/Berlin",
										},
									},
									DisplayValue: "Europe/Berlin",
								},
								{
									Value: &settings.ListOptionValue{
										Option: &settings.ListOptionValue_StringValue{
											StringValue: "Asia/Kathmandu",
										},
									},
									DisplayValue: "Asia/Kathmandu",
								},
							},
						},
					},
				},
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

func generateSettingsBundleNotificationsRequest() settings.SaveSettingsBundleRequest {
	return settings.SaveSettingsBundleRequest{
		SettingsBundle: &settings.SettingsBundle{
			Identifier: &settings.Identifier{
				Extension: "ocis-accounts",
				BundleKey: "notifications",
			},
			DisplayName: "Notifications",
			Settings: []*settings.Setting{
				{
					SettingKey:  "email",
					DisplayName: "Email",
					Value: &settings.Setting_BoolValue{
						BoolValue: &settings.BoolSetting{
							Default: false,
							Label:   "Send via email",
						},
					},
				},
				{
					SettingKey:  "stream",
					DisplayName: "Stream",
					Value: &settings.Setting_BoolValue{
						BoolValue: &settings.BoolSetting{
							Default: true,
							Label:   "Show in stream",
						},
					},
				},
				{
					SettingKey:  "transport",
					DisplayName: "Transport",
					Value: &settings.Setting_MultiChoiceValue{
						MultiChoiceValue: &settings.MultiChoiceListSetting{
							Options: []*settings.ListOption{
								{
									Value: &settings.ListOptionValue{
										Option: &settings.ListOptionValue_StringValue{
											StringValue: "email",
										},
									},
									DisplayValue: "Send via email",
								},
								{
									Value: &settings.ListOptionValue{
										Option: &settings.ListOptionValue_StringValue{
											StringValue: "stream",
										},
									},
									DisplayValue: "Show in stream",
									Default:      true,
								},
							},
						},
					},
				},
			},
		},
	}
}

