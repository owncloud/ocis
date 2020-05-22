package runtime

import (
	"fmt"
	golog "log"
	"net/rpc"
	"time"

	"github.com/micro/cli/v2"

	"github.com/micro/micro/v2/client/api"
	"github.com/micro/micro/v2/client/proxy"
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
		"thumbnails",
	}
)

// Runtime represents an oCIS runtime environment.
type Runtime struct{}

// New creates a new ocis + micro runtime
func New() Runtime {
	return Runtime{}
}

// Start rpc runtime
func (r *Runtime) Start() error {
	go r.Launch()
	return service.Start()
}

// Launch ocis Extensions
func (r *Runtime) Launch() {
	client, err := rpc.DialHTTP("tcp", "localhost:10666")
	if err != nil {
		fmt.Println("rpc service not available, retrying in 1 second...")
		time.Sleep(1 * time.Second)
		r.Launch()
	}

	all := append(Extensions, MicroServices...)
	for i := range all {
		arg0 := process.NewProcEntry(
			all[i],
			[]string{all[i]}...,
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
}
