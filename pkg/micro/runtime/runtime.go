package runtime

import (
	"fmt"
	golog "log"
	"net/rpc"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/micro/cli/v2"
	gorun "github.com/micro/go-micro/v2/runtime"
	"github.com/micro/micro/v2/api"
	"github.com/micro/micro/v2/proxy"
	"github.com/micro/micro/v2/registry"
	"github.com/micro/micro/v2/runtime"
	"github.com/micro/micro/v2/web"
	"github.com/owncloud/ocis-pkg/v2/log"

	"github.com/refs/pman/pkg/process"
	"github.com/refs/pman/pkg/service"
)

var (
	// OwncloudNamespace is the base path for micro' services to use
	OwncloudNamespace = "com.owncloud."

	// MicroServices to start as part of the fullstack option
	MicroServices = []string{
		"api",      // :8080
		"web",      // :8082
		"registry", // :8000
	}

	// Extensions are ocis extension services
	Extensions = []string{
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
		//"reva-storage-root",
		"reva-storage-home",
		"reva-storage-home-data",
		"reva-storage-eos",
		"reva-storage-eos-data",
		"reva-storage-oc",
		"reva-storage-oc-data",
		"accounts",
		"glauth",
		"konnectd",
		"proxy",
		"thumbnails",
	}
)

// Runtime is a wrapper around micro's own runtime
type Runtime struct {
	Logger log.Logger
	R      *gorun.Runtime
	Ctx    *cli.Context

	services []*gorun.Service
}

// New creates a new ocis + micro runtime
func New(opts ...Option) Runtime {
	options := newOptions(opts...)

	r := Runtime{
		Logger: options.Logger,
		R:      options.MicroRuntime,
		Ctx:    options.Context,
	}

	for _, v := range append(MicroServices, Extensions...) {
		r.services = append(r.services, &gorun.Service{Name: v})
	}

	return r
}

// Trap listen and blocks for termination signals
func (r Runtime) Trap() {
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	if err := (*r.R).Start(); err != nil {
		os.Exit(1)
	}

	for range shutdown {
		r.Logger.Info().Msg("shutdown signal received")
		close(shutdown)
	}

	if err := (*r.R).Stop(); err != nil {
		r.Logger.Err(err).Msgf("error while shutting down")
	}

	for _, s := range r.services {
		r.Logger.Info().Msgf("gracefully stopping service %v", s.Name)
		if err := (*r.R).Delete(s); err != nil {
			r.Logger.Err(err).Msgf("error while deleting service: %v", s.Name)
		}
	}

	os.Exit(0)
}

// Start rpc runtime
func (r *Runtime) Start() {
	go r.Launch()
	service.Start()
}

// Launch ocis Extensions
func (r *Runtime) Launch() {
	client, err := rpc.DialHTTP("tcp", "localhost:10666")
	if err != nil {
		// ensure the rpc service is running before attempting to start any extension
		fmt.Println("rpc service not available, retrying in 1 second...")
		time.Sleep(1 * time.Second)
		r.Launch()
	}

	// loop over extensions starting them
	for i := range Extensions {
		arg0 := process.NewProcEntry(
			Extensions[i],
			[]string{Extensions[i]}...,
		)
		var arg1 int

		if err := client.Call("Service.Start", arg0, &arg1); err != nil {
			golog.Fatal(err)
		}
	}
}

// AddMicroPlatform adds the micro subcommands to the cli app
func AddMicroPlatform(app *cli.App) {
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
