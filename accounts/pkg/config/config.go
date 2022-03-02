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

	HTTP HTTP `ocisConfig:"http"`
	GRPC GRPC `ocisConfig:"grpc"`

	TokenManager TokenManager `ocisConfig:"token_manager"`

	Asset              Asset       `ocisConfig:"asset"`
	Repo               Repo        `ocisConfig:"repo"`
	Index              Index       `ocisConfig:"index"`
	ServiceUser        ServiceUser `ocisConfig:"service_user"`
	HashDifficulty     int         `ocisConfig:"hash_difficulty" env:"ACCOUNTS_HASH_DIFFICULTY" desc:"The hash difficulty makes sure that validating a password takes at least a certain amount of time."`
	DemoUsersAndGroups bool        `ocisConfig:"demo_users_and_groups" env:"ACCOUNTS_DEMO_USERS_AND_GROUPS" desc:"If this flag is set the service will setup the demo users and groups."`

	Context context.Context
}

// Asset defines the available asset configuration.
type Asset struct {
	Path string `ocisConfig:"path" env:"ACCOUNTS_ASSET_PATH" desc:"The path to the ui assets."`
}

// TokenManager is the config for using the reva token manager
type TokenManager struct {
	JWTSecret string `ocisConfig:"jwt_secret" env:"OCIS_JWT_SECRET;ACCOUNTS_JWT_SECRET" desc:"The secret to mint jwt tokens."`
}

// Repo defines which storage implementation is to be used.
type Repo struct {
	Backend string `ocisConfig:"backend" env:"ACCOUNTS_STORAGE_BACKEND" desc:"Defines which storage implementation is to be used"`
	Disk    Disk
	CS3     CS3
}

// Disk is the local disk implementation of the storage.
type Disk struct {
	Path string `ocisConfig:"path" env:"ACCOUNTS_STORAGE_DISK_PATH" desc:"The path where the accounts data is stored."`
}

// CS3 is the cs3 implementation of the storage.
type CS3 struct {
	ProviderAddr string `ocisConfig:"provider_addr" env:"ACCOUNTS_STORAGE_CS3_PROVIDER_ADDR" desc:"The address to the storage provider."`
}

// ServiceUser defines the user required for EOS.
type ServiceUser struct {
	UUID     string `ocisConfig:"uuid" env:"ACCOUNTS_SERVICE_USER_UUID" desc:"The id of the accounts service user."`
	Username string `ocisConfig:"username" env:"ACCOUNTS_SERVICE_USER_USERNAME" desc:"The username of the accounts service user."`
	UID      int64  `ocisConfig:"uid" env:"ACCOUNTS_SERVICE_USER_UID" desc:"The uid of the accounts service user."`
	GID      int64  `ocisConfig:"gid" env:"ACCOUNTS_SERVICE_USER_GID" desc:"The gid of the accounts service user."`
}

// Index defines config for indexes.
type Index struct {
	UID UIDBound `ocisConfig:"uid"`
	GID GIDBound `ocisConfig:"gid"`
}

// GIDBound defines a lower and upper bound.
type GIDBound struct {
	Lower int64 `ocisConfig:"lower" env:"ACCOUNTS_GID_INDEX_LOWER_BOUND" desc:"The lowest possible gid value for the indexer."`
	Upper int64 `ocisConfig:"upper" env:"ACCOUNTS_GID_INDEX_UPPER_BOUND" desc:"The highest possible gid value for the indexer."`
}

// UIDBound defines a lower and upper bound.
type UIDBound struct {
	Lower int64 `ocisConfig:"lower" env:"ACCOUNTS_UID_INDEX_LOWER_BOUND" desc:"The lowest possible uid value for the indexer."`
	Upper int64 `ocisConfig:"upper" env:"ACCOUNTS_UID_INDEX_UPPER_BOUND" desc:"The highest possible uid value for the indexer."`
}
