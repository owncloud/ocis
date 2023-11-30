package unifiedrole

import (
	"cmp"
	"errors"
	"slices"

	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	libregraph "github.com/owncloud/libre-graph-api-go"
	"google.golang.org/protobuf/proto"

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
	// UnifiedRoleCoownerID Unified role coowner id.
	UnifiedRoleCoownerID = "3a4ba8e9-6a0d-4235-9140-0e7a34007abe"
	// UnifiedRoleUploaderID Unified role uploader id.
	UnifiedRoleUploaderID = "1c996275-f1c9-4e71-abdf-a42f6495e960"
	// UnifiedRoleManagerID Unified role manager id.
	UnifiedRoleManagerID = "312c0871-5ef7-4b3a-85b6-0e4074c64049"

	// UnifiedRoleConditionSelf defines constraint where the principal matches the target resource
	UnifiedRoleConditionSelf = "@Subject.objectId == @Resource.objectId"
	// UnifiedRoleConditionOwner defines constraints when the principal is the owner of the target resource
	UnifiedRoleConditionOwner = "@Subject.objectId Any_of @Resource.owners"
	// UnifiedRoleConditionGrantee does not exist in MS Graph, but we use it to express permissions on shared resources
	UnifiedRoleConditionGrantee = "@Subject.objectId Any_of @Resource.grantee"

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

// NewViewerUnifiedRole creates a viewer role. `sharing` indicates if sharing permission should be added
func NewViewerUnifiedRole(sharing bool) *libregraph.UnifiedRoleDefinition {
	r := conversions.NewViewerRole(sharing)
	return &libregraph.UnifiedRoleDefinition{
		Id:          proto.String(UnifiedRoleViewerID),
		Description: proto.String("Allows reading the shared file or folder"),
		DisplayName: displayName(r),
		RolePermissions: []libregraph.UnifiedRolePermission{
			{
				AllowedResourceActions: convert(r),
				Condition:              proto.String(UnifiedRoleConditionGrantee),
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
		Description: proto.String("Allows reading the shared space"),
		DisplayName: displayName(r),
		RolePermissions: []libregraph.UnifiedRolePermission{
			{
				AllowedResourceActions: convert(r),
				Condition:              proto.String(UnifiedRoleConditionOwner),
			},
		},
		LibreGraphWeight: proto.Int32(0),
	}
}

// NewEditorUnifiedRole creates an editor role. `sharing` indicates if sharing permission should be added
func NewEditorUnifiedRole(sharing bool) *libregraph.UnifiedRoleDefinition {
	r := conversions.NewEditorRole(sharing)
	return &libregraph.UnifiedRoleDefinition{
		Id:          proto.String(UnifiedRoleEditorID),
		Description: proto.String("Allows creating, reading, updating and deleting the shared file or folder"),
		DisplayName: displayName(r),
		RolePermissions: []libregraph.UnifiedRolePermission{
			{
				AllowedResourceActions: convert(r),
				Condition:              proto.String(UnifiedRoleConditionGrantee),
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
		Description: proto.String("Allows creating, reading, updating and deleting file or folder in the shared space"),
		DisplayName: displayName(r),
		RolePermissions: []libregraph.UnifiedRolePermission{
			{
				AllowedResourceActions: convert(r),
				Condition:              proto.String(UnifiedRoleConditionOwner),
			},
		},
		LibreGraphWeight: proto.Int32(0),
	}
}

// NewFileEditorUnifiedRole creates a file-editor role
func NewFileEditorUnifiedRole(sharing bool) *libregraph.UnifiedRoleDefinition {
	r := conversions.NewFileEditorRole(sharing)
	return &libregraph.UnifiedRoleDefinition{
		Id:          proto.String(UnifiedRoleFileEditorID),
		Description: proto.String("Allows reading and updating file"),
		DisplayName: displayName(r),
		RolePermissions: []libregraph.UnifiedRolePermission{
			{
				AllowedResourceActions: convert(r),
				Condition:              proto.String(UnifiedRoleConditionGrantee),
			},
		},
		LibreGraphWeight: proto.Int32(0),
	}
}

// NewCoownerUnifiedRole creates a coowner role.
func NewCoownerUnifiedRole() *libregraph.UnifiedRoleDefinition {
	r := conversions.NewCoownerRole()
	return &libregraph.UnifiedRoleDefinition{
		Id:          proto.String(UnifiedRoleCoownerID),
		Description: proto.String("Grants co-owner permissions on a resource"),
		DisplayName: displayName(r),
		RolePermissions: []libregraph.UnifiedRolePermission{
			{
				AllowedResourceActions: convert(r),
				Condition:              proto.String(UnifiedRoleConditionGrantee),
			},
		},
		LibreGraphWeight: proto.Int32(0),
	}
}

// NewUploaderUnifiedRole creates an uploader role
func NewUploaderUnifiedRole() *libregraph.UnifiedRoleDefinition {
	r := conversions.NewUploaderRole()
	return &libregraph.UnifiedRoleDefinition{
		Id:          proto.String(UnifiedRoleUploaderID),
		Description: proto.String("Allows upload file or folder"),
		DisplayName: displayName(r),
		RolePermissions: []libregraph.UnifiedRolePermission{
			{
				AllowedResourceActions: convert(r),
				Condition:              proto.String(UnifiedRoleConditionGrantee),
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
		Description: proto.String("Grants manager permissions on a resource. Semantically equivalent to co-owner"),
		DisplayName: displayName(r),
		RolePermissions: []libregraph.UnifiedRolePermission{
			{
				AllowedResourceActions: convert(r),
				Condition:              proto.String(UnifiedRoleConditionGrantee),
			},
			{
				AllowedResourceActions: convert(r),
				Condition:              proto.String(UnifiedRoleConditionOwner),
			},
		},
		LibreGraphWeight: proto.Int32(0),
	}
}

// NewUnifiedRoleFromID returns a unified role definition from the provided id
func NewUnifiedRoleFromID(id string, resharing bool) (*libregraph.UnifiedRoleDefinition, error) {
	for _, definition := range GetBuiltinRoleDefinitionList(resharing) {
		if definition.GetId() != id {
			continue
		}

		return definition, nil
	}

	return nil, errors.New("role not found")
}

func GetBuiltinRoleDefinitionList(resharing bool) []*libregraph.UnifiedRoleDefinition {
	return []*libregraph.UnifiedRoleDefinition{
		NewViewerUnifiedRole(resharing),
		NewSpaceViewerUnifiedRole(),
		NewEditorUnifiedRole(resharing),
		NewSpaceEditorUnifiedRole(),
		NewFileEditorUnifiedRole(resharing),
		NewCoownerUnifiedRole(),
		NewUploaderUnifiedRole(),
		NewManagerUnifiedRole(),
	}
}

// GetApplicableRoleDefinitionsForActions returns a list of role definitions
// that match the provided actions and constraints
func GetApplicableRoleDefinitionsForActions(actions []string, constraints string, resharing, descending bool) []*libregraph.UnifiedRoleDefinition {
	var definitions []*libregraph.UnifiedRoleDefinition

	for _, definition := range GetBuiltinRoleDefinitionList(resharing) {
		definitionMatch := true

		for _, permission := range definition.GetRolePermissions() {
			if permission.GetCondition() != constraints {
				definitionMatch = false
				continue
			}

			for _, action := range permission.GetAllowedResourceActions() {
				if !slices.Contains(actions, action) {
					definitionMatch = false
					break
				}
			}

			if definitionMatch {
				break
			}
		}

		if !definitionMatch {
			continue
		}

		definitions = append(definitions, definition)
	}

	return WeightRoleDefinitions(definitions, descending)
}

// WeightRoleDefinitions sorts the provided role definitions by the number of permissions[n].actions they grant,
// the implementation is optimistic and assumes that the weight relies on the number of available actions.
// descending - false - sorts the roles from least to most permissions
// descending - true - sorts the roles from most to least permissions
func WeightRoleDefinitions(roleDefinitions []*libregraph.UnifiedRoleDefinition, descending bool) []*libregraph.UnifiedRoleDefinition {
	slices.SortFunc(roleDefinitions, func(i, j *libregraph.UnifiedRoleDefinition) int {
		var ia []string
		for _, rp := range i.GetRolePermissions() {
			ia = append(ia, rp.GetAllowedResourceActions()...)
		}

		var ja []string
		for _, rp := range j.GetRolePermissions() {
			ja = append(ja, rp.GetAllowedResourceActions()...)
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
func CS3ResourcePermissionsToLibregraphActions(p provider.ResourcePermissions) (actions []string) {
	if p.AddGrant {
		actions = append(actions, DriveItemPermissionsCreate)
	}
	if p.CreateContainer {
		actions = append(actions, DriveItemChildrenCreate)
	}
	if p.Delete {
		actions = append(actions, DriveItemStandardDelete)
	}
	if p.GetPath {
		actions = append(actions, DriveItemPathRead)
	}
	if p.GetQuota {
		actions = append(actions, DriveItemQuotaRead)
	}
	if p.InitiateFileDownload {
		actions = append(actions, DriveItemContentRead)
	}
	if p.InitiateFileUpload {
		actions = append(actions, DriveItemUploadCreate)
	}
	if p.ListGrants {
		actions = append(actions, DriveItemPermissionsRead)
	}
	if p.ListContainer {
		actions = append(actions, DriveItemChildrenRead)
	}
	if p.ListFileVersions {
		actions = append(actions, DriveItemVersionsRead)
	}
	if p.ListRecycle {
		actions = append(actions, DriveItemDeletedRead)
	}
	if p.Move {
		actions = append(actions, DriveItemPathUpdate)
	}
	if p.RemoveGrant {
		actions = append(actions, DriveItemPermissionsDelete)
	}
	if p.PurgeRecycle {
		actions = append(actions, DriveItemDeletedDelete)
	}
	if p.RestoreFileVersion {
		actions = append(actions, DriveItemVersionsUpdate)
	}
	if p.RestoreRecycleItem {
		actions = append(actions, DriveItemDeletedUpdate)
	}
	if p.Stat {
		actions = append(actions, DriveItemBasicRead)
	}
	if p.UpdateGrant {
		actions = append(actions, DriveItemPermissionsUpdate)
	}
	if p.DenyGrant {
		actions = append(actions, DriveItemPermissionsDeny)
	}
	return actions
}

// CS3ResourcePermissionsToUnifiedRole tries to find the UnifiedRoleDefinition that matches the supplied
// CS3 ResourcePermissions and constraints.
func CS3ResourcePermissionsToUnifiedRole(p provider.ResourcePermissions, constraints string, resharing bool) *libregraph.UnifiedRoleDefinition {
	actionSet := map[string]struct{}{}
	for _, action := range CS3ResourcePermissionsToLibregraphActions(p) {
		actionSet[action] = struct{}{}
	}

	var res *libregraph.UnifiedRoleDefinition
	for _, uRole := range GetBuiltinRoleDefinitionList(resharing) {
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
		displayName = "Viewer"
	case conversions.RoleSpaceViewer:
		displayName = "Space Viewer"
	case conversions.RoleEditor:
		displayName = "Editor"
	case conversions.RoleSpaceEditor:
		displayName = "Space Editor"
	case conversions.RoleFileEditor:
		displayName = "File Editor"
	case conversions.RoleCoowner:
		displayName = "Co Owner"
	case conversions.RoleUploader:
		displayName = "Uploader"
	case conversions.RoleManager:
		displayName = "Manager"
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
	return CS3ResourcePermissionsToLibregraphActions(*role.CS3ResourcePermissions())
}

func GetAllowedResourceActions(role *libregraph.UnifiedRoleDefinition, condition string) []string {
	for _, p := range role.GetRolePermissions() {
		if p.GetCondition() == condition {
			return p.GetAllowedResourceActions()
		}
	}
	return []string{}
}
