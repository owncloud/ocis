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
