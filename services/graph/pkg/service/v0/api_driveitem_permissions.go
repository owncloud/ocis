package svc

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"slices"
	"strings"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	grouppb "github.com/cs3org/go-cs3apis/cs3/identity/group/v1beta1"
	userpb "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	ocmprovider "github.com/cs3org/go-cs3apis/cs3/ocm/provider/v1beta1"
	collaboration "github.com/cs3org/go-cs3apis/cs3/sharing/collaboration/v1beta1"
	link "github.com/cs3org/go-cs3apis/cs3/sharing/link/v1beta1"
	ocm "github.com/cs3org/go-cs3apis/cs3/sharing/ocm/v1beta1"
	storageprovider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	types "github.com/cs3org/go-cs3apis/cs3/types/v1beta1"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	libregraph "github.com/owncloud/libre-graph-api-go"
	"go.opentelemetry.io/otel/trace"

	revactx "github.com/owncloud/reva/v2/pkg/ctx"
	"github.com/owncloud/reva/v2/pkg/errtypes"
	"github.com/owncloud/reva/v2/pkg/publicshare"
	"github.com/owncloud/reva/v2/pkg/rgrpc/todo/pool"
	"github.com/owncloud/reva/v2/pkg/share"
	"github.com/owncloud/reva/v2/pkg/storagespace"
	"github.com/owncloud/reva/v2/pkg/utils"

	"github.com/owncloud/ocis/v2/ocis-pkg/l10n"
	settingssvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/settings/v0"
	l10n_pkg "github.com/owncloud/ocis/v2/services/graph/pkg/l10n"

	"github.com/owncloud/ocis/v2/ocis-pkg/conversions"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/services/graph/pkg/config"
	"github.com/owncloud/ocis/v2/services/graph/pkg/errorcode"
	"github.com/owncloud/ocis/v2/services/graph/pkg/identity"
	"github.com/owncloud/ocis/v2/services/graph/pkg/unifiedrole"
	"github.com/owncloud/ocis/v2/services/graph/pkg/validate"
)

const (
	invalidIdMsg       = "invalid driveID or itemID"
	parseDriveIDErrMsg = "could not parse driveID"
)

// DriveItemPermissionsProvider contains the methods related to handling permissions on drive items
type DriveItemPermissionsProvider interface {
	Invite(ctx context.Context, resourceId *storageprovider.ResourceId, invite libregraph.DriveItemInvite) (libregraph.Permission, error)
	SpaceRootInvite(ctx context.Context, driveID *storageprovider.ResourceId, invite libregraph.DriveItemInvite) (libregraph.Permission, error)
	ListPermissions(ctx context.Context, itemID *storageprovider.ResourceId, listFederatedRoles, selectRoles bool) (libregraph.CollectionOfPermissionsWithAllowedValues, error)
	ListSpaceRootPermissions(ctx context.Context, driveID *storageprovider.ResourceId) (libregraph.CollectionOfPermissionsWithAllowedValues, error)
	DeletePermission(ctx context.Context, itemID *storageprovider.ResourceId, permissionID string) error
	DeleteSpaceRootPermission(ctx context.Context, driveID *storageprovider.ResourceId, permissionID string) error
	UpdatePermission(ctx context.Context, itemID *storageprovider.ResourceId, permissionID string, newPermission libregraph.Permission) (libregraph.Permission, error)
	UpdateSpaceRootPermission(ctx context.Context, driveID *storageprovider.ResourceId, permissionID string, newPermission libregraph.Permission) (libregraph.Permission, error)
	CreateLink(ctx context.Context, driveItemID *storageprovider.ResourceId, createLink libregraph.DriveItemCreateLink) (libregraph.Permission, error)
	CreateSpaceRootLink(ctx context.Context, driveID *storageprovider.ResourceId, createLink libregraph.DriveItemCreateLink) (libregraph.Permission, error)
	SetPublicLinkPassword(ctx context.Context, driveItemID *storageprovider.ResourceId, permissionID string, password string) (libregraph.Permission, error)
	SetPublicLinkPasswordOnSpaceRoot(ctx context.Context, driveID *storageprovider.ResourceId, permissionID string, password string) (libregraph.Permission, error)
	GetPermission(ctx context.Context, itemID *storageprovider.ResourceId, permissionID string) (libregraph.Permission, error)
	GetSpaceRootPermission(ctx context.Context, driveID *storageprovider.ResourceId, permissionID string) (libregraph.Permission, error)
}

// DriveItemPermissionsService contains the production business logic for everything that relates to permissions on drive items.
type DriveItemPermissionsService struct {
	BaseGraphService
	tp              trace.TracerProvider
	identityBackend identity.Backend
}

type permissionType int

const (
	Unknown permissionType = iota
	Public
	User
	Space
	OCM
)

// NewDriveItemPermissionsService creates a new DriveItemPermissionsService
func NewDriveItemPermissionsService(logger log.Logger, gatewaySelector pool.Selectable[gateway.GatewayAPIClient], identityCache identity.IdentityCache, config *config.Config, tp trace.TracerProvider, be identity.Backend) (DriveItemPermissionsService, error) {
	return DriveItemPermissionsService{
		BaseGraphService: BaseGraphService{
			logger:          &log.Logger{Logger: logger.With().Str("graph api", "DrivesDriveItemService").Logger()},
			gatewaySelector: gatewaySelector,
			identityCache:   identityCache,
			config:          config,
		},
		tp:              tp,
		identityBackend: be,
	}, nil
}

