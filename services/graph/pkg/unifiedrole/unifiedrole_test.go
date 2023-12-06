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

			r := unifiedrole.CS3ResourcePermissionsToUnifiedRole(*cs3perm, constraints, true)
			Expect(r.GetId()).To(Equal(unifiedRole.GetId()))

		},
		Entry(rConversions.RoleViewer, rConversions.NewViewerRole(true), unifiedrole.NewViewerUnifiedRole(true), unifiedrole.UnifiedRoleConditionGrantee),
		Entry(rConversions.RoleEditor, rConversions.NewEditorRole(true), unifiedrole.NewEditorUnifiedRole(true), unifiedrole.UnifiedRoleConditionGrantee),
		Entry(rConversions.RoleFileEditor, rConversions.NewFileEditorRole(true), unifiedrole.NewFileEditorUnifiedRole(true), unifiedrole.UnifiedRoleConditionGrantee),
		Entry(rConversions.RoleCoowner, rConversions.NewCoownerRole(), unifiedrole.NewCoownerUnifiedRole(), unifiedrole.UnifiedRoleConditionGrantee),
		Entry(rConversions.RoleManager, rConversions.NewManagerRole(), unifiedrole.NewManagerUnifiedRole(), unifiedrole.UnifiedRoleConditionGrantee),
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
		Entry(rConversions.RoleViewer, rConversions.NewViewerRole(true), unifiedrole.NewViewerUnifiedRole(true), true),
		Entry(rConversions.RoleEditor, rConversions.NewEditorRole(true), unifiedrole.NewEditorUnifiedRole(true), true),
		Entry(rConversions.RoleFileEditor, rConversions.NewFileEditorRole(true), unifiedrole.NewFileEditorUnifiedRole(true), true),
		Entry(rConversions.RoleCoowner, rConversions.NewCoownerRole(), unifiedrole.NewCoownerUnifiedRole(), true),
		Entry(rConversions.RoleManager, rConversions.NewManagerRole(), unifiedrole.NewManagerUnifiedRole(), true),
		Entry("no match", rConversions.NewFileEditorRole(true), unifiedrole.NewManagerUnifiedRole(), false),
	)

	DescribeTable("WeightRoleDefinitions",
		func(roleDefinitions []*libregraph.UnifiedRoleDefinition, descending bool, expectedDefinitions []*libregraph.UnifiedRoleDefinition) {

			for i, generatedDefinition := range unifiedrole.WeightRoleDefinitions(roleDefinitions, descending) {
				Expect(generatedDefinition.Id).To(Equal(expectedDefinitions[i].Id))
			}
		},

		Entry("ascending",
			[]*libregraph.UnifiedRoleDefinition{
				unifiedrole.NewViewerUnifiedRole(false),
				unifiedrole.NewFileEditorUnifiedRole(false),
			},
			false,
			[]*libregraph.UnifiedRoleDefinition{
				unifiedrole.NewViewerUnifiedRole(false),
				unifiedrole.NewFileEditorUnifiedRole(false),
			},
		),

		Entry("descending",
			[]*libregraph.UnifiedRoleDefinition{
				unifiedrole.NewViewerUnifiedRole(false),
				unifiedrole.NewFileEditorUnifiedRole(false),
			},
			true,
			[]*libregraph.UnifiedRoleDefinition{
				unifiedrole.NewFileEditorUnifiedRole(false),
				unifiedrole.NewViewerUnifiedRole(false),
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
			func(givenActions []string, constraints string, resharing bool, expectedDefinitions []*libregraph.UnifiedRoleDefinition) {

				generatedDefinitions := unifiedrole.GetApplicableRoleDefinitionsForActions(givenActions, constraints, resharing, false)

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
				rolesToAction(unifiedrole.NewViewerUnifiedRole(false)),
				unifiedrole.UnifiedRoleConditionGrantee,
				false,
				[]*libregraph.UnifiedRoleDefinition{
					unifiedrole.NewViewerUnifiedRole(false),
				},
			),

			Entry(
				"ViewerUnifiedRole | share",
				rolesToAction(unifiedrole.NewViewerUnifiedRole(true)),
				unifiedrole.UnifiedRoleConditionGrantee,
				true,
				[]*libregraph.UnifiedRoleDefinition{
					unifiedrole.NewViewerUnifiedRole(true),
				},
			),

			Entry(
				"NewFileEditorUnifiedRole",
				rolesToAction(unifiedrole.NewFileEditorUnifiedRole(false)),
				unifiedrole.UnifiedRoleConditionGrantee,
				false,
				[]*libregraph.UnifiedRoleDefinition{
					unifiedrole.NewViewerUnifiedRole(false),
					unifiedrole.NewFileEditorUnifiedRole(false),
				},
			),

			Entry(
				"NewFileEditorUnifiedRole - share",
				rolesToAction(unifiedrole.NewFileEditorUnifiedRole(true)),
				unifiedrole.UnifiedRoleConditionGrantee,
				true,
				[]*libregraph.UnifiedRoleDefinition{
					unifiedrole.NewViewerUnifiedRole(true),
					unifiedrole.NewFileEditorUnifiedRole(true),
				},
			),

			Entry(
				"NewEditorUnifiedRole",
				rolesToAction(unifiedrole.NewEditorUnifiedRole(false)),
				unifiedrole.UnifiedRoleConditionGrantee,
				false,
				[]*libregraph.UnifiedRoleDefinition{
					unifiedrole.NewUploaderUnifiedRole(),
					unifiedrole.NewViewerUnifiedRole(false),
					unifiedrole.NewFileEditorUnifiedRole(false),
					unifiedrole.NewEditorUnifiedRole(false),
				},
			),

			Entry(
				"NewEditorUnifiedRole - share",
				rolesToAction(unifiedrole.NewEditorUnifiedRole(true)),
				unifiedrole.UnifiedRoleConditionGrantee,
				true,
				[]*libregraph.UnifiedRoleDefinition{
					unifiedrole.NewUploaderUnifiedRole(),
					unifiedrole.NewViewerUnifiedRole(true),
					unifiedrole.NewFileEditorUnifiedRole(true),
					unifiedrole.NewEditorUnifiedRole(true),
				},
			),

			Entry(
				"GetBuiltinRoleDefinitionList",
				rolesToAction(unifiedrole.GetBuiltinRoleDefinitionList(false)...),
				unifiedrole.UnifiedRoleConditionGrantee,
				false,
				[]*libregraph.UnifiedRoleDefinition{
					unifiedrole.NewUploaderUnifiedRole(),
					unifiedrole.NewViewerUnifiedRole(false),
					unifiedrole.NewFileEditorUnifiedRole(false),
					unifiedrole.NewEditorUnifiedRole(false),
					unifiedrole.NewCoownerUnifiedRole(),
					unifiedrole.NewManagerUnifiedRole(),
				},
			),

			Entry(
				"GetBuiltinRoleDefinitionList - share",
				rolesToAction(unifiedrole.GetBuiltinRoleDefinitionList(true)...),
				unifiedrole.UnifiedRoleConditionGrantee,
				true,
				[]*libregraph.UnifiedRoleDefinition{
					unifiedrole.NewUploaderUnifiedRole(),
					unifiedrole.NewViewerUnifiedRole(true),
					unifiedrole.NewFileEditorUnifiedRole(true),
					unifiedrole.NewEditorUnifiedRole(true),
					unifiedrole.NewCoownerUnifiedRole(),
					unifiedrole.NewManagerUnifiedRole(),
				},
			),

			Entry(
				"single",
				[]string{unifiedrole.DriveItemQuotaRead},
				unifiedrole.UnifiedRoleConditionGrantee,
				true,
				[]*libregraph.UnifiedRoleDefinition{},
			),

			Entry(
				"mixed",
				append(rolesToAction(unifiedrole.NewUploaderUnifiedRole()), unifiedrole.DriveItemQuotaRead),
				unifiedrole.UnifiedRoleConditionGrantee,
				true,
				[]*libregraph.UnifiedRoleDefinition{
					unifiedrole.NewUploaderUnifiedRole(),
				},
			),
		)
	}

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
