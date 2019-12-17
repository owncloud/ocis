package command

import (
	"context"

	"github.com/micro/cli"
	"github.com/oklog/run"
	"github.com/owncloud/ocis-hello/pkg/command"
	svcconfig "github.com/owncloud/ocis-hello/pkg/config"
	"github.com/owncloud/ocis-hello/pkg/flagset"
	"github.com/owncloud/ocis-hello/pkg/metrics"
	"github.com/owncloud/ocis-hello/pkg/server/grpc"
	"github.com/owncloud/ocis-hello/pkg/server/http"
	"github.com/owncloud/ocis/pkg/config"
	"github.com/owncloud/ocis/pkg/register"
)

// HelloCommand is the entrypoint for the hello command.
func HelloCommand(cfg *config.Config) cli.Command {
	return cli.Command{
		Name:     "hello",
		Usage:    "Start hello server",
		Category: "Extensions",
		Flags:    flagset.ServerWithConfig(cfg.Hello),
		Action: func(c *cli.Context) error {
			scfg := configureHello(cfg)

			return cli.HandleAction(
				command.Server(scfg).Action,
				c,
			)
		},
	}
}

// HelloHandler defines the direct server handler.
func HelloHandler(ctx context.Context, cancel context.CancelFunc, gr *run.Group, cfg *config.Config) error {
	scfg := configureHello(cfg)
	logger := command.NewLogger(scfg)
	m := metrics.New()

	{
		server, err := http.Server(
			http.Logger(logger),
			http.Context(ctx),
			http.Config(scfg),
			http.Metrics(m),
		)

		if err != nil {
			logger.Info().
				Err(err).
				Str("transport", "http").
				Msg("Failed to initialize server")

			return err
		}

		gr.Add(func() error {
			return server.Run()
		}, func(_ error) {
			logger.Info().
				Str("transport", "http").
				Msg("Shutting down server")

			cancel()
		})
	}

	{
		server, err := grpc.Server(
			grpc.Logger(logger),
			grpc.Context(ctx),
			grpc.Config(scfg),
			grpc.Metrics(m),
		)

		if err != nil {
			logger.Info().
				Err(err).
				Str("transport", "grpc").
				Msg("Failed to initialize server")

			return err
		}

		gr.Add(func() error {
			return server.Run()
		}, func(_ error) {
			logger.Info().
				Str("transport", "grpc").
				Msg("Shutting down server")

			cancel()
		})
	}

	return nil
}

func configureHello(cfg *config.Config) *svcconfig.Config {
	cfg.Hello.Log.Level = cfg.Log.Level
	cfg.Hello.Log.Pretty = cfg.Log.Pretty
	cfg.Hello.Log.Color = cfg.Log.Color
	cfg.Hello.Tracing.Enabled = false
	cfg.Hello.HTTP.Addr = "localhost:9105"
	cfg.Hello.HTTP.Root = "/"
	cfg.Hello.GRPC.Addr = "localhost:9106"

	return cfg.Hello
}

func init() {
	register.AddCommand(HelloCommand)
	register.AddHandler(HelloHandler)
}
