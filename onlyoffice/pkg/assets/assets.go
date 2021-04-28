package assets

import (
	"github.com/owncloud/ocis/onlyoffice"
	"net/http"
	"os"
	"path"

	"github.com/owncloud/ocis/ocis-pkg/log"
	"github.com/owncloud/ocis/onlyoffice/pkg/config"
)

// assets gets initialized by New and provides the handler.
type assets struct {
	logger log.Logger
	config *config.Config
}

// Open just implements the HTTP filesystem interface.
func (a assets) Open(original string) (http.File, error) {
	if a.config.Asset.Path != "" {
		if stat, err := os.Stat(a.config.Asset.Path); err == nil && stat.IsDir() {
			custom := path.Join(
				a.config.Asset.Path,
				original,
			)

			if _, err := os.Stat(custom); !os.IsNotExist(err) {
				f, err := os.Open(custom)

				if err != nil {
					return nil, err
				}

				return f, nil
			}
		} else {
			a.logger.Warn().
				Str("path", a.config.Asset.Path).
				Msg("Assets directory doesn't exist")
		}
	}

	return onlyoffice.Assets.Open(original)
}

// New returns a new http filesystem to serve assets.
func New(opts ...Option) http.FileSystem {
	options := newOptions(opts...)

	return assets{
		config: options.Config,
		logger: options.Logger,
	}
}
