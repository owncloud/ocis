package idm

import (
	_ "embed"
)

//go:embed ldif/base.ldif.tmpl
var BaseLDIF string

//go:embed ldif/demousers.ldif
var DemoUsersLDIF string
