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
		expected bool
		path     string
		name     string
	}{
		// API endpoints that should be allowed
		{expected: true, path: "/graph/v1.0/shares", name: "graph api endpoint"},
		{expected: true, path: "/graph/v1.0/shares/123", name: "graph api with path params"},
		{expected: true, path: "/graph/v1.0/users/me/drive/items/123:/testfile.txt:/permissions", name: "graph api with complex path"},
		{expected: true, path: "/ocs/v1.php/cloud/users", name: "ocs api endpoint"},
		{expected: true, path: "/ocs/v2.php/apps/files_sharing/api/v1/shares", name: "ocs api with subpath"},
		{expected: true, path: "/remote.php/webdav/testfile.txt", name: "remote.php endpoint"},
		{expected: true, path: "/remote.php/dav/files/alice/testfile.txt", name: "remote.php with dav"},

		// File serving paths that should be validated
		{expected: true, path: "/index.html", name: "simple file path"},
		{expected: true, path: "/assets/css/style.css", name: "subdirectory file"},
		{expected: true, path: "/js/vendor/jquery.min.js", name: "nested directory"},

		// Path traversal attempts that should be blocked
		{expected: false, path: "/assets/../secret.txt", name: "traversal with dots"},
		{expected: false, path: "/assets/%2e%2e/secret.txt", name: "traversal with encoded dots"},
		{expected: false, path: "/assets/./secret.txt", name: "current directory"},
		{expected: false, path: "/assets/%2e/secret.txt", name: "encoded current directory"},
		{expected: false, path: "/assets/css/../../secret.txt", name: "double traversal"},
		{expected: false, path: "/assets/../css/../../secret.txt", name: "mixed traversal"},

		// Edge cases
		{expected: true, path: "", name: "empty path"},
		{expected: true, path: "/", name: "root path"},
		{expected: true, path: "//assets/file.txt", name: "double slash"},
		{expected: true, path: "/assets/", name: "trailing slash"},

		// Additional test cases based on Go blog post about os.Root
		// URL-encoded traversal attempts
		{expected: false, path: "/assets/%2e%2e%2fsecret.txt", name: "URL encoded traversal"},
		{expected: false, path: "/assets/%2e%2e%5csecret.txt", name: "URL encoded backslash traversal"},
		{expected: false, path: "/assets/%252e%252e/secret.txt", name: "double URL encoded dots"},
		{expected: false, path: "/assets/%252e%252e%252fsecret.txt", name: "double URL encoded traversal"},

		// Unicode normalization attacks
		{expected: false, path: "/assets/..%c0%afsecret.txt", name: "UTF-8 encoded traversal"},
		{expected: false, path: "/assets/..%ef%bc%8fsecret.txt", name: "fullwidth slash traversal"},
		{expected: false, path: "/assets/..%c1%9csecret.txt", name: "UTF-8 encoded backslash"},

		// Multiple encoding layers
		{expected: false, path: "/assets/%252e%252e%252fsecret.txt", name: "double percent encoding"},
		{expected: false, path: "/assets/%252e%252e%255csecret.txt", name: "double percent encoding backslash"},

		// Null byte injection attempts
		{expected: false, path: "/assets/..%00/secret.txt", name: "null byte injection"},
		{expected: false, path: "/assets/%00../secret.txt", name: "null byte prefix"},

		// Mixed case and encoding
		{expected: false, path: "/assets/..%2Fsecret.txt", name: "mixed case URL encoding"},
		{expected: false, path: "/assets/..%2fsecret.txt", name: "lowercase URL encoding"},
		{expected: false, path: "/assets/..%5Csecret.txt", name: "uppercase backslash encoding"},

		// Path normalization edge cases
		{expected: false, path: "/assets/././../secret.txt", name: "redundant current dir traversal"},
		{expected: false, path: "/assets/.././secret.txt", name: "mixed traversal current dir"},
		{expected: false, path: "/assets/.../secret.txt", name: "triple dots"},
		{expected: false, path: "/assets/..../secret.txt", name: "quadruple dots"},

		// Paths that should be allowed (no traversal)
		{expected: true, path: "/assets/dots.in.filename.txt", name: "dots in filename"},
		{expected: true, path: "/assets/file..txt", name: "dots at end of filename"},
		{expected: true, path: "/assets/.hidden", name: "hidden file"},
		{expected: true, path: "/assets/..hidden", name: "file starting with dots"},
		{expected: true, path: "/assets/...hidden", name: "file starting with triple dots"},
		{expected: true, path: "/assets/normal/path/with/dots.txt", name: "normal path with dots"},
		{expected: true, path: "/assets/path/with/encoded%20spaces.txt", name: "URL encoded spaces"},
		{expected: true, path: "/assets/path/with/unicode%c3%a9.txt", name: "URL encoded unicode"},

		// API paths that might contain dots but should be allowed
		{expected: true, path: "/graph/v1.0/users/me/drive/items/123:/file..txt:/permissions", name: "graph api with dots in filename"},
		{expected: true, path: "/ocs/v1.php/cloud/users/..user", name: "ocs api with dots in username"},
		{expected: true, path: "/remote.php/dav/files/alice/.config", name: "remote.php with hidden file"},

		// e2e
		//   - navigating to "https://ocis-server:9200/files/spaces/personal/alice/parent/folder%252Fwith%252FSlashes?fileId=048eb01c-483c-4a9e-a3ac-17345d767500%24137d4fd3-afc7-4059-924b-834b0622f8c9%2158c2219f-2b84-444e-b814-8e077d4c496e&items-per-page=100&files-spaces-generic-view-mode=resource-table&tiles-size=2", waiting until "load"
		{expected: true, path: "files/spaces/personal/alice/parent/folder%252Fwith%252FSlashes", name: "e2e"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := IsSafePath(tc.path)
			if result != tc.expected {
				t.Errorf("isSafePath(%q) = %v, want %v", tc.path, result, tc.expected)
			}
		})
	}
}

