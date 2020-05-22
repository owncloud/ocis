package main

import (
	"os"

	"github.com/owncloud/ocis-settings/pkg/command"
)

func main() {
	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
