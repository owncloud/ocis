// Package store implements the go-micro store interface
package store

import (
	"context"
	"os"

	olog "github.com/owncloud/ocis/ocis-pkg/log"
	"github.com/owncloud/ocis/settings/pkg/config"
	"github.com/owncloud/ocis/settings/pkg/settings"
)

var (
	// Name is the default name for the settings store
	Name        = "ocis-settings"
	managerName = "metadata"
)

// MetadataClient is the interface to talk to metadata service
type MetadataClient interface {
	SimpleDownload(ctx context.Context, id string) ([]byte, error)
	SimpleUpload(ctx context.Context, id string, content []byte) error
	Delete(ctx context.Context, id string) error
	ReadDir(ctx context.Context, id string) ([]string, error)
}

// Store interacts with the filesystem to manage settings information
type Store struct {
	Logger olog.Logger

	mdc MetadataClient
}

// New creates a new store
func New(cfg *config.Config) settings.Manager {
	s := Store{
		//Logger: olog.NewLogger(
		//	olog.Color(cfg.Log.Color),
		//	olog.Pretty(cfg.Log.Pretty),
		//	olog.Level(cfg.Log.Level),
		//	olog.File(cfg.Log.File),
		//),
	}

	if _, err := os.Stat(cfg.DataPath); err != nil {
		s.Logger.Info().Msgf("creating container on %v", cfg.DataPath)
		err = os.MkdirAll(cfg.DataPath, 0700)

		if err != nil {
			s.Logger.Err(err).Msgf("providing container on %v", cfg.DataPath)
		}
	}

	return &s
}

func init() {
	settings.Registry[managerName] = New
}
