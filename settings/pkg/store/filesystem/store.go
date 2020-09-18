// Package store implements the go-micro store interface
package store

import (
	"os"

	olog "github.com/owncloud/ocis-pkg/v2/log"
	"github.com/owncloud/ocis/settings/pkg/config"
	"github.com/owncloud/ocis/settings/pkg/settings"
)

var (
	// Name is the default name for the settings store
	Name        = "ocis-settings"
	managerName = "filesystem"
)

// Store interacts with the filesystem to manage settings information
type Store struct {
	dataPath string
	Logger   olog.Logger
}

// New creates a new store
func New(cfg *config.Config) settings.Manager {
	s := Store{
		Logger: olog.NewLogger(
			olog.Color(cfg.Log.Color),
			olog.Pretty(cfg.Log.Pretty),
			olog.Level(cfg.Log.Level),
		),
	}

	if _, err := os.Stat(cfg.Storage.DataPath); err != nil {
		s.Logger.Info().Msgf("creating container on %v", cfg.Storage.DataPath)
		err := os.MkdirAll(cfg.Storage.DataPath, 0700)

		if err != nil {
			s.Logger.Err(err).Msgf("providing container on %v", cfg.Storage.DataPath)
		}
	}

	s.dataPath = cfg.Storage.DataPath
	return &s
}

func init() {
	settings.Registry[managerName] = New
}
