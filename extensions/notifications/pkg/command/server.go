package command

import (
	"fmt"

	"github.com/cs3org/reva/v2/pkg/events"
	"github.com/cs3org/reva/v2/pkg/events/server"
	"github.com/go-micro/plugins/v4/events/natsjs"
	"github.com/owncloud/ocis/v2/extensions/notifications/pkg/channels"
	"github.com/owncloud/ocis/v2/extensions/notifications/pkg/config"
	"github.com/owncloud/ocis/v2/extensions/notifications/pkg/config/parser"
	"github.com/owncloud/ocis/v2/extensions/notifications/pkg/logging"
	"github.com/owncloud/ocis/v2/extensions/notifications/pkg/service"
	"github.com/urfave/cli/v2"
)

// Server is the entrypoint for the server command.
func Server(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "server",
		Usage:    fmt.Sprintf("start %s extension without runtime (unsupervised mode)", cfg.Service.Name),
		Category: "server",
		Before: func(c *cli.Context) error {
			err := parser.ParseConfig(cfg)
			if err != nil {
				fmt.Printf("%v", err)
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
