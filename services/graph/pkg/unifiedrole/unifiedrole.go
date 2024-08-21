package unifiedrole

import (
	"cmp"
	"errors"
	"slices"
	"strings"

	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	libregraph "github.com/owncloud/libre-graph-api-go"
	"google.golang.org/protobuf/proto"

	"github.com/owncloud/ocis/v2/ocis-pkg/l10n"

	"github.com/cs3org/reva/v2/pkg/conversions"
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
	// UnifiedRoleFederatedViewerID Unified role federated viewer id.
	UnifiedRoleFederatedViewerID = "be531789-063c-48bf-a9fe-857e6fbee7da"
	// UnifiedRoleFederatedEditorID Unified role federated editor id.
	UnifiedRoleFederatedEditorID = "36279a93-e4e3-4bbb-8a23-53b05b560963"

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

var legacyNames map[string]string = map[string]string{
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
)

// NewViewerUnifiedRole creates a viewer role.
func NewViewerUnifiedRole() *libregraph.UnifiedRoleDefinition {
	r := conversions.NewViewerRole()
	return &libregraph.UnifiedRoleDefinition{
		Id:          proto.String(UnifiedRoleViewerID),
		Description: proto.String(_viewerUnifiedRoleDescription),
		DisplayName: displayName(r),
		RolePermissions: []libregraph.UnifiedRolePermission{
			{
				AllowedResourceActions: convert(r),
				Condition:              proto.String(UnifiedRoleConditionFile),
			},
			{
				AllowedResourceActions: convert(r),
				Condition:              proto.String(UnifiedRoleConditionFolder),
			},
			{
				AllowedResourceActions: convert(r),
				Condition:              proto.String(UnifiedRoleConditionFileFederatedUser),
			},
			{
				AllowedResourceActions: convert(r),
				Condition:              proto.String(UnifiedRoleConditionFolderFederatedUser),
			},
		},
		LibreGraphWeight: proto.Int32(0),
	}
}

// NewSpaceViewerUnifiedRole creates a spaceviewer role
func NewSpaceViewerUnifiedRole() *libregraph.UnifiedRoleDefinition {
	r := conversions.NewSpaceViewerRole()
	return &libregraph.UnifiedRoleDefinition{
		Id:          proto.String(UnifiedRoleSpaceViewerID),
		Description: proto.String(_spaceViewerUnifiedRoleDescription),
		DisplayName: displayName(r),
		RolePermissions: []libregraph.UnifiedRolePermission{
			{
				AllowedResourceActions: convert(r),
				Condition:              proto.String(UnifiedRoleConditionDrive),
			},
		},
		LibreGraphWeight: proto.Int32(0),
	}
}

// NewEditorUnifiedRole creates an editor role.
func NewEditorUnifiedRole() *libregraph.UnifiedRoleDefinition {
	r := conversions.NewEditorRole()
	return &libregraph.UnifiedRoleDefinition{
		Id:          proto.String(UnifiedRoleEditorID),
		Description: proto.String(_editorUnifiedRoleDescription),
		DisplayName: displayName(r),
		RolePermissions: []libregraph.UnifiedRolePermission{
			{
				AllowedResourceActions: convert(r),
				Condition:              proto.String(UnifiedRoleConditionFolder),
			},
			{
				AllowedResourceActions: convert(r),
				Condition:              proto.String(UnifiedRoleConditionFolderFederatedUser),
			},
		},
		LibreGraphWeight: proto.Int32(0),
	}
}

// NewSpaceEditorUnifiedRole creates an editor role
func NewSpaceEditorUnifiedRole() *libregraph.UnifiedRoleDefinition {
	r := conversions.NewSpaceEditorRole()
	return &libregraph.UnifiedRoleDefinition{
		Id:          proto.String(UnifiedRoleSpaceEditorID),
		Description: proto.String(_spaceEditorUnifiedRoleDescription),
		DisplayName: displayName(r),
		RolePermissions: []libregraph.UnifiedRolePermission{
			{
				AllowedResourceActions: convert(r),
				Condition:              proto.String(UnifiedRoleConditionDrive),
			},
		},
		LibreGraphWeight: proto.Int32(0),
	}
}

// NewFileEditorUnifiedRole creates a file-editor role
func NewFileEditorUnifiedRole() *libregraph.UnifiedRoleDefinition {
	r := conversions.NewFileEditorRole()
	return &libregraph.UnifiedRoleDefinition{
		Id:          proto.String(UnifiedRoleFileEditorID),
		Description: proto.String(_fileEditorUnifiedRoleDescription),
		DisplayName: displayName(r),
		RolePermissions: []libregraph.UnifiedRolePermission{
			{
				AllowedResourceActions: convert(r),
				Condition:              proto.String(UnifiedRoleConditionFile),
			},
			{
				AllowedResourceActions: convert(r),
				Condition:              proto.String(UnifiedRoleConditionFileFederatedUser),
			},
		},
		LibreGraphWeight: proto.Int32(0),
	}
}

// NewEditorLiteUnifiedRole creates an editor-lite role
func NewEditorLiteUnifiedRole() *libregraph.UnifiedRoleDefinition {
	r := conversions.NewEditorLiteRole()
	return &libregraph.UnifiedRoleDefinition{
		Id:          proto.String(UnifiedRoleEditorLiteID),
		Description: proto.String(_editorLiteUnifiedRoleDescription),
		DisplayName: displayName(r),
		RolePermissions: []libregraph.UnifiedRolePermission{
			{
				AllowedResourceActions: convert(r),
				Condition:              proto.String(UnifiedRoleConditionFolder),
			},
		},
		LibreGraphWeight: proto.Int32(0),
	}
}

// NewManagerUnifiedRole creates a manager role
func NewManagerUnifiedRole() *libregraph.UnifiedRoleDefinition {
	r := conversions.NewManagerRole()
	return &libregraph.UnifiedRoleDefinition{
		Id:          proto.String(UnifiedRoleManagerID),
		Description: proto.String(_managerUnifiedRoleDescription),
		DisplayName: displayName(r),
		RolePermissions: []libregraph.UnifiedRolePermission{
			{
				AllowedResourceActions: convert(r),
				Condition:              proto.String(UnifiedRoleConditionDrive),
			},
		},
		LibreGraphWeight: proto.Int32(0),
	}
}

// NewSecureViewerUnifiedRole creates a secure viewer role
func NewSecureViewerUnifiedRole() *libregraph.UnifiedRoleDefinition {
	r := conversions.NewSecureViewerRole()
	return &libregraph.UnifiedRoleDefinition{
		Id:          proto.String(UnifiedRoleSecureViewerID),
		Description: proto.String(_secureViewerUnifiedRoleDescription),
		DisplayName: displayName(r),
		RolePermissions: []libregraph.UnifiedRolePermission{
			{
				AllowedResourceActions: convert(r),
				Condition:              proto.String(UnifiedRoleConditionFile),
			},
			{
				AllowedResourceActions: convert(r),
				Condition:              proto.String(UnifiedRoleConditionFolder),
			},
		},
		LibreGraphWeight: proto.Int32(0),
	}
}

// NewUnifiedRoleFromID returns a unified role definition from the provided id
func NewUnifiedRoleFromID(id string) (*libregraph.UnifiedRoleDefinition, error) {
	for _, definition := range GetBuiltinRoleDefinitionList() {
		if definition.GetId() != id {
			continue
		}

		return definition, nil
	}

	return nil, errors.New("role not found")
}

func GetBuiltinRoleDefinitionList() []*libregraph.UnifiedRoleDefinition {
	return []*libregraph.UnifiedRoleDefinition{
		NewViewerUnifiedRole(),
		NewSpaceViewerUnifiedRole(),
		NewEditorUnifiedRole(),
		NewSpaceEditorUnifiedRole(),
		NewFileEditorUnifiedRole(),
		NewEditorLiteUnifiedRole(),
		NewManagerUnifiedRole(),
		NewSecureViewerUnifiedRole(),
	}
}

// GetApplicableRoleDefinitionsForActions returns a list of role definitions
// that match the provided actions and constraints
func GetApplicableRoleDefinitionsForActions(actions []string, constraints string, listFederatedRoles, descending bool) []*libregraph.UnifiedRoleDefinition {
	builtin := GetBuiltinRoleDefinitionList()
	definitions := make([]*libregraph.UnifiedRoleDefinition, 0, len(builtin))

	for _, definition := range builtin {
		var definitionMatch bool

		for _, permission := range definition.GetRolePermissions() {
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
					definitionMatch = true
				}
			}

			if definitionMatch {
				break
			}
		}

		if definitionMatch {
			definitions = append(definitions, definition)
		}

	}

	return WeightRoleDefinitions(definitions, constraints, descending)
}

