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
	"net/url"
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
	err := c.doJSON(http.MethodPost, "/rest/files/"+sourceID+"/actions/move",
		&FileCopyMove{DestFolderID: destFolderID, Replace: replace}, nil)
	if err != nil {
		if ce, ok := err.(*ClientError); ok && ce.StatusCode == 404 {
			return c.doJSON(http.MethodPost, "/rest/folders/"+sourceID+"/actions/move",
				&FileCopyMove{DestFolderID: destFolderID, Replace: replace}, nil)
		}
		return err
	}
	return nil
}

// CopyResource copies a file or folder to a destination folder
func (c *Client) CopyResource(sourceID, destFolderID string, replace bool) error {
	err := c.doJSON(http.MethodPost, "/rest/files/"+sourceID+"/actions/copy",
		&FileCopyMove{DestFolderID: destFolderID, Replace: replace}, nil)
	if err != nil {
		if ce, ok := err.(*ClientError); ok && ce.StatusCode == 404 {
			return c.doJSON(http.MethodPost, "/rest/folders/"+sourceID+"/actions/copy",
				&FileCopyMove{DestFolderID: destFolderID, Replace: replace}, nil)
		}
		return err
	}
	return nil
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
	err := c.doJSON(http.MethodGet, "/rest/query?query="+url.QueryEscape(path), nil, &result)
	if err != nil {
		return nil, err
	}
	items := append(result.Folders, result.Files...)
	return items, nil
}
