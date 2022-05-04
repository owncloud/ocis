package main

import (
	"os"

	"github.com/owncloud/ocis/v2/extensions/graph-explorer/pkg/command"
	"github.com/owncloud/ocis/v2/extensions/graph-explorer/pkg/config/defaults"
)

func main() {
	if err := command.Execute(defaults.DefaultConfig()); err != nil {
		os.Exit(1)
	}
}
