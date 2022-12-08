package provider

import (
	"context"
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"google.golang.org/grpc/metadata"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	user "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	rpcv1beta1 "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	revactx "github.com/cs3org/reva/v2/pkg/ctx"
	"github.com/cs3org/reva/v2/pkg/errtypes"
	"github.com/cs3org/reva/v2/pkg/events"
	sdk "github.com/cs3org/reva/v2/pkg/sdk/common"
	"github.com/cs3org/reva/v2/pkg/storage/utils/walker"
	"github.com/cs3org/reva/v2/pkg/storagespace"
	"github.com/cs3org/reva/v2/pkg/utils"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/services/search/pkg/search"
	searchTracing "github.com/owncloud/ocis/v2/services/search/pkg/tracing"

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

var ListenEvents = []events.Unmarshaller{
	events.ItemTrashed{},
	events.ItemRestored{},
	events.ItemMoved{},
	events.ContainerCreated{},
	events.FileUploaded{},
	events.FileTouched{},
	events.FileVersionRestored{},
}

type Provider struct {
	logger            log.Logger
	gwClient          gateway.GatewayAPIClient
	indexClient       search.IndexClient
	machineAuthAPIKey string

	indexSpaceDebouncer *SpaceDebouncer
}

type MatchArray []*searchmsg.Match

func (s MatchArray) Len() int {
	return len(s)
}
func (s MatchArray) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s MatchArray) Less(i, j int) bool {
	return s[i].Score > s[j].Score
}

func New(gwClient gateway.GatewayAPIClient, indexClient search.IndexClient, machineAuthAPIKey string, eventsChan <-chan interface{}, debounceDuration int, logger log.Logger) *Provider {
	p := &Provider{
		gwClient:          gwClient,
		indexClient:       indexClient,
		machineAuthAPIKey: machineAuthAPIKey,
		logger:            logger,
	}

	p.indexSpaceDebouncer = NewSpaceDebouncer(time.Duration(debounceDuration)*time.Millisecond, func(id *provider.StorageSpaceId, userID *user.UserId) {
		err := p.doIndexSpace(context.Background(), id, userID)
		if err != nil {
			p.logger.Error().Err(err).Interface("spaceID", id).Interface("userID", userID).Msg("error while indexing a space")
		}
	})

	go func() {
		for {
			ev := <-eventsChan
			go func() {
				time.Sleep(1 * time.Second) // Give some time to let everything settle down before trying to access it when indexing
				p.handleEvent(ev)
			}()
		}
	}()

	return p
}

// NewWithDebouncer returns a new provider with a customer index space debouncer
func NewWithDebouncer(gwClient gateway.GatewayAPIClient, indexClient search.IndexClient, machineAuthAPIKey string, eventsChan <-chan interface{}, logger log.Logger, debouncer *SpaceDebouncer) *Provider {
	p := New(gwClient, indexClient, machineAuthAPIKey, eventsChan, 0, logger)
	p.indexSpaceDebouncer = debouncer
	return p
}

