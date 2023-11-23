package validate_test

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	libregraph "github.com/owncloud/libre-graph-api-go"

	"github.com/owncloud/ocis/v2/services/graph/pkg/unifiedrole"
	"github.com/owncloud/ocis/v2/services/graph/pkg/validate"
)

var _ = Describe("libregraph", func() {

	var driveItemInvite libregraph.DriveItemInvite

	BeforeEach(func() {
		driveItemInvite = libregraph.DriveItemInvite{
			Recipients:                   []libregraph.DriveRecipient{{ObjectId: libregraph.PtrString("1")}},
			Roles:                        []string{unifiedrole.UnifiedRoleEditorID},
			LibreGraphPermissionsActions: []string{unifiedrole.DriveItemVersionsUpdate},
			ExpirationDateTime:           libregraph.PtrTime(time.Now().Add(time.Hour)),
		}
	})

	DescribeTable("DriveItemInvite",
		func(factory func() libregraph.DriveItemInvite, expectError bool) {
			f := factory()
			switch err := validate.StructCtx(context.Background(), f); expectError {
			case true:
				Expect(err).To(HaveOccurred())
			default:
				Expect(err).ToNot(HaveOccurred())
			}

		},
		Entry("succeed: roles", func() libregraph.DriveItemInvite {
			driveItemInvite.LibreGraphPermissionsActions = nil
			return driveItemInvite
		}, false),
		Entry("succeed: permission actions", func() libregraph.DriveItemInvite {
			driveItemInvite.Roles = nil
			return driveItemInvite
		}, false),
		Entry("succeed: without ExpirationDateTime", func() libregraph.DriveItemInvite {
			driveItemInvite.Roles = nil
			driveItemInvite.ExpirationDateTime = nil
			return driveItemInvite
		}, false),
		Entry("fail: multiple role assignment", func() libregraph.DriveItemInvite {
			driveItemInvite.Roles = []string{
				unifiedrole.UnifiedRoleEditorID,
				unifiedrole.UnifiedRoleManagerID,
			}
			driveItemInvite.LibreGraphPermissionsActions = nil
			return driveItemInvite
		}, true),
		Entry("fail: unknown role", func() libregraph.DriveItemInvite {
			driveItemInvite.Roles = []string{"foo"}
			driveItemInvite.LibreGraphPermissionsActions = nil
			return driveItemInvite
		}, true),
		Entry("fail: unknown action", func() libregraph.DriveItemInvite {
			driveItemInvite.Roles = nil
			driveItemInvite.LibreGraphPermissionsActions = []string{"foo"}
			return driveItemInvite
		}, true),
		Entry("fail: missing roles or permission actions", func() libregraph.DriveItemInvite {
			driveItemInvite.Roles = nil
			driveItemInvite.LibreGraphPermissionsActions = nil
			return driveItemInvite
		}, true),
		Entry("fail: different number of roles and actions", func() libregraph.DriveItemInvite {
			driveItemInvite.LibreGraphPermissionsActions = []string{
				unifiedrole.DriveItemVersionsUpdate,
				unifiedrole.DriveItemChildrenCreate,
			}
			return driveItemInvite
		}, true),
		Entry("fail: missing recipients", func() libregraph.DriveItemInvite {
			driveItemInvite.Roles = nil
			driveItemInvite.Recipients = nil
			return driveItemInvite
		}, true),
		Entry("fail: expirationDateTime in the past", func() libregraph.DriveItemInvite {
			driveItemInvite.Roles = nil
			driveItemInvite.ExpirationDateTime = libregraph.PtrTime(time.Now().Add(-time.Hour))
			return driveItemInvite
		}, true),
	)
})
