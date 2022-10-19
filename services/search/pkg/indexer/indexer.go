package indexer

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	user "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	ctxpkg "github.com/cs3org/reva/v2/pkg/ctx"
	"github.com/cs3org/reva/v2/pkg/storage/utils/walker"
	"github.com/cs3org/reva/v2/pkg/storagespace"
	"github.com/cs3org/reva/v2/pkg/utils"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	searchsvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/search/v0"
	"github.com/owncloud/ocis/v2/services/search/pkg/content"
	"github.com/owncloud/ocis/v2/services/search/pkg/engine"
	"google.golang.org/grpc/metadata"
)

//go:generate mockery --name=Indexer

type Indexer interface {
	IndexSpace(ctx context.Context, spaceID *provider.StorageSpaceId, userID *user.UserId) error
}

type indexer struct {
	gateway   gateway.GatewayAPIClient
	engine    engine.Engine
	extractor content.Extractor
	logger    log.Logger
	secret    string
}

func NewIndexer(gw gateway.GatewayAPIClient, eng engine.Engine, extractor content.Extractor, logger log.Logger, secret string) *indexer {
	return &indexer{
		gateway:   gw,
		engine:    eng,
		extractor: extractor,
		logger:    logger,
		secret:    secret,
	}
}

func (i *indexer) IndexSpace(ctx context.Context, spaceID *provider.StorageSpaceId, userID *user.UserId) error {
	authRes, err := i.gateway.Authenticate(ctx, &gateway.AuthenticateRequest{
		Type:         "machine",
		ClientId:     "userid:" + userID.OpaqueId,
		ClientSecret: i.secret,
	})
	if err != nil || authRes.GetStatus().GetCode() != rpc.Code_CODE_OK {
		return err
	}
	if authRes.GetStatus().GetCode() != rpc.Code_CODE_OK {
		return fmt.Errorf("could not get authenticated context for user")
	}
	ownerCtx := metadata.AppendToOutgoingContext(ctx, ctxpkg.TokenHeader, authRes.Token)

	// Walk the space and index all files
	walker := walker.NewWalker(i.gateway)
	rootID, err := storagespace.ParseID(spaceID.OpaqueId)
	if err != nil {
		i.logger.Error().Err(err).Msg("invalid space id")
		return err
	}
	if rootID.StorageId == "" || rootID.SpaceId == "" {
		i.logger.Error().Err(err).Msg("invalid space id")
		return fmt.Errorf("invalid space id")
	}
	rootID.OpaqueId = rootID.SpaceId

	err = walker.Walk(ownerCtx, &rootID, func(wd string, info *provider.ResourceInfo, err error) error {
		if err != nil {
			i.logger.Error().Err(err).Msg("error walking the tree")
			return err
		}

		if info == nil {
			return nil
		}

		ref := &provider.Reference{
			Path:       utils.MakeRelativePath(filepath.Join(wd, info.Path)),
			ResourceId: &rootID,
		}
		i.logger.Debug().Str("path", ref.Path).Msg("Walking tree")

		// Has this item/subtree changed?
		searchRes, err := i.engine.Search(ownerCtx, &searchsvc.SearchIndexRequest{
			Query: "+ID:" + storagespace.FormatResourceID(*info.Id) + ` +Mtime:>="` + utils.TSToTime(info.Mtime).Format(time.RFC3339Nano) + `"`,
		})
		if err == nil && len(searchRes.Matches) >= 1 {
			if info.Type == provider.ResourceType_RESOURCE_TYPE_CONTAINER {
				i.logger.Debug().Str("path", ref.Path).Msg("subtree hasn't changed. Skipping.")
				return filepath.SkipDir
			}
			i.logger.Debug().Str("path", ref.Path).Msg("element hasn't changed. Skipping.")
			return nil
		}

		doc, err := i.extractor.Extract(ownerCtx, info)
		if err != nil {
			i.logger.Error().Err(err).Msg("error extracting content")
		}

		var pid string
		if info.ParentId != nil {
			pid = storagespace.FormatResourceID(*info.ParentId)
		}
		r := engine.Resource{
			ID: storagespace.FormatResourceID(*info.Id),
			RootID: storagespace.FormatResourceID(provider.ResourceId{
				StorageId: info.Id.StorageId,
				OpaqueId:  info.Id.SpaceId,
				SpaceId:   info.Id.SpaceId,
			}),
			ParentID: pid,
			Path:     ref.Path,
			Type:     uint64(info.Type),
			Document: doc,
		}

		err = i.engine.Upsert(r.ID, r)
		if err != nil {
			i.logger.Error().Err(err).Msg("error adding resource to the index")
		} else {
			i.logger.Debug().Interface("ref", ref).Msg("added resource to index")
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
