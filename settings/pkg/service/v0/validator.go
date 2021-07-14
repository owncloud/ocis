package svc

import (
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/owncloud/ocis/settings/pkg/proto/v0"
)

var (
	regexForAccountUUID = regexp.MustCompile(`^[A-Za-z0-9\-_.+@:]+$`)
	requireAccountID    = []validation.Rule{
		// use rule for validation error message consistency (".. must not be blank" on empty strings)
		validation.Required,
		validation.Match(regexForAccountUUID),
	}
	regexForKeys        = regexp.MustCompile(`^[A-Za-z0-9\-_]*$`)
	requireAlphanumeric = []validation.Rule{
		validation.Required,
		validation.Match(regexForKeys),
	}
)

func validateSaveBundle(req *proto.SaveBundleRequest) error {
	if err := validation.ValidateStruct(
		req.Bundle,
		validation.Field(&req.Bundle.Id, validation.When(req.Bundle.Id != "", is.UUID)),
		validation.Field(&req.Bundle.Name, requireAlphanumeric...),
		validation.Field(&req.Bundle.Type, validation.NotIn(proto.Bundle_TYPE_UNKNOWN)),
		validation.Field(&req.Bundle.Extension, requireAlphanumeric...),
		validation.Field(&req.Bundle.DisplayName, validation.Required),
		validation.Field(&req.Bundle.Settings, validation.Required),
	); err != nil {
		return err
	}
	if err := validateResource(req.Bundle.Resource); err != nil {
		return err
	}
	for i := range req.Bundle.Settings {
		if err := validateSetting(req.Bundle.Settings[i]); err != nil {
			return err
		}
	}
	return nil
}

func validateGetBundle(req *proto.GetBundleRequest) error {
	return validation.Validate(&req.BundleId, is.UUID)
}

func validateListBundles(req *proto.ListBundlesRequest) error {
	return nil
}

func validateAddSettingToBundle(req *proto.AddSettingToBundleRequest) error {
	if err := validation.ValidateStruct(req, validation.Field(&req.BundleId, is.UUID)); err != nil {
		return err
	}
	return validateSetting(req.Setting)
}

func validateRemoveSettingFromBundle(req *proto.RemoveSettingFromBundleRequest) error {
	return validation.ValidateStruct(
		req,
		validation.Field(&req.BundleId, is.UUID),
		validation.Field(&req.SettingId, is.UUID),
	)
}

func validateSaveValue(req *proto.SaveValueRequest) error {
	if err := validation.ValidateStruct(
		req.Value,
		validation.Field(&req.Value.Id, validation.When(req.Value.Id != "", is.UUID)),
		validation.Field(&req.Value.BundleId, is.UUID),
		validation.Field(&req.Value.SettingId, is.UUID),
		validation.Field(&req.Value.AccountUuid, requireAccountID...),
	); err != nil {
		return err
	}

	if err := validateResource(req.Value.Resource); err != nil {
		return err
	}

	// TODO: validate values against the respective setting. need to check if constraints of the setting are fulfilled.
	return nil
}

func validateGetValue(req *proto.GetValueRequest) error {
	return validation.Validate(req.Id, is.UUID)
}

func validateGetValueByUniqueIdentifiers(req *proto.GetValueByUniqueIdentifiersRequest) error {
	return validation.ValidateStruct(
		req,
		validation.Field(&req.SettingId, is.UUID),
		validation.Field(&req.AccountUuid, requireAccountID...),
	)
}

func validateListValues(req *proto.ListValuesRequest) error {
	return validation.ValidateStruct(
		req,
		validation.Field(&req.BundleId, validation.When(req.BundleId != "", is.UUID)),
		validation.Field(&req.AccountUuid, validation.When(req.AccountUuid != "", validation.Match(regexForAccountUUID))),
	)
}

func validateListRoles(req *proto.ListBundlesRequest) error {
	return nil
}

func validateListRoleAssignments(req *proto.ListRoleAssignmentsRequest) error {
	return validation.Validate(req.AccountUuid, requireAccountID...)
}

func validateAssignRoleToUser(req *proto.AssignRoleToUserRequest) error {
	return validation.ValidateStruct(
		req,
		validation.Field(&req.AccountUuid, requireAccountID...),
		validation.Field(&req.RoleId, is.UUID),
	)
}

func validateRemoveRoleFromUser(req *proto.RemoveRoleFromUserRequest) error {
	return validation.ValidateStruct(
		req,
		validation.Field(&req.Id, is.UUID),
	)
}

func validateListPermissionsByResource(req *proto.ListPermissionsByResourceRequest) error {
	return validateResource(req.Resource)
}

func validateGetPermissionByID(req *proto.GetPermissionByIDRequest) error {
	return validation.ValidateStruct(
		req,
		validation.Field(&req.PermissionId, requireAlphanumeric...),
	)
}

// validateResource is an internal helper for validating the content of a resource.
func validateResource(resource *proto.Resource) error {
	if err := validation.Validate(&resource, validation.Required); err != nil {
		return err
	}
	return validation.Validate(&resource, validation.NotIn(proto.Resource_TYPE_UNKNOWN))
}

// validateSetting is an internal helper for validating the content of a setting.
func validateSetting(setting *proto.Setting) error {
	// TODO: make sanity checks, like for int settings, min <= default <= max.
	if err := validation.ValidateStruct(
		setting,
		validation.Field(&setting.Id, validation.When(setting.Id != "", is.UUID)),
		validation.Field(&setting.Name, requireAlphanumeric...),
	); err != nil {
		return err
	}
	return validateResource(setting.Resource)
}
