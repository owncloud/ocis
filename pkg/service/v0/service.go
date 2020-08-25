package svc

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/ptypes/empty"
	merrors "github.com/micro/go-micro/v2/errors"
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
	service := Service{
		config:  cfg,
		logger:  logger,
		manager: store.New(cfg),
	}
	// FIXME: we're writing default roles per service start (i.e. twice at the moment, for http and grpc server).
	for _, role := range generateBundlesDefaultRoles() {
		bundleID := role.Extension + "." + role.Id
		// check if the role already exists
		bundle, _ := service.manager.ReadBundle(role.Id)
		if bundle != nil {
			logger.Debug().Msgf("Settings bundle %v already exists. Skipping.", bundleID)
			continue
		}
		// create the role
		_, err := service.manager.WriteBundle(role)
		if err != nil {
			logger.Error().Err(err).Msgf("Failed to register settings bundle %v", bundleID)
		}
		logger.Debug().Msgf("Successfully registered settings bundle %v", bundleID)
	}
	return service
}

// TODO: check permissions on every request

// SaveBundle implements the BundleServiceHandler interface
func (g Service) SaveBundle(c context.Context, req *proto.SaveBundleRequest, res *proto.SaveBundleResponse) error {
	cleanUpResource(c, req.Bundle.Resource)
	if validationError := validateSaveBundle(req); validationError != nil {
		return merrors.BadRequest("ocis-settings", "%s", validationError)
	}
	r, err := g.manager.WriteBundle(req.Bundle)
	if err != nil {
		return merrors.BadRequest("ocis-settings", "%s", err)
	}
	res.Bundle = r
	return nil
}

// GetBundle implements the BundleServiceHandler interface
func (g Service) GetBundle(c context.Context, req *proto.GetBundleRequest, res *proto.GetBundleResponse) error {
	accountUUID := getValidatedAccountUUID(c, "me")
	if validationError := validateGetBundle(req); validationError != nil {
		return merrors.BadRequest("ocis-settings", "%s", validationError)
	}
	bundle, err := g.manager.ReadBundle(req.BundleId)
	if err != nil {
		return merrors.NotFound("ocis-settings", "%s", err)
	}
	roleIDs := g.getRoleIDs(c, accountUUID)
	filteredBundle := g.getFilteredBundle(roleIDs, bundle)
	if len(filteredBundle.Settings) == 0 {
		err = fmt.Errorf("could not read bundle: %s", req.BundleId)
		return merrors.NotFound("ocis-settings", "%s", err)
	}
	res.Bundle = filteredBundle
	return nil
}

// ListBundles implements the BundleServiceHandler interface
func (g Service) ListBundles(c context.Context, req *proto.ListBundlesRequest, res *proto.ListBundlesResponse) error {
	// fetch all bundles
	accountUUID := getValidatedAccountUUID(c, "me")
	if validationError := validateListBundles(req); validationError != nil {
		return merrors.BadRequest("ocis-settings", "%s", validationError)
	}
	bundles, err := g.manager.ListBundles(proto.Bundle_TYPE_DEFAULT)
	if err != nil {
		return merrors.NotFound("ocis-settings", "%s", err)
	}
	roleIDs := g.getRoleIDs(c, accountUUID)

	// filter settings in bundles that are allowed according to roles
	var filteredBundles []*proto.Bundle
	for _, bundle := range bundles {
		filteredBundle := g.getFilteredBundle(roleIDs, bundle)
		if len(filteredBundle.Settings) > 0 {
			filteredBundles = append(filteredBundles, filteredBundle)
		}
	}

	res.Bundles = filteredBundles
	return nil
}

func (g Service) getFilteredBundle(roleIDs []string, bundle *proto.Bundle) *proto.Bundle {
	// check if full bundle is whitelisted
	bundleResource := &proto.Resource{
		Type: proto.Resource_TYPE_BUNDLE,
		Id:   bundle.Id,
	}
	if g.hasPermission(
		roleIDs,
		bundleResource,
		[]proto.Permission_Operation{proto.Permission_OPERATION_READ, proto.Permission_OPERATION_READWRITE},
		proto.Permission_CONSTRAINT_OWN,
	) {
		return bundle
	}

	// filter settings based on permissions
	var filteredSettings []*proto.Setting
	for _, setting := range bundle.Settings {
		settingResource := &proto.Resource{
			Type: proto.Resource_TYPE_SETTING,
			Id:   setting.Id,
		}
		if g.hasPermission(
			roleIDs,
			settingResource,
			[]proto.Permission_Operation{proto.Permission_OPERATION_READ, proto.Permission_OPERATION_READWRITE},
			proto.Permission_CONSTRAINT_OWN,
		) {
			filteredSettings = append(filteredSettings, setting)
		}
	}
	bundle.Settings = filteredSettings
	return bundle
}

