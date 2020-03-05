package storage

type StorageContext struct {
	ETag   string
	Types  []string
	Width  int
	Height int
}

// Storage defines the interface for a thumbnail store.
type Storage interface {
	Get(string) []byte
	Set(string, []byte) error
	BuildKey(StorageContext) string
}
