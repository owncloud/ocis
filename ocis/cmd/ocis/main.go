package main

import (
	"fmt"
	"os"

	"github.com/owncloud/ocis/v2/ocis/pkg/command"
)

func main() {
	if err := command.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
