package main

import (
	"fmt"
	"os"

	"github.com/owncloud/ocis/v2/ocis/pkg/command"
)

func main() {
	if err := command.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
