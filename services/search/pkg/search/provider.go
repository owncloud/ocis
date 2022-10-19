package search

import (
	"context"
	"fmt"
	"sort"
	"strings"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	user "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	rpcv1beta1 "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/cs3org/reva/v2/pkg/errtypes"
	"github.com/cs3org/reva/v2/pkg/events"
	sdk "github.com/cs3org/reva/v2/pkg/sdk/common"
	"github.com/cs3org/reva/v2/pkg/storagespace"
	"github.com/cs3org/reva/v2/pkg/utils"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/services/search/pkg/content"
	"github.com/owncloud/ocis/v2/services/search/pkg/engine"
	"github.com/owncloud/ocis/v2/services/search/pkg/indexer"

	searchmsg "github.com/owncloud/ocis/v2/protogen/gen/ocis/messages/search/v0"
	searchsvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/search/v0"
)

// Permissions is copied from reva internal conversion pkg
type Permissions uint

// consts are copied from reva internal conversion pkg
const (
	// PermissionInvalid represents an invalid permission
	PermissionInvalid Permissions = 0
	// PermissionRead grants read permissions on a resource
	PermissionRead Permissions = 1 << (iota - 1)
	// PermissionWrite grants write permissions on a resource
	PermissionWrite
	// PermissionCreate grants create permissions on a resource
	PermissionCreate
	// PermissionDelete grants delete permissions on a resource
	PermissionDelete
	// PermissionShare grants share permissions on a resource
	PermissionShare
)

// ListenEvents are the events the search service is listening to
var ListenEvents = []events.Unmarshaller{
	events.ItemTrashed{},
	events.ItemRestored{},
	events.ItemMoved{},
	events.ContainerCreated{},
	events.FileUploaded{},
	events.FileTouched{},
	events.FileVersionRestored{},
}

// Provider is responsible for indexing spaces and pass on a search
// to it's underlying engine.
type Provider struct {
	logger    log.Logger
	gateway   gateway.GatewayAPIClient
	engine    engine.Engine
	extractor content.Extractor
	indexer   indexer.Indexer
}

// NewProvider creates a new Provider instance.
func NewProvider(gw gateway.GatewayAPIClient, eng engine.Engine, extractor content.Extractor, indexer indexer.Indexer, logger log.Logger) *Provider {
	return &Provider{
		gateway:   gw,
		engine:    eng,
		logger:    logger,
		extractor: extractor,
		indexer:   indexer,
	}
}

