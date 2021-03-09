package runtime

import (
	"github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis/pkg/runtime/service"
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
