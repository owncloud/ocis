package main

import (
	"os"

	"github.com/owncloud/ocis/extensions/nats/pkg/command"
	"github.com/owncloud/ocis/extensions/nats/pkg/config/defaults"
)

func main() {
	if err := command.Execute(defaults.DefaultConfig()); err != nil {
		os.Exit(1)
	}
}
