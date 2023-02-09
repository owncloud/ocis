package command

import (
	"fmt"
	"time"

	"github.com/cs3org/reva/v2/pkg/events"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/configlog"
	"github.com/owncloud/ocis/v2/services/storage-users/pkg/config"
	"github.com/owncloud/ocis/v2/services/storage-users/pkg/config/parser"
	"github.com/owncloud/ocis/v2/services/storage-users/pkg/event"
	"github.com/urfave/cli/v2"
)

var (
	flagUserID = &cli.StringFlag{
		Name:  "user-id",
		Usage: "The user-id of the user who should be used to list and delete expired space trash-bin items",
	}
	flagPurgeBefore = &cli.StringFlag{
		Name:  "purge-before",
		Usage: "Specifies the period of time in which items that have been in the trash bin for longer than this should be deleted",
		Value: "720h", // 30 days
	}
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
		Flags: []cli.Flag{
			flagUserID,
			flagPurgeBefore,
		},
		Before: func(c *cli.Context) error {
			err := parser.ParseConfig(cfg)

			if c.Value(flagUserID.Name) == "" {
				_ = c.Set(flagUserID.Name, cfg.Commons.AdminUserID)
			}

			return configlog.ReturnFatal(err)
		},
		Action: func(c *cli.Context) error {
			userID := c.String(flagUserID.Name)
			if userID == "" {
				return cli.Exit(fmt.Errorf("%s must be set", flagUserID.Name), 1)
			}

			purgeBefore, err := time.ParseDuration(c.String(flagPurgeBefore.Name))
			if err != nil {
				return cli.Exit(err, 1)
			}

			stream, err := event.NewStream(cfg.Events)
			if err != nil {
				return err
			}

			if err := events.Publish(stream, event.PurgeTrashBin{
				ExecutantID:  userID,
				RemoveBefore: time.Now().Add(-purgeBefore),
			}); err != nil {
				return err
			}

			// go-micro nats implementation uses async publishing,
			// therefor we need to manually wait.
			//
			// fixMe: upstream pr
			//
			// https://github.com/go-micro/plugins/blob/3e77393890683be4bacfb613bc5751867d584692/v4/events/natsjs/nats.go#L115
			time.Sleep(5 * time.Second)

			return nil
		},
	}
}