// Invite invites a user to a drive item.
func (s DriveItemPermissionsService) Invite(ctx context.Context, resourceId *storageprovider.ResourceId, invite libregraph.DriveItemInvite) (libregraph.Permission, error) {

	ctx, span := s.tp.Tracer("graph").Start(ctx, "Invite")
	defer span.End()

	gatewayClient, err := s.gatewaySelector.Next()
	if err != nil {
		s.logger.Error().Err(err).Str("resourceId", resourceId.GetStorageId()).Interface("gateway.res", gatewayClient).Msg("gatewaySelector failed")
		return libregraph.Permission{}, err
	}

	statResponse, err := gatewayClient.Stat(ctx, &storageprovider.StatRequest{Ref: &storageprovider.Reference{ResourceId: resourceId}})
	if err := errorcode.FromStat(statResponse, err); err != nil {
		s.logger.Error().Err(err).Str("resourceId", resourceId.GetStorageId()).Interface("stat.res", statResponse).Msg("stat failed")
		return libregraph.Permission{}, err
	}

	var condition string
	if condition, err = roleConditionForResourceType(statResponse.GetInfo()); err != nil {
		s.logger.Error().Err(err).Str("resourceId", resourceId.GetStorageId()).Interface("condition.res", condition).Msg("role condition failed")
		return libregraph.Permission{}, err
	}

	unifiedRolePermissions := []*libregraph.UnifiedRolePermission{{AllowedResourceActions: invite.LibreGraphPermissionsActions}}
	for _, roleID := range invite.GetRoles() {
		// only allow roles that are enabled in the config
		if !slices.Contains(s.config.UnifiedRoles.AvailableRoles, roleID) {
			err = unifiedrole.ErrUnknownRole
			s.logger.Error().Err(err).Str("resourceId", resourceId.GetStorageId()).Str("roleID", roleID).Msg("unknown role")
			return libregraph.Permission{}, err
		}

		role, err := unifiedrole.GetRole(unifiedrole.RoleFilterIDs(roleID))
		if err != nil {
			s.logger.Error().Err(err).Str("resourceId", resourceId.GetStorageId()).Str("roleID", roleID).Interface("role", invite.GetRoles()[0]).Msg("unable to convert requested role")
			return libregraph.Permission{}, err
		}

		allowedResourceActions := unifiedrole.GetAllowedResourceActions(role, condition)
		if len(allowedResourceActions) == 0 && role.GetId() != unifiedrole.UnifiedRoleDeniedID {
			err = errorcode.New(errorcode.InvalidRequest, "role not applicable to this resource")
			s.logger.Error().Err(err).Str("resourceId", resourceId.GetStorageId()).Str("roleID", roleID).Msg(err.Error())
			return libregraph.Permission{}, err
		}

		unifiedRolePermissions = append(unifiedRolePermissions, conversions.ToPointerSlice(role.GetRolePermissions())...)
	}

	driveRecipient := invite.GetRecipients()[0]

	objectID := driveRecipient.GetObjectId()
	cs3ResourcePermissions := unifiedrole.PermissionsToCS3ResourcePermissions(unifiedRolePermissions)

	permission := &libregraph.Permission{}
	availableRoles := unifiedrole.GetRoles(unifiedrole.RoleFilterIDs(s.config.UnifiedRoles.AvailableRoles...))
	if role := unifiedrole.CS3ResourcePermissionsToRole(availableRoles, cs3ResourcePermissions, condition, false); role != nil {
		permission.Roles = []string{role.GetId()}
	}

	if len(permission.GetRoles()) == 0 {
		permission.LibreGraphPermissionsActions = unifiedrole.CS3ResourcePermissionsToLibregraphActions(cs3ResourcePermissions)
	}

	var shareid string
	var expiration *types.Timestamp
	var cTime *types.Timestamp
	switch driveRecipient.GetLibreGraphRecipientType() {
	case "group":
		group, err := s.identityCache.GetGroup(ctx, objectID)
		if err != nil {
			err = errorcode.New(errorcode.InvalidRequest, err.Error())
			s.logger.Error().Err(err).Str("resourceId", resourceId.GetStorageId()).Interface("userId", objectID).Msg("failed group lookup")
			return libregraph.Permission{}, err
		}
		permission.GrantedToV2 = &libregraph.SharePointIdentitySet{
			Group: &libregraph.Identity{
				DisplayName: group.GetDisplayName(),
				Id:          conversions.ToPointer(group.GetId()),
			},
		}
		createShareRequest := createShareRequestToGroup(group, statResponse.GetInfo(), cs3ResourcePermissions)
		if invite.ExpirationDateTime != nil {
			createShareRequest.GetGrant().Expiration = utils.TimeToTS(*invite.ExpirationDateTime)
		}
		createShareResponse, err := gatewayClient.CreateShare(ctx, createShareRequest)
		if err := errorcode.FromCS3Status(createShareResponse.GetStatus(), err); err != nil {
			s.logger.Error().Err(err).Str("resourceId", resourceId.GetStorageId()).Interface("userId", objectID).Msg("share creation failed")
			return libregraph.Permission{}, err
		}
		shareid = createShareResponse.GetShare().GetId().GetOpaqueId()
		cTime = createShareResponse.GetShare().GetCtime()
		expiration = createShareResponse.GetShare().GetExpiration()
	default:
		user, err := s.identityCache.GetUser(ctx, objectID)
		if errors.Is(err, identity.ErrNotFound) {
			if s.config.MultiInstance.Enabled {
				user, err = s.identityBackend.AddUser(ctx, objectID, s.config.MultiInstance.InstanceID)
			}
			if s.config.IncludeOCMSharees && err != nil {
				user, err = s.identityCache.GetAcceptedUser(ctx, objectID)
				if err == nil && IsSpaceRoot(statResponse.GetInfo().GetId()) {
					err = errorcode.New(errorcode.InvalidRequest, "federated user can not become a space member")
					s.logger.Error().Err(err).Str("resourceId", resourceId.GetStorageId()).Interface("userId", objectID).Msg(err.Error())
					return libregraph.Permission{}, err
				}
			}
		}
		if err != nil {
			err = errorcode.New(errorcode.InvalidRequest, err.Error())
			s.logger.Error().Err(err).Str("resourceId", resourceId.GetStorageId()).Interface("userId", objectID).Msg("failed user lookup")
			return libregraph.Permission{}, err
		}
		permission.GrantedToV2 = &libregraph.SharePointIdentitySet{
			User: &libregraph.Identity{
				DisplayName:        user.GetDisplayName(),
				Id:                 conversions.ToPointer(user.GetId()),
				LibreGraphUserType: conversions.ToPointer(user.GetUserType()),
			},
		}

		if user.GetUserType() == identity.UserTypeFederated {
			if len(user.Identities) < 1 {
				err = errorcode.New(errorcode.InvalidRequest, "user has no federated identity")
				s.logger.Error().Err(err).Str("resourceId", resourceId.GetStorageId()).Interface("userId", objectID).Msg(err.Error())
				return libregraph.Permission{}, err
			}
			providerInfoResp, err := gatewayClient.GetInfoByDomain(ctx, &ocmprovider.GetInfoByDomainRequest{
				Domain: *user.Identities[0].Issuer,
			})
			if err := errorcode.FromCS3Status(providerInfoResp.GetStatus(), err); err != nil {
				s.logger.Error().Err(err).Str("resourceId", resourceId.GetStorageId()).Interface("userId", objectID).Msg("getting provider info failed")
				return libregraph.Permission{}, err
			}

			createShareRequest := createShareRequestToFederatedUser(user, statResponse.GetInfo().GetId(), providerInfoResp.ProviderInfo, cs3ResourcePermissions)
			if invite.ExpirationDateTime != nil {
				createShareRequest.Expiration = utils.TimeToTS(*invite.ExpirationDateTime)
			}
			createShareResponse, err := gatewayClient.CreateOCMShare(ctx, createShareRequest)
			if err := errorcode.FromCS3Status(createShareResponse.GetStatus(), err); err != nil {
				s.logger.Error().Err(err).Str("resourceId", resourceId.GetStorageId()).Interface("userId", objectID).Msg("share creation failed")
				return libregraph.Permission{}, err
			}
			shareid = createShareResponse.GetShare().GetId().GetOpaqueId()
			cTime = createShareResponse.GetShare().GetCtime()
			expiration = createShareResponse.GetShare().GetExpiration()
		} else {
			createShareRequest := createShareRequestToUser(user, statResponse.GetInfo(), cs3ResourcePermissions)
			if invite.ExpirationDateTime != nil {
				createShareRequest.GetGrant().Expiration = utils.TimeToTS(*invite.ExpirationDateTime)
			}
			createShareResponse, err := gatewayClient.CreateShare(ctx, createShareRequest)
			if err := errorcode.FromCS3Status(createShareResponse.GetStatus(), err); err != nil {
				s.logger.Error().Err(err).Str("resourceId", resourceId.GetStorageId()).Interface("userId", objectID).Msg("share creation failed")
				return libregraph.Permission{}, err
			}
			shareid = createShareResponse.GetShare().GetId().GetOpaqueId()
			cTime = createShareResponse.GetShare().GetCtime()
			expiration = createShareResponse.GetShare().GetExpiration()
		}

	}

	if shareid != "" {
		permission.Id = conversions.ToPointer(shareid)
	} else if IsSpaceRoot(statResponse.GetInfo().GetId()) {
		// permissions on a space root are not handled by a share manager so
		// they don't get a share-id
		permission.SetId(identitySetToSpacePermissionID(permission.GetGrantedToV2()))
	}

	if expiration != nil {
		permission.SetExpirationDateTime(utils.TSToTime(expiration))
	}

	// set cTime
	if cTime != nil {
		permission.SetCreatedDateTime(cs3TimestampToTime(cTime))
	}

	if user, ok := revactx.ContextGetUser(ctx); ok {
		identity, err := userIdToIdentity(ctx, s.identityCache, user.GetId().GetOpaqueId())
		if err != nil {
			err = errorcode.New(errorcode.InvalidRequest, err.Error())
			s.logger.Error().Err(err).Str("resourceId", resourceId.GetStorageId()).Interface("userId", objectID).Msg("identity lookup failed")
			return libregraph.Permission{}, err
		}
		permission.SetInvitation(libregraph.SharingInvitation{
			InvitedBy: &libregraph.IdentitySet{
				User: &identity,
			},
		})
	}

	return *permission, nil
}

