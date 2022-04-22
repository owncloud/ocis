package config

import (
	"context"

	"github.com/owncloud/ocis/ocis-pkg/shared"
)

// Config combines all available configuration parts.
type Config struct {
	*shared.Commons `yaml:"-"`

	Service Service `yaml:"-"`

	Tracing *Tracing `yaml:"tracing,omitempty"`
	Log     *Log     `yaml:"log,omitempty"`
	Debug   Debug    `yaml:"debug,omitempty"`

	HTTP HTTP `yaml:"http,omitempty"`
	GRPC GRPC `yaml:"grpc,omitempty"`

	TokenManager *shared.TokenManager `yaml:"token_manager,omitempty"`

	Asset              Asset       `yaml:"asset,omitempty"`
	Repo               Repo        `yaml:"repo,omitempty"`
	Index              Index       `yaml:"index,omitempty"`
	ServiceUser        ServiceUser `yaml:"service_user,omitempty"`
	HashDifficulty     int         `yaml:"hash_difficulty,omitempty" env:"ACCOUNTS_HASH_DIFFICULTY" desc:"The hash difficulty makes sure that validating a password takes at least a certain amount of time."`
	DemoUsersAndGroups bool        `yaml:"demo_users_and_groups,omitempty" env:"ACCOUNTS_DEMO_USERS_AND_GROUPS" desc:"If this flag is set the service will setup the demo users and groups."`

	Context context.Context `yaml:"-"`
}

// Asset defines the available asset configuration.
type Asset struct {
	Path string `yaml:"path" env:"ACCOUNTS_ASSET_PATH" desc:"The path to the ui assets."`
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
