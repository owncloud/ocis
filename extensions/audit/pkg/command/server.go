package command

import (
	"context"
	"fmt"

	"github.com/asim/go-micro/plugins/events/natsjs/v4"
	"github.com/cs3org/reva/v2/pkg/events"
	"github.com/cs3org/reva/v2/pkg/events/server"
	"github.com/owncloud/ocis/extensions/audit/pkg/config"
	"github.com/owncloud/ocis/extensions/audit/pkg/config/parser"
	"github.com/owncloud/ocis/extensions/audit/pkg/logging"
	svc "github.com/owncloud/ocis/extensions/audit/pkg/service"
	"github.com/owncloud/ocis/extensions/audit/pkg/types"
	"github.com/urfave/cli/v2"
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

			ctx := cfg.Context
			if ctx == nil {
				ctx = context.Background()
			}
			ctx, cancel := context.WithCancel(ctx)
			defer cancel()

			evtsCfg := cfg.Events
			client, err := server.NewNatsStream(
				natsjs.Address(evtsCfg.Endpoint),
				natsjs.ClusterID(evtsCfg.Cluster),
			)
			if err != nil {
				return err
			}
			evts, err := events.Consume(client, evtsCfg.ConsumerGroup, types.RegisteredEvents()...)
			if err != nil {
				return err
			}

			svc.AuditLoggerFromConfig(ctx, cfg.Auditlog, evts, logger)
			return nil
		},
	}
}
