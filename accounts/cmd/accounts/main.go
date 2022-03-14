package main

import (
	"os"

	"github.com/owncloud/ocis/accounts/pkg/command"
	"github.com/owncloud/ocis/accounts/pkg/config/defaults"
)

func main() {
	if err := command.Execute(defaults.DefaultConfig()); err != nil {
		os.Exit(1)
	}
}
