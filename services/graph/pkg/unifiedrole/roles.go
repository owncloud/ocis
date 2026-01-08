package unifiedrole

import (
	"cmp"
	"slices"
	"strings"

	libregraph "github.com/owncloud/libre-graph-api-go"
	"google.golang.org/protobuf/proto"

	"github.com/owncloud/reva/v2/pkg/conversions"

	"github.com/owncloud/ocis/v2/ocis-pkg/l10n"
)

const (
	// UnifiedRoleViewerID Unified role viewer id.
	UnifiedRoleViewerID = "b1e2218d-eef8-4d4c-b82d-0f1a1b48f3b5"
	// UnifiedRoleViewerListGrantsID Unified role viewer id.
	UnifiedRoleViewerListGrantsID = "d5041006-ebb3-4b4a-b6a4-7c180ecfb17d"
	// UnifiedRoleSpaceViewerID Unified role space viewer id.
	UnifiedRoleSpaceViewerID = "a8d5fe5e-96e3-418d-825b-534dbdf22b99"
	// UnifiedRoleEditorID Unified role editor id.
	UnifiedRoleEditorID = "fb6c3e19-e378-47e5-b277-9732f9de6e21"
	// UnifiedRoleEditorListGrantsID Unified role editor id.
	UnifiedRoleEditorListGrantsID = "e8ea8b21-abd4-45d2-b893-8d1546378e9e"
	// UnifiedRoleEditorListGrantsWithVersionsID Unified role editor with versions id.
	UnifiedRoleEditorListGrantsWithVersionsID = "0911d62b-1e3f-4778-8b1b-903b7e4e8476"
	// UnifiedRoleSpaceEditorID Unified role space editor id.
	UnifiedRoleSpaceEditorID = "58c63c02-1d89-4572-916a-870abc5a1b7d"
	// UnifiedRoleSpaceEditorWithoutVersionsID Unified role space editor without list/restore versions id.
	UnifiedRoleSpaceEditorWithoutVersionsID = "3284f2d5-0070-4ad8-ac40-c247f7c1fb27"
	// UnifiedRoleSpaceEditorWithoutTrashbinID Unified role space editor without list/restore resources in trashbin id.
	UnifiedRoleSpaceEditorWithoutTrashbinID = "8f4701d9-c68f-4109-a482-88e22ee32805"
	// UnifiedRoleFileEditorID Unified role file editor id.
	UnifiedRoleFileEditorID = "2d00ce52-1fc2-4dbc-8b95-a73b73395f5a"
	// UnifiedRoleFileEditorListGrantsID Unified role file editor id.
	UnifiedRoleFileEditorListGrantsID = "c1235aea-d106-42db-8458-7d5610fb0a67"
	// UnifiedRoleFileEditorListGrantsWithVersionsID Unified role file editor id.
	UnifiedRoleFileEditorListGrantsWithVersionsID = "b173329d-cf2e-42f0-a595-ee410645d840"
	// UnifiedRoleEditorLiteID Unified role editor-lite id.
	UnifiedRoleEditorLiteID = "1c996275-f1c9-4e71-abdf-a42f6495e960"
	// UnifiedRoleManagerID Unified role manager id.
	UnifiedRoleManagerID = "312c0871-5ef7-4b3a-85b6-0e4074c64049"
	// UnifiedRoleSecureViewerID Unified role secure viewer id.
	UnifiedRoleSecureViewerID = "aa97fe03-7980-45ac-9e50-b325749fd7e6"
	// UnifiedRoleDeniedID Unified role to deny all access.
	UnifiedRoleDeniedID = "63e64e19-8d43-42ec-a738-2b6af2610efa"

	// Wile the below conditions follow the SDDL syntax, they are not parsed anywhere. We use them as strings to
	// represent the constraints that a role definition applies to. For the actual syntax, see the SDDL documentation
	// at https://learn.microsoft.com/en-us/windows/win32/secauthz/security-descriptor-definition-language-for-conditional-aces-#conditional-expressions

	// Some roles apply to a specific type of resource, for example, a role that applies to a file or a folder.
	// @Resource is the placeholder for the resource that the role is applied to
	// .Root, .Folder and .File are facets of the driveItem resource that indicate the type of the resource if they are present.

	// UnifiedRoleConditionDrive defines constraint that matches a Driveroot/Spaceroot
	UnifiedRoleConditionDrive = "exists @Resource.Root"
	// UnifiedRoleConditionFolder defines constraints that matches a DriveItem representing a Folder
	UnifiedRoleConditionFolder = "exists @Resource.Folder"
	// UnifiedRoleConditionFile defines a constraint that matches a DriveItem representing a File
	UnifiedRoleConditionFile = "exists @Resource.File"

	// Some roles apply to a specific type of user, for example, a role that applies to a federated user.
	// @Subject is the placeholder for the subject that the role is applied to. For sharing roles this is the user that the resource is shared with.
	// .UserType is the type of the user: 'Member' for a member of the organization, 'Guest' for a guest user, 'Federated' for a federated user.

	// UnifiedRoleConditionFederatedUser defines a constraint that matches a federated user
	UnifiedRoleConditionFederatedUser = "@Subject.UserType==\"Federated\""

	// For federated sharing we need roles that combine the constraints for the resource and the user.
	// UnifiedRoleConditionFileFederatedUser defines a constraint that matches a File and a federated user
	UnifiedRoleConditionFileFederatedUser = UnifiedRoleConditionFile + " && " + UnifiedRoleConditionFederatedUser
	// UnifiedRoleConditionFolderFederatedUser defines a constraint that matches a Folder and a federated user
	UnifiedRoleConditionFolderFederatedUser = UnifiedRoleConditionFolder + " && " + UnifiedRoleConditionFederatedUser

	DriveItemPermissionsCreate = "libre.graph/driveItem/permissions/create"
	DriveItemChildrenCreate    = "libre.graph/driveItem/children/create"
	DriveItemStandardDelete    = "libre.graph/driveItem/standard/delete"
	DriveItemPathRead          = "libre.graph/driveItem/path/read"
	DriveItemQuotaRead         = "libre.graph/driveItem/quota/read"
	DriveItemContentRead       = "libre.graph/driveItem/content/read"
	DriveItemUploadCreate      = "libre.graph/driveItem/upload/create"
	DriveItemPermissionsRead   = "libre.graph/driveItem/permissions/read"
	DriveItemChildrenRead      = "libre.graph/driveItem/children/read"
	DriveItemVersionsRead      = "libre.graph/driveItem/versions/read"
	DriveItemDeletedRead       = "libre.graph/driveItem/deleted/read"
	DriveItemPathUpdate        = "libre.graph/driveItem/path/update"
	DriveItemPermissionsDelete = "libre.graph/driveItem/permissions/delete"
	DriveItemDeletedDelete     = "libre.graph/driveItem/deleted/delete"
	DriveItemVersionsUpdate    = "libre.graph/driveItem/versions/update"
	DriveItemDeletedUpdate     = "libre.graph/driveItem/deleted/update"
	DriveItemBasicRead         = "libre.graph/driveItem/basic/read"
	DriveItemPermissionsUpdate = "libre.graph/driveItem/permissions/update"
	DriveItemPermissionsDeny   = "libre.graph/driveItem/permissions/deny"
)