func createShareRequestToGroup(group libregraph.Group, info *storageprovider.ResourceInfo, cs3ResourcePermissions *storageprovider.ResourcePermissions) *collaboration.CreateShareRequest {
	return &collaboration.CreateShareRequest{
		ResourceInfo: info,
		Grant: &collaboration.ShareGrant{
			Grantee: &storageprovider.Grantee{
				Type: storageprovider.GranteeType_GRANTEE_TYPE_GROUP,
				Id: &storageprovider.Grantee_GroupId{GroupId: &grouppb.GroupId{
					OpaqueId: group.GetId(),
				}},
			},
			Permissions: &collaboration.SharePermissions{
				Permissions: cs3ResourcePermissions,
			},
		},
	}
}
func createShareRequestToUser(user libregraph.User, info *storageprovider.ResourceInfo, cs3ResourcePermissions *storageprovider.ResourcePermissions) *collaboration.CreateShareRequest {
	return &collaboration.CreateShareRequest{
		ResourceInfo: info,
		Grant: &collaboration.ShareGrant{
			Grantee: &storageprovider.Grantee{
				Type: storageprovider.GranteeType_GRANTEE_TYPE_USER,
				Id: &storageprovider.Grantee_UserId{UserId: &userpb.UserId{
					OpaqueId: user.GetId(),
				}},
			},
			Permissions: &collaboration.SharePermissions{
				Permissions: cs3ResourcePermissions,
			},
		},
	}
}
func createShareRequestToFederatedUser(user libregraph.User, resourceId *storageprovider.ResourceId, providerInfo *ocmprovider.ProviderInfo, cs3ResourcePermissions *storageprovider.ResourcePermissions) *ocm.CreateOCMShareRequest {
	return &ocm.CreateOCMShareRequest{
		ResourceId: resourceId,
		Grantee: &storageprovider.Grantee{
			Type: storageprovider.GranteeType_GRANTEE_TYPE_USER,
			Id: &storageprovider.Grantee_UserId{UserId: &userpb.UserId{
				Type:     userpb.UserType_USER_TYPE_FEDERATED,
				OpaqueId: user.GetId(),
				Idp:      *user.GetIdentities()[0].Issuer, // the domain is persisted in the grant as u:{opaqueid}:{domain}
			}},
		},
		RecipientMeshProvider: providerInfo,
		AccessMethods: []*ocm.AccessMethod{
			{
				Term: &ocm.AccessMethod_WebdavOptions{
					WebdavOptions: &ocm.WebDAVAccessMethod{
						Permissions: cs3ResourcePermissions,
					},
				},
			},
		},
	}
}

