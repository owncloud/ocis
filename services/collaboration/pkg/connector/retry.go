package connector

import (
	"context"
	"time"

	rpcv1beta1 "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	"github.com/rs/zerolog"
)

const (
	// tooEarlyInitialWait is the waiting time before the first retry
	tooEarlyInitialWait = 20 * time.Millisecond
	// tooEarlyMaxWait caps the waiting time between retries
	tooEarlyMaxWait = 5 * time.Second
	// tooEarlyMaxRetries limits the retries, so the request can't hang forever
	tooEarlyMaxRetries = 10
)

// statusResponse is implemented by every CS3 response carrying a status
type statusResponse interface {
	GetStatus() *rpcv1beta1.Status
}

// retryWhileTooEarly calls the request until the response status isn't
// "TOO_EARLY", which is returned while the target resource is still being
// processed (e.g. right after an upload). The waiting time starts at 20ms and
// doubles on each retry, capped at 5 seconds.
//
// The last response is returned if the resource is still processing after all
// the retries, so the caller must handle the "TOO_EARLY" code too. Errors abort
// the retries immediately, and so does a cancelled context, in which case the
// context error is returned.
func retryWhileTooEarly[T statusResponse](ctx context.Context, logger zerolog.Logger, msg string, request func() (T, error)) (T, error) {
	for retries := 0; ; retries++ {
		res, err := request()
		if err != nil || res.GetStatus().GetCode() != rpcv1beta1.Code_CODE_TOO_EARLY || retries >= tooEarlyMaxRetries {
			return res, err
		}

		waitingTime := min(tooEarlyInitialWait<<retries, tooEarlyMaxWait)
		logger.Warn().
			Str("StatusCode", res.GetStatus().GetCode().String()).
			Str("StatusMsg", res.GetStatus().GetMessage()).
			Dur("WaitingTime", waitingTime).
			Int("Retries", retries).
			Msg(msg + ": file isn't ready yet. Retrying")

		select {
		case <-ctx.Done():
			return res, ctx.Err()
		case <-time.After(waitingTime):
			// nothing to do, the waiting time elapsed so we just retry the request
		}
	}
}
