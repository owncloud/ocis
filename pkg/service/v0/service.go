package svc

import (
	"context"

	"github.com/owncloud/ocis-settings/pkg/settings"
	store "github.com/owncloud/ocis-settings/pkg/store/filesystem"

	"github.com/owncloud/ocis-settings/pkg/config"
	"github.com/owncloud/ocis-settings/pkg/proto/v0"
)

type Service struct {
	config       *config.Config
	manager      settings.Manager
}

// NewService returns a service implementation for Service.
func NewService(cfg *config.Config) Service {
	return Service{
		config:  cfg,
		manager: store.New(cfg),
	}
}

func (g Service) CreateSettingsBundle(c context.Context, req *proto.CreateSettingsBundleRequest, res *proto.CreateSettingsBundleResponse) error {
	r, err := g.manager.WriteBundle(req.SettingsBundle)
	if err != nil {
		return err
	}
	res.SettingsBundle = r
	return nil
}

func (g Service) GetSettingsBundle(c context.Context, req *proto.GetSettingsBundleRequest, res *proto.GetSettingsBundleResponse) error {
	r, err := g.manager.ReadBundle(req.Extension, req.BundleKey)
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

func (g Service) SaveSettingsValue(c context.Context, req *proto.SaveSettingsValueRequest, res *proto.SaveSettingsValueResponse) error {
	r, err := g.manager.WriteValue(req.SettingsValue)
	if err != nil {
		return err
	}
	res.SettingsValue = r
	return nil
}

func (g Service) GetSettingsValue(c context.Context, req *proto.GetSettingsValueRequest, res *proto.GetSettingsValueResponse) error {
	r, err := g.manager.ReadValue(req.AccountUuid, req.Extension, req.BundleKey, req.SettingKey)
	if err != nil {
		return err
	}
	res.SettingsValue = r
	return nil
}
