package unifiedrole

import (
	"cmp"
	"slices"

	libregraph "github.com/owncloud/libre-graph-api-go"
	"google.golang.org/protobuf/proto"

	"github.com/cs3org/reva/v2/pkg/conversions"

	"github.com/owncloud/ocis/v2/ocis-pkg/l10n"
)

const (
	// UnifiedRoleViewerID Unified role viewer id.
	UnifiedRoleViewerID = "b1e2218d-eef8-4d4c-b82d-0f1a1b48f3b5"
	// UnifiedRoleSpaceViewerID Unified role space viewer id.
	UnifiedRoleSpaceViewerID = "a8d5fe5e-96e3-418d-825b-534dbdf22b99"
	// UnifiedRoleEditorID Unified role editor id.
	UnifiedRoleEditorID = "fb6c3e19-e378-47e5-b277-9732f9de6e21"
	// UnifiedRoleSpaceEditorID Unified role space editor id.
	UnifiedRoleSpaceEditorID = "58c63c02-1d89-4572-916a-870abc5a1b7d"
	// UnifiedRoleFileEditorID Unified role file editor id.
	UnifiedRoleFileEditorID = "2d00ce52-1fc2-4dbc-8b95-a73b73395f5a"
	// UnifiedRoleEditorLiteID Unified role editor-lite id.
	UnifiedRoleEditorLiteID = "1c996275-f1c9-4e71-abdf-a42f6495e960"
	// UnifiedRoleManagerID Unified role manager id.
	UnifiedRoleManagerID = "312c0871-5ef7-4b3a-85b6-0e4074c64049"
	// UnifiedRoleSecureViewerID Unified role secure viewer id.
	UnifiedRoleSecureViewerID = "aa97fe03-7980-45ac-9e50-b325749fd7e6"

	// UnifiedRoleConditionDrive defines constraint that matches a Driveroot/Spaceroot
	UnifiedRoleConditionDrive = "exists @Resource.Root"
	// UnifiedRoleConditionFolder defines constraints that matches a DriveItem representing a Folder
	UnifiedRoleConditionFolder = "exists @Resource.Folder"
	// UnifiedRoleConditionFile defines a constraint that matches a DriveItem representing a File
	UnifiedRoleConditionFile = "exists @Resource.File"

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
	// UnifiedRole Viewer, Role Description (resolves directly)
	_viewerUnifiedRoleDescription = l10n.Template("View and download.")

	// UnifiedRole Viewer, Role DisplayName (resolves directly)
	_viewerUnifiedRoleDisplayName = l10n.Template("Can view")

	// UnifiedRole SpaceViewer, Role Description (resolves directly)
	_spaceViewerUnifiedRoleDescription = l10n.Template("View and download.")

	// UnifiedRole SpaseViewer, Role DisplayName (resolves directly)
	_spaceViewerUnifiedRoleDisplayName = l10n.Template("Can view")

	// UnifiedRole Editor, Role Description (resolves directly)
	_editorUnifiedRoleDescription = l10n.Template("View, download, upload, edit, add and delete.")

	// UnifiedRole Editor, Role DisplayName (resolves directly)
	_editorUnifiedRoleDisplayName = l10n.Template("Can edit")

	// UnifiedRole SpaseEditor, Role Description (resolves directly)
	_spaceEditorUnifiedRoleDescription = l10n.Template("View, download, upload, edit, add and delete.")

	// UnifiedRole SpaseEditor, Role DisplayName (resolves directly)
	_spaceEditorUnifiedRoleDisplayName = l10n.Template("Can edit")

	// UnifiedRole FileEditor, Role Description (resolves directly)
	_fileEditorUnifiedRoleDescription = l10n.Template("View, download and edit.")

	// UnifiedRole FileEditor, Role DisplayName (resolves directly)
	_fileEditorUnifiedRoleDisplayName = l10n.Template("Can edit")

	// UnifiedRole EditorLite, Role Description (resolves directly)
	_editorLiteUnifiedRoleDescription = l10n.Template("View, download and upload.")

	// UnifiedRole EditorLite, Role DisplayName (resolves directly)
	_editorLiteUnifiedRoleDisplayName = l10n.Template("Can upload")

	// UnifiedRole Manager, Role Description (resolves directly)
	_managerUnifiedRoleDescription = l10n.Template("View, download, upload, edit, add, delete and manage members.")

	// UnifiedRole Manager, Role DisplayName (resolves directly)
	_managerUnifiedRoleDisplayName = l10n.Template("Can manage")

	// UnifiedRole SecureViewer, Role Description (resolves directly)
	_secureViewerUnifiedRoleDescription = l10n.Template("View only documents, images and PDFs. Watermarks will be applied.")

	// UnifiedRole SecureViewer, Role DisplayName (resolves directly)
	_secureViewerUnifiedRoleDisplayName = l10n.Template("Can view (secure)")

	// legacyNames contains the legacy role names.
	legacyNames = map[string]string{
		UnifiedRoleViewerID: conversions.RoleViewer,
		// one V1 api the "spaceviewer" role was call "viewer" and the "spaceeditor" was "editor",
		// we need to stay compatible with that
		UnifiedRoleSpaceViewerID:  "viewer",
		UnifiedRoleSpaceEditorID:  "editor",
		UnifiedRoleEditorID:       conversions.RoleEditor,
		UnifiedRoleFileEditorID:   conversions.RoleFileEditor,
		UnifiedRoleEditorLiteID:   conversions.RoleEditorLite,
		UnifiedRoleManagerID:      conversions.RoleManager,
		UnifiedRoleSecureViewerID: conversions.RoleSecureViewer,
	}

	// buildInRoles contains the built-in roles.
	buildInRoles = []*libregraph.UnifiedRoleDefinition{
		roleViewer,
		roleSpaceViewer,
		roleEditor,
		roleSpaceEditor,
		roleFileEditor,
		roleEditorLite,
		roleManager,
		roleSecureViewer,
	}

	// roleViewer creates a viewer role.
	roleViewer = func() *libregraph.UnifiedRoleDefinition {
		r := conversions.NewViewerRole()
		return &libregraph.UnifiedRoleDefinition{
			Id:          proto.String(UnifiedRoleViewerID),
			Description: proto.String(_viewerUnifiedRoleDescription),
			DisplayName: proto.String(cs3RoleToDisplayName(r)),
			RolePermissions: []libregraph.UnifiedRolePermission{
				{
					AllowedResourceActions: CS3ResourcePermissionsToLibregraphActions(r.CS3ResourcePermissions()),
					Condition:              proto.String(UnifiedRoleConditionFile),
				},
				{
					AllowedResourceActions: CS3ResourcePermissionsToLibregraphActions(r.CS3ResourcePermissions()),
					Condition:              proto.String(UnifiedRoleConditionFolder),
				},
			},
			LibreGraphWeight: proto.Int32(0),
		}
	}()

	// roleSpaceViewer creates a spaceviewer role
	roleSpaceViewer = func() *libregraph.UnifiedRoleDefinition {
		r := conversions.NewSpaceViewerRole()
		return &libregraph.UnifiedRoleDefinition{
			Id:          proto.String(UnifiedRoleSpaceViewerID),
			Description: proto.String(_spaceViewerUnifiedRoleDescription),
			DisplayName: proto.String(cs3RoleToDisplayName(r)),
			RolePermissions: []libregraph.UnifiedRolePermission{
				{
					AllowedResourceActions: CS3ResourcePermissionsToLibregraphActions(r.CS3ResourcePermissions()),
					Condition:              proto.String(UnifiedRoleConditionDrive),
				},
			},
			LibreGraphWeight: proto.Int32(0),
		}
	}()

	// roleEditor creates an editor role.
	roleEditor = func() *libregraph.UnifiedRoleDefinition {
		r := conversions.NewEditorRole()
		return &libregraph.UnifiedRoleDefinition{
			Id:          proto.String(UnifiedRoleEditorID),
			Description: proto.String(_editorUnifiedRoleDescription),
			DisplayName: proto.String(cs3RoleToDisplayName(r)),
			RolePermissions: []libregraph.UnifiedRolePermission{
				{
					AllowedResourceActions: CS3ResourcePermissionsToLibregraphActions(r.CS3ResourcePermissions()),
					Condition:              proto.String(UnifiedRoleConditionFolder),
				},
			},
			LibreGraphWeight: proto.Int32(0),
		}
	}()

	// roleSpaceEditor creates an editor role
	roleSpaceEditor = func() *libregraph.UnifiedRoleDefinition {
		r := conversions.NewSpaceEditorRole()
		return &libregraph.UnifiedRoleDefinition{
			Id:          proto.String(UnifiedRoleSpaceEditorID),
			Description: proto.String(_spaceEditorUnifiedRoleDescription),
			DisplayName: proto.String(cs3RoleToDisplayName(r)),
			RolePermissions: []libregraph.UnifiedRolePermission{
				{
					AllowedResourceActions: CS3ResourcePermissionsToLibregraphActions(r.CS3ResourcePermissions()),
					Condition:              proto.String(UnifiedRoleConditionDrive),
				},
			},
			LibreGraphWeight: proto.Int32(0),
		}
	}()

	// roleFileEditor creates a file-editor role
	roleFileEditor = func() *libregraph.UnifiedRoleDefinition {
		r := conversions.NewFileEditorRole()
		return &libregraph.UnifiedRoleDefinition{
			Id:          proto.String(UnifiedRoleFileEditorID),
			Description: proto.String(_fileEditorUnifiedRoleDescription),
			DisplayName: proto.String(cs3RoleToDisplayName(r)),
			RolePermissions: []libregraph.UnifiedRolePermission{
				{
					AllowedResourceActions: CS3ResourcePermissionsToLibregraphActions(r.CS3ResourcePermissions()),
					Condition:              proto.String(UnifiedRoleConditionFile),
				},
			},
			LibreGraphWeight: proto.Int32(0),
		}
	}()

	// roleEditorLite creates an editor-lite role
	roleEditorLite = func() *libregraph.UnifiedRoleDefinition {
		r := conversions.NewEditorLiteRole()
		return &libregraph.UnifiedRoleDefinition{
			Id:          proto.String(UnifiedRoleEditorLiteID),
			Description: proto.String(_editorLiteUnifiedRoleDescription),
			DisplayName: proto.String(cs3RoleToDisplayName(r)),
			RolePermissions: []libregraph.UnifiedRolePermission{
				{
					AllowedResourceActions: CS3ResourcePermissionsToLibregraphActions(r.CS3ResourcePermissions()),
					Condition:              proto.String(UnifiedRoleConditionFolder),
				},
			},
			LibreGraphWeight: proto.Int32(0),
		}
	}()

	// roleManager creates a manager role
	roleManager = func() *libregraph.UnifiedRoleDefinition {
		r := conversions.NewManagerRole()
		return &libregraph.UnifiedRoleDefinition{
			Id:          proto.String(UnifiedRoleManagerID),
			Description: proto.String(_managerUnifiedRoleDescription),
			DisplayName: proto.String(cs3RoleToDisplayName(r)),
			RolePermissions: []libregraph.UnifiedRolePermission{
				{
					AllowedResourceActions: CS3ResourcePermissionsToLibregraphActions(r.CS3ResourcePermissions()),
					Condition:              proto.String(UnifiedRoleConditionDrive),
				},
			},
			LibreGraphWeight: proto.Int32(0),
		}
	}()

	// roleSecureViewer creates a secure viewer role
	roleSecureViewer = func() *libregraph.UnifiedRoleDefinition {
		r := conversions.NewSecureViewerRole()
		return &libregraph.UnifiedRoleDefinition{
			Id:          proto.String(UnifiedRoleSecureViewerID),
			Description: proto.String(_secureViewerUnifiedRoleDescription),
			DisplayName: proto.String(cs3RoleToDisplayName(r)),
			RolePermissions: []libregraph.UnifiedRolePermission{
				{
					AllowedResourceActions: CS3ResourcePermissionsToLibregraphActions(r.CS3ResourcePermissions()),
					Condition:              proto.String(UnifiedRoleConditionFile),
				},
				{
					AllowedResourceActions: CS3ResourcePermissionsToLibregraphActions(r.CS3ResourcePermissions()),
					Condition:              proto.String(UnifiedRoleConditionFolder),
				},
			},
			LibreGraphWeight: proto.Int32(0),
		}
	}()
)

