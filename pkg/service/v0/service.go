package svc

import (
	"context"

	"github.com/owncloud/ocis-settings/pkg/settings"
	store "github.com/owncloud/ocis-settings/pkg/store/filesystem"

	"github.com/owncloud/ocis-settings/pkg/config"
	"github.com/owncloud/ocis-settings/pkg/proto/v0"
)

type Service struct {
	config  *config.Config
	manager settings.Manager
}

// NewService returns a service implementation for Service.
func NewService(cfg *config.Config) Service {
	return Service{
		config:  cfg,
		manager: store.New(cfg),
	}
}

func (g Service) CreateSettingsBundle(c context.Context, req *proto.CreateSettingsBundleRequest, res *proto.CreateSettingsBundleResponse) error {
	r, err := g.manager.Write(req.SettingsBundle)
	if err != nil {
		return err
	}
	res.SettingsBundle = r
	return nil
}

func (g Service) GetSettingsBundle(c context.Context, req *proto.GetSettingsBundleRequest, res *proto.GetSettingsBundleResponse) error {
	r, err := g.manager.Read(req.Extension, req.Key)
	if err != nil {
		return err
	}
	res.SettingsBundle = r
	return nil
}

func (g Service) ListSettingsBundles(c context.Context, req *proto.ListSettingsBundlesRequest, res *proto.ListSettingsBundlesResponse) error {
	r, err := g.manager.ListByExtension(req.Extension)
	if err != nil {
		return err
	}
	res.SettingsBundles = r
	return nil
}
