package service

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
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

// Service implements the SettingsServiceHandler interface
type Service struct {
	Config  *config.Config
	Manager account.Manager
}

// Set implements the SettingsServiceHandler interface
// This implementation replaces the existent data with the requested. It does not calculate diff
func (s Service) Set(c context.Context, req *proto.Record, res *proto.Record) error {
	r, err := s.Manager.Write(req)
	if err != nil {
		return err
	}

	res.Payload = r.GetPayload()
	return nil
}

// Get implements the SettingsServiceHandler interface
func (s Service) Get(c context.Context, req *proto.Query, res *proto.Record) error {
	r, err := s.Manager.Read(req.GetKey())
	if err != nil {
		return err
	}

	res.Payload = r.GetPayload()
	return nil
}

// List implements the SettingsServiceHandler interface
func (s Service) List(ctx context.Context, in *empty.Empty, res *proto.Records) error {
	r, err := s.Manager.List()
	if err != nil {
		return err
	}

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
	requests := []settings.CreateSettingsBundleRequest{
		generateSettingsBundleProfileRequest(),
		generateSettingsBundleNotificationsRequest(),
	}

	for _, request := range requests {
		res, err := service.CreateSettingsBundle(context.Background(), &request)
		if err != nil {
			l.Err(err).
				Msg("Error registering settings bundle")
		} else {
			l.Info().
				Str("bundle key", res.GetSettingsBundle().GetKey()).
				Msg("Successfully registered settings bundle")
		}
	}
}

func generateSettingsBundleProfileRequest() settings.CreateSettingsBundleRequest {
	return settings.CreateSettingsBundleRequest{
		SettingsBundle: &settings.SettingsBundle{
			Extension:   "ocis-accounts",
			Key:         "profile",
			DisplayName: "Profile",
			Settings: []*settings.Setting{
				{
					Key:         "firstname",
					DisplayName: "Firstname",
					Description: "Input for firstname",
					Value: &settings.Setting_StringValue{
						StringValue: &settings.StringSetting{
							Placeholder: "Set firstname",
						},
					},
				},
				{
					Key:         "lastname",
					DisplayName: "Lastname",
					Description: "Input for lastname",
					Value: &settings.Setting_StringValue{
						StringValue: &settings.StringSetting{
							Placeholder: "Set lastname",
						},
					},
				},
				{
					Key:         "age",
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
					Key:         "timezone",
					DisplayName: "Timezone",
					Description: "User timezone",
					Value: &settings.Setting_SingleChoiceValue{
						SingleChoiceValue: &settings.SingleChoiceListSetting{
							Options: []*settings.ListOption{
								{
									Option: &settings.ListOption_StringValue{
										StringValue: "Europe/Berlin",
									},
									DisplayValue: "Europe/Berlin",
								},
								{
									Option: &settings.ListOption_StringValue{
										StringValue: "Asia/Kathmandu",
									},
									DisplayValue: "Asia/Kathmandu",
								},
							},
						},
					},
				},
				{
					Key:         "language",
					DisplayName: "Language",
					Description: "User language",
					Value: &settings.Setting_SingleChoiceValue{
						SingleChoiceValue: &settings.SingleChoiceListSetting{
							Options: []*settings.ListOption{
								{
									Option: &settings.ListOption_StringValue{
										StringValue: "de_DE",
									},
									DisplayValue: "Deutsch",
								},
								{
									Option: &settings.ListOption_StringValue{
										StringValue: "en_EN",
									},
									DisplayValue: "English",
								},
							},
						},
					},
				},
			},
		},
	}
}

func generateSettingsBundleNotificationsRequest() settings.CreateSettingsBundleRequest {
	return settings.CreateSettingsBundleRequest{
		SettingsBundle: &settings.SettingsBundle{
			Extension:   "ocis-accounts",
			Key:         "notifications",
			DisplayName: "Notifications",
			Settings: []*settings.Setting{
				{
					Key:         "email",
					DisplayName: "Email",
					Value: &settings.Setting_BoolValue{
						BoolValue: &settings.BoolSetting{
							Default: false,
							Label:   "Send via email",
						},
					},
				},
				{
					Key:         "stream",
					DisplayName: "Stream",
					Value: &settings.Setting_BoolValue{
						BoolValue: &settings.BoolSetting{
							Default: true,
							Label:   "Show in stream",
						},
					},
				},
			},
		},
	}
}
