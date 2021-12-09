package main

import (
	"os"

	"github.com/owncloud/ocis/idp/pkg/command"
	"github.com/owncloud/ocis/idp/pkg/config"
)

func main() {
	if err := command.Execute(config.DefaultConfig()); err != nil {
		os.Exit(1)
	}
}
