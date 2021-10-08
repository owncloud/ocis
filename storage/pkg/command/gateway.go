package command

import (
	"context"
	"encoding/json"
	"flag"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/owncloud/ocis/storage/pkg/tracing"

	"github.com/owncloud/ocis/ocis-pkg/sync"
	"github.com/owncloud/ocis/ocis-pkg/version"

	"github.com/cs3org/reva/cmd/revad/runtime"
	"github.com/gofrs/uuid"
	"github.com/oklog/run"
	ociscfg "github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/log"
	"github.com/owncloud/ocis/storage/pkg/config"
	"github.com/owncloud/ocis/storage/pkg/flagset"
	"github.com/owncloud/ocis/storage/pkg/server/debug"
	"github.com/owncloud/ocis/storage/pkg/service/external"
	"github.com/thejerf/suture/v4"
	"github.com/urfave/cli/v2"
)

// Gateway is the entrypoint for the gateway command.
func Gateway(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "gateway",
		Usage: "Start gateway",
		Flags: flagset.GatewayWithConfig(cfg),
		Before: func(c *cli.Context) error {
			cfg.Reva.Gateway.Services = c.StringSlice("service")
			cfg.Reva.StorageRegistry.Rules = c.StringSlice("storage-registry-rule")

			if cfg.Reva.DataGateway.PublicURL == "" {
				cfg.Reva.DataGateway.PublicURL = strings.TrimRight(cfg.Reva.Frontend.PublicURL, "/") + "/data"
			}

			return nil
		},
		Action: func(c *cli.Context) error {
			logger := NewLogger(cfg)
			tracing.Configure(cfg, logger)
			gr := run.Group{}
			ctx, cancel := context.WithCancel(context.Background())
			uuid := uuid.Must(uuid.NewV4())
			pidFile := path.Join(os.TempDir(), "revad-"+c.Command.Name+"-"+uuid.String()+".pid")
			rcfg := gatewayConfigFromStruct(c, cfg, logger)
			logger.Debug().
				Str("server", "gateway").
				Interface("reva-config", rcfg).
				Msg("config")

			defer cancel()

			gr.Add(func() error {
				err := external.RegisterGRPCEndpoint(
					ctx,
					"com.owncloud.storage",
					uuid.String(),
					cfg.Reva.Gateway.GRPCAddr,
					version.String,
					logger,
				)

				if err != nil {
					return err
				}

				runtime.RunWithOptions(
					rcfg,
					pidFile,
					runtime.WithLogger(&logger.Logger),
				)
				return nil
			}, func(_ error) {
				logger.Info().
					Str("server", c.Command.Name).
					Msg("Shutting down server")

				cancel()
			})

			debugServer, err := debug.Server(
				debug.Name(c.Command.Name+"-debug"),
				debug.Addr(cfg.Reva.Gateway.DebugAddr),
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

			if !cfg.Reva.Gateway.Supervised {
				sync.Trap(&gr, cancel)
			}

			return gr.Run()
		},
	}
}

