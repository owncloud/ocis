package command

import (
	"context"
	"encoding/json"
	"flag"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/cs3org/reva/v2/cmd/revad/runtime"
	"github.com/gofrs/uuid"
	"github.com/mitchellh/mapstructure"
	"github.com/oklog/run"
	"github.com/owncloud/ocis/extensions/gateway/pkg/config"
	"github.com/owncloud/ocis/extensions/storage/pkg/server/debug"
	"github.com/owncloud/ocis/extensions/storage/pkg/service/external"
	ociscfg "github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/log"
	"github.com/owncloud/ocis/ocis-pkg/sync"
	"github.com/owncloud/ocis/ocis-pkg/tracing"
	"github.com/owncloud/ocis/ocis-pkg/version"
	"github.com/thejerf/suture/v4"
	"github.com/urfave/cli/v2"
)

// Gateway is the entrypoint for the gateway command.
func Gateway(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "gateway",
		Usage: "start gateway",
		Before: func(c *cli.Context) error {
			if cfg.DataGatewayPublicURL == "" {
				cfg.DataGatewayPublicURL = strings.TrimRight(cfg.FrontendPublicURL, "/") + "/data"
			}

			return nil
		},
		Action: func(c *cli.Context) error {
			logCfg := cfg.Logging
			logger := log.NewLogger(
				log.Level(logCfg.Level),
				log.File(logCfg.File),
				log.Pretty(logCfg.Pretty),
				log.Color(logCfg.Color),
			)
			tracing.Configure(cfg.Tracing.Enabled, cfg.Tracing.Type, logger)
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
					cfg.GRPC.Addr,
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
				debug.Addr(cfg.Debug.Addr),
				debug.Logger(logger),
				debug.Context(ctx),
				debug.Pprof(cfg.Debug.Pprof),
				debug.Zpages(cfg.Debug.Zpages),
				debug.Token(cfg.Debug.Token),
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

			return gr.Run()
		},
	}
}