// SpaceRootInvite handles invitation request on project spaces
func (s DriveItemPermissionsService) SpaceRootInvite(ctx context.Context, driveID *storageprovider.ResourceId, invite libregraph.DriveItemInvite) (libregraph.Permission, error) {

	ctx, span := s.tp.Tracer("graph").Start(ctx, "SpaceRootInvite")
	defer span.End()

	gatewayClient, err := s.gatewaySelector.Next()
	if err != nil {
		s.logger.Error().Err(err).Str("storageID", driveID.GetStorageId()).Msg("gateway selection failed")
		return libregraph.Permission{}, err
	}

	space, err := utils.GetSpace(ctx, storagespace.FormatResourceID(driveID), gatewayClient)
	if err != nil {
		s.logger.Error().Err(err).Str("storageID", driveID.GetStorageId()).Msg("GetSpace failed")
		return libregraph.Permission{}, errorcode.FromUtilsStatusCodeError(err)
	}

	if space.SpaceType != _spaceTypeProject {
		err = errorcode.New(errorcode.InvalidRequest, "unsupported space type")
		s.logger.Error().Err(err).Str("storageID", driveID.GetStorageId()).Msg(err.Error())
		return libregraph.Permission{}, err
	}

	rootResourceID := space.GetRoot()
	return s.Invite(ctx, rootResourceID, invite)
}

// ListPermissions lists the permissions of a driveItem
func (s DriveItemPermissionsService) ListPermissions(ctx context.Context, itemID *storageprovider.ResourceId, listFederatedRoles, selectRoles bool) (libregraph.CollectionOfPermissionsWithAllowedValues, error) {

	ctx, span := s.tp.Tracer("graph").Start(ctx, "ListPermissions")
	defer span.End()

	collectionOfPermissions := libregraph.CollectionOfPermissionsWithAllowedValues{}
	gatewayClient, err := s.gatewaySelector.Next()
	if err != nil {
		s.logger.Error().Err(err).Str("storageID", itemID.GetStorageId()).Msg("gateway selection failed")
		return collectionOfPermissions, err
	}

	statResponse, err := gatewayClient.Stat(ctx, &storageprovider.StatRequest{Ref: &storageprovider.Reference{ResourceId: itemID}})
	if err := errorcode.FromStat(statResponse, err); err != nil {
		s.logger.Warn().Err(err).Str("storageID", itemID.GetStorageId()).Interface("stat.res", statResponse).Msg("stat failed")
		return collectionOfPermissions, err
	}

	var condition string
	if condition, err = roleConditionForResourceType(statResponse.GetInfo()); err != nil {
		s.logger.Error().Err(err).Str("storageID", itemID.GetStorageId()).Msg("roleConditionForResourceType failed")
		return collectionOfPermissions, err
	}

	permissionSet := statResponse.GetInfo().GetPermissionSet()
	allowedActions := unifiedrole.CS3ResourcePermissionsToLibregraphActions(permissionSet)

	collectionOfPermissions = libregraph.CollectionOfPermissionsWithAllowedValues{
		LibreGraphPermissionsActionsAllowedValues: allowedActions,
		LibreGraphPermissionsRolesAllowedValues: conversions.ToValueSlice(
			unifiedrole.GetRolesByPermissions(
				unifiedrole.GetRoles(unifiedrole.RoleFilterIDs(s.config.UnifiedRoles.AvailableRoles...)),
				allowedActions,
				condition,
				listFederatedRoles,
				false,
			),
		),
	}

	for i, definition := range collectionOfPermissions.LibreGraphPermissionsRolesAllowedValues {
		// the openapi spec defines that the rolePermissions should not be part of the response
		definition.RolePermissions = nil
		collectionOfPermissions.LibreGraphPermissionsRolesAllowedValues[i] = definition
	}

	if selectRoles {
		// drop the actions
		collectionOfPermissions.LibreGraphPermissionsActionsAllowedValues = nil
		// no need to fetch shares, we are only interested in the roles
		return collectionOfPermissions, nil
	}

	driveItems := make(driveItemsByResourceID, 1)
	// we can use the statResponse to build the drive item before fetching the shares
	item, err := cs3ResourceToDriveItem(s.logger, statResponse.GetInfo())
	if err != nil {
		return collectionOfPermissions, err
	}
	driveItems[storagespace.FormatResourceID(statResponse.GetInfo().GetId())] = *item

	if IsSpaceRoot(statResponse.GetInfo().GetId()) {
		permissions, err := s.getSpaceRootPermissions(ctx, statResponse.GetInfo().GetSpace().GetId())
		if err != nil {
			s.logger.Error().Err(err).Str("storageID", itemID.GetStorageId()).Msg("getSpaceRootPermissions failed")
			return collectionOfPermissions, err
		}
		collectionOfPermissions.Value = permissions
	} else {
		// "normal" driveItem, populate user  permissions via share providers
		driveItems, err = s.listUserShares(ctx, []*collaboration.Filter{
			share.ResourceIDFilter(itemID),
		}, driveItems)
		if err != nil {
			s.logger.Error().Err(err).Str("storageID", itemID.GetStorageId()).Msg("listUserShares failed")
			return collectionOfPermissions, err
		}
		if s.config.IncludeOCMSharees {
			driveItems, err = s.listOCMShares(ctx, []*ocm.ListOCMSharesRequest_Filter{
				{
					Type: ocm.ListOCMSharesRequest_Filter_TYPE_RESOURCE_ID,
					Term: &ocm.ListOCMSharesRequest_Filter_ResourceId{ResourceId: itemID},
				},
			}, driveItems)
			if err != nil {
				s.logger.Error().Err(err).Str("storageID", itemID.GetStorageId()).Msg("listOCMShares failed")
				return collectionOfPermissions, err
			}
		}
	}
	// finally get public shares, which are possible for spaceroots and "normal" resources
	driveItems, err = s.listPublicShares(ctx, []*link.ListPublicSharesRequest_Filter{
		publicshare.ResourceIDFilter(itemID),
	}, driveItems)
	if err != nil && strings.Contains(err.Error(), errtypes.ERR_ALREADY_EXISTS) {
		s.logger.Info().Err(err).Str("storageID", itemID.GetStorageId()).Msg("listPublicShares error: " + err.Error())
	}
	if err != nil && !strings.Contains(err.Error(), errtypes.ERR_ALREADY_EXISTS) {
		s.logger.Error().Err(err).Str("storageID", itemID.GetStorageId()).Msg("listPublicShares failed")
		return collectionOfPermissions, err
	}

	for _, driveItem := range driveItems {
		collectionOfPermissions.Value = append(collectionOfPermissions.Value, driveItem.Permissions...)
	}

	return collectionOfPermissions, nil
}

