package config

import (
	"context"

	"github.com/owncloud/ocis/ocis-pkg/shared"
)

// Config combines all available configuration parts.
type Config struct {
	*shared.Commons

	Service Service

	Tracing *Tracing `ocisConfig:"tracing"`
	Log     *Log     `ocisConfig:"log"`
	Debug   Debug    `ocisConfig:"debug"`

	IDM Settings `ocisConfig:"idm"`

	Context context.Context
}

type Settings struct {
	LDAPSAddr     string `ocisConfig:"ldaps_addr" env:"IDM_LDAPS_ADDR"`
	Cert          string `ocisConfig:"cert" env:"IDM_LDAPS_CERT"`
	Key           string `ocisConfig:"cert" env:"IDM_LDAPS_KEY"`
	DatabasePath  string `ocisConfig:"database" env:"IDM_DATABASE_PATH"`
	AdminPassword string `ocisConfig:"admin_password" env:"IDM_ADMIN_PASSWORD"`
}
