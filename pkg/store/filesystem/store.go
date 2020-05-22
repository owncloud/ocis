// Package store implements the go-micro store interface
package store

import (
	"os"
	"path"

	olog "github.com/owncloud/ocis-pkg/v2/log"
	"github.com/owncloud/ocis-settings/pkg/config"
	"github.com/owncloud/ocis-settings/pkg/settings"
)

var (
	// Name is the default name for the settings store
	Name        = "ocis-settings-store"
	managerName = "filesystem"
)

// Store interacts with the filesystem to manage settings information
type Store struct {
	mountPath string
	Logger    olog.Logger
}

// New creates a new store
func New(cfg *config.Config) settings.Manager {
	s := Store{}

	dest := path.Join(cfg.Storage.RootMountPath, Name)
	if _, err := os.Stat(dest); err != nil {
		s.Logger.Info().Msgf("creating container on %v", dest)
		err := os.MkdirAll(dest, 0700)
		if err != nil {
			s.Logger.Err(err).Msgf("providing container on %v", dest)
		}
	}

	s.mountPath = dest
	return &s
}

func init() {
	settings.Registry[managerName] = New
}
