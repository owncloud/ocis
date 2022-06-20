package main

import (
	"os"

	"github.com/owncloud/ocis/v2/services/nats/pkg/command"
	"github.com/owncloud/ocis/v2/services/nats/pkg/config/defaults"
)

func main() {
	if err := command.Execute(defaults.DefaultConfig()); err != nil {
		os.Exit(1)
	}
}