// AddSettingToBundle implements the BundleServiceHandler interface
func (g Service) AddSettingToBundle(c context.Context, req *proto.AddSettingToBundleRequest, res *proto.AddSettingToBundleResponse) error {
	cleanUpResource(c, req.Setting.Resource)
	if validationError := validateAddSettingToBundle(req); validationError != nil {
		return merrors.BadRequest("ocis-settings", "%s", validationError)
	}
	r, err := g.manager.AddSettingToBundle(req.BundleId, req.Setting)
	if err != nil {
		return merrors.BadRequest("ocis-settings", "%s", err)
	}
	res.Setting = r
	return nil
}

// RemoveSettingFromBundle implements the BundleServiceHandler interface
func (g Service) RemoveSettingFromBundle(c context.Context, req *proto.RemoveSettingFromBundleRequest, _ *empty.Empty) error {
	if validationError := validateRemoveSettingFromBundle(req); validationError != nil {
		return merrors.BadRequest("ocis-settings", "%s", validationError)
	}
	if err := g.manager.RemoveSettingFromBundle(req.BundleId, req.SettingId); err != nil {
		return merrors.BadRequest("ocis-settings", "%s", err)
	}

	return nil
}

// SaveValue implements the ValueServiceHandler interface
func (g Service) SaveValue(c context.Context, req *proto.SaveValueRequest, res *proto.SaveValueResponse) error {
	req.Value.AccountUuid = getValidatedAccountUUID(c, req.Value.AccountUuid)
	cleanUpResource(c, req.Value.Resource)
	// TODO: we need to check, if the authenticated user has permission to write the value for the specified resource (e.g. global, file with id xy, ...)
	if validationError := validateSaveValue(req); validationError != nil {
		return merrors.BadRequest("ocis-settings", "%s", validationError)
	}
	r, err := g.manager.WriteValue(req.Value)
	if err != nil {
		return merrors.BadRequest("ocis-settings", "%s", err)
	}
	valueWithIdentifier, err := g.getValueWithIdentifier(r)
	if err != nil {
		return merrors.NotFound("ocis-settings", "%s", err)
	}
	res.Value = valueWithIdentifier
	return nil
}

// GetValue implements the ValueServiceHandler interface
func (g Service) GetValue(c context.Context, req *proto.GetValueRequest, res *proto.GetValueResponse) error {
	if validationError := validateGetValue(req); validationError != nil {
		return merrors.BadRequest("ocis-settings", "%s", validationError)
	}
	r, err := g.manager.ReadValue(req.Id)
	if err != nil {
		return merrors.NotFound("ocis-settings", "%s", err)
	}
	valueWithIdentifier, err := g.getValueWithIdentifier(r)
	if err != nil {
		return merrors.NotFound("ocis-settings", "%s", err)
	}
	res.Value = valueWithIdentifier
	return nil
}

// GetValueByUniqueIdentifiers implements the ValueService interface
func (g Service) GetValueByUniqueIdentifiers(ctx context.Context, in *proto.GetValueByUniqueIdentifiersRequest, res *proto.GetValueResponse) error {
	v, err := g.manager.ReadValueByUniqueIdentifiers(in.AccountUuid, in.SettingId)
	if err != nil {
		return merrors.NotFound("ocis-settings", "%s", err)
	}

	if v.BundleId != "" {
		valueWithIdentifier, err := g.getValueWithIdentifier(v)
		if err != nil {
			return merrors.NotFound("ocis-settings", "%s", err)
		}

		res.Value = valueWithIdentifier
	}
	return nil
}

