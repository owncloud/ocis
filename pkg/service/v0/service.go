package service

import (
	"context"

	"github.com/micro/go-micro/v2"
	"github.com/owncloud/ocis-accounts/pkg/account"
	"github.com/owncloud/ocis-accounts/pkg/config"
	"github.com/owncloud/ocis-accounts/pkg/proto/v0"
	olog "github.com/owncloud/ocis-pkg/v2/log"
	settings "github.com/owncloud/ocis-settings/pkg/proto/v0"
)

// New returns a new instance of Service
func New(cfg *config.Config) Service {
	s := Service{
		Config: cfg,
	}

	if newReg, ok := account.Registry[cfg.Manager]; ok {
		s.Manager = newReg(cfg)
	} else {
		l := olog.NewLogger(olog.Name("ocis-accounts"))
		l.Fatal().Msgf("unknown manager: %v", cfg.Manager)
	}

	return s
}

// Service implements the AccountsServiceHandler interface
type Service struct {
	Config  *config.Config
	Manager account.Manager
}

// Set implements the AccountsServiceHandler interface
// This implementation replaces the existent data with the requested. It does not calculate diff
func (s Service) Set(c context.Context, req *proto.Record, res *proto.Record) error {
	r, err := s.Manager.Write(req)
	if err != nil {
		return err
	}

	res.Payload = r.GetPayload()
	return nil
}

// Get implements the AccountsServiceHandler interface
func (s Service) Get(c context.Context, req *proto.GetRequest, res *proto.Record) error {
	// TODO implement other GetRequest properties: Identity, username&password, email
	r, err := s.Manager.Read(req.GetUuid())
	if err != nil {
		return err
	}

	res.Payload = r.GetPayload()
	return nil
}

// Search implements the AccountsServiceHandler interface
func (s Service) Search(ctx context.Context, in *proto.Query, res *proto.Records) error {
	r, err := s.Manager.List()
	if err != nil {
		return err
	}

	// TODO implement filter
	// TODO implement pagination

	res.Records = r
	return nil
}

// RegisterSettingsBundles pushes the settings bundle definitions for this extension to the ocis-settings service.
func RegisterSettingsBundles(l *olog.Logger) {
	// TODO it's ok if this fails. But show a warning that the settings service is not reachable. Make sure that init doesn't die if the settings service is not reachable.
	svc := micro.NewService()
	svc.Init()
	service := settings.NewBundleService("com.owncloud.api.settings", svc.Client()) // TODO fetch service name instead of hardcoding it.

	// TODO avoid hardcoding these values, perhaps load them from a file and using jsonpb's type Marshal.
	requests := []settings.SaveSettingsBundleRequest{
		generateSettingsBundleProfileRequest(),
		generateSettingsBundleNotificationsRequest(),
	}

	for _, request := range requests {
		res, err := service.SaveSettingsBundle(context.Background(), &request)
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
