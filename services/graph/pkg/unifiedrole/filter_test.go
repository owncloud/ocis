package unifiedrole_test

import (
	"testing"

	. "github.com/onsi/gomega"
	libregraph "github.com/owncloud/libre-graph-api-go"
	"google.golang.org/protobuf/proto"

	"github.com/owncloud/ocis/v2/services/graph/pkg/unifiedrole"
)

func TestRoleFilterIDs(t *testing.T) {
	NewWithT(t).Expect(
		unifiedrole.RoleFilterIDs(
			unifiedrole.UnifiedRoleEditorLiteID,
			unifiedrole.UnifiedRoleSpaceEditorID,
		)(unifiedrole.RoleEditorLite),
	).To(BeTrue())
}

func TestRoleFilterInvert(t *testing.T) {
	NewWithT(t).Expect(
		unifiedrole.RoleFilterInvert(
			unifiedrole.RoleFilterAll(),
		)(unifiedrole.RoleEditorLite),
	).To(BeFalse())
}

func TestRoleFilterAll(t *testing.T) {
	NewWithT(t).Expect(
		unifiedrole.RoleFilterAll()(unifiedrole.RoleEditorLite),
	).To(BeTrue())
}

func TestRoleFilterPermissions(t *testing.T) {
	tests := map[string]struct {
		unifiedRolePermission []libregraph.UnifiedRolePermission
		filterCondition       string
		filterActions         []string
		filterMatch           bool
	}{
		"true | single": {
			unifiedRolePermission: []libregraph.UnifiedRolePermission{
				{
					Condition: proto.String(unifiedrole.UnifiedRoleConditionDrive),
					AllowedResourceActions: []string{
						unifiedrole.DriveItemPermissionsCreate,
					},
				},
			},
			filterCondition: unifiedrole.UnifiedRoleConditionDrive,
			filterActions: []string{
				unifiedrole.DriveItemPermissionsCreate,
			},
			filterMatch: true,
		},
		"true | multiple": {
			unifiedRolePermission: []libregraph.UnifiedRolePermission{
				{
					Condition: proto.String(unifiedrole.UnifiedRoleConditionFolder),
					AllowedResourceActions: []string{
						unifiedrole.DriveItemDeletedRead,
					},
				},
				{
					Condition: proto.String(unifiedrole.UnifiedRoleConditionDrive),
					AllowedResourceActions: []string{
						unifiedrole.DriveItemPermissionsCreate,
					},
				},
			},
			filterCondition: unifiedrole.UnifiedRoleConditionDrive,
			filterActions: []string{
				unifiedrole.DriveItemPermissionsCreate,
			},
			filterMatch: true,
		},
		"false | cross match": {
			unifiedRolePermission: []libregraph.UnifiedRolePermission{
				{
					Condition: proto.String(unifiedrole.UnifiedRoleConditionDrive),
					AllowedResourceActions: []string{
						unifiedrole.DriveItemDeletedRead,
					},
				},
				{
					Condition: proto.String(unifiedrole.UnifiedRoleConditionFolder),
					AllowedResourceActions: []string{
						unifiedrole.DriveItemPermissionsCreate,
					},
				},
			},
			filterCondition: unifiedrole.UnifiedRoleConditionDrive,
			filterActions:   []string{unifiedrole.DriveItemPermissionsCreate},
			filterMatch:     false,
		},
		"false | too many actions": {
			unifiedRolePermission: []libregraph.UnifiedRolePermission{
				{
					Condition: proto.String(unifiedrole.UnifiedRoleConditionDrive),
					AllowedResourceActions: []string{
						unifiedrole.DriveItemDeletedRead,
						unifiedrole.DriveItemPermissionsCreate,
					},
				},
			},
			filterCondition: unifiedrole.UnifiedRoleConditionDrive,
			filterActions: []string{
				unifiedrole.DriveItemPermissionsCreate,
			},
			filterMatch: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			NewWithT(t).Expect(
				unifiedrole.RoleFilterPermission(
					unifiedrole.RoleFilterMatchExact,
					tc.filterCondition,
					tc.filterActions...,
				)(&libregraph.UnifiedRoleDefinition{
					RolePermissions: tc.unifiedRolePermission,
				}),
			).To(Equal(tc.filterMatch))
		})
	}
}
