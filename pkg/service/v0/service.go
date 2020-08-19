package svc

import (
	"context"

	"github.com/owncloud/ocis-pkg/v2/log"
	"github.com/owncloud/ocis-pkg/v2/middleware"
	"github.com/owncloud/ocis-settings/pkg/config"
	"github.com/owncloud/ocis-settings/pkg/proto/v0"
	"github.com/owncloud/ocis-settings/pkg/settings"
	store "github.com/owncloud/ocis-settings/pkg/store/filesystem"
)

// Service represents a service.
type Service struct {
	config  *config.Config
	logger  log.Logger
	manager settings.Manager
}

// NewService returns a service implementation for Service.
func NewService(cfg *config.Config, logger log.Logger) Service {
	return Service{
		config:  cfg,
		logger:  logger,
		manager: store.New(cfg),
	}
}

// SaveBundle implements the BundleServiceHandler interface
func (g Service) SaveBundle(c context.Context, req *proto.SaveBundleRequest, res *proto.SaveBundleResponse) error {
	cleanUpResource(c, req.Bundle.Resource)
	if validationError := validateSaveBundle(req); validationError != nil {
		return validationError
	}
	r, err := g.manager.WriteBundle(req.Bundle)
	if err != nil {
		return err
	}
	res.Bundle = r
	return nil
}

// GetBundle implements the BundleServiceHandler interface
func (g Service) GetBundle(c context.Context, req *proto.GetBundleRequest, res *proto.GetBundleResponse) error {
	if validationError := validateGetBundle(req); validationError != nil {
		return validationError
	}
	r, err := g.manager.ReadBundle(req.BundleId)
	if err != nil {
		return err
	}
	res.Bundle = r
	return nil
}

// ListBundles implements the BundleServiceHandler interface
func (g Service) ListBundles(c context.Context, req *proto.ListBundlesRequest, res *proto.ListBundlesResponse) error {
	// fetch all bundles
	req.AccountUuid = getValidatedAccountUUID(c, req.AccountUuid)
	if validationError := validateListBundles(req); validationError != nil {
		return validationError
	}
	bundles, err := g.manager.ListBundles(proto.Bundle_TYPE_DEFAULT)
	if err != nil {
		return err
	}
	res.Bundles = bundles
	return nil
}

// SaveValue implements the ValueServiceHandler interface
func (g Service) SaveValue(c context.Context, req *proto.SaveValueRequest, res *proto.SaveValueResponse) error {
	req.Value.AccountUuid = getValidatedAccountUUID(c, req.Value.AccountUuid)
	cleanUpResource(c, req.Value.Resource)
	// TODO: we need to check, if the authenticated user has permission to write the value for the specified resource (e.g. global, file with id xy, ...)
	if validationError := validateSaveValue(req); validationError != nil {
		return validationError
	}
	r, err := g.manager.WriteValue(req.Value)
	if err != nil {
		return err
	}
	valueWithIdentifier, err := g.getValueWithIdentifier(r)
	if err != nil {
		return err
	}
	res.Value = valueWithIdentifier
	return nil
}

// GetValue implements the ValueServiceHandler interface
func (g Service) GetValue(c context.Context, req *proto.GetValueRequest, res *proto.GetValueResponse) error {
	if validationError := validateGetValue(req); validationError != nil {
		return validationError
	}
	r, err := g.manager.ReadValue(req.Id)
	if err != nil {
		return err
	}
	valueWithIdentifier, err := g.getValueWithIdentifier(r)
	if err != nil {
		return err
	}
	res.Value = valueWithIdentifier
	return nil
}

// GetValueByUniqueIdentifiers implements the ValueService interface
func (g Service) GetValueByUniqueIdentifiers(ctx context.Context, in *proto.GetValueByUniqueIdentifiersRequest, res *proto.GetValueResponse) error {
	v, err := g.manager.ReadValueByUniqueIdentifiers(in.AccountUuid, in.SettingId)
	if err != nil {
		return err
	}

	if v.BundleId != "" {
		valueWithIdentifier, err := g.getValueWithIdentifier(v)
		if err != nil {
			return err
		}

		res.Value = valueWithIdentifier
	}
	return nil
}

// ListValues implements the ValueServiceHandler interface
func (g Service) ListValues(c context.Context, req *proto.ListValuesRequest, res *proto.ListValuesResponse) error {
	req.AccountUuid = getValidatedAccountUUID(c, req.AccountUuid)
	if validationError := validateListValues(req); validationError != nil {
		return validationError
	}
	r, err := g.manager.ListValues(req.BundleId, req.AccountUuid)
	if err != nil {
		return err
	}
	var result []*proto.ValueWithIdentifier
	for _, value := range r {
		valueWithIdentifier, err := g.getValueWithIdentifier(value)
		if err == nil {
			result = append(result, valueWithIdentifier)
		}
	}
	res.Values = result
	return nil
}

func (g Service) getValueWithIdentifier(value *proto.Value) (*proto.ValueWithIdentifier, error) {
	bundle, err := g.manager.ReadBundle(value.BundleId)
	if err != nil {
		return nil, err
	}
	setting, err := g.manager.ReadSetting(value.SettingId)
	if err != nil {
		return nil, err
	}
	return &proto.ValueWithIdentifier{
		Identifier: &proto.Identifier{
			Extension: bundle.Extension,
			Bundle:    bundle.Name,
			Setting:   setting.Name,
		},
		Value: value,
	}, nil
}

// cleanUpResource makes sure that the account uuid of the authenticated user is injected if needed.
func cleanUpResource(c context.Context, resource *proto.Resource) {
	if resource != nil && resource.Type == proto.Resource_TYPE_USER {
		resource.Id = getValidatedAccountUUID(c, resource.Id)
	}
}

// getValidatedAccountUUID converts `me` into an actual account uuid from the context, if possible.
// the result of this function will always be a valid lower-case UUID or an empty string.
func getValidatedAccountUUID(c context.Context, accountUUID string) string {
	if accountUUID == "me" {
		if ownAccountUUID, ok := c.Value(middleware.UUIDKey).(string); ok {
			accountUUID = ownAccountUUID
		}
	}
	return accountUUID
}
