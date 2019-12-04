package main

import (
	"os"

	"github.com/owncloud/ocis-graph/pkg/command"
)

func main() {
	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
