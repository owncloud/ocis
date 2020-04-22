package flagset

import "github.com/micro/cli/v2"
import "github.com/owncloud/ocis-accounts/pkg/config"

// ServerWithConfig applies cfg to the root flagset
func ServerWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "manager",
			DefaultText: "filesystem",
			Usage:       "accounts backend manager",
			Value:       "filesystem",
			EnvVars:     []string{"ACCOUNTS_MANAGER"},
			Destination: &cfg.Manager,
		},
		&cli.StringFlag{
			Name:        "mount-path",
			Usage:       "mounting point (necessary when manager=filesystem)",
			EnvVars:     []string{"ACCOUNTS_MOUNT_PATH"},
			Destination: &cfg.MountPath,
		},
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
			Value:       "com.owncloud",
			DefaultText: "com.owncloud",
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
	}
}
