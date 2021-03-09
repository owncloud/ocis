package runtime

import (
	"os"

	mzlog "github.com/asim/go-micro/plugins/logger/zerolog/v3"
	"github.com/asim/go-micro/v3/logger"
	"github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis/pkg/runtime/service"
	"github.com/rs/zerolog"
)

// Runtime represents an oCIS runtime environment.
type Runtime struct {
	c *config.Config
}

// New creates a new oCIS + micro runtime
func New(cfg *config.Config) Runtime {
	return Runtime{
		c: cfg,
	}
}

// Start rpc runtime
func (r *Runtime) Start() error {
	return service.Start(service.WithConfig(r.c))
}

// for logging reasons we don't want the same logging level on both oCIS and micro. As a framework builder we do not
// want to expose to the end user the internal framework logs unless explicitly specified.
func setMicroLogger(log config.Log) {
	if os.Getenv("MICRO_LOG_LEVEL") == "" {
		os.Setenv("MICRO_LOG_LEVEL", "error")
	}

	lev, err := zerolog.ParseLevel(os.Getenv("MICRO_LOG_LEVEL"))
	if err != nil {
		lev = zerolog.ErrorLevel
	}
	logger.DefaultLogger = mzlog.NewLogger(logger.WithLevel(logger.Level(lev)))
}
