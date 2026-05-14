package svc

import (
	"context"
	"encoding/xml"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"

	gatewayv1beta1 "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	rpcv1beta1 "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	grpcmetadata "google.golang.org/grpc/metadata"

	"github.com/owncloud/ocis/v2/services/webdav/pkg/constants"
	"github.com/owncloud/ocis/v2/services/webdav/pkg/net"
	"github.com/owncloud/ocis/v2/services/webdav/pkg/prop"
	"github.com/owncloud/ocis/v2/services/webdav/pkg/propfind"
	revactx "github.com/owncloud/reva/v2/pkg/ctx"
	"github.com/owncloud/reva/v2/pkg/conversions"
	"github.com/owncloud/reva/v2/pkg/permission"
	"github.com/owncloud/reva/v2/pkg/storagespace"
	"github.com/owncloud/reva/v2/pkg/utils"
)

const propOcFavorite = "http://owncloud.org/ns/favorite"

// favoriteInfo holds a ResourceInfo and the resolved path for href construction.
type favoriteInfo struct {
	info *provider.ResourceInfo
	// href-ready path relative to the space root, e.g. "Documents/notes.md"
	relativePath string
	// hrefPrefix is the DAV prefix for constructing hrefs, e.g.
	// "/dav/files/admin" for personal/share spaces or
	// "/dav/spaces/<storageId>$<spaceId>" for project spaces.
	hrefPrefix string
}

// handleFilterFiles handles REPORT requests with oc:filter-files / oc:filter-rules.
func (g Webdav) handleFilterFiles(w http.ResponseWriter, r *http.Request, ff *reportFilterFiles) {
	logger := g.log.SubloggerWithRequestID(r.Context())

	if !ff.Rules.Favorite {
		// Only favorites filtering is supported; return empty 207.
		g.sendFavoritesResponse(nil, w, r)
		return
	}

	t := r.Header.Get(revactx.TokenHeader)
	ctx := revactx.ContextSetToken(r.Context(), t)
	ctx = grpcmetadata.AppendToOutgoingContext(ctx, revactx.TokenHeader, t)

	gwClient, err := g.gatewaySelector.Next()
	if err != nil {
		logger.Error().Err(err).Msg("error selecting gateway client")
		renderError(w, r, errInternalError("could not get gateway client"))
		return
	}

	// Get current user — needed both for CheckPermission (which reads the
	// user from the context) and for href construction later.
	whoAmI, err := gwClient.WhoAmI(ctx, &gatewayv1beta1.WhoAmIRequest{Token: t})
	if err != nil {
		logger.Error().Err(err).Msg("error getting current user")
		renderError(w, r, errInternalError("could not get current user"))
		return
	}
	if whoAmI.Status.Code != rpcv1beta1.Code_CODE_OK {
		logger.Error().Str("status", whoAmI.Status.Message).Msg("could not get current user")
		renderError(w, r, errInternalError("could not get current user"))
		return
	}
	ctx = revactx.ContextSetUser(ctx, whoAmI.User)
	username := whoAmI.User.Username

	// Check permission
	ok, err := utils.CheckPermission(ctx, permission.ListFavorites, gwClient)
	if err != nil {
		logger.Error().Err(err).Msg("error checking list favorites permission")
		renderError(w, r, errInternalError("error checking permission"))
		return
	}
	if !ok {
		logger.Debug().Msg("user not allowed to list favorites")
		renderError(w, r, errPermissionDenied("permission denied"))
		return
	}

	// List user's storage spaces
	spacesResp, err := gwClient.ListStorageSpaces(ctx, &provider.ListStorageSpacesRequest{
		Filters: []*provider.ListStorageSpacesRequest_Filter{
			{
				Type: provider.ListStorageSpacesRequest_Filter_TYPE_USER,
				Term: &provider.ListStorageSpacesRequest_Filter_Owner{
					Owner: whoAmI.User.Id,
				},
			},
		},
	})
	if err != nil {
		logger.Error().Err(err).Msg("error listing storage spaces")
		renderError(w, r, errInternalError("could not list storage spaces"))
		return
	}
	if spacesResp.Status.Code != rpcv1beta1.Code_CODE_OK {
		logger.Error().Str("status", spacesResp.Status.Message).Msg("could not list storage spaces")
		renderError(w, r, errInternalError("could not list storage spaces"))
		return
	}

	// Build the /dav/files/<user> href prefix.
	// The frontend expects all hrefs under this prefix when it sends
	// REPORT to /dav/files/<user>. Project spaces are not addressable
	// under this path and are skipped for now.
	filesPrefix := path.Join("/dav/files", username)
	if strings.HasPrefix(r.URL.Path, "/remote.php/") {
		filesPrefix = path.Join("/remote.php/dav/files", username)
	}

	// Collect favorites across personal and share spaces
	var favorites []favoriteInfo
	for _, space := range spacesResp.StorageSpaces {
		if space.Root == nil {
			continue
		}

		var pathPrefix string
		switch space.SpaceType {
		case "personal":
			pathPrefix = ""
		case "mountpoint", "grant":
			// Mounted shares appear under "Shares/<name>"
			name := space.Name
			if name == "" {
				name = space.Id.OpaqueId
			}
			pathPrefix = path.Join("Shares", name)
		default:
			// Project spaces and other types don't appear under
			// /dav/files/<user>/ — skip for now.
			continue
		}

		g.collectFavorites(ctx, gwClient, &provider.Reference{ResourceId: space.Root}, pathPrefix, filesPrefix, &favorites)
	}

	g.sendFavoritesResponse(favorites, w, r)
}

