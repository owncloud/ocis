package command

import (
	"context"
	"time"

	"github.com/cs3org/reva/v2/pkg/events"
	"github.com/cs3org/reva/v2/pkg/events/stream"
	"github.com/cs3org/reva/v2/pkg/utils"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/configlog"
	"github.com/owncloud/ocis/v2/services/postprocessing/pkg/config"
	"github.com/owncloud/ocis/v2/services/postprocessing/pkg/config/parser"
	"github.com/urfave/cli/v2"
)

// RestartPostprocessing cli command to restart postprocessing
func RestartPostprocessing(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "restart",
		Usage: "restart postprocessing for an uploadID",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "upload-id",
				Aliases:  []string{"u"},
				Required: true,
				Usage:    "the uploadid to restart",
			},
		},
		Before: func(c *cli.Context) error {
			return configlog.ReturnFatal(parser.ParseConfig(cfg))
		},
		Action: func(c *cli.Context) error {
			stream, err := stream.NatsFromConfig(cfg.Service.Name, stream.NatsConfig(cfg.Postprocessing.Events))
			if err != nil {
				return err
			}

			ev := events.ResumePostprocessing{
				UploadID:  c.String("upload-id"),
				Timestamp: utils.TSNow(),
			}

			if err := events.Publish(context.Background(), stream, ev); err != nil {
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
