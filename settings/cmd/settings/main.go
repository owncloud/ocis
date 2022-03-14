package main

import (
	"os"

	"github.com/owncloud/ocis/settings/pkg/command"
	"github.com/owncloud/ocis/settings/pkg/config/defaults"
)

func main() {
	if err := command.Execute(defaults.DefaultConfig()); err != nil {
		os.Exit(1)
	}
}
