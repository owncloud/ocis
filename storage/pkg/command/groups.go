package command

import (
	"context"
	"flag"
	"os"
	"path"
	"path/filepath"

	"github.com/cs3org/reva/cmd/revad/runtime"
	"github.com/gofrs/uuid"
	"github.com/oklog/run"
	ociscfg "github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/sync"
	"github.com/owncloud/ocis/storage/pkg/config"
	"github.com/owncloud/ocis/storage/pkg/flagset"
	"github.com/owncloud/ocis/storage/pkg/server/debug"
	"github.com/owncloud/ocis/storage/pkg/tracing"
	"github.com/thejerf/suture/v4"
	"github.com/urfave/cli/v2"
)

// Groups is the entrypoint for the sharing command.
func Groups(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "groups",
		Usage: "Start groups service",
		Flags: flagset.GroupsWithConfig(cfg),
		Before: func(c *cli.Context) error {
			cfg.Reva.Groups.Services = c.StringSlice("service")

			return nil
		},
		Action: func(c *cli.Context) error {
			logger := NewLogger(cfg)
			tracing.Configure(cfg, logger)
			gr := run.Group{}
			ctx, cancel := context.WithCancel(context.Background())

			// pre-create folders
			if cfg.Reva.Groups.Driver == "json" && cfg.Reva.Groups.JSON != "" {
				if err := os.MkdirAll(filepath.Dir(cfg.Reva.Groups.JSON), os.FileMode(0700)); err != nil {
					return err
				}
			}

			uuid := uuid.Must(uuid.NewV4())
			pidFile := path.Join(os.TempDir(), "revad-"+c.Command.Name+"-"+uuid.String()+".pid")
			defer cancel()

			rcfg := groupsConfigFromStruct(c, cfg)

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
						"ldap": map[string]interface{}{
							"hostname":        cfg.Reva.LDAP.Hostname,
							"port":            cfg.Reva.LDAP.Port,
							"cacert":          cfg.Reva.LDAP.CACert,
							"insecure":        cfg.Reva.LDAP.Insecure,
							"base_dn":         cfg.Reva.LDAP.BaseDN,
							"groupfilter":     cfg.Reva.LDAP.GroupFilter,
							"attributefilter": cfg.Reva.LDAP.GroupAttributeFilter,
							"findfilter":      cfg.Reva.LDAP.GroupFindFilter,
							"memberfilter":    cfg.Reva.LDAP.GroupMemberFilter,
							"bind_username":   cfg.Reva.LDAP.BindDN,
							"bind_password":   cfg.Reva.LDAP.BindPassword,
							"idp":             cfg.Reva.LDAP.IDP,
							"schema": map[string]interface{}{
								"dn":          "dn",
								"gid":         cfg.Reva.LDAP.GroupSchema.GID,
								"mail":        cfg.Reva.LDAP.GroupSchema.Mail,
								"displayName": cfg.Reva.LDAP.GroupSchema.DisplayName,
								"cn":          cfg.Reva.LDAP.GroupSchema.CN,
								"gidNumber":   cfg.Reva.LDAP.GroupSchema.GIDNumber,
							},
						},
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
	cfg.Storage.Log = cfg.Commons.Log
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
