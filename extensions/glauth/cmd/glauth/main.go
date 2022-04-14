package main

import (
	"os"

	"github.com/owncloud/ocis/extensions/glauth/pkg/command"
	"github.com/owncloud/ocis/extensions/glauth/pkg/config/defaults"
)

func main() {
	if err := command.Execute(defaults.DefaultConfig()); err != nil {
		os.Exit(1)
	}
}
