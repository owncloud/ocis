package main

import (
	"os"

	"github.com/owncloud/reva-phoenix/service"
)

func main() {

	revaphoenixCommand := service.NewRevaPhoenixCommand("reva-phoenix")

	if err := revaphoenixCommand.Execute(); err != nil {
		os.Exit(1)
	}

}
