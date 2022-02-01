package config

import (
	"context"

	"github.com/owncloud/ocis/ocis-pkg/shared"
)

// Config combines all available configuration parts.
type Config struct {
	*shared.Commons `yaml:"-"`

	Service Service `yaml:"-"`

	Tracing *Tracing
	Log     *Log
	Debug   Debug

	HTTP HTTP
	GRPC GRPC

	TokenManager TokenManager

	Asset              Asset
	Repo               Repo
	Index              Index
	ServiceUser        ServiceUser
	HashDifficulty     int  `env:"ACCOUNTS_HASH_DIFFICULTY"`
	DemoUsersAndGroups bool `env:"ACCOUNTS_DEMO_USERS_AND_GROUPS"`

	Context context.Context `yaml:"-"`
}

// Asset defines the available asset configuration.
type Asset struct {
	Path string `env:"ACCOUNTS_ASSET_PATH"`
}

// TokenManager is the config for using the reva token manager
type TokenManager struct {
	JWTSecret string `env:"OCIS_JWT_SECRET;ACCOUNTS_JWT_SECRET"`
}

// Repo defines which storage implementation is to be used.
type Repo struct {
	Backend string `env:"ACCOUNTS_STORAGE_BACKEND"`
	Disk    Disk
	CS3     CS3
}

// Disk is the local disk implementation of the storage.
type Disk struct {
	Path string `env:"ACCOUNTS_STORAGE_DISK_PATH"`
}

// CS3 is the cs3 implementation of the storage.
type CS3 struct {
	ProviderAddr string `env:"ACCOUNTS_STORAGE_CS3_PROVIDER_ADDR"`
	JWTSecret    string `env:"ACCOUNTS_STORAGE_CS3_JWT_SECRET"`
}

// ServiceUser defines the user required for EOS.
type ServiceUser struct {
	UUID     string `env:"ACCOUNTS_SERVICE_USER_UUID"`
	Username string `env:"ACCOUNTS_SERVICE_USER_USERNAME"`
	UID      int64  `env:"ACCOUNTS_SERVICE_USER_UID"`
	GID      int64  `env:"ACCOUNTS_SERVICE_USER_GID"`
}

// Index defines config for indexes.
type Index struct {
	UID UIDBound
	GID GIDBound
}

// GIDBound defines a lower and upper bound.
type GIDBound struct {
	Lower int64 `env:"ACCOUNTS_GID_INDEX_LOWER_BOUND"`
	Upper int64 `env:"ACCOUNTS_GID_INDEX_UPPER_BOUND"`
}

// UIDBound defines a lower and upper bound.
type UIDBound struct {
	Lower int64 `env:"ACCOUNTS_UID_INDEX_LOWER_BOUND"`
	Upper int64 `env:"ACCOUNTS_UID_INDEX_UPPER_BOUND"`
}
