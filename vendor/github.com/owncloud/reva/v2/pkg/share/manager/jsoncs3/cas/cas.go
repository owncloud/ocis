// Package cas classifies CAS-conflict errors shared by providercache, sharecache, and receivedsharecache.
package cas

import (
	"context"
	"time"

	backoff "github.com/cenkalti/backoff/v5"
	"github.com/owncloud/reva/v2/pkg/errtypes"
	"github.com/owncloud/reva/v2/pkg/storage/utils/metadata"
	"github.com/rs/zerolog"
	grpccodes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// NewBackoff returns the backoff policy shared by every CAS-retry loop.
func NewBackoff() *backoff.ExponentialBackOff {
	bo := backoff.NewExponentialBackOff()
	bo.InitialInterval = 500 * time.Microsecond
	bo.Multiplier = 2.0
	bo.RandomizationFactor = 1.0
	bo.MaxInterval = 50 * time.Millisecond
	return bo
}

// IsConflict reports whether err is a CAS conflict the caller should resync and retry on.
func IsConflict(err error) bool {
	switch err.(type) {
	case errtypes.Aborted:
		// If-Match etag check failed
		return true
	case errtypes.PreconditionFailed:
		// same as Aborted; some server paths return this instead
		return true
	case errtypes.AlreadyExists:
		// If-None-Match=* conflict: cache thought there was no file yet
		return true
	case errtypes.TooEarly:
		// an upload is already in progress for this path
		return true
	default:
		return false
	}
}

// IsSyncTransient reports whether a sync's Download error is worth retrying.
func IsSyncTransient(err error) bool {
	_, isTooEarly := err.(errtypes.IsTooEarly)
	return isTooEarly || IsTransientGRPCStatus(err)
}

// IsTransientGRPCStatus catches raw gRPC transport errors that metadata.CS3 never wraps in errtypes.
func IsTransientGRPCStatus(err error) bool {
	st, ok := status.FromError(err)
	if !ok {
		return false
	}
	switch st.Code() {
	case grpccodes.Unavailable, grpccodes.DeadlineExceeded, grpccodes.Canceled, grpccodes.ResourceExhausted:
		return true
	default:
		return false
	}
}

// RetryPersist retries persistFunc on a CAS conflict, calling resyncFunc between attempts to
// refresh in-memory state before the next persist, up to maxRetries times with this package's
// shared backoff. onConflict/onPersistFailed/onResyncFailed let the caller record its own
// span/log side effects for each outcome without this package depending on otel/zerolog.
func RetryPersist(ctx context.Context, maxRetries int, persistFunc, resyncFunc func() error,
	onConflict, onPersistFailed, onResyncFailed func(err error)) error {
	bo := NewBackoff()
	var err error
	for retries := maxRetries; retries > 0; retries-- {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		err = persistFunc()
		switch {
		case err == nil:
			return nil
		case IsConflict(err):
			onConflict(err)
		default:
			onPersistFailed(err)
			return err
		}
		if rerr := resyncFunc(); rerr != nil {
			onResyncFailed(rerr)
			return rerr
		}
		timer := time.NewTimer(bo.NextBackOff())
		select {
		case <-timer.C:
		case <-ctx.Done():
			timer.Stop()
			return ctx.Err()
		}
	}
	return err
}

// DecideNotFoundReset reports whether a NotFound download result should be treated as a
// legitimately empty resource (reset) rather than a lost update (error). A resync with a
// prior etag never resets silently; otherwise the backend's trash state decides, failing
// open (reset=true) when the backend can't tell (e.g. Disk). trashed is only meaningful
// when reset is false, to pick the right log reason.
func DecideNotFoundReset(ctx context.Context, storage metadata.Storage, path string, resetOnNotFound, hadPriorEtag bool, log zerolog.Logger) (reset, trashed bool) {
	if !resetOnNotFound && hadPriorEtag {
		return false, false
	}
	trashed, terr := storage.WasRecentlyDeleted(ctx, path)
	if terr != nil {
		log.Warn().Err(terr).Msg("could not check trash state, assuming not deleted")
		return true, false
	}
	return !trashed, trashed
}
