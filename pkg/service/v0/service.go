package service

import (
	"context"

	"github.com/owncloud/ocis-accounts/pkg/account"
	"github.com/owncloud/ocis-accounts/pkg/config"
	"github.com/owncloud/ocis-accounts/pkg/proto/v0"
	olog "github.com/owncloud/ocis-pkg/v2/log"
	mclient "github.com/micro/go-micro/v2/client"
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
	// TODO this won't work with a registry other than mdns. Look into Micro's client initialization.
	// https://github.com/owncloud/ocis-proxy/issues/38
	service := settings.NewBundleService("com.owncloud.api.settings", mclient.DefaultClient)

	requests := []settings.SaveSettingsBundleRequest{
		generateSettingsBundleProfileRequest(),
		generateSettingsBundleNotificationsRequest(),
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