// ListSpaceRootPermissions handles ListPermissions request on project spaces
func (s DriveItemPermissionsService) ListSpaceRootPermissions(ctx context.Context, driveID *storageprovider.ResourceId) (libregraph.CollectionOfPermissionsWithAllowedValues, error) {

	ctx, span := s.tp.Tracer("graph").Start(ctx, "ListSpaceRootPermissions")
	defer span.End()

	collectionOfPermissions := libregraph.CollectionOfPermissionsWithAllowedValues{}
	gatewayClient, err := s.gatewaySelector.Next()
	if err != nil {
		s.logger.Error().Err(err).Str("storageID", driveID.GetStorageId()).Msg("gateway selection failed")
		return collectionOfPermissions, err
	}

	space, err := utils.GetSpace(ctx, storagespace.FormatResourceID(driveID), gatewayClient)
	if err != nil {
		err = errorcode.FromUtilsStatusCodeError(err)
		s.logger.Error().Err(err).Str("storageID", driveID.GetStorageId()).Interface("space.res", space).Msg("GetSpace failed")
		return collectionOfPermissions, err
	}

	isSupportedSpaceType := slices.Contains([]string{_spaceTypeProject, _spaceTypePersonal, _spaceTypeVirtual}, space.GetSpaceType())
	if !isSupportedSpaceType {
		err = errorcode.New(errorcode.InvalidRequest, "unsupported space type")
		s.logger.Error().Err(err).Str("storageID", driveID.GetStorageId()).Msg(err.Error())
		return collectionOfPermissions, err
	}

	rootResourceID := space.GetRoot()
	return s.ListPermissions(ctx, rootResourceID, false, false) // federated roles are not supported for spaces
}

// GetPermission returns a single permission of a drive item identified by permissionID
func (s DriveItemPermissionsService) GetPermission(ctx context.Context, itemID *storageprovider.ResourceId, permissionID string) (libregraph.Permission, error) {
	ctx, span := s.tp.Tracer("graph").Start(ctx, "GetPermission")
	defer span.End()

	// Reuse ListPermissions to obtain all permissions and select the requested one
	collection, err := s.ListPermissions(ctx, itemID, false, false)
	if err != nil {
		return libregraph.Permission{}, err
	}

	for _, p := range collection.Value {
		if p.GetId() == permissionID {
			return p, nil
		}
	}

	return libregraph.Permission{}, errorcode.New(errorcode.ItemNotFound, "permission not found")
}

// DeletePermission deletes a permission from a drive item
func (s DriveItemPermissionsService) DeletePermission(ctx context.Context, itemID *storageprovider.ResourceId, permissionID string) error {

	ctx, span := s.tp.Tracer("graph").Start(ctx, "DeletePermission")
	defer span.End()

	var permissionType permissionType

	sharedResourceID, err := s.getLinkPermissionResourceID(ctx, permissionID)
	switch {
	// Check if the ID is referring to a public share
	case err == nil:
		permissionType = Public
	// If the item id is referring to a space root and this is not a public share
	// we have to deal with space permissions
	case IsSpaceRoot(itemID):
		permissionType = Space
		sharedResourceID = itemID
		err = nil
	// If this is neither a public share nor a space permission, check if this is a
	// user share
	default:
		sharedResourceID, err = s.getUserPermissionResourceID(ctx, permissionID)
		if err == nil {
			permissionType = User
		}
	}

	if sharedResourceID == nil && s.config.IncludeOCMSharees {
		sharedResourceID, err = s.getOCMPermissionResourceID(ctx, permissionID)
		if err == nil {
			permissionType = OCM
		}
	}

	switch {
	case err != nil:
		s.logger.Error().Err(err).Msg("permission error")
		return err
	case permissionType == Unknown:
		s.logger.Error().Err(err).Msg("permission not found")
		return errorcode.New(errorcode.ItemNotFound, "permission not found")
	case sharedResourceID == nil:
		s.logger.Error().Err(err).Msg("failed to resolve resource id for shared resource")
		return errorcode.New(errorcode.ItemNotFound, "failed to resolve resource id for shared resource")
	}

	// The resourceID of the shared resource need to match the item ID from the Request Path
	// otherwise this is an invalid Request.
	if !utils.ResourceIDEqual(sharedResourceID, itemID) {
		err = errorcode.New(errorcode.InvalidRequest, "permissionID and itemID do not match")
		s.logger.Error().Err(err).Str("resourceId", itemID.GetStorageId()).Str("permissionID", permissionID).Msg(err.Error())
		return err
	}

	switch permissionType {
	case User:
		return s.removeUserShare(ctx, permissionID)
	case Public:
		return s.removePublicShare(ctx, permissionID)
	case Space:
		return s.removeSpacePermission(ctx, permissionID, sharedResourceID)
	case OCM:
		return s.removeOCMPermission(ctx, permissionID)
	default:
		// This should never be reached
		err = errorcode.New(errorcode.GeneralException, "failed to delete permission")
		s.logger.Error().Err(err).Str("resourceId", itemID.GetStorageId()).Str("permissionID", permissionID).Msg(err.Error())
		return err
	}
}

// DeleteSpaceRootPermission deletes a permission on the root item of a project space
func (s DriveItemPermissionsService) DeleteSpaceRootPermission(ctx context.Context, driveID *storageprovider.ResourceId, permissionID string) error {

	ctx, span := s.tp.Tracer("graph").Start(ctx, "DeleteSpaceRootPermission")
	defer span.End()

	gatewayClient, err := s.gatewaySelector.Next()
	if err != nil {
		s.logger.Error().Err(err).Str("permissionID", permissionID).Msg("selecting gatewaySelector failed")
		return err
	}

	space, err := utils.GetSpace(ctx, storagespace.FormatResourceID(driveID), gatewayClient)
	if err != nil {
		s.logger.Error().Err(err).Str("permissionID", permissionID).Interface("space.res", space).Msg("GetSpace failed")
		return errorcode.FromUtilsStatusCodeError(err)
	}

	if space.SpaceType != _spaceTypeProject {
		err = errorcode.New(errorcode.InvalidRequest, "unsupported space type")
		s.logger.Error().Err(err).Str("resourceId", driveID.GetStorageId()).Str("permissionID", permissionID).Msg(err.Error())
		return err
	}

	rootResourceID := space.GetRoot()
	return s.DeletePermission(ctx, rootResourceID, permissionID)
}

