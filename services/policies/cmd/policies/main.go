package main

import (
	"os"

	"github.com/owncloud/ocis/v2/services/policies/pkg/command"
	"github.com/owncloud/ocis/v2/services/policies/pkg/config/defaults"
)

func main() {
	if err := command.Execute(defaults.DefaultConfig()); err != nil {
		os.Exit(1)
	}
}
