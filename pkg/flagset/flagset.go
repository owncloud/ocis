package flagset

import (
	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis-accounts/pkg/config"
)

// RootWithConfig applies cfg to the root flagset
func RootWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "log-level",
			Value:       "info",
			Usage:       "Set logging level",
			EnvVars:     []string{"ACCOUNTS_LOG_LEVEL"},
			Destination: &cfg.Log.Level,
		},
		&cli.BoolFlag{
			Value:       true,
			Name:        "log-pretty",
			Usage:       "Enable pretty logging",
			EnvVars:     []string{"ACCOUNTS_LOG_PRETTY"},
			Destination: &cfg.Log.Pretty,
		},
		&cli.BoolFlag{
			Value:       true,
			Name:        "log-color",
			Usage:       "Enable colored logging",
			EnvVars:     []string{"ACCOUNTS_LOG_COLOR"},
			Destination: &cfg.Log.Color,
		},
	}
}

// ServerWithConfig applies cfg to the root flagset
func ServerWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "name",
			Value:       "accounts",
			DefaultText: "accounts",
			Usage:       "service name",
			EnvVars:     []string{"ACCOUNTS_NAME"},
			Destination: &cfg.Server.Name,
		},
		&cli.StringFlag{
			Name:        "namespace",
			Aliases:     []string{"ns"},
			Value:       "com.owncloud.api",
			DefaultText: "com.owncloud.api",
			Usage:       "namespace",
			EnvVars:     []string{"ACCOUNTS_NAMESPACE"},
			Destination: &cfg.Server.Namespace,
		},
		&cli.StringFlag{
			Name:        "address",
			Aliases:     []string{"addr"},
			Value:       "localhost:9180",
			DefaultText: "localhost:9180",
			Usage:       "service endpoint",
			EnvVars:     []string{"ACCOUNTS_ADDRESS"},
			Destination: &cfg.Server.Address,
		},
		&cli.StringFlag{
			Name:        "accounts-data-path",
			Value:       "/var/tmp/ocis-accounts",
			DefaultText: "/var/tmp/ocis-accounts",
			Usage:       "accounts folder",
			EnvVars:     []string{"ACCOUNTS_DATA_PATH"},
			Destination: &cfg.Server.AccountsDataPath,
		},
	}
}
