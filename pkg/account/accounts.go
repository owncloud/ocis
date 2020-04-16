package account

import (
	"github.com/owncloud/ocis-accounts/pkg/config"
	"github.com/owncloud/ocis-accounts/pkg/proto/v0"
)

var (
	// Registry uses the strategy pattern as a registry
	Registry = map[string]RegisterFunc{}

	// DefaultManager defines the default accounts manager
	DefaultManager = "filesystem"
)

// RegisterFunc stores store constructors
type RegisterFunc func(*config.Config) Manager

// Manager is an accounts service interface
type Manager interface {
	// Read a record
	Read(key string) (*proto.Record, error)
	// Write a record
	Write(*proto.Record) (*proto.Record, error)
	// List all records
	List() ([]*proto.Record, error)
}
