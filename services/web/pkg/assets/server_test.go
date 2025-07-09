package assets

import (
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/fstest"

	"github.com/onsi/gomega"
)

func TestFileServer(t *testing.T) {
	g := gomega.NewWithT(t)
	recorderStatus := func(s int) string {
		return fmt.Sprintf("%03d %s", s, http.StatusText(s))
	}

	{
		s := FileServer(fstest.MapFS{})
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/foo", nil)
		//defer req.Body.Close()
		s.ServeHTTP(w, req)
		res := w.Result()
		defer res.Body.Close()

		g.Expect(res.Status).To(gomega.Equal(recorderStatus(http.StatusNotFound)))
	}

	for _, tt := range []struct {
		name     string
		url      string
		fs       fstest.MapFS
		expected string
	}{
		{
			name: "not found fallback",
			url:  "/index.txt",
			fs: fstest.MapFS{
				"index.html": &fstest.MapFile{
					Data: []byte("index file content"),
				},
			},
			expected: `<html><head><base href="/"/></head><body>index file content</body></html>`,
		},
		{
			name: "directory fallback",
			url:  "/some-folder",
			fs: fstest.MapFS{
				"some-folder": &fstest.MapFile{
					Mode: fs.ModeDir,
				},
				"index.html": &fstest.MapFile{
					Data: []byte("index file content"),
				},
			},
			expected: `<html><head><base href="/"/></head><body>index file content</body></html>`,
		},
		{
			name: "index.html",
			url:  "/index.html",
			fs: fstest.MapFS{
				"index.html": &fstest.MapFile{
					Data: []byte("index file content"),
				},
			},
			expected: `<html><head><base href="/"/></head><body>index file content</body></html>`,
		},
		{
			name: "oidc-callback.html",
			url:  "/oidc-callback.html",
			fs: fstest.MapFS{
				"index.html": &fstest.MapFile{
					Data: []byte("oidc-callback file content"),
				},
			},
			expected: `<html><head><base href="/"/></head><body>oidc-callback file content</body></html>`,
		},
		{
			name: "oidc-silent-redirect.html",
			url:  "/oidc-silent-redirect.html",
			fs: fstest.MapFS{
				"index.html": &fstest.MapFile{
					Data: []byte("oidc-silent-redirect file content"),
				},
			},
			expected: `<html><head><base href="/"/></head><body>oidc-silent-redirect file content</body></html>`,
		},
		{
			name: "some-file.txt",
			url:  "/some-file.txt",
			fs: fstest.MapFS{
				"some-file.txt": &fstest.MapFile{
					Data: []byte("some file content"),
				},
			},
			expected: "some file content",
		},
	} {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", tt.url, nil)
			FileServer(tt.fs).ServeHTTP(w, req)
			res := w.Result()
			defer res.Body.Close()

			g.Expect(res.Status).To(gomega.Equal(recorderStatus(http.StatusOK)))

			data, err := io.ReadAll(res.Body)
			g.Expect(err).ToNot(gomega.HaveOccurred())
			g.Expect(string(data)).To(gomega.Equal(tt.expected))

		})
	}
}

func TestIsSafePath(t *testing.T) {
	cases := []struct {
		permit bool
		path   string
	}{
		// simple cases
		{true, "/"},
		{true, "/index.html"},
		{true, "/some-file.txt"},

		// WebDAV double-slash variants that must be accepted
		{true, "//dav/spaces/123/textfile0.txt"},
		{true, "//dav//spaces/123/PARENT/parent.txt"},
		{true, "/dav//spaces/123/PARENT"},
		{true, "//dav/spaces/123//FOLDER"},

		// WebDAV files API variants
		{true, "//dav//files/alice/textfile1.txt"},
		{true, "/dav//files/alice/PARENT/parent.txt"},
		{true, "//dav/files/alice//FOLDER"},
		{true, "/dav/files/alice//PARENT5"},

		// WebDAV root variants
		{true, "//webdav/textfile0.txt"},
		{true, "/webdav//PARENT"},
		{true, "//webdav//PARENT3"},
		{true, "/webdav//textfile1.txt"},

		// Spaces API with double slashes
		{true, "//dav/spaces//SPACEID/PARENT4"},
		{true, "/dav/spaces//SPACEID/textfile7.txt"},
		{true, "//dav/spaces//SPACEID/PARENT//parent.txt"},

		// MOVE/COPY source patterns
		{true, "//dav/spaces/ID//PARENT1"},
		{true, "/dav//spaces/ID/textfile1.txt"},
		{true, "//dav/files//alice//PARENT1"},

		// Traversal attempts that must be rejected
		{false, "/dav/spaces/123/../../passwd"},
		{false, "../secret.txt"},
		{false, "/dav/%2e%2e/secret"},
		{false, "/dav/%2e/secret"},
		{false, "/dav/spaces/ID/../secret"},
		{false, "/dav/files/alice/././secret"},
	}

	for _, c := range cases {
		got := isSafePath(c.path)
		if got != c.permit {
			t.Errorf("isSafePath(%q) = %v, want %v", c.path, got, c.permit)
		}
	}
}

func TestFileServerPathTraversal(t *testing.T) {
	g := gomega.NewWithT(t)

	// setup in-memory filesystem that represents a "public" dir and a secret file at the root
	fsys := fstest.MapFS{
		"public": &fstest.MapFile{
			Mode: fs.ModeDir,
		},
		"secret.txt": &fstest.MapFile{ // file outside the intended public directory
			Data: []byte("super-secret"),
		},
	}

	// Request targets ../ traversal to reach secret.txt
	req := httptest.NewRequest("GET", "/public/../secret.txt", nil)
	w := httptest.NewRecorder()

	FileServer(fsys).ServeHTTP(w, req)
	res := w.Result()
	defer res.Body.Close()

	// Expect that the request is rejected due to path traversal attempt
	g.Expect(res.StatusCode).To(gomega.Equal(http.StatusNotFound))

	body, err := io.ReadAll(res.Body)
	g.Expect(err).ToNot(gomega.HaveOccurred())
	g.Expect(string(body)).ToNot(gomega.ContainSubstring("super-secret"))
}