func (p *Provider) Search(ctx context.Context, req *searchsvc.SearchRequest) (*searchsvc.SearchResponse, error) {
	ctx, span := searchTracing.TraceProvider.Tracer("search").Start(ctx, "search")
	defer span.End()
	span.SetAttributes(attribute.String("query", req.GetQuery()))
	if req.Query == "" {
		return nil, errtypes.BadRequest("empty query provided")
	}
	p.logger.Debug().Str("query", req.Query).Msg("performing a search")

	listSpacesRes, err := p.gwClient.ListStorageSpaces(ctx, &provider.ListStorageSpacesRequest{
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
		grantSpaceId := storagespace.FormatResourceID(provider.ResourceId{
			StorageId: opaqueMap["grantStorageID"],
			SpaceId:   opaqueMap["grantSpaceID"],
			OpaqueId:  opaqueMap["grantOpaqueID"],
		})
		mountpointMap[grantSpaceId] = space.Id.OpaqueId
	}

	matches := MatchArray{}
	total := int32(0)
	for _, space := range listSpacesRes.StorageSpaces {
		searchRootId := &searchmsg.ResourceID{
			StorageId: space.Root.StorageId,
			SpaceId:   space.Root.SpaceId,
			OpaqueId:  space.Root.OpaqueId,
		}

		if req.Ref != nil &&
			(req.Ref.ResourceId.StorageId != searchRootId.StorageId ||
				req.Ref.ResourceId.SpaceId != searchRootId.SpaceId ||
				req.Ref.ResourceId.OpaqueId != searchRootId.OpaqueId) {
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
			searchRootId.OpaqueId = space.Root.SpaceId
			mountpointID, ok := mountpointMap[space.Id.OpaqueId]
			if !ok {
				p.logger.Warn().Interface("space", space).Msg("could not find mountpoint space for grant space")
				continue
			}
			gpRes, err := p.gwClient.GetPath(ctx, &provider.GetPathRequest{
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

			rootName = filepath.Join("/", filepath.Base(gpRes.GetPath()))
			permissions = space.GetRootInfo().GetPermissionSet()
			p.logger.Debug().Interface("grantSpace", space).Interface("mountpointRootId", mountpointRootID).Msg("searching a grant")
		case "personal":
			permissions = space.GetRootInfo().GetPermissionSet()
		}

		res, err := p.indexClient.Search(ctx, &searchsvc.SearchIndexRequest{
			Query: formatQuery(req.Query),
			Ref: &searchmsg.Reference{
				ResourceId: searchRootId,
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
	span.SetAttributes(attribute.Int("num_matches", len(matches)))
	span.SetAttributes(attribute.Int("total_matches", int(total)))
	limit := req.PageSize
	if limit == 0 {
		limit = 200
	}
	if int32(len(matches)) > limit {
		matches = matches[0:limit]
	}

	return &searchsvc.SearchResponse{
		Matches:      matches,
		TotalMatches: total,
	}, nil
}

func (p *Provider) IndexSpace(ctx context.Context, req *searchsvc.IndexSpaceRequest) (*searchsvc.IndexSpaceResponse, error) {
	err := p.doIndexSpace(ctx, &provider.StorageSpaceId{OpaqueId: req.SpaceId}, &user.UserId{OpaqueId: req.UserId})
	if err != nil {
		return nil, err
	}
	return &searchsvc.IndexSpaceResponse{}, nil
}

func (p *Provider) doIndexSpace(ctx context.Context, spaceID *provider.StorageSpaceId, userID *user.UserId) error {
	ctx, span := searchTracing.TraceProvider.Tracer("search").Start(ctx, "index space")
	defer span.End()
	authRes, err := p.gwClient.Authenticate(ctx, &gateway.AuthenticateRequest{
		Type:         "machine",
		ClientId:     "userid:" + userID.OpaqueId,
		ClientSecret: p.machineAuthAPIKey,
	})
	if err != nil || authRes.GetStatus().GetCode() != rpc.Code_CODE_OK {
		return err
	}
	span.SetAttributes(attribute.String("user_id", userID.GetOpaqueId()))
	span.SetAttributes(attribute.String("space_id", spaceID.GetOpaqueId()))

	if authRes.GetStatus().GetCode() != rpc.Code_CODE_OK {
		return fmt.Errorf("could not get authenticated context for user")
	}
	ownerCtx := metadata.AppendToOutgoingContext(ctx, revactx.TokenHeader, authRes.Token)

	// Walk the space and index all files
	walker := walker.NewWalker(p.gwClient)
	rootID, err := storagespace.ParseID(spaceID.OpaqueId)
	if err != nil {
		p.logger.Error().Err(err).Msg("invalid space id")
		return err
	}
	if rootID.StorageId == "" || rootID.SpaceId == "" {
		p.logger.Error().Err(err).Msg("invalid space id")
		return fmt.Errorf("invalid space id")
	}
	rootID.OpaqueId = rootID.SpaceId

	err = walker.Walk(ownerCtx, &rootID, func(wd string, info *provider.ResourceInfo, err error) error {
		if err != nil {
			p.logger.Error().Err(err).Msg("error walking the tree")
			return err
		}

		if info == nil {
			return nil
		}

		ref := &provider.Reference{
			Path:       utils.MakeRelativePath(filepath.Join(wd, info.Path)),
			ResourceId: &rootID,
		}
		p.logger.Debug().Str("path", ref.Path).Msg("Walking tree")

		// Has this item/subtree changed?
		searchRes, err := p.indexClient.Search(ownerCtx, &searchsvc.SearchIndexRequest{
			Query: "+ID:" + storagespace.FormatResourceID(*info.Id) + ` +Mtime:>="` + utils.TSToTime(info.Mtime).Format(time.RFC3339Nano) + `"`,
		})
		if err == nil && len(searchRes.Matches) >= 1 {
			if info.Type == provider.ResourceType_RESOURCE_TYPE_CONTAINER {
				p.logger.Debug().Str("path", ref.Path).Msg("subtree hasn't changed. Skipping.")
				return filepath.SkipDir
			}
			p.logger.Debug().Str("path", ref.Path).Msg("element hasn't changed. Skipping.")
			return nil
		}

		err = p.indexClient.Add(ref, info)
		if err != nil {
			p.logger.Error().Err(err).Msg("error adding resource to the index")
		} else {
			p.logger.Debug().Interface("ref", ref).Msg("added resource to index")
		}
		return nil
	})
	if err != nil {
		return err
	}

	p.logDocCount()
	return nil
}

func (p *Provider) logDocCount() {
	c, err := p.indexClient.DocCount()
	if err != nil {
		p.logger.Error().Err(err).Msg("error getting document count from the index")
	}
	p.logger.Debug().Interface("count", c).Msg("new document count")
}

func formatQuery(q string) string {
	query := q
	fields := []string{"RootID", "Path", "ID", "Name", "Size", "Mtime", "MimeType", "Type"}
	for _, field := range fields {
		query = strings.ReplaceAll(query, strings.ToLower(field)+":", field+":")
	}

	if strings.Contains(query, ":") {
		return query // Sophisticated field based search
	}

	// this is a basic filename search
	return "Name:*" + strings.ReplaceAll(strings.ToLower(query), " ", `\ `) + "*"
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
