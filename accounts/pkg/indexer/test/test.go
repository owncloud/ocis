package test

import (
	"context"
	"flag"
	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis/storage/pkg/command"
	mcfg "github.com/owncloud/ocis/storage/pkg/config"
)

func init() {
	go setupMetadataStorage()
}

func setupMetadataStorage() {
	cfg := mcfg.New()
	app := cli.App{
		Name:     "storage-metadata-for-tests",
		Commands: []*cli.Command{command.StorageMetadata(cfg)},
	}

	_ = app.Command("storage-metadata").Run(cli.NewContext(&app, &flag.FlagSet{}, &cli.Context{Context: context.Background()}))
}
