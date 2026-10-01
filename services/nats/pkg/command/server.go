package command

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"os/signal"

	"github.com/urfave/cli/v2"

	"github.com/owncloud/ocis/v2/ocis-pkg/config/configlog"
	pkgcrypto "github.com/owncloud/ocis/v2/ocis-pkg/crypto"
	"github.com/owncloud/ocis/v2/ocis-pkg/runner"
	"github.com/owncloud/ocis/v2/services/nats/pkg/config"
	"github.com/owncloud/ocis/v2/services/nats/pkg/config/parser"
	"github.com/owncloud/ocis/v2/services/nats/pkg/logging"
	"github.com/owncloud/ocis/v2/services/nats/pkg/server/debug"
	"github.com/owncloud/ocis/v2/services/nats/pkg/server/nats"
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
			logger := logging.Configure(cfg.Service.Name, cfg.Log)

			var cancel context.CancelFunc
			if cfg.Context == nil {
				cfg.Context, cancel = signal.NotifyContext(context.Background(), runner.StopSignals...)
				defer cancel()
			}
			ctx := cfg.Context

			gr := runner.NewGroup()
			{
				debugServer, err := debug.Server(
					debug.Logger(logger),
					debug.Context(ctx),
					debug.Config(cfg),
				)
				if err != nil {
					logger.Info().Err(err).Str("transport", "debug").Msg("Failed to initialize server")
					return err
				}

				gr.Add(runner.NewGolangHttpServerRunner(cfg.Service.Name+".debug", debugServer))
			}

			var tlsConf *tls.Config
			if cfg.Nats.EnableTLS {
				// Generate a self-signing cert if no certificate is present
				if err := pkgcrypto.GenCert(cfg.Nats.TLSCert, cfg.Nats.TLSKey, logger); err != nil {
					logger.Fatal().Err(err).Msgf("Could not generate test-certificate")
				}

				crt, err := tls.LoadX509KeyPair(cfg.Nats.TLSCert, cfg.Nats.TLSKey)
				if err != nil {
					return err
				}

				clientAuth := tls.RequireAndVerifyClientCert
				if cfg.Nats.TLSSkipVerifyClientCert {
					clientAuth = tls.NoClientCert
				}

				tlsConf = &tls.Config{
					MinVersion:   tls.VersionTLS12,
					ClientAuth:   clientAuth,
					Certificates: []tls.Certificate{crt},
				}
			}
			authConfigured := cfg.Nats.AuthUsername != "" || cfg.Nats.AuthPassword != ""
			if !authConfigured && !isLoopbackHost(cfg.Nats.Host) {
				logger.Warn().
					Str("host", cfg.Nats.Host).
					Msg("the NATS event bus is bound to a non-loopback address without authentication: " +
						"any client able to reach this port can publish forged internal events. " +
						"Set OCIS_EVENTS_AUTH_USERNAME and OCIS_EVENTS_AUTH_PASSWORD (and enable TLS), " +
						"or bind the broker to a loopback/isolated network only")
			}

			natsServer, err := nats.NewNATSServer(
				logging.NewLogWrapper(logger),
				nats.Host(cfg.Nats.Host),
				nats.Port(cfg.Nats.Port),
				nats.ClusterID(cfg.Nats.ClusterID),
				nats.StoreDir(cfg.Nats.StoreDir),
				nats.TLSConfig(tlsConf),
				nats.AllowNonTLS(!cfg.Nats.EnableTLS),
				nats.Auth(cfg.Nats.AuthUsername, cfg.Nats.AuthPassword),
			)
			if err != nil {
				return err
			}

			gr.Add(runner.New(cfg.Service.Name+".svc", func() error {
				return natsServer.ListenAndServe()
			}, func() {
				natsServer.Shutdown()
			}))

			logger.Warn().Msgf("starting service %s", cfg.Service.Name)
			grResults := gr.Run(ctx)

			if err := runner.ProcessResults(grResults); err != nil {
				logger.Error().Err(err).Msgf("service %s stopped with error", cfg.Service.Name)
				return err
			}
			logger.Warn().Msgf("service %s stopped without error", cfg.Service.Name)
			return nil
		},
	}
}

// isLoopbackHost reports whether host refers to the local loopback interface,
// in which case binding the broker without authentication is acceptable.
// Unparsable or non-IP hosts are treated as non-loopback so the warning errs
// on the side of caution.
func isLoopbackHost(host string) bool {
	if host == "localhost" {
		return true
	}
	// An empty host makes the NATS server listen on all interfaces, so it must
	// not be treated as loopback.
	if ip := net.ParseIP(host); ip != nil {
		return ip.IsLoopback()
	}
	return false
}