// UpdatePermission updates a permission on a drive item
func (s DriveItemPermissionsService) UpdatePermission(ctx context.Context, itemID *storageprovider.ResourceId, permissionID string, newPermission libregraph.Permission) (libregraph.Permission, error) {

	ctx, span := s.tp.Tracer("graph").Start(ctx, "UpdatePermission")
	defer span.End()

	oldPermission, sharedResourceID, err := s.getPermissionByID(ctx, permissionID, itemID)

	// try to get the permission from ocm if the permission was not found first place
	if err != nil && s.config.IncludeOCMSharees {
		oldPermission, sharedResourceID, err = s.getOCMPermissionByID(ctx, permissionID, itemID)
	}

	// if we still can't find the permission, return an error
	if err != nil {
		s.logger.Error().Err(err).Str("permissionID", permissionID).Msg("getPermissionByID failed")
		return libregraph.Permission{}, err
	}

	// The resourceID of the shared resource need to match the item ID from the Request Path
	// otherwise this is an invalid Request.
	if !utils.ResourceIDEqual(sharedResourceID, itemID) {
		err = errorcode.New(errorcode.InvalidRequest, "permissionID and itemID do not match")
		s.logger.Error().Str("permissionID", permissionID).Msg(err.Error())
		return libregraph.Permission{}, err
	}

	// This is a public link
	if _, ok := oldPermission.GetLinkOk(); ok {
		updatedPermission, err := s.updatePublicLinkPermission(ctx, permissionID, itemID, &newPermission)
		if err != nil {
			s.logger.Error().Err(err).Str("permissionID", permissionID).Msg("updatePublicLinkPermission failed")
			return libregraph.Permission{}, err
		}
		return *updatedPermission, nil
	}

	// This is a user share
	updatedPermission, err := s.updateUserShare(ctx, permissionID, sharedResourceID, &newPermission)
	if err == nil && updatedPermission != nil {
		return *updatedPermission, nil
	}

	// This is an ocm share
	if s.config.IncludeOCMSharees {
		updatePermission, err := s.updateOCMPermission(ctx, permissionID, itemID, &newPermission)
		if err == nil {
			return *updatePermission, nil
		}
	}

	s.logger.Error().Err(err).Str("permissionID", permissionID).Msg("UpdatePermission failed")
	return libregraph.Permission{}, err
}

// UpdateSpaceRootPermission updates a permission on the root item of a project space
func (s DriveItemPermissionsService) UpdateSpaceRootPermission(ctx context.Context, driveID *storageprovider.ResourceId, permissionID string, newPermission libregraph.Permission) (libregraph.Permission, error) {

	ctx, span := s.tp.Tracer("graph").Start(ctx, "UpdateSpaceRootPermission")
	defer span.End()

	gatewayClient, err := s.gatewaySelector.Next()
	if err != nil {
		s.logger.Error().Err(err).Str("permissionID", permissionID).Str("storageID", driveID.GetStorageId()).Msg("selecting gateway failed")
		return libregraph.Permission{}, err
	}

	space, err := utils.GetSpace(ctx, storagespace.FormatResourceID(driveID), gatewayClient)
	if err != nil {
		err = errorcode.FromUtilsStatusCodeError(err)
		s.logger.Error().Err(err).Str("permissionID", permissionID).Str("storageID", driveID.GetStorageId()).Msg("GetSpace failed")
		return libregraph.Permission{}, err
	}

	if space.SpaceType != _spaceTypeProject {
		err = errorcode.New(errorcode.InvalidRequest, "unsupported space type")
		s.logger.Error().Err(err).Str("permissionID", permissionID).Str("storageID", driveID.GetStorageId()).Msg(err.Error())
		return libregraph.Permission{}, err
	}

	rootResourceID := space.GetRoot()
	return s.UpdatePermission(ctx, rootResourceID, permissionID, newPermission)
}

// GetSpaceRootPermission handles requests to fetch a single permission on a space root
func (s DriveItemPermissionsService) GetSpaceRootPermission(ctx context.Context, driveID *storageprovider.ResourceId, permissionID string) (libregraph.Permission, error) {

	ctx, span := s.tp.Tracer("graph").Start(ctx, "GetSpaceRootPermission")
	defer span.End()

	gatewayClient, err := s.gatewaySelector.Next()
	if err != nil {
		s.logger.Error().Err(err).Str("storageID", driveID.GetStorageId()).Msg("gateway selection failed")
		return libregraph.Permission{}, err
	}

	space, err := utils.GetSpace(ctx, storagespace.FormatResourceID(driveID), gatewayClient)
	if err != nil {
		s.logger.Error().Err(err).Str("storageID", driveID.GetStorageId()).Msg("GetSpace failed")
		return libregraph.Permission{}, errorcode.FromUtilsStatusCodeError(err)
	}

	if space.SpaceType != _spaceTypeProject {
		err = errorcode.New(errorcode.InvalidRequest, "unsupported space type")
		s.logger.Error().Err(err).Str("storageID", driveID.GetStorageId()).Msg(err.Error())
		return libregraph.Permission{}, err
	}

	rootResourceID := space.GetRoot()
	return s.GetPermission(ctx, rootResourceID, permissionID)
}

// DriveItemPermissionsApi is the api that registers the http endpoints which expose needed operation to the graph api.
// the business logic is delegated to the permissions service and further down to the cs3 client.
type DriveItemPermissionsApi struct {
	logger                      log.Logger
	driveItemPermissionsService DriveItemPermissionsProvider
	config                      *config.Config
	valueService                settingssvc.ValueService
}

// NewDriveItemPermissionsApi creates a new DriveItemPermissionsApi
func NewDriveItemPermissionsApi(driveItemPermissionService DriveItemPermissionsProvider, logger log.Logger, c *config.Config, valueService settingssvc.ValueService) (DriveItemPermissionsApi, error) {
	return DriveItemPermissionsApi{
		logger:                      log.Logger{Logger: logger.With().Str("graph api", "DrivesDriveItemApi").Logger()},
		driveItemPermissionsService: driveItemPermissionService,
		config:                      c,
		valueService:                valueService,
	}, nil
}

