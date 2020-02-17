package runtime

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/micro/cli/v2"
	gorun "github.com/micro/go-micro/v2/runtime"
	"github.com/micro/micro/v2/api"
	"github.com/micro/micro/v2/proxy"
	"github.com/micro/micro/v2/registry"
	"github.com/micro/micro/v2/runtime"
	"github.com/micro/micro/v2/web"
	"github.com/owncloud/ocis-pkg/v2/log"
)

// OwncloudNamespace is the base path for micro' services to use
var OwncloudNamespace = "com.owncloud."

// RuntimeServices to start as part of the fullstack option
var RuntimeServices = []string{
	"api",      // :8080
	"proxy",    // :8081
	"web",      // :8082
	"registry", // :8000
	"runtime",  // :8088 (future proof. We want to be able to control extensions through a runtime)
}

// Extensions are ocis extension services
var Extensions = []string{
	"hello",
	"phoenix",
	"graph",
	"graph-explorer",
	"ocs",
	"webdav",
	"reva-frontend",
	"reva-gateway",
	"reva-users",
	"reva-auth-basic",
	"reva-auth-bearer",
	"reva-sharing",
	"reva-storage-root",
	"reva-storage-home",
	"reva-storage-home-data",
	"reva-storage-oc",
	"reva-storage-oc-data",
	"devldap",
	"konnectd",
}

// Runtime is a wrapper around micro's own runtime
type Runtime struct {
	Logger log.Logger
	R      *gorun.Runtime

	services []*gorun.Service
}

// New creates a new ocis + micro runtime
func New(opts ...Option) Runtime {
	options := newOptions(opts...)

	r := Runtime{
		Logger: options.Logger,
		R:      options.MicroRuntime,
	}

	for _, v := range append(RuntimeServices, Extensions...) {
		r.services = append(r.services, &gorun.Service{Name: v})
	}

	return r
}

// Trap waits for a sigkill to stop the runtime
func (r *Runtime) Trap() {
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	if err := (*r.R).Start(); err != nil {
		os.Exit(1)
	}

	// block until there is a value
	for range shutdown {
		r.Logger.Info().Msg("shutdown signal received")
		close(shutdown)
	}

	if err := (*r.R).Stop(); err != nil {
		r.Logger.Err(err)
	}

	for _, v := range r.services {
		r.Logger.Info().Msgf("gracefully stopping service %v", v.Name)
		(*r.R).Delete(v)
	}

	os.Exit(0)
}

// Start starts preconfigured services
func (r *Runtime) Start() {
	env := os.Environ()

	for _, service := range r.services {
		r.Logger.Info().Msgf("args: %v %v", os.Args[0], service.Name) // TODO uncommenting this line causes some issues where the binary calls itself with the `server` as argument
		args := []gorun.CreateOption{
			gorun.WithCommand(os.Args[0], service.Name),
			gorun.WithEnv(env),
			gorun.WithOutput(os.Stdout),
		}

		if err := (*r.R).Create(service, args...); err != nil {
			r.Logger.Error().Msgf("Failed to create runtime enviroment: %v", err)
		}
	}
}

// AddRuntime adds the micro subcommands to the cli app
func AddRuntime(app *cli.App) {
	setDefaults()

	app.Commands = append(app.Commands, api.Commands()...)
	app.Commands = append(app.Commands, proxy.Commands()...)
	app.Commands = append(app.Commands, web.Commands()...)
	app.Commands = append(app.Commands, registry.Commands()...)
	app.Commands = append(app.Commands, runtime.Commands()...)
}

// provide a config.Config with default values?
func setDefaults() {
	// api
	api.Name = OwncloudNamespace + "api"
	api.Namespace = OwncloudNamespace + "api"
	api.HeaderPrefix = "X-Micro-Owncloud-"

	// proxy
	proxy.Name = OwncloudNamespace + "proxy"

	// web
	web.Name = OwncloudNamespace + "web"
	web.Namespace = OwncloudNamespace + "web"

	// registry
	registry.Name = OwncloudNamespace + "registry"

	// runtime
	runtime.Name = OwncloudNamespace + "runtime"
}
