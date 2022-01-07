package config

import (
	"context"

	"github.com/owncloud/ocis/ocis-pkg/shared"
)

// Config combines all available configuration parts.
type Config struct {
	*shared.Commons

	Service Service

	Tracing Tracing `ocisConfig:"tracing"`
	Log     *Log    `ocisConfig:"log"`
	Debug   Debug   `ocisConfig:"debug"`

	HTTP HTTP `ocisConfig:"http"`

	TokenManager TokenManager `ocisConfig:"token_manager"`
	Reva         Reva         `ocisConfig:"reva"`

	IdentityManagement IdentityManagement `ocisConfig:"identity_management"`

	AccountBackend     string `ocisConfig:"account_backend" env:"OCS_ACCOUNT_BACKEND_TYPE"`
	StorageUsersDriver string `ocisConfig:"storage_users_driver" env:"STORAGE_USERS_DRIVER;OCS_STORAGE_USERS_DRIVER"`
	MachineAuthAPIKey  string `ocisConfig:"machine_auth_api_key" env:"OCIS_MACHINE_AUTH_API_KEY;OCS_MACHINE_AUTH_API_KEY"`

	Context context.Context
}

// IdentityManagement keeps track of the OIDC address. This is because Reva requisite of uniqueness for users
// is based in the combination of IDP hostname + UserID. For more information see:
// https://github.com/cs3org/reva/blob/4fd0229f13fae5bc9684556a82dbbd0eced65ef9/pkg/storage/utils/decomposedfs/node/node.go#L856-L865
type IdentityManagement struct {
	Address string `ocisConfig:"address" env:"OCIS_URL;OCS_IDM_ADDRESS"`
}