// gatewayConfigFromStruct will adapt an oCIS config struct into a reva mapstructure to start a reva service.
func gatewayConfigFromStruct(c *cli.Context, cfg *config.Config, logger log.Logger) map[string]interface{} {
	rcfg := map[string]interface{}{
		"core": map[string]interface{}{
			"max_cpus":             cfg.Reva.Users.MaxCPUs,
			"tracing_enabled":      cfg.Tracing.Enabled,
			"tracing_endpoint":     cfg.Tracing.Endpoint,
			"tracing_collector":    cfg.Tracing.Collector,
			"tracing_service_name": c.Command.Name,
		},
		"shared": map[string]interface{}{
			"jwt_secret": cfg.Reva.JWTSecret,
			"gatewaysvc": cfg.Reva.Gateway.Endpoint,
		},
		"grpc": map[string]interface{}{
			"network": cfg.Reva.Gateway.GRPCNetwork,
			"address": cfg.Reva.Gateway.GRPCAddr,
			// TODO build services dynamically
			"services": map[string]interface{}{
				"gateway": map[string]interface{}{
					// registries is located on the gateway
					"authregistrysvc":    cfg.Reva.Gateway.Endpoint,
					"storageregistrysvc": cfg.Reva.Gateway.Endpoint,
					"appregistrysvc":     cfg.Reva.Gateway.Endpoint,
					// user metadata is located on the users services
					"preferencessvc":   cfg.Reva.Users.Endpoint,
					"userprovidersvc":  cfg.Reva.Users.Endpoint,
					"groupprovidersvc": cfg.Reva.Groups.Endpoint,
					// sharing is located on the sharing service
					"usershareprovidersvc":          cfg.Reva.Sharing.Endpoint,
					"publicshareprovidersvc":        cfg.Reva.Sharing.Endpoint,
					"ocmshareprovidersvc":           cfg.Reva.Sharing.Endpoint,
					"commit_share_to_storage_grant": cfg.Reva.Gateway.CommitShareToStorageGrant,
					"commit_share_to_storage_ref":   cfg.Reva.Gateway.CommitShareToStorageRef,
					"share_folder":                  cfg.Reva.Gateway.ShareFolder, // ShareFolder is the location where to create shares in the recipient's storage provider.
					// other
					"disable_home_creation_on_login": cfg.Reva.Gateway.DisableHomeCreationOnLogin,
					"datagateway":                    cfg.Reva.DataGateway.PublicURL,
					"transfer_shared_secret":         cfg.Reva.TransferSecret,
					"transfer_expires":               cfg.Reva.TransferExpires,
					"home_mapping":                   cfg.Reva.Gateway.HomeMapping,
					"etag_cache_ttl":                 cfg.Reva.Gateway.EtagCacheTTL,
				},
				"authregistry": map[string]interface{}{
					"driver": "static",
					"drivers": map[string]interface{}{
						"static": map[string]interface{}{
							"rules": map[string]interface{}{
								"basic":        cfg.Reva.AuthBasic.Endpoint,
								"bearer":       cfg.Reva.AuthBearer.Endpoint,
								"publicshares": cfg.Reva.StoragePublicLink.Endpoint,
							},
						},
					},
				},
				"appregistry": map[string]interface{}{
					"driver": "static",
					"drivers": map[string]interface{}{
						"static": map[string]interface{}{
							"mime_types": mimetypes(cfg, logger),
						},
					},
				},
				"storageregistry": map[string]interface{}{
					"driver": cfg.Reva.StorageRegistry.Driver,
					"drivers": map[string]interface{}{
						"static": map[string]interface{}{
							"home_provider": cfg.Reva.StorageRegistry.HomeProvider,
							"rules":         rules(cfg, logger),
						},
					},
				},
			},
		},
	}
	return rcfg
}

func rules(cfg *config.Config, logger log.Logger) map[string]map[string]interface{} {

	// if a list of rules is given it overrides the generated rules from below
	if len(cfg.Reva.StorageRegistry.Rules) > 0 {
		rules := map[string]map[string]interface{}{}
		for i := range cfg.Reva.StorageRegistry.Rules {
			parts := strings.SplitN(cfg.Reva.StorageRegistry.Rules[i], "=", 2)
			rules[parts[0]] = map[string]interface{}{"address": parts[1]}
		}
		return rules
	}

	// check if the rules have to be read from a json file
	if cfg.Reva.StorageRegistry.JSON != "" {
		data, err := ioutil.ReadFile(cfg.Reva.StorageRegistry.JSON)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to read storage registry rules from JSON file: " + cfg.Reva.StorageRegistry.JSON)
			return nil
		}
		var rules map[string]map[string]interface{}
		if err = json.Unmarshal(data, &rules); err != nil {
			logger.Error().Err(err).Msg("Failed to unmarshal storage registry rules")
			return nil
		}
		return rules
	}

	// generate rules based on default config
	return map[string]map[string]interface{}{
		cfg.Reva.StorageHome.MountPath:       {"address": cfg.Reva.StorageHome.Endpoint},
		cfg.Reva.StorageHome.MountID:         {"address": cfg.Reva.StorageHome.Endpoint},
		cfg.Reva.StorageUsers.MountPath:      {"address": cfg.Reva.StorageUsers.Endpoint},
		cfg.Reva.StorageUsers.MountID + ".*": {"address": cfg.Reva.StorageUsers.Endpoint},
		cfg.Reva.StoragePublicLink.MountPath: {"address": cfg.Reva.StoragePublicLink.Endpoint},
		// public link storage returns the mount id of the actual storage
		// medatada storage not part of the global namespace
	}
}

