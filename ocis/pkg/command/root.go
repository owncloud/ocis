package command

import (
	"os"

	gofig "github.com/gookit/config/v2"
	gooyaml "github.com/gookit/config/v2/yaml"

	"github.com/owncloud/ocis/ocis-pkg/log"
	"github.com/urfave/cli/v2"

	"github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/version"
	"github.com/owncloud/ocis/ocis/pkg/register"
)

// Execute is the entry point for the ocis command.
func Execute() error {
	cfg := config.New()

	app := &cli.App{
		Name:     "ocis",
		Version:  version.String,
		Usage:    "ownCloud Infinite Scale Stack",
		Compiled: version.Compiled(),
		Before: func(c *cli.Context) error {
			return ParseConfig(c, cfg)
		},
		Authors: []*cli.Author{
			{
				Name:  "ownCloud GmbH",
				Email: "support@owncloud.com",
			},
		},
	}

	for _, fn := range register.Commands {
		app.Commands = append(
			app.Commands,
			fn(cfg),
		)
	}

	cli.HelpFlag = &cli.BoolFlag{
		Name:  "help,h",
		Usage: "Show the help",
	}

	cli.VersionFlag = &cli.BoolFlag{
		Name:  "version,v",
		Usage: "Print the version",
	}

	return app.Run(os.Args)
}

// NewLogger initializes a service-specific logger instance
func NewLogger(cfg *config.Config) log.Logger {
	return log.NewLogger(
		log.Name("ocis"),
		log.Level(cfg.Log.Level),
		log.Pretty(cfg.Log.Pretty),
		log.Color(cfg.Log.Color),
		log.File(cfg.Log.File),
	)
}

// ParseConfig loads ocis configuration from known paths.
func ParseConfig(c *cli.Context, cfg *config.Config) error {
	cnf := gofig.NewWithOptions("ocis", gofig.ParseEnv)
	cnf.AddDriver(gooyaml.Driver)
	err := cnf.LoadFiles("/Users/aunger/code/owncloud/ocis/ocis/pkg/command/ocis_example_config.yaml")
	if err != nil {
		return err
	}

	err = cnf.BindStruct("", cfg)
	if err != nil {
		return err
	}

	return nil
}
