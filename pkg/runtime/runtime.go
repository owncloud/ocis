package runtime

import (
	"fmt"
	golog "log"
	"net/rpc"
	"time"

	"github.com/micro/cli/v2"

	"github.com/micro/micro/v2/client/api"
	"github.com/micro/micro/v2/client/web"
	"github.com/micro/micro/v2/service/registry"

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
		"proxy",
		"settings",
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
		// "reva-storage-root",
		"reva-storage-home",
		"reva-storage-public-link",
		"reva-storage-public-link-data",
		"reva-storage-home-data",
		"reva-storage-eos",
		"reva-storage-eos-data",
		"reva-storage-oc",
		"reva-storage-oc-data",
		"accounts",
		"glauth",
		"konnectd",
		"thumbnails",
	}

	// There seem to be a race condition when reva-sharing needs to read the sharing.json file and the parent folder is not present.
	dependants = []string{
		"reva-sharing",
	}

	// Maximum number of retries until getting a connection to the rpc runtime service.
	maxRetries int = 10
)

// Runtime represents an oCIS runtime environment.
type Runtime struct{}

// New creates a new ocis + micro runtime
func New() Runtime {
	return Runtime{}
}

// Start rpc runtime
func (r *Runtime) Start(services ...string) error {
	go r.Launch(services)
	return service.Start()
}

// Launch ocis default ocis extensions.
func (r *Runtime) Launch(services []string) {
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
			fmt.Println("runtime not available, retrying in 1 second...")
			time.Sleep(1 * time.Second)
		} else {
			goto OUT
		}
	}

OUT:
	for _, v := range services {
		args := process.NewProcEntry(v, []string{v}...)
		var reply int

		if err := client.Call("Service.Start", args, &reply); err != nil {
			golog.Fatal(err)
		}
	}

	// TODO(refs) this should disappear and tackled at the runtime (pman) level.
	// see https://github.com/cs3org/reva/issues/795 for race condition.
	// dependants might not be needed on a ocis_simple build, therefore
	// it should not be started under these circumstances.
	if len(services) >= len(Extensions) { // it will not run for ocis_simple builds.
		time.Sleep(2 * time.Second)
		for _, v := range dependants {
			args := process.NewProcEntry(v, []string{v}...)
			var reply int

			if err := client.Call("Service.Start", args, &reply); err != nil {
				golog.Fatal(err)
			}
		}
	}
}

// AddMicroPlatform adds the micro subcommands to the cli app
func AddMicroPlatform(app *cli.App) {
	setDefaults()

	app.Commands = append(app.Commands, api.Commands()...)
	app.Commands = append(app.Commands, web.Commands()...)
	app.Commands = append(app.Commands, registry.Commands()...)
}

// provide a config.Config with default values?
func setDefaults() {
	// api
	api.Name = OwncloudNamespace + "api"
	api.Namespace = OwncloudNamespace + "api"
	api.HeaderPrefix = "X-Micro-Owncloud-"

	// web
	web.Name = OwncloudNamespace + "web"
	web.Namespace = OwncloudNamespace + "web"

	// registry
	registry.Name = OwncloudNamespace + "registry"
}
