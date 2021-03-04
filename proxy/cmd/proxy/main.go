package main

import (
	"os"

	"github.com/owncloud/ocis/proxy/pkg/command"
	"github.com/owncloud/ocis/proxy/pkg/config"
)

func main() {
	if err := command.Execute(config.New()); err != nil {
		os.Exit(1)
	}
}
