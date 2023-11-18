package unifiedrole_test

import (
	"fmt"

	"github.com/cs3org/reva/v2/pkg/conversions"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
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

	DescribeTable("UnifiedRolePermissionsToCS3ResourcePermissions",
		func(cs3Role *conversions.Role, libregraphRole *libregraph.UnifiedRoleDefinition, match bool) {
			permsFromCS3 := cs3Role.CS3ResourcePermissions()
			permsFromUnifiedRole := unifiedrole.PermissionsToCS3ResourcePermissions(libregraphRole.RolePermissions)

			var matcher types.GomegaMatcher

			if match {
				matcher = Equal(permsFromUnifiedRole)
			} else {
				matcher = Not(Equal(permsFromUnifiedRole))
			}

			Expect(permsFromCS3).To(matcher)
		},
		Entry(conversions.RoleViewer, conversions.NewViewerRole(true), unifiedrole.NewViewerUnifiedRole(true), true),
		Entry(conversions.RoleEditor, conversions.NewEditorRole(true), unifiedrole.NewEditorUnifiedRole(true), true),
		Entry(conversions.RoleFileEditor, conversions.NewFileEditorRole(true), unifiedrole.NewFileEditorUnifiedRole(true), true),
		Entry(conversions.RoleCoowner, conversions.NewCoownerRole(), unifiedrole.NewCoownerUnifiedRole(), true),
		Entry(conversions.RoleManager, conversions.NewManagerRole(), unifiedrole.NewManagerUnifiedRole(), true),
		Entry("no match", conversions.NewFileEditorRole(true), unifiedrole.NewManagerUnifiedRole(), false),
	)

	{
		var newUnifiedRoleFromIDEntries []TableEntry
		for _, resharing := range []bool{true, false} {
			attachEntry := func(name, id string, definition *libregraph.UnifiedRoleDefinition, errors bool) {
				e := Entry(
					fmt.Sprintf("%s - resharing: %t", name, resharing),
					id,
					resharing,
					definition,
					errors,
				)

				newUnifiedRoleFromIDEntries = append(newUnifiedRoleFromIDEntries, e)
			}

			for _, definition := range unifiedrole.GetBuiltinRoleDefinitionList(resharing) {
				attachEntry(definition.GetDisplayName(), definition.GetId(), definition, false)
			}

			attachEntry("unknown", "123", nil, true)
		}

		DescribeTable("NewUnifiedRoleFromID",
			func(id string, resharing bool, expectedRole *libregraph.UnifiedRoleDefinition, expectError bool) {
				role, err := unifiedrole.NewUnifiedRoleFromID(id, resharing)

				if expectError {
					Expect(err).To(HaveOccurred())
				} else {
					Expect(err).NotTo(HaveOccurred())
					Expect(role).To(Equal(expectedRole))
				}
			},
			newUnifiedRoleFromIDEntries,
		)
	}
})
