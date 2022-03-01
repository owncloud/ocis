package main

import (
	"github.com/owncloud/ocis/ocis-pkg/docs"
	"github.com/owncloud/ocis/thumbnails/pkg/config"
)

func main() {
	cfg := config.DefaultConfig()
	docs.Display(*cfg)
}
