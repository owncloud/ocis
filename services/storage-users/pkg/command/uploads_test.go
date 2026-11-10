package command

import (
	"testing"

	"github.com/owncloud/reva/v2/pkg/storage"
	"github.com/test-go/testify/require"
)

func TestBuildInfo(t *testing.T) {
	testCases := []struct {
		alias        string
		filter       storage.UploadSessionFilter
		expectedInfo string
	}{
		{
			alias:        "empty filter",
			filter:       storage.UploadSessionFilter{},
			expectedInfo: "Sessions:",
		},
		{
			alias:        "processing",
			filter:       storage.UploadSessionFilter{Processing: new(true)},
			expectedInfo: "Processing sessions:",
		},
		{
			alias:        "processing and not expired",
			filter:       storage.UploadSessionFilter{Processing: new(true), Expired: new(false)},
			expectedInfo: "Processing, not expired sessions:",
		},
		{
			alias:        "processing and expired",
			filter:       storage.UploadSessionFilter{Processing: new(true), Expired: new(true)},
			expectedInfo: "Processing, expired sessions:",
		},
		{
			alias:        "with id",
			filter:       storage.UploadSessionFilter{ID: new("123")},
			expectedInfo: "Session with id '123':",
		},
		{
			alias:        "processing, not expired and not virus infected",
			filter:       storage.UploadSessionFilter{Processing: new(true), Expired: new(false), HasVirus: new(false)},
			expectedInfo: "Processing, not expired, not virusinfected sessions:",
		},
		{
			alias:        "not virusinfected",
			filter:       storage.UploadSessionFilter{HasVirus: new(false)},
			expectedInfo: "Not virusinfected sessions:",
		},
		{
			alias:        "expired and virusinfected",
			filter:       storage.UploadSessionFilter{Expired: new(true), HasVirus: new(true)},
			expectedInfo: "Expired, virusinfected sessions:",
		},
		{
			alias:        "expired and not virus infected",
			filter:       storage.UploadSessionFilter{Expired: new(true), HasVirus: new(false)},
			expectedInfo: "Expired, not virusinfected sessions:",
		},
		{
			alias:        "processing, not expired, virus infected and with id (note: this makes no sense)",
			filter:       storage.UploadSessionFilter{Processing: new(true), Expired: new(false), HasVirus: new(true), ID: new("123")},
			expectedInfo: "Processing, not expired, virusinfected session with id '123':",
		},
		{
			alias:        "orphaned",
			filter:       storage.UploadSessionFilter{Orphaned: new(true)},
			expectedInfo: "Orphaned sessions:",
		},
		{
			alias:        "not orphaned",
			filter:       storage.UploadSessionFilter{Orphaned: new(false)},
			expectedInfo: "Not orphaned sessions:",
		},
		{
			alias:        "processing and orphaned",
			filter:       storage.UploadSessionFilter{Processing: new(true), Orphaned: new(true)},
			expectedInfo: "Processing, orphaned sessions:",
		},
		{
			alias:        "expired, not virus infected and orphaned",
			filter:       storage.UploadSessionFilter{Expired: new(true), HasVirus: new(false), Orphaned: new(true)},
			expectedInfo: "Expired, not virusinfected, orphaned sessions:",
		},
	}

	for _, tc := range testCases {
		alias := tc.alias
		filter := tc.filter
		expectedInfo := tc.expectedInfo

		t.Run(alias, func(t *testing.T) {
			require.Equal(t, expectedInfo, buildInfo(filter))
		})
	}
}

//go:fix inline
func boolPtr(b bool) *bool {
	return new(b)
}

//go:fix inline
func strPtr(s string) *string {
	return new(s)
}
