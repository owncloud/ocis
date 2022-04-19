package command

import (
	"fmt"

	"github.com/asim/go-micro/plugins/events/natsjs/v4"
	"github.com/cs3org/reva/v2/pkg/events"
	"github.com/cs3org/reva/v2/pkg/events/server"
	"github.com/owncloud/ocis/extensions/notifications/pkg/channels"
	"github.com/owncloud/ocis/extensions/notifications/pkg/config"
	"github.com/owncloud/ocis/extensions/notifications/pkg/config/parser"
	"github.com/owncloud/ocis/extensions/notifications/pkg/logging"
	"github.com/owncloud/ocis/extensions/notifications/pkg/service"
	"github.com/urfave/cli/v2"
)

// Server is the entrypoint for the server command.
func Server(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "server",
		Usage:    fmt.Sprintf("start %s extension without runtime (unsupervised mode)", cfg.Service.Name),
		Category: "server",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "config-file",
				Value:       cfg.ConfigFile,
				Usage:       "config file to be loaded by the extension",
				Destination: &cfg.ConfigFile,
			},
		},
		Before: func(c *cli.Context) error {
			err := parser.ParseConfig(cfg)
			if err != nil {
				logger := logging.Configure(cfg.Service.Name, &config.Log{})
				logger.Error().Err(err).Msg("couldn't find the specified config file")
			}
			return err
		},
		Action: func(c *cli.Context) error {
			logger := logging.Configure(cfg.Service.Name, cfg.Log)

			evs := []events.Unmarshaller{
				events.ShareCreated{},
			}

			evtsCfg := cfg.Notifications.Events
			client, err := server.NewNatsStream(
				natsjs.Address(evtsCfg.Endpoint),
				natsjs.ClusterID(evtsCfg.Cluster),
			)
			if err != nil {
				return err
			}
			evts, err := events.Consume(client, evtsCfg.ConsumerGroup, evs...)
			if err != nil {
				return err
			}
			channel, err := channels.NewMailChannel(*cfg, logger)
			if err != nil {
				return err
			}
			svc := service.NewEventsNotifier(evts, channel, logger)
			return svc.Run()
		},
	}
}
