package search

import "time"

// SetExtractionRetryDelay shortens the delay between extraction retries for a
// test.
func (s *Service) SetExtractionRetryDelay(d time.Duration) {
	s.extractionRetryDelay = d
}
