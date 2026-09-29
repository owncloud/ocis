package kiteworks

import (
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"

	"github.com/owncloud/reva/v2/pkg/conversions"
	"github.com/owncloud/reva/v2/pkg/storage/fs/kiteworks/kwlib"
)

// permissionSet returns the CS3 ResourcePermissions for fi based on its KW permissions.
func permissionSet(fi *kwlib.FileInfo) *provider.ResourcePermissions {
	if fi.IsDir() {
		return folderPermissions(fi)
	}
	return filePermissions(fi)
}

func folderPermissions(fi *kwlib.FileInfo) *provider.ResourcePermissions {
	canView := fi.HasPermission(kwlib.PermPropertiesView)
	return &provider.ResourcePermissions{
		Stat:                 canView,
		GetPath:              canView,
		ListContainer:        canView,
		InitiateFileDownload: fi.HasPermission(kwlib.PermDownload),
		InitiateFileUpload:   fi.HasPermission(kwlib.PermFileAdd),
		CreateContainer:      fi.HasPermission(kwlib.PermFolderAdd),
		Delete:               fi.HasPermission(kwlib.PermFolderDelete),
		Move:                 fi.HasPermission(kwlib.PermFolderMove),
		ListGrants:           fi.HasPermission(kwlib.PermUserView),
		AddGrant:             fi.HasPermission(kwlib.PermUserAdd),
		UpdateGrant:          fi.HasPermission(kwlib.PermUserEdit),
		RemoveGrant:          fi.HasPermission(kwlib.PermUserRemove),
	}
}

func filePermissions(fi *kwlib.FileInfo) *provider.ResourcePermissions {
	canView := fi.HasPermission(kwlib.PermView)
	return &provider.ResourcePermissions{
		Stat:                 canView,
		GetPath:              canView,
		InitiateFileDownload: fi.HasPermission(kwlib.PermDownload),
		InitiateFileUpload:   fi.HasPermission(kwlib.PermVersionCreate),
		ListFileVersions:     fi.HasPermission(kwlib.PermVersionView),
		RestoreFileVersion:   fi.HasPermission(kwlib.PermVersionPromote),
		Delete:               fi.HasPermission(kwlib.PermFileDelete),
		Move:                 fi.HasPermission(kwlib.PermFileMove),
	}
}

// spaceRole returns the CS3 role for the current user on a space root folder.
// Used to build the grants opaque in ListStorageSpaces.
// Precedence: Manager > Editor > Viewer.
func spaceRole(fi *kwlib.FileInfo) *provider.ResourcePermissions {
	switch {
	case fi.HasPermission(kwlib.PermUserAdd):
		return conversions.NewManagerRole().CS3ResourcePermissions()
	case fi.HasPermission(kwlib.PermFileAdd):
		return conversions.NewEditorRole().CS3ResourcePermissions()
	default:
		return conversions.NewViewerRole().CS3ResourcePermissions()
	}
}
