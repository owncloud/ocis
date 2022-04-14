package command

import (
	"context"
	"flag"
	"os"
	"path"
	"path/filepath"

	"github.com/cs3org/reva/v2/cmd/revad/runtime"
	"github.com/gofrs/uuid"
	"github.com/oklog/run"
	"github.com/owncloud/ocis/extensions/storage/pkg/config"
	"github.com/owncloud/ocis/extensions/storage/pkg/server/debug"
	"github.com/owncloud/ocis/extensions/storage/pkg/tracing"
	ociscfg "github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/sync"
	"github.com/thejerf/suture/v4"
	"github.com/urfave/cli/v2"
)

// Groups is the entrypoint for the sharing command.
func Groups(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "groups",
		Usage: "start groups service",
		Before: func(c *cli.Context) error {
			return ParseConfig(c, cfg, "storage-groups")
		},
		Action: func(c *cli.Context) error {
			logger := NewLogger(cfg)
			tracing.Configure(cfg, logger)
			gr := run.Group{}
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			// pre-create folders
			if cfg.Reva.Groups.Driver == "json" && cfg.Reva.Groups.JSON != "" {
				if err := os.MkdirAll(filepath.Dir(cfg.Reva.Groups.JSON), os.FileMode(0700)); err != nil {
					return err
				}
			}

			cuuid := uuid.Must(uuid.NewV4())
			pidFile := path.Join(os.TempDir(), "revad-"+c.Command.Name+"-"+cuuid.String()+".pid")

			rcfg := groupsConfigFromStruct(c, cfg)

			if cfg.Reva.Groups.Driver == "ldap" {
				if err := waitForLDAPCA(logger, &cfg.Reva.LDAP); err != nil {
					logger.Error().Err(err).Msg("The configured LDAP CA cert does not exist")
					return err
				}
			}

			gr.Add(func() error {
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
				debug.Addr(cfg.Reva.Groups.DebugAddr),
				debug.Logger(logger),
				debug.Context(ctx),
				debug.Config(cfg),
			)

			if err != nil {
				logger.Info().Err(err).Str("server", c.Command.Name+"-debug").Msg("Failed to initialize server")
				return err
			}

			gr.Add(debugServer.ListenAndServe, func(_ error) {
				cancel()
			})

			if !cfg.Reva.Groups.Supervised {
				sync.Trap(&gr, cancel)
			}

			return gr.Run()
		},
	}
}

// groupsConfigFromStruct will adapt an oCIS config struct into a reva mapstructure to start a reva service.
func groupsConfigFromStruct(c *cli.Context, cfg *config.Config) map[string]interface{} {
	return map[string]interface{}{
		"core": map[string]interface{}{
			"max_cpus":             cfg.Reva.Groups.MaxCPUs,
			"tracing_enabled":      cfg.Tracing.Enabled,
			"tracing_endpoint":     cfg.Tracing.Endpoint,
			"tracing_collector":    cfg.Tracing.Collector,
			"tracing_service_name": c.Command.Name,
		},
		"shared": map[string]interface{}{
			"jwt_secret":                cfg.Reva.JWTSecret,
			"gatewaysvc":                cfg.Reva.Gateway.Endpoint,
			"skip_user_groups_in_token": cfg.Reva.SkipUserGroupsInToken,
		},
		"grpc": map[string]interface{}{
			"network": cfg.Reva.Groups.GRPCNetwork,
			"address": cfg.Reva.Groups.GRPCAddr,
			// TODO build services dynamically
			"services": map[string]interface{}{
				"groupprovider": map[string]interface{}{
					"driver": cfg.Reva.Groups.Driver,
					"drivers": map[string]interface{}{
						"json": map[string]interface{}{
							"groups": cfg.Reva.Groups.JSON,
						},
						"ldap": ldapConfigFromString(cfg),
						"rest": map[string]interface{}{
							"client_id":                      cfg.Reva.UserGroupRest.ClientID,
							"client_secret":                  cfg.Reva.UserGroupRest.ClientSecret,
							"redis_address":                  cfg.Reva.UserGroupRest.RedisAddress,
							"redis_username":                 cfg.Reva.UserGroupRest.RedisUsername,
							"redis_password":                 cfg.Reva.UserGroupRest.RedisPassword,
							"group_members_cache_expiration": cfg.Reva.Groups.GroupMembersCacheExpiration,
							"id_provider":                    cfg.Reva.UserGroupRest.IDProvider,
							"api_base_url":                   cfg.Reva.UserGroupRest.APIBaseURL,
							"oidc_token_endpoint":            cfg.Reva.UserGroupRest.OIDCTokenEndpoint,
							"target_api":                     cfg.Reva.UserGroupRest.TargetAPI,
						},
					},
				},
			},
		},
	}
}

// GroupSutureService allows for the storage-groupprovider command to be embedded and supervised by a suture supervisor tree.
type GroupSutureService struct {
	cfg *config.Config
}

// NewGroupProviderSutureService creates a new storage.GroupProvider
func NewGroupProvider(cfg *ociscfg.Config) suture.Service {
	cfg.Storage.Commons = cfg.Commons
	return GroupSutureService{
		cfg: cfg.Storage,
	}
}

func (s GroupSutureService) Serve(ctx context.Context) error {
	s.cfg.Reva.Groups.Context = ctx
	f := &flag.FlagSet{}
	cmdFlags := Groups(s.cfg).Flags
	for k := range cmdFlags {
		if err := cmdFlags[k].Apply(f); err != nil {
			return err
		}
	}
	cliCtx := cli.NewContext(nil, f, nil)
	if Groups(s.cfg).Before != nil {
		if err := Groups(s.cfg).Before(cliCtx); err != nil {
			return err
		}
	}
	if err := Groups(s.cfg).Action(cliCtx); err != nil {
		return err
	}

	return nil
}
