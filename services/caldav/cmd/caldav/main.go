package main

import (
	"os"

	"github.com/owncloud/ocis/v2/services/caldav/pkg/command"
	"github.com/owncloud/ocis/v2/services/caldav/pkg/config/defaults"
)

func main() {
	if err := command.Execute(defaults.DefaultConfig()); err != nil {
		os.Exit(1)
	}
}
