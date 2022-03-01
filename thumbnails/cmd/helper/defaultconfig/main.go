package main

import (
	"fmt"

	"github.com/owncloud/ocis/accounts/pkg/config"
	"github.com/owncloud/ocis/accounts/pkg/config/parser"
	"gopkg.in/yaml.v2"
)

func main() {

	cfg := config.DefaultConfig()

	parser.EnsureDefaults(cfg)
	parser.Sanitize(cfg)

	b, err := yaml.Marshal(cfg)
	if err != nil {
		return
	}
	fmt.Println(string(b))
}