var (
	// Sharing Viewer (Files + Folders)

	// SecureViewer
	// UnifiedRole SecureViewer, Role DisplayName (resolves directly)
	_secureViewerUnifiedRoleDisplayName = l10n.Template("Can view (secure)")
	// UnifiedRole SecureViewer, Role Description (resolves directly)
	_secureViewerUnifiedRoleDescription = l10n.Template("View only documents, images and PDFs. Watermarks will be applied.")
	// UnifiedRole SecureViewer, Permissions
	_secureViewerRole = conversions.NewSecureViewerRole()

	// Viewer
	// UnifiedRole Viewer, Role DisplayName (resolves directly)
	_viewerUnifiedRoleDisplayName = l10n.Template("Can view")
	// UnifiedRole Viewer, Role Description (resolves directly)
	_viewerUnifiedRoleDescription = l10n.Template("View and download.")
	// UnifiedRole Viewer, Permissions
	_viewerRole = conversions.NewViewerRole()

	// Viewer + ListGrants
	// UnifiedRole ViewerListGrants, Role DisplayName (resolves directly)
	_viewerListGrantsUnifiedRoleDisplayName = l10n.Template("Can view and show invitees")
	// UnifiedRole ViewerListGrants, Role Description (resolves directly)
	_viewerListGrantsUnifiedRoleDescription = l10n.Template("View, download and show all invited people.")
	// UnifiedRole ViewerListGrants, PermissionsReference
	_viewerListGrantsRole = conversions.NewViewerListGrantsRole()

	// Sharing Folder

	// Denial
	// UnifiedRole FullDenial, Role DisplayName (resolves directly)
	_deniedUnifiedRoleDisplayName = l10n.Template("Cannot access")
	// UnifiedRole FullDenial, Role Description (resolves directly)
	_deniedUnifiedRoleDescription = l10n.Template("Deny all access.")
	// UnifiedRole FullDenial, Permissions
	_deniedRole = conversions.NewDeniedRole()

	// EditorLite
	// UnifiedRole EditorLite, Role DisplayName (resolves directly)
	_editorLiteUnifiedRoleDisplayName = l10n.Template("Can upload")
	// UnifiedRole EditorLite, Role Description (resolves directly)
	_editorLiteUnifiedRoleDescription = l10n.Template("View, download, upload, edit and add.")
	// UnifiedRole EditorLite, Permissions
	_editorLiteRole = conversions.NewEditorLiteRole()

	// Editor
	// UnifiedRole Editor, Role DisplayName (resolves directly)
	_editorUnifiedRoleDisplayName = l10n.Template("Can upload with trashbin")
	// UnifiedRole Editor, Role Description (resolves directly)
	_editorUnifiedRoleDescription = l10n.Template("View, download, upload, edit, add and delete.")
	// UnifiedRole Editor, Permissions
	_editorRole = conversions.NewEditorRole()

	// Editor + ListGrants
	// UnifiedRole EditorListGrants, Role DisplayName (resolves directly)
	_editorListGrantsUnifiedRoleDisplayName = l10n.Template("Can edit")
	// UnifiedRole EditorListGrants Editor, Role Description (resolves directly)
	_editorListGrantsUnifiedRoleDescription = l10n.Template("View, download, upload, edit, add, delete and show all invited people.")
	// UnifiedRole EditorListGrants Editor, Permissions
	_editorListGrantsRole = conversions.NewEditorListGrantsRole()

	// Editor + ListGrants + Version
	// UnifiedRole EditorListGrantsWithVersions, Role DisplayName (resolves directly)
	_editorListGrantsWithVersionsUnifiedRoleDisplayName = l10n.Template("Can edit with versions")
	// UnifiedRole EditorListGrantsWithVersions, Role Description (resolves directly)
	_editorListGrantsWithVersionsUnifiedRoleDescription = l10n.Template("View, download, upload, edit, delete, show all invited people and show all versions.")
	// UnifiedRole EditorListGrantsWithVersions, Permissions
	_editorListGrantsWithVersionsRole = conversions.NewEditorListGrantsWithVersionsRole()

	// Sharing File

	// FileEditor
	// UnifiedRole FileEditor, Role DisplayName (resolves directly)
	_fileEditorUnifiedRoleDisplayName = l10n.Template("Can upload")
	// UnifiedRole FileEditor, Role Description (resolves directly)
	_fileEditorUnifiedRoleDescription = l10n.Template("View, download and edit.")
	// UnifiedRole FileEditor, Permissions
	_fileEditorRole = conversions.NewFileEditorRole()

	// FileEditor + ListGrants
	// UnifiedRole FileEditorListGrants, Role DisplayName (resolves directly)
	_fileEditorListGrantsUnifiedRoleDisplayName = l10n.Template("Can edit")
	// UnifiedRole FileEditorListGrants, Role Description (resolves directly)
	_fileEditorListGrantsUnifiedRoleDescription = l10n.Template("View, download, edit and show all invited people.")
	// UnifiedRole FileEditorListGrants, Permissions
	_fileEditorListGrantsRole = conversions.NewFileEditorListGrantsRole()

	// FileEditor + ListGrants + Versions
	// UnifiedRole FileEditorListGrantsWithVersions, Role DisplayName (resolves directly)
	_fileEditorListGrantsWithVersionsUnifiedRoleDisplayName = l10n.Template("Can edit with versions")
	// UnifiedRole FileEditorListGrantsWithVersions, Role Description (resolves directly)
	_fileEditorListGrantsWithVersionsUnifiedRoleDescription = l10n.Template("View, download, edit, show all invited people and show all versions.")
	// UnifiedRole FileEditorListGrantsWithVersion, Role Permissions
	_fileEditorListGrantsWithVersionsRole = conversions.NewFileEditorListGrantsWithVersionsRole()

	// Space Membership

	// Viewer
	// UnifiedRole SpaceViewer, Role DisplayName (resolves directly)
	_spaceViewerUnifiedRoleDisplayName = l10n.Template("Can view")
	// UnifiedRole SpaceViewer, Role Description (resolves directly)
	_spaceViewerUnifiedRoleDescription = l10n.Template("View and download.")
	// UnifiedRole SpaceViewer, Permissions
	_spaceViewerRole = conversions.NewSpaceViewerRole()

	// Editor without Trashbin
	// UnifiedRole SpaceEditorWithoutTrashbin, Role DisplayName (resolves directly)
	_spaceEditorWithoutTrashbinUnifiedRoleDisplayName = l10n.Template("Can edit")
	// UnifiedRole SpaceEditorWithoutTrashbin, Role Description (resolves directly)
	_spaceEditorWithoutTrashbinUnifiedRoleDescription = l10n.Template("View, download, upload, edit and add.")
	// UnifiedRole SpaceEditorWithoutTrashbin, Permissions
	_spaceEditorWithoutTrashbinRole = conversions.NewSpaceEditorWithoutTrashbinRole()

	// Editor without Versions
	// UnifiedRole SpaceEditorWithoutVersions, Role DisplayName (resolves directly)
	_spaceEditorWithoutVersionsUnifiedRoleDisplayName = l10n.Template("Can edit with trashbin")
	// UnifiedRole SpaceEditorWithoutVersions, Role Description (resolves directly)
	_spaceEditorWithoutVersionsUnifiedRoleDescription = l10n.Template("View, download, upload, edit, add and delete.")
	// UnifiedRole SpaceEditorWithoutVersions, Permissions
	_spaceEditorWithoutVersionsRole = conversions.NewSpaceEditorWithoutVersionsRole()

	// Editor
	// UnifiedRole SpaceEditor, Role DisplayName (resolves directly)
	_spaceEditorUnifiedRoleDisplayName = l10n.Template("Can edit with trashbin and versions")
	// UnifiedRole SpaceEditor, Role Description (resolves directly)
	_spaceEditorUnifiedRoleDescription = l10n.Template("View, download, upload, edit, add, delete and show all versions.")
	// UnifiedRole SpaceEditor, Permissions
	_spaceEditorRole = conversions.NewSpaceEditorRole()

	// Manager
	// UnifiedRole Manager, Role DisplayName (resolves directly)
	_managerUnifiedRoleDisplayName = l10n.Template("Can manage")
	// UnifiedRole Manager, Role Description (resolves directly)
	_managerUnifiedRoleDescription = l10n.Template("View, download, upload, edit, add, delete, show all versions and manage members.")
	// UnifiedRole Manager, Permissions
	_managerRole = conversions.NewManagerRole()

	// unifiedRoleLabel contains the mapping of unified role IDs to their labels.
	unifiedRoleLabel = map[string]string{
		UnifiedRoleViewerID:                           "Viewer",
		UnifiedRoleViewerListGrantsID:                 "ViewerListGrants",
		UnifiedRoleSpaceViewerID:                      "SpaceViewer",
		UnifiedRoleEditorID:                           "Editor",
		UnifiedRoleEditorListGrantsID:                 "EditorListGrants",
		UnifiedRoleEditorListGrantsWithVersionsID:     "EditorListGrantsWithVersions",
		UnifiedRoleSpaceEditorID:                      "SpaceEditor",
		UnifiedRoleSpaceEditorWithoutVersionsID:       "SpaceEditorWithoutVersions",
		UnifiedRoleSpaceEditorWithoutTrashbinID:       "SpaceEditorWithoutTrashbin",
		UnifiedRoleFileEditorID:                       "FileEditor",
		UnifiedRoleFileEditorListGrantsID:             "FileEditorListGrants",
		UnifiedRoleFileEditorListGrantsWithVersionsID: "FileEditorListGrantsWithVersions",
		UnifiedRoleEditorLiteID:                       "EditorLite",
		UnifiedRoleManagerID:                          "Manager",
		UnifiedRoleSecureViewerID:                     "SecureViewer",
		UnifiedRoleDeniedID:                           "Denied",
	}

	// legacyNames contains the legacy role names.
	legacyNames = map[string]string{
		UnifiedRoleViewerID: conversions.RoleViewer,
		// one V1 api the "spaceviewer" role was call "viewer" and the "spaceeditor" was "editor",
		// we need to stay compatible with that
		UnifiedRoleSpaceViewerID:                "viewer",
		UnifiedRoleSpaceEditorID:                "editor",
		UnifiedRoleSpaceEditorWithoutVersionsID: conversions.RoleSpaceEditorWithoutVersions,
		UnifiedRoleEditorID:                     conversions.RoleEditor,
		UnifiedRoleFileEditorID:                 conversions.RoleFileEditor,
		UnifiedRoleEditorLiteID:                 conversions.RoleEditorLite,
		UnifiedRoleManagerID:                    conversions.RoleManager,
		UnifiedRoleSecureViewerID:               conversions.RoleSecureViewer,
	}

	// roleSecureViewer creates a secure viewer role
	roleSecureViewer = func() *libregraph.UnifiedRoleDefinition {
		return &libregraph.UnifiedRoleDefinition{
			Id:          proto.String(UnifiedRoleSecureViewerID),
			DisplayName: proto.String(_secureViewerUnifiedRoleDisplayName),
			Description: proto.String(_secureViewerUnifiedRoleDescription),
			RolePermissions: []libregraph.UnifiedRolePermission{
				{
					AllowedResourceActions: CS3ResourcePermissionsToLibregraphActions(_secureViewerRole.CS3ResourcePermissions()),
					Condition:              proto.String(UnifiedRoleConditionFile),
				},
				{
					AllowedResourceActions: CS3ResourcePermissionsToLibregraphActions(_secureViewerRole.CS3ResourcePermissions()),
					Condition:              proto.String(UnifiedRoleConditionFolder),
				},
			},
			LibreGraphWeight: proto.Int32(0),
		}
	}
	// roleViewer creates a viewer role.
	roleViewer = func() *libregraph.UnifiedRoleDefinition {
		return &libregraph.UnifiedRoleDefinition{
			Id:          proto.String(UnifiedRoleViewerID),
			DisplayName: proto.String(_viewerUnifiedRoleDisplayName),
			Description: proto.String(_viewerUnifiedRoleDescription),
			RolePermissions: []libregraph.UnifiedRolePermission{
				{
					AllowedResourceActions: CS3ResourcePermissionsToLibregraphActions(_viewerRole.CS3ResourcePermissions()),
					Condition:              proto.String(UnifiedRoleConditionFile),
				},
				{
					AllowedResourceActions: CS3ResourcePermissionsToLibregraphActions(_viewerRole.CS3ResourcePermissions()),
					Condition:              proto.String(UnifiedRoleConditionFolder),
				},
				{
					AllowedResourceActions: CS3ResourcePermissionsToLibregraphActions(_viewerRole.CS3ResourcePermissions()),
					Condition:              proto.String(UnifiedRoleConditionFileFederatedUser),
				},
				{
					AllowedResourceActions: CS3ResourcePermissionsToLibregraphActions(_viewerRole.CS3ResourcePermissions()),
					Condition:              proto.String(UnifiedRoleConditionFolderFederatedUser),
				},
			},
			LibreGraphWeight: proto.Int32(0),
		}
	}

	// roleViewerListGrants creates a viewer role.
	roleViewerListGrants = func() *libregraph.UnifiedRoleDefinition {
		return &libregraph.UnifiedRoleDefinition{
			Id:          proto.String(UnifiedRoleViewerListGrantsID),
			DisplayName: proto.String(_viewerListGrantsUnifiedRoleDisplayName),
			Description: proto.String(_viewerListGrantsUnifiedRoleDescription),
			RolePermissions: []libregraph.UnifiedRolePermission{
				{
					AllowedResourceActions: CS3ResourcePermissionsToLibregraphActions(_viewerListGrantsRole.CS3ResourcePermissions()),
					Condition:              proto.String(UnifiedRoleConditionFile),
				},
				{
					AllowedResourceActions: CS3ResourcePermissionsToLibregraphActions(_viewerListGrantsRole.CS3ResourcePermissions()),
					Condition:              proto.String(UnifiedRoleConditionFolder),
				},
				{
					AllowedResourceActions: CS3ResourcePermissionsToLibregraphActions(_viewerListGrantsRole.CS3ResourcePermissions()),
					Condition:              proto.String(UnifiedRoleConditionFileFederatedUser),
				},
				{
					AllowedResourceActions: CS3ResourcePermissionsToLibregraphActions(_viewerListGrantsRole.CS3ResourcePermissions()),
					Condition:              proto.String(UnifiedRoleConditionFolderFederatedUser),
				},
			},
			LibreGraphWeight: proto.Int32(0),
		}
	}

	// roleDenied creates a secure viewer role
	roleDenied = func() *libregraph.UnifiedRoleDefinition {
		return &libregraph.UnifiedRoleDefinition{
			Id:          proto.String(UnifiedRoleDeniedID),
			DisplayName: proto.String(_deniedUnifiedRoleDisplayName),
			Description: proto.String(_deniedUnifiedRoleDescription),
			RolePermissions: []libregraph.UnifiedRolePermission{
				{
					AllowedResourceActions: CS3ResourcePermissionsToLibregraphActions(_deniedRole.CS3ResourcePermissions()),
					Condition:              proto.String(UnifiedRoleConditionFolder),
				},
			},
			LibreGraphWeight: proto.Int32(0),
		}
	}

	// roleEditorLite creates an editor-lite role
	roleEditorLite = func() *libregraph.UnifiedRoleDefinition {
		return &libregraph.UnifiedRoleDefinition{
			Id:          proto.String(UnifiedRoleEditorLiteID),
			DisplayName: proto.String(_editorLiteUnifiedRoleDisplayName),
			Description: proto.String(_editorLiteUnifiedRoleDescription),
			RolePermissions: []libregraph.UnifiedRolePermission{
				{
					AllowedResourceActions: CS3ResourcePermissionsToLibregraphActions(_editorLiteRole.CS3ResourcePermissions()),
					Condition:              proto.String(UnifiedRoleConditionFolder),
				},
			},
			LibreGraphWeight: proto.Int32(0),
		}
	}

	// roleEditor creates an editor role.
	roleEditor = func() *libregraph.UnifiedRoleDefinition {
		return &libregraph.UnifiedRoleDefinition{
			Id:          proto.String(UnifiedRoleEditorID),
			DisplayName: proto.String(_editorUnifiedRoleDisplayName),
			Description: proto.String(_editorUnifiedRoleDescription),
			RolePermissions: []libregraph.UnifiedRolePermission{
				{
					AllowedResourceActions: CS3ResourcePermissionsToLibregraphActions(_editorRole.CS3ResourcePermissions()),
					Condition:              proto.String(UnifiedRoleConditionFolder),
				},
				{
					AllowedResourceActions: CS3ResourcePermissionsToLibregraphActions(_editorRole.CS3ResourcePermissions()),
					Condition:              proto.String(UnifiedRoleConditionFolderFederatedUser),
				},
			},
			LibreGraphWeight: proto.Int32(0),
		}
	}

	// roleEditorListGrants creates an editor role.
	roleEditorListGrants = func() *libregraph.UnifiedRoleDefinition {
		return &libregraph.UnifiedRoleDefinition{
			Id:          proto.String(UnifiedRoleEditorListGrantsID),
			DisplayName: proto.String(_editorListGrantsUnifiedRoleDisplayName),
			Description: proto.String(_editorListGrantsUnifiedRoleDescription),
			RolePermissions: []libregraph.UnifiedRolePermission{
				{
					AllowedResourceActions: CS3ResourcePermissionsToLibregraphActions(_editorListGrantsRole.CS3ResourcePermissions()),
					Condition:              proto.String(UnifiedRoleConditionFolder),
				},
				{
					AllowedResourceActions: CS3ResourcePermissionsToLibregraphActions(_editorListGrantsRole.CS3ResourcePermissions()),
					Condition:              proto.String(UnifiedRoleConditionFolderFederatedUser),
				},
			},
			LibreGraphWeight: proto.Int32(0),
		}
	}

	// roleEditorListGrantsWithVersions creates an editor-list-grants-with-versions role.
	roleEditorListGrantsWithVersions = func() *libregraph.UnifiedRoleDefinition {
		return &libregraph.UnifiedRoleDefinition{
			Id:          proto.String(UnifiedRoleEditorListGrantsWithVersionsID),
			DisplayName: proto.String(_editorListGrantsWithVersionsUnifiedRoleDisplayName),
			Description: proto.String(_editorListGrantsWithVersionsUnifiedRoleDescription),
			RolePermissions: []libregraph.UnifiedRolePermission{
				{
					AllowedResourceActions: CS3ResourcePermissionsToLibregraphActions(_editorListGrantsWithVersionsRole.CS3ResourcePermissions()),
					Condition:              proto.String(UnifiedRoleConditionFolder),
				},
				{
					AllowedResourceActions: CS3ResourcePermissionsToLibregraphActions(_editorListGrantsWithVersionsRole.CS3ResourcePermissions()),
					Condition:              proto.String(UnifiedRoleConditionFolderFederatedUser),
				},
			},
			LibreGraphWeight: proto.Int32(0),
		}
	}

	// roleFileEditor creates a file-editor role
	roleFileEditor = func() *libregraph.UnifiedRoleDefinition {
		return &libregraph.UnifiedRoleDefinition{
			Id:          proto.String(UnifiedRoleFileEditorID),
			DisplayName: proto.String(_fileEditorUnifiedRoleDisplayName),
			Description: proto.String(_fileEditorUnifiedRoleDescription),
			RolePermissions: []libregraph.UnifiedRolePermission{
				{
					AllowedResourceActions: CS3ResourcePermissionsToLibregraphActions(_fileEditorRole.CS3ResourcePermissions()),
					Condition:              proto.String(UnifiedRoleConditionFile),
				},
				{
					AllowedResourceActions: CS3ResourcePermissionsToLibregraphActions(_fileEditorRole.CS3ResourcePermissions()),
					Condition:              proto.String(UnifiedRoleConditionFileFederatedUser),
				},
			},
			LibreGraphWeight: proto.Int32(0),
		}
	}

	// roleFileEditorListGrants creates a file-editor role
	roleFileEditorListGrants = func() *libregraph.UnifiedRoleDefinition {
		return &libregraph.UnifiedRoleDefinition{
			Id:          proto.String(UnifiedRoleFileEditorListGrantsID),
			DisplayName: proto.String(_fileEditorListGrantsUnifiedRoleDisplayName),
			Description: proto.String(_fileEditorListGrantsUnifiedRoleDescription),
			RolePermissions: []libregraph.UnifiedRolePermission{
				{
					AllowedResourceActions: CS3ResourcePermissionsToLibregraphActions(_fileEditorListGrantsRole.CS3ResourcePermissions()),
					Condition:              proto.String(UnifiedRoleConditionFile),
				},
				{
					AllowedResourceActions: CS3ResourcePermissionsToLibregraphActions(_fileEditorListGrantsRole.CS3ResourcePermissions()),
					Condition:              proto.String(UnifiedRoleConditionFileFederatedUser),
				},
			},
			LibreGraphWeight: proto.Int32(0),
		}
	}

	// roleFileEditorListGrantsWithVersions creates a file-editor role
	roleFileEditorListGrantsWithVersions = func() *libregraph.UnifiedRoleDefinition {
		return &libregraph.UnifiedRoleDefinition{
			Id:          proto.String(UnifiedRoleFileEditorListGrantsWithVersionsID),
			DisplayName: proto.String(_fileEditorListGrantsWithVersionsUnifiedRoleDisplayName),
			Description: proto.String(_fileEditorListGrantsWithVersionsUnifiedRoleDescription),
			RolePermissions: []libregraph.UnifiedRolePermission{
				{
					AllowedResourceActions: CS3ResourcePermissionsToLibregraphActions(_fileEditorListGrantsWithVersionsRole.CS3ResourcePermissions()),
					Condition:              proto.String(UnifiedRoleConditionFile),
				},
				{
					AllowedResourceActions: CS3ResourcePermissionsToLibregraphActions(_fileEditorListGrantsWithVersionsRole.CS3ResourcePermissions()),
					Condition:              proto.String(UnifiedRoleConditionFileFederatedUser),
				},
			},
			LibreGraphWeight: proto.Int32(0),
		}
	}

	// roleSpaceViewer creates a spaceviewer role
	roleSpaceViewer = func() *libregraph.UnifiedRoleDefinition {
		return &libregraph.UnifiedRoleDefinition{
			Id:          proto.String(UnifiedRoleSpaceViewerID),
			DisplayName: proto.String(_spaceViewerUnifiedRoleDisplayName),
			Description: proto.String(_spaceViewerUnifiedRoleDescription),
			RolePermissions: []libregraph.UnifiedRolePermission{
				{
					AllowedResourceActions: CS3ResourcePermissionsToLibregraphActions(_spaceViewerRole.CS3ResourcePermissions()),
					Condition:              proto.String(UnifiedRoleConditionDrive),
				},
			},
			LibreGraphWeight: proto.Int32(0),
		}
	}

	// roleSpaceEditorWithoutTrashbin creates an editor without trashbin role
	roleSpaceEditorWithoutTrashbin = func() *libregraph.UnifiedRoleDefinition {
		return &libregraph.UnifiedRoleDefinition{
			Id:          proto.String(UnifiedRoleSpaceEditorWithoutTrashbinID),
			DisplayName: proto.String(_spaceEditorWithoutTrashbinUnifiedRoleDisplayName),
			Description: proto.String(_spaceEditorWithoutTrashbinUnifiedRoleDescription),
			RolePermissions: []libregraph.UnifiedRolePermission{
				{
					AllowedResourceActions: CS3ResourcePermissionsToLibregraphActions(_spaceEditorWithoutTrashbinRole.CS3ResourcePermissions()),
					Condition:              proto.String(UnifiedRoleConditionDrive),
				},
			},
			LibreGraphWeight: proto.Int32(0),
		}
	}

	// roleSpaceEditorWithoutVersions creates an editor without versions role
	roleSpaceEditorWithoutVersions = func() *libregraph.UnifiedRoleDefinition {
		return &libregraph.UnifiedRoleDefinition{
			Id:          proto.String(UnifiedRoleSpaceEditorWithoutVersionsID),
			DisplayName: proto.String(_spaceEditorWithoutVersionsUnifiedRoleDisplayName),
			Description: proto.String(_spaceEditorWithoutVersionsUnifiedRoleDescription),
			RolePermissions: []libregraph.UnifiedRolePermission{
				{
					AllowedResourceActions: CS3ResourcePermissionsToLibregraphActions(_spaceEditorWithoutVersionsRole.CS3ResourcePermissions()),
					Condition:              proto.String(UnifiedRoleConditionDrive),
				},
			},
			LibreGraphWeight: proto.Int32(0),
		}
	}

	// roleSpaceEditor creates an editor role
	roleSpaceEditor = func() *libregraph.UnifiedRoleDefinition {
		return &libregraph.UnifiedRoleDefinition{
			Id:          proto.String(UnifiedRoleSpaceEditorID),
			DisplayName: proto.String(_spaceEditorUnifiedRoleDisplayName),
			Description: proto.String(_spaceEditorUnifiedRoleDescription),
			RolePermissions: []libregraph.UnifiedRolePermission{
				{
					AllowedResourceActions: CS3ResourcePermissionsToLibregraphActions(_spaceEditorRole.CS3ResourcePermissions()),
					Condition:              proto.String(UnifiedRoleConditionDrive),
				},
			},
			LibreGraphWeight: proto.Int32(0),
		}
	}

	// roleManager creates a manager role
	roleManager = func() *libregraph.UnifiedRoleDefinition {
		return &libregraph.UnifiedRoleDefinition{
			Id:          proto.String(UnifiedRoleManagerID),
			DisplayName: proto.String(_managerUnifiedRoleDisplayName),
			Description: proto.String(_managerUnifiedRoleDescription),
			RolePermissions: []libregraph.UnifiedRolePermission{
				{
					AllowedResourceActions: CS3ResourcePermissionsToLibregraphActions(_managerRole.CS3ResourcePermissions()),
					Condition:              proto.String(UnifiedRoleConditionDrive),
				},
			},
			LibreGraphWeight: proto.Int32(0),
		}
	}
)