// collectFavorites recursively walks a storage space, collecting resources
// that have the oc:favorite metadata set.
func (g Webdav) collectFavorites(
	ctx context.Context,
	client gatewayv1beta1.GatewayAPIClient,
	ref *provider.Reference,
	pathPrefix string,
	hrefPrefix string,
	results *[]favoriteInfo,
) {
	resp, err := client.ListContainer(ctx, &provider.ListContainerRequest{
		Ref:                   ref,
		ArbitraryMetadataKeys: []string{propOcFavorite},
	})
	if err != nil {
		g.log.Error().Err(err).Msg("error listing container for favorites")
		return
	}
	if resp.Status.Code != rpcv1beta1.Code_CODE_OK {
		// Skip spaces/directories we cannot access
		return
	}

	for _, info := range resp.Infos {
		childPath := path.Join(pathPrefix, info.GetName())

		// Check if this resource is favorited
		if md := info.GetArbitraryMetadata().GetMetadata(); md != nil {
			if fav, ok := md[propOcFavorite]; ok && fav != "" && fav != "0" {
				*results = append(*results, favoriteInfo{
					info:         info,
					relativePath: childPath,
					hrefPrefix:   hrefPrefix,
				})
			}
		}

		// Recurse into directories
		if info.Type == provider.ResourceType_RESOURCE_TYPE_CONTAINER {
			g.collectFavorites(ctx, client, &provider.Reference{ResourceId: info.Id}, childPath, hrefPrefix, results)
		}
	}
}

// sendFavoritesResponse writes a 207 Multi-Status response for the collected favorites.
func (g Webdav) sendFavoritesResponse(favorites []favoriteInfo, w http.ResponseWriter, r *http.Request) {
	logger := g.log.SubloggerWithRequestID(r.Context())

	responses := make([]*propfind.ResponseXML, 0, len(favorites))
	for i := range favorites {
		resp := favoriteInfoToPropResponse(&favorites[i])
		responses = append(responses, resp)
	}

	msr := propfind.NewMultiStatusResponseXML()
	msr.Responses = responses

	msg, err := xml.Marshal(msr)
	if err != nil {
		logger.Error().Err(err).Msg("error marshaling favorites response")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set(net.HeaderDav, "1, 3, extended-mkcol")
	w.Header().Set(net.HeaderContentType, "application/xml; charset=utf-8")
	w.WriteHeader(http.StatusMultiStatus)
	if _, err := w.Write(msg); err != nil {
		logger.Err(err).Msg("error writing favorites response")
	}
}

// favoriteInfoToPropResponse converts a favoriteInfo into a ResponseXML.
func favoriteInfoToPropResponse(fav *favoriteInfo) *propfind.ResponseXML {
	info := fav.info

	response := &propfind.ResponseXML{
		Href:     net.EncodePath(path.Join(fav.hrefPrefix, fav.relativePath)),
		Propstat: []propfind.PropstatXML{},
	}

	propstatOK := propfind.PropstatXML{
		Status: "HTTP/1.1 200 OK",
		Prop:   []prop.PropertyXML{},
	}

	// oc:fileid
	if info.Id != nil {
		propstatOK.Prop = append(propstatOK.Prop, prop.Escaped("oc:fileid", storagespace.FormatResourceID(info.Id)))
	}

	// oc:name
	propstatOK.Prop = append(propstatOK.Prop, prop.Escaped("oc:name", info.GetName()))

	// d:getlastmodified
	if info.Mtime != nil {
		t := time.Unix(int64(info.Mtime.Seconds), int64(info.Mtime.Nanos))
		propstatOK.Prop = append(propstatOK.Prop, prop.Escaped("d:getlastmodified", t.UTC().Format(constants.RFC1123)))
	}

	// d:getcontenttype
	propstatOK.Prop = append(propstatOK.Prop, prop.Escaped("d:getcontenttype", info.GetMimeType()))

	// d:getetag
	if info.Etag != "" {
		propstatOK.Prop = append(propstatOK.Prop, prop.Escaped("d:getetag", "\""+info.Etag+"\""))
	}

	// oc:permissions
	if info.PermissionSet != nil {
		role := conversions.RoleFromResourcePermissions(info.PermissionSet, false)
		isDir := info.Type == provider.ResourceType_RESOURCE_TYPE_CONTAINER
		wdp := role.WebDAVPermissions(isDir, false, false, false)
		propstatOK.Prop = append(propstatOK.Prop, prop.Escaped("oc:permissions", wdp))
	}

	// oc:favorite (always "1" since we only return favorites)
	propstatOK.Prop = append(propstatOK.Prop, prop.Escaped("oc:favorite", "1"))

	// d:resourcetype + size
	size := strconv.FormatUint(info.Size, 10)
	if info.Type == provider.ResourceType_RESOURCE_TYPE_CONTAINER {
		propstatOK.Prop = append(propstatOK.Prop, prop.Raw("d:resourcetype", "<d:collection/>"))
		propstatOK.Prop = append(propstatOK.Prop, prop.Escaped("oc:size", size))
	} else {
		propstatOK.Prop = append(propstatOK.Prop,
			prop.Escaped("d:resourcetype", ""),
			prop.Escaped("d:getcontentlength", size),
		)
	}

	if len(propstatOK.Prop) > 0 {
		response.Propstat = append(response.Propstat, propstatOK)
	}

	return response
}
