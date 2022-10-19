package main

import (
	"os"

	"github.com/KimMachineGun/automemlimit/memlimit"
	"github.com/owncloud/ocis/v2/ocis/pkg/command"
)

func main() {

	if _, present := os.LookupEnv("GOMEMLIMIT"); !present {
		memlimit.SetGoMemLimit(0.9)
	}

	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
