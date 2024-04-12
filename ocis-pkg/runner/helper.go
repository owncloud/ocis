package runner

import (
	"fmt"
	"time"
)

// InterruptedTimeoutRunner will create a new runner (R2) based an original
// runner (R1).
// The new runner (R2) will monitor the original (R1). Once the `Interrupt`
// method is called in the new (R2), the interruption will be delivered to
// the original (R1), but a timeout will start. If we reach the timeout
// before the original runner (R1) is finished, the new runner (R2) will
// return an error.
//
// Any valid duration can be provided for the timeout, but you should give
// enough time for the task to finish in order to get the error from the
// original task (R1) and not the timeout one from the new (R2).
// Depending on the task, 5s, 10s or 30s might be reasonable timeout values.
//
// The timeout will start once the new (R2) runner is interrupted, either
// manually or via context
//
// Note that R2 can't stop R1 in any way. Even if R2 returns a "timeout" error,
// R1 might still be running and consuming resources.
// This method is intended to provide a way to ensure that the main thread
// won't be blocked forever.
func InterruptedTimeoutRunner(r *Runner, d time.Duration) *Runner {
	timeoutCh := make(chan time.Time)
	return New(r.ID, func() error {
		ch := make(chan *Result)
		r.RunAsync(ch)

		select {
		case result := <-ch:
			return result.RunnerError // forward the runner error
		case t := <-timeoutCh:
			// timeout reached. We can't stop the task, but we'll return
			// an error instead to prevent blocking the thread.
			return fmt.Errorf("Timeout reached at %s after waiting for %s after being interrupted", t.String(), d.String())
		}
	}, func() {
		go func() {
			select {
			case <-r.Finished():
				// Task finished -> runner should be delivering the result
			case t := <-time.After(d):
				// timeout reached -> send it through the channel so our runner
				// can abort
				timeoutCh <- t
			}
		}()
		r.Interrupt()
	})
}
