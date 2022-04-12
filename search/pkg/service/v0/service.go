package service

import (
	"github.com/blevesearch/bleve/v2"
	"github.com/cs3org/reva/v2/pkg/rgrpc/todo/pool"

	"github.com/owncloud/ocis/ocis-pkg/log"
	"github.com/owncloud/ocis/search/pkg/config"
	"github.com/owncloud/ocis/search/pkg/search"
	"github.com/owncloud/ocis/search/pkg/search/index"
	searchprovider "github.com/owncloud/ocis/search/pkg/search/provider"
)

// userDefaultGID is the default integer representing the "users" group.
const userDefaultGID = 30000

// New returns a new instance of Service
func New(opts ...Option) (*Service, error) {
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
		logger.Fatal().Msgf("could not get reva client at address %s", cfg.Reva.Address)
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
