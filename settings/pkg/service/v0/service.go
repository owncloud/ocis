package svc

import (
	"context"
	"fmt"

	merrors "github.com/asim/go-micro/v3/errors"
	"github.com/asim/go-micro/v3/metadata"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/owncloud/ocis/ocis-pkg/log"
	"github.com/owncloud/ocis/ocis-pkg/middleware"
	"github.com/owncloud/ocis/ocis-pkg/roles"
	"github.com/owncloud/ocis/settings/pkg/config"
	"github.com/owncloud/ocis/settings/pkg/proto/v0"
	"github.com/owncloud/ocis/settings/pkg/settings"
	store "github.com/owncloud/ocis/settings/pkg/store/filesystem"
)

// Service represents a service.
type Service struct {
	id      string
	config  *config.Config
	logger  log.Logger
	manager settings.Manager
}

// NewService returns a service implementation for Service.
func NewService(cfg *config.Config, logger log.Logger) Service {
	service := Service{
		id:      "ocis-settings",
		config:  cfg,
		logger:  logger,
		manager: store.New(cfg),
	}
	service.RegisterDefaultRoles()
	return service
}

// RegisterDefaultRoles composes default roles and saves them. Skipped if the roles already exist.
func (g Service) RegisterDefaultRoles() {
	// FIXME: we're writing default roles per service start (i.e. twice at the moment, for http and grpc server). has to happen only once.
	for _, role := range generateBundlesDefaultRoles() {
		bundleID := role.Extension + "." + role.Id
		// check if the role already exists
		bundle, _ := g.manager.ReadBundle(role.Id)
		if bundle != nil {
			g.logger.Debug().Str("bundleID", bundleID).Msg("bundle already exists. skipping.")
			continue
		}
		// create the role
		_, err := g.manager.WriteBundle(role)
		if err != nil {
			g.logger.Error().Err(err).Str("bundleID", bundleID).Msg("failed to register bundle")
		}
		g.logger.Debug().Str("bundleID", bundleID).Msg("successfully registered bundle")
	}

	for _, req := range generatePermissionRequests() {
		_, err := g.manager.AddSettingToBundle(req.GetBundleId(), req.GetSetting())
		if err != nil {
			g.logger.Error().
				Err(err).
				Str("bundleID", req.GetBundleId()).
				Interface("setting", req.GetSetting()).
				Msg("failed to register permission")
		}
	}
}

// TODO: check permissions on every request

// SaveBundle implements the BundleServiceHandler interface
func (g Service) SaveBundle(ctx context.Context, req *proto.SaveBundleRequest, res *proto.SaveBundleResponse) error {
	cleanUpResource(ctx, req.Bundle.Resource)
	if err := g.checkStaticPermissionsByBundleType(ctx, req.Bundle.Type); err != nil {
		return err
	}
	if validationError := validateSaveBundle(req); validationError != nil {
		return merrors.BadRequest(g.id, "%s", validationError)
	}

	r, err := g.manager.WriteBundle(req.Bundle)
	if err != nil {
		return merrors.BadRequest(g.id, "%s", err)
	}
	res.Bundle = r
	return nil
}

// GetBundle implements the BundleServiceHandler interface
func (g Service) GetBundle(ctx context.Context, req *proto.GetBundleRequest, res *proto.GetBundleResponse) error {
	if validationError := validateGetBundle(req); validationError != nil {
		return merrors.BadRequest(g.id, "%s", validationError)
	}
	bundle, err := g.manager.ReadBundle(req.BundleId)
	if err != nil {
		return merrors.NotFound(g.id, "%s", err)
	}
	filteredBundle := g.getFilteredBundle(g.getRoleIDs(ctx), bundle)
	if len(filteredBundle.Settings) == 0 {
		err = fmt.Errorf("could not read bundle: %s", req.BundleId)
		return merrors.NotFound(g.id, "%s", err)
	}
	res.Bundle = filteredBundle
	return nil
}

