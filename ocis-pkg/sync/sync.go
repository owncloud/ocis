package sync

import "sync"

var (
	// ParsingViperConfig addresses the fact that config parsing using Viper is not thread safe.
	ParsingViperConfig sync.Mutex
)
