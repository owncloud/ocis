package main

import (
	"os"

	"github.com/owncloud/ocis/glauth/pkg/command"
	"github.com/owncloud/ocis/glauth/pkg/config"
)

func main() {
	if err := command.Execute(config.New()); err != nil {
		os.Exit(1)
	}
}
