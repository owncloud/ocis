package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	user "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/jellydator/ttlcache/v2"
	revactx "github.com/owncloud/reva/v2/pkg/ctx"
	"github.com/owncloud/reva/v2/pkg/errtypes"
	"github.com/owncloud/reva/v2/pkg/events/stream"
	"github.com/owncloud/reva/v2/pkg/rgrpc/todo/pool"
	"github.com/owncloud/reva/v2/pkg/token"
	"github.com/owncloud/reva/v2/pkg/token/manager/jwt"
	"github.com/owncloud/reva/v2/pkg/utils"
	merrors "go-micro.dev/v4/errors"
	"go-micro.dev/v4/metadata"
	grpcmetadata "google.golang.org/grpc/metadata"

	"github.com/owncloud/ocis/v2/ocis-pkg/generators"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/ocis-pkg/registry"
	v0 "github.com/owncloud/ocis/v2/protogen/gen/ocis/messages/search/v0"
	searchsvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/search/v0"
	"github.com/owncloud/ocis/v2/services/search/pkg/config"
	"github.com/owncloud/ocis/v2/services/search/pkg/content"
	"github.com/owncloud/ocis/v2/services/search/pkg/engine"
	bleveEngine "github.com/owncloud/ocis/v2/services/search/pkg/engine/bleve"
	"github.com/owncloud/ocis/v2/services/search/pkg/query/bleve"
	"github.com/owncloud/ocis/v2/services/search/pkg/search"
)

// NewHandler returns a service implementation for Service.
func NewHandler(opts ...Option) (searchsvc.SearchProviderHandler, func(), error) {
	teardown := func() {}
	options := newOptions(opts...)
	logger := options.Logger
	cfg := options.Config

	// initialize search engine
	var eng engine.Engine
	switch cfg.Engine.Type {
	case "bleve":
		bleveMapping, err := engine.BuildBleveMapping()
		if err != nil {
			return nil, teardown, err
		}

		var indexGetter bleveEngine.IndexGetter
		indexGetter = bleveEngine.NewIndexGetterPersistent(cfg.Engine.Bleve.Datapath, bleveMapping)
		if cfg.Engine.Bleve.Scale {
			indexGetter = bleveEngine.NewIndexGetterPersistentScale(cfg.Engine.Bleve.Datapath, bleveMapping)
		}

		bleveEngine := engine.NewBleveEngine(indexGetter, bleve.DefaultCreator)

		teardown = func() {
			_ = bleveEngine.Close()
		}
		eng = bleveEngine

	default:
		return nil, teardown, fmt.Errorf("unknown search engine: %s", cfg.Engine.Type)
	}

	// initialize gateway
	selector, err := pool.GatewaySelector(cfg.Reva.Address, pool.WithRegistry(registry.GetRegistry()), pool.WithTracerProvider(options.TracerProvider))
	if err != nil {
		logger.Fatal().Err(err).Msg("could not get reva gateway selector")
		return nil, teardown, err
	}
	// initialize search content extractor
	var extractor content.Extractor
	switch cfg.Extractor.Type {
	case "basic":
		if extractor, err = content.NewBasicExtractor(logger); err != nil {
			return nil, teardown, err
		}
	case "tika":
		if extractor, err = content.NewTikaExtractor(selector, logger, cfg); err != nil {
			return nil, teardown, err
		}
	default:
		return nil, teardown, fmt.Errorf("unknown search extractor: %s", cfg.Extractor.Type)
	}

	connName := generators.GenerateConnectionName(cfg.Service.Name, generators.NTypeBus)
	bus, err := stream.NatsFromConfig(connName, false, stream.NatsConfig{
		Endpoint:             cfg.Events.Endpoint,
		Cluster:              cfg.Events.Cluster,
		EnableTLS:            cfg.Events.EnableTLS,
		TLSInsecure:          cfg.Events.TLSInsecure,
		TLSRootCACertificate: cfg.Events.TLSRootCACertificate,
		AuthUsername:         cfg.Events.AuthUsername,
		AuthPassword:         cfg.Events.AuthPassword,
	})
	if err != nil {
		return nil, teardown, err
	}

	ss := search.NewService(selector, eng, extractor, logger, cfg)

	// setup event handling
	if err := search.HandleEvents(ss, bus, logger, cfg); err != nil {
		return nil, teardown, err
	}

	cache := ttlcache.NewCache()
	if err := cache.SetTTL(time.Second); err != nil {
		return nil, teardown, err
	}

	tokenManager, err := jwt.New(map[string]interface{}{
		"secret":  options.JWTSecret,
		"expires": int64(24 * 60 * 60),
	})
	if err != nil {
		return nil, teardown, err
	}

	return &Service{
		id:           cfg.GRPC.Namespace + "." + cfg.Service.Name,
		log:          logger,
		searcher:     ss,
		cache:        cache,
		tokenManager: tokenManager,
		gws:          selector,
		cfg:          cfg,
	}, teardown, nil
}

