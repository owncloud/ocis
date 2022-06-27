package timing

import (
	"context"
	"time"
)

// key is an unexported type for keys defined in this package.
type key int

// elapsedKey is the key for elapsedRecord which holds start end elapsed time
// in Contexts.
const elapsedKey key = 0

type elapsedRecord struct {
	start   time.Time
	elapsed time.Duration
	done    chan bool
	cancel  context.CancelFunc
}

// NewContext returns a new Context that carries the start time and calls the
// provided stopped function then the created Context is cancelled. The stopped
// function can be nil in which case nothing is called.
func NewContext(parent context.Context, stopped func(elapsed time.Duration)) context.Context {
	ctx, cancel := context.WithCancel(parent)
	recordPtr := &elapsedRecord{
		start:  time.Now(),
		done:   make(chan bool),
		cancel: cancel,
	}
	ctx = context.WithValue(ctx, elapsedKey, recordPtr)
	go func() {
		<-ctx.Done()
		recordPtr.elapsed = time.Since(recordPtr.start)
		close(recordPtr.done)
		if stopped != nil {
			stopped(recordPtr.elapsed)
		}
	}()

	return ctx
}

// StartFromContext returns the start time from the provided Context.
func StartFromContext(ctx context.Context) time.Time {
	return ctx.Value(elapsedKey).(*elapsedRecord).start
}

// ElapsedFromContext returns the elapsed time from the provided Context. If the
// provided context is not yet cancelled the duration from now since start is
// returned.
func ElapsedFromContext(ctx context.Context) time.Duration {
	elapsed := ctx.Value(elapsedKey).(*elapsedRecord).elapsed
	if elapsed > 0 {
		return elapsed
	}

	return time.Since(StartFromContext(ctx))
}

// CancelContext cancels the provided Context if it carries start time.
func CancelContext(ctx context.Context) {
	recordPtr := ctx.Value(elapsedKey).(*elapsedRecord)
	if recordPtr != nil {
		recordPtr.cancel()
		<-recordPtr.done // Wait until done is complete.
	}
}