func mimetypes(cfg *config.Config, logger log.Logger) map[string]map[string]interface{} {

	type mimeTypeConfig struct {
		Extension     string `json:"extension" mapstructure:"extension"`
		Name          string `json:"name" mapstructure:"name"`
		Description   string `json:"description" mapstructure:"description"`
		Icon          string `json:"icon" mapstructure:"icon"`
		DefaultApp    string `json:"default_app" mapstructure:"default_app"`
		AllowCreation bool   `json:"allow_creation" mapstructure:"allow_creation"`
	}
	var mimetypes map[string]mimeTypeConfig
	var m map[string]map[string]interface{}

	// load default app mimetypes from a json file
	if cfg.Reva.AppRegistry.MimetypesJSON != "" {
		data, err := ioutil.ReadFile(cfg.Reva.AppRegistry.MimetypesJSON)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to read app registry mimetypes from JSON file: " + cfg.Reva.AppRegistry.MimetypesJSON)
			return nil
		}
		if err = json.Unmarshal(data, &mimetypes); err != nil {
			logger.Error().Err(err).Msg("Failed to unmarshal storage registry rules")
			return nil
		}
		if err := mapstructure.Decode(mimetypes, &m); err != nil {
			logger.Error().Err(err).Msg("Failed to decode defaultapp registry mimetypes to mapstructure")
			return nil
		}
		return m
	}

	logger.Info().Msg("No app registry mimetypes JSON file provided, loading default configuration")

	mimetypes = map[string]mimeTypeConfig{
		"application/pdf": {
			Extension:   "pdf",
			Name:        "PDF",
			Description: "PDF document",
		},
		"application/vnd.oasis.opendocument.text": {
			Extension:     "odt",
			Name:          "OpenDocument",
			Description:   "OpenDocument text document",
			AllowCreation: true,
		},
		"application/vnd.oasis.opendocument.spreadsheet": {
			Extension:     "ods",
			Name:          "OpenSpreadsheet",
			Description:   "OpenDocument spreadsheet document",
			AllowCreation: true,
		},
		"application/vnd.oasis.opendocument.presentation": {
			Extension:     "odp",
			Name:          "OpenPresentation",
			Description:   "OpenDocument presentation document",
			AllowCreation: true,
		},
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document": {
			Extension:     "docx",
			Name:          "Microsoft Word",
			Description:   "Microsoft Word document",
			AllowCreation: true,
		},
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet": {
			Extension:     "xlsx",
			Name:          "Microsoft Excel",
			Description:   "Microsoft Excel document",
			AllowCreation: true,
		},
		"application/vnd.openxmlformats-officedocument.presentationml.presentation": {
			Extension:     "pptx",
			Name:          "Microsoft PowerPoint",
			Description:   "Microsoft PowerPoint document",
			AllowCreation: true,
		},
		"application/vnd.jupyter": {
			Extension:   "ipynb",
			Name:        "Jupyter Notebook",
			Description: "Jupyter Notebook",
		},
		"text/markdown": {
			Extension:     "md",
			Name:          "Markdown file",
			Description:   "Markdown file",
			AllowCreation: true,
		},
		"application/compressed-markdown": {
			Extension:   "zmd",
			Name:        "Compressed markdown file",
			Description: "Compressed markdown file",
		},
	}

	if err := mapstructure.Decode(mimetypes, &m); err != nil {
		logger.Error().Err(err).Msg("Failed to decode defaultapp registry mimetypes to mapstructure")
		return nil
	}
	return m

}

// GatewaySutureService allows for the storage-gateway command to be embedded and supervised by a suture supervisor tree.
type GatewaySutureService struct {
	cfg *config.Config
}

// NewGatewaySutureService creates a new gateway.GatewaySutureService
func NewGateway(cfg *ociscfg.Config) suture.Service {
	if cfg.Mode == 0 {
		cfg.Storage.Reva.Gateway.Supervised = true
	}
	return GatewaySutureService{
		cfg: cfg.Storage,
	}
}

func (s GatewaySutureService) Serve(ctx context.Context) error {
	s.cfg.Reva.Gateway.Context = ctx
	f := &flag.FlagSet{}
	cmdFlags := Gateway(s.cfg).Flags
	for k := range cmdFlags {
		if err := cmdFlags[k].Apply(f); err != nil {
			return err
		}
	}
	cliCtx := cli.NewContext(nil, f, nil)
	if Gateway(s.cfg).Before != nil {
		if err := Gateway(s.cfg).Before(cliCtx); err != nil {
			return err
		}
	}
	if err := Gateway(s.cfg).Action(cliCtx); err != nil {
		return err
	}

	return nil
}
