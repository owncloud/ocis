package command

import (
	"context"
	"os"
	"strings"

	"github.com/owncloud/ocis/ocis-pkg/sync"

	"github.com/owncloud/ocis/accounts/pkg/flagset"

	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis/accounts/pkg/config"
	"github.com/owncloud/ocis/accounts/pkg/version"
	ociscfg "github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/log"
	"github.com/spf13/viper"
	"github.com/thejerf/suture/v4"
)

var (
	defaultConfigPaths = []string{"/etc/ocis", "$HOME/.ocis", "./config"}
	defaultFilename    = "accounts"
)

// Execute is the entry point for the ocis-accounts command.
func Execute(cfg *config.Config) error {
	app := &cli.App{
		Name:     "ocis-accounts",
		Version:  version.String,
		Usage:    "Provide accounts and groups for oCIS",
		Compiled: version.Compiled(),

		Authors: []*cli.Author{
			{
				Name:  "ownCloud GmbH",
				Email: "support@owncloud.com",
			},
		},

		Flags: flagset.RootWithConfig(cfg),

		Before: func(c *cli.Context) error {
			cfg.Server.Version = version.String
			return ParseConfig(c, cfg)
		},

		Commands: []*cli.Command{
			Server(cfg),
			AddAccount(cfg),
			UpdateAccount(cfg),
			ListAccounts(cfg),
			InspectAccount(cfg),
			RemoveAccount(cfg),
			PrintVersion(cfg),
			RebuildIndex(cfg),
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
		log.Name("accounts"),
		log.Level(cfg.Log.Level),
		log.Pretty(cfg.Log.Pretty),
		log.Color(cfg.Log.Color),
	)
}

// ParseConfig loads accounts configuration from Viper known paths.
func ParseConfig(c *cli.Context, cfg *config.Config) error {
	sync.ParsingViperConfig.Lock()
	defer sync.ParsingViperConfig.Unlock()
	logger := NewLogger(cfg)

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetEnvPrefix("ACCOUNTS")
	viper.AutomaticEnv()

	if c.IsSet("config-file") {
		viper.SetConfigFile(c.String("config-file"))
	} else {
		viper.SetConfigName(defaultFilename)

		for _, v := range defaultConfigPaths {
			viper.AddConfigPath(v)
		}
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

// SutureService allows for the accounts command to be embedded and supervised by a suture supervisor tree.
type SutureService struct {
	cfg *config.Config
}

// NewSutureService creates a new accounts.SutureService
func NewSutureService(cfg *ociscfg.Config) suture.Service {
	if cfg.Mode == 0 {
		cfg.Accounts.Supervised = true
	}
	return SutureService{
		cfg: cfg.Accounts,
	}
}

func (s SutureService) Serve(ctx context.Context) error {
	s.cfg.Context = ctx
	if err := Execute(s.cfg); err != nil {
		return err
	}

	return nil
}