// Invite handles DriveItemInvite requests
func (api DriveItemPermissionsApi) Invite(w http.ResponseWriter, r *http.Request) {
	_, itemID, err := GetDriveAndItemIDParam(r, &api.logger)
	if err != nil {
		api.logger.Error().Err(err).Msg(invalidIdMsg)
		errorcode.InvalidRequest.Render(w, r, http.StatusUnprocessableEntity, invalidIdMsg)
		return
	}

	driveItemInvite := &libregraph.DriveItemInvite{}
	if err = StrictJSONUnmarshal(r.Body, driveItemInvite); err != nil {
		api.logger.Error().Err(err).Interface("Body", r.Body).Msg("failed unmarshalling request body")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "invalid request body")
		return
	}

	ctx := validate.ContextWithAllowedRoleIDs(r.Context(), api.config.UnifiedRoles.AvailableRoles)
	if err = validate.StructCtx(ctx, driveItemInvite); err != nil {
		api.logger.Error().Err(err).Interface("Body", r.Body).Msg("invalid request body")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, err.Error())
		return
	}

	permission, err := api.driveItemPermissionsService.Invite(ctx, itemID, *driveItemInvite)
	if err != nil {
		errorcode.RenderError(w, r, err)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, &ListResponse{Value: []interface{}{permission}})
}

// SpaceRootInvite handles DriveItemInvite requests on a space root
func (api DriveItemPermissionsApi) SpaceRootInvite(w http.ResponseWriter, r *http.Request) {
	driveID, err := parseIDParam(r, "driveID")
	if err != nil {
		api.logger.Error().Err(err).Msg(parseDriveIDErrMsg)
		errorcode.InvalidRequest.Render(w, r, http.StatusUnprocessableEntity, parseDriveIDErrMsg)
		return
	}

	driveItemInvite := &libregraph.DriveItemInvite{}
	if err = StrictJSONUnmarshal(r.Body, driveItemInvite); err != nil {
		api.logger.Error().Err(err).Interface("Body", r.Body).Msg("failed unmarshalling request body")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "invalid request body")
		return
	}

	ctx := validate.ContextWithAllowedRoleIDs(r.Context(), api.config.UnifiedRoles.AvailableRoles)
	if err = validate.StructCtx(ctx, driveItemInvite); err != nil {
		api.logger.Error().Err(err).Interface("Body", r.Body).Msg("invalid request body")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, err.Error())
		return
	}
	permission, err := api.driveItemPermissionsService.SpaceRootInvite(ctx, &driveID, *driveItemInvite)

	if err != nil {
		errorcode.RenderError(w, r, err)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, &ListResponse{Value: []interface{}{permission}})
}

// ListPermissions handles ListPermissions requests
func (api DriveItemPermissionsApi) ListPermissions(w http.ResponseWriter, r *http.Request) {
	_, itemID, err := GetDriveAndItemIDParam(r, &api.logger)
	if err != nil {
		api.logger.Error().Err(err).Msg(invalidIdMsg)
		errorcode.RenderError(w, r, err)
		return
	}

	var listFederatedRoles bool
	if GetFilterParam(r) == "@libre.graph.permissions.roles.allowedValues/rolePermissions/any(p:contains(p/condition, '@Subject.UserType==\"Federated\"'))" {
		listFederatedRoles = true
	}

	var selectRoles bool
	if GetSelectParam(r) == "@libre.graph.permissions.roles.allowedValues" {
		selectRoles = true
	}

	ctx := r.Context()

	permissions, err := api.driveItemPermissionsService.ListPermissions(ctx, itemID, listFederatedRoles, selectRoles)
	if err != nil {
		errorcode.RenderError(w, r, err)
		return
	}

	userID := revactx.ContextMustGetUser(ctx).GetId().GetOpaqueId()
	loc := l10n.MustGetUserLocale(ctx, userID, r.Header.Get(l10n.HeaderAcceptLanguage), api.valueService)
	w.Header().Add("Content-Language", loc)
	if loc != "" && loc != "en" {
		err := l10n_pkg.TranslateEntity(loc, "en", permissions,
			l10n.TranslateEach("LibreGraphPermissionsRolesAllowedValues",
				l10n.TranslateField("Description"),
				l10n.TranslateField("DisplayName"),
			),
		)
		if err != nil {
			api.logger.Error().Err(err).Msg("tranlation error")
		}
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, permissions)
}

// ListSpaceRootPermissions handles ListPermissions requests on a space root
func (api DriveItemPermissionsApi) ListSpaceRootPermissions(w http.ResponseWriter, r *http.Request) {
	driveID, err := parseIDParam(r, "driveID")
	if err != nil {
		api.logger.Error().Err(err).Msg(parseDriveIDErrMsg)
		errorcode.InvalidRequest.Render(w, r, http.StatusUnprocessableEntity, parseDriveIDErrMsg)
		return
	}

	ctx := r.Context()
	permissions, err := api.driveItemPermissionsService.ListSpaceRootPermissions(ctx, &driveID)

	if err != nil {
		errorcode.RenderError(w, r, err)
		return
	}

	userID := revactx.ContextMustGetUser(ctx).GetId().GetOpaqueId()
	loc := l10n.MustGetUserLocale(ctx, userID, r.Header.Get(l10n.HeaderAcceptLanguage), api.valueService)
	w.Header().Add("Content-Language", loc)
	if loc != "" && loc != "en" {
		err := l10n_pkg.TranslateEntity(loc, "en", permissions,
			l10n.TranslateEach("LibreGraphPermissionsRolesAllowedValues",
				l10n.TranslateField("Description"),
				l10n.TranslateField("DisplayName"),
			),
		)
		if err != nil {
			api.logger.Error().Err(err).Msg("tranlation error")
		}
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, permissions)
}

