package log

import (
	mdlog "go-micro.dev/v4/debug/log"
)

// FIXME: nolint
// nolint: unused
type logStream struct {
	stream <-chan mdlog.Record
	stop   chan bool
}

// Chan
// FIXME: nolint
// nolint: unused
func (l *logStream) Chan() <-chan mdlog.Record {
	return l.stream
}

// Stop
// FIXME: nolint
// nolint: unused
func (l *logStream) Stop() error {
	select {
	case <-l.stop:
		return nil
	default:
		close(l.stop)
	}
	return nil
}
