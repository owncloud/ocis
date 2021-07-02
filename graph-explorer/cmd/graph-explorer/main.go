package main

import (
	"os"

	"github.com/owncloud/ocis/graph-explorer/pkg/command"
	"github.com/owncloud/ocis/graph-explorer/pkg/config"
)

func main() {
	if err := command.Execute(config.New()); err != nil {
		os.Exit(1)
	}
}
