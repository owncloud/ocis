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

	GRPC GRPC `ocisConfig:"grpc"`

	Reva Reva `ocisConfig:"reva"`

	TokenManager TokenManager `ocisConfig:"token_manager"`

	ServiceMap ServiceMap `ocisConfig:"service_map"`

	StorageRegistry StorageRegistry `ocisConfig:"storage_registry"`

	Context context.Context
}

type StorageRegistry struct {
	Driver       string `ocisConfig:"driver"`
	HomeProvider string `ocisConfig:"home_provider"`
	Rules        Rules  `ocisConfig:"rules"`
	Storages     Storages
}

type Rules struct {
}

type Storages struct {
	StorageHome        StorageHome
	StorageUsers       StorageUsers
	StoragePublicShare StoragePublicShare
}
type StorageHome struct {
	MountPath     string
	AlternativeID string
}

type StorageUsers struct {
	MountPath string
	MountID   string
}

type StoragePublicShare struct {
	MountPath string
	MountID   string
}
