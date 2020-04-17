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

// ReportSettingsBundle writes the settings bundle for the extension on the ocis-settings service.
// TODO implement retry logic.
// TODO ensure the settings service is available.
func ReportSettingsBundle(l *olog.Logger) {
	// TODO wrap this in an infinite for loop with a BREAK label to GOTO once the request succeed.
	svc := micro.NewService()
	svc.Init()
	service := settings.NewBundleService("com.owncloud.ocis-settings", svc.Client()) // TODO fetch service name instead of hardcoding it.

	// TODO avoid hardcoding these values, perhaps load them from a file and using jsonpb's type Marshal.
	createBundleRequest := settings.CreateSettingsBundleRequest{
		SettingsBundle: &settings.SettingsBundle{
			Extension:   "ocis-accounts",
			Key:         "profile",
			DisplayName: "Profile",
			Settings: []*settings.Setting{
				&settings.Setting{
					Key:         "timezone",
					DisplayName: "Timezone",
					Description: "User timezone",
					Value: &settings.Setting_SingleChoiceValue{
						SingleChoiceValue: &settings.SingleChoiceListSetting{
							Options: []*settings.ListOption{
								&settings.ListOption{
									Selected: true,
									Option: &settings.ListOption_StringValue{
										StringValue: "Europe/Berlin",
									},
								},
								&settings.ListOption{
									Selected: false,
									Option: &settings.ListOption_StringValue{
										StringValue: "Asia/Kathmandu",
									},
								},
							},
						},
					},
				},
				&settings.Setting{
					Key:         "language",
					DisplayName: "Language",
					Description: "User language",
					Value: &settings.Setting_SingleChoiceValue{
						SingleChoiceValue: &settings.SingleChoiceListSetting{
							Options: []*settings.ListOption{
								&settings.ListOption{
									Selected: true,
									Option: &settings.ListOption_StringValue{
										StringValue: "de_DE",
									},
								},
								&settings.ListOption{
									Selected: false,
									Option: &settings.ListOption_StringValue{
										StringValue: "en_EN",
									},
								},
							},
						},
					},
				},
			},
		},
	}

	res, err := service.CreateSettingsBundle(context.Background(), &createBundleRequest)
	if err != nil {
		l.Err(err).
			Msg("Error reporting settings bundle")
	} else {
		l.Info().
			Str("bundle key", res.GetSettingsBundle().GetKey()).
			Msg("Succesfully reported settings bundle")
	}

}
