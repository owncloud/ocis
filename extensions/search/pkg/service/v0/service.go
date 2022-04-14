package service

import (
	"context"

	"github.com/blevesearch/bleve/v2"
	"github.com/cs3org/reva/v2/pkg/rgrpc/todo/pool"

	"github.com/owncloud/ocis/extensions/search/pkg/config"
	"github.com/owncloud/ocis/extensions/search/pkg/search"
	"github.com/owncloud/ocis/extensions/search/pkg/search/index"
	searchprovider "github.com/owncloud/ocis/extensions/search/pkg/search/provider"
	"github.com/owncloud/ocis/ocis-pkg/log"
	searchsvc "github.com/owncloud/ocis/protogen/gen/ocis/services/search/v0"
)

// NewHandler returns a service implementation for Service.
func NewHandler(opts ...Option) (searchsvc.SearchProviderHandler, error) {
	options := newOptions(opts...)
	logger := options.Logger
	cfg := options.Config

	bleveIndex, err := bleve.NewMemOnly(index.BuildMapping())
	if err != nil {
		return nil, err
	}
	index, err := index.New(bleveIndex)
	if err != nil {
		return nil, err
	}

	gwclient, err := pool.GetGatewayServiceClient(cfg.Reva.Address)
	if err != nil {
		logger.Fatal().Err(err).Str("addr", cfg.Reva.Address).Msg("could not get reva client")
	}

	provider := searchprovider.New(gwclient, index)

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
	res, err := s.provider.Search(ctx, &searchsvc.SearchRequest{
		Query: in.Query,
	})
	if err != nil {
		return nil
	}

	out.Matches = res.Matches
	out.NextPageToken = res.NextPageToken
	return nil
}
