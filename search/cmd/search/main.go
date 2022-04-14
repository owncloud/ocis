package main

import (
	"os"

	"github.com/owncloud/ocis/search/pkg/command"
	"github.com/owncloud/ocis/search/pkg/config/defaults"
)

func main() {
	if err := command.Execute(defaults.DefaultConfig()); err != nil {
		os.Exit(1)
	}
}