// ListBundles implements the BundleServiceHandler interface
func (g Service) ListBundles(ctx context.Context, req *proto.ListBundlesRequest, res *proto.ListBundlesResponse) error {
	// fetch all bundles
	if validationError := validateListBundles(req); validationError != nil {
		return merrors.BadRequest(g.id, "%s", validationError)
	}
	bundles, err := g.manager.ListBundles(proto.Bundle_TYPE_DEFAULT, req.BundleIds)
	if err != nil {
		return merrors.NotFound(g.id, "%s", err)
	}
	roleIDs := g.getRoleIDs(ctx)

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
func (g Service) AddSettingToBundle(ctx context.Context, req *proto.AddSettingToBundleRequest, res *proto.AddSettingToBundleResponse) error {
	cleanUpResource(ctx, req.Setting.Resource)
	if err := g.checkStaticPermissionsByBundleID(ctx, req.BundleId); err != nil {
		return err
	}
	if validationError := validateAddSettingToBundle(req); validationError != nil {
		return merrors.BadRequest(g.id, "%s", validationError)
	}

	r, err := g.manager.AddSettingToBundle(req.BundleId, req.Setting)
	if err != nil {
		return merrors.BadRequest(g.id, "%s", err)
	}
	res.Setting = r
	return nil
}

// RemoveSettingFromBundle implements the BundleServiceHandler interface
func (g Service) RemoveSettingFromBundle(ctx context.Context, req *proto.RemoveSettingFromBundleRequest, _ *empty.Empty) error {
	if err := g.checkStaticPermissionsByBundleID(ctx, req.BundleId); err != nil {
		return err
	}
	if validationError := validateRemoveSettingFromBundle(req); validationError != nil {
		return merrors.BadRequest(g.id, "%s", validationError)
	}

	if err := g.manager.RemoveSettingFromBundle(req.BundleId, req.SettingId); err != nil {
		return merrors.BadRequest(g.id, "%s", err)
	}

	return nil
}

// SaveValue implements the ValueServiceHandler interface
func (g Service) SaveValue(ctx context.Context, req *proto.SaveValueRequest, res *proto.SaveValueResponse) error {
	req.Value.AccountUuid = getValidatedAccountUUID(ctx, req.Value.AccountUuid)
	cleanUpResource(ctx, req.Value.Resource)
	// TODO: we need to check, if the authenticated user has permission to write the value for the specified resource (e.g. global, file with id xy, ...)
	if validationError := validateSaveValue(req); validationError != nil {
		return merrors.BadRequest(g.id, "%s", validationError)
	}
	r, err := g.manager.WriteValue(req.Value)
	if err != nil {
		return merrors.BadRequest(g.id, "%s", err)
	}
	valueWithIdentifier, err := g.getValueWithIdentifier(r)
	if err != nil {
		return merrors.NotFound(g.id, "%s", err)
	}
	res.Value = valueWithIdentifier
	return nil
}

// GetValue implements the ValueServiceHandler interface
func (g Service) GetValue(ctx context.Context, req *proto.GetValueRequest, res *proto.GetValueResponse) error {
	if validationError := validateGetValue(req); validationError != nil {
		return merrors.BadRequest(g.id, "%s", validationError)
	}
	r, err := g.manager.ReadValue(req.Id)
	if err != nil {
		return merrors.NotFound(g.id, "%s", err)
	}
	valueWithIdentifier, err := g.getValueWithIdentifier(r)
	if err != nil {
		return merrors.NotFound(g.id, "%s", err)
	}
	res.Value = valueWithIdentifier
	return nil
}

// GetValueByUniqueIdentifiers implements the ValueService interface
func (g Service) GetValueByUniqueIdentifiers(ctx context.Context, req *proto.GetValueByUniqueIdentifiersRequest, res *proto.GetValueResponse) error {
	if validationError := validateGetValueByUniqueIdentifiers(req); validationError != nil {
		return merrors.BadRequest(g.id, "%s", validationError)
	}
	v, err := g.manager.ReadValueByUniqueIdentifiers(req.AccountUuid, req.SettingId)
	if err != nil {
		return merrors.NotFound(g.id, "%s", err)
	}

	if v.BundleId != "" {
		valueWithIdentifier, err := g.getValueWithIdentifier(v)
		if err != nil {
			return merrors.NotFound(g.id, "%s", err)
		}

		res.Value = valueWithIdentifier
	}
	return nil
}

// ListValues implements the ValueServiceHandler interface
func (g Service) ListValues(ctx context.Context, req *proto.ListValuesRequest, res *proto.ListValuesResponse) error {
	req.AccountUuid = getValidatedAccountUUID(ctx, req.AccountUuid)
	if validationError := validateListValues(req); validationError != nil {
		return merrors.BadRequest(g.id, "%s", validationError)
	}
	r, err := g.manager.ListValues(req.BundleId, req.AccountUuid)
	if err != nil {
		return merrors.NotFound(g.id, "%s", err)
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
		return merrors.BadRequest(g.id, "%s", validationError)
	}
	r, err := g.manager.ListBundles(proto.Bundle_TYPE_ROLE, req.BundleIds)
	if err != nil {
		return merrors.NotFound(g.id, "%s", err)
	}
	// TODO: only allow to list roles when user has account/role/... management permissions
	res.Bundles = r
	return nil
}

// ListRoleAssignments implements the RoleServiceHandler interface
func (g Service) ListRoleAssignments(ctx context.Context, req *proto.ListRoleAssignmentsRequest, res *proto.ListRoleAssignmentsResponse) error {
	req.AccountUuid = getValidatedAccountUUID(ctx, req.AccountUuid)
	if validationError := validateListRoleAssignments(req); validationError != nil {
		return merrors.BadRequest(g.id, "%s", validationError)
	}
	r, err := g.manager.ListRoleAssignments(req.AccountUuid)
	if err != nil {
		return merrors.NotFound(g.id, "%s", err)
	}
	res.Assignments = r
	return nil
}

// AssignRoleToUser implements the RoleServiceHandler interface
func (g Service) AssignRoleToUser(ctx context.Context, req *proto.AssignRoleToUserRequest, res *proto.AssignRoleToUserResponse) error {
	if err := g.checkStaticPermissionsByBundleType(ctx, proto.Bundle_TYPE_ROLE); err != nil {
		return err
	}

	req.AccountUuid = getValidatedAccountUUID(ctx, req.AccountUuid)
	if validationError := validateAssignRoleToUser(req); validationError != nil {
		return merrors.BadRequest(g.id, "%s", validationError)
	}
	r, err := g.manager.WriteRoleAssignment(req.AccountUuid, req.RoleId)
	if err != nil {
		return merrors.BadRequest(g.id, "%s", err)
	}
	res.Assignment = r
	return nil
}

// RemoveRoleFromUser implements the RoleServiceHandler interface
func (g Service) RemoveRoleFromUser(ctx context.Context, req *proto.RemoveRoleFromUserRequest, _ *empty.Empty) error {
	if err := g.checkStaticPermissionsByBundleType(ctx, proto.Bundle_TYPE_ROLE); err != nil {
		return err
	}

	if validationError := validateRemoveRoleFromUser(req); validationError != nil {
		return merrors.BadRequest(g.id, "%s", validationError)
	}
	if err := g.manager.RemoveRoleAssignment(req.Id); err != nil {
		return merrors.BadRequest(g.id, "%s", err)
	}
	return nil
}

// ListPermissionsByResource implements the PermissionServiceHandler interface
func (g Service) ListPermissionsByResource(ctx context.Context, req *proto.ListPermissionsByResourceRequest, res *proto.ListPermissionsByResourceResponse) error {
	if validationError := validateListPermissionsByResource(req); validationError != nil {
		return merrors.BadRequest(g.id, "%s", validationError)
	}
	permissions, err := g.manager.ListPermissionsByResource(req.Resource, g.getRoleIDs(ctx))
	if err != nil {
		return merrors.BadRequest(g.id, "%s", err)
	}
	res.Permissions = permissions
	return nil
}

// GetPermissionByID implements the PermissionServiceHandler interface
func (g Service) GetPermissionByID(ctx context.Context, req *proto.GetPermissionByIDRequest, res *proto.GetPermissionByIDResponse) error {
	if validationError := validateGetPermissionByID(req); validationError != nil {
		return merrors.BadRequest(g.id, "%s", validationError)
	}
	permission, err := g.manager.ReadPermissionByID(req.PermissionId, g.getRoleIDs(ctx))
	if err != nil {
		return merrors.BadRequest(g.id, "%s", err)
	}
	if permission == nil {
		return merrors.NotFound(g.id, "%s", fmt.Errorf("permission %s not found in roles", req.PermissionId))
	}
	res.Permission = permission
	return nil
}

// cleanUpResource makes sure that the account uuid of the authenticated user is injected if needed.
func cleanUpResource(ctx context.Context, resource *proto.Resource) {
	if resource != nil && resource.Type == proto.Resource_TYPE_USER {
		resource.Id = getValidatedAccountUUID(ctx, resource.Id)
	}
}

// getValidatedAccountUUID converts `me` into an actual account uuid from the context, if possible.
// the result of this function will always be a valid lower-case UUID or an empty string.
func getValidatedAccountUUID(ctx context.Context, accountUUID string) string {
	if accountUUID == "me" {
		if ownAccountUUID, ok := metadata.Get(ctx, middleware.AccountID); ok {
			accountUUID = ownAccountUUID
		}
	}
	if accountUUID == "me" {
		// no matter what happens above, an accountUUID of `me` must not be passed on. Clear it instead.
		accountUUID = ""
	}
	return accountUUID
}

// getRoleIDs extracts the roleIDs of the authenticated user from the context.
func (g Service) getRoleIDs(ctx context.Context) []string {
	if ownRoleIDs, ok := roles.ReadRoleIDsFromContext(ctx); ok {
		return ownRoleIDs
	}
	return []string{}
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

func (g Service) hasStaticPermission(ctx context.Context, permissionID string) bool {
	roleIDs, ok := roles.ReadRoleIDsFromContext(ctx)
	if !ok {
		/**
		* FIXME: with this we are skipping permission checks on all requests that are coming in without roleIDs in the
		* metadata context. This is a huge security impairment, as that's the case not only for grpc requests but also
		* for unauthenticated http requests and http requests coming in without hitting the ocis-proxy first.
		 */
		// TODO add system role for internal requests.
		// - at least the proxy needs to look up account info
		// - glauth needs to make bind requests
		// tracked as OCIS-454
		return true
	}
	p, err := g.manager.ReadPermissionByID(permissionID, roleIDs)
	return err == nil && p != nil
}

func (g Service) checkStaticPermissionsByBundleID(ctx context.Context, bundleID string) error {
	bundle, err := g.manager.ReadBundle(bundleID)
	if err != nil {
		return merrors.NotFound(g.id, "bundle not found: %s", err)
	}
	return g.checkStaticPermissionsByBundleType(ctx, bundle.Type)
}

func (g Service) checkStaticPermissionsByBundleType(ctx context.Context, bundleType proto.Bundle_Type) error {
	if bundleType == proto.Bundle_TYPE_ROLE {
		if !g.hasStaticPermission(ctx, RoleManagementPermissionID) {
			return merrors.Forbidden(g.id, "user has no role management permission")
		}
		return nil
	}
	if !g.hasStaticPermission(ctx, SettingsManagementPermissionID) {
		return merrors.Forbidden(g.id, "user has no settings management permission")
	}
	return nil
}
