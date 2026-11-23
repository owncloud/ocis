package command

import (
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/owncloud/reva/v2/pkg/errtypes"
	"github.com/owncloud/reva/v2/pkg/storage/utils/decomposedfs/node"
)

// fakeBlobstore is a tree.Blobstore that returns a fixed result from Download.
type fakeBlobstore struct {
	rc  io.ReadCloser
	err error
}

func (f fakeBlobstore) Upload(*node.Node, string) error                     { return nil }
func (f fakeBlobstore) UploadFromReader(*node.Node, io.Reader, int64) error { return nil }
func (f fakeBlobstore) Download(*node.Node) (io.ReadCloser, error)          { return f.rc, f.err }
func (f fakeBlobstore) Delete(*node.Node) error                             { return nil }

func Test_downloadBlob(t *testing.T) {
	tests := []struct {
		name         string
		bs           fakeBlobstore
		wantErr      bool
		wantMessage  string
		wantNotFound bool
	}{
		{
			name: "blob exists",
			bs:   fakeBlobstore{rc: io.NopCloser(strings.NewReader("hello world"))},
		},
		{
			name:         "blob does not exist",
			bs:           fakeBlobstore{err: errtypes.NotFound("/blobs/00/00-0000")},
			wantErr:      true,
			wantMessage:  "download failed: could not read blob: error: not found: /blobs/00/00-0000",
			wantNotFound: true,
		},
		{
			name:        "other download error",
			bs:          fakeBlobstore{err: errors.New("permission denied")},
			wantErr:     true,
			wantMessage: "download failed: permission denied",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := downloadBlob(tt.bs, "blob-id", "space-id", 11)
			if (err != nil) != tt.wantErr {
				t.Fatalf("downloadBlob() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				return
			}
			if err.Error() != tt.wantMessage {
				t.Errorf("downloadBlob() message = %q, want %q", err.Error(), tt.wantMessage)
			}
			var nf errtypes.IsNotFound
			if errors.As(err, &nf) != tt.wantNotFound {
				t.Errorf("downloadBlob() NotFound = %v, want %v", !tt.wantNotFound, tt.wantNotFound)
			}
		})
	}
}

func Test_depathify(t *testing.T) {
	tests := []struct {
		name  string
		path  string
		depth int
		want  string
	}{
		{
			name:  "blob id depth 4",
			path:  "61/03/ab/c3/-b08a-4556-9937-2bf3065c1202",
			depth: 4,
			want:  "6103abc3-b08a-4556-9937-2bf3065c1202",
		},
		{
			name:  "space id depth 1",
			path:  "b1/9ec764-5398-458a-8ff1-1925bd906999",
			depth: 1,
			want:  "b19ec764-5398-458a-8ff1-1925bd906999",
		},
		{
			name:  "depth 0 is a no-op",
			path:  "abcd-1234",
			depth: 0,
			want:  "abcd-1234",
		},
		{
			name:  "fewer segments than depth leaves remainder intact",
			path:  "ab/cd",
			depth: 4,
			want:  "abcd",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := depathify(tt.path, tt.depth)
			if got != tt.want {
				t.Errorf("depathify(%q, %d) = %q, want %q", tt.path, tt.depth, got, tt.want)
			}
		})
	}
}

func Test_parseBlobPath(t *testing.T) {
	tests := []struct {
		name        string
		path        string
		wantSpaceID string
		wantBlobID  string
		wantErr     bool
	}{
		{
			// s3ng: <spaceID>/<pathified_blobID>
			name:        "s3ng format",
			path:        "b19ec764-5398-458a-8ff1-1925bd906999/61/03/ab/c3/-b08a-4556-9937-2bf3065c1202",
			wantSpaceID: "b19ec764-5398-458a-8ff1-1925bd906999",
			wantBlobID:  "6103abc3-b08a-4556-9937-2bf3065c1202",
		},
		{
			// ocis: …/spaces/<pathified_spaceID>/blobs/<pathified_blobID>
			name:        "ocis filesystem format",
			path:        "/var/lib/ocis/storage/users/spaces/b1/9ec764-5398-458a-8ff1-1925bd906999/blobs/61/03/ab/c3/-b08a-4556-9937-2bf3065c1202",
			wantSpaceID: "b19ec764-5398-458a-8ff1-1925bd906999",
			wantBlobID:  "6103abc3-b08a-4556-9937-2bf3065c1202",
		},
		{
			name:    "ocis format missing /blobs/ segment",
			path:    "/var/lib/ocis/storage/users/spaces/b1/9ec764-5398-458a-8ff1-1925bd906999/noblobs/61/03/ab/c3",
			wantErr: true,
		},
		{
			name:    "s3ng format missing blob id (no slash after spaceID)",
			path:    "b19ec764-5398-458a-8ff1-1925bd906999",
			wantErr: true,
		},
		{
			name:    "empty path",
			path:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSpaceID, gotBlobID, err := parseBlobPath(tt.path)
			if (err != nil) != tt.wantErr {
				t.Fatalf("parseBlobPath(%q) error = %v, wantErr %v", tt.path, err, tt.wantErr)
			}
			if err != nil {
				return
			}
			if gotSpaceID != tt.wantSpaceID {
				t.Errorf("parseBlobPath(%q) spaceID = %q, want %q", tt.path, gotSpaceID, tt.wantSpaceID)
			}
			if gotBlobID != tt.wantBlobID {
				t.Errorf("parseBlobPath(%q) blobID = %q, want %q", tt.path, gotBlobID, tt.wantBlobID)
			}
		})
	}
}
