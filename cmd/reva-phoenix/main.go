package main

import (
	"os"

	"github.com/owncloud/reva-phoenix/pkg/service"
)

func main() {
	if err := service.RootCommand().Execute(); err != nil {
		os.Exit(1)
	}
}
