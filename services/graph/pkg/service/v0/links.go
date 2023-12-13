package svc

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"time"

	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	link "github.com/cs3org/go-cs3apis/cs3/sharing/link/v1beta1"
	providerv1beta1 "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	types "github.com/cs3org/go-cs3apis/cs3/types/v1beta1"
	"github.com/cs3org/reva/v2/pkg/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	libregraph "github.com/owncloud/libre-graph-api-go"
	"github.com/owncloud/ocis/v2/services/graph/pkg/errorcode"
	"github.com/owncloud/ocis/v2/services/graph/pkg/linktype"
)

// CreateLink creates a public link on the cs3 api
func (g Graph) CreateLink(w http.ResponseWriter, r *http.Request) {
	logger := g.logger.SubloggerWithRequestID(r.Context())
	logger.Info().Msg("calling create link")

	_, driveItemID, err := g.GetDriveAndItemIDParam(r)
	if err != nil {
		errorcode.RenderError(w, r, err)
		return
	}

	var createLink libregraph.DriveItemCreateLink
	body, err := io.ReadAll(r.Body)
	if err := json.Unmarshal(body, &createLink); err != nil {
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

	render.Status(r, http.StatusOK)
	render.JSON(w, r, *perm)
}

// SetLinkPassword sets public link password on the cs3 api
func (g Graph) SetLinkPassword(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	_, itemID, err := g.GetDriveAndItemIDParam(r)
	if err != nil {
		errorcode.RenderError(w, r, err)
		return
	}

	permissionID, err := url.PathUnescape(chi.URLParam(r, "permissionID"))
	if err != nil {
		g.logger.Debug().Err(err).Msg("could not parse permissionID")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "invalid permissionID")
		return
	}

	password := &libregraph.SharingLinkPassword{}
	if err := StrictJSONUnmarshal(r.Body, password); err != nil {
		g.logger.Debug().Err(err).Interface("Body", r.Body).Msg("failed unmarshalling request body")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "invalid request body")
		return
	}

	publicShare, err := g.getCS3PublicShareByID(ctx, permissionID)
	if err != nil {
		errorcode.RenderError(w, r, err)
		return
	}

	// The resourceID of the shared resource need to match the item ID from the Request Path
	// otherwise this is an invalid Request.
	if !utils.ResourceIDEqual(publicShare.GetResourceId(), &itemID) {
		g.logger.Debug().Msg("resourceID of shared does not match itemID")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "permissionID and itemID do not match")
		return
	}

	newPermission, err := g.updatePublicLinkPassword(ctx, permissionID, password.GetPassword())
	if err != nil {
		errorcode.RenderError(w, r, err)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, *newPermission)
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
	expirationDate, isSet := createLink.GetExpirationDateTimeOk()
	if isSet {
		expireTime := parseAndFillUpTime(expirationDate)
		if expireTime == nil {
			g.logger.Debug().Interface("createLink", createLink).Msg(err.Error())
			return nil, errorcode.New(errorcode.InvalidRequest, "invalid expiration date")
		}
		req.GetGrant().Expiration = expireTime
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
		perm.SetExpirationDateTime(cs3TimestampToTime(createdLink.GetExpiration()).UTC())
	}

	perm.SetHasPassword(createdLink.GetPasswordProtected())

	return perm, nil
}

func parseAndFillUpTime(t *time.Time) *types.Timestamp {
	if t == nil {
		return nil
	}

	// the link needs to be valid for the whole day
	tLink := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	tLink = tLink.Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	final := tLink.UnixNano()

	return &types.Timestamp{
		Seconds: uint64(final / 1000000000),
		Nanos:   uint32(final % 1000000000),
	}
}

