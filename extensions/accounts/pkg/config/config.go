package config

import (
	"context"

	"github.com/owncloud/ocis/ocis-pkg/shared"
)

// Config combines all available configuration parts.
type Config struct {
	*shared.Commons `yaml:"-"`

	Service Service `yaml:"-"`

	Tracing *Tracing `yaml:"tracing"`
	Log     *Log     `yaml:"log"`
	Debug   Debug    `yaml:"debug"`

	HTTP HTTP `yaml:"http"`
	GRPC GRPC `yaml:"grpc"`

	TokenManager TokenManager `yaml:"token_manager"`

	Asset              Asset       `yaml:"asset"`
	Repo               Repo        `yaml:"repo"`
	Index              Index       `yaml:"index"`
	ServiceUser        ServiceUser `yaml:"service_user"`
	HashDifficulty     int         `yaml:"hash_difficulty" env:"ACCOUNTS_HASH_DIFFICULTY" desc:"The hash difficulty makes sure that validating a password takes at least a certain amount of time."`
	DemoUsersAndGroups bool        `yaml:"demo_users_and_groups" env:"ACCOUNTS_DEMO_USERS_AND_GROUPS" desc:"If this flag is set the service will setup the demo users and groups."`

	ConfigFile           string `yaml:"-" env:"ACCOUNTS_CONFIG_FILE"` // config file to be used by the accounts extension
	ConfigFileHasBeenSet bool   `yaml:"-"`

	Context context.Context `yaml:"-"`
}

// Asset defines the available asset configuration.
type Asset struct {
	Path string `yaml:"path" env:"ACCOUNTS_ASSET_PATH" desc:"The path to the ui assets."`
}

// TokenManager is the config for using the reva token manager
type TokenManager struct {
	JWTSecret string `yaml:"jwt_secret" env:"OCIS_JWT_SECRET;ACCOUNTS_JWT_SECRET" desc:"The secret to mint jwt tokens."`
}

// Repo defines which storage implementation is to be used.
type Repo struct {
	Backend string `yaml:"backend" env:"ACCOUNTS_STORAGE_BACKEND" desc:"Defines which storage implementation is to be used"`
	Disk    Disk   `yaml:"disk"`
	CS3     CS3    `yaml:"cs3"`
}

// Disk is the local disk implementation of the storage.
type Disk struct {
	Path string `yaml:"path" env:"ACCOUNTS_STORAGE_DISK_PATH" desc:"The path where the accounts data is stored."`
}

// CS3 is the cs3 implementation of the storage.
type CS3 struct {
	ProviderAddr string `yaml:"provider_addr" env:"ACCOUNTS_STORAGE_CS3_PROVIDER_ADDR" desc:"The address to the storage provider."`
}

// ServiceUser defines the user required for EOS.
type ServiceUser struct {
	UUID     string `yaml:"uuid" env:"ACCOUNTS_SERVICE_USER_UUID" desc:"The id of the accounts service user."`
	Username string `yaml:"username" env:"ACCOUNTS_SERVICE_USER_USERNAME" desc:"The username of the accounts service user."`
	UID      int64  `yaml:"uid" env:"ACCOUNTS_SERVICE_USER_UID" desc:"The uid of the accounts service user."`
	GID      int64  `yaml:"gid" env:"ACCOUNTS_SERVICE_USER_GID" desc:"The gid of the accounts service user."`
}

// Index defines config for indexes.
type Index struct {
	UID UIDBound `yaml:"uid"`
	GID GIDBound `yaml:"gid"`
}

// GIDBound defines a lower and upper bound.
type GIDBound struct {
	Lower int64 `yaml:"lower" env:"ACCOUNTS_GID_INDEX_LOWER_BOUND" desc:"The lowest possible gid value for the indexer."`
	Upper int64 `yaml:"upper" env:"ACCOUNTS_GID_INDEX_UPPER_BOUND" desc:"The highest possible gid value for the indexer."`
}

// UIDBound defines a lower and upper bound.
type UIDBound struct {
	Lower int64 `yaml:"lower" env:"ACCOUNTS_UID_INDEX_LOWER_BOUND" desc:"The lowest possible uid value for the indexer."`
	Upper int64 `yaml:"upper" env:"ACCOUNTS_UID_INDEX_UPPER_BOUND" desc:"The highest possible uid value for the indexer."`
}
