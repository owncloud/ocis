package linktype

import (
	"errors"

	linkv1beta1 "github.com/cs3org/go-cs3apis/cs3/sharing/link/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/cs3org/reva/v2/pkg/storage/utils/grants"
	libregraph "github.com/owncloud/libre-graph-api-go"
	"github.com/owncloud/ocis/v2/services/graph/pkg/unifiedrole"
)

const NoPermissionMatchError = "no matching permission set found"

type LinkType struct {
	Permissions *provider.ResourcePermissions
	linkType    libregraph.SharingLinkType
}

func (l *LinkType) GetPermissions() *provider.ResourcePermissions {
	if l != nil {
		return l.Permissions
	}
	return nil
}

// SharingLinkTypeFromCS3Permissions creates a libregraph link type
// It returns a list of libregraph actions when the conversion is not possible
func SharingLinkTypeFromCS3Permissions(permissions *linkv1beta1.PublicSharePermissions) (*libregraph.SharingLinkType, []string) {
	linkTypes := GetAvailableLinkTypes()
	for _, linkType := range linkTypes {
		if grants.PermissionsEqual(linkType.GetPermissions(), permissions.GetPermissions()) {
			return &linkType.linkType, nil
		}
	}
	return nil, unifiedrole.CS3ResourcePermissionsToLibregraphActions(*permissions.GetPermissions())
}

func CS3ResourcePermissionsFromSharingLink(createLink libregraph.DriveItemCreateLink, info provider.ResourceType) (*provider.ResourcePermissions, error) {
	switch createLink.GetType() {
	case "":
		return nil, errors.New("link type is empty")
	case libregraph.VIEW:
		return NewViewLinkPermissionSet().GetPermissions(), nil
	case libregraph.EDIT:
		if info == provider.ResourceType_RESOURCE_TYPE_FILE {
			return NewFileEditLinkPermissionSet().GetPermissions(), nil
		}
		return NewFolderEditLinkPermissionSet().GetPermissions(), nil
	case libregraph.CREATE_ONLY:
		if info == provider.ResourceType_RESOURCE_TYPE_FILE {
			return nil, errors.New(NoPermissionMatchError)
		}
		return NewFolderDropLinkPermissionSet().GetPermissions(), nil
	case libregraph.UPLOAD:
		if info == provider.ResourceType_RESOURCE_TYPE_FILE {
			return nil, errors.New(NoPermissionMatchError)
		}
		return NewFolderUploadLinkPermissionSet().GetPermissions(), nil
	case libregraph.INTERNAL:
		return NewInternalLinkPermissionSet().GetPermissions(), nil
	default:
		return nil, errors.New(NoPermissionMatchError)
	}
}

func NewInternalLinkPermissionSet() *LinkType {
	return &LinkType{
		Permissions: &provider.ResourcePermissions{},
		linkType:    libregraph.INTERNAL,
	}
}

func NewViewLinkPermissionSet() *LinkType {
	return &LinkType{
		Permissions: &provider.ResourcePermissions{
			GetPath:              true,
			GetQuota:             true,
			InitiateFileDownload: true,
			ListContainer:        true,
			// why is this needed?
			ListRecycle: true,
			Stat:        true,
		},
		linkType: libregraph.VIEW,
	}
}

func NewFileEditLinkPermissionSet() *LinkType {
	return &LinkType{
		Permissions: &provider.ResourcePermissions{
			GetPath:              true,
			GetQuota:             true,
			InitiateFileDownload: true,
			InitiateFileUpload:   true,
			ListContainer:        true,
			// why is this needed?
			ListRecycle: true,
			// why is this needed?
			RestoreRecycleItem: true,
			Stat:               true,
		},
		linkType: libregraph.EDIT,
	}
}

func NewFolderEditLinkPermissionSet() *LinkType {
	return &LinkType{
		Permissions: &provider.ResourcePermissions{
			CreateContainer:      true,
			Delete:               true,
			GetPath:              true,
			GetQuota:             true,
			InitiateFileDownload: true,
			InitiateFileUpload:   true,
			ListContainer:        true,
			// why is this needed?
			ListRecycle: true,
			Move:        true,
			// why is this needed?
			RestoreRecycleItem: true,
			Stat:               true,
		},
		linkType: libregraph.EDIT,
	}
}

func NewFolderDropLinkPermissionSet() *LinkType {
	return &LinkType{
		Permissions: &provider.ResourcePermissions{
			Stat:               true,
			GetPath:            true,
			CreateContainer:    true,
			InitiateFileUpload: true,
		},
		linkType: libregraph.CREATE_ONLY,
	}
}

func NewFolderUploadLinkPermissionSet() *LinkType {
	return &LinkType{
		Permissions: &provider.ResourcePermissions{
			CreateContainer:      true,
			GetPath:              true,
			GetQuota:             true,
			InitiateFileDownload: true,
			InitiateFileUpload:   true,
			ListContainer:        true,
			ListRecycle:          true,
			Stat:                 true,
		},
		linkType: libregraph.UPLOAD,
	}
}

func GetAvailableLinkTypes() []*LinkType {
	return []*LinkType{
		NewInternalLinkPermissionSet(),
		NewViewLinkPermissionSet(),
		NewFolderUploadLinkPermissionSet(),
		NewFileEditLinkPermissionSet(),
		NewFolderEditLinkPermissionSet(),
		NewFolderDropLinkPermissionSet(),
	}
}