// GetRoles returns a role filter that matches the provided resources
func GetRoles(f RoleFilter) []*libregraph.UnifiedRoleDefinition {
	return filterRoles(buildInRoles(), f)
}

// GetRole returns a role filter that matches the provided resources
func GetRole(f RoleFilter) (*libregraph.UnifiedRoleDefinition, error) {
	roles := filterRoles(buildInRoles(), f)
	if len(roles) == 0 {
		return nil, ErrUnknownRole
	}

	return roles[0], nil
}

// GetRolesByPermissions returns a list of role definitions
// that match the provided actions and constraints
func GetRolesByPermissions(roleSet []*libregraph.UnifiedRoleDefinition, actions []string, constraints string, listFederatedRoles, descending bool) []*libregraph.UnifiedRoleDefinition {
	roles := make([]*libregraph.UnifiedRoleDefinition, 0, len(roleSet))

	for _, role := range roleSet {
		var match bool

		for _, permission := range role.GetRolePermissions() {
			// this is a dirty comparison because we are not really parsing the SDDL, but as long as we && the conditions we are good
			isFederatedRole := strings.Contains(permission.GetCondition(), UnifiedRoleConditionFederatedUser)
			switch {
			case !strings.Contains(permission.GetCondition(), constraints):
				continue
			case listFederatedRoles && !isFederatedRole:
				continue
			case !listFederatedRoles && isFederatedRole:
				continue
			}

			for i, action := range permission.GetAllowedResourceActions() {
				if !slices.Contains(actions, action) {
					break
				}
				if i == len(permission.GetAllowedResourceActions())-1 {
					match = true
				}
			}
			if role.GetId() == UnifiedRoleDeniedID && slices.Contains(actions, DriveItemPermissionsDeny) {
				match = true
			}
			if match {
				break
			}
		}

		if match {
			roles = append(roles, role)
		}

	}

	return weightRoles(roles, constraints, descending)
}

