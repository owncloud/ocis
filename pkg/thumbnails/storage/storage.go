package storage

// Context combines different attributes needed for storage operations.
type Context struct {
	ETag   string
	Types  []string
	Width  int
	Height int
}

// Storage defines the interface for a thumbnail store.
type Storage interface {
	Get(string) []byte
	Set(string, []byte) error
	BuildKey(Context) string
}
