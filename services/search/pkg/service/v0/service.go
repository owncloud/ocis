package service

import (
	"context"
	"errors"
	"path/filepath"

	"github.com/blevesearch/bleve/v2"
	revactx "github.com/cs3org/reva/v2/pkg/ctx"
	"github.com/cs3org/reva/v2/pkg/errtypes"
	"github.com/cs3org/reva/v2/pkg/events"
	"github.com/cs3org/reva/v2/pkg/events/server"
	"github.com/cs3org/reva/v2/pkg/rgrpc/todo/pool"
	"github.com/go-micro/plugins/v4/events/natsjs"
	merrors "go-micro.dev/v4/errors"
	"go-micro.dev/v4/metadata"
	grpcmetadata "google.golang.org/grpc/metadata"

	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	searchsvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/search/v0"
	"github.com/owncloud/ocis/v2/services/search/pkg/config"
	"github.com/owncloud/ocis/v2/services/search/pkg/search"
	"github.com/owncloud/ocis/v2/services/search/pkg/search/index"
	searchprovider "github.com/owncloud/ocis/v2/services/search/pkg/search/provider"
)

// NewHandler returns a service implementation for Service.
func NewHandler(opts ...Option) (searchsvc.SearchProviderHandler, error) {
	options := newOptions(opts...)
	logger := options.Logger
	cfg := options.Config

	// Connect to nats to listen for changes that need to trigger an index update
	evtsCfg := cfg.Events
	client, err := server.NewNatsStream(
		natsjs.Address(evtsCfg.Endpoint),
		natsjs.ClusterID(evtsCfg.Cluster),
	)
	if err != nil {
		return nil, err
	}
	evts, err := events.Consume(client, evtsCfg.ConsumerGroup, searchprovider.ListenEvents...)
	if err != nil {
		return nil, err
	}

	indexDir := filepath.Join(cfg.Datapath, "index.bleve")
	bleveIndex, err := bleve.Open(indexDir)
	if err != nil {
		mapping, err := index.BuildMapping()
		if err != nil {
			return nil, err
		}
		bleveIndex, err = bleve.New(indexDir, mapping)
		if err != nil {
			return nil, err
		}
	}
	index, err := index.New(bleveIndex)
	if err != nil {
		return nil, err
	}

	gwclient, err := pool.GetGatewayServiceClient(cfg.Reva.Address)
	if err != nil {
		logger.Fatal().Err(err).Str("addr", cfg.Reva.Address).Msg("could not get reva client")
	}

	provider := searchprovider.New(gwclient, index, cfg.MachineAuthAPIKey, evts, logger)

	return &Service{
		id:       cfg.GRPC.Namespace + "." + cfg.Service.Name,
		log:      logger,
		Config:   cfg,
		provider: provider,
	}, nil
}

// Service implements the searchServiceHandler interface
type Service struct {
	id       string
	log      log.Logger
	Config   *config.Config
	provider search.ProviderClient
}

func (s Service) Search(ctx context.Context, in *searchsvc.SearchRequest, out *searchsvc.SearchResponse) error {
	// Get token from the context (go-micro) and make it known to the reva client too (grpc)
	t, ok := metadata.Get(ctx, revactx.TokenHeader)
	if !ok {
		s.log.Error().Msg("Could not get token from context")
		return errors.New("could not get token from context")
	}
	ctx = grpcmetadata.AppendToOutgoingContext(ctx, revactx.TokenHeader, t)

	res, err := s.provider.Search(ctx, &searchsvc.SearchRequest{
		Query:    in.Query,
		PageSize: in.PageSize,
	})
	if err != nil {
		switch err.(type) {
		case errtypes.BadRequest:
			return merrors.BadRequest(s.id, err.Error())
		default:
			return merrors.InternalServerError(s.id, err.Error())
		}
	}

	out.Matches = res.Matches
	out.NextPageToken = res.NextPageToken
	return nil
}

func (s Service) IndexSpace(ctx context.Context, in *searchsvc.IndexSpaceRequest, out *searchsvc.IndexSpaceResponse) error {
	_, err := s.provider.IndexSpace(ctx, in)
	return err
}
