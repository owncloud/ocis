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
	HashDifficulty     int         `ocisConfig:"hash_difficulty" env:"ACCOUNTS_HASH_DIFFICULTY"`
	DemoUsersAndGroups bool        `ocisConfig:"demo_users_and_groups" env:"ACCOUNTS_DEMO_USERS_AND_GROUPS"`

	Context context.Context
}

// Asset defines the available asset configuration.
type Asset struct {
	Path string `ocisConfig:"path" env:"ACCOUNTS_ASSET_PATH"`
}

// TokenManager is the config for using the reva token manager
type TokenManager struct {
	JWTSecret string `ocisConfig:"jwt_secret" env:"OCIS_JWT_SECRET;ACCOUNTS_JWT_SECRET"`
}

// Repo defines which storage implementation is to be used.
type Repo struct {
	Backend string `ocisConfig:"backend"  env:"ACCOUNTS_STORAGE_BACKEND"`
	Disk    Disk   `ocisConfig:"disk"`
	CS3     CS3    `ocisConfig:"cs3"`
}

// Disk is the local disk implementation of the storage.
type Disk struct {
	Path string `ocisConfig:"path" env:"ACCOUNTS_STORAGE_DISK_PATH"`
}

// CS3 is the cs3 implementation of the storage.
type CS3 struct {
	ProviderAddr string `ocisConfig:"provider_addr" env:"ACCOUNTS_STORAGE_CS3_PROVIDER_ADDR"`
	JWTSecret    string `ocisConfig:"jwt_secret" env:"ACCOUNTS_STORAGE_CS3_JWT_SECRET"`
}

// ServiceUser defines the user required for EOS.
type ServiceUser struct {
	UUID     string `ocisConfig:"uuid" env:"ACCOUNTS_SERVICE_USER_UUID"`
	Username string `ocisConfig:"username" env:"ACCOUNTS_SERVICE_USER_USERNAME"`
	UID      int64  `ocisConfig:"uid" env:"ACCOUNTS_SERVICE_USER_UID"`
	GID      int64  `ocisConfig:"gid" env:"ACCOUNTS_SERVICE_USER_GID"`
}

// Index defines config for indexes.
type Index struct {
	UID UIDBound `ocisConfig:"uid"`
	GID GIDBound `ocisConfig:"gid"`
}

// GIDBound defines a lower and upper bound.
type GIDBound struct {
	Lower int64 `ocisConfig:"lower" env:"ACCOUNTS_GID_INDEX_LOWER_BOUND"`
	Upper int64 `ocisConfig:"upper" env:"ACCOUNTS_GID_INDEX_UPPER_BOUND"`
}

// UIDBound defines a lower and upper bound.
type UIDBound struct {
	Lower int64 `ocisConfig:"lower" env:"ACCOUNTS_UID_INDEX_LOWER_BOUND"`
	Upper int64 `ocisConfig:"upper" env:"ACCOUNTS_UID_INDEX_UPPER_BOUND"`
}
