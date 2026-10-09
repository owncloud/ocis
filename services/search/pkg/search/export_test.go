package search

import "time"

// SetExtractionRetryDelay shortens the delay between extraction retries for a
// test and returns a function that restores the previous value.
func SetExtractionRetryDelay(d time.Duration) func() {
	prev := _extractionRetryDelay
	_extractionRetryDelay = d
	return func() { _extractionRetryDelay = prev }
}
