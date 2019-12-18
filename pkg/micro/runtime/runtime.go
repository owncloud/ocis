package runtime

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/micro/cli"
	gorun "github.com/micro/go-micro/runtime"
	"github.com/micro/micro/api"
	"github.com/micro/micro/broker"
	"github.com/micro/micro/health"
	"github.com/micro/micro/monitor"
	"github.com/micro/micro/network"
	"github.com/micro/micro/plugin/build"
	"github.com/micro/micro/proxy"
	"github.com/micro/micro/registry"
	"github.com/micro/micro/router"
	"github.com/micro/micro/runtime"
	"github.com/micro/micro/server"
	"github.com/micro/micro/service"
	"github.com/micro/micro/store"
	"github.com/micro/micro/token"
	"github.com/micro/micro/tunnel"
	"github.com/micro/micro/web"
	"github.com/owncloud/ocis-pkg/log"
)

// OwncloudNamespace is the base path for micro' services to use
var OwncloudNamespace = "com.owncloud."

// RuntimeServices to start as part of the fullstack option
var RuntimeServices = []string{
	// "network",  // :8085
	"runtime",  // :8088
	"registry", // :8000
	"broker",   // :8001
	"store",    // :8002
	"tunnel",   // :8083
	"router",   // :8084
	"proxy",    // :8081
	"api",      // :8080
	"web",      // :8082
	"monitor",  // :????
	// "debug",    // :????
	// "bot",      // :????
}

// Extensions are ocis extension services
var Extensions = []string{
	"hello",
	"phoenix",
	"graph",
	"ocs",
	"webdav",
	"reva",
}

// Runtime is a micro' runtime
type Runtime struct {
	Services []string
	Logger   log.Logger
	R        *gorun.Runtime
}

// New creates a new ocis + micro runtime
func New(opts ...Option) Runtime {
	options := newOptions(opts...)

	return Runtime{
		Services: options.Services,
		Logger:   options.Logger,
		R:        options.MicroRuntime,
	}
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

	r.Logger.Info().Msgf("Service runtime shutdown")
	os.Exit(0)
}

// Start starts preconfigured services
func (r *Runtime) Start() {
	env := os.Environ()

	for _, service := range r.Services {
		args := []gorun.CreateOption{
			// the binary calls itself with the micro service as a subcommand as first argument
			gorun.WithCommand(os.Args[0], service),
			gorun.WithEnv(env),
			// and logs to STDOUT. Perhaps this can be overridden to use a log.Logger
			gorun.WithOutput(os.Stdout),
		}

		muService := &gorun.Service{Name: service}
		if err := (*r.R).Create(muService, args...); err != nil {
			r.Logger.Error().Msgf("Failed to create runtime enviroment: %v", err)
		}
	}
}

// AddRuntime adds the micro subcommands to the cli app
func AddRuntime(app *cli.App) {
	// load this from disk / env?
	setDefaults()

	app.Commands = append(app.Commands, api.Commands()...)
	app.Commands = append(app.Commands, broker.Commands()...)
	app.Commands = append(app.Commands, health.Commands()...)
	app.Commands = append(app.Commands, proxy.Commands()...)
	app.Commands = append(app.Commands, monitor.Commands()...)
	app.Commands = append(app.Commands, router.Commands()...)
	app.Commands = append(app.Commands, tunnel.Commands()...)
	app.Commands = append(app.Commands, network.Commands()...)
	app.Commands = append(app.Commands, registry.Commands()...)
	app.Commands = append(app.Commands, runtime.Commands()...)
	app.Commands = append(app.Commands, server.Commands()...)
	app.Commands = append(app.Commands, service.Commands()...)
	app.Commands = append(app.Commands, store.Commands()...)
	app.Commands = append(app.Commands, token.Commands()...)
	app.Commands = append(app.Commands, build.Commands()...)
	app.Commands = append(app.Commands, web.Commands()...)
	// app.Commands = append(app.Commands, bot.Commands()...)
	// app.Commands = append(app.Commands, new.Commands()...)
	// app.Commands = append(app.Commands, debug.Commands()...)
}

// provide a config.Config with default values?
func setDefaults() {
	// api
	api.Name = OwncloudNamespace + "api"
	api.Namespace = OwncloudNamespace + "api"
	api.HeaderPrefix = "X-Micro-Owncloud-"
	// api.EnableRPC = false
	// api.RPCPath
	// api.Address
	// api.Handler
	// api.Resolver
	// api.ProxyPath
	// api.APIPath

	// broker
	broker.Name = OwncloudNamespace + "http.broker"
	// broker.Address

	// proxy
	proxy.Name = OwncloudNamespace + "proxy"
	// proxy.Address
	// proxy.Protocol
	// proxy.Endpoint

	// monitor
	monitor.Name = OwncloudNamespace + "monitor"

	// router
	router.Name = OwncloudNamespace + "router"
	// router.Address
	// router.Network
	// router.Topic

	// tunnel
	tunnel.Name = OwncloudNamespace + "tunnel"
	// tunnel.Address
	// tunnel.Token
	// tunnel.Tunnel

	// network
	network.Name = OwncloudNamespace + "network"
	// network.Address
	// network.Advertise
	// network.Resolver
	// network.Token

	// registry
	registry.Name = OwncloudNamespace + "registry"
	// registry.Address
	// registry.SyncTime
	// registry.Topic

	// runtime
	runtime.Name = OwncloudNamespace + "runtime"
	// runtime.Address

	// server
	server.Name = OwncloudNamespace + "server"
	// server.Address
	// server.Network
	// server.Router

	// store
	store.Name = OwncloudNamespace + "store"
	// store.Address

	// web
	web.Name = OwncloudNamespace + "web"
	web.Namespace = OwncloudNamespace + "web"
	// web.ACMECA
	// web.ACMEChallengeProvider
	// web.ACMEProvider
	// web.Address
	// web.BasePathHeader
	// web.DefaultIcon

}
