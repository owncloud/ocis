package main

import (
	"os"

	"github.com/owncloud/ocis/nats/pkg/command"
	"github.com/owncloud/ocis/nats/pkg/config/defaults"
)

func main() {
	if err := command.Execute(defaults.DefaultConfig()); err != nil {
		os.Exit(1)
	}
}