// ListValues implements the ValueServiceHandler interface
func (g Service) ListValues(c context.Context, req *proto.ListValuesRequest, res *proto.ListValuesResponse) error {
	req.AccountUuid = getValidatedAccountUUID(c, req.AccountUuid)
	if validationError := validateListValues(req); validationError != nil {
		return merrors.BadRequest("ocis-settings", "%s", validationError)
	}
	r, err := g.manager.ListValues(req.BundleId, req.AccountUuid)
	if err != nil {
		return merrors.NotFound("ocis-settings", "%s", err)
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

// ListRoles implements the RoleServiceHandler interface
func (g Service) ListRoles(c context.Context, req *proto.ListBundlesRequest, res *proto.ListBundlesResponse) error {
	//accountUUID := getValidatedAccountUUID(c, "me")
	if validationError := validateListRoles(req); validationError != nil {
		return merrors.BadRequest("ocis-settings", "%s", validationError)
	}
	r, err := g.manager.ListBundles(proto.Bundle_TYPE_ROLE)
	if err != nil {
		return merrors.NotFound("ocis-settings", "%s", err)
	}
	// TODO: only allow to list roles when user has account management permissions
	res.Bundles = r
	return nil
}

// ListRoleAssignments implements the RoleServiceHandler interface
func (g Service) ListRoleAssignments(c context.Context, req *proto.ListRoleAssignmentsRequest, res *proto.ListRoleAssignmentsResponse) error {
	req.AccountUuid = getValidatedAccountUUID(c, req.AccountUuid)
	if validationError := validateListRoleAssignments(req); validationError != nil {
		return merrors.BadRequest("ocis-settings", "%s", validationError)
	}
	r, err := g.manager.ListRoleAssignments(req.AccountUuid)
	if err != nil {
		return merrors.NotFound("ocis-settings", "%s", err)
	}
	res.Assignments = r
	return nil
}

// AssignRoleToUser implements the RoleServiceHandler interface
func (g Service) AssignRoleToUser(c context.Context, req *proto.AssignRoleToUserRequest, res *proto.AssignRoleToUserResponse) error {
	req.AccountUuid = getValidatedAccountUUID(c, req.AccountUuid)
	if validationError := validateAssignRoleToUser(req); validationError != nil {
		return merrors.BadRequest("ocis-settings", "%s", validationError)
	}
	r, err := g.manager.WriteRoleAssignment(req.AccountUuid, req.RoleId)
	if err != nil {
		return merrors.BadRequest("ocis-settings", "%s", err)
	}
	res.Assignment = r
	return nil
}

// RemoveRoleFromUser implements the RoleServiceHandler interface
func (g Service) RemoveRoleFromUser(c context.Context, req *proto.RemoveRoleFromUserRequest, _ *empty.Empty) error {
	if validationError := validateRemoveRoleFromUser(req); validationError != nil {
		return merrors.BadRequest("ocis-settings", "%s", validationError)
	}
	if err := g.manager.RemoveRoleAssignment(req.Id); err != nil {
		return merrors.BadRequest("ocis-settings", "%s", err)
	}
	return nil
}

// ListPermissionsByResource implements the PermissionServiceHandler interface
func (g Service) ListPermissionsByResource(c context.Context, req *proto.ListPermissionsByResourceRequest, res *proto.ListPermissionsByResourceResponse) error {
	if validationError := validateListPermissionsByResource(req); validationError != nil {
		return merrors.BadRequest("ocis-settings", "%s", validationError)
	}
	permissions, err := g.manager.ListPermissionsByResource(req.Resource, req.RoleIds)
	if err != nil {
		return merrors.BadRequest("ocis-settings", "%s", err)
	}
	res.Permissions = permissions
	return nil
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

// getRoleIDs loads the role assignments for the given accountUUID
// TODO: this should work on the context in the future, as roles are supposed to be sent within the context.
func (g Service) getRoleIDs(c context.Context, accountUUID string) []string {
	// TODO: replace this with role ids from the context
	// WIP PR: https://github.com/owncloud/ocis-proxy/pull/70
	rolesResponse := &proto.ListRoleAssignmentsResponse{}
	err := g.ListRoleAssignments(c, &proto.ListRoleAssignmentsRequest{AccountUuid: accountUUID}, rolesResponse)
	if err != nil {
		g.logger.Err(err).Str("accountUUID", accountUUID).Msg("failed to list role assignments")
		return []string{}
	}
	roleIDs := make([]string, 0)
	for _, assignment := range rolesResponse.Assignments {
		roleIDs = append(roleIDs, assignment.RoleId)
	}
	return roleIDs
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
