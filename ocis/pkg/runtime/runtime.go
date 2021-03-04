package runtime

import (
	"context"
	"os"
	"os/signal"

	accounts "github.com/owncloud/ocis/accounts/pkg/command"
	settings "github.com/owncloud/ocis/settings/pkg/command"

	"github.com/thejerf/suture"

	"github.com/rs/zerolog"

	mzlog "github.com/asim/go-micro/plugins/logger/zerolog/v3"
	"github.com/asim/go-micro/v3/logger"

	"github.com/owncloud/ocis/ocis/pkg/config"

	storage "github.com/owncloud/ocis/storage/pkg/command"
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
		"settings", // done
		"store",
		"storage-metadata", // done
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
		"accounts", // done
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

// tokens are used to keep track of the services
var tokens = serviceTokens{}

// Start rpc runtime
func (r *Runtime) Start() error {
	setMicroLogger(r.c.Log)
	halt := make(chan os.Signal, 1)
	signal.Notify(halt, os.Interrupt)

	supervisor := suture.NewSimple("ocis")
	globalCtx, globalCancel := context.WithCancel(context.Background())

	// TODO(refs + jfd)
	// - to avoid this getting out of hands, a supervisor would need to be injected on each supervised service.
	// - each service would then add its execute func to the supervisor, and return its token (?)
	// - this runtime should only care about start / stop services, for that we use serviceTokens.
	// - normalize the use of panics so that suture can restart services that die.
	// - shutting down a service implies iterating over all the serviceToken for the given service name and terminating them.
	// - config file parsing with Viper is no longer possible as viper is not thread-safe (https://github.com/spf13/viper/issues/19)
	// - replace occurrences of log.Fatal in favor of panic() since the supervisor relies on panics.
	// - the runtime should ideally run as an rpc service one can do requests, like the good ol' pman, rest in pieces.
	// - establish on suture a max number of retries before all initialization comes to a halt.

	// propagate reva log config to storage services
	r.c.Storage.Log.Level = r.c.Log.Level
	r.c.Storage.Log.Color = r.c.Log.Color
	r.c.Storage.Log.Pretty = r.c.Log.Pretty

	addServiceToken("settings", supervisor.Add(settings.NewSutureService(globalCtx, r.c.Settings)))
	addServiceToken("storagemetadata", supervisor.Add(storage.NewStorageMetadata(globalCtx, r.c.Storage)))
	addServiceToken("accounts", supervisor.Add(accounts.NewSutureService(globalCtx, r.c.Accounts)))

	go supervisor.ServeBackground()

	select {
	case <-halt:
		globalCancel()
		close(halt)
		return nil
	}
}

// addServiceToken adds a service token to a global slice of service tokens that contains services managed by the supervisor.
func addServiceToken(service string, token suture.ServiceToken) {
	tokens[service] = append(tokens[service], token)
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
