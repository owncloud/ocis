package kwlib

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

const uploadChunkSize = 8 << 20 // 8 MiB

func NewClientFactory(server, agentString string, insecure bool) *APIClientFactory {
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   32,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	if insecure {
		// #nosec
		transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	uploadTransport := transport.Clone()
	uploadTransport.ResponseHeaderTimeout = 30 * time.Second

	return &APIClientFactory{
		server:       server,
		agentString:  agentString,
		httpClient:   &http.Client{Transport: transport, Timeout: 15 * time.Second},
		uploadClient: &http.Client{Transport: uploadTransport},
	}
}

type APIClientFactory struct {
	server       string
	agentString  string
	httpClient   *http.Client
	uploadClient *http.Client
}

type APIClient struct {
	server       string
	agentString  string
	logger       *zerolog.Logger
	host         string
	token        string
	requestId    string
	remoteAddr   string
	httpClient   *http.Client
	uploadClient *http.Client
}

func decodeJSON(body io.ReadCloser, out any) error {
	defer body.Close()
	return json.NewDecoder(body).Decode(out)
}

func (f *APIClientFactory) Build(host, requestId, remoteAddr, token string, l *zerolog.Logger) *APIClient {
	return &APIClient{
		token:        token,
		server:       f.server,
		host:         host,
		logger:       l,
		requestId:    requestId,
		remoteAddr:   remoteAddr,
		httpClient:   f.httpClient,
		uploadClient: f.uploadClient,
	}
}

func (c *APIClient) GetTopFolders() (*DirectoryInfo, error) {
	request, err := c.NewGetRequest("/rest/folders/top?deleted=false&with=(permissions)")
	if err != nil {
		return nil, err
	}
	response, err := c.SendRequest(request)
	if err != nil {
		return nil, err
	}
	out := &DirectoryInfo{}
	if err := decodeJSON(response.Body, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *APIClient) GetFolderByID(id string) (*FileInfo, error) {
	request, err := c.NewGetRequest(fmt.Sprintf("/rest/folders/%s?with=(permissions)", id))
	if err != nil {
		return nil, err
	}
	response, err := c.SendRequest(request)
	if err != nil {
		return nil, err
	}
	out := &FileInfo{}
	if err := decodeJSON(response.Body, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *APIClient) ListFolderContents(id string) ([]FileInfo, error) {
	request, err := c.NewGetRequest(fmt.Sprintf("/rest/folders/%s/children?deleted=false&with=(permissions)", id))
	if err != nil {
		return nil, err
	}
	response, err := c.SendRequest(request)
	if err != nil {
		return nil, err
	}
	dir := &DirectoryInfo{}
	if err := decodeJSON(response.Body, dir); err != nil {
		return nil, err
	}
	return dir.Data, nil
}

func (c *APIClient) Search(path string) (*FileInfo, error) {
	req, err := c.NewGetRequest("/rest/query")
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Add("includeContent", "false")
	q.Add("searchType", "f,d")
	q.Add("path", path)
	q.Add("with", "(permissions)")
	req.URL.RawQuery = q.Encode()
	response, err := c.SendRequest(req)
	if err != nil {
		return nil, err
	}
	search := &FileSearch{}
	if err := decodeJSON(response.Body, search); err != nil {
		return nil, err
	}
	return search.FindByParent(""), nil
}

func (c *APIClient) GetFileByID(id string) (*FileInfo, error) {
	request, err := c.NewGetRequest(fmt.Sprintf("/rest/files/%s?with=(permissions,lockUser)", id))
	if err != nil {
		return nil, err
	}
	response, err := c.SendRequest(request)
	if err != nil {
		return nil, err
	}
	out := &FileInfo{}
	if err := decodeJSON(response.Body, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *APIClient) GetFileContents(id, httpRange string) (*http.Response, error) {
	request, err := c.NewGetRequest(fmt.Sprintf("/rest/files/%s/content", id))
	if err != nil {
		return nil, err
	}
	if httpRange != "" {
		request.Header.Set("Range", httpRange)
	}
	return c.SendRequest(request)
}

func (c *APIClient) GetMe() (*User, error) {
	return c.GetUser("me")
}

func (c *APIClient) GetUser(id string) (*User, error) {
	request, err := c.NewGetRequest(fmt.Sprintf("/rest/users/%s", id))
	if err != nil {
		return nil, err
	}
	response, err := c.SendRequest(request)
	if err != nil {
		return nil, err
	}
	out := &User{}
	if err := decodeJSON(response.Body, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *APIClient) GetFolderQuota(folderID string) (*FolderQuota, error) {
	request, err := c.NewGetRequest(fmt.Sprintf("/rest/folders/%s/quota", folderID))
	if err != nil {
		return nil, err
	}
	response, err := c.SendRequest(request)
	if err != nil {
		return nil, err
	}
	out := &FolderQuota{}
	if err := decodeJSON(response.Body, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *APIClient) GetGroups(limit, offset int) (*ContactList, error) {
	q := url.Values{}
	if limit > 0 {
		q.Set("limit", strconv.Itoa(limit))
	}
	if offset > 0 {
		q.Set("offset", strconv.Itoa(offset))
	}
	path := "/rest/groups"
	if len(q) > 0 {
		path += "?" + q.Encode()
	}
	request, err := c.NewGetRequest(path)
	if err != nil {
		return nil, err
	}
	response, err := c.SendRequest(request)
	if err != nil {
		return nil, err
	}
	out := &ContactList{}
	if err := decodeJSON(response.Body, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *APIClient) CreateFolder(id string, payload CreateDirRequest) (string, error) {
	request, err := c.NewPostRequest(fmt.Sprintf("/rest/folders/%s/folders", id), payload)
	if err != nil {
		return "", err
	}
	response, err := c.SendRequest(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	return ParseFolderIDFromPath(response.Header.Get("X-Accellion-Location")), nil
}

func (c *APIClient) InitializeUpload(parentID, name string, size int64, numberOfChunks int) (*UploadResult, error) {
	return c.initializeUpload(fmt.Sprintf("/rest/folders/%s/actions/initiateUpload", parentID), name, size, numberOfChunks)
}

func (c *APIClient) InitializeVersionUpload(fileID, name string, size int64, numberOfChunks int) (*UploadResult, error) {
	return c.initializeUpload(fmt.Sprintf("/rest/files/%s/actions/initiateUpload", fileID), name, size, numberOfChunks)
}

func (c *APIClient) initializeUpload(path, name string, size int64, numberOfChunks int) (*UploadResult, error) {
	payload := InitializeUpload{
		FileName:    name,
		TotalSize:   size,
		TotalChunks: numberOfChunks,
	}
	request, err := c.NewPostRequest(path, payload)
	if err != nil {
		return nil, err
	}
	response, err := c.SendRequest(request)
	if err != nil {
		return nil, err
	}
	out := &UploadResult{}
	if err := decodeJSON(response.Body, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *APIClient) TerminateUpload(uploadID int64) error {
	request, err := c.newRequest("DELETE", fmt.Sprintf("/rest/uploads/%d", uploadID), nil)
	if err != nil {
		return err
	}
	response, err := c.SendRequest(request)
	if err != nil {
		return err
	}
	response.Body.Close()
	return nil
}

func (c *APIClient) UploadChunk(ctx context.Context, uploadURI, name string, file io.Reader, chunkIndex int, chunk int64, isLastChunk bool) (*FileInfo, error) {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("content", name)
	if err != nil {
		return nil, err
	}
	if _, err = io.Copy(part, io.LimitReader(file, chunk)); err != nil {
		return nil, err
	}
	for k, v := range map[string]string{
		"compressionMode": "NORMAL",
		"compressionSize": strconv.FormatInt(chunk, 10),
		"originalSize":    strconv.FormatInt(chunk, 10),
		"index":           strconv.Itoa(chunkIndex + 1),
	} {
		_ = writer.WriteField(k, v)
	}
	if err = writer.Close(); err != nil {
		return nil, err
	}
	request, err := c.newRequest("POST", uploadURI, body)
	if err != nil {
		return nil, err
	}
	request = request.WithContext(ctx)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	if isLastChunk {
		q := request.URL.Query()
		q.Add("returnEntity", "true")
		request.URL.RawQuery = q.Encode()
	}
	response, err := c.sendWith(c.uploadClient, request)
	if err != nil {
		return nil, err
	}
	if isLastChunk {
		out := &FileInfo{}
		if err := decodeJSON(response.Body, out); err != nil {
			return nil, err
		}
		return out, nil
	}
	// drained so the connection can be reused across chunks
	_, _ = io.Copy(io.Discard, response.Body)
	response.Body.Close()
	return nil, nil
}

func (c *APIClient) MoveFolder(id, destinationFolderID string) error {
	request, err := c.NewPostRequest(
		fmt.Sprintf("/rest/folders/%s/actions/move", id),
		MoveFolderRequest{DestinationFolderID: destinationFolderID},
	)
	if err != nil {
		return err
	}
	_, err = c.SendRequest(request)
	return err
}

// name must be the file's current name: the server requires it even for a version upload.
func (c *APIClient) UploadFileVersion(ctx context.Context, fileID, name string, body io.Reader, length int64) error {
	chunks := chunkCount(length)
	session, err := c.InitializeVersionUpload(fileID, name, length, chunks)
	if err != nil {
		return err
	}
	for i := 0; i < chunks; i++ {
		size := int64(uploadChunkSize)
		if remaining := length - int64(i)*uploadChunkSize; remaining < size {
			size = remaining
		}
		if _, err := c.UploadChunk(ctx, session.URI, name, body, i, size, i == chunks-1); err != nil {
			if termErr := c.TerminateUpload(session.ID); termErr != nil {
				c.logger.Warn().Err(termErr).Int64("uploadID", session.ID).Msg("could not terminate kiteworks upload session")
			}
			return err
		}
	}
	return nil
}

func chunkCount(length int64) int {
	if length <= 0 {
		return 1
	}
	return int((length + uploadChunkSize - 1) / uploadChunkSize)
}

func (c *APIClient) GetFileVersions(fileID string) ([]Version, error) {
	req, err := c.NewGetRequest(fmt.Sprintf("/rest/files/%s/versions", fileID))
	if err != nil {
		return nil, err
	}
	resp, err := c.SendRequest(req)
	if err != nil {
		return nil, err
	}
	out := &VersionList{}
	if err := decodeJSON(resp.Body, out); err != nil {
		return nil, err
	}
	return out.Data, nil
}

func (c *APIClient) DeleteFileVersion(fileID, versionID string) error {
	req, err := c.newRequest("DELETE", fmt.Sprintf("/rest/files/%s/versions/%s", fileID, versionID), nil)
	if err != nil {
		return err
	}
	_, err = c.SendRequest(req)
	return err
}

func (c *APIClient) PromoteFileVersion(fileID, versionID string) error {
	req, err := c.newRequest("POST", fmt.Sprintf("/rest/files/%s/versions/%s/actions/promote", fileID, versionID), nil)
	if err != nil {
		return err
	}
	_, err = c.SendRequest(req)
	return err
}

func (c *APIClient) GetVersionContents(fileID, versionID string) (*http.Response, error) {
	req, err := c.NewGetRequest(fmt.Sprintf("/rest/files/%s/versions/%s/content", fileID, versionID))
	if err != nil {
		return nil, err
	}
	return c.SendRequest(req)
}

func (c *APIClient) LockFile(fileID string) error {
	req, err := c.newRequest("PATCH", fmt.Sprintf("/rest/files/%s/actions/lock", fileID), nil)
	if err != nil {
		return err
	}
	resp, err := c.SendRequest(req)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

func (c *APIClient) UnlockFile(fileID string) error {
	req, err := c.newRequest("PATCH", fmt.Sprintf("/rest/files/%s/actions/unlock", fileID), nil)
	if err != nil {
		return err
	}
	resp, err := c.SendRequest(req)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

func (c *APIClient) DeleteFolder(id string) error {
	request, err := c.newRequest("DELETE", fmt.Sprintf("/rest/folders/%s", id), nil)
	if err != nil {
		return err
	}
	_, err = c.SendRequest(request)
	return err
}

func (c *APIClient) DeleteFile(id string) error {
	request, err := c.newRequest("DELETE", fmt.Sprintf("/rest/files/%s", id), nil)
	if err != nil {
		return err
	}
	_, err = c.SendRequest(request)
	return err
}

func (c *APIClient) Move(source *FileInfo, parent *FileInfo, replace bool) (bool, error) {
	return c.moveOrCopy(moveOp, source, parent, replace)
}

func (c *APIClient) Copy(source *FileInfo, parent *FileInfo, replace bool) (bool, error) {
	return c.moveOrCopy(copyOp, source, parent, replace)
}

type moveOrCopy int

const (
	moveOp moveOrCopy = iota + 1
	copyOp
)

func (c *APIClient) moveOrCopy(op moveOrCopy, source *FileInfo, dest *FileInfo, replace bool) (bool, error) {
	api := "/rest/files/actions/move"
	if op == copyOp {
		api = "/rest/files/actions/copy"
	}
	request, err := c.NewPostRequest(api, FileCopyMove{DestinationFolderID: dest.ID, Replace: replace})
	if err != nil {
		return false, err
	}
	q := request.URL.Query()
	q.Add("id:in", source.ID)
	request.URL.RawQuery = q.Encode()
	_, err = c.SendRequest(request)
	return err == nil, err
}

func (c *APIClient) RenameFolder(source *FileInfo, name string) (bool, error) {
	request, err := c.NewPutRequest(fmt.Sprintf("/rest/folders/%s", source.ID), FolderUpdatePutRequest{Name: name})
	if err != nil {
		return false, err
	}
	q := request.URL.Query()
	q.Add("id:in", source.ID)
	request.URL.RawQuery = q.Encode()
	_, err = c.SendRequest(request)
	return err == nil, err
}

func (c *APIClient) RenameFile(source *FileInfo, name string, replace bool) (bool, error) {
	request, err := c.NewPutRequest(fmt.Sprintf("/rest/files/%s", source.ID), FileUpdateRequest{Name: name, Replace: replace})
	if err != nil {
		return false, err
	}
	_, err = c.SendRequest(request)
	return err == nil, err
}

func (c *APIClient) NewGetRequest(path string) (*http.Request, error) {
	return c.newRequest("GET", path, nil)
}

func (c *APIClient) NewPostRequest(path string, v any) (*http.Request, error) {
	return c.newJSONRequest("POST", path, v)
}

func (c *APIClient) NewPutRequest(path string, v any) (*http.Request, error) {
	return c.newJSONRequest("PUT", path, v)
}

func (c *APIClient) newJSONRequest(method, path string, v any) (*http.Request, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	req, err := c.newRequest(method, path, bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

func (c *APIClient) newRequest(method, path string, body io.Reader) (*http.Request, error) {
	path = strings.TrimLeft(path, "/")
	req, err := http.NewRequest(method, fmt.Sprintf("%s/%s", c.server, path), body)
	if err != nil {
		return nil, err
	}
	req.Host = c.host
	if c.agentString != "" {
		req.Header.Set("User-Agent", c.agentString)
	}
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
	req.Header.Set("X-Accellion-Version", "28")
	req.Header.Set("X-Request-Context", c.requestId)
	req.Header.Set("X-Forwarded-For", c.remoteAddr)
	return req, nil
}

func (c *APIClient) SendRequest(req *http.Request) (*http.Response, error) {
	return c.sendWith(c.httpClient, req)
}

func (c *APIClient) sendWith(client *http.Client, req *http.Request) (*http.Response, error) {
	response, err := client.Do(req)
	if err != nil {
		c.logger.Debug().Str("method", req.Method).Str("path", req.URL.String()).Err(err).Msg("kiteworks API call errored")
		return nil, err
	}
	if response.StatusCode >= http.StatusBadRequest {
		defer response.Body.Close()
		b, _ := io.ReadAll(response.Body)
		c.logger.Debug().Str("method", req.Method).Str("path", req.URL.String()).Str("body", string(b)).Int("status", response.StatusCode).Msg("kiteworks API call failed")
		return response, NewClientError(response.StatusCode)
	}
	c.logger.Debug().Str("method", req.Method).Str("path", req.URL.String()).Int("status", response.StatusCode).Msg("kiteworks API call success")
	return response, nil
}

type ClientError struct {
	StatusCode int
}

func (e *ClientError) Error() string {
	return fmt.Sprintf("received non 200 response code: %d", e.StatusCode)
}

func NewClientError(statusCode int) *ClientError {
	return &ClientError{statusCode}
}

func AsClientError(err error) ClientError {
	var httpError *ClientError
	if errors.As(err, &httpError) {
		return *httpError
	}
	return *NewClientError(http.StatusInternalServerError)
}

func ParseFolderIDFromPath(path string) string {
	i := strings.LastIndex(path, "/")
	if i == -1 {
		return ""
	}
	return path[i+1:]
}
