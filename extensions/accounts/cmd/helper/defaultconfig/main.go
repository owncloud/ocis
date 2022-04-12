package main

import (
	"fmt"

	accountsdefaults "github.com/owncloud/ocis/extensions/accounts/pkg/config/defaults"
	idpdefaults "github.com/owncloud/ocis/extensions/idp/pkg/config/defaults"
	"gopkg.in/yaml.v2"
)

func main() {

	fn1 := accountsdefaults.FullDefaultConfig
	fn2 := idpdefaults.FullDefaultConfig

	b, err := yaml.Marshal(fn1())
	if err != nil {
		return
	}
	fmt.Println(string(b))

	b, err = yaml.Marshal(fn2())
	if err != nil {
		return
	}
	fmt.Println(string(b))
}
