package command

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/oklog/run"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/configlog"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/ocis-pkg/service/debug"
	"github.com/owncloud/ocis/v2/ocis-pkg/version"
	"github.com/owncloud/ocis/v2/services/antivirus/pkg/config"
	"github.com/owncloud/ocis/v2/services/antivirus/pkg/config/parser"
	"github.com/owncloud/ocis/v2/services/antivirus/pkg/service"
	"github.com/urfave/cli/v2"
)

// Server is the entrypoint for the server command.
func Server(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "server",
		Usage:    fmt.Sprintf("start the %s service without runtime (unsupervised mode)", cfg.Service.Name),
		Category: "server",
		Before: func(c *cli.Context) error {
			return configlog.ReturnFatal(parser.ParseConfig(cfg))
		},
		Action: func(c *cli.Context) error {
			var (
				gr          = run.Group{}
				ctx, cancel = func() (context.Context, context.CancelFunc) {
					if cfg.Context == nil {
						return context.WithCancel(context.Background())
					}
					return context.WithCancel(cfg.Context)
				}()
				logger = log.NewLogger(
					log.Name(cfg.Service.Name),
					log.Level(cfg.Log.Level),
					log.Pretty(cfg.Log.Pretty),
					log.Color(cfg.Log.Color),
					log.File(cfg.Log.File),
				)
			)
			defer cancel()

			{
				svc, err := service.NewAntivirus(cfg, logger)
				if err != nil {
					return err
				}

				gr.Add(svc.Run, func(_ error) {
					cancel()
				})
			}

			{
				server := debug.NewService(
					debug.Logger(logger),
					debug.Name(cfg.Service.Name),
					debug.Version(version.GetString()),
					debug.Address(cfg.Debug.Addr),
					debug.Token(cfg.Debug.Token),
					debug.Pprof(cfg.Debug.Pprof),
					debug.Zpages(cfg.Debug.Zpages),
					debug.Health(
						func(w http.ResponseWriter, r *http.Request) {
							w.Header().Set("Content-Type", "text/plain")
							w.WriteHeader(http.StatusOK)

							// TODO: check if services are up and running

							_, err := io.WriteString(w, http.StatusText(http.StatusOK))
							// io.WriteString should not fail but if it does we want to know.
							if err != nil {
								panic(err)
							}
						},
					),
					debug.Ready(
						func(w http.ResponseWriter, r *http.Request) {
							w.Header().Set("Content-Type", "text/plain")
							w.WriteHeader(http.StatusOK)

							// TODO: check if services are up and running

							_, err := io.WriteString(w, http.StatusText(http.StatusOK))
							// io.WriteString should not fail but if it does we want to know.
							if err != nil {
								panic(err)
							}
						},
					),
				)

				gr.Add(server.ListenAndServe, func(_ error) {
					_ = server.Shutdown(ctx)
					cancel()
				})
			}

			return gr.Run()
		},
	}
}
