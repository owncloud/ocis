package runtime

import (
	"os"

	mzlog "github.com/asim/go-micro/plugins/logger/zerolog/v3"
	"github.com/asim/go-micro/v3/logger"
	"github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis/pkg/runtime/service"
	"github.com/rs/zerolog"
	"github.com/thejerf/suture"
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
// TODO(refs) initialize a new Service struct
// - run the set of default services under the server command.
func (r *Runtime) Start() error {
	// 1. create new Service struct
	// 2. start

	return service.Start(service.WithConfig(r.c))

	//setMicroLogger(r.c.Log)
	//halt := make(chan os.Signal, 1)
	//signal.Notify(halt, os.Interrupt)
	//
	//supervisor := suture.NewSimple("ocis")
	//globalCtx, globalCancel := context.WithCancel(context.Background())
	//
	//// propagate reva log config to storage services
	//inheritedOptions := []storage.Option{
	//	storage.WithLogPretty(r.c.Log.Pretty),
	//	storage.WithLogColor(r.c.Log.Color),
	//	storage.WithLogLevel(r.c.Log.Level),
	//}
	//
	//addServiceToken("settings", supervisor.Add(settings.NewSutureService(globalCtx, r.c.Settings)))
	//addServiceToken("storage-metadata", supervisor.Add(storage.NewStorageMetadata(globalCtx, inheritedOptions...)))
	//addServiceToken("accounts", supervisor.Add(accounts.NewSutureService(globalCtx, r.c.Accounts)))
	//addServiceToken("glauth", supervisor.Add(glauth.NewSutureService(globalCtx, r.c.GLAuth)))
	//addServiceToken("idp", supervisor.Add(idp.NewSutureService(globalCtx, r.c.IDP)))
	//addServiceToken("ocs", supervisor.Add(ocs.NewSutureService(globalCtx, r.c.OCS)))
	//addServiceToken("onlyoffice", supervisor.Add(onlyoffice.NewSutureService(globalCtx, r.c.Onlyoffice)))
	//addServiceToken("proxy", supervisor.Add(proxy.NewSutureService(globalCtx, r.c.Proxy)))
	//addServiceToken("store", supervisor.Add(store.NewSutureService(globalCtx, r.c.Store)))
	//addServiceToken("thumbnails", supervisor.Add(thumbnails.NewSutureService(globalCtx, r.c.Thumbnails)))
	//addServiceToken("web", supervisor.Add(web.NewSutureService(globalCtx, r.c.Web)))
	//addServiceToken("webdav", supervisor.Add(webdav.NewSutureService(globalCtx, r.c.WebDAV)))
	//addServiceToken("storage-frontend", supervisor.Add(storage.NewFrontend(globalCtx, inheritedOptions...)))
	//addServiceToken("storage-gateway", supervisor.Add(storage.NewGateway(globalCtx, inheritedOptions...)))
	//addServiceToken("storage-users", supervisor.Add(storage.NewUsersProviderService(globalCtx, inheritedOptions...)))
	//addServiceToken("storage-groupsprovider", supervisor.Add(storage.NewGroupsProvider(globalCtx, inheritedOptions...)))
	//addServiceToken("storage-authbasic", supervisor.Add(storage.NewAuthBasic(globalCtx, inheritedOptions...)))
	//addServiceToken("storage-authbearer", supervisor.Add(storage.NewAuthBearer(globalCtx, inheritedOptions...)))
	//addServiceToken("storage-home", supervisor.Add(storage.NewStorageHome(globalCtx, inheritedOptions...)))
	//addServiceToken("storage-users", supervisor.Add(storage.NewStorageUsers(globalCtx, inheritedOptions...)))
	//addServiceToken("storage-public-link", supervisor.Add(storage.NewStoragePublicLink(globalCtx, inheritedOptions...)))
	//addServiceToken("storage-sharing", supervisor.Add(storage.NewSharing(globalCtx, inheritedOptions...)))
	//
	//// TODO(refs) debug line with supervised services.
	//go supervisor.ServeBackground()
	//
	//<-halt
	//
	//globalCancel()
	//close(halt)
	return nil
}

// addServiceToken adds a service token to a global slice of service tokens that contains services managed by the supervisor.
func addServiceToken(service string, token suture.ServiceToken) {
	tokens[service] = append(tokens[service], token)
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