// TestPathTraversalVulnerabilities tests various path traversal attack vectors
// mentioned in the Go blog post about os.Root https://go.dev/blog/osroot
func TestPathTraversalVulnerabilities(t *testing.T) {
	g := gomega.NewWithT(t)

	// Setup filesystem with sensitive files outside intended directory
	fsys := fstest.MapFS{
		"public": &fstest.MapFile{
			Mode: fs.ModeDir,
		},
		"public/legitimate.txt": &fstest.MapFile{
			Data: []byte("legitimate content"),
		},
		"secret.txt": &fstest.MapFile{
			Data: []byte("super-secret-content"),
		},
		"etc/passwd": &fstest.MapFile{
			Data: []byte("root:x:0:0:root:/root:/bin/bash"),
		},
		"config/database.yml": &fstest.MapFile{
			Data: []byte("database_password: secret123"),
		},
	}

	testCases := []struct {
		name            string
		path            string
		isBlockExpected bool
		description     string
	}{
		// Basic traversal attempts
		{
			name:            "basic_traversal",
			path:            "/public/../secret.txt",
			isBlockExpected: true,
			description:     "Basic directory traversal with ..",
		},
		{
			name:            "double_traversal",
			path:            "/public/subdir/../../secret.txt",
			isBlockExpected: true,
			description:     "Double directory traversal",
		},
		{
			name:            "mixed_traversal",
			path:            "/public/../config/../secret.txt",
			isBlockExpected: true,
			description:     "Mixed traversal with current directory",
		},

		// URL-encoded traversal attempts
		{
			name:            "url_encoded_traversal",
			path:            "/public/%2e%2e/secret.txt",
			isBlockExpected: true,
			description:     "URL-encoded dots",
		},
		{
			name:            "url_encoded_slash",
			path:            "/public/%2e%2e%2fsecret.txt",
			isBlockExpected: true,
			description:     "URL-encoded dots and slash",
		},
		{
			name:            "double_encoded",
			path:            "/public/%252e%252e/secret.txt",
			isBlockExpected: true,
			description:     "Double URL-encoded dots",
		},

		// Unicode normalization attacks
		{
			name:            "utf8_encoded_slash",
			path:            "/public/..%c0%afsecret.txt",
			isBlockExpected: true,
			description:     "UTF-8 encoded slash",
		},
		{
			name:            "fullwidth_slash",
			path:            "/public/..%ef%bc%8fsecret.txt",
			isBlockExpected: true,
			description:     "Fullwidth slash",
		},

		// Windows-specific attacks
		{
			name:            "windows_backslash",
			path:            "/public/..\\secret.txt",
			isBlockExpected: true,
			description:     "Windows backslash traversal",
		},
		{
			name:            "windows_encoded_backslash",
			path:            "/public/%2e%2e%5csecret.txt",
			isBlockExpected: true,
			description:     "URL-encoded Windows backslash",
		},

		// Null byte injection
		{
			name:            "null_byte_injection",
			path:            "/public/..%00/secret.txt",
			isBlockExpected: true,
			description:     "Null byte injection",
		},

		// Path normalization edge cases
		{
			name:            "redundant_current_dir",
			path:            "/public/././../secret.txt",
			isBlockExpected: true,
			description:     "Redundant current directory references",
		},
		{
			name:            "triple_dots",
			path:            "/public/.../secret.txt",
			isBlockExpected: true,
			description:     "Triple dots (should be treated as traversal)",
		},

		// Legitimate paths that should be allowed
		{
			name:            "legitimate_file",
			path:            "/public/legitimate.txt",
			isBlockExpected: false,
			description:     "Legitimate file access",
		},
		{
			name:            "dots_in_filename",
			path:            "/public/file..txt",
			isBlockExpected: true,
			description:     "Dots in filename (not traversal)",
		},
		{
			name:            "hidden_file",
			path:            "/public/.hidden",
			isBlockExpected: true,
			description:     "Hidden file (not traversal)",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tc.path, nil)
			w := httptest.NewRecorder()

			FileServer(fsys).ServeHTTP(w, req)
			res := w.Result()
			defer res.Body.Close()

			body, err := io.ReadAll(res.Body)
			g.Expect(err).ToNot(gomega.HaveOccurred())

			if tc.isBlockExpected {
				// Should be blocked - return 404 and not contain sensitive data
				g.Expect(res.StatusCode).To(gomega.Equal(http.StatusNotFound))
				g.Expect(string(body)).ToNot(gomega.ContainSubstring("super-secret"))
				g.Expect(string(body)).ToNot(gomega.ContainSubstring("database_password"))
				g.Expect(string(body)).ToNot(gomega.ContainSubstring("root:x:0:0"))
			} else {
				// Should be allowed - return 200 and contain expected content
				g.Expect(res.StatusCode).To(gomega.Equal(http.StatusOK))
				if tc.path == "/public/legitimate.txt" {
					g.Expect(string(body)).To(gomega.ContainSubstring("legitimate content"))
				}
			}
		})
	}
}
