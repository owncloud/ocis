package idm

import (
	_ "embed"
)

// BaseLDIF
// FIXME: nolint
// nolint: revive
//
//go:embed ldif/base.ldif.tmpl
var BaseLDIF string

// DemoUsersLDIF
// FIXME: nolint
// nolint: revive
//
//go:embed ldif/demousers.ldif
var DemoUsersLDIF string
