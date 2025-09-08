package assets

import (
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
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

// Combined OS-backed sanity tests
func TestFileServerSanity(t *testing.T) {
	g := gomega.NewWithT(t)

	cases := []struct {
		name       string
		setup      func(root string)
		path       string
		wantStatus int
		wantBody   string
	}{
		{
			name: "serve regular file",
			setup: func(root string) {
				_ = os.WriteFile(filepath.Join(root, "file.txt"), []byte("hello"), 0o644)
			},
			path:       "/file.txt",
			wantStatus: http.StatusOK,
			wantBody:   "hello",
		},
		{
			name: "reject symlink",
			setup: func(root string) {
				outside := t.TempDir()
				outsideFile := filepath.Join(outside, "link.txt")
				_ = os.WriteFile(outsideFile, []byte("link"), 0o644)
				_ = os.Symlink(outsideFile, filepath.Join(root, "link"))
			},
			path:       "/link",
			wantStatus: http.StatusNotFound,
		},
		{
			name: "parent file not accessible via traversal",
			setup: func(root string) {
				parent := filepath.Dir(root)
				_ = os.WriteFile(filepath.Join(parent, "outside.txt"), []byte("outside"), 0o644)
			},
			path:       "/../outside.txt",
			wantStatus: http.StatusNotFound,
		},
		{
			name: "index fallback serves index",
			setup: func(root string) {
				_ = os.WriteFile(filepath.Join(root, "index.html"), []byte("<html><head><title>index</title></head><body>index file content</body></html>"), 0o644)
			},
			path:       "/random",
			wantStatus: http.StatusOK,
			wantBody:   `<html><head><base href="/"/><title>index</title></head><body>index file content</body></html>`,
		},
		{
			name:       "resiliant",
			setup:      func(root string) {},
			path:       "random",
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			if tc.setup != nil {
				tc.setup(root)
			}
			h := &fileServer{fsys: http.Dir(root)}
			w := httptest.NewRecorder()
			req := &http.Request{Method: http.MethodGet, URL: &url.URL{Path: tc.path}, Header: make(http.Header)}
			h.ServeHTTP(w, req)
			res := w.Result()
			g.Expect(res.StatusCode).To(gomega.Equal(tc.wantStatus))
			if tc.wantBody != "" {
				data, _ := io.ReadAll(res.Body)
				g.Expect(string(data)).To(gomega.Equal(tc.wantBody))
			}
		})
	}
}
