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

type validatableFactory[T any] func() (T, bool)

var _ = Describe("libregraph", func() {
	var ctx context.Context
	var driveItemInvite libregraph.DriveItemInvite
	var driveRecipient libregraph.DriveRecipient

	BeforeEach(func() {
		ctx = context.Background()
		driveRecipient = libregraph.DriveRecipient{
			ObjectId:                new("1"),
			LibreGraphRecipientType: new("user"),
		}

		driveItemInvite = libregraph.DriveItemInvite{
			Recipients:                   []libregraph.DriveRecipient{driveRecipient},
			Roles:                        []string{unifiedrole.UnifiedRoleEditorID},
			LibreGraphPermissionsActions: []string{unifiedrole.DriveItemVersionsUpdate},
			ExpirationDateTime:           new(time.Now().Add(time.Hour)),
		}
	})

	DescribeTable("DriveItemInvite",
		func(factories ...validatableFactory[libregraph.DriveItemInvite]) {
			for _, factory := range factories {
				s, pass := factory()
				switch err := validate.StructCtx(ctx, s); pass {
				case false:
					Expect(err).To(HaveOccurred())
				default:
					Expect(err).ToNot(HaveOccurred())
				}
			}
		},
		Entry("succeed: roles only", func() (libregraph.DriveItemInvite, bool) {
			driveItemInvite.LibreGraphPermissionsActions = nil
			return driveItemInvite, true
		}),
		Entry("succeed: actions only", func() (libregraph.DriveItemInvite, bool) {
			driveItemInvite.Roles = nil
			return driveItemInvite, true
		}),
		Entry("succeed: no ExpirationDateTime", func() (libregraph.DriveItemInvite, bool) {
			driveItemInvite.Roles = nil
			driveItemInvite.ExpirationDateTime = nil
			return driveItemInvite, true
		}),
		Entry("fail: multiple role assignment", func() (libregraph.DriveItemInvite, bool) {
			driveItemInvite.Roles = []string{
				unifiedrole.UnifiedRoleEditorID,
				unifiedrole.UnifiedRoleManagerID,
			}
			driveItemInvite.LibreGraphPermissionsActions = nil
			return driveItemInvite, false
		}),
		Entry("fail: unknown role", func() (libregraph.DriveItemInvite, bool) {
			driveItemInvite.Roles = []string{"foo"}
			driveItemInvite.LibreGraphPermissionsActions = nil
			return driveItemInvite, false
		}),
		Entry("fail: disabled role", func() (libregraph.DriveItemInvite, bool) {
			ctx = validate.ContextWithAllowedRoleIDs(ctx, []string{unifiedrole.UnifiedRoleEditorID})
			driveItemInvite.Roles = []string{
				unifiedrole.UnifiedRoleSecureViewerID,
			}
			driveItemInvite.LibreGraphPermissionsActions = nil
			return driveItemInvite, false
		}),
		Entry("fail: unknown action", func() (libregraph.DriveItemInvite, bool) {
			driveItemInvite.Roles = nil
			driveItemInvite.LibreGraphPermissionsActions = []string{"foo"}
			return driveItemInvite, false
		}),
		Entry("fail: missing roles or permission actions", func() (libregraph.DriveItemInvite, bool) {
			driveItemInvite.Roles = nil
			driveItemInvite.LibreGraphPermissionsActions = nil
			return driveItemInvite, false
		}),
		Entry("fail: different number of roles and actions", func() (libregraph.DriveItemInvite, bool) {
			driveItemInvite.LibreGraphPermissionsActions = []string{
				unifiedrole.DriveItemVersionsUpdate,
				unifiedrole.DriveItemChildrenCreate,
			}
			return driveItemInvite, false
		}),
		Entry("fail: missing recipients", func() (libregraph.DriveItemInvite, bool) {
			driveItemInvite.Roles = nil
			driveItemInvite.Recipients = nil
			return driveItemInvite, false
		}),
		Entry("fail: dive recipients", func() (libregraph.DriveItemInvite, bool) {
			driveItemInvite.Roles = nil
			driveItemInvite.Recipients = []libregraph.DriveRecipient{
				{ObjectId: new("1"), LibreGraphRecipientType: new("dive")},
			}
			return driveItemInvite, false
		}),
		Entry("fail: more than 1 recipient", func() (libregraph.DriveItemInvite, bool) {
			driveItemInvite.Roles = nil
			driveItemInvite.Recipients = []libregraph.DriveRecipient{
				driveRecipient,
				{ObjectId: new("2"), LibreGraphRecipientType: new("group")},
			}
			return driveItemInvite, false
		}),
		Entry("fail: expirationDateTime in the past", func() (libregraph.DriveItemInvite, bool) {
			driveItemInvite.Roles = nil
			driveItemInvite.ExpirationDateTime = new(time.Now().Add(-time.Hour))
			return driveItemInvite, false
		}),
	)

	// Permission is validated on the update paths (PATCH). It must honour the
	// admin role allowlist just like DriveItemInvite does.
	DescribeTable("Permission",
		func(factory validatableFactory[libregraph.Permission]) {
			s, pass := factory()
			switch err := validate.StructCtx(ctx, s); pass {
			case false:
				Expect(err).To(HaveOccurred())
			default:
				Expect(err).ToNot(HaveOccurred())
			}
		},
		Entry("fail: disabled role rejected when allowlist is set", func() (libregraph.Permission, bool) {
			ctx = validate.ContextWithAllowedRoleIDs(ctx, []string{unifiedrole.UnifiedRoleEditorID})
			return libregraph.Permission{Roles: []string{unifiedrole.UnifiedRoleSecureViewerID}}, false
		}),
		Entry("succeed: allowed role accepted when allowlist is set", func() (libregraph.Permission, bool) {
			ctx = validate.ContextWithAllowedRoleIDs(ctx, []string{unifiedrole.UnifiedRoleEditorID})
			return libregraph.Permission{Roles: []string{unifiedrole.UnifiedRoleEditorID}}, true
		}),
	)

	DescribeTable("DriveRecipient",
		func(factories ...validatableFactory[libregraph.DriveRecipient]) {
			for _, factory := range factories {
				s, pass := factory()
				switch err := validate.StructCtx(context.Background(), s); pass {
				case false:
					Expect(err).To(HaveOccurred())
				default:
					Expect(err).ToNot(HaveOccurred())
				}
			}
		},
		Entry("fail: invalid objectId",
			func() (libregraph.DriveRecipient, bool) {
				driveRecipient.ObjectId = nil
				return driveRecipient, false
			},
			func() (libregraph.DriveRecipient, bool) {
				driveRecipient.ObjectId = new("")
				return driveRecipient, false
			},
		),
		Entry("succeed: valid role",
			func() (libregraph.DriveRecipient, bool) {
				driveRecipient.LibreGraphRecipientType = new("user")
				return driveRecipient, true
			},
			func() (libregraph.DriveRecipient, bool) {
				driveRecipient.LibreGraphRecipientType = new("group")
				return driveRecipient, true
			},
		),
		Entry("fail: invalid role",
			func() (libregraph.DriveRecipient, bool) {
				driveRecipient.LibreGraphRecipientType = new("foo")
				return driveRecipient, false
			},
			func() (libregraph.DriveRecipient, bool) {
				driveRecipient.LibreGraphRecipientType = nil
				return driveRecipient, false
			},
		),
	)
})
