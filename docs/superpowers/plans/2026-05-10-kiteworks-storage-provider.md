# Kiteworks Storage Provider Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Create a new `services/storage-kiteworks` oCIS service that implements the reva `storage.FS` interface backed by the Kiteworks REST API, exposing each top-level Kiteworks folder as a CS3 StorageSpace.

**Architecture:** A new independent oCIS service (`services/storage-kiteworks`) following the `storage-users` pattern: it wraps a reva gRPC+HTTP server with a new reva storage driver registered as `"kiteworks"`. The driver translates CS3 operations to Kiteworks REST calls using a copied/adapted client from `github.com/owncloud/kwdav`. OIDC tokens are passed through from the CS3 context to Kiteworks Bearer auth.

**Tech Stack:** Go 1.22, reva v2, CS3 APIs (go-cs3apis), github.com/urfave/cli/v2, github.com/rs/zerolog, github.com/mitchellh/mapstructure, Ginkgo v2 + Gomega for tests.

---

## File Map

### New files — reva driver (vendored)
- `vendor/github.com/owncloud/reva/v2/pkg/storage/fs/kiteworks/types.go` — Kiteworks data model (FileInfo, DirectoryInfo, etc.)
- `vendor/github.com/owncloud/reva/v2/pkg/storage/fs/kiteworks/client.go` — REST API client
- `vendor/github.com/owncloud/reva/v2/pkg/storage/fs/kiteworks/upload.go` — chunked upload logic
- `vendor/github.com/owncloud/reva/v2/pkg/storage/fs/kiteworks/kiteworks.go` — `storage.FS` implementation + `init()` registration
- `vendor/github.com/owncloud/reva/v2/pkg/storage/fs/kiteworks/kiteworks_test.go` — Ginkgo suite bootstrap
- `vendor/github.com/owncloud/reva/v2/pkg/storage/fs/kiteworks/mock_server_test.go` — httptest mock server

### Modified files — reva driver loader
- `vendor/github.com/owncloud/reva/v2/pkg/storage/fs/loader/loader.go` — add blank import for kiteworks driver

### New files — oCIS service
- `services/storage-kiteworks/cmd/storage-kiteworks/main.go` — binary entry point
- `services/storage-kiteworks/pkg/command/root.go` — GetCommands, Execute
- `services/storage-kiteworks/pkg/command/server.go` — Server command, reva runtime wiring
- `services/storage-kiteworks/pkg/command/health.go` — Health command
- `services/storage-kiteworks/pkg/command/version.go` — Version command
- `services/storage-kiteworks/pkg/config/config.go` — Config, KiteworksDriver structs
- `services/storage-kiteworks/pkg/config/reva.go` — TokenManager struct
- `services/storage-kiteworks/pkg/config/defaults/defaultconfig.go` — DefaultConfig, EnsureDefaults, Sanitize, FullDefaultConfig
- `services/storage-kiteworks/pkg/config/parser/parse.go` — ParseConfig, Validate
- `services/storage-kiteworks/pkg/revaconfig/config.go` — StorageKiteworksConfigFromStruct
- `services/storage-kiteworks/pkg/logging/logging.go` — Configure
- `services/storage-kiteworks/pkg/server/debug/server.go` — debug HTTP server
- `services/storage-kiteworks/pkg/server/debug/options.go` — debug server options

### Modified files — oCIS integration
- `ocis-pkg/config/config.go` — add `StorageKiteworks *storagekiteworks.Config`
- `ocis-pkg/config/defaults/defaultconfig.go` — add StorageKiteworks default
- `ocis/pkg/command/services.go` — register storage-kiteworks CLI command
- `ocis/pkg/runtime/service/service.go` — register storage-kiteworks runtime service
- `services/storage-kiteworks/go.mod` — module definition (if separate module; otherwise use repo root)

---

## Task 1: Kiteworks data types

**Files:**
- Create: `vendor/github.com/owncloud/reva/v2/pkg/storage/fs/kiteworks/types.go`

- [ ] **Step 1: Create types.go with all Kiteworks data model types**

```go
// vendor/github.com/owncloud/reva/v2/pkg/storage/fs/kiteworks/types.go
package kiteworks

import (
	"strings"
	"time"
)

// Time is a custom time type that handles Kiteworks date format
type Time struct {
	time.Time
}

func (t *Time) UnmarshalJSON(data []byte) error {
	s := strings.Trim(string(data), `"`)
	if s == "null" || s == "" {
		return nil
	}
	parsed, err := time.Parse("2006-01-02T15:04:05+0000", s)
	if err != nil {
		parsed, err = time.Parse(time.RFC3339, s)
		if err != nil {
			return err
		}
	}
	t.Time = parsed
	return nil
}

// FileFingerPrint holds a single hash entry
type FileFingerPrint struct {
	Algo string `json:"algo"`
	Hash string `json:"hash"`
}

// FileFingerPrints holds all checksums for a file
type FileFingerPrints struct {
	FingerPrints []FileFingerPrint `json:"fingerprints"`
}

// Permission represents a Kiteworks permission entry
type Permission struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Allowed bool   `json:"allowed"`
}

// FileInfo is the metadata for a Kiteworks file or folder
type FileInfo struct {
	ID           string            `json:"id"`
	Name         string            `json:"name"`
	Path         string            `json:"path"`
	Type         string            `json:"type"` // "d" = directory, "f" = file
	Size         int64             `json:"size"`
	Modified     *Time             `json:"modified"`
	Created      *Time             `json:"created"`
	ParentID     *string           `json:"parentId"`
	IsShared     *bool             `json:"shared"`
	Permissions  []Permission      `json:"currentUserPermissions"`
	FingerPrints *FileFingerPrints `json:"fileFingerprints"`
}

// IsDir returns true if the FileInfo represents a directory
func (fi *FileInfo) IsDir() bool {
	return fi.Type == "d"
}

// IsSharedWithUser returns true when this top-level folder is a received share
func (fi *FileInfo) IsSharedWithUser() bool {
	if fi == nil || fi.IsShared == nil || !*fi.IsShared {
		return false
	}
	if fi.ParentID == nil {
		return true
	}
	return *fi.ParentID == "0"
}

// DirectoryInfo wraps a listing of files and folders
type DirectoryInfo struct {
	Files   []FileInfo `json:"files"`
	Folders []FileInfo `json:"folders"`
}

// CreateDir is the request body for creating a folder
type CreateDir struct {
	Name     string `json:"name"`
	Syncable bool   `json:"syncable"`
}

// InitializeUpload is the request body for initiating a chunked upload
type InitializeUpload struct {
	Filename      string `json:"filename"`
	TotalChunks   int    `json:"totalChunks"`
	TotalFileSize int64  `json:"totalFileSize"`
}

// UploadResult is the response from initiating an upload
type UploadResult struct {
	ID  string `json:"id"`
	URI string `json:"uri"`
}

// FileCopyMove is the request body for move/copy operations
type FileCopyMove struct {
	DestFolderID string `json:"destFolderId"`
	Replace      bool   `json:"replace"`
}

// FileUpdateRequest is the request body for renaming a file
type FileUpdateRequest struct {
	Name    string `json:"name"`
	Replace bool   `json:"replace"`
}

// FolderUpdateRequest is the request body for renaming a folder
type FolderUpdateRequest struct {
	Name string `json:"name"`
}

// QuotaInfo holds user storage quota information
type QuotaInfo struct {
	Allowed int64 `json:"storageQuota"`
	Used    int64 `json:"storageUsed"`
}

