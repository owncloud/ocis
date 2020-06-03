package svc

import (
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/owncloud/ocis-settings/pkg/proto/v0"
)

var (
	regexForKeys = regexp.MustCompile("^[A-Za-z0-9\\-_]*$")
	keyRule = []validation.Rule{
		validation.Required,
		validation.Match(regexForKeys),
	}
	settingKeyRule = []validation.Rule{
		validation.Required,
		validation.Match(regexForKeys),
	}
	accountUuidRule = []validation.Rule{
		validation.Required,
		is.UUID,
	}
)

func validateSaveSettingsBundle(req *proto.SaveSettingsBundleRequest) error {
	if err := validateBundleIdentifier(req.SettingsBundle.Identifier); err != nil {
		return err
	}
	return nil
}

func validateGetSettingsBundle(req *proto.GetSettingsBundleRequest) error {
	if err := validateBundleIdentifier(req.Identifier); err != nil {
		return err
	}
	return nil
}

func validateListSettingsBundles(req *proto.ListSettingsBundlesRequest) error {
	return validation.ValidateStruct(
		req.Identifier,
		validation.Field(&req.Identifier.Extension, validation.Match(regexForKeys)),
	)
}

func validateSaveSettingsValue(req *proto.SaveSettingsValueRequest) error {
	if err := validateValueIdentifier(req.SettingsValue.Identifier); err != nil {
		return err
	}
	return nil
}

func validateGetSettingsValue(req *proto.GetSettingsValueRequest) error {
	if err := validateValueIdentifier(req.Identifier); err != nil {
		return err
	}
	return nil
}

func validateListSettingsValues(req *proto.ListSettingsValuesRequest) error {
	return validation.ValidateStruct(
		req.Identifier,
		validation.Field(&req.Identifier.AccountUuid, is.UUID),
		validation.Field(&req.Identifier.Extension, validation.Match(regexForKeys)),
		validation.Field(&req.Identifier.Extension, validation.When(req.Identifier.BundleKey != "", validation.Required).Else(validation.Nil)),
		validation.Field(&req.Identifier.BundleKey, validation.Match(regexForKeys)),
		validation.Field(&req.Identifier.BundleKey, validation.When(req.Identifier.SettingKey != "", validation.Required).Else(validation.Nil)),
		validation.Field(&req.Identifier.SettingKey, validation.Match(regexForKeys)),
	)
}

func validateBundleIdentifier(identifier *proto.Identifier) error {
	return validation.ValidateStruct(
		identifier,
		validation.Field(&identifier.Extension, keyRule...),
		validation.Field(&identifier.BundleKey, keyRule...),
	)
}

func validateValueIdentifier(identifier *proto.Identifier) error {
	return validation.ValidateStruct(
		identifier,
		validation.Field(&identifier.Extension, keyRule...),
		validation.Field(&identifier.BundleKey, keyRule...),
		validation.Field(&identifier.SettingKey, settingKeyRule...),
		validation.Field(&identifier.AccountUuid, accountUuidRule...),
	)
}

