package main

import (
	"os"

	"github.com/owncloud/ocis/v2/extensions/ocdav/pkg/command"
	"github.com/owncloud/ocis/v2/extensions/ocdav/pkg/config/defaults"
)

func main() {
	if err := command.Execute(defaults.DefaultConfig()); err != nil {
		os.Exit(1)
	}
}
