package command

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/owncloud/ocis/v2/ocis-pkg/config/configlog"
	"github.com/owncloud/ocis/v2/services/search/pkg/config"
	"github.com/owncloud/ocis/v2/services/search/pkg/config/parser"
	"github.com/owncloud/ocis/v2/services/search/pkg/engine"
)

// Optimize is the entrypoint for the optimize command.
func Optimize(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "optimize",
		Usage:    "compact the search index by merging segments, without re-indexing content",
		Category: "index management",
		Before: func(_ *cli.Context) error {
			return configlog.ReturnFatal(parser.ParseConfig(cfg))
		},
		Action: func(_ *cli.Context) error {
			eng, closer, err := engine.NewEngineFromConfig(cfg)
			if err != nil {
				return err
			}
			defer closer.Close()

			fmt.Println("optimizing search index...")
			if err := eng.Optimize(context.Background()); err != nil {
				fmt.Println("failed to optimize index: " + err.Error())
				return err
			}
			fmt.Println("index optimization complete")
			return nil
		},
	}
}