// Service implements the searchServiceHandler interface
type Service struct {
	id           string
	log          log.Logger
	searcher     search.Searcher
	cache        *ttlcache.Cache
	tokenManager token.Manager
	gws          *pool.Selector[gateway.GatewayAPIClient]
	cfg          *config.Config
}

// Search handles the search
func (s Service) Search(ctx context.Context, in *searchsvc.SearchRequest, out *searchsvc.SearchResponse) error {
	// Get token from the context (go-micro) and make it known to the reva client too (grpc)
	t, ok := metadata.Get(ctx, revactx.TokenHeader)
	if !ok {
		s.log.Error().Msg("Could not get token from context")
		return errors.New("could not get token from context")
	}
	ctx = grpcmetadata.AppendToOutgoingContext(ctx, revactx.TokenHeader, t)

	// unpack user
	u, _, err := s.tokenManager.DismantleToken(ctx, t)
	if err != nil {
		return err
	}
	ctx = revactx.ContextSetUser(ctx, u)

	key := cacheKey(in.Query, in.PageSize, in.Ref, u)
	res, ok := s.FromCache(key)
	if !ok {
		var err error
		res, err = s.searcher.Search(ctx, &searchsvc.SearchRequest{
			Query:    in.Query,
			PageSize: in.PageSize,
			Ref:      in.Ref,
		})
		if err != nil {
			switch err.(type) {
			case errtypes.BadRequest:
				return merrors.BadRequest(s.id, "%s", err.Error())
			default:
				return merrors.InternalServerError(s.id, "%s", err.Error())
			}
		}

		s.Cache(key, res)
	}

	out.Matches = res.Matches
	out.TotalMatches = res.TotalMatches
	out.NextPageToken = res.NextPageToken
	return nil
}

// IndexSpace (re)indexes all resources of a given space.
func (s Service) IndexSpace(_ context.Context, in *searchsvc.IndexSpaceRequest, _ *searchsvc.IndexSpaceResponse) error {
	if in.GetSpaceId() != "" {
		return s.searcher.IndexSpace(&provider.StorageSpaceId{OpaqueId: in.GetSpaceId()})
	}

	// index all spaces instead
	gwc, err := s.gws.Next()
	if err != nil {
		return err
	}

	ctx, err := utils.GetServiceUserContext(s.cfg.ServiceAccount.ServiceAccountID, gwc, s.cfg.ServiceAccount.ServiceAccountSecret)
	if err != nil {
		return err
	}

	resp, err := gwc.ListStorageSpaces(ctx, &provider.ListStorageSpacesRequest{})
	if err != nil {
		return err
	}

	if resp.GetStatus().GetCode() != rpc.Code_CODE_OK {
		return errors.New(resp.GetStatus().GetMessage())
	}

	for _, space := range resp.GetStorageSpaces() {
		if err := s.searcher.IndexSpace(space.GetId()); err != nil {
			return err
		}
	}

	return nil
}

// FromCache pulls a search result from cache
func (s Service) FromCache(key string) (*searchsvc.SearchResponse, bool) {
	v, err := s.cache.Get(key)
	if err != nil {
		return nil, false
	}

	sr, ok := v.(*searchsvc.SearchResponse)
	return sr, ok
}

// Cache caches the search result
func (s Service) Cache(key string, res *searchsvc.SearchResponse) {
	// lets ignore the error
	_ = s.cache.Set(key, res)
}

func cacheKey(query string, pagesize int32, ref *v0.Reference, user *user.User) string {
	return fmt.Sprintf("%s|%d|%s$%s!%s/%s|%s", query, pagesize, ref.GetResourceId().GetStorageId(), ref.GetResourceId().GetSpaceId(), ref.GetResourceId().GetOpaqueId(), ref.GetPath(), user.GetId().GetOpaqueId())
}
