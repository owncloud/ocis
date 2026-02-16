package idm

import (
	_ "embed"
)

// BaseLDIF is a template for base LDAP entries
//
//go:embed ldif/base.ldif.tmpl
var BaseLDIF string

// DemoUsersLDIF is a template for demo users
//
//go:embed ldif/demousers.ldif.tmpl
var DemoUsersLDIF string