func (p *Provider) Search(ctx context.Context, req *searchsvc.SearchRequest) (*searchsvc.SearchResponse, error) {
	if req.Query == "" {
		return nil, errtypes.BadRequest("empty query provided")
	}
	p.logger.Debug().Str("query", req.Query).Msg("performing a search")

	listSpacesRes, err := p.gateway.ListStorageSpaces(ctx, &provider.ListStorageSpacesRequest{
		Filters: []*provider.ListStorageSpacesRequest_Filter{
			{
				Type: provider.ListStorageSpacesRequest_Filter_TYPE_SPACE_TYPE,
				Term: &provider.ListStorageSpacesRequest_Filter_SpaceType{SpaceType: "+grant"},
			},
		},
	})
	if err != nil {
		p.logger.Error().Err(err).Msg("failed to list the user's storage spaces")
		return nil, err
	}

	mountpointMap := map[string]string{}
	for _, space := range listSpacesRes.StorageSpaces {
		if space.SpaceType != "mountpoint" {
			continue
		}
		opaqueMap := sdk.DecodeOpaqueMap(space.Opaque)
		grantSpaceID := storagespace.FormatResourceID(provider.ResourceId{
			StorageId: opaqueMap["grantStorageID"],
			SpaceId:   opaqueMap["grantSpaceID"],
			OpaqueId:  opaqueMap["grantOpaqueID"],
		})
		mountpointMap[grantSpaceID] = space.Id.OpaqueId
	}

	matches := matchArray{}
	total := int32(0)
	for _, space := range listSpacesRes.StorageSpaces {
		searchRootID := &searchmsg.ResourceID{
			StorageId: space.Root.StorageId,
			SpaceId:   space.Root.SpaceId,
			OpaqueId:  space.Root.OpaqueId,
		}

		if req.Ref != nil &&
			(req.Ref.ResourceId.StorageId != searchRootID.StorageId ||
				req.Ref.ResourceId.SpaceId != searchRootID.SpaceId ||
				req.Ref.ResourceId.OpaqueId != searchRootID.OpaqueId) {
			continue
		}

		var (
			mountpointRootID *searchmsg.ResourceID
			rootName         string
			permissions      *provider.ResourcePermissions
		)
		mountpointPrefix := ""
		switch space.SpaceType {
		case "mountpoint":
			continue // mountpoint spaces are only "links" to the shared spaces. we have to search the shared "grant" space instead
		case "grant":
			// In case of grant spaces we search the root of the outer space and translate the paths to the according mountpoint
			searchRootID.OpaqueId = space.Root.SpaceId
			mountpointID, ok := mountpointMap[space.Id.OpaqueId]
			if !ok {
				p.logger.Warn().Interface("space", space).Msg("could not find mountpoint space for grant space")
				continue
			}
			gpRes, err := p.gateway.GetPath(ctx, &provider.GetPathRequest{
				ResourceId: space.Root,
			})
			if err != nil {
				p.logger.Error().Err(err).Str("space", space.Id.OpaqueId).Msg("failed to get path for grant space root")
				continue
			}
			if gpRes.Status.Code != rpcv1beta1.Code_CODE_OK {
				p.logger.Error().Interface("status", gpRes.Status).Str("space", space.Id.OpaqueId).Msg("failed to get path for grant space root")
				continue
			}
			mountpointPrefix = utils.MakeRelativePath(gpRes.Path)
			sid, spid, oid, err := storagespace.SplitID(mountpointID)
			if err != nil {
				p.logger.Error().Err(err).Str("space", space.Id.OpaqueId).Str("mountpointId", mountpointID).Msg("invalid mountpoint space id")
				continue
			}
			mountpointRootID = &searchmsg.ResourceID{
				StorageId: sid,
				SpaceId:   spid,
				OpaqueId:  oid,
			}
			rootName = space.GetRootInfo().GetPath()
			permissions = space.GetRootInfo().GetPermissionSet()
			p.logger.Debug().Interface("grantSpace", space).Interface("mountpointRootId", mountpointRootID).Msg("searching a grant")
		case "personal":
			permissions = space.GetRootInfo().GetPermissionSet()
		}

		res, err := p.engine.Search(ctx, &searchsvc.SearchIndexRequest{
			Query: req.Query,
			Ref: &searchmsg.Reference{
				ResourceId: searchRootID,
				Path:       mountpointPrefix,
			},
			PageSize: req.PageSize,
		})
		if err != nil {
			p.logger.Error().Err(err).Str("space", space.Id.OpaqueId).Msg("failed to search the index")
			return nil, err
		}
		p.logger.Debug().Str("space", space.Id.OpaqueId).Int("hits", len(res.Matches)).Msg("space search done")

		total += res.TotalMatches
		for _, match := range res.Matches {
			if mountpointPrefix != "" {
				match.Entity.Ref.Path = utils.MakeRelativePath(strings.TrimPrefix(match.Entity.Ref.Path, mountpointPrefix))
			}
			if mountpointRootID != nil {
				match.Entity.Ref.ResourceId = mountpointRootID
			}
			match.Entity.ShareRootName = rootName

			isShared := match.GetEntity().GetRef().GetResourceId().GetSpaceId() == utils.ShareStorageSpaceID
			isMountpoint := isShared && match.GetEntity().GetRef().GetPath() == "."
			isDir := match.GetEntity().GetMimeType() == "httpd/unix-directory"
			match.Entity.Permissions = convertToWebDAVPermissions(isShared, isMountpoint, isDir, permissions)
			matches = append(matches, match)
		}
	}

	// compile one sorted list of matches from all spaces and apply the limit if needed
	sort.Sort(matches)
	limit := req.PageSize
	if limit == 0 {
		limit = 200
	}
	if int32(len(matches)) > limit && limit != -1 {
		matches = matches[0:limit]
	}

	return &searchsvc.SearchResponse{
		Matches:      matches,
		TotalMatches: total,
	}, nil
}

// IndexSpace (re)indexes all resources of a given space.
func (p *Provider) IndexSpace(ctx context.Context, req *searchsvc.IndexSpaceRequest) (*searchsvc.IndexSpaceResponse, error) {
	err := p.indexer.IndexSpace(ctx, &provider.StorageSpaceId{OpaqueId: req.SpaceId}, &user.UserId{OpaqueId: req.UserId})
	if err != nil {
		return nil, err
	}
	return &searchsvc.IndexSpaceResponse{}, nil
}

// NOTE: this converts CS3 to WebDAV permissions
// since conversions pkg is reva internal we have no other choice than to duplicate the logic
func convertToWebDAVPermissions(isShared, isMountpoint, isDir bool, p *provider.ResourcePermissions) string {
	if p == nil {
		return ""
	}
	var b strings.Builder
	if isShared {
		fmt.Fprintf(&b, "S")
	}
	if p.ListContainer &&
		p.ListFileVersions &&
		p.ListRecycle &&
		p.Stat &&
		p.GetPath &&
		p.GetQuota &&
		p.InitiateFileDownload {
		fmt.Fprintf(&b, "R")
	}
	if isMountpoint {
		fmt.Fprintf(&b, "M")
	}
	if p.Delete &&
		p.PurgeRecycle {
		fmt.Fprintf(&b, "D")
	}
	if p.InitiateFileUpload &&
		p.RestoreFileVersion &&
		p.RestoreRecycleItem {
		fmt.Fprintf(&b, "NV")
		if !isDir {
			fmt.Fprintf(&b, "W")
		}
	}
	if isDir &&
		p.ListContainer &&
		p.Stat &&
		p.CreateContainer &&
		p.InitiateFileUpload {
		fmt.Fprintf(&b, "CK")
	}
	return b.String()
}
