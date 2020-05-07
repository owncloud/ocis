package service

import (
	"context"
	"errors"

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
func (s Service) Get(c context.Context, req *proto.GetRequest, res *proto.Record) (err error) {
	// TODO implement other GetRequest properties: Identity, username&password, email
	var r, ruuid, rname, rmail *proto.Record
	if req.GetIdentity() != nil {
		r, err = s.Manager.ReadByIdentity(req.GetIdentity())
		if err != nil {
			return err
		}
	}
	if req.GetUuid() != "" {
		ruuid, err = s.Manager.Read(req.GetUuid())
		if err != nil {
			return err
		}

		if r == nil {
			r = ruuid
		} else if r.Key != ruuid.Key {
			r = nil
			return errors.New("uuid mismatch")
		}
	}
	if req.GetUsername() != "" {
		rname, err = s.Manager.ReadByUsername(req.GetUsername())
		if err != nil {
			return err
		}
		if r == nil {
			r = rname
		} else if r.Key != rname.Key {
			r = nil
			return errors.New("username mismatch")
		}
	}
	if req.GetEmail() != "" {
		rmail, err = s.Manager.ReadByEmail(req.GetEmail())
		if err != nil {
			return err
		}
		if r == nil {
			r = rmail
		} else if r.Key != rmail.Key {
			r = nil
			return errors.New("email mismatch")
		}
	}

	if r != nil {
		// TODO store only salted hash
		if req.GetPassword() != "" {
			if r.Payload.Account.Password != req.GetPassword() {
				return errors.New("wrong password")
			}
		}
		res.Key = r.Key
		res.Payload = r.GetPayload()
		// password never leaves
		res.Payload.Account.Password = ""
		return nil
	}

	return errors.New("at least one request param must be set")
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
