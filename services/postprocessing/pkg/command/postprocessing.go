package command

import (
	"context"

	"github.com/owncloud/ocis/v2/ocis-pkg/config/configlog"
	"github.com/owncloud/ocis/v2/ocis-pkg/generators"
	"github.com/owncloud/ocis/v2/services/postprocessing/pkg/config"
	"github.com/owncloud/ocis/v2/services/postprocessing/pkg/config/parser"
	"github.com/owncloud/reva/v2/pkg/events"
	"github.com/owncloud/reva/v2/pkg/events/stream"
	"github.com/owncloud/reva/v2/pkg/utils"
	"github.com/urfave/cli/v2"
)

// RestartPostprocessing cli command to restart postprocessing
func RestartPostprocessing(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:    "resume",
		Aliases: []string{"restart"},
		Usage:   "resume postprocessing for an uploadID",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "upload-id",
				Aliases: []string{"u"},
				Usage:   "the uploadid to resume. Ignored if unset.",
			},
			&cli.StringFlag{
				Name:    "step",
				Aliases: []string{"s"},
				Usage:   "resume all uploads in the given postprocessing step. Ignored if upload-id is set.",
				Value:   "finished",
			},
			&cli.BoolFlag{
				Name:    "restart",
				Aliases: []string{"r"},
				Usage:   "restart postprocessing for the given uploadID. Ignores the step flag.",
			},
		},
		Before: func(c *cli.Context) error {
			return configlog.ReturnFatal(parser.ParseConfig(cfg))
		},
		Action: func(c *cli.Context) error {
			connName := generators.GenerateConnectionName(cfg.Service.Name, generators.NTypeBus)
			stream, err := stream.NatsFromConfig(connName, false, stream.NatsConfig(cfg.Postprocessing.Events))
			if err != nil {
				return err
			}

			uid, step := c.String("upload-id"), ""
			if uid == "" {
				step = c.String("step")
			}

			var ev events.Unmarshaller
			switch {
			case c.Bool("retrigger"):
				ev = events.RestartPostprocessing{
					UploadID:  uid,
					Timestamp: utils.TSNow(),
				}
			default:
				ev = events.ResumePostprocessing{
					UploadID:  uid,
					Step:      events.Postprocessingstep(step),
					Timestamp: utils.TSNow(),
				}
			}

			return events.Publish(context.Background(), stream, ev)
		},
	}
}
