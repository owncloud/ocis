package command

import (
	"github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/shared"
	"github.com/urfave/cli/v2"
)

func ParseStorageCommon(ctx *cli.Context, cfg *config.Config) error {
	if err := ParseConfig(ctx, cfg); err != nil {
		return err
	}

	if cfg.Storage.Log == nil && cfg.Commons != nil && cfg.Commons.Log != nil {
		cfg.Storage.Log = &shared.Log{
			Level:  cfg.Commons.Log.Level,
			Pretty: cfg.Commons.Log.Pretty,
			Color:  cfg.Commons.Log.Color,
			File:   cfg.Commons.Log.File,
		}
	} else if cfg.Storage.Log == nil && cfg.Commons == nil {
		cfg.Storage.Log = &shared.Log{}
	}

	return nil
}
