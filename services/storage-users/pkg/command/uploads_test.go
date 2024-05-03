package command

import (
	"testing"

	"github.com/cs3org/reva/v2/pkg/storage"
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
			filter:       storage.UploadSessionFilter{Processing: boolPtr(true)},
			expectedInfo: "Processing sessions:",
		},
		{
			alias:        "processing and not expired",
			filter:       storage.UploadSessionFilter{Processing: boolPtr(true), Expired: boolPtr(false)},
			expectedInfo: "Processing, not expired sessions:",
		},
		{
			alias:        "processing and expired",
			filter:       storage.UploadSessionFilter{Processing: boolPtr(true), Expired: boolPtr(true)},
			expectedInfo: "Processing, expired sessions:",
		},
		{
			alias:        "with id",
			filter:       storage.UploadSessionFilter{ID: strPtr("123")},
			expectedInfo: "Session with id '123':",
		},
		{
			alias:        "processing, not expired and not virus infected",
			filter:       storage.UploadSessionFilter{Processing: boolPtr(true), Expired: boolPtr(false), HasVirus: boolPtr(false)},
			expectedInfo: "Processing, not expired, not virusinfected sessions:",
		},
		{
			alias:        "not virusinfected",
			filter:       storage.UploadSessionFilter{HasVirus: boolPtr(false)},
			expectedInfo: "Not virusinfected sessions:",
		},
		{
			alias:        "expired and virusinfected",
			filter:       storage.UploadSessionFilter{Expired: boolPtr(true), HasVirus: boolPtr(true)},
			expectedInfo: "Expired, virusinfected sessions:",
		},
		{
			alias:        "expired and not virus infected",
			filter:       storage.UploadSessionFilter{Expired: boolPtr(true), HasVirus: boolPtr(false)},
			expectedInfo: "Expired, not virusinfected sessions:",
		},
		{
			alias:        "processing, not expired, virus infected and with id (note: this makes no sense)",
			filter:       storage.UploadSessionFilter{Processing: boolPtr(true), Expired: boolPtr(false), HasVirus: boolPtr(true), ID: strPtr("123")},
			expectedInfo: "Processing, not expired, virusinfected session with id '123':",
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

func boolPtr(b bool) *bool {
	return &b
}

func strPtr(s string) *string {
	return &s
}
