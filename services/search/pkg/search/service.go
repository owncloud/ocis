package search

import (
	"context"
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"time"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	user "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	rpcv1beta1 "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/cs3org/reva/v2/pkg/errtypes"
	sdk "github.com/cs3org/reva/v2/pkg/sdk/common"
	"github.com/cs3org/reva/v2/pkg/storage/utils/walker"
	"github.com/cs3org/reva/v2/pkg/storagespace"
	"github.com/cs3org/reva/v2/pkg/utils"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	searchmsg "github.com/owncloud/ocis/v2/protogen/gen/ocis/messages/search/v0"
	searchsvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/search/v0"
	"github.com/owncloud/ocis/v2/services/search/pkg/content"
	"github.com/owncloud/ocis/v2/services/search/pkg/engine"
)

//go:generate mockery --name=Searcher

// Searcher is the interface to the SearchService
type Searcher interface {
	Search(ctx context.Context, req *searchsvc.SearchRequest) (*searchsvc.SearchResponse, error)
	IndexSpace(rid *provider.StorageSpaceId, uId *user.UserId) error
	TrashItem(rid *provider.ResourceId)
	UpsertItem(ref *provider.Reference, uid *user.UserId)
	RestoreItem(ref *provider.Reference, uid *user.UserId)
	MoveItem(ref *provider.Reference, uid *user.UserId)
}

// SearchService is responsible for indexing spaces and pass on a search
// to it's underlying engine.
type Service struct {
	logger    log.Logger
	gateway   gateway.GatewayAPIClient
	engine    engine.Engine
	extractor content.Extractor
	secret    string
}

// NewService creates a new Provider instance.
func NewService(gw gateway.GatewayAPIClient, eng engine.Engine, extractor content.Extractor, logger log.Logger, secret string) *Service {
	return &Service{
		gateway:   gw,
		engine:    eng,
		secret:    secret,
		logger:    logger,
		extractor: extractor,
	}
}

// Search processes a search request and passes it down to the engine.
func (s *Service) Search(ctx context.Context, req *searchsvc.SearchRequest) (*searchsvc.SearchResponse, error) {
	if req.Query == "" {
		return nil, errtypes.BadRequest("empty query provided")
	}
	s.logger.Debug().Str("query", req.Query).Msg("performing a search")

	listSpacesRes, err := s.gateway.ListStorageSpaces(ctx, &provider.ListStorageSpacesRequest{
		Filters: []*provider.ListStorageSpacesRequest_Filter{
			{
				Type: provider.ListStorageSpacesRequest_Filter_TYPE_SPACE_TYPE,
				Term: &provider.ListStorageSpacesRequest_Filter_SpaceType{SpaceType: "+grant"},
			},
		},
	})
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to list the user's storage spaces")
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
				s.logger.Warn().Interface("space", space).Msg("could not find mountpoint space for grant space")
				continue
			}
			gpRes, err := s.gateway.GetPath(ctx, &provider.GetPathRequest{
				ResourceId: space.Root,
			})
			if err != nil {
				s.logger.Error().Err(err).Str("space", space.Id.OpaqueId).Msg("failed to get path for grant space root")
				continue
			}
			if gpRes.Status.Code != rpcv1beta1.Code_CODE_OK {
				s.logger.Error().Interface("status", gpRes.Status).Str("space", space.Id.OpaqueId).Msg("failed to get path for grant space root")
				continue
			}
			mountpointPrefix = utils.MakeRelativePath(gpRes.Path)
			sid, spid, oid, err := storagespace.SplitID(mountpointID)
			if err != nil {
				s.logger.Error().Err(err).Str("space", space.Id.OpaqueId).Str("mountpointId", mountpointID).Msg("invalid mountpoint space id")
				continue
			}
			mountpointRootID = &searchmsg.ResourceID{
				StorageId: sid,
				SpaceId:   spid,
				OpaqueId:  oid,
			}
			rootName = space.GetRootInfo().GetPath()
			permissions = space.GetRootInfo().GetPermissionSet()
			s.logger.Debug().Interface("grantSpace", space).Interface("mountpointRootId", mountpointRootID).Msg("searching a grant")
		case "personal":
			permissions = space.GetRootInfo().GetPermissionSet()
		}

		res, err := s.engine.Search(ctx, &searchsvc.SearchIndexRequest{
			Query: req.Query,
			Ref: &searchmsg.Reference{
				ResourceId: searchRootID,
				Path:       mountpointPrefix,
			},
			PageSize: req.PageSize,
		})
		if err != nil {
			s.logger.Error().Err(err).Str("space", space.Id.OpaqueId).Msg("failed to search the index")
			return nil, err
		}
		s.logger.Debug().Str("space", space.Id.OpaqueId).Int("hits", len(res.Matches)).Msg("space search done")

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
func (s *Service) IndexSpace(spaceID *provider.StorageSpaceId, uId *user.UserId) error {
	ownerCtx, err := getAuthContext(&user.User{Id: uId}, s.gateway, s.secret, s.logger)
	if err != nil {
		return err
	}

	rootID, err := storagespace.ParseID(spaceID.OpaqueId)
	if err != nil {
		s.logger.Error().Err(err).Msg("invalid space id")
		return err
	}
	if rootID.StorageId == "" || rootID.SpaceId == "" {
		s.logger.Error().Err(err).Msg("invalid space id")
		return fmt.Errorf("invalid space id")
	}
	rootID.OpaqueId = rootID.SpaceId

	w := walker.NewWalker(s.gateway)
	err = w.Walk(ownerCtx, &rootID, func(wd string, info *provider.ResourceInfo, err error) error {
		if err != nil {
			s.logger.Error().Err(err).Msg("error walking the tree")
			return err
		}

		if info == nil {
			return nil
		}

		ref := &provider.Reference{
			Path:       utils.MakeRelativePath(filepath.Join(wd, info.Path)),
			ResourceId: &rootID,
		}
		s.logger.Debug().Str("path", ref.Path).Msg("Walking tree")

		searchRes, err := s.engine.Search(ownerCtx, &searchsvc.SearchIndexRequest{
			Query: "+ID:" + storagespace.FormatResourceID(*info.Id) + ` +Mtime:>="` + utils.TSToTime(info.Mtime).Format(time.RFC3339Nano) + `"`,
		})

		if err == nil && len(searchRes.Matches) >= 1 {
			if info.Type == provider.ResourceType_RESOURCE_TYPE_CONTAINER {
				s.logger.Debug().Str("path", ref.Path).Msg("subtree hasn't changed. Skipping.")
				return filepath.SkipDir
			}
			s.logger.Debug().Str("path", ref.Path).Msg("element hasn't changed. Skipping.")
			return nil
		}

		s.UpsertItem(ref, uId)

		return nil
	})

	if err != nil {
		return err
	}

	logDocCount(s.engine, s.logger)

	return nil
}

