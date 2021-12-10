package main

import (
	"os"

	"github.com/owncloud/ocis/webdav/pkg/command"
	"github.com/owncloud/ocis/webdav/pkg/config"
)

func main() {
	if err := command.Execute(config.DefaultConfig()); err != nil {
		os.Exit(1)
	}
}
