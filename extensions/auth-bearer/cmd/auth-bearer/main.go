package main

import (
	"os"

	"github.com/owncloud/ocis/extensions/auth-bearer/pkg/command"
	"github.com/owncloud/ocis/extensions/auth-bearer/pkg/config/defaults"
)

func main() {
	if err := command.Execute(defaults.DefaultConfig()); err != nil {
		os.Exit(1)
	}
}
