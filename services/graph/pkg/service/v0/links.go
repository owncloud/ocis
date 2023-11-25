package svc

import (
	"context"
	"net/http"
	"net/url"
	"path"
	"strconv"

	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	link "github.com/cs3org/go-cs3apis/cs3/sharing/link/v1beta1"
	providerv1beta1 "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/cs3org/reva/v2/pkg/utils"
	"github.com/go-chi/render"
	libregraph "github.com/owncloud/libre-graph-api-go"
	"github.com/owncloud/ocis/v2/services/graph/pkg/linktype"
	"github.com/owncloud/ocis/v2/services/graph/pkg/service/v0/errorcode"
)

func (g Graph) CreateLink(w http.ResponseWriter, r *http.Request) {
	logger := g.logger.SubloggerWithRequestID(r.Context())
	logger.Info().Msg("calling create link")
	driveID, err := parseIDParam(r, "driveID")
	if err != nil {
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, err.Error())
		return
	}
	driveItemID, err := parseIDParam(r, "itemID")
	if err != nil {
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, err.Error())
		return
	}
	if driveID.StorageId != driveItemID.StorageId || driveID.SpaceId != driveItemID.SpaceId {
		errorcode.ItemNotFound.Render(w, r, http.StatusNotFound, "Item does not exist")
		return
	}
	var createLink libregraph.DriveItemCreateLink
	if err := StrictJSONUnmarshal(r.Body, &createLink); err != nil {
		logger.Error().Err(err).Interface("body", r.Body).Msg("could not create link: invalid body schema definition")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "invalid body schema definition")
		return
	}

	createdLink, err := g.createLink(r.Context(), &driveItemID, createLink)
	if err != nil {
		errorcode.RenderError(w, r, err)
		return
	}

	perm, err := g.libreGraphPermissionFromCS3PublicShare(createdLink)
	if err != nil {
		errorcode.GeneralException.Render(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, []libregraph.Permission{*perm})
}

func (g Graph) createLink(ctx context.Context, driveItemID *providerv1beta1.ResourceId, createLink libregraph.DriveItemCreateLink) (*link.PublicShare, error) {
	gatewayClient, err := g.gatewaySelector.Next()
	if err != nil {
		g.logger.Error().Err(err).Msg("could not select next gateway client")
		return nil, errorcode.New(errorcode.GeneralException, err.Error())
	}

	statResp, err := gatewayClient.Stat(
		ctx,
		&providerv1beta1.StatRequest{
			Ref: &providerv1beta1.Reference{
				ResourceId: driveItemID,
				Path:       ".",
			},
		})
	if err != nil {
		g.logger.Error().Err(err).Msg("transport error, could not stat resource")
		return nil, errorcode.New(errorcode.GeneralException, err.Error())
	}
	if code := statResp.GetStatus().GetCode(); code != rpc.Code_CODE_OK {
		g.logger.Debug().Interface("itemID", driveItemID).Msg(statResp.GetStatus().GetMessage())
		return nil, errorcode.New(cs3StatusToErrCode(code), statResp.GetStatus().GetMessage())
	}
	permissions, err := linktype.CS3ResourcePermissionsFromSharingLink(createLink, statResp.GetInfo().GetType())
	if err != nil {
		g.logger.Debug().Interface("createLink", createLink).Msg(err.Error())
		return nil, errorcode.New(errorcode.InvalidRequest, "invalid link type")
	}
	req := link.CreatePublicShareRequest{
		ResourceInfo: statResp.GetInfo(),
		Grant: &link.Grant{
			Permissions: &link.PublicSharePermissions{
				Permissions: permissions,
			},
			Password: createLink.GetPassword(),
		},
	}
	if createLink.ExpirationDateTime != nil {
		req.GetGrant().Expiration = utils.TimeToTS(createLink.GetExpirationDateTime())
	}

	// set displayname and password protected as arbitrary metadata
	req.ResourceInfo.ArbitraryMetadata = &providerv1beta1.ArbitraryMetadata{
		Metadata: map[string]string{
			"name":      createLink.GetDisplayName(),
			"quicklink": strconv.FormatBool(createLink.GetLibreGraphQuickLink()),
		},
	}
	createResp, err := gatewayClient.CreatePublicShare(ctx, &req)
	if err != nil {
		g.logger.Error().Err(err).Msg("transport error, could not create link")
		return nil, errorcode.New(errorcode.GeneralException, err.Error())
	}
	if statusCode := createResp.GetStatus().GetCode(); statusCode != rpc.Code_CODE_OK {
		return nil, errorcode.New(cs3StatusToErrCode(statusCode), createResp.Status.Message)
	}
	return createResp.GetShare(), nil
}

func (g Graph) libreGraphPermissionFromCS3PublicShare(createdLink *link.PublicShare) (*libregraph.Permission, error) {
	webURL, err := url.Parse(g.config.Spaces.WebDavBase)
	if err != nil {
		g.logger.Error().
			Err(err).
			Str("url", g.config.Spaces.WebDavBase).
			Msg("failed to parse webURL base url")
		return nil, err
	}
	lt, actions := linktype.SharingLinkTypeFromCS3Permissions(createdLink.GetPermissions())
	perm := libregraph.NewPermission()
	perm.Id = libregraph.PtrString(createdLink.GetId().GetOpaqueId())
	perm.Link = &libregraph.SharingLink{
		Type:                  lt,
		PreventsDownload:      libregraph.PtrBool(false),
		LibreGraphDisplayName: libregraph.PtrString(createdLink.GetDisplayName()),
		LibreGraphQuickLink:   libregraph.PtrBool(createdLink.GetQuicklink()),
	}
	perm.LibreGraphPermissionsActions = actions
	webURL.Path = path.Join(webURL.Path, "s", createdLink.GetToken())
	perm.Link.SetWebUrl(webURL.String())

	// set expiration date
	if createdLink.GetExpiration() != nil {
		perm.SetExpirationDateTime(cs3TimestampToTime(createdLink.GetExpiration()))
	}
	return perm, nil
}
