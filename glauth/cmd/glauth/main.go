package main

import (
	"os"

	"github.com/owncloud/ocis/glauth/pkg/command"
	"github.com/owncloud/ocis/glauth/pkg/config/defaults"
)

func main() {
	if err := command.Execute(defaults.DefaultConfig()); err != nil {
		os.Exit(1)
	}
}
