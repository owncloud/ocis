package search

import "time"

// SetExtractionRetryDelay shortens the delay between extraction retries for a
// test.
func (s *Service) SetExtractionRetryDelay(d time.Duration) {
	s.extractionRetryDelay = d
}

// SetNow replaces the clock the service uses to date failed extractions.
func (s *Service) SetNow(now func() time.Time) {
	s.now = now
}

// AddPhotoMetadata exposes addPhotoMetadata to the tests.
var AddPhotoMetadata = addPhotoMetadata
