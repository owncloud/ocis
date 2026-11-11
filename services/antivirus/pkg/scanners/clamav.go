package scanners

import (
	"fmt"
	"time"

	"github.com/dutchcoders/go-clamd"
)

// defaultTimeout is used when a non-positive timeout is configured, so a
// misconfiguration can't silently turn every scan into an instant timeout.
const defaultTimeout = 5 * time.Minute

// NewClamAV returns a Scanner talking to clamAV via socket
func NewClamAV(socket string, timeout time.Duration) *ClamAV {
	if timeout <= 0 {
		timeout = defaultTimeout
	}
	return &ClamAV{
		clamd:   clamd.NewClamd(socket),
		timeout: timeout,
	}
}

// ClamAV is a Scanner based on clamav
type ClamAV struct {
	clamd   *clamd.Clamd
	timeout time.Duration
}

type clamavOutcome struct {
	result Result
	err    error
}

// Scan to fulfill Scanner interface
func (s ClamAV) Scan(in Input) (Result, error) {
	abort := make(chan bool)
	// Closing (never sending on) abort is what makes go-clamd's own watcher
	// goroutine tear down the underlying connection (see ScanStream) - on every
	// return path, not just the timeout one, or that goroutine blocks forever.
	defer close(abort)

	done := make(chan clamavOutcome, 1)
	go func() {
		ch, err := s.clamd.ScanStream(in.Body, abort)
		if err != nil {
			done <- clamavOutcome{err: err}
			return
		}

		r := <-ch
		if r == nil {
			// ch is closed without a value when the connection breaks before a
			// result line is read (e.g. clamd crashes or resets mid-scan).
			done <- clamavOutcome{err: fmt.Errorf("clamav: connection closed before a scan result was received")}
			return
		}

		done <- clamavOutcome{result: Result{
			Infected:    r.Status == clamd.RES_FOUND,
			Description: r.Description,
			ScanTime:    time.Now(),
		}}
	}()

	select {
	case o := <-done:
		return o.result, o.err
	case <-time.After(s.timeout):
		// Runs the whole ScanStream call (dial, writes, EOF, awaiting the
		// result) above, not just the final receive, so a backend that stops
		// draining the socket is bounded too, not only one that stops replying.
		return Result{}, fmt.Errorf("clamav scan timed out after %s", s.timeout)
	}
}
