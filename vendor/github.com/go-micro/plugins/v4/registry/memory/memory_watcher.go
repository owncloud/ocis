package memory

import (
	"errors"

	"go-micro.dev/v4/registry"
)

type memoryWatcher struct {
	exit chan bool
	opts registry.WatchOptions
}

func (m *memoryWatcher) Next() (*registry.Result, error) {
	// not implement so we just block until exit
	<-m.exit
	return nil, errors.New("watcher stopped")
}

func (m *memoryWatcher) Stop() {
	select {
	case <-m.exit:
		return
	default:
		close(m.exit)
	}
}
