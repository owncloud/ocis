package unifiedrole_test

import (
	"github.com/cs3org/reva/v2/pkg/conversions"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	libregraph "github.com/owncloud/libre-graph-api-go"
	"github.com/owncloud/ocis/v2/services/graph/pkg/unifiedrole"
)

var _ = Describe("unifiedroles", func() {
	DescribeTable("CS3ResourcePermissionsToUnifiedRole",
		func(legacyRole *conversions.Role, unifiedRole *libregraph.UnifiedRoleDefinition) {
			cs3perm := legacyRole.CS3ResourcePermissions()

			r := unifiedrole.CS3ResourcePermissionsToUnifiedRole(*cs3perm, unifiedrole.UnifiedRoleConditionGrantee, true)
			Expect(r.GetId()).To(Equal(unifiedRole.GetId()))

		},
		Entry(conversions.RoleViewer, conversions.NewViewerRole(true), unifiedrole.NewViewerUnifiedRole(true)),
		Entry(conversions.RoleEditor, conversions.NewEditorRole(true), unifiedrole.NewEditorUnifiedRole(true)),
		Entry(conversions.RoleFileEditor, conversions.NewFileEditorRole(true), unifiedrole.NewFileEditorUnifiedRole(true)),
		Entry(conversions.RoleCoowner, conversions.NewCoownerRole(), unifiedrole.NewCoownerUnifiedRole()),
		Entry(conversions.RoleManager, conversions.NewManagerRole(), unifiedrole.NewManagerUnifiedRole()),
	)
})