// GetSpaceRootPermission handles requests to fetch a single permission on a space root
func (api DriveItemPermissionsApi) GetSpaceRootPermission(w http.ResponseWriter, r *http.Request) {
	driveID, err := parseIDParam(r, "driveID")
	if err != nil {
		api.logger.Error().Err(err).Msg(parseDriveIDErrMsg)
		errorcode.InvalidRequest.Render(w, r, http.StatusUnprocessableEntity, parseDriveIDErrMsg)
		return
	}

	permissionID, err := url.PathUnescape(chi.URLParam(r, "permissionID"))
	if err != nil {
		api.logger.Error().Err(err).Msg("could not parse permissionID")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "invalid permissionID")
		return
	}

	if permissionID == "" {
		api.logger.Error().Msg("permissionID cannot be empty")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "invalid permissionID")
		return
	}

	ctx := r.Context()
	permission, err := api.driveItemPermissionsService.GetSpaceRootPermission(ctx, &driveID, permissionID)
	if err != nil {
		errorcode.RenderError(w, r, err)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, &permission)
}

// GetPermission handles requests to fetch a single permission on a drive item
func (api DriveItemPermissionsApi) GetPermission(w http.ResponseWriter, r *http.Request) {
	_, itemID, err := GetDriveAndItemIDParam(r, &api.logger)
	if err != nil {
		api.logger.Error().Err(err).Msg(invalidIdMsg)
		errorcode.RenderError(w, r, err)
		return
	}

	permissionID, err := url.PathUnescape(chi.URLParam(r, "permissionID"))
	if err != nil {
		api.logger.Error().Err(err).Msg("could not parse permissionID")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "invalid permissionID")
		return
	}

	if permissionID == "" {
		api.logger.Error().Msg("permissionID cannot be empty")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "invalid permissionID")
		return
	}

	ctx := r.Context()
	permission, err := api.driveItemPermissionsService.GetPermission(ctx, itemID, permissionID)
	if err != nil {
		errorcode.RenderError(w, r, err)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, &permission)
}

// DeletePermission handles DeletePermission requests
func (api DriveItemPermissionsApi) DeletePermission(w http.ResponseWriter, r *http.Request) {
	_, itemID, err := GetDriveAndItemIDParam(r, &api.logger)
	if err != nil {
		api.logger.Error().Err(err).Msg(invalidIdMsg)
		errorcode.RenderError(w, r, err)
		return
	}

	permissionID, err := url.PathUnescape(chi.URLParam(r, "permissionID"))
	if err != nil {
		api.logger.Error().Err(err).Msg("could not parse permissionID")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "invalid permissionID")
		return
	}

	ctx := r.Context()
	err = api.driveItemPermissionsService.DeletePermission(ctx, itemID, permissionID)
	if err != nil {
		errorcode.RenderError(w, r, err)
		return
	}

	render.Status(r, http.StatusNoContent)
	render.NoContent(w, r)
}

// DeleteSpaceRootPermission handles DeletePermission requests on a space root
func (api DriveItemPermissionsApi) DeleteSpaceRootPermission(w http.ResponseWriter, r *http.Request) {
	driveID, err := parseIDParam(r, "driveID")
	if err != nil {
		api.logger.Error().Err(err).Msg(parseDriveIDErrMsg)
		errorcode.InvalidRequest.Render(w, r, http.StatusUnprocessableEntity, parseDriveIDErrMsg)
		return
	}

	permissionID, err := url.PathUnescape(chi.URLParam(r, "permissionID"))
	if err != nil {
		api.logger.Error().Err(err).Msg("could not parse permissionID")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "invalid permissionID")
		return
	}

	ctx := r.Context()
	err = api.driveItemPermissionsService.DeleteSpaceRootPermission(ctx, &driveID, permissionID)
	if err != nil {
		errorcode.RenderError(w, r, err)
		return
	}

	render.Status(r, http.StatusNoContent)
	render.NoContent(w, r)
}

// UpdatePermission handles UpdatePermission requests
func (api DriveItemPermissionsApi) UpdatePermission(w http.ResponseWriter, r *http.Request) {
	_, itemID, err := GetDriveAndItemIDParam(r, &api.logger)
	if err != nil {
		api.logger.Error().Err(err).Msg(invalidIdMsg)
		errorcode.RenderError(w, r, err)
		return
	}

	permissionID, err := url.PathUnescape(chi.URLParam(r, "permissionID"))
	if err != nil {
		api.logger.Error().Err(err).Msg("could not parse permissionID")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "invalid permissionID")
		return
	}

	permission := libregraph.Permission{}
	if err = StrictJSONUnmarshal(r.Body, &permission); err != nil {
		api.logger.Error().Err(err).Interface("Body", r.Body).Msg("failed unmarshalling request body")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "invalid request body")
		return
	}

	ctx := r.Context()
	if err = validate.StructCtx(ctx, permission); err != nil {
		api.logger.Error().Err(err).Interface("Body", r.Body).Msg("invalid request body")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, err.Error())
		return
	}

	updatedPermission, err := api.driveItemPermissionsService.UpdatePermission(ctx, itemID, permissionID, permission)
	if err != nil {
		errorcode.RenderError(w, r, err)
		return
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, &updatedPermission)
}

// UpdateSpaceRootPermission handles UpdatePermission requests on a space root
func (api DriveItemPermissionsApi) UpdateSpaceRootPermission(w http.ResponseWriter, r *http.Request) {
	driveID, err := parseIDParam(r, "driveID")
	if err != nil {
		api.logger.Error().Err(err).Msg(parseDriveIDErrMsg)
		errorcode.InvalidRequest.Render(w, r, http.StatusUnprocessableEntity, parseDriveIDErrMsg)
		return
	}

	permissionID, err := url.PathUnescape(chi.URLParam(r, "permissionID"))
	if err != nil {
		api.logger.Error().Err(err).Msg("could not parse permissionID")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "invalid permissionID")
		return
	}

	permission := libregraph.Permission{}
	if err = StrictJSONUnmarshal(r.Body, &permission); err != nil {
		api.logger.Error().Err(err).Interface("Body", r.Body).Msg("failed unmarshalling request body")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "invalid request body")
		return
	}

	ctx := r.Context()
	if err = validate.StructCtx(ctx, permission); err != nil {
		api.logger.Error().Err(err).Interface("Body", r.Body).Msg("invalid request body")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, err.Error())
		return
	}

	updatedPermission, err := api.driveItemPermissionsService.UpdateSpaceRootPermission(ctx, &driveID, permissionID, permission)
	if err != nil {
		errorcode.RenderError(w, r, err)
		return
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, &updatedPermission)
}
