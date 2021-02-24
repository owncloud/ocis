package assets

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/owncloud/ocis/ocis-pkg/log"
	"github.com/owncloud/ocis/settings/pkg/config"

	// Fake the import to make the dep tree happy.
	_ "golang.org/x/net/context"

	// Fake the import to make the dep tree happy.
	_ "golang.org/x/net/webdav"
)

//go:generate make -C ../.. embed.yml

// assets gets initialized by New and provides the handler.
type assets struct {
	logger log.Logger
	config *config.Config
}

// Open just implements the HTTP filesystem interface.
func (a assets) Open(original string) (http.File, error) {
	if a.config.Asset.Path != "" {
		if stat, err := os.Stat(a.config.Asset.Path); err == nil && stat.IsDir() {
			custom := filepath.Join(
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

	return FS.OpenFile(
		CTX,
		original,
		os.O_RDONLY,
		0644,
	)
}

// New returns a new http filesystem to serve assets.
func New(opts ...Option) http.FileSystem {
	options := newOptions(opts...)

	return assets{
		config: options.Config,
	}
}
