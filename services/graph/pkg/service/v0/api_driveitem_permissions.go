package svc

import (
	"context"
	"net/http"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	grouppb "github.com/cs3org/go-cs3apis/cs3/identity/group/v1beta1"
	userpb "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	collaboration "github.com/cs3org/go-cs3apis/cs3/sharing/collaboration/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	storageprovider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/cs3org/reva/v2/pkg/rgrpc/todo/pool"
	"github.com/cs3org/reva/v2/pkg/utils"
	"github.com/go-chi/render"
	libregraph "github.com/owncloud/libre-graph-api-go"
	"github.com/owncloud/ocis/v2/ocis-pkg/conversions"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/services/graph/pkg/errorcode"
	"github.com/owncloud/ocis/v2/services/graph/pkg/identity"
	"github.com/owncloud/ocis/v2/services/graph/pkg/unifiedrole"
	"github.com/owncloud/ocis/v2/services/graph/pkg/validate"
)

type DriveItemPermissionsProvider interface {
	Invite(ctx context.Context, resourceId provider.ResourceId, invite libregraph.DriveItemInvite) (libregraph.Permission, error)
}

// DriveItemPermissionsService contains the production business logic for everything that relates to permissions on drive items.
type DriveItemPermissionsService struct {
	logger           log.Logger
	gatewaySelector  pool.Selectable[gateway.GatewayAPIClient]
	identityCache    identity.IdentityCache
	resharingEnabled bool
}

// NewDriveItemPermissionsService creates a new DriveItemPermissionsService
func NewDriveItemPermissionsService(logger log.Logger, gatewaySelector pool.Selectable[gateway.GatewayAPIClient], identityCache identity.IdentityCache, resharing bool) (DriveItemPermissionsService, error) {
	return DriveItemPermissionsService{
		logger:           log.Logger{Logger: logger.With().Str("graph api", "DrivesDriveItemService").Logger()},
		gatewaySelector:  gatewaySelector,
		identityCache:    identityCache,
		resharingEnabled: resharing,
	}, nil
}