// User holds Kiteworks user information
type User struct {
	ID    string    `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
	Quota QuotaInfo `json:"quota"`
}
```

- [ ] **Step 2: Compile check**

```bash
cd /home/deepdiver/Development/claude-work/ocis
go build ./vendor/github.com/owncloud/reva/v2/pkg/storage/fs/kiteworks/...
```
Expected: no output (success).

- [ ] **Step 3: Commit**

```bash
git add vendor/github.com/owncloud/reva/v2/pkg/storage/fs/kiteworks/types.go
git commit -s -m "feat(kiteworks): add Kiteworks data model types"
```

---

## Task 2: Kiteworks REST client

**Files:**
- Create: `vendor/github.com/owncloud/reva/v2/pkg/storage/fs/kiteworks/client.go`

- [ ] **Step 1: Create client.go**

```go
// vendor/github.com/owncloud/reva/v2/pkg/storage/fs/kiteworks/client.go
package kiteworks

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
)

const kiteworksAPIVersion = "28"

// ClientError wraps an unexpected HTTP status
type ClientError struct {
	StatusCode int
	Body       []byte
}

func (e *ClientError) Error() string {
	return fmt.Sprintf("kiteworks: unexpected status %d: %s", e.StatusCode, e.Body)
}

// Client is a Kiteworks REST API client scoped to a single user token
type Client struct {
	endpoint   string
	token      string
	httpClient *http.Client
}

// NewClient creates a Client. Set insecure=true only for development/testing.
func NewClient(endpoint, token string, insecure bool) *Client {
	transport := http.DefaultTransport
	if insecure {
		transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, //nolint:gosec
		}
	}
	return &Client{
		endpoint:   endpoint,
		token:      token,
		httpClient: &http.Client{Transport: transport},
	}
}

func (c *Client) newRequest(method, path string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, c.endpoint+path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("X-Accellion-Version", kiteworksAPIVersion)
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

func (c *Client) do(req *http.Request) (*http.Response, error) {
	return c.httpClient.Do(req)
}

func (c *Client) doJSON(method, path string, body interface{}, out interface{}) error {
	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return err
		}
		bodyReader = bytes.NewReader(data)
	}
	req, err := c.newRequest(method, path, bodyReader)
	if err != nil {
		return err
	}
	resp, err := c.do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode >= 300 {
		return &ClientError{StatusCode: resp.StatusCode, Body: respBody}
	}
	if out != nil {
		return json.Unmarshal(respBody, out)
	}
	return nil
}

// GetTopFolders returns the user's top-level folders
func (c *Client) GetTopFolders() ([]FileInfo, error) {
	var result struct {
		Data []FileInfo `json:"data"`
	}
	err := c.doJSON(http.MethodGet, "/rest/folders/top", nil, &result)
	return result.Data, err
}

// GetFolder returns folder metadata by ID
func (c *Client) GetFolder(id string) (*FileInfo, error) {
	var fi FileInfo
	err := c.doJSON(http.MethodGet, "/rest/folders/"+id, nil, &fi)
	return &fi, err
}

// ListFolder returns the children of a folder
func (c *Client) ListFolder(id string) (*DirectoryInfo, error) {
	var result DirectoryInfo
	err := c.doJSON(http.MethodGet, "/rest/folders/"+id+"/children", nil, &result)
	return &result, err
}

// CreateFolder creates a subfolder inside parent
func (c *Client) CreateFolder(parentID, name string) (*FileInfo, error) {
	var fi FileInfo
	err := c.doJSON(http.MethodPost, "/rest/folders/"+parentID+"/folders",
		&CreateDir{Name: name, Syncable: true}, &fi)
	return &fi, err
}

// DeleteFolder deletes a folder by ID
func (c *Client) DeleteFolder(id string) error {
	return c.doJSON(http.MethodDelete, "/rest/folders/"+id, nil, nil)
}

// RenameFolder renames a folder
func (c *Client) RenameFolder(id, name string) error {
	return c.doJSON(http.MethodPut, "/rest/folders/"+id,
		&FolderUpdateRequest{Name: name}, nil)
}

// GetFile returns file metadata by ID
func (c *Client) GetFile(id string) (*FileInfo, error) {
	var fi FileInfo
	err := c.doJSON(http.MethodGet, "/rest/files/"+id, nil, &fi)
	return &fi, err
}

// DownloadFile returns a ReadCloser for the file content. The caller must close it.
// rangeHeader is optional (e.g. "bytes=0-1023").
func (c *Client) DownloadFile(id, rangeHeader string) (*http.Response, error) {
	req, err := c.newRequest(http.MethodGet, "/rest/files/"+id+"/content", nil)
	if err != nil {
		return nil, err
	}
	if rangeHeader != "" {
		req.Header.Set("Range", rangeHeader)
	}
	return c.do(req)
}

// DeleteFile deletes a file by ID
func (c *Client) DeleteFile(id string) error {
	return c.doJSON(http.MethodDelete, "/rest/files/"+id, nil, nil)
}

// RenameFile renames a file
func (c *Client) RenameFile(id, name string, replace bool) error {
	return c.doJSON(http.MethodPut, "/rest/files/"+id,
		&FileUpdateRequest{Name: name, Replace: replace}, nil)
}

// MoveResource moves a file or folder to a new parent folder
func (c *Client) MoveResource(sourceID, destFolderID string, replace bool) error {
	return c.doJSON(http.MethodPost, "/rest/files/actions/move",
		&FileCopyMove{DestFolderID: destFolderID, Replace: replace}, nil)
}

// CopyResource copies a file or folder to a destination folder
func (c *Client) CopyResource(sourceID, destFolderID string, replace bool) error {
	return c.doJSON(http.MethodPost, "/rest/files/actions/copy",
		&FileCopyMove{DestFolderID: destFolderID, Replace: replace}, nil)
}

// InitiateUpload starts a chunked upload session
func (c *Client) InitiateUpload(parentID, filename string, size int64, numChunks int) (*UploadResult, error) {
	var result UploadResult
	err := c.doJSON(http.MethodPost, "/rest/folders/"+parentID+"/actions/initiateUpload",
		&InitializeUpload{
			Filename:      filename,
			TotalChunks:   numChunks,
			TotalFileSize: size,
		}, &result)
	return &result, err
}

// UploadChunk uploads a single chunk. Returns FileInfo only on the last chunk.
func (c *Client) UploadChunk(uploadURI, filename string, data []byte, chunkIndex int, isLastChunk bool) (*FileInfo, error) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	_ = w.WriteField("index", strconv.Itoa(chunkIndex))
	_ = w.WriteField("compressionMode", "NORMAL")
	if isLastChunk {
		_ = w.WriteField("lastChunk", "1")
	}
	part, err := w.CreateFormFile("file", filename)
	if err != nil {
		return nil, err
	}
	if _, err := part.Write(data); err != nil {
		return nil, err
	}
	w.Close()

	req, err := http.NewRequest(http.MethodPost, c.endpoint+uploadURI, &b)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("X-Accellion-Version", kiteworksAPIVersion)
	req.Header.Set("Content-Type", w.FormDataContentType())

	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 300 {
		return nil, &ClientError{StatusCode: resp.StatusCode, Body: respBody}
	}
	if !isLastChunk {
		return nil, nil
	}
	var fi FileInfo
	if err := json.Unmarshal(respBody, &fi); err != nil {
		return nil, err
	}
	return &fi, nil
}

// GetMe returns the current user's info including quota
func (c *Client) GetMe() (*User, error) {
	var u User
	err := c.doJSON(http.MethodGet, "/rest/users/me", nil, &u)
	return &u, err
}

// Search searches for a file/folder by path
func (c *Client) Search(path string) ([]FileInfo, error) {
	var result struct {
		Files   []FileInfo `json:"files"`
		Folders []FileInfo `json:"folders"`
	}
	err := c.doJSON(http.MethodGet, "/rest/query?query="+path, nil, &result)
	if err != nil {
		return nil, err
	}
	items := append(result.Folders, result.Files...)
	return items, nil
}
```

- [ ] **Step 2: Compile check**

```bash
go build ./vendor/github.com/owncloud/reva/v2/pkg/storage/fs/kiteworks/...
```
Expected: no output.

- [ ] **Step 3: Commit**

```bash
git add vendor/github.com/owncloud/reva/v2/pkg/storage/fs/kiteworks/client.go
git commit -s -m "feat(kiteworks): add Kiteworks REST API client"
```

---

## Task 3: Client unit tests

**Files:**
- Create: `vendor/github.com/owncloud/reva/v2/pkg/storage/fs/kiteworks/mock_server_test.go`
- Create: `vendor/github.com/owncloud/reva/v2/pkg/storage/fs/kiteworks/kiteworks_test.go`
- Create: `vendor/github.com/owncloud/reva/v2/pkg/storage/fs/kiteworks/client_test.go`

- [ ] **Step 1: Write the Ginkgo suite bootstrap**

```go
// vendor/github.com/owncloud/reva/v2/pkg/storage/fs/kiteworks/kiteworks_test.go
package kiteworks_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestKiteworks(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Kiteworks Storage Driver Suite")
}
```

- [ ] **Step 2: Write the mock server**

```go
// vendor/github.com/owncloud/reva/v2/pkg/storage/fs/kiteworks/mock_server_test.go
package kiteworks_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"

	. "github.com/owncloud/reva/v2/pkg/storage/fs/kiteworks"
)

func newMockServer() *httptest.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/rest/users/me", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(User{
			ID:    "u1",
			Name:  "Alice",
			Email: "alice@example.com",
			Quota: QuotaInfo{Allowed: 10737418240, Used: 1073741824},
		})
	})

	mux.HandleFunc("/rest/folders/top", func(w http.ResponseWriter, r *http.Request) {
		isShared := false
		parentID := "0"
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": []FileInfo{
				{ID: "f1", Name: "MyFiles", Type: "d", IsShared: &isShared},
				{ID: "f2", Name: "SharedFolder", Type: "d", IsShared: func() *bool { b := true; return &b }(), ParentID: &parentID},
			},
		})
	})

	mux.HandleFunc("/rest/folders/f1", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		isShared := false
		json.NewEncoder(w).Encode(FileInfo{ID: "f1", Name: "MyFiles", Type: "d", IsShared: &isShared})
	})

	mux.HandleFunc("/rest/folders/f1/children", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(DirectoryInfo{
			Files:   []FileInfo{{ID: "file1", Name: "hello.txt", Type: "f", Size: 5}},
			Folders: []FileInfo{},
		})
	})

	mux.HandleFunc("/rest/files/file1", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(FileInfo{ID: "file1", Name: "hello.txt", Type: "f", Size: 5})
	})

	mux.HandleFunc("/rest/files/file1/content", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("hello"))
	})

	return httptest.NewServer(mux)
}
```

- [ ] **Step 3: Write client tests**

```go
// vendor/github.com/owncloud/reva/v2/pkg/storage/fs/kiteworks/client_test.go
package kiteworks_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/owncloud/reva/v2/pkg/storage/fs/kiteworks"
)

var _ = Describe("Client", func() {
	var (
		srv    *httptest.Server
		client *Client
	)

	BeforeEach(func() {
		srv = newMockServer()
		client = NewClient(srv.URL, "test-token", false)
	})

	AfterEach(func() {
		srv.Close()
	})

	Describe("GetTopFolders", func() {
		It("returns two top-level folders", func() {
			folders, err := client.GetTopFolders()
			Expect(err).ToNot(HaveOccurred())
			Expect(folders).To(HaveLen(2))
			Expect(folders[0].ID).To(Equal("f1"))
			Expect(folders[1].ID).To(Equal("f2"))
		})
	})

	Describe("GetFolder", func() {
		It("returns folder metadata", func() {
			fi, err := client.GetFolder("f1")
			Expect(err).ToNot(HaveOccurred())
			Expect(fi.ID).To(Equal("f1"))
			Expect(fi.Name).To(Equal("MyFiles"))
		})
	})

	Describe("ListFolder", func() {
		It("returns folder children", func() {
			dir, err := client.ListFolder("f1")
			Expect(err).ToNot(HaveOccurred())
			Expect(dir.Files).To(HaveLen(1))
			Expect(dir.Files[0].Name).To(Equal("hello.txt"))
		})
	})

	Describe("GetMe", func() {
		It("returns current user with quota", func() {
			u, err := client.GetMe()
			Expect(err).ToNot(HaveOccurred())
			Expect(u.ID).To(Equal("u1"))
			Expect(u.Quota.Allowed).To(Equal(int64(10737418240)))
		})
	})

	Describe("IsSharedWithUser", func() {
		It("returns false for owned folder", func() {
			folders, _ := client.GetTopFolders()
			Expect(folders[0].IsSharedWithUser()).To(BeFalse())
		})
		It("returns true for received share", func() {
			folders, _ := client.GetTopFolders()
			Expect(folders[1].IsSharedWithUser()).To(BeTrue())
		})
	})
})
```

Note: add `"net/http/httptest"` import to `client_test.go`.

- [ ] **Step 4: Run tests to verify they pass**

```bash
go test ./vendor/github.com/owncloud/reva/v2/pkg/storage/fs/kiteworks/... -v
```
Expected: All tests PASS.

- [ ] **Step 5: Commit**

```bash
git add vendor/github.com/owncloud/reva/v2/pkg/storage/fs/kiteworks/
git commit -s -m "test(kiteworks): add client unit tests with mock server"
```

---

## Task 4: Chunked upload logic

**Files:**
- Create: `vendor/github.com/owncloud/reva/v2/pkg/storage/fs/kiteworks/upload.go`

- [ ] **Step 1: Write the upload helper**

```go
// vendor/github.com/owncloud/reva/v2/pkg/storage/fs/kiteworks/upload.go
package kiteworks

import (
	"io"
	"math"
)

const defaultChunkSize = 5 * 1024 * 1024 // 5 MB

// uploadFile performs a chunked upload of r into parentFolderID.
// chunkSize of 0 uses defaultChunkSize.
func uploadFile(c *Client, parentFolderID, filename string, size int64, r io.Reader, chunkSize int64) (*FileInfo, error) {
	if chunkSize <= 0 {
		chunkSize = defaultChunkSize
	}

	numChunks := 1
	if size > 0 {
		numChunks = int(math.Ceil(float64(size) / float64(chunkSize)))
	}

	result, err := c.InitiateUpload(parentFolderID, filename, size, numChunks)
	if err != nil {
		return nil, err
	}

	buf := make([]byte, chunkSize)
	var fi *FileInfo
	for i := 0; i < numChunks; i++ {
		n, err := io.ReadFull(r, buf)
		if err != nil && err != io.ErrUnexpectedEOF && err != io.EOF {
			return nil, err
		}
		isLast := i == numChunks-1
		fi, err = c.UploadChunk(result.URI, filename, buf[:n], i, isLast)
		if err != nil {
			return nil, err
		}
	}
	return fi, nil
}
```

- [ ] **Step 2: Compile check**

```bash
go build ./vendor/github.com/owncloud/reva/v2/pkg/storage/fs/kiteworks/...
```
Expected: no output.

- [ ] **Step 3: Commit**

```bash
git add vendor/github.com/owncloud/reva/v2/pkg/storage/fs/kiteworks/upload.go
git commit -s -m "feat(kiteworks): add chunked upload helper"
```

---

## Task 5: storage.FS implementation — spaces and read operations

**Files:**
- Create: `vendor/github.com/owncloud/reva/v2/pkg/storage/fs/kiteworks/kiteworks.go`

This is the main driver file. We implement it in two steps: first the read path (spaces, GetMD, ListFolder, Download), then the write path in Task 6.

- [ ] **Step 1: Create kiteworks.go with registration, struct, and read operations**

```go
// vendor/github.com/owncloud/reva/v2/pkg/storage/fs/kiteworks/kiteworks.go
package kiteworks

import (
	"context"
	"io"
	"net/url"
	"time"

	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	typespb "github.com/cs3org/go-cs3apis/cs3/types/v1beta1"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	ctxpkg "github.com/owncloud/reva/v2/pkg/ctx"
	"github.com/owncloud/reva/v2/pkg/errtypes"
	"github.com/owncloud/reva/v2/pkg/events"
	"github.com/owncloud/reva/v2/pkg/storage"
	"github.com/owncloud/reva/v2/pkg/storage/fs/registry"
)

func init() {
	registry.Register("kiteworks", New)
}

// Config holds the driver configuration decoded from reva mapstructure
type Config struct {
	Endpoint  string `mapstructure:"endpoint"`
	Insecure  bool   `mapstructure:"insecure"`
	ChunkSize int64  `mapstructure:"chunk_size"`
}

// Driver implements storage.FS backed by the Kiteworks REST API
type Driver struct {
	cfg *Config
}

// New creates a new kiteworks storage driver
func New(m map[string]interface{}, _ events.Stream, _ *zerolog.Logger) (storage.FS, error) {
	cfg := &Config{}
	if err := mapstructure.Decode(m, cfg); err != nil {
		return nil, errors.Wrap(err, "kiteworks: error decoding config")
	}
	if cfg.Endpoint == "" {
		return nil, errors.New("kiteworks: 'endpoint' must be set")
	}
	if cfg.ChunkSize <= 0 {
		cfg.ChunkSize = defaultChunkSize
	}
	return &Driver{cfg: cfg}, nil
}

func (d *Driver) client(ctx context.Context) (*Client, error) {
	token, ok := ctxpkg.ContextGetToken(ctx)
	if !ok || token == "" {
		return nil, errtypes.PermissionDenied("kiteworks: no token in context")
	}
	return NewClient(d.cfg.Endpoint, token, d.cfg.Insecure), nil
}

// fileInfoToResourceInfo converts a Kiteworks FileInfo to a CS3 ResourceInfo
func fileInfoToResourceInfo(fi *FileInfo) *provider.ResourceInfo {
	ri := &provider.ResourceInfo{
		Id: &provider.ResourceId{
			StorageId: "kiteworks",
			OpaqueId:  fi.ID,
		},
		Path: fi.Name,
		Type: provider.ResourceType_RESOURCE_TYPE_FILE,
		Size: uint64(fi.Size),
	}
	if fi.IsDir() {
		ri.Type = provider.ResourceType_RESOURCE_TYPE_CONTAINER
	}
	if fi.Modified != nil {
		ri.Mtime = &typespb.Timestamp{
			Seconds: uint64(fi.Modified.Unix()),
		}
	}
	if fi.FingerPrints != nil {
		for _, fp := range fi.FingerPrints.FingerPrints {
			if fp.Algo == "sha256" {
				ri.Checksum = &provider.ResourceChecksum{
					Type: provider.ResourceChecksumType_RESOURCE_CHECKSUM_TYPE_SHA256,
					Sum:  fp.Hash,
				}
			}
		}
	}
	return ri
}

// spaceFromFileInfo converts a top-level Kiteworks folder to a CS3 StorageSpace
func spaceFromFileInfo(fi *FileInfo) *provider.StorageSpace {
	spaceType := "project"
	if fi.IsSharedWithUser() {
		spaceType = "mountpoint"
	}
	sp := &provider.StorageSpace{
		Id: &provider.StorageSpaceId{
			OpaqueId: fi.ID,
		},
		Root: &provider.ResourceId{
			StorageId: "kiteworks",
			OpaqueId:  fi.ID,
		},
		Name:      fi.Name,
		SpaceType: spaceType,
	}
	if fi.Modified != nil {
		sp.Mtime = &typespb.Timestamp{
			Seconds: uint64(fi.Modified.Unix()),
		}
	}
	return sp
}

// Shutdown implements storage.FS
func (d *Driver) Shutdown(_ context.Context) error { return nil }

// ListStorageSpaces implements storage.FS — returns each top-level Kiteworks folder as a space
func (d *Driver) ListStorageSpaces(ctx context.Context, _ []*provider.ListStorageSpacesRequest_Filter, _ bool) ([]*provider.StorageSpace, error) {
	c, err := d.client(ctx)
	if err != nil {
		return nil, err
	}
	folders, err := c.GetTopFolders()
	if err != nil {
		return nil, err
	}
	spaces := make([]*provider.StorageSpace, 0, len(folders))
	for i := range folders {
		spaces = append(spaces, spaceFromFileInfo(&folders[i]))
	}
	return spaces, nil
}

// GetQuota implements storage.FS
func (d *Driver) GetQuota(ctx context.Context, _ *provider.Reference) (uint64, uint64, uint64, error) {
	c, err := d.client(ctx)
	if err != nil {
		return 0, 0, 0, err
	}
	u, err := c.GetMe()
	if err != nil {
		return 0, 0, 0, err
	}
	total := uint64(u.Quota.Allowed)
	used := uint64(u.Quota.Used)
	var remaining uint64
	if total > used {
		remaining = total - used
	}
	return total, used, remaining, nil
}

// GetMD implements storage.FS
func (d *Driver) GetMD(ctx context.Context, ref *provider.Reference, _, _ []string) (*provider.ResourceInfo, error) {
	c, err := d.client(ctx)
	if err != nil {
		return nil, err
	}
	id := ref.GetResourceId().GetOpaqueId()
	if id == "" && ref.GetPath() != "" {
		// path-based lookup via search
		results, err := c.Search(ref.GetPath())
		if err != nil {
			return nil, err
		}
		if len(results) == 0 {
			return nil, errtypes.NotFound(ref.GetPath())
		}
		return fileInfoToResourceInfo(&results[0]), nil
	}
	// Try folder first, fall back to file
	fi, err := c.GetFolder(id)
	if err != nil {
		if ce, ok := err.(*ClientError); ok && ce.StatusCode == 404 {
			fi, err = c.GetFile(id)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	return fileInfoToResourceInfo(fi), nil
}

// ListFolder implements storage.FS
func (d *Driver) ListFolder(ctx context.Context, ref *provider.Reference, _, _ []string) ([]*provider.ResourceInfo, error) {
	c, err := d.client(ctx)
	if err != nil {
		return nil, err
	}
	id := ref.GetResourceId().GetOpaqueId()
	dir, err := c.ListFolder(id)
	if err != nil {
		return nil, err
	}
	var infos []*provider.ResourceInfo
	for i := range dir.Folders {
		infos = append(infos, fileInfoToResourceInfo(&dir.Folders[i]))
	}
	for i := range dir.Files {
		infos = append(infos, fileInfoToResourceInfo(&dir.Files[i]))
	}
	return infos, nil
}

// Download implements storage.FS
func (d *Driver) Download(ctx context.Context, ref *provider.Reference, openReaderFunc func(*provider.ResourceInfo) bool) (*provider.ResourceInfo, io.ReadCloser, error) {
	c, err := d.client(ctx)
	if err != nil {
		return nil, nil, err
	}
	id := ref.GetResourceId().GetOpaqueId()
	fi, err := c.GetFile(id)
	if err != nil {
		return nil, nil, err
	}
	ri := fileInfoToResourceInfo(fi)
	if openReaderFunc != nil && !openReaderFunc(ri) {
		return ri, nil, nil
	}
	resp, err := c.DownloadFile(id, "")
	if err != nil {
		return nil, nil, err
	}
	return ri, resp.Body, nil
}

// GetPathByID implements storage.FS
func (d *Driver) GetPathByID(ctx context.Context, id *provider.ResourceId) (string, error) {
	c, err := d.client(ctx)
	if err != nil {
		return "", err
	}
	// Try folder first, fall back to file
	fi, err := c.GetFolder(id.OpaqueId)
	if err != nil {
		if ce, ok := err.(*ClientError); ok && ce.StatusCode == 404 {
			fi, err = c.GetFile(id.OpaqueId)
			if err != nil {
				return "", err
			}
		} else {
			return "", err
		}
	}
	return fi.Path, nil
}

// --- Stubbed / not-supported operations ---

func (d *Driver) CreateStorageSpace(_ context.Context, _ *provider.CreateStorageSpaceRequest) (*provider.CreateStorageSpaceResponse, error) {
	return nil, errtypes.NotSupported("kiteworks: CreateStorageSpace")
}
func (d *Driver) UpdateStorageSpace(_ context.Context, _ *provider.UpdateStorageSpaceRequest) (*provider.UpdateStorageSpaceResponse, error) {
	return nil, errtypes.NotSupported("kiteworks: UpdateStorageSpace")
}
func (d *Driver) DeleteStorageSpace(_ context.Context, _ *provider.DeleteStorageSpaceRequest) error {
	return errtypes.NotSupported("kiteworks: DeleteStorageSpace")
}
func (d *Driver) CreateHome(_ context.Context) error {
	return errtypes.NotSupported("kiteworks: CreateHome")
}
func (d *Driver) GetHome(_ context.Context) (string, error) {
	return "", errtypes.NotSupported("kiteworks: GetHome")
}
func (d *Driver) CreateReference(_ context.Context, _ string, _ *url.URL) error {
	return errtypes.NotSupported("kiteworks: CreateReference")
}
func (d *Driver) ListRevisions(_ context.Context, _ *provider.Reference) ([]*provider.FileVersion, error) {
	return nil, errtypes.NotSupported("kiteworks: ListRevisions")
}
func (d *Driver) DownloadRevision(_ context.Context, _ *provider.Reference, _ string, _ func(*provider.ResourceInfo) bool) (*provider.ResourceInfo, io.ReadCloser, error) {
	return nil, nil, errtypes.NotSupported("kiteworks: DownloadRevision")
}
func (d *Driver) RestoreRevision(_ context.Context, _ *provider.Reference, _ string) error {
	return errtypes.NotSupported("kiteworks: RestoreRevision")
}
func (d *Driver) ListRecycle(_ context.Context, _ *provider.Reference, _, _ string) ([]*provider.RecycleItem, error) {
	return nil, errtypes.NotSupported("kiteworks: ListRecycle")
}
func (d *Driver) RestoreRecycleItem(_ context.Context, _ *provider.Reference, _, _ string, _ *provider.Reference) error {
	return errtypes.NotSupported("kiteworks: RestoreRecycleItem")
}
func (d *Driver) PurgeRecycleItem(_ context.Context, _ *provider.Reference, _, _ string) error {
	return errtypes.NotSupported("kiteworks: PurgeRecycleItem")
}
func (d *Driver) EmptyRecycle(_ context.Context, _ *provider.Reference) error {
	return errtypes.NotSupported("kiteworks: EmptyRecycle")
}
func (d *Driver) DenyGrant(_ context.Context, _ *provider.Reference, _ *provider.Grantee) error {
	return errtypes.NotSupported("kiteworks: DenyGrant")
}
func (d *Driver) SetArbitraryMetadata(_ context.Context, _ *provider.Reference, _ *provider.ArbitraryMetadata) error {
	return errtypes.NotSupported("kiteworks: SetArbitraryMetadata")
}
func (d *Driver) UnsetArbitraryMetadata(_ context.Context, _ *provider.Reference, _ []string) error {
	return errtypes.NotSupported("kiteworks: UnsetArbitraryMetadata")
}
func (d *Driver) GetLock(_ context.Context, _ *provider.Reference) (*provider.Lock, error) {
	return nil, errtypes.NotSupported("kiteworks: GetLock")
}
func (d *Driver) SetLock(_ context.Context, _ *provider.Reference, _ *provider.Lock) error {
	return errtypes.NotSupported("kiteworks: SetLock")
}
func (d *Driver) RefreshLock(_ context.Context, _ *provider.Reference, _ *provider.Lock, _ string) error {
	return errtypes.NotSupported("kiteworks: RefreshLock")
}
func (d *Driver) Unlock(_ context.Context, _ *provider.Reference, _ *provider.Lock) error {
	return errtypes.NotSupported("kiteworks: Unlock")
}
```

- [ ] **Step 2: Compile check**

```bash
go build ./vendor/github.com/owncloud/reva/v2/pkg/storage/fs/kiteworks/...
```
Expected: no output.

- [ ] **Step 3: Commit**

```bash
git add vendor/github.com/owncloud/reva/v2/pkg/storage/fs/kiteworks/kiteworks.go
git commit -s -m "feat(kiteworks): implement storage.FS read path and space listing"
```

---

## Task 6: storage.FS — write operations

**Files:**
- Modify: `vendor/github.com/owncloud/reva/v2/pkg/storage/fs/kiteworks/kiteworks.go`

Add the write methods to the Driver. Open the file and append these methods before the last `}`.

- [ ] **Step 1: Add write methods to kiteworks.go**

Append to `kiteworks.go` (before the final closing brace or as standalone functions in the package):

```go
// CreateDir implements storage.FS
func (d *Driver) CreateDir(ctx context.Context, ref *provider.Reference) error {
	c, err := d.client(ctx)
	if err != nil {
		return err
	}
	parentID := ref.GetResourceId().GetOpaqueId()
	name := ref.GetPath()
	_, err = c.CreateFolder(parentID, name)
	return err
}

// TouchFile implements storage.FS — creates an empty file
func (d *Driver) TouchFile(ctx context.Context, ref *provider.Reference, _ bool, _ string) error {
	c, err := d.client(ctx)
	if err != nil {
		return err
	}
	parentID := ref.GetResourceId().GetOpaqueId()
	name := ref.GetPath()
	_, err = uploadFile(c, parentID, name, 0, io.NopCloser(nil), d.cfg.ChunkSize)
	return err
}

// Delete implements storage.FS
func (d *Driver) Delete(ctx context.Context, ref *provider.Reference) error {
	c, err := d.client(ctx)
	if err != nil {
		return err
	}
	id := ref.GetResourceId().GetOpaqueId()
	// Try folder first, fall back to file
	err = c.DeleteFolder(id)
	if err != nil {
		if ce, ok := err.(*ClientError); ok && ce.StatusCode == 404 {
			return c.DeleteFile(id)
		}
		return err
	}
	return nil
}

// Move implements storage.FS
func (d *Driver) Move(ctx context.Context, oldRef, newRef *provider.Reference) error {
	c, err := d.client(ctx)
	if err != nil {
		return err
	}
	sourceID := oldRef.GetResourceId().GetOpaqueId()
	destFolderID := newRef.GetResourceId().GetOpaqueId()
	// If same parent — this is a rename
	if destFolderID == "" {
		// rename: newRef has path only
		return c.RenameFile(sourceID, newRef.GetPath(), false)
	}
	return c.MoveResource(sourceID, destFolderID, false)
}

// InitiateUpload implements storage.FS
func (d *Driver) InitiateUpload(ctx context.Context, ref *provider.Reference, uploadLength int64, metadata map[string]string) (map[string]string, error) {
	c, err := d.client(ctx)
	if err != nil {
		return nil, err
	}
	parentID := ref.GetResourceId().GetOpaqueId()
	filename := metadata["filename"]
	if filename == "" {
		filename = ref.GetPath()
	}
	numChunks := 1
	if uploadLength > 0 {
		numChunks = int(math.Ceil(float64(uploadLength) / float64(d.cfg.ChunkSize)))
	}
	result, err := c.InitiateUpload(parentID, filename, uploadLength, numChunks)
	if err != nil {
		return nil, err
	}
	return map[string]string{
		"uploadID":  result.ID,
		"uploadURI": result.URI,
		"filename":  filename,
		"parentID":  parentID,
	}, nil
}

// Upload implements storage.FS
func (d *Driver) Upload(ctx context.Context, req storage.UploadRequest, uploadFunc storage.UploadFinishedFunc) (*provider.ResourceInfo, error) {
	c, err := d.client(ctx)
	if err != nil {
		return nil, err
	}
	parentID := req.Ref.GetResourceId().GetOpaqueId()
	filename := req.Ref.GetPath()
	fi, err := uploadFile(c, parentID, filename, -1, req.Body, d.cfg.ChunkSize)
	if err != nil {
		return nil, err
	}
	ri := fileInfoToResourceInfo(fi)
	if uploadFunc != nil {
		uploadFunc(ctx, ctxpkg.ContextMustGetUser(ctx), ri)
	}
	return ri, nil
}

// AddGrant implements storage.FS
func (d *Driver) AddGrant(ctx context.Context, ref *provider.Reference, g *provider.Grant) error {
	// Kiteworks permissions are managed via their own endpoint;
	// map the CS3 grant role to a Kiteworks permission ID at implementation time
	// once the OpenAPI spec permission IDs are confirmed.
	return errtypes.NotSupported("kiteworks: AddGrant — permission mapping not yet implemented")
}

// RemoveGrant implements storage.FS
func (d *Driver) RemoveGrant(ctx context.Context, ref *provider.Reference, g *provider.Grant) error {
	return errtypes.NotSupported("kiteworks: RemoveGrant — permission mapping not yet implemented")
}

// UpdateGrant implements storage.FS
func (d *Driver) UpdateGrant(ctx context.Context, ref *provider.Reference, g *provider.Grant) error {
	return errtypes.NotSupported("kiteworks: UpdateGrant — permission mapping not yet implemented")
}

// ListGrants implements storage.FS — returns Kiteworks permissions as CS3 grants
func (d *Driver) ListGrants(ctx context.Context, ref *provider.Reference) ([]*provider.Grant, error) {
	c, err := d.client(ctx)
	if err != nil {
		return nil, err
	}
	id := ref.GetResourceId().GetOpaqueId()
	fi, err := c.GetFolder(id)
	if err != nil {
		if ce, ok := err.(*ClientError); ok && ce.StatusCode == 404 {
			fi, err = c.GetFile(id)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	var grants []*provider.Grant
	for _, perm := range fi.Permissions {
		if !perm.Allowed {
			continue
		}
		grants = append(grants, &provider.Grant{
			Grantee: &provider.Grantee{
				Type: provider.GranteeType_GRANTEE_TYPE_USER,
				Id: &provider.Grantee_UserId{
					UserId: &typespb.UserId{OpaqueId: perm.Name},
				},
			},
			Permissions: &provider.ResourcePermissions{
				// Minimal: treat every allowed permission as read+write
				GetPath:              true,
				InitiateFileDownload: true,
				InitiateFileUpload:   true,
				ListContainer:        true,
				Stat:                 true,
			},
		})
	}
	return grants, nil
}
```

Also add `"math"` to the import block in `kiteworks.go` and `storage "github.com/owncloud/reva/v2/pkg/storage"`.

- [ ] **Step 2: Compile check**

```bash
go build ./vendor/github.com/owncloud/reva/v2/pkg/storage/fs/kiteworks/...
```
Expected: no output.

- [ ] **Step 3: Commit**

```bash
git add vendor/github.com/owncloud/reva/v2/pkg/storage/fs/kiteworks/kiteworks.go
git commit -s -m "feat(kiteworks): implement storage.FS write path, upload, grants"
```

---

## Task 7: Register driver in reva loader

**Files:**
- Modify: `vendor/github.com/owncloud/reva/v2/pkg/storage/fs/loader/loader.go`

- [ ] **Step 1: Add blank import for kiteworks driver**

In `loader.go`, add this line to the import block (after the `nextcloud` line):

```go
_ "github.com/owncloud/reva/v2/pkg/storage/fs/kiteworks"
```

- [ ] **Step 2: Compile check**

```bash
go build ./vendor/github.com/owncloud/reva/v2/pkg/storage/fs/loader/...
```
Expected: no output.

- [ ] **Step 3: Commit**

```bash
git add vendor/github.com/owncloud/reva/v2/pkg/storage/fs/loader/loader.go
git commit -s -m "feat(kiteworks): register kiteworks driver in reva fs loader"
```

---

## Task 8: oCIS service — config and defaults

**Files:**
- Create: `services/storage-kiteworks/pkg/config/config.go`
- Create: `services/storage-kiteworks/pkg/config/reva.go`
- Create: `services/storage-kiteworks/pkg/config/defaults/defaultconfig.go`
- Create: `services/storage-kiteworks/pkg/config/parser/parse.go`

- [ ] **Step 1: Create config.go**

```go
// services/storage-kiteworks/pkg/config/config.go
package config

import (
	"context"

	"github.com/owncloud/ocis/v2/ocis-pkg/shared"
)

// Config is the configuration for the storage-kiteworks service
type Config struct {
	Commons *shared.Commons `yaml:"-"`
	Service Service         `yaml:"-"`
	Tracing *Tracing        `yaml:"tracing"`
	Log     *Log            `yaml:"log"`
	Debug   Debug           `yaml:"debug"`

	GRPC GRPCConfig `yaml:"grpc"`

	TokenManager *TokenManager `yaml:"token_manager"`
	Reva         *shared.Reva  `yaml:"reva"`

	GracefulShutdownTimeout int `yaml:"graceful_shutdown_timeout" env:"STORAGE_KITEWORKS_GRACEFUL_SHUTDOWN_TIMEOUT" desc:"Seconds to wait for graceful shutdown." introductionVersion:"1.0.0"`

	Driver  KiteworksDriver `yaml:"driver"`
	MountID string          `yaml:"mount_id" env:"STORAGE_KITEWORKS_MOUNT_ID" desc:"Mount ID of this storage provider." introductionVersion:"1.0.0"`

	Context context.Context `yaml:"-"`
}

// Service holds service name config
type Service struct {
	Name string `yaml:"-" env:"STORAGE_KITEWORKS_SERVICE_NAME"`
}

// Log configures logging
type Log struct {
	Level  string `yaml:"level" env:"OCIS_LOG_LEVEL;STORAGE_KITEWORKS_LOG_LEVEL"`
	Pretty bool   `yaml:"pretty" env:"OCIS_LOG_PRETTY;STORAGE_KITEWORKS_LOG_PRETTY"`
	Color  bool   `yaml:"color" env:"OCIS_LOG_COLOR;STORAGE_KITEWORKS_LOG_COLOR"`
	File   string `yaml:"file" env:"OCIS_LOG_FILE;STORAGE_KITEWORKS_LOG_FILE"`
}

// Tracing configures distributed tracing
type Tracing struct {
	Enabled   bool   `yaml:"enabled" env:"OCIS_TRACING_ENABLED;STORAGE_KITEWORKS_TRACING_ENABLED"`
	Type      string `yaml:"type" env:"OCIS_TRACING_TYPE;STORAGE_KITEWORKS_TRACING_TYPE"`
	Endpoint  string `yaml:"endpoint" env:"OCIS_TRACING_ENDPOINT;STORAGE_KITEWORKS_TRACING_ENDPOINT"`
	Collector string `yaml:"collector" env:"OCIS_TRACING_COLLECTOR;STORAGE_KITEWORKS_TRACING_COLLECTOR"`
}

// Debug configures the debug/metrics server
type Debug struct {
	Addr   string `yaml:"addr" env:"STORAGE_KITEWORKS_DEBUG_ADDR"`
	Token  string `yaml:"token" env:"STORAGE_KITEWORKS_DEBUG_TOKEN"`
	Pprof  bool   `yaml:"pprof" env:"STORAGE_KITEWORKS_DEBUG_PPROF"`
	Zpages bool   `yaml:"zpages" env:"STORAGE_KITEWORKS_DEBUG_ZPAGES"`
}

// GRPCConfig configures the gRPC server
type GRPCConfig struct {
	Addr      string                 `yaml:"addr" env:"STORAGE_KITEWORKS_GRPC_ADDR"`
	TLS       *shared.GRPCServiceTLS `yaml:"tls"`
	Namespace string                 `yaml:"-"`
	Protocol  string                 `yaml:"protocol" env:"OCIS_GRPC_PROTOCOL;STORAGE_KITEWORKS_GRPC_PROTOCOL"`
}

// KiteworksDriver holds the Kiteworks-specific driver config
type KiteworksDriver struct {
	Endpoint  string `yaml:"endpoint"   env:"STORAGE_KITEWORKS_ENDPOINT"    desc:"Base URL of the Kiteworks server, e.g. https://kiteworks.example.com" introductionVersion:"1.0.0"`
	Insecure  bool   `yaml:"insecure"   env:"STORAGE_KITEWORKS_INSECURE"     desc:"Skip TLS certificate verification (development only)." introductionVersion:"1.0.0"`
	ChunkSize int64  `yaml:"chunk_size" env:"STORAGE_KITEWORKS_CHUNK_SIZE"   desc:"Upload chunk size in bytes. Default 5242880 (5 MB)." introductionVersion:"1.0.0"`
}
```

- [ ] **Step 2: Create reva.go**

```go
// services/storage-kiteworks/pkg/config/reva.go
package config

// TokenManager is the config for the reva token manager
type TokenManager struct {
	JWTSecret string `yaml:"jwt_secret" env:"OCIS_JWT_SECRET;STORAGE_KITEWORKS_JWT_SECRET" desc:"The secret to mint and validate JWT tokens." introductionVersion:"1.0.0"`
}
```

- [ ] **Step 3: Create defaults/defaultconfig.go**

```go
// services/storage-kiteworks/pkg/config/defaults/defaultconfig.go
package defaults

import (
	"github.com/owncloud/ocis/v2/ocis-pkg/shared"
	"github.com/owncloud/ocis/v2/services/storage-kiteworks/pkg/config"
)

// FullDefaultConfig returns a fully initialized default configuration
func FullDefaultConfig() *config.Config {
	cfg := DefaultConfig()
	EnsureDefaults(cfg)
	Sanitize(cfg)
	return cfg
}

// DefaultConfig returns a basic default configuration
func DefaultConfig() *config.Config {
	return &config.Config{
		Debug: config.Debug{
			Addr:   "127.0.0.1:9179",
			Token:  "",
			Pprof:  false,
			Zpages: false,
		},
		GRPC: config.GRPCConfig{
			Addr:      "127.0.0.1:9177",
			Namespace: "com.owncloud.api",
			Protocol:  "tcp",
		},
		Service: config.Service{
			Name: "storage-kiteworks",
		},
		Reva:                    shared.DefaultRevaConfig(),
		GracefulShutdownTimeout: 30,
		Driver: config.KiteworksDriver{
			ChunkSize: 5 * 1024 * 1024,
		},
	}
}

// EnsureDefaults sets values that depend on other config fields
func EnsureDefaults(cfg *config.Config) {
	if cfg.Reva == nil {
		cfg.Reva = shared.DefaultRevaConfig()
	}
	if cfg.TokenManager == nil {
		cfg.TokenManager = &config.TokenManager{}
	}
	if cfg.Log == nil {
		cfg.Log = &config.Log{}
	}
	if cfg.Tracing == nil {
		cfg.Tracing = &config.Tracing{}
	}
}

// Sanitize cleans up any config values
func Sanitize(_ *config.Config) {}
```

- [ ] **Step 4: Create parser/parse.go**

```go
// services/storage-kiteworks/pkg/config/parser/parse.go
package parser

import (
	"errors"
	"fmt"

	ociscfg "github.com/owncloud/ocis/v2/ocis-pkg/config"
	defaults2 "github.com/owncloud/ocis/v2/ocis-pkg/config/defaults"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/envdecode"
	"github.com/owncloud/ocis/v2/ocis-pkg/shared"
	"github.com/owncloud/ocis/v2/services/storage-kiteworks/pkg/config"
	"github.com/owncloud/ocis/v2/services/storage-kiteworks/pkg/config/defaults"
)

// ParseConfig loads configuration from known paths
func ParseConfig(cfg *config.Config) error {
	err := ociscfg.BindSourcesToStructs(cfg.Service.Name, cfg)
	if err != nil {
		return err
	}

	defaults.EnsureDefaults(cfg)

	if err := envdecode.Decode(cfg); err != nil {
		if !errors.Is(err, envdecode.ErrNoTargetFieldsAreSet) {
			return err
		}
	}

	defaults.Sanitize(cfg)

	return Validate(cfg)
}

// Validate checks that required config fields are set
func Validate(cfg *config.Config) error {
	if cfg.TokenManager.JWTSecret == "" {
		return shared.MissingJWTTokenError(cfg.Service.Name)
	}
	if cfg.MountID == "" {
		return fmt.Errorf("The storage-kiteworks mount ID has not been configured. "+
			"Run ocis init or set STORAGE_KITEWORKS_MOUNT_ID. See %s for details.",
			defaults2.BaseConfigPath())
	}
	if cfg.Driver.Endpoint == "" {
		return fmt.Errorf("STORAGE_KITEWORKS_ENDPOINT must be set to the base URL of the Kiteworks server")
	}
	return nil
}
```

- [ ] **Step 5: Compile check**

```bash
go build ./services/storage-kiteworks/...
```
Expected: no output (some files still missing — that's OK at this stage, just the config packages must compile).

- [ ] **Step 6: Commit**

```bash
git add services/storage-kiteworks/pkg/config/
git commit -s -m "feat(kiteworks): add oCIS service config and defaults"
```

---

## Task 9: oCIS service — logging, debug server, revaconfig

**Files:**
- Create: `services/storage-kiteworks/pkg/logging/logging.go`
- Create: `services/storage-kiteworks/pkg/server/debug/server.go`
- Create: `services/storage-kiteworks/pkg/server/debug/options.go`
- Create: `services/storage-kiteworks/pkg/revaconfig/config.go`

- [ ] **Step 1: Create logging/logging.go**

```go
// services/storage-kiteworks/pkg/logging/logging.go
package logging

import (
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/services/storage-kiteworks/pkg/config"
)

// Configure initializes a service-specific logger instance
func Configure(name string, cfg *config.Log) log.Logger {
	return log.NewLogger(
		log.Name(name),
		log.Level(cfg.Level),
		log.Pretty(cfg.Pretty),
		log.Color(cfg.Color),
		log.File(cfg.File),
	)
}
```

- [ ] **Step 2: Create debug server options**

```go
// services/storage-kiteworks/pkg/server/debug/options.go
package debug

import (
	"context"

	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/services/storage-kiteworks/pkg/config"
)

// Option is a function to configure debug server options
type Option func(*Options)

// Options holds debug server configuration
type Options struct {
	Logger  log.Logger
	Context context.Context
	Config  *config.Config
}

func newOptions(opts ...Option) Options {
	o := Options{}
	for _, opt := range opts {
		opt(&o)
	}
	return o
}

// Logger sets the logger
func Logger(l log.Logger) Option {
	return func(o *Options) { o.Logger = l }
}

// Context sets the context
func Context(ctx context.Context) Option {
	return func(o *Options) { o.Context = ctx }
}

// Config sets the config
func Config(cfg *config.Config) Option {
	return func(o *Options) { o.Config = cfg }
}
```

- [ ] **Step 3: Create debug/server.go**

```go
// services/storage-kiteworks/pkg/server/debug/server.go
package debug

import (
	"net/http"

	"github.com/owncloud/ocis/v2/ocis-pkg/checks"
	"github.com/owncloud/ocis/v2/ocis-pkg/handlers"
	"github.com/owncloud/ocis/v2/ocis-pkg/service/debug"
	"github.com/owncloud/ocis/v2/ocis-pkg/version"
)

// Server initializes the debug service and server
func Server(opts ...Option) (*http.Server, error) {
	options := newOptions(opts...)

	readyHandler := handlers.NewCheckHandler(
		handlers.NewCheckHandlerConfiguration().
			WithLogger(options.Logger).
			WithCheck("grpc reachability", checks.NewGRPCCheck(options.Config.GRPC.Addr)),
	)

	return debug.NewService(
		debug.Logger(options.Logger),
		debug.Context(options.Context),
		debug.Name(options.Config.Service.Name),
		debug.Version(version.GetString()),
		debug.Address(options.Config.Debug.Addr),
		debug.Token(options.Config.Debug.Token),
		debug.Pprof(options.Config.Debug.Pprof),
		debug.Zpages(options.Config.Debug.Zpages),
		debug.Ready(readyHandler),
	), nil
}
```

- [ ] **Step 4: Create revaconfig/config.go**

```go
// services/storage-kiteworks/pkg/revaconfig/config.go
package revaconfig

import (
	"github.com/owncloud/ocis/v2/services/storage-kiteworks/pkg/config"
)

// StorageKiteworksConfigFromStruct adapts an oCIS Config into a reva mapstructure
func StorageKiteworksConfigFromStruct(cfg *config.Config) map[string]interface{} {
	return map[string]interface{}{
		"core": map[string]interface{}{
			"graceful_shutdown_timeout": cfg.GracefulShutdownTimeout,
		},
		"shared": map[string]interface{}{
			"jwt_secret":  cfg.TokenManager.JWTSecret,
			"gatewaysvc":  cfg.Reva.Address,
			"grpc_client_options": cfg.Reva.GetGRPCClientConfig(),
		},
		"grpc": map[string]interface{}{
			"network": cfg.GRPC.Protocol,
			"address": cfg.GRPC.Addr,
			"tls_settings": map[string]interface{}{
				"enabled":     cfg.GRPC.TLS.Enabled,
				"certificate": cfg.GRPC.TLS.Cert,
				"key":         cfg.GRPC.TLS.Key,
			},
			"services": map[string]interface{}{
				"storageprovider": map[string]interface{}{
					"driver":  "kiteworks",
					"drivers": map[string]interface{}{
						"kiteworks": map[string]interface{}{
							"endpoint":   cfg.Driver.Endpoint,
							"insecure":   cfg.Driver.Insecure,
							"chunk_size": cfg.Driver.ChunkSize,
						},
					},
					"mount_id": cfg.MountID,
				},
			},
		},
	}
}
```

- [ ] **Step 5: Compile check**

```bash
go build ./services/storage-kiteworks/...
```
Expected: packages compile (binary entry point still missing).

- [ ] **Step 6: Commit**

```bash
git add services/storage-kiteworks/pkg/logging/ services/storage-kiteworks/pkg/server/ services/storage-kiteworks/pkg/revaconfig/
git commit -s -m "feat(kiteworks): add logging, debug server, reva config mapping"
```

---

## Task 10: oCIS service — commands and binary entry point

**Files:**
- Create: `services/storage-kiteworks/pkg/command/root.go`
- Create: `services/storage-kiteworks/pkg/command/server.go`
- Create: `services/storage-kiteworks/pkg/command/health.go`
- Create: `services/storage-kiteworks/pkg/command/version.go`
- Create: `services/storage-kiteworks/cmd/storage-kiteworks/main.go`

- [ ] **Step 1: Create command/root.go**

```go
// services/storage-kiteworks/pkg/command/root.go
package command

import (
	"os"

	"github.com/owncloud/ocis/v2/ocis-pkg/clihelper"
	"github.com/owncloud/ocis/v2/services/storage-kiteworks/pkg/config"
	"github.com/urfave/cli/v2"
)

// GetCommands provides all commands for this service
func GetCommands(cfg *config.Config) cli.Commands {
	return []*cli.Command{
		Server(cfg),
		Health(cfg),
		Version(cfg),
	}
}

// Execute is the entry point for the storage-kiteworks command
func Execute(cfg *config.Config) error {
	app := clihelper.DefaultApp(&cli.App{
		Name:     "storage-kiteworks",
		Usage:    "Provide Kiteworks storage integration for oCIS",
		Commands: GetCommands(cfg),
	})
	return app.RunContext(cfg.Context, os.Args)
}
```

- [ ] **Step 2: Create command/server.go**

```go
// services/storage-kiteworks/pkg/command/server.go
package command

import (
	"context"
	"fmt"
	"os/signal"

	"github.com/owncloud/ocis/v2/ocis-pkg/config/configlog"
	"github.com/owncloud/ocis/v2/ocis-pkg/registry"
	"github.com/owncloud/ocis/v2/ocis-pkg/runner"
	"github.com/owncloud/ocis/v2/ocis-pkg/tracing"
	"github.com/owncloud/ocis/v2/ocis-pkg/version"
	"github.com/owncloud/ocis/v2/services/storage-kiteworks/pkg/config"
	"github.com/owncloud/ocis/v2/services/storage-kiteworks/pkg/config/parser"
	"github.com/owncloud/ocis/v2/services/storage-kiteworks/pkg/logging"
	"github.com/owncloud/ocis/v2/services/storage-kiteworks/pkg/revaconfig"
	"github.com/owncloud/ocis/v2/services/storage-kiteworks/pkg/server/debug"
	"github.com/owncloud/reva/v2/cmd/revad/runtime"
	"github.com/urfave/cli/v2"
)

// Server is the entry point for the server command
func Server(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "server",
		Usage:    fmt.Sprintf("start the %s service without runtime (unsupervised mode)", cfg.Service.Name),
		Category: "server",
		Before: func(c *cli.Context) error {
			return configlog.ReturnFatal(parser.ParseConfig(cfg))
		},
		Action: func(c *cli.Context) error {
			logger := logging.Configure(cfg.Service.Name, cfg.Log)
			traceProvider, err := tracing.GetServiceTraceProvider(cfg.Tracing, cfg.Service.Name)
			if err != nil {
				return err
			}

			var cancel context.CancelFunc
			if cfg.Context == nil {
				cfg.Context, cancel = signal.NotifyContext(context.Background(), runner.StopSignals...)
				defer cancel()
			}
			ctx := cfg.Context

			gr := runner.NewGroup()

			rCfg := revaconfig.StorageKiteworksConfigFromStruct(cfg)
			if rServer := runtime.NewDrivenGRPCServerWithOptions(rCfg,
				runtime.WithLogger(&logger.Logger),
				runtime.WithRegistry(registry.GetRegistry()),
				runtime.WithTraceProvider(traceProvider),
			); rServer != nil {
				gr.Add(runner.NewRevaServiceRunner(cfg.Service.Name+".rgrpc", rServer))
			}

			debugServer, err := debug.Server(
				debug.Logger(logger),
				debug.Context(ctx),
				debug.Config(cfg),
			)
			if err != nil {
				logger.Info().Err(err).Str("server", "debug").Msg("Failed to initialize server")
				return err
			}
			gr.Add(runner.NewGolangHttpServerRunner("storage-kiteworks_debug", debugServer))

			grpcSvc := registry.BuildGRPCService(cfg.GRPC.Namespace+"."+cfg.Service.Name, cfg.GRPC.Protocol, cfg.GRPC.Addr, version.GetString())
			if err := registry.RegisterService(ctx, logger, grpcSvc, cfg.Debug.Addr); err != nil {
				logger.Fatal().Err(err).Msg("failed to register the grpc service")
			}

			logger.Info().Msgf("starting service %s", cfg.Service.Name)
			grResults := gr.Run(ctx)
			if err := runner.ProcessResults(grResults); err != nil {
				logger.Error().Err(err).Msgf("service %s stopped with error", cfg.Service.Name)
				return err
			}
			logger.Info().Msgf("service %s stopped without error", cfg.Service.Name)
			return nil
		},
	}
}
```

- [ ] **Step 3: Create command/health.go**

```go
// services/storage-kiteworks/pkg/command/health.go
package command

import (
	"fmt"
	"net/http"

	"github.com/owncloud/ocis/v2/ocis-pkg/config/configlog"
	"github.com/owncloud/ocis/v2/services/storage-kiteworks/pkg/config"
	"github.com/owncloud/ocis/v2/services/storage-kiteworks/pkg/config/parser"
	"github.com/owncloud/ocis/v2/services/storage-kiteworks/pkg/logging"
	"github.com/urfave/cli/v2"
)

// Health is the entrypoint for the health command
func Health(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "health",
		Usage:    "check health status",
		Category: "info",
		Before: func(c *cli.Context) error {
			return configlog.ReturnError(parser.ParseConfig(cfg))
		},
		Action: func(c *cli.Context) error {
			logger := logging.Configure(cfg.Service.Name, cfg.Log)
			resp, err := http.Get(fmt.Sprintf("http://%s/healthz", cfg.Debug.Addr))
			if err != nil {
				logger.Fatal().Err(err).Msg("Failed to request health check")
			}
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				logger.Fatal().Int("code", resp.StatusCode).Msg("Health seems to be in bad state")
			}
			logger.Debug().Int("code", resp.StatusCode).Msg("Health got a good state")
			return nil
		},
	}
}
```

- [ ] **Step 4: Create command/version.go**

```go
// services/storage-kiteworks/pkg/command/version.go
package command

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/owncloud/ocis/v2/ocis-pkg/registry"
	"github.com/owncloud/ocis/v2/ocis-pkg/version"
	"github.com/owncloud/ocis/v2/services/storage-kiteworks/pkg/config"
	"github.com/urfave/cli/v2"
)

// Version prints the service versions of all running instances
func Version(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "version",
		Usage:    "print the version of this binary and the running service instances",
		Category: "info",
		Action: func(c *cli.Context) error {
			fmt.Println("Version: " + version.GetString())
			fmt.Printf("Compiled: %s\n", version.Compiled())
			fmt.Println("")

			reg := registry.GetRegistry()
			services, err := reg.GetService(cfg.GRPC.Namespace + "." + cfg.Service.Name)
			if err != nil {
				fmt.Println(fmt.Errorf("could not get %s services from the registry: %v", cfg.Service.Name, err))
				return err
			}
			if len(services) == 0 {
				fmt.Println("No running " + cfg.Service.Name + " service found.")
				return nil
			}
			table := tablewriter.NewTable(os.Stdout)
			table.Header("Version", "Address", "Id")
			for _, s := range services {
				for _, n := range s.Nodes {
					table.Append([]string{s.Version, n.Address, n.Id})
				}
			}
			table.Render()
			return nil
		},
	}
}
```

- [ ] **Step 5: Create cmd/storage-kiteworks/main.go**

```go
// services/storage-kiteworks/cmd/storage-kiteworks/main.go
package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/owncloud/ocis/v2/services/storage-kiteworks/pkg/command"
	"github.com/owncloud/ocis/v2/services/storage-kiteworks/pkg/config/defaults"
)

func main() {
	cfg := defaults.DefaultConfig()
	cfg.Context, _ = signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)
	if err := command.Execute(cfg); err != nil {
		os.Exit(1)
	}
}
```

- [ ] **Step 6: Full service compile check**

```bash
go build ./services/storage-kiteworks/...
```
Expected: no output.

- [ ] **Step 7: Commit**

```bash
git add services/storage-kiteworks/
git commit -s -m "feat(kiteworks): add oCIS service commands and binary entry point"
```

---

## Task 11: Wire service into main oCIS binary

**Files:**
- Modify: `ocis-pkg/config/config.go`
- Modify: `ocis-pkg/config/defaults/defaultconfig.go` (or wherever global defaults are set)
- Modify: `ocis/pkg/command/services.go`
- Modify: `ocis/pkg/runtime/service/service.go`

- [ ] **Step 1: Add StorageKiteworks to global config**

In `ocis-pkg/config/config.go`, add to the import block:
```go
storagekiteworks "github.com/owncloud/ocis/v2/services/storage-kiteworks/pkg/config"
```

In the `Config` struct (near line 116–119 where other storage services are listed), add:
```go
StorageKiteworks *storagekiteworks.Config `yaml:"storage_kiteworks"`
```

- [ ] **Step 2: Add default in global defaultconfig**

Find `ocis-pkg/config/defaults/defaultconfig.go` (or the equivalent file that initializes all service configs) and add:
```go
// in the return/init block where StorageUsers is initialized:
cfg.StorageKiteworks = skdefaults.FullDefaultConfig()
```
Also add the import:
```go
skdefaults "github.com/owncloud/ocis/v2/services/storage-kiteworks/pkg/config/defaults"
```

- [ ] **Step 3: Register CLI command in services.go**

In `ocis/pkg/command/services.go`, add to the import block:
```go
storagekiteworks "github.com/owncloud/ocis/v2/services/storage-kiteworks/pkg/command"
```

Add to the `GetCommands` or equivalent function (near the storage-users entry):
```go
ServiceCommand(cfg, cfg.StorageKiteworks.Service.Name, storagekiteworks.GetCommands(cfg.StorageKiteworks), func(c *config.Config) {
    cfg.StorageKiteworks.Commons = cfg.Commons
}),
```

- [ ] **Step 4: Register runtime service in service.go**

In `ocis/pkg/runtime/service/service.go`, add to the import block:
```go
storagekiteworks "github.com/owncloud/ocis/v2/services/storage-kiteworks/pkg/command"
```

Add to the `RegisterServices` function (near the storage-users entry — use priority level 3):
```go
reg(3, opts.Config.StorageKiteworks.Service.Name, func(ctx context.Context, cfg *ociscfg.Config) error {
    cfg.StorageKiteworks.Context = ctx
    cfg.StorageKiteworks.Commons = cfg.Commons
    return runServerCommand(ctx, storagekiteworks.Server(cfg.StorageKiteworks))
})
```

- [ ] **Step 5: Full binary compile check**

```bash
go build ./ocis/...
```
Expected: no output.

- [ ] **Step 6: Commit**

```bash
git add ocis-pkg/config/ ocis/pkg/command/services.go ocis/pkg/runtime/service/service.go
git commit -s -m "feat(kiteworks): wire storage-kiteworks service into main oCIS binary"
```

---

## Task 12: Run all tests

- [ ] **Step 1: Run driver unit tests**

```bash
go test ./vendor/github.com/owncloud/reva/v2/pkg/storage/fs/kiteworks/... -v
```
Expected: All Ginkgo specs PASS.

- [ ] **Step 2: Run full oCIS test suite (kiteworks packages)**

```bash
go test ./services/storage-kiteworks/... -v
```
Expected: PASS (config/parser tests if any, no failures).

- [ ] **Step 3: Lint check**

```bash
make golangci-lint
```
Expected: no new lint errors in `services/storage-kiteworks/` or `vendor/.../kiteworks/`.

- [ ] **Step 4: Final commit**

```bash
git add .
git commit -s -m "test(kiteworks): verify all tests pass"
```

---

## Self-Review

**Spec coverage:**
- ✅ Section 1 (service structure): Tasks 8–11
- ✅ Section 2 (CS3 → Kiteworks mapping): Tasks 5–6 (read + write path)
- ✅ Space type mapping (`IsSharedWithUser` → mountpoint/project): Task 5 `spaceFromFileInfo`
- ✅ Section 3 (config + auth): Tasks 8–9, token passthrough in `client()` method
- ✅ Section 4 (upload flow): Tasks 4 + 6
- ✅ Section 5 (integration + testing): Tasks 7, 11, 12
- ✅ Grants (ListGrants implemented, Add/Remove/Update stubbed pending permission ID mapping): Task 6
- ✅ GetMD folder-first-then-file fallback: Task 5

**Type consistency check:**
- `Client` defined in Task 2, used in Tasks 5–6 ✅
- `uploadFile` defined in Task 4, called in Task 6 ✅
- `fileInfoToResourceInfo` defined in Task 5, used in Tasks 5–6 ✅
- `spaceFromFileInfo` defined in Task 5, used in Task 5 ✅
- `Config` struct in Task 8, used in Tasks 9–10 ✅
- `StorageKiteworksConfigFromStruct` in Task 9, called in Task 10 ✅

**No placeholders:** All code is complete. Grant permission mapping is explicitly stubbed with a clear comment noting it requires confirming Kiteworks permission IDs from the OpenAPI spec.
