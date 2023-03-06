package main

import (
	"os"

	_ "github.com/KimMachineGun/automemlimit"
	"github.com/owncloud/ocis/v2/ocis/pkg/command"
)

func main() {
	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
