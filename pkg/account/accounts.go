package account

import "github.com/owncloud/ocis-accounts/pkg/config"

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
	Read(key string) *Record
	// Write a record
	Write(*Record) *Record
	// List all records
	List() []*Record
}

// Record is an entry in the account storage
type Record struct {
	Key   string
	Value []byte
}
