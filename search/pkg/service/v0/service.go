package service

import (
	"time"

	"github.com/pkg/errors"

	"github.com/owncloud/ocis/ocis-pkg/service/grpc"

	"github.com/owncloud/ocis/ocis-pkg/indexer"

	"github.com/owncloud/ocis/ocis-pkg/log"
	oreg "github.com/owncloud/ocis/ocis-pkg/registry"
	"github.com/owncloud/ocis/ocis-pkg/roles"
	settingssvc "github.com/owncloud/ocis/protogen/gen/ocis/services/settings/v0"
	"github.com/owncloud/ocis/search/pkg/config"
)

// userDefaultGID is the default integer representing the "users" group.
const userDefaultGID = 30000

// New returns a new instance of Service
func New(opts ...Option) (s *Service, err error) {
	options := newOptions(opts...)
	logger := options.Logger
	cfg := options.Config

	roleService := options.RoleService
	if roleService == nil {
		roleService = settingssvc.NewRoleService("com.owncloud.api.settings", grpc.DefaultClient)
	}
	roleManager := options.RoleManager
	if roleManager == nil {
		m := roles.NewManager(
			roles.CacheSize(1024),
			roles.CacheTTL(time.Hour*24*7),
			roles.Logger(options.Logger),
			roles.RoleService(roleService),
		)
		roleManager = &m
	}

	storage, err := createMetadataStorage(cfg, logger)
	if err != nil {
		return nil, errors.Wrap(err, "could not create metadata storage")
	}

	s = &Service{
		id:     cfg.GRPC.Namespace + "." + cfg.Service.Name,
		log:    logger,
		Config: cfg,
	}

	r := oreg.GetRegistry()
	if cfg.Repo.Backend == "cs3" {
		if _, err := r.GetService("com.owncloud.storage.metadata"); err != nil {
			logger.Error().Err(err).Msg("index: storage-metadata service not present")
			return nil, err
		}
	}

	return
}

// Service implements the searchServiceHandler interface
type Service struct {
	id     string
	log    log.Logger
	Config *config.Config
	index  *indexer.Indexer
}