func (g Graph) updatePublicLinkPassword(ctx context.Context, permissionID string, password string) (*libregraph.Permission, error) {
	gatewayClient, err := g.gatewaySelector.Next()
	if err != nil {
		return nil, err
	}

	changeLinkRes, err := gatewayClient.UpdatePublicShare(ctx, &link.UpdatePublicShareRequest{
		Update: &link.UpdatePublicShareRequest_Update{
			Type: link.UpdatePublicShareRequest_Update_TYPE_PASSWORD,
			Grant: &link.Grant{
				Password: password,
			},
		},
		Ref: &link.PublicShareReference{
			Spec: &link.PublicShareReference_Id{
				Id: &link.PublicShareId{
					OpaqueId: permissionID,
				},
			},
		},
	})
	if errCode := errorcode.FromCS3Status(changeLinkRes.GetStatus(), err); errCode != nil {
		return nil, *errCode
	}
	permission, err := g.libreGraphPermissionFromCS3PublicShare(changeLinkRes.GetShare())
	if err != nil {
		return nil, err
	}
	return permission, nil
}

func (g Graph) updatePublicLinkPermission(ctx context.Context, permissionID string, itemID *providerv1beta1.ResourceId, newPermission *libregraph.Permission) (perm *libregraph.Permission, err error) {
	gatewayClient, err := g.gatewaySelector.Next()
	if err != nil {
		g.logger.Error().Err(err).Msg("could not select next gateway client")
		return nil, errorcode.New(errorcode.GeneralException, err.Error())
	}

	statResp, err := gatewayClient.Stat(
		ctx,
		&providerv1beta1.StatRequest{
			Ref: &providerv1beta1.Reference{
				ResourceId: itemID,
				Path:       ".",
			},
		})

	if errCode := errorcode.FromCS3Status(statResp.GetStatus(), err); errCode != nil {
		return nil, *errCode
	}

	if newPermission.HasExpirationDateTime() {
		expirationDate := newPermission.GetExpirationDateTime()
		update := &link.UpdatePublicShareRequest_Update{
			Type:  link.UpdatePublicShareRequest_Update_TYPE_EXPIRATION,
			Grant: &link.Grant{Expiration: parseAndFillUpTime(&expirationDate)},
		}
		perm, err = g.updatePublicLink(ctx, permissionID, update)
		if err != nil {
			return nil, err
		}
	}

	if newPermission.HasLink() && newPermission.Link.HasLibreGraphDisplayName() {
		changedLink := newPermission.GetLink()
		update := &link.UpdatePublicShareRequest_Update{
			Type:        link.UpdatePublicShareRequest_Update_TYPE_DISPLAYNAME,
			DisplayName: changedLink.GetLibreGraphDisplayName(),
		}
		perm, err = g.updatePublicLink(ctx, permissionID, update)
		if err != nil {
			return nil, err
		}
	}

	if newPermission.HasLink() && newPermission.Link.HasType() {
		changedLink := newPermission.Link.GetType()
		permissions, err := linktype.CS3ResourcePermissionsFromSharingLink(
			libregraph.DriveItemCreateLink{
				Type: &changedLink,
			},
			statResp.GetInfo().GetType(),
		)
		update := &link.UpdatePublicShareRequest_Update{
			Type: link.UpdatePublicShareRequest_Update_TYPE_PERMISSIONS,
			Grant: &link.Grant{
				Permissions: &link.PublicSharePermissions{Permissions: permissions},
			},
		}
		perm, err = g.updatePublicLink(ctx, permissionID, update)
		if err != nil {
			return nil, err
		}
	}

	return perm, err
}

func (g Graph) updatePublicLink(ctx context.Context, permissionID string, update *link.UpdatePublicShareRequest_Update) (*libregraph.Permission, error) {
	gatewayClient, err := g.gatewaySelector.Next()
	if err != nil {
		return nil, err
	}

	changeLinkRes, err := gatewayClient.UpdatePublicShare(ctx, &link.UpdatePublicShareRequest{
		Update: update,
		Ref: &link.PublicShareReference{
			Spec: &link.PublicShareReference_Id{
				Id: &link.PublicShareId{
					OpaqueId: permissionID,
				},
			},
		},
	})

	if errCode := errorcode.FromCS3Status(changeLinkRes.GetStatus(), err); errCode != nil {
		return nil, *errCode
	}

	permission, err := g.libreGraphPermissionFromCS3PublicShare(changeLinkRes.GetShare())
	if err != nil {
		return nil, err
	}
	return permission, nil
}
