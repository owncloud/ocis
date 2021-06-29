package flagset

import (
	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis/storage/pkg/config"
)

// RuntimeConfig applies common debug config cfg to the flagset
func RuntimeConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:  "extensions",
			Usage: "Run specific extensions during supervised mode",
		},
	}
}
