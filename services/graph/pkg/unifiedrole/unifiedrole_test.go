package unifiedrole_test

import (
	"fmt"
	"slices"

	rConversions "github.com/cs3org/reva/v2/pkg/conversions"
	"github.com/owncloud/ocis/v2/ocis-pkg/conversions"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
	libregraph "github.com/owncloud/libre-graph-api-go"

	"github.com/owncloud/ocis/v2/services/graph/pkg/unifiedrole"
)

var _ = Describe("unifiedroles", func() {
	DescribeTable("CS3ResourcePermissionsToUnifiedRole",
		func(legacyRole *rConversions.Role, unifiedRole *libregraph.UnifiedRoleDefinition, constraints string) {
			cs3perm := legacyRole.CS3ResourcePermissions()

			r := unifiedrole.CS3ResourcePermissionsToUnifiedRole(*cs3perm, constraints)
			Expect(r.GetId()).To(Equal(unifiedRole.GetId()))

		},
		Entry(rConversions.RoleViewer, rConversions.NewViewerRole(), unifiedrole.NewViewerUnifiedRole(), unifiedrole.UnifiedRoleConditionGrantee),
		Entry(rConversions.RoleEditor, rConversions.NewEditorRole(), unifiedrole.NewEditorUnifiedRole(), unifiedrole.UnifiedRoleConditionGrantee),
		Entry(rConversions.RoleFileEditor, rConversions.NewFileEditorRole(), unifiedrole.NewFileEditorUnifiedRole(), unifiedrole.UnifiedRoleConditionGrantee),
		Entry(rConversions.RoleManager, rConversions.NewManagerRole(), unifiedrole.NewManagerUnifiedRole(), unifiedrole.UnifiedRoleConditionOwner),
		Entry(rConversions.RoleSpaceViewer, rConversions.NewSpaceViewerRole(), unifiedrole.NewSpaceViewerUnifiedRole(), unifiedrole.UnifiedRoleConditionOwner),
		Entry(rConversions.RoleSpaceEditor, rConversions.NewSpaceEditorRole(), unifiedrole.NewSpaceEditorUnifiedRole(), unifiedrole.UnifiedRoleConditionOwner),
	)

	DescribeTable("UnifiedRolePermissionsToCS3ResourcePermissions",
		func(cs3Role *rConversions.Role, libregraphRole *libregraph.UnifiedRoleDefinition, match bool) {
			permsFromCS3 := cs3Role.CS3ResourcePermissions()
			permsFromUnifiedRole := unifiedrole.PermissionsToCS3ResourcePermissions(
				conversions.ToPointerSlice(libregraphRole.RolePermissions),
			)

			var matcher types.GomegaMatcher

			if match {
				matcher = Equal(permsFromUnifiedRole)
			} else {
				matcher = Not(Equal(permsFromUnifiedRole))
			}

			Expect(permsFromCS3).To(matcher)
		},
		Entry(rConversions.RoleViewer, rConversions.NewViewerRole(), unifiedrole.NewViewerUnifiedRole(), true),
		Entry(rConversions.RoleEditor, rConversions.NewEditorRole(), unifiedrole.NewEditorUnifiedRole(), true),
		Entry(rConversions.RoleFileEditor, rConversions.NewFileEditorRole(), unifiedrole.NewFileEditorUnifiedRole(), true),
		Entry(rConversions.RoleManager, rConversions.NewManagerRole(), unifiedrole.NewManagerUnifiedRole(), true),
		Entry("no match", rConversions.NewFileEditorRole(), unifiedrole.NewManagerUnifiedRole(), false),
	)

	DescribeTable("WeightRoleDefinitions",
		func(roleDefinitions []*libregraph.UnifiedRoleDefinition, descending bool, expectedDefinitions []*libregraph.UnifiedRoleDefinition) {

			for i, generatedDefinition := range unifiedrole.WeightRoleDefinitions(roleDefinitions, descending) {
				Expect(generatedDefinition.Id).To(Equal(expectedDefinitions[i].Id))
			}
		},

		Entry("ascending",
			[]*libregraph.UnifiedRoleDefinition{
				unifiedrole.NewViewerUnifiedRole(),
				unifiedrole.NewFileEditorUnifiedRole(),
			},
			false,
			[]*libregraph.UnifiedRoleDefinition{
				unifiedrole.NewViewerUnifiedRole(),
				unifiedrole.NewFileEditorUnifiedRole(),
			},
		),

		Entry("descending",
			[]*libregraph.UnifiedRoleDefinition{
				unifiedrole.NewViewerUnifiedRole(),
				unifiedrole.NewFileEditorUnifiedRole(),
			},
			true,
			[]*libregraph.UnifiedRoleDefinition{
				unifiedrole.NewFileEditorUnifiedRole(),
				unifiedrole.NewViewerUnifiedRole(),
			},
		),
	)

	{
		rolesToAction := func(definitions ...*libregraph.UnifiedRoleDefinition) []string {
			var actions []string

			for _, definition := range definitions {
				for _, permission := range definition.GetRolePermissions() {
					for _, action := range permission.GetAllowedResourceActions() {
						if slices.Contains(actions, action) {
							continue
						}
						actions = append(actions, action)
					}
				}
			}

			return actions
		}

		DescribeTable("GetApplicableRoleDefinitionsForActions",
			func(givenActions []string, constraints string, expectedDefinitions []*libregraph.UnifiedRoleDefinition) {

				generatedDefinitions := unifiedrole.GetApplicableRoleDefinitionsForActions(givenActions, constraints, false)

				Expect(len(generatedDefinitions)).To(Equal(len(expectedDefinitions)))

				for i, generatedDefinition := range generatedDefinitions {
					Expect(generatedDefinition.Id).To(Equal(expectedDefinitions[i].Id))
					Expect(*generatedDefinition.LibreGraphWeight).To(Equal(int32(i + 1)))
				}

				generatedActions := rolesToAction(generatedDefinitions...)
				Expect(len(givenActions) >= len(generatedActions)).To(BeTrue())

				for _, generatedAction := range generatedActions {
					Expect(slices.Contains(givenActions, generatedAction)).To(BeTrue())
				}
			},

			Entry(
				"ViewerUnifiedRole",
				rolesToAction(unifiedrole.NewViewerUnifiedRole()),
				unifiedrole.UnifiedRoleConditionGrantee,
				[]*libregraph.UnifiedRoleDefinition{
					unifiedrole.NewViewerUnifiedRole(),
				},
			),

			Entry(
				"ViewerUnifiedRole | share",
				rolesToAction(unifiedrole.NewViewerUnifiedRole()),
				unifiedrole.UnifiedRoleConditionGrantee,
				[]*libregraph.UnifiedRoleDefinition{
					unifiedrole.NewViewerUnifiedRole(),
				},
			),

			Entry(
				"NewFileEditorUnifiedRole",
				rolesToAction(unifiedrole.NewFileEditorUnifiedRole()),
				unifiedrole.UnifiedRoleConditionGrantee,
				[]*libregraph.UnifiedRoleDefinition{
					unifiedrole.NewViewerUnifiedRole(),
					unifiedrole.NewFileEditorUnifiedRole(),
				},
			),

			Entry(
				"NewFileEditorUnifiedRole - share",
				rolesToAction(unifiedrole.NewFileEditorUnifiedRole()),
				unifiedrole.UnifiedRoleConditionGrantee,
				[]*libregraph.UnifiedRoleDefinition{
					unifiedrole.NewViewerUnifiedRole(),
					unifiedrole.NewFileEditorUnifiedRole(),
				},
			),

			Entry(
				"NewEditorUnifiedRole",
				rolesToAction(unifiedrole.NewEditorUnifiedRole()),
				unifiedrole.UnifiedRoleConditionGrantee,
				[]*libregraph.UnifiedRoleDefinition{
					unifiedrole.NewUploaderUnifiedRole(),
					unifiedrole.NewViewerUnifiedRole(),
					unifiedrole.NewFileEditorUnifiedRole(),
					unifiedrole.NewEditorUnifiedRole(),
				},
			),

			Entry(
				"NewEditorUnifiedRole - share",
				rolesToAction(unifiedrole.NewEditorUnifiedRole()),
				unifiedrole.UnifiedRoleConditionGrantee,
				[]*libregraph.UnifiedRoleDefinition{
					unifiedrole.NewUploaderUnifiedRole(),
					unifiedrole.NewViewerUnifiedRole(),
					unifiedrole.NewFileEditorUnifiedRole(),
					unifiedrole.NewEditorUnifiedRole(),
				},
			),

			Entry(
				"GetBuiltinRoleDefinitionList",
				rolesToAction(unifiedrole.GetBuiltinRoleDefinitionList()...),
				unifiedrole.UnifiedRoleConditionGrantee,
				[]*libregraph.UnifiedRoleDefinition{
					unifiedrole.NewUploaderUnifiedRole(),
					unifiedrole.NewViewerUnifiedRole(),
					unifiedrole.NewFileEditorUnifiedRole(),
					unifiedrole.NewEditorUnifiedRole(),
				},
			),

			Entry(
				"GetBuiltinRoleDefinitionList",
				rolesToAction(unifiedrole.GetBuiltinRoleDefinitionList()...),
				unifiedrole.UnifiedRoleConditionOwner,
				[]*libregraph.UnifiedRoleDefinition{
					unifiedrole.NewSpaceViewerUnifiedRole(),
					unifiedrole.NewSpaceEditorUnifiedRole(),
					unifiedrole.NewManagerUnifiedRole(),
				},
			),

			Entry(
				"GetBuiltinRoleDefinitionList - share",
				rolesToAction(unifiedrole.GetBuiltinRoleDefinitionList()...),
				unifiedrole.UnifiedRoleConditionGrantee,
				[]*libregraph.UnifiedRoleDefinition{
					unifiedrole.NewUploaderUnifiedRole(),
					unifiedrole.NewViewerUnifiedRole(),
					unifiedrole.NewFileEditorUnifiedRole(),
					unifiedrole.NewEditorUnifiedRole(),
				},
			),

			Entry(
				"single",
				[]string{unifiedrole.DriveItemQuotaRead},
				unifiedrole.UnifiedRoleConditionGrantee,
				[]*libregraph.UnifiedRoleDefinition{},
			),

			Entry(
				"mixed",
				append(rolesToAction(unifiedrole.NewUploaderUnifiedRole()), unifiedrole.DriveItemQuotaRead),
				unifiedrole.UnifiedRoleConditionGrantee,
				[]*libregraph.UnifiedRoleDefinition{
					unifiedrole.NewUploaderUnifiedRole(),
				},
			),
		)
	}

	{
		var newUnifiedRoleFromIDEntries []TableEntry
		attachEntry := func(name, id string, definition *libregraph.UnifiedRoleDefinition, errors bool) {
			e := Entry(
				fmt.Sprintf("%s", name),
				id,
				definition,
				errors,
			)

			newUnifiedRoleFromIDEntries = append(newUnifiedRoleFromIDEntries, e)
		}

		for _, definition := range unifiedrole.GetBuiltinRoleDefinitionList() {
			attachEntry(definition.GetDisplayName(), definition.GetId(), definition, false)
		}

		attachEntry("unknown", "123", nil, true)

		DescribeTable("NewUnifiedRoleFromID",
			func(id string, expectedRole *libregraph.UnifiedRoleDefinition, expectError bool) {
				role, err := unifiedrole.NewUnifiedRoleFromID(id)

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
