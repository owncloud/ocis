package main

import (
	"github.com/owncloud/ocis/v2/services/authz/pkg/command"
	"github.com/owncloud/ocis/v2/services/authz/pkg/config/defaults"
	"os"
)

func main() {
	if err := command.Execute(defaults.DefaultConfig()); err != nil {
		os.Exit(1)
	}
}
