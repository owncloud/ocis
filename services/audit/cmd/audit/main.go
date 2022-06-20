package main

import (
	"os"

	"github.com/owncloud/ocis/v2/extensions/audit/pkg/command"
	"github.com/owncloud/ocis/v2/extensions/audit/pkg/config/defaults"
)

func main() {
	if err := command.Execute(defaults.DefaultConfig()); err != nil {
		os.Exit(1)
	}
}
