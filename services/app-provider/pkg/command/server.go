package command

import (
	"context"
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/cs3org/reva/v2/cmd/revad/runtime"
	"github.com/google/uuid"
	"github.com/oklog/run"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/configlog"
	"github.com/owncloud/ocis/v2/ocis-pkg/registry"
	"github.com/owncloud/ocis/v2/ocis-pkg/sync"
	"github.com/owncloud/ocis/v2/ocis-pkg/tracing"
	"github.com/owncloud/ocis/v2/ocis-pkg/version"
	"github.com/owncloud/ocis/v2/services/app-provider/pkg/config"
	"github.com/owncloud/ocis/v2/services/app-provider/pkg/config/parser"
	"github.com/owncloud/ocis/v2/services/app-provider/pkg/logging"
	"github.com/owncloud/ocis/v2/services/app-provider/pkg/revaconfig"
	"github.com/owncloud/ocis/v2/services/app-provider/pkg/server/debug"
	"github.com/urfave/cli/v2"

	registrypb "github.com/cs3org/go-cs3apis/cs3/app/registry/v1beta1"
	oreg "github.com/owncloud/ocis/v2/ocis-pkg/registry"
	mreg "go-micro.dev/v4/registry"
)

// Server is the entry point for the server command.
func Server(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "server",
		Usage:    fmt.Sprintf("start the %s service without runtime (unsupervised mode)", cfg.Service.Name),
		Category: "server",
		Before: func(c *cli.Context) error {
			return configlog.ReturnFatal(parser.ParseConfig(cfg))
		},
		Action: func(c *cli.Context) error {
			logger := logging.Configure(cfg.Service.Name, cfg.Log)
			tracingProvider, err := tracing.GetServiceTraceProvider(cfg.Tracing, cfg.Service.Name)
			if err != nil {
				return err
			}
			gr := run.Group{}
			ctx, cancel := defineContext(cfg)

			defer cancel()

			gr.Add(func() error {
				pidFile := path.Join(os.TempDir(), "revad-"+cfg.Service.Name+"-"+uuid.New().String()+".pid")
				rCfg := revaconfig.AppProviderConfigFromStruct(cfg)
				reg := registry.GetRegistry()

				runtime.RunWithOptions(rCfg, pidFile,
					runtime.WithLogger(&logger.Logger),
					runtime.WithRegistry(reg),
					runtime.WithTraceProvider(tracingProvider),
				)

				return nil
			}, func(err error) {
				logger.Error().
					Str("server", cfg.Service.Name).
					Err(err).
					Msg("Shutting down server")

				cancel()
			})

			debugServer, err := debug.Server(
				debug.Logger(logger),
				debug.Context(ctx),
				debug.Config(cfg),
			)
			if err != nil {
				logger.Info().Err(err).Str("server", "debug").Msg("Failed to initialize server")
				return err
			}

			gr.Add(debugServer.ListenAndServe, func(_ error) {
				cancel()
			})

			if !cfg.Supervised {
				sync.Trap(&gr, cancel)
			}

			grpcSvc := registry.BuildGRPCService(
				cfg.GRPC.Namespace+"."+cfg.Service.Name,
				uuid.New().String(),
				cfg.GRPC.Addr,
				version.GetString(),
				nil,
			)
			updateRegistryNode(grpcSvc.Nodes[0],
				"cs3", //TODO: make configurable
				cfg.Drivers.WOPI.AppAddress,
				cfg.Drivers.WOPI.AppName,
				cfg.Drivers.WOPI.AppDescription,
				cfg.Drivers.WOPI.AppIconURI,
				cfg.Drivers.WOPI.AppPriority,
				cfg.Drivers.WOPI.AppCapabilities,
				cfg.Drivers.WOPI.AppDesktopOnly,
				cfg.Drivers.WOPI.AppSupportedMimeTypes,
			)

			if err := registry.RegisterService(ctx, grpcSvc, logger); err != nil {
				logger.Fatal().Err(err).Msg("failed to register the grpc service")
			}

			return gr.Run()
		},
	}
}

func updateRegistryNode(node *mreg.Node, ns, address, name, desc, icon, prio string, cap int32, desktopOnly bool, mimeTypes []string) {
	// TODO: the string "app-provider" needs to be supplied through cfg.Service.Name
	//node.Address = address
	//node.Id = ns + "api.app-provider"
	node.Metadata[ns+".app-provider.mime_type"] = joinMimeTypes(mimeTypes)
	node.Metadata[ns+".app-provider.name"] = name
	node.Metadata[ns+".app-provider.description"] = desc
	node.Metadata[ns+".app-provider.icon"] = icon

	node.Metadata[ns+".app-provider.allow_creation"] = registrypb.ProviderInfo_Capability_name[cap]
	node.Metadata[ns+".app-provider.priority"] = prio
	if desktopOnly {
		node.Metadata[ns+".app-provider.desktop_only"] = "true"
	}
}

// This is to mock registering with the go-micro registry and at the same time the reference implementation
func registerWithMicroReg(ns, address, name, desc, icon, prio string, cap int32, desktopOnly bool, mimeTypes []string) error {
	reg := oreg.GetRegistry()

	serviceID := ns + ".api.app-provider"

	node := &mreg.Node{
		//		Id:       serviceID + "-" + uuid.New().String(),
		Id:       serviceID + "-" + strings.ToLower(name),
		Address:  address,
		Metadata: make(map[string]string),
	}

	node.Metadata["registry"] = reg.String()
	node.Metadata["server"] = "grpc"
	node.Metadata["transport"] = "grpc"
	node.Metadata["protocol"] = "grpc"

	node.Metadata[ns+".app-provider.mime_type"] = joinMimeTypes(mimeTypes)
	node.Metadata[ns+".app-provider.name"] = name
	node.Metadata[ns+".app-provider.description"] = desc
	node.Metadata[ns+".app-provider.icon"] = icon

	node.Metadata[ns+".app-provider.allow_creation"] = registrypb.ProviderInfo_Capability_name[cap]
	node.Metadata[ns+".app-provider.priority"] = prio
	if desktopOnly {
		node.Metadata[ns+".app-provider.desktop_only"] = "true"
	}

	service := &mreg.Service{
		Name: serviceID,
		//Version:   version,
		Nodes:     []*mreg.Node{node},
		Endpoints: make([]*mreg.Endpoint, 0),
	}

	rOpts := []mreg.RegisterOption{mreg.RegisterTTL(time.Minute)}
	if err := reg.Register(service, rOpts...); err != nil {
		return err
	}

	return nil
}

// defineContext sets the context for the service. If there is a context configured it will create a new child from it,
// if not, it will create a root context that can be cancelled.
func defineContext(cfg *config.Config) (context.Context, context.CancelFunc) {
	return func() (context.Context, context.CancelFunc) {
		if cfg.Context == nil {
			return context.WithCancel(context.Background())
		}
		return context.WithCancel(cfg.Context)
	}()
}

// use the UTF-8 record separator
func splitMimeTypes(s string) []string {
	return strings.Split(s, "␞")
}

func joinMimeTypes(mimetypes []string) string {
	return strings.Join(mimetypes, "␞")
}
