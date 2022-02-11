package command

import (
	"fmt"

	"github.com/cs3org/reva/pkg/events/server"
	"github.com/owncloud/ocis/nats/pkg/config"
	"github.com/owncloud/ocis/nats/pkg/config/parser"
	"github.com/owncloud/ocis/nats/pkg/logging"
	"github.com/owncloud/ocis/ocis-pkg/log"
	"github.com/urfave/cli/v2"

	// TODO: .Logger Option on events/server would make this import redundant
	stanServer "github.com/nats-io/nats-streaming-server/server"
)

// Server is the entrypoint for the server command.
func Server(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "server",
		Usage:    fmt.Sprintf("start %s extension without runtime (unsupervised mode)", cfg.Service.Name),
		Category: "server",
		Before: func(c *cli.Context) error {
			return parser.ParseConfig(cfg)
		},
		Action: func(c *cli.Context) error {
			logger := logging.Configure(cfg.Service.Name, cfg.Log)
			err := server.RunNatsServer(server.Host(cfg.Nats.Host), server.Port(cfg.Nats.Port), server.StanOpts(func(o *stanServer.Options) {
				o.CustomLogger = &logWrapper{logger}
			}))
			if err != nil {
				return err
			}
			for {
			}
		},
	}
}

// we need to wrap our logger so we can pass it to the nats server
type logWrapper struct {
	logger log.Logger
}

// Noticef logs a notice statement
func (l *logWrapper) Noticef(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	l.logger.Info().Msg(msg)
}

// Warnf logs a warning statement
func (l *logWrapper) Warnf(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	l.logger.Warn().Msg(msg)
}

// Fatalf logs a fatal statement
func (l *logWrapper) Fatalf(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	l.logger.Fatal().Msg(msg)
}

// Errorf logs an error statement
func (l *logWrapper) Errorf(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	l.logger.Error().Msg(msg)
}

// Debugf logs a debug statement
func (l *logWrapper) Debugf(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	l.logger.Debug().Msg(msg)
}

// Tracef logs a trace statement
func (l *logWrapper) Tracef(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	l.logger.Trace().Msg(msg)
}
