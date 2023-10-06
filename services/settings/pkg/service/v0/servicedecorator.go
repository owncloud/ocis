package svc

import (
	"context"
	"strings"

	settingsmsg "github.com/owncloud/ocis/v2/protogen/gen/ocis/messages/settings/v0"
	settingssvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/settings/v0"
	"github.com/owncloud/ocis/v2/services/settings/pkg/config"
	"github.com/owncloud/ocis/v2/services/settings/pkg/settings"
	"github.com/owncloud/ocis/v2/services/settings/pkg/store/defaults"
)

var _defaultLanguage = "en"

// NewDefaultLanguageService returns a default language decorator for ServiceHandler.
func NewDefaultLanguageService(cfg *config.Config, serviceHandler settings.ServiceHandler) settings.ServiceHandler {
	defaultLanguage := cfg.DefaultLanguage
	if defaultLanguage == "" {
		defaultLanguage = _defaultLanguage
	}
	return &defaultLanguageDecorator{defaultLanguage: defaultLanguage, ServiceHandler: serviceHandler}
}

type defaultLanguageDecorator struct {
	defaultLanguage string
	settings.ServiceHandler
}

// GetValueByUniqueIdentifiers implements the ValueService interface
func (s *defaultLanguageDecorator) GetValueByUniqueIdentifiers(ctx context.Context, req *settingssvc.GetValueByUniqueIdentifiersRequest, res *settingssvc.GetValueResponse) error {
	err := s.ServiceHandler.GetValueByUniqueIdentifiers(ctx, req, res)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "not found") && req.GetSettingId() == defaults.SettingUUIDProfileLanguage && res.GetValue() == nil {
			res.Value = s.withDefaultLanguageSetting(req.AccountUuid)
			return nil
		}
		return err
	}
	return nil
}

// ListValues implements the ValueServiceHandler interface
func (s *defaultLanguageDecorator) ListValues(ctx context.Context, req *settingssvc.ListValuesRequest, res *settingssvc.ListValuesResponse) error {
	err := s.ServiceHandler.ListValues(ctx, req, res)
	if err != nil {
		return err
	}
	for _, v := range res.Values {
		if v.GetValue().GetSettingId() == defaults.SettingUUIDProfileLanguage {
			return nil
		}
	}

	res.Values = append(res.Values, s.withDefaultLanguageSetting(req.AccountUuid))
	return nil
}

func (s *defaultLanguageDecorator) withDefaultLanguageSetting(accountUUID string) *settingsmsg.ValueWithIdentifier {
	return &settingsmsg.ValueWithIdentifier{
		Identifier: &settingsmsg.Identifier{
			Extension: "ocis-accounts",
			Bundle:    "profile",
			Setting:   "language",
		},
		Value: &settingsmsg.Value{
			BundleId:    defaults.BundleUUIDProfile,
			SettingId:   defaults.SettingUUIDProfileLanguage,
			AccountUuid: accountUUID,
			Resource: &settingsmsg.Resource{
				Type: settingsmsg.Resource_TYPE_USER,
			},
			Value: &settingsmsg.Value_ListValue{
				ListValue: &settingsmsg.ListValue{Values: []*settingsmsg.ListOptionValue{
					{
						Option: &settingsmsg.ListOptionValue_StringValue{
							StringValue: s.defaultLanguage,
						},
					},
				}},
			},
		},
	}
}
