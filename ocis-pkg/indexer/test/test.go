package test

import (
	"context"
	"flag"
	"net"
	"time"

	"github.com/owncloud/ocis/storage/pkg/command"
	mcfg "github.com/owncloud/ocis/storage/pkg/config"
	"github.com/urfave/cli/v2"
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

	// wait until port is open
	d := net.Dialer{Timeout: 5 * time.Second}
	conn, err := d.Dial("tcp", "localhost:9125")
	if err != nil {
		panic("timeout waiting for storage")
	}
	conn.Close()
}
