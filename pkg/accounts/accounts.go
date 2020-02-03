package accounts

// Account is an accounts service interface
type Account interface {
	// Read a record
	Read(key string) (*Record, error)
	// Write a record
	Write(Record) Record
	// List all records
	List() []*Record
}

// Record is an entry in the account storage
type Record struct {
	Key   string
	Value []byte
}