// WeightRoleDefinitions sorts the provided role definitions by the number of permissions[n].actions they grant,
// the implementation is optimistic and assumes that the weight relies on the number of available actions.
// descending - false - sorts the roles from least to most permissions
// descending - true - sorts the roles from most to least permissions
func WeightRoleDefinitions(roleDefinitions []*libregraph.UnifiedRoleDefinition, constraints string, descending bool) []*libregraph.UnifiedRoleDefinition {
	slices.SortFunc(roleDefinitions, func(i, j *libregraph.UnifiedRoleDefinition) int {
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

	for i, definition := range roleDefinitions {
		definition.LibreGraphWeight = libregraph.PtrInt32(int32(i) + 1)
	}

	// return for the sage of consistency, optional because the slice is modified in place
	return roleDefinitions
}

// PermissionsToCS3ResourcePermissions converts the provided libregraph UnifiedRolePermissions to a cs3 ResourcePermissions
func PermissionsToCS3ResourcePermissions(unifiedRolePermissions []*libregraph.UnifiedRolePermission) *provider.ResourcePermissions {
	p := &provider.ResourcePermissions{}

	for _, permission := range unifiedRolePermissions {
		for _, allowedResourceAction := range permission.AllowedResourceActions {
			switch allowedResourceAction {
			case DriveItemPermissionsCreate:
				p.AddGrant = true
			case DriveItemChildrenCreate:
				p.CreateContainer = true
			case DriveItemStandardDelete:
				p.Delete = true
			case DriveItemPathRead:
				p.GetPath = true
			case DriveItemQuotaRead:
				p.GetQuota = true
			case DriveItemContentRead:
				p.InitiateFileDownload = true
			case DriveItemUploadCreate:
				p.InitiateFileUpload = true
			case DriveItemPermissionsRead:
				p.ListGrants = true
			case DriveItemChildrenRead:
				p.ListContainer = true
			case DriveItemVersionsRead:
				p.ListFileVersions = true
			case DriveItemDeletedRead:
				p.ListRecycle = true
			case DriveItemPathUpdate:
				p.Move = true
			case DriveItemPermissionsDelete:
				p.RemoveGrant = true
			case DriveItemDeletedDelete:
				p.PurgeRecycle = true
			case DriveItemVersionsUpdate:
				p.RestoreFileVersion = true
			case DriveItemDeletedUpdate:
				p.RestoreRecycleItem = true
			case DriveItemBasicRead:
				p.Stat = true
			case DriveItemPermissionsUpdate:
				p.UpdateGrant = true
			case DriveItemPermissionsDeny:
				p.DenyGrant = true
			}
		}
	}

	return p
}

// CS3ResourcePermissionsToLibregraphActions converts the provided cs3 ResourcePermissions to a list of
// libregraph actions
func CS3ResourcePermissionsToLibregraphActions(p *provider.ResourcePermissions) (actions []string) {
	if p.GetAddGrant() {
		actions = append(actions, DriveItemPermissionsCreate)
	}
	if p.GetCreateContainer() {
		actions = append(actions, DriveItemChildrenCreate)
	}
	if p.GetDelete() {
		actions = append(actions, DriveItemStandardDelete)
	}
	if p.GetGetPath() {
		actions = append(actions, DriveItemPathRead)
	}
	if p.GetGetQuota() {
		actions = append(actions, DriveItemQuotaRead)
	}
	if p.GetInitiateFileDownload() {
		actions = append(actions, DriveItemContentRead)
	}
	if p.GetInitiateFileUpload() {
		actions = append(actions, DriveItemUploadCreate)
	}
	if p.GetListGrants() {
		actions = append(actions, DriveItemPermissionsRead)
	}
	if p.GetListContainer() {
		actions = append(actions, DriveItemChildrenRead)
	}
	if p.GetListFileVersions() {
		actions = append(actions, DriveItemVersionsRead)
	}
	if p.GetListRecycle() {
		actions = append(actions, DriveItemDeletedRead)
	}
	if p.GetMove() {
		actions = append(actions, DriveItemPathUpdate)
	}
	if p.GetRemoveGrant() {
		actions = append(actions, DriveItemPermissionsDelete)
	}
	if p.GetPurgeRecycle() {
		actions = append(actions, DriveItemDeletedDelete)
	}
	if p.GetRestoreFileVersion() {
		actions = append(actions, DriveItemVersionsUpdate)
	}
	if p.GetRestoreRecycleItem() {
		actions = append(actions, DriveItemDeletedUpdate)
	}
	if p.GetStat() {
		actions = append(actions, DriveItemBasicRead)
	}
	if p.GetUpdateGrant() {
		actions = append(actions, DriveItemPermissionsUpdate)
	}
	if p.GetDenyGrant() {
		actions = append(actions, DriveItemPermissionsDeny)
	}
	return actions
}

func GetLegacyName(role libregraph.UnifiedRoleDefinition) string {
	return legacyNames[role.GetId()]
}

// CS3ResourcePermissionsToUnifiedRole tries to find the UnifiedRoleDefinition that matches the supplied
// CS3 ResourcePermissions and constraints.
func CS3ResourcePermissionsToUnifiedRole(p *provider.ResourcePermissions, constraints string) *libregraph.UnifiedRoleDefinition {
	actionSet := map[string]struct{}{}
	for _, action := range CS3ResourcePermissionsToLibregraphActions(p) {
		actionSet[action] = struct{}{}
	}

	var res *libregraph.UnifiedRoleDefinition
	for _, uRole := range GetBuiltinRoleDefinitionList() {
		matchFound := false
		for _, uPerm := range uRole.GetRolePermissions() {
			if uPerm.GetCondition() != constraints {
				// the requested constraints don't match, this isn't our role
				continue
			}

			// if the actions converted from the ResourcePermissions equal the action the defined for the role, we have match
			if resourceActionsEqual(actionSet, uPerm.GetAllowedResourceActions()) {
				matchFound = true
				break
			}
		}
		if matchFound {
			res = uRole
			break
		}
	}
	return res
}

func resourceActionsEqual(targetActionSet map[string]struct{}, actions []string) bool {
	if len(targetActionSet) != len(actions) {
		return false
	}

	for _, action := range actions {
		if _, ok := targetActionSet[action]; !ok {
			return false
		}
	}
	return true
}

func displayName(role *conversions.Role) *string {
	if role == nil {
		return nil
	}

	var displayName string
	switch role.Name {
	case conversions.RoleViewer:
		displayName = _viewerUnifiedRoleDisplayName
	case conversions.RoleSpaceViewer:
		displayName = _spaceViewerUnifiedRoleDisplayName
	case conversions.RoleEditor:
		displayName = _editorUnifiedRoleDisplayName
	case conversions.RoleSpaceEditor:
		displayName = _spaceEditorUnifiedRoleDisplayName
	case conversions.RoleFileEditor:
		displayName = _fileEditorUnifiedRoleDisplayName
	case conversions.RoleEditorLite:
		displayName = _editorLiteUnifiedRoleDisplayName
	case conversions.RoleManager:
		displayName = _managerUnifiedRoleDisplayName
	case conversions.RoleSecureViewer:
		displayName = _secureViewerUnifiedRoleDisplayName
	default:
		return nil
	}
	return proto.String(displayName)
}

func convert(role *conversions.Role) []string {
	actions := make([]string, 0, 8)
	if role == nil && role.CS3ResourcePermissions() == nil {
		return actions
	}
	return CS3ResourcePermissionsToLibregraphActions(role.CS3ResourcePermissions())
}

func GetAllowedResourceActions(role *libregraph.UnifiedRoleDefinition, condition string) []string {
	for _, p := range role.GetRolePermissions() {
		if p.GetCondition() == condition {
			return p.GetAllowedResourceActions()
		}
	}
	return []string{}
}
