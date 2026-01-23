package publicshares

import (
	"context"
	"encoding/json"
	"errors"
	"math/rand"
	"sync"
	"time"

	"github.com/owncloud/reva/v2/pkg/appctx"
	"github.com/owncloud/reva/v2/pkg/rgrpc"
	"github.com/rs/zerolog"
	microstore "go-micro.dev/v4/store"
	"google.golang.org/grpc/metadata"
)

const (
	maxWriteRetries = 10
)

// attemptData contains the data we need to store for each failed attempt.
// Right now, only the timestamp of the failed attempt is needed
type attemptData struct {
	ID        int   `json:"attemptId"`
	Timestamp int64 `json:"timestamp"`
}

// BruteForceProtection implements a rate-limit-based protection for the
// public shares.
// Given a time duration (10 minutes, for example), a maximum of X failed
// attempts are allowed. If that rate is reached, access to the public link
// should be blocked until the rate decreases.
// Note that the time the link should be blocked is undefined and will be
// somewhere between 0 and the given duration
type BruteForceProtection struct {
	ID          int
	mutex       *sync.Mutex
	timeGap     time.Duration
	maxAttempts int
	store       microstore.Store
}

// NewBruteForceProtection creates a new instance of BruteForceProtection
// If either the timeGap or maxAttempts are 0, the BruteForceProtection
// won't register any failed attempt and it will act as if it is disabled.
func NewBruteForceProtection(store microstore.Store, timeGap time.Duration, maxAttempts int) *BruteForceProtection {
	return &BruteForceProtection{
		ID:          rand.Int(),
		mutex:       &sync.Mutex{},
		timeGap:     timeGap,
		maxAttempts: maxAttempts,
		store:       store,
	}
}

// AddAttemptAndCheckAllow register a new failed attempt for the provided public share
// If the time gap or the max attempts are 0, the failed attempt won't be
// registered.
// The function returns a boolean whether there are more failed attempts
// available (allowing access), and an error if something wrong happens.
func (bfp *BruteForceProtection) AddAttemptAndCheckAllow(ctx context.Context, shareToken string) (bool, error) {
	if bfp.timeGap <= 0 || bfp.maxAttempts <= 0 {
		return true, nil
	}

	bfp.mutex.Lock()
	defer bfp.mutex.Unlock()

	attempt := &attemptData{
		ID:        rand.Int(),
		Timestamp: time.Now().Unix(),
	}

	log := appctx.GetLogger(ctx)
	sublog := log.With().
		Int("instanceID", bfp.ID).
		Str("shareToken", shareToken).
		Int("attemptID", attempt.ID).
		Int64("attemptTimestamp", attempt.Timestamp).
		Logger()

	attemptCount, err := bfp.writeWithRetry(shareToken, attempt, sublog)
	if err != nil {
		sublog.Error().Err(err).Msg("Could not include the failed attempt for brute force protection")
		return false, err
	}

	sublog.Debug().
		Int("attemptCount", attemptCount).
		Bool("stillAccessible", attemptCount <= bfp.maxAttempts).
		Msg("Failed attempt registered for brute force protection")

	return attemptCount <= bfp.maxAttempts, nil
}

// Verify will check if access to the share is available.
// It will also update the stored information (expiring old data).
// In case of errors, the verification will return false.
func (bfp *BruteForceProtection) Verify(ctx context.Context, shareToken string) bool {
	bfp.mutex.Lock()
	defer bfp.mutex.Unlock()

	log := appctx.GetLogger(ctx)
	sublog := log.With().
		Int("instanceID", bfp.ID).
		Str("shareToken", shareToken).
		Logger()

	attemptList, err := bfp.readFromStore(shareToken)
	if err != nil {
		sublog.Error().Err(err).Msg("Could not read from the cache store")
		return false
	}

	updatedList := bfp.cleanAttempts(attemptList)
	if len(attemptList) != len(updatedList) {
		sublog.Debug().Msg("Cleaning obsolete attempts")
		if _, err := bfp.writeWithRetry(shareToken, nil, sublog); err != nil {
			sublog.Error().Err(err).Msg("Failed to update the cache store")
			return false
		}
	}

	updatedListCount := len(updatedList)
	sublog.Debug().
		Int("attemptCount", updatedListCount).
		Bool("stillAccessible", updatedListCount <= bfp.maxAttempts).
		Msg("Verification for brute force protection done")

	return updatedListCount <= bfp.maxAttempts
}

// Ensure the attempt is added in the right possition. Since multiple attempts
// can happen very closely, the attempt at time 73 might have been registered
// before the attempt at time 72
func (bfp *BruteForceProtection) insertAttempt(attemptList []*attemptData, attempt *attemptData) []*attemptData {
	if attempt == nil {
		return attemptList
	}

	var i int
	for i = 0; i < len(attemptList); i++ {
		if attemptList[i].Timestamp > attempt.Timestamp {
			break
		}
	}

	return append(attemptList[:i], append([]*attemptData{attempt}, attemptList[i:]...)...)
}

// cleanAttempts will remove obsolete attempt data
func (bfp *BruteForceProtection) cleanAttempts(attemptList []*attemptData) []*attemptData {
	minTimestamp := time.Now().Add(-1 * bfp.timeGap).Unix()

	var index int
	for index = 0; index < len(attemptList); index++ {
		if attemptList[index].Timestamp >= minTimestamp {
			break
		}
	}

	return attemptList[index:]
}

// areEqualLists checks that both lists have the same data
func (bfp *BruteForceProtection) areEqualLists(a, b []*attemptData) bool {
	if len(a) != len(b) {
		return false
	}

	for i := 0; i < len(a); i++ { // already checked that the length of both lists are the same
		if a[i].Timestamp != b[i].Timestamp {
			return false
		}
	}
	return true
}