func (s *Service) TrashItem(rid *provider.ResourceId) {
	err := s.engine.Delete(storagespace.FormatResourceID(*rid))
	if err != nil {
		s.logger.Error().Err(err).Interface("Id", rid).Msg("failed to remove item from index")
	}
}

func (s *Service) UpsertItem(ref *provider.Reference, uid *user.UserId) {
	ctx, stat, path := s.resInfo(uid, ref)
	if ctx == nil || stat == nil || path == "" {
		return
	}

	doc, err := s.extractor.Extract(ctx, stat.Info)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to extract resource content")
		return
	}

	r := engine.Resource{
		ID: storagespace.FormatResourceID(*stat.Info.Id),
		RootID: storagespace.FormatResourceID(provider.ResourceId{
			StorageId: stat.Info.Id.StorageId,
			OpaqueId:  stat.Info.Id.SpaceId,
			SpaceId:   stat.Info.Id.SpaceId,
		}),
		Path:     utils.MakeRelativePath(path),
		Type:     uint64(stat.Info.Type),
		Document: doc,
	}

	if parentId := stat.GetInfo().GetParentId(); parentId != nil {
		r.ParentID = storagespace.FormatResourceID(*parentId)
	}

	if err = s.engine.Upsert(r.ID, r); err != nil {
		s.logger.Error().Err(err).Msg("error adding updating the resource in the index")
	} else {
		logDocCount(s.engine, s.logger)
	}
}

func (s *Service) RestoreItem(ref *provider.Reference, uid *user.UserId) {
	ctx, stat, path := s.resInfo(uid, ref)
	if ctx == nil || stat == nil || path == "" {
		return
	}

	if err := s.engine.Restore(storagespace.FormatResourceID(*stat.Info.Id)); err != nil {
		s.logger.Error().Err(err).Msg("failed to restore the changed resource in the index")
	}
}

func (s *Service) MoveItem(ref *provider.Reference, uid *user.UserId) {
	ctx, stat, path := s.resInfo(uid, ref)
	if ctx == nil || stat == nil || path == "" {
		return
	}

	if err := s.engine.Move(storagespace.FormatResourceID(*stat.GetInfo().GetId()), storagespace.FormatResourceID(*stat.GetInfo().GetParentId()), path); err != nil {
		s.logger.Error().Err(err).Msg("failed to move the changed resource in the index")
	}
}

func (s *Service) resInfo(uid *user.UserId, ref *provider.Reference) (context.Context, *provider.StatResponse, string) {
	ownerCtx, err := getAuthContext(&user.User{Id: uid}, s.gateway, s.secret, s.logger)
	if err != nil {
		return nil, nil, ""
	}

	statRes, err := statResource(ownerCtx, ref, s.gateway, s.logger)
	if err != nil {
		return nil, nil, ""
	}

	r, err := ResolveReference(ownerCtx, ref, statRes.GetInfo(), s.gateway)
	if err != nil {
		return nil, nil, ""
	}

	return ownerCtx, statRes, r.GetPath()
}