// Invite invites a user to a drive item.
func (s DriveItemPermissionsService) Invite(ctx context.Context, resourceId provider.ResourceId, invite libregraph.DriveItemInvite) (libregraph.Permission, error) {
	gatewayClient, err := s.gatewaySelector.Next()
	if err != nil {
		return libregraph.Permission{}, err
	}

	statResponse, err := gatewayClient.Stat(ctx, &storageprovider.StatRequest{Ref: &storageprovider.Reference{ResourceId: &resourceId}})
	if errCode := errorcode.FromStat(statResponse, err); errCode != nil {
		s.logger.Warn().Err(errCode).Interface("stat.res", statResponse).Msg("stat failed")
		return libregraph.Permission{}, *errCode
	}

	resourceInfo := statResponse.GetInfo()
	condition := unifiedrole.UnifiedRoleConditionGrantee
	if IsSpaceRoot(resourceInfo.GetId()) {
		condition = unifiedrole.UnifiedRoleConditionOwner
	}

	unifiedRolePermissions := []*libregraph.UnifiedRolePermission{{AllowedResourceActions: invite.LibreGraphPermissionsActions}}
	for _, roleID := range invite.GetRoles() {
		role, err := unifiedrole.NewUnifiedRoleFromID(roleID, s.resharingEnabled)
		if err != nil {
			s.logger.Debug().Err(err).Interface("role", invite.GetRoles()[0]).Msg("unable to convert requested role")
			return libregraph.Permission{}, err
		}

		allowedResourceActions := unifiedrole.GetAllowedResourceActions(role, condition)
		if len(allowedResourceActions) == 0 {
			return libregraph.Permission{}, errorcode.New(errorcode.InvalidRequest, "role not applicable to this resource")
		}

		unifiedRolePermissions = append(unifiedRolePermissions, conversions.ToPointerSlice(role.GetRolePermissions())...)
	}

	driveRecipient := invite.GetRecipients()[0]

	objectID := driveRecipient.GetObjectId()
	cs3ResourcePermissions := unifiedrole.PermissionsToCS3ResourcePermissions(unifiedRolePermissions)

	createShareRequest := &collaboration.CreateShareRequest{
		ResourceInfo: resourceInfo,
		Grant: &collaboration.ShareGrant{
			Permissions: &collaboration.SharePermissions{
				Permissions: cs3ResourcePermissions,
			},
		},
	}

	permission := &libregraph.Permission{}
	if role := unifiedrole.CS3ResourcePermissionsToUnifiedRole(*cs3ResourcePermissions, condition, s.resharingEnabled); role != nil {
		permission.Roles = []string{role.GetId()}
	}

	if len(permission.GetRoles()) == 0 {
		permission.LibreGraphPermissionsActions = unifiedrole.CS3ResourcePermissionsToLibregraphActions(*cs3ResourcePermissions)
	}

	switch driveRecipient.GetLibreGraphRecipientType() {
	case "group":
		group, err := s.identityCache.GetGroup(ctx, objectID)
		if err != nil {
			s.logger.Debug().Err(err).Interface("groupId", objectID).Msg("failed group lookup")
			return libregraph.Permission{}, errorcode.New(errorcode.InvalidRequest, err.Error())
		}
		createShareRequest.GetGrant().Grantee = &storageprovider.Grantee{
			Type: storageprovider.GranteeType_GRANTEE_TYPE_GROUP,
			Id: &storageprovider.Grantee_GroupId{GroupId: &grouppb.GroupId{
				OpaqueId: group.GetId(),
			}},
		}
		permission.GrantedToV2 = &libregraph.SharePointIdentitySet{
			Group: &libregraph.Identity{
				DisplayName: group.GetDisplayName(),
				Id:          conversions.ToPointer(group.GetId()),
			},
		}
	default:
		user, err := s.identityCache.GetUser(ctx, objectID)
		if err != nil {
			s.logger.Debug().Err(err).Interface("userId", objectID).Msg("failed user lookup")
			return libregraph.Permission{}, errorcode.New(errorcode.InvalidRequest, err.Error())
		}

		createShareRequest.GetGrant().Grantee = &storageprovider.Grantee{
			Type: storageprovider.GranteeType_GRANTEE_TYPE_USER,
			Id: &storageprovider.Grantee_UserId{UserId: &userpb.UserId{
				OpaqueId: user.GetId(),
			}},
		}
		permission.GrantedToV2 = &libregraph.SharePointIdentitySet{
			User: &libregraph.Identity{
				DisplayName: user.GetDisplayName(),
				Id:          conversions.ToPointer(user.GetId()),
			},
		}
	}

	if invite.ExpirationDateTime != nil {
		createShareRequest.GetGrant().Expiration = utils.TimeToTS(*invite.ExpirationDateTime)
	}

	createShareResponse, err := gatewayClient.CreateShare(ctx, createShareRequest)
	if errCode := errorcode.FromCS3Status(createShareResponse.GetStatus(), err); errCode != nil {
		s.logger.Debug().Err(err).Msg("share creation failed")
		return libregraph.Permission{}, *errCode
	}

	if id := createShareResponse.GetShare().GetId().GetOpaqueId(); id != "" {
		permission.Id = conversions.ToPointer(id)
	} else if IsSpaceRoot(resourceInfo.GetId()) {
		// permissions on a space root are not handled by a share manager so
		// they don't get a share-id
		permission.SetId(identitySetToSpacePermissionID(permission.GetGrantedToV2()))
	}

	if expiration := createShareResponse.GetShare().GetExpiration(); expiration != nil {
		permission.SetExpirationDateTime(utils.TSToTime(expiration))
	}

	return *permission, nil
}

// DriveItemPermissionsService is the api that registers the http endpoints which expose needed operation to the graph api.
// the business logic is delegated to the permissions service and further down to the cs3 client.
type DriveItemPermissionsApi struct {
	logger                      log.Logger
	driveItemPermissionsService DriveItemPermissionsProvider
}

// NewDriveItemPermissionsApi creates a new DriveItemPermissionsApi
func NewDriveItemPermissionsApi(driveItemPermissionService DriveItemPermissionsProvider, logger log.Logger) (DriveItemPermissionsApi, error) {
	return DriveItemPermissionsApi{
		logger:                      log.Logger{Logger: logger.With().Str("graph api", "DrivesDriveItemApi").Logger()},
		driveItemPermissionsService: driveItemPermissionService,
	}, nil
}

func (api DriveItemPermissionsApi) Invite(w http.ResponseWriter, r *http.Request) {
	_, itemID, err := GetDriveAndItemIDParam(r, &api.logger)
	if err != nil {
		msg := "invalid driveID or itemID"
		api.logger.Debug().Err(err).Msg(msg)
		errorcode.InvalidRequest.Render(w, r, http.StatusUnprocessableEntity, msg)
		return
	}

	driveItemInvite := &libregraph.DriveItemInvite{}
	if err = StrictJSONUnmarshal(r.Body, driveItemInvite); err != nil {
		api.logger.Debug().Err(err).Interface("Body", r.Body).Msg("failed unmarshalling request body")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "invalid request body")
		return
	}

	ctx := r.Context()

	if err = validate.StructCtx(ctx, driveItemInvite); err != nil {
		api.logger.Debug().Err(err).Interface("Body", r.Body).Msg("invalid request body")
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
