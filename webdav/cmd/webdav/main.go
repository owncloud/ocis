package main

import (
	"os"

	"github.com/owncloud/ocis/webdav/pkg/command"
	"github.com/owncloud/ocis/webdav/pkg/config/defaults"
)

func main() {
	if err := command.Execute(defaults.DefaultConfig()); err != nil {
		os.Exit(1)
	}
}
