package main

import (
	"os"

	"github.com/owncloud/ocis/ocs/pkg/command"
	"github.com/owncloud/ocis/ocs/pkg/config"
)

func main() {
	if err := command.Execute(config.New()); err != nil {
		os.Exit(1)
	}
}
