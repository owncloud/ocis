package command

import (
	"context"

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
				Name:    "upload-id",
				Aliases: []string{"u"},
				Usage:   "the uploadid to restart. Ignored if unset.",
			},
			&cli.StringFlag{
				Name:    "step",
				Aliases: []string{"s"},
				Usage:   "restarts all uploads in the given postprocessing step. Ignored if upload-id is set.",
				Value:   "finished", // Calling `ocis postprocessing restart` without any arguments will restart all uploads that are finished but failed to move the uploed from the upload area to the blobstore.
			},
		},
		Before: func(c *cli.Context) error {
			return configlog.ReturnFatal(parser.ParseConfig(cfg))
		},
		Action: func(c *cli.Context) error {
			stream, err := stream.NatsFromConfig(cfg.Service.Name, false, stream.NatsConfig(cfg.Postprocessing.Events))
			if err != nil {
				return err
			}

			uid, step := c.String("upload-id"), ""
			if uid == "" {
				step = c.String("step")
			}

			ev := events.ResumePostprocessing{
				UploadID:  uid,
				Step:      events.Postprocessingstep(step),
				Timestamp: utils.TSNow(),
			}

			return events.Publish(context.Background(), stream, ev)
		},
	}
}