// gatewayConfigFromStruct will adapt an oCIS config struct into a reva mapstructure to start a reva service.
func gatewayConfigFromStruct(c *cli.Context, cfg *config.Config, logger log.Logger) map[string]interface{} {
	rcfg := map[string]interface{}{
		"core": map[string]interface{}{
			"tracing_enabled":      cfg.Tracing.Enabled,
			"tracing_endpoint":     cfg.Tracing.Endpoint,
			"tracing_collector":    cfg.Tracing.Collector,
			"tracing_service_name": c.Command.Name,
		},
		"shared": map[string]interface{}{
			"jwt_secret":                cfg.JWTSecret,
			"gatewaysvc":                cfg.GatewayEndpoint,
			"skip_user_groups_in_token": cfg.SkipUserGroupsInToken,
		},
		"grpc": map[string]interface{}{
			"network": cfg.GRPC.Protocol,
			"address": cfg.GRPC.Addr,
			// TODO build services dynamically
			"services": map[string]interface{}{
				"gateway": map[string]interface{}{
					// registries is located on the gateway
					"authregistrysvc":    cfg.GatewayEndpoint,
					"storageregistrysvc": cfg.GatewayEndpoint,
					"appregistrysvc":     cfg.GatewayEndpoint,
					// user metadata is located on the users services
					"preferencessvc":   cfg.UsersEndpoint,
					"userprovidersvc":  cfg.UsersEndpoint,
					"groupprovidersvc": cfg.GroupsEndpoint,
					"permissionssvc":   cfg.PermissionsEndpoint,
					// sharing is located on the sharing service
					"usershareprovidersvc":          cfg.SharingEndpoint,
					"publicshareprovidersvc":        cfg.SharingEndpoint,
					"ocmshareprovidersvc":           cfg.SharingEndpoint,
					"commit_share_to_storage_grant": cfg.CommitShareToStorageGrant,
					"commit_share_to_storage_ref":   cfg.CommitShareToStorageRef,
					"share_folder":                  cfg.ShareFolder, // ShareFolder is the location where to create shares in the recipient's storage provider.
					// other
					"disable_home_creation_on_login": cfg.DisableHomeCreationOnLogin,
					"datagateway":                    cfg.DataGatewayPublicURL,
					"transfer_shared_secret":         cfg.TransferSecret,
					"transfer_expires":               cfg.TransferExpires,
					"home_mapping":                   cfg.HomeMapping,
					"etag_cache_ttl":                 cfg.EtagCacheTTL,
				},
				"authregistry": map[string]interface{}{
					"driver": "static",
					"drivers": map[string]interface{}{
						"static": map[string]interface{}{
							"rules": map[string]interface{}{
								"basic":        cfg.AuthBasicEndpoint,
								"bearer":       cfg.AuthBearerEndpoint,
								"machine":      cfg.AuthMachineEndpoint,
								"publicshares": cfg.StoragePublicLinkEndpoint,
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
					"driver": cfg.StorageRegistry.Driver,
					"drivers": map[string]interface{}{
						"static": map[string]interface{}{
							"home_provider": "/home", // TODO use /users/{{.OpaqueId}} ?
							"rules": map[string]interface{}{
								"/home": map[string]interface{}{
									"address":       "127.0.0.1:9157",
									"provider_id":   "1284d238-aa92-42ce-bdc4-0b0000009157",
									"provider_path": "/users",
									"space_type":    "personal",
								},
								//"/users/{{.Id.OpaqueId}}": map[string]interface{}{
								"/users": map[string]interface{}{
									"address":       "127.0.0.1:9157",
									"provider_id":   "1284d238-aa92-42ce-bdc4-0b0000009157",
									"provider_path": "/users",
								},
								"1284d238-aa92-42ce-bdc4-0b0000009157": map[string]interface{}{
									"address":       "127.0.0.1:9157",
									"provider_id":   "1284d238-aa92-42ce-bdc4-0b0000009157",
									"provider_path": "/users",
								},
								"/project": map[string]interface{}{
									"address":       "127.0.0.1:9157",
									"provider_id":   "df7debc5-3491-4b7a-8b0d-6888009adb28",
									"provider_path": "/project",
									"space_type":    "project",
								},
								"df7debc5-3491-4b7a-8b0d-6888009adb28": map[string]interface{}{
									"address":       "127.0.0.1:9157",
									"provider_id":   "df7debc5-3491-4b7a-8b0d-6888009adb28",
									"provider_path": "/project",
									"space_type":    "project",
								},
								"/public": map[string]interface{}{
									"address":       "localhost:9178",
									"provider_id":   "7993447f-687f-490d-875c-ac95e89a62a4",
									"provider_path": "/public",
								},
								"7993447f-687f-490d-875c-ac95e89a62a4": map[string]interface{}{
									"address":       "localhost:9178",
									"provider_id":   "7993447f-687f-490d-875c-ac95e89a62a4",
									"provider_path": "/public",
								},
								/*
									"metadata": map[string]interface{}{
										"address":     "127.0.0.1:9215",
										"provider_id": "0dba9855-3ab1-432f-ace7-e01224fe2c65",
									},
									"0dba9855-3ab1-432f-ace7-e01224fe2c65": map[string]interface{}{
										"address":       "127.0.0.1:9215",
										"provider_id":   "0dba9855-3ab1-432f-ace7-e01224fe2c65",
										"provider_path": "metadata",
									},
								*/
							},
						},
						"spaces": map[string]interface{}{
							"providers": spacesProviders(cfg, logger),
						},
					},
				},
			},
		},
	}
	return rcfg
}

func spacesProviders(cfg *config.Config, logger log.Logger) map[string]map[string]interface{} {

	// if a list of rules is given it overrides the generated rules from below
	if len(cfg.StorageRegistry.Rules) > 0 {
		rules := map[string]map[string]interface{}{}
		for i := range cfg.StorageRegistry.Rules {
			parts := strings.SplitN(cfg.StorageRegistry.Rules[i], "=", 2)
			rules[parts[0]] = map[string]interface{}{"address": parts[1]}
		}
		return rules
	}

	// check if the rules have to be read from a json file
	if cfg.StorageRegistry.JSON != "" {
		data, err := ioutil.ReadFile(cfg.StorageRegistry.JSON)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to read storage registry rules from JSON file: " + cfg.StorageRegistry.JSON)
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
		cfg.StorageUsersEndpoint: {
			"spaces": map[string]interface{}{
				"personal": map[string]interface{}{
					"mount_point":   "/users",
					"path_template": "/users/{{.Space.Owner.Id.OpaqueId}}",
				},
				"project": map[string]interface{}{
					"mount_point":   "/projects",
					"path_template": "/projects/{{.Space.Name}}",
				},
			},
		},
		cfg.StorageSharesEndpoint: {
			"spaces": map[string]interface{}{
				"virtual": map[string]interface{}{
					// The root of the share jail is mounted here
					"mount_point": "/users/{{.CurrentUser.Id.OpaqueId}}/Shares",
				},
				"grant": map[string]interface{}{
					// Grants are relative to a space root that the gateway will determine with a stat
					"mount_point": ".",
				},
				"mountpoint": map[string]interface{}{
					// The jail needs to be filled with mount points
					// .Space.Name is a path relative to the mount point
					"mount_point":   "/users/{{.CurrentUser.Id.OpaqueId}}/Shares",
					"path_template": "/users/{{.CurrentUser.Id.OpaqueId}}/Shares/{{.Space.Name}}",
				},
			},
		},
		// public link storage returns the mount id of the actual storage
		cfg.StoragePublicLinkEndpoint: {
			"spaces": map[string]interface{}{
				"grant": map[string]interface{}{
					"mount_point": ".",
				},
				"mountpoint": map[string]interface{}{
					"mount_point":   "/public",
					"path_template": "/public/{{.Space.Root.OpaqueId}}",
				},
			},
		},
		// medatada storage not part of the global namespace
	}
}

func mimetypes(cfg *config.Config, logger log.Logger) []map[string]interface{} {

	type mimeTypeConfig struct {
		MimeType      string `json:"mime_type" mapstructure:"mime_type"`
		Extension     string `json:"extension" mapstructure:"extension"`
		Name          string `json:"name" mapstructure:"name"`
		Description   string `json:"description" mapstructure:"description"`
		Icon          string `json:"icon" mapstructure:"icon"`
		DefaultApp    string `json:"default_app" mapstructure:"default_app"`
		AllowCreation bool   `json:"allow_creation" mapstructure:"allow_creation"`
	}
	var mimetypes []mimeTypeConfig
	var m []map[string]interface{}

	// load default app mimetypes from a json file
	if cfg.AppRegistry.MimetypesJSON != "" {
		data, err := ioutil.ReadFile(cfg.AppRegistry.MimetypesJSON)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to read app registry mimetypes from JSON file: " + cfg.AppRegistry.MimetypesJSON)
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

	mimetypes = []mimeTypeConfig{
		{
			MimeType:    "application/pdf",
			Extension:   "pdf",
			Name:        "PDF",
			Description: "PDF document",
		},
		{
			MimeType:      "application/vnd.oasis.opendocument.text",
			Extension:     "odt",
			Name:          "OpenDocument",
			Description:   "OpenDocument text document",
			AllowCreation: true,
		},
		{
			MimeType:      "application/vnd.oasis.opendocument.spreadsheet",
			Extension:     "ods",
			Name:          "OpenSpreadsheet",
			Description:   "OpenDocument spreadsheet document",
			AllowCreation: true,
		},
		{
			MimeType:      "application/vnd.oasis.opendocument.presentation",
			Extension:     "odp",
			Name:          "OpenPresentation",
			Description:   "OpenDocument presentation document",
			AllowCreation: true,
		},
		{
			MimeType:      "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
			Extension:     "docx",
			Name:          "Microsoft Word",
			Description:   "Microsoft Word document",
			AllowCreation: true,
		},
		{
			MimeType:      "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
			Extension:     "xlsx",
			Name:          "Microsoft Excel",
			Description:   "Microsoft Excel document",
			AllowCreation: true,
		},
		{
			MimeType:      "application/vnd.openxmlformats-officedocument.presentationml.presentation",
			Extension:     "pptx",
			Name:          "Microsoft PowerPoint",
			Description:   "Microsoft PowerPoint document",
			AllowCreation: true,
		},
		{
			MimeType:    "application/vnd.jupyter",
			Extension:   "ipynb",
			Name:        "Jupyter Notebook",
			Description: "Jupyter Notebook",
		},
		{
			MimeType:      "text/markdown",
			Extension:     "md",
			Name:          "Markdown file",
			Description:   "Markdown file",
			AllowCreation: true,
		},
		{
			MimeType:    "application/compressed-markdown",
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
	cfg.Gateway.Commons = cfg.Commons
	return GatewaySutureService{
		cfg: cfg.Gateway,
	}
}

func (s GatewaySutureService) Serve(ctx context.Context) error {
	cmd := Gateway(s.cfg)
	f := &flag.FlagSet{}
	cmdFlags := cmd.Flags
	for k := range cmdFlags {
		if err := cmdFlags[k].Apply(f); err != nil {
			return err
		}
	}
	cliCtx := cli.NewContext(nil, f, nil)
	if cmd.Before != nil {
		if err := cmd.Before(cliCtx); err != nil {
			return err
		}
	}
	if err := cmd.Action(cliCtx); err != nil {
		return err
	}

	return nil
}
