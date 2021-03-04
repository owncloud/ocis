package command

import (
	"context"
	"os"
	"strings"

	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis/ocis-pkg/log"
	"github.com/owncloud/ocis/webdav/pkg/config"
	"github.com/owncloud/ocis/webdav/pkg/version"
	"github.com/spf13/viper"
)

// Execute is the entry point for the ocis-webdav command.
func Execute(cfg *config.Config) error {
	app := &cli.App{
		Name:     "webdav",
		Version:  version.String,
		Usage:    "Serve WebDAV API for oCIS",
		Compiled: version.Compiled(),

		Authors: []*cli.Author{
			{
				Name:  "ownCloud GmbH",
				Email: "support@owncloud.com",
			},
		},
		Before: func(c *cli.Context) error {
			cfg.Service.Version = version.String
			return nil
		},

		Commands: []*cli.Command{
			Server(cfg),
			Health(cfg),
		},
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

// NewLogger initializes a service-specific logger instance.
func NewLogger(cfg *config.Config) log.Logger {
	return log.NewLogger(
		log.Name("webdav"),
		log.Level(cfg.Log.Level),
		log.Pretty(cfg.Log.Pretty),
		log.Color(cfg.Log.Color),
	)
}

// ParseConfig loads webdav configuration from Viper known paths.
func ParseConfig(c *cli.Context, cfg *config.Config) error {
	logger := NewLogger(cfg)

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetEnvPrefix("WEBDAV")
	viper.AutomaticEnv()

	if c.IsSet("config-file") {
		viper.SetConfigFile(c.String("config-file"))
	} else {
		viper.SetConfigName("webdav")

		viper.AddConfigPath("/etc/ocis")
		viper.AddConfigPath("$HOME/.ocis")
		viper.AddConfigPath("./config")
	}

	if err := viper.ReadInConfig(); err != nil {
		switch err.(type) {
		case viper.ConfigFileNotFoundError:
			logger.Debug().
				Msg("no config found on preconfigured location")
		case viper.UnsupportedConfigError:
			logger.Fatal().
				Err(err).
				Msg("Unsupported config type")
		default:
			logger.Fatal().
				Err(err).
				Msg("Failed to read config")
		}
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		logger.Fatal().
			Err(err).
			Msg("Failed to parse config")
	}

	return nil
}

// SutureService allows for the webdav command to be embedded and supervised by a suture supervisor tree.
type SutureService struct {
	ctx    context.Context
	cancel context.CancelFunc // used to cancel the context go-micro services used to shutdown a service.
	cfg    *config.Config
}

// NewSutureService creates a new webdav.SutureService
func NewSutureService(ctx context.Context, cfg *config.Config) SutureService {
	sctx, cancel := context.WithCancel(ctx)
	cfg.Context = sctx // propagate the context down to the go-micro services.
	return SutureService{
		ctx:    sctx,
		cancel: cancel,
		cfg:    cfg,
	}
}

func (s SutureService) Serve() {
	if err := Execute(s.cfg); err != nil {
		return
	}
}

func (s SutureService) Stop() {
	s.cancel()
}
