package runtime

import (
	"context"
	"fmt"
	golog "log"
	"net/rpc"
	"os"
	"os/signal"
	"time"

	"github.com/thejerf/suture"

	"github.com/rs/zerolog"

	mzlog "github.com/asim/go-micro/plugins/logger/zerolog/v3"
	"github.com/asim/go-micro/v3/logger"

	"github.com/owncloud/ocis/ocis/pkg/config"

	"github.com/owncloud/ocis/ocis/pkg/runtime/process"
	settings "github.com/owncloud/ocis/settings/pkg/command"
)

var (
	// OwncloudNamespace is the base path for micro' services to use
	OwncloudNamespace = "com.owncloud."

	// MicroServices to start as part of the fullstack option
	MicroServices = []string{
		"api",      // :8080
		"registry", // :8000
	}

	// Extensions are oCIS extension services
	Extensions = []string{
		"glauth",
		"idp",
		"ocs",
		"onlyoffice",
		"proxy",
		"settings",
		"store",
		"storage-metadata",
		"storage-frontend",
		"storage-gateway",
		"storage-userprovider",
		"storage-groupprovider",
		"storage-auth-basic",
		"storage-auth-bearer",
		"storage-home",
		"storage-users",
		"storage-public-link",
		"thumbnails",
		"web",
		"webdav",
	}

	dependants = []string{
		"storage-sharing",
		"accounts",
	}
	// Maximum number of retries until getting a connection to the rpc runtime service.
	maxRetries = 10
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

// serviceTokens keeps in memory a set of [service name] = []suture.ServiceToken that is used to shutdown services.
// Shutting down a service implies removing it from the supervisor AND cancelling its context, this should be done
// within the service Stop() method. Services should cancel their context.
type serviceTokens map[string][]suture.ServiceToken

// Start rpc runtime
func (r *Runtime) Start() error {
	setMicroLogger(r.c.Log)
	halt := make(chan os.Signal, 1)
	signal.Notify(halt, os.Interrupt)

	// tokens are used to keep track of the services
	tokens := serviceTokens{}
	supervisor := suture.NewSimple("ocis")
	globalCtx, globalCancel := context.WithCancel(context.Background())

	// TODO(refs + jfd)
	// - to avoid this getting out of hands, a supervisor would need to be injected on each supervised service.
	// - each service would then add its execute func to the supervisor, and return its token (?)
	// - this runtime should only care about start / stop services, for that we use serviceTokens.

	tokens["settings"] = append(tokens["settings"], supervisor.Add(settings.NewSutureService(globalCtx, r.c.Settings)))

	go supervisor.ServeBackground()

	<-halt
	globalCancel()
	close(halt)
	return nil
}

// for logging reasons we don't want the same logging level on both oCIS and micro. As a framework builder we do not
// want to expose to the end user the internal framework logs unless explicitly specified.a
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

// Launch oCIS default oCIS extensions.
func (r *Runtime) Launch() {
	var client *rpc.Client
	var err error
	var try int

	for {
		if try >= maxRetries {
			golog.Fatal("could not get a connection to rpc runtime on localhost:10666")
		}
		client, err = rpc.DialHTTP("tcp", "localhost:10666")
		if err != nil {
			try++
			fmt.Println("runtime not available, retrying...")
			time.Sleep(1 * time.Second)
		} else {
			goto OUT
		}
	}

OUT:
	for _, v := range Extensions {
		RunService(client, v)
	}

	if len(dependants) > 0 {
		time.Sleep(2 * time.Second)
		for _, v := range dependants {
			RunService(client, v)
		}
	}
}

// RunService sends a Service.Start command with the given service name  to pman
func RunService(client *rpc.Client, service string) int {
	args := process.NewProcEntry(service, os.Environ(), []string{service}...)

	all := append(Extensions, dependants...)
	if !contains(all, service) {
		return 1
	}

	var reply int
	if err := client.Call("Service.Start", args, &reply); err != nil {
		golog.Fatal(err)
	}
	return reply
}

func contains(a []string, b string) bool {
	for i := range a {
		if a[i] == b {
			return true
		}
	}
	return false
}