// writeWithRetry will write the attempt data for the share. You can use a
// nil attempt in order to just update and remove obsolete attempts.
// This involves several steps:
// 1. Read the current stored data.
// 2. Update and insert the attempt in its right position (although unlikely,
// replicas might have written a different attempt which happens later than
// the attempt we're writing). Cleanup obsolete data also happens here.
// 3. Write the data into the store.
// 4. Re-read the written data to check it wasn't overwritten. This implies
// checking and comparing the data written and read. If the check is wrong,
// go back to step 2.
//
// This method includes retries (up to maxWriteRetries = 10) in case the check
// fails.
//
// This method will return the number of the elements read for the check (at
// step 4) and an error if any.
func (bfp *BruteForceProtection) writeWithRetry(shareToken string, attempt *attemptData, logger zerolog.Logger) (int, error) {
	// read the stored data
	attemptList, err := bfp.readFromStore(shareToken)
	if err != nil {
		logger.Error().Err(err).Msg("Could not read from store before writing on it")
		return 0, err
	}

	tries := 0
	for tries = 0; tries < maxWriteRetries; tries++ {
		updatedList := bfp.cleanAttempts(bfp.insertAttempt(attemptList, attempt))
		if len(updatedList) == 0 {
			// if all attempts expired, delete the info
			// the new attempt shouldn't have expired, so this shouldn't happen
			if derr := bfp.deleteFromStore(shareToken); derr != nil {
				logger.Error().Err(err).Int("tryNo", tries).Msg("Error deleting from the cache store")
				return 0, derr
			}
		} else {

			// write the updated data into the store
			if werr := bfp.writeToStore(shareToken, updatedList); werr != nil {
				logger.Error().Err(err).Int("tryNo", tries).Msg("Error writing to the cache store")
				return 0, werr
			}
		}

		// re-read the data to ensure the data has been correctly updated
		attemptList, err = bfp.readFromStore(shareToken)
		if err != nil {
			logger.Error().Err(err).Int("tryNo", tries).Msg("Error reading from the cache store")
			return 0, err
		}

		// TODO: the areEqualLists might be too strict. We might just need
		// to check that the data we want to add is present despite there
		// could be additional data from other replicas
		if bfp.areEqualLists(attemptList, updatedList) {
			// if both lists are equal, the write was successful
			break
		}

		logger.Info().Int("tryNo", tries).Msg("Attempt data seems to have been overwritten. Retrying")
	}

	if tries >= maxWriteRetries {
		// couldn't write the data
		return 0, errors.New("Could not ensure the data to be written in the store. Retries spent")
	}
	return len(attemptList), nil
}

// readFromStore will read the data from store and return a list with the failed
// attempts. If there no failed attempts registered, an empty list will be returned.
func (bfp *BruteForceProtection) readFromStore(shareToken string) ([]*attemptData, error) {
	records, err := bfp.store.Read(shareToken)
	if errors.Is(err, microstore.ErrNotFound) {
		// if the key isn't found, use an empty list
		return make([]*attemptData, 0), nil
	} else if err != nil {
		return nil, err
	}

	attemptList := make([]*attemptData, 0)
	if jerr := json.Unmarshal(records[0].Value, &attemptList); jerr != nil {
		return nil, jerr
	}

	return attemptList, nil
}

// writeToStore will write the attempt list into the store
func (bfp *BruteForceProtection) writeToStore(shareToken string, attemptList []*attemptData) error {
	marshalledList, merr := json.Marshal(attemptList)
	if merr != nil {
		return merr
	}

	newRecord := &microstore.Record{
		Key:   shareToken,
		Value: marshalledList,
	}
	if werr := bfp.store.Write(newRecord); werr != nil {
		return werr
	}
	return nil
}

// deleteFromStore will delete the attempt list from the store. This should
// only be called if the attempt list is empty.
func (bfp *BruteForceProtection) deleteFromStore(shareToken string) error {
	return bfp.store.Delete(shareToken)
}

// outgoingContextKey is the key that will be used to mark that the brute
// force protection shouldn't consider the password for the share token
// (context value) as a failed attempt even if the password is wrong.
// The key is intended to auto-propagate across all the GRPC services.
const outgoingContextKey = rgrpc.AutoPropPrefix + "bfp-skip"

// MarkSkipAttemptContext will mark the share token so it will be skipped
// for the brute force protection. This means that the password for the
// share token won't be counted as a failed attempt even if the password
// is wrong.
// This "skip" will be valid within the returned context.
// The context key used should auto-propagate across all the GRPC services,
// assuming the metadata interceptors are in place (check
// internal/grpc/interceptors/metadata/metadata.go)
func MarkSkipAttemptContext(ctx context.Context, shareToken string) context.Context {
	return metadata.AppendToOutgoingContext(ctx, outgoingContextKey, shareToken)
}

// CheckSkipAttempt will check whether we should skip the brute force
// protection for the share token based on the context data.
// If you want to skip the protection, the MarkSkipAttemptContext method
// should have been called for the provided share token, and the returned
// context needs to be used.
// This method will return true if the context contains data marking the
// share token as "to skip" (from the MarkSkipAttemptContext method). If there
// is no such data, it will return false.
func CheckSkipAttempt(ctx context.Context, shareToken string) bool {
	possibleValues := metadata.ValueFromIncomingContext(ctx, outgoingContextKey)
	if possibleValues == nil || len(possibleValues) < 1 {
		return false
	}

	for _, value := range possibleValues {
		if value == shareToken {
			return true
		}
	}
	return false
}
