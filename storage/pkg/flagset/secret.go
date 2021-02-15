package flagset

import (
	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis/storage/pkg/config"
)

// SecretWithConfig applies cfg to the root flagset
func SecretWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "jwt-secret",
			Value:       "Pive-Fumkiu4",
			Usage:       "Shared jwt secret for reva service communication",
			EnvVars:     []string{"STORAGE_JWT_SECRET", "OCIS_JWT_SECRET"},
			Destination: &cfg.Reva.JWTSecret,
		},
	}
}
