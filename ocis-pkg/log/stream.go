package log

import (
	mdlog "go-micro.dev/v4/debug/log"
)

type logStream struct {
	stream <-chan mdlog.Record
	stop   chan bool
}

func (l *logStream) Chan() <-chan mdlog.Record {
	return l.stream
}

func (l *logStream) Stop() error {
	select {
	case <-l.stop:
		return nil
	default:
		close(l.stop)
	}
	return nil
}