// GetUnifiedRoleLabel returns the label for the provided unified role ID
func GetUnifiedRoleLabel(unifiedRoleId string) string {
	return unifiedRoleLabel[unifiedRoleId]
}

// GetLegacyRoleName returns the legacy role name for the provided role
func GetLegacyRoleName(role libregraph.UnifiedRoleDefinition) string {
	return legacyNames[role.GetId()]
}

// weightRoles sorts the provided role definitions by the number of permissions[n].actions they grant,
// the implementation is optimistic and assumes that the weight relies on the number of available actions.
// descending - false - sorts the roles from least to most permissions
// descending - true - sorts the roles from most to least permissions
func weightRoles(roleSet []*libregraph.UnifiedRoleDefinition, constraints string, descending bool) []*libregraph.UnifiedRoleDefinition {
	slices.SortFunc(roleSet, func(i, j *libregraph.UnifiedRoleDefinition) int {
		var ia []string
		for _, rp := range i.GetRolePermissions() {
			if rp.GetCondition() == constraints {
				ia = append(ia, rp.GetAllowedResourceActions()...)
			}
		}

		var ja []string
		for _, rp := range j.GetRolePermissions() {
			if rp.GetCondition() == constraints {
				ja = append(ja, rp.GetAllowedResourceActions()...)
			}
		}

		switch descending {
		case true:
			return cmp.Compare(len(ja), len(ia))
		default:
			return cmp.Compare(len(ia), len(ja))
		}
	})

	for i, role := range roleSet {
		role.LibreGraphWeight = libregraph.PtrInt32(int32(i) + 1)
	}

	// return for the sake of consistency, optional because the slice is modified in place
	return roleSet
}

// GetAllowedResourceActions returns the allowed resource actions for the provided role by condition
func GetAllowedResourceActions(role *libregraph.UnifiedRoleDefinition, condition string) []string {
	if role == nil {
		return []string{}
	}

	for _, p := range role.GetRolePermissions() {
		if p.GetCondition() == condition {
			return p.GetAllowedResourceActions()
		}
	}

	return []string{}
}
