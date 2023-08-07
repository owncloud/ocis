package command

import (
	"time"

	"github.com/cs3org/reva/v2/pkg/events"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/configlog"
	"github.com/owncloud/ocis/v2/services/storage-users/pkg/config"
	"github.com/owncloud/ocis/v2/services/storage-users/pkg/config/parser"
	"github.com/owncloud/ocis/v2/services/storage-users/pkg/event"
	"github.com/urfave/cli/v2"
)

// TrashBin wraps trash-bin related sub-commands.
func TrashBin(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "trash-bin",
		Usage: "manage trash-bin's",
		Subcommands: []*cli.Command{
			PurgeExpiredResources(cfg),
		},
	}
}

// PurgeExpiredResources cli command removes old trash-bin items.
func PurgeExpiredResources(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "purge-expired",
		Usage: "Purge expired trash-bin items",
		Flags: []cli.Flag{},
		Before: func(c *cli.Context) error {
			return configlog.ReturnFatal(parser.ParseConfig(cfg))
		},
		Action: func(c *cli.Context) error {
			stream, err := event.NewStream(cfg)
			if err != nil {
				return err
			}

			if err := events.Publish(c.Context, stream, event.PurgeTrashBin{ExecutionTime: time.Now()}); err != nil {
				return err
			}

			// go-micro nats implementation uses async publishing,
			// therefore we need to manually wait.
			//
			// FIXME: upstream pr
			//
			// https://github.com/go-micro/plugins/blob/3e77393890683be4bacfb613bc5751867d584692/v4/events/natsjs/nats.go#L115
			time.Sleep(5 * time.Second)

			return nil
		},
	}
}
