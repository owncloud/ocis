package flagset

import (
	"github.com/owncloud/ocis/ocis-pkg/flags"
	"github.com/owncloud/ocis/storage/pkg/config"
	"github.com/urfave/cli/v2"
)

// SecretWithConfig applies cfg to the root flagset
func SecretWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "jwt-secret",
			Value:       flags.OverrideDefaultString(cfg.Reva.JWTSecret, "Pive-Fumkiu4"),
			Usage:       "Shared jwt secret for reva service communication",
			EnvVars:     []string{"STORAGE_JWT_SECRET", "OCIS_JWT_SECRET"},
			Destination: &cfg.Reva.JWTSecret,
		},
		&cli.BoolFlag{
			Name:        "skip-user-groups-in-token",
			Value:       flags.OverrideDefaultBool(cfg.Reva.SkipUserGroupsInToken, false),
			Usage:       "Whether to skip encoding user groups in reva's JWT token",
			EnvVars:     []string{"STORAGE_SKIP_USER_GROUPS_IN_TOKEN"},
			Destination: &cfg.Reva.SkipUserGroupsInToken,
		},
	}
}