// GetDefinitions returns a role filter that matches the provided resources
func GetDefinitions(filter RoleFilter) []*libregraph.UnifiedRoleDefinition {
	return filterRoles(buildInRoles, filter)
}

// GetDefinition returns a role filter that matches the provided resources
func GetDefinition(filter RoleFilter) (*libregraph.UnifiedRoleDefinition, error) {
	definitions := filterRoles(buildInRoles, filter)
	if len(definitions) == 0 {
		return nil, ErrUnknownUnifiedRole
	}

	return definitions[0], nil
}

// GetRolesByPermissions returns a list of role definitions
// that match the provided actions and constraints
func GetRolesByPermissions(actions []string, constraints string, descending bool) []*libregraph.UnifiedRoleDefinition {
	roles := GetDefinitions(RoleFilterPermission(RoleFilterMatchSome, constraints, actions...))
	roles = weightDefinitions(roles, constraints, descending)

	return roles
}

// GetLegacyDefinitionName returns the legacy role name for the provided role
func GetLegacyDefinitionName(definition libregraph.UnifiedRoleDefinition) string {
	return legacyNames[definition.GetId()]
}

// weightDefinitions sorts the provided role definitions by the number of permissions[n].actions they grant,
// the implementation is optimistic and assumes that the weight relies on the number of available actions.
// descending - false - sorts the roles from least to most permissions
// descending - true - sorts the roles from most to least permissions
func weightDefinitions(definitions []*libregraph.UnifiedRoleDefinition, constraints string, descending bool) []*libregraph.UnifiedRoleDefinition {
	slices.SortFunc(definitions, func(i, j *libregraph.UnifiedRoleDefinition) int {
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

	for i, definition := range definitions {
		definition.LibreGraphWeight = libregraph.PtrInt32(int32(i) + 1)
	}

	// return for the sake of consistency, optional because the slice is modified in place
	return definitions
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
