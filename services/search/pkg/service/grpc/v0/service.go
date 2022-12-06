package service

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"github.com/cs3org/reva/v2/pkg/storagespace"
	"os"
	"time"

	revactx "github.com/cs3org/reva/v2/pkg/ctx"
	"github.com/cs3org/reva/v2/pkg/errtypes"
	"github.com/cs3org/reva/v2/pkg/events/server"
	"github.com/cs3org/reva/v2/pkg/rgrpc/todo/pool"
	"github.com/go-micro/plugins/v4/events/natsjs"
	"github.com/jellydator/ttlcache/v2"
	"github.com/owncloud/ocis/v2/services/search/pkg/content"
	"github.com/owncloud/ocis/v2/services/search/pkg/engine"
	"github.com/owncloud/ocis/v2/services/search/pkg/search"
	merrors "go-micro.dev/v4/errors"
	"go-micro.dev/v4/metadata"
	grpcmetadata "google.golang.org/grpc/metadata"

	user "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	ctxpkg "github.com/cs3org/reva/v2/pkg/ctx"
	ociscrypto "github.com/owncloud/ocis/v2/ocis-pkg/crypto"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	v0 "github.com/owncloud/ocis/v2/protogen/gen/ocis/messages/search/v0"
	searchsvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/search/v0"
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
		idx, err := engine.NewBleveIndex(cfg.Engine.Bleve.Datapath)
		if err != nil {
			return nil, teardown, err
		}

		teardown = func() {
			_ = idx.Close()
		}

		eng = engine.NewBleveEngine(idx)
	default:
		return nil, teardown, fmt.Errorf("unknown search engine: %s", cfg.Engine.Type)
	}

	// initialize gateway
	gw, err := pool.GetGatewayServiceClient(cfg.Reva.Address)
	if err != nil {
		logger.Fatal().Err(err).Str("addr", cfg.Reva.Address).Msg("could not get reva client")
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
		if extractor, err = content.NewTikaExtractor(gw, logger, cfg); err != nil {
			return nil, teardown, err
		}
	default:
		return nil, teardown, fmt.Errorf("unknown search extractor: %s", cfg.Extractor.Type)
	}

	var rootCAPool *x509.CertPool
	if cfg.Events.TLSRootCACertificate != "" {
		rootCrtFile, err := os.Open(cfg.Events.TLSRootCACertificate)
		if err != nil {
			return nil, teardown, err
		}

		rootCAPool, err = ociscrypto.NewCertPoolFromPEM(rootCrtFile)
		if err != nil {
			return nil, teardown, err
		}
		cfg.Events.TLSInsecure = false
	}

	tlsConf := &tls.Config{
		InsecureSkipVerify: cfg.Events.TLSInsecure, //nolint:gosec
		RootCAs:            rootCAPool,
	}
	bus, err := server.NewNatsStream(
		natsjs.TLSConfig(tlsConf),
		natsjs.Address(cfg.Events.Endpoint),
		natsjs.ClusterID(cfg.Events.Cluster),
	)
	if err != nil {
		return nil, teardown, err
	}

	ss := search.NewService(gw, eng, extractor, logger, cfg.MachineAuthAPIKey)

	// setup event handling
	if err := search.HandleEvents(ss, bus, logger, cfg); err != nil {
		return nil, teardown, err
	}

	cache := ttlcache.NewCache()
	if err := cache.SetTTL(time.Second); err != nil {
		return nil, teardown, err
	}

	return &Service{
		id:       cfg.GRPC.Namespace + "." + cfg.Service.Name,
		log:      logger,
		searcher: ss,
		cache:    cache,
	}, teardown, nil
}

// Service implements the searchServiceHandler interface
type Service struct {
	id       string
	log      log.Logger
	searcher search.Searcher
	cache    *ttlcache.Cache
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

	u, _ := ctxpkg.ContextGetUser(ctx)
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
				return merrors.BadRequest(s.id, err.Error())
			default:
				return merrors.InternalServerError(s.id, err.Error())
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
func (s Service) IndexSpace(ctx context.Context, in *searchsvc.IndexSpaceRequest, _ *searchsvc.IndexSpaceResponse) error {
	rid, err := storagespace.ParseID(in.SpaceId)
	if err != nil {
		return err
	}

	return s.searcher.IndexSpace(&rid, &user.UserId{OpaqueId: in.UserId})
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
