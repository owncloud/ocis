package command

import (
	"context"
	"reflect"
	"strings"

	"github.com/imdario/mergo"
	"github.com/jinzhu/copier"
	"github.com/mitchellh/cli"
	"github.com/owncloud/ocis/ocis-pkg/sync"

	"github.com/micro/cli/v2"
	"github.com/oklog/run"
	"github.com/owncloud/ocis/accounts/pkg/config"
	"github.com/owncloud/ocis/accounts/pkg/flagset"
	"github.com/owncloud/ocis/accounts/pkg/metrics"
	"github.com/owncloud/ocis/accounts/pkg/server/grpc"
	"github.com/owncloud/ocis/accounts/pkg/server/http"
	svc "github.com/owncloud/ocis/accounts/pkg/service/v0"
	"github.com/owncloud/ocis/accounts/pkg/tracing"
)

// Server is the entry point for the server command.
func Server(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:        "server",
		Usage:       "Start ocis accounts service",
		Description: "uses an LDAP server as the storage backend",
		Flags:       flagset.ServerWithConfig(cfg),
		Before: func(ctx *cli.Context) error {
			if cfg.HTTP.Root != "/" {
				cfg.HTTP.Root = strings.TrimSuffix(cfg.HTTP.Root, "/")
			}

			// beforeOverride contains cfg with values parsed by urfavecli,
			// this should take precedence when merging as they are more explicit.
			// beforeOverride has the highest priority, as they are inherited values.
			beforeOverride := config.Config{}
			if err := copier.Copy(&beforeOverride, cfg); err != nil {
				return err
			}

			defaultConfig := config.DefaultConfig()

			// By the time we unmarshal viper parsed values onto cfg, any value having been set by the cli framework
			// will get overridden, in order to ensure that this values are accounted for we have to perform a 3-way merge:
			// 1. merge viper onto cfg
			// 2. merge defaults onto cfg
			// 3. merge parsed flags onto cfg
			// the result of this is the same order of precedence as the cli framework claims, except a new "artificial"
			// source which accounts for structured configuration.
			if !cfg.Supervised {
				if err := ParseConfig(ctx, cfg); err != nil {
					return err
				}
			}

			fromAccountsConfigFile := config.Config{}
			if err := ParseConfig(ctx, &fromAccountsConfigFile); err != nil {
				return err
			}

			if err := mergo.Merge(cfg, defaultConfig); err != nil {
				panic(err)
			}

			// When an extension is running supervised, we have the use case where executing `ocis run extension`
			// we want to ONLY take into consideration fhe existing config file. This is a hard requirement.
			if !reflect.DeepEqual(fromAccountsConfigFile, config.Config{}) {
				if err := mergo.Merge(cfg, fromAccountsConfigFile); err != nil {
					panic(err)
				}
				return nil
			}

			if err := mergo.Merge(cfg, fromAccountsConfigFile); err != nil {
				panic(err)
			}

			// preserves the original order from inherited values. This has the drawback that also persists values
			// inherited from an ocis.yaml global config file, with the side effect of these global values overriding
			// concrete values from a closer-to-the-process proxy.yaml file.
			if err := mergo.Merge(cfg, beforeOverride, mergo.WithOverride); err != nil {
				panic(err)
			}

			return nil
		},
		Action: func(c *cli.Context) error {
			logger := NewLogger(cfg)
			err := tracing.Configure(cfg, logger)
			if err != nil {
				return err
			}
			gr := run.Group{}
			ctx, cancel := defineContext(cfg)
			mtrcs := metrics.New()

			defer cancel()

			mtrcs.BuildInfo.WithLabelValues(cfg.Server.Version).Set(1)

			handler, err := svc.New(svc.Logger(logger), svc.Config(cfg))
			if err != nil {
				logger.Error().Err(err).Msg("handler init")
				return err
			}

			httpServer := http.Server(
				http.Config(cfg),
				http.Logger(logger),
				http.Name(cfg.Server.Name),
				http.Context(ctx),
				http.Metrics(mtrcs),
				http.Handler(handler),
			)

			gr.Add(httpServer.Run, func(_ error) {
				logger.Info().Str("server", "http").Msg("shutting down server")
				cancel()
			})

			grpcServer := grpc.Server(
				grpc.Config(cfg),
				grpc.Logger(logger),
				grpc.Name(cfg.Server.Name),
				grpc.Context(ctx),
				grpc.Metrics(mtrcs),
				grpc.Handler(handler),
			)

			gr.Add(grpcServer.Run, func(_ error) {
				logger.Info().Str("server", "grpc").Msg("shutting down server")
				cancel()
			})

			if !cfg.Supervised {
				sync.Trap(&gr, cancel)
			}

			return gr.Run()
		},
	}
}

// defineContext sets the context for the extension. If there is a context configured it will create a new child from it,
// if not, it will create a root context that can be cancelled.
func defineContext(cfg *config.Config) (context.Context, context.CancelFunc) {
	return func() (context.Context, context.CancelFunc) {
		if cfg.Context == nil {
			return context.WithCancel(context.Background())
		}
		return context.WithCancel(cfg.Context)
	}()
}
