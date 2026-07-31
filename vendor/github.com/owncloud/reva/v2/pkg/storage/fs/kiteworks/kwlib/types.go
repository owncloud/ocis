package kwlib

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

const (
	DirectoryType = "d"
	FileType      = "f"

	dateFormat = "2006-01-02T15:04:05+0000"
)

const (
	DownloadPermission = "download"
)

type FileSearch struct {
	Files        []FileInfo `json:"files"`
	Folders      []FileInfo `json:"folders"`
	TotalFolders int        `json:"totalFolders"`
	TotalFiles   int        `json:"totalFiles"`
}

func (fs *FileSearch) FindByParent(parentID string) *FileInfo {
	if parentID == "" {
		if len(fs.Files) > 0 {
			return &fs.Files[0]
		}
		if len(fs.Folders) > 0 {
			return &fs.Folders[0]
		}
	}
	for _, f := range fs.Files {
		if f.ParentID != nil && *f.ParentID == parentID {
			return &f
		}
	}
	for _, f := range fs.Folders {
		if f.ParentID != nil && *f.ParentID == parentID {
			return &f
		}
	}
	return nil
}

type Permission struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Allowed bool   `json:"allowed"`
}

type FileInfo struct {
	ID             string           `json:"id"`
	ParentID       *string          `json:"parentId"`
	Type           string           `json:"type"`
	Name           string           `json:"name"`
	Path           string           `json:"path"`
	Size           *int64           `json:"size"`
	Modified       Time             `json:"modified"`
	ClientModified *Time            `json:"clientModified"`
	FingerPrints   FileFingerPrints `json:"fingerprints"`
	SyncAble       *bool            `json:"syncable"`
	IsFavorite     *bool            `json:"isFavorite"`
	IsShared       *bool            `json:"isShared"`
	PermaLink      string           `json:"permalink"`
	Permissions    []Permission     `json:"permissions"`
	Creator        User             `json:"creator"`
}

type FileFingerPrints []FileFingerPrint

func (fp FileFingerPrints) FindHash(algo string) string {
	for _, f := range fp {
		if f.Algo == algo {
			return f.Hash
		}
	}
	return ""
}

type FileFingerPrint struct {
	Algo string `json:"algo"`
	Hash string `json:"hash"`
}

type Time time.Time

func (t *Time) MarshalJSON() ([]byte, error) {
	if t == nil {
		return []byte("null"), nil
	}
	return json.Marshal(time.Time(*t).Format(dateFormat))
}

func (t *Time) UnmarshalJSON(data []byte) error {
	if string(data) == "" {
		return nil
	}
	parsed, err := time.Parse(dateFormat, strings.Trim(string(data), "\""))
	if err != nil {
		return err
	}
	*t = Time(parsed)
	return nil
}

type DirectoryInfo struct {
	Data []FileInfo `json:"data"`
}

func (fi *FileInfo) IsFile() bool {
	return fi != nil && fi.Type == FileType
}

func (fi *FileInfo) IsDir() bool {
	return fi != nil && fi.Type == DirectoryType
}

func (fi *FileInfo) IsSharedWithUser() bool {
	if fi == nil || fi.IsShared == nil || !*fi.IsShared {
		return false
	}
	if fi.ParentID == nil {
		return true
	}
	return *fi.ParentID == "0"
}

func (fi *FileInfo) IsSyncAble() bool {
	if fi == nil {
		return false
	}
	if fi.IsFile() {
		return true
	}
	if fi.SyncAble == nil {
		return true
	}
	return *fi.SyncAble
}

func (fi *FileInfo) HasDownloadPermission() bool {
	if fi == nil {
		return false
	}
	for i := range fi.Permissions {
		if fi.Permissions[i].Name == DownloadPermission {
			return true
		}
	}
	return false
}

func (fi *FileInfo) ETag() string {
	return quoteEtag(fmt.Sprintf("%x", fi.MTime().UnixNano()))
}

func (fi *FileInfo) MTime() time.Time {
	return time.Time(fi.Modified)
}

func quoteEtag(etag string) string {
	if strings.HasPrefix(etag, "W/") {
		return `W/"` + strings.Trim(etag[2:], `"`) + `"`
	}
	return `"` + strings.Trim(etag, `"`) + `"`
}

type CreateDirRequest struct {
	Name     string `json:"name"`
	SyncAble *bool  `json:"syncable,omitempty"`
}

type InitializeUpload struct {
	FileName       string `json:"filename"`
	TotalChunks    int    `json:"totalChunks,omitempty"`
	TotalSize      int64  `json:"totalSize"`
	ClientModified string `json:"clientModified,omitempty"`
}

type UploadResult struct {
	ID        int64  `json:"id"`
	URI       string `json:"uri"`
	TotalSize int64  `json:"totalSize"`
}

type QuotaInfo struct {
	FolderQuotaAllowed int64 `json:"folder_quota_allowed"`
	FolderQuotaUsed    int64 `json:"folder_quota_used"`
}

type User struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	ProfileIcon string `json:"profileIcon"`
	AdminRoleId *int   `json:"adminRoleId"`
	UserTypeId  int    `json:"userTypeId"`
}

type Contact struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
}

type ContactList struct {
	Data []Contact `json:"data"`
}

type FileCopyMove struct {
	DestinationFolderID string `json:"destinationFolderId"`
	Replace             bool   `json:"replace,omitempty"`
}

type FileUpdateRequest struct {
	Name    string `json:"name,omitempty"`
	Replace bool   `json:"replace,omitempty"`
}

type FolderUpdatePutRequest struct {
	Name string `json:"name,omitempty"`
}
