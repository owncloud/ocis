package connector_test

import (
	"context"
	"encoding/hex"
	"errors"

	appproviderv1beta1 "github.com/cs3org/go-cs3apis/cs3/app/provider/v1beta1"
	userv1beta1 "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	providerv1beta1 "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	typesv1beta1 "github.com/cs3org/go-cs3apis/cs3/types/v1beta1"
	"github.com/cs3org/reva/v2/pkg/rgrpc/status"
	cs3mocks "github.com/cs3org/reva/v2/tests/cs3mocks/mocks"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/owncloud/ocis/v2/services/collaboration/pkg/config"
	"github.com/owncloud/ocis/v2/services/collaboration/pkg/connector"
	"github.com/owncloud/ocis/v2/services/collaboration/pkg/connector/fileinfo"
	"github.com/owncloud/ocis/v2/services/collaboration/pkg/middleware"
	"github.com/stretchr/testify/mock"
)

var _ = Describe("FileConnector", func() {
	var (
		fc            *connector.FileConnector
		gatewayClient *cs3mocks.GatewayAPIClient
		cfg           *config.Config
		wopiCtx       middleware.WopiContext
	)

	BeforeEach(func() {
		cfg = &config.Config{
			App: config.App{
				LockName: "testName_for_unittests", // Only the LockName is used
			},
		}
		gatewayClient = &cs3mocks.GatewayAPIClient{}
		fc = connector.NewFileConnector(gatewayClient, cfg)

		wopiCtx = middleware.WopiContext{
			AccessToken: "abcdef123456",
			FileReference: providerv1beta1.Reference{
				ResourceId: &providerv1beta1.ResourceId{
					StorageId: "abc",
					OpaqueId:  "12345",
					SpaceId:   "zzz",
				},
				Path: ".",
			},
			User: &userv1beta1.User{
				Id: &userv1beta1.UserId{
					Idp:      "inmemory",
					OpaqueId: "opaqueId",
					Type:     userv1beta1.UserType_USER_TYPE_PRIMARY,
				},
				Username:    "Shaft",
				DisplayName: "Pet Shaft",
				Mail:        "shaft@example.com",
				// Opaque is here for reference, not used by default but might be needed for some tests
				//Opaque: &typesv1beta1.Opaque{
				//	Map: map[string]*typesv1beta1.OpaqueEntry{
				//		"public-share-role": &typesv1beta1.OpaqueEntry{
				//			Decoder: "plain",
				//			Value:   []byte("viewer"),
				//		},
				//	},
				//},
			},
			ViewMode:   appproviderv1beta1.ViewMode_VIEW_MODE_READ_WRITE,
			EditAppUrl: "http://test.ex.prv/edit",
			ViewAppUrl: "http://test.ex.prv/view",
		}
	})

	Describe("GetLock", func() {
		It("No valid context", func() {
			ctx := context.Background()
			newLockId, err := fc.GetLock(ctx)
			Expect(err).To(HaveOccurred())
			Expect(newLockId).To(Equal(""))
		})

		It("Get lock failed", func() {
			ctx := middleware.WopiContextToCtx(context.Background(), wopiCtx)

			targetErr := errors.New("Something went wrong")
			gatewayClient.On("GetLock", mock.Anything, mock.Anything).Times(1).Return(&providerv1beta1.GetLockResponse{
				Status: status.NewInternal(ctx, "Something failed"),
			}, targetErr)

			newLockId, err := fc.GetLock(ctx)
			Expect(err).To(Equal(targetErr))
			Expect(newLockId).To(Equal(""))
		})

		It("Get lock failed status not ok", func() {
			// assume failure happens because the target file doesn't exists
			ctx := middleware.WopiContextToCtx(context.Background(), wopiCtx)

			gatewayClient.On("GetLock", mock.Anything, mock.Anything).Times(1).Return(&providerv1beta1.GetLockResponse{
				Status: status.NewNotFound(ctx, "File is missing"),
			}, nil)

			newLockId, err := fc.GetLock(ctx)
			Expect(err).To(HaveOccurred())
			conErr := err.(*connector.ConnectorError)
			Expect(conErr.HttpCodeOut).To(Equal(404))
			Expect(newLockId).To(Equal(""))
		})

		It("Get lock success", func() {
			// assume failure happens because the target file doesn't exists
			ctx := middleware.WopiContextToCtx(context.Background(), wopiCtx)

			gatewayClient.On("GetLock", mock.Anything, mock.Anything).Times(1).Return(&providerv1beta1.GetLockResponse{
				Status: status.NewOK(ctx),
				Lock: &providerv1beta1.Lock{
					LockId: "zzz999",
					Type:   providerv1beta1.LockType_LOCK_TYPE_WRITE,
				},
			}, nil)

			newLockId, err := fc.GetLock(ctx)
			Expect(err).To(Succeed())
			Expect(newLockId).To(Equal("zzz999"))
		})
	})

	Describe("Lock", func() {
		Describe("Lock", func() {
			It("No valid context", func() {
				ctx := context.Background()
				newLockId, err := fc.Lock(ctx, "newLock", "")
				Expect(err).To(HaveOccurred())
				Expect(newLockId).To(Equal(""))
			})

			It("Empty lockId", func() {
				ctx := middleware.WopiContextToCtx(context.Background(), wopiCtx)

				newLockId, err := fc.Lock(ctx, "", "")
				Expect(err).To(HaveOccurred())
				conErr := err.(*connector.ConnectorError)
				Expect(conErr.HttpCodeOut).To(Equal(400))
				Expect(newLockId).To(Equal(""))
			})

			It("Set lock failed", func() {
				ctx := middleware.WopiContextToCtx(context.Background(), wopiCtx)

				targetErr := errors.New("Something went wrong")
				gatewayClient.On("SetLock", mock.Anything, mock.Anything).Times(1).Return(&providerv1beta1.SetLockResponse{
					Status: status.NewInternal(ctx, "Something failed"),
				}, targetErr)

				newLockId, err := fc.Lock(ctx, "abcdef123", "")
				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(targetErr))
				Expect(newLockId).To(Equal(""))
			})

			It("Set lock success", func() {
				ctx := middleware.WopiContextToCtx(context.Background(), wopiCtx)

				gatewayClient.On("SetLock", mock.Anything, mock.Anything).Times(1).Return(&providerv1beta1.SetLockResponse{
					Status: status.NewOK(ctx),
				}, nil)

				newLockId, err := fc.Lock(ctx, "abcdef123", "")
				Expect(err).To(Succeed())
				Expect(newLockId).To(Equal(""))
			})

			It("Set lock mismatches error getting lock", func() {
				ctx := middleware.WopiContextToCtx(context.Background(), wopiCtx)

				gatewayClient.On("SetLock", mock.Anything, mock.Anything).Times(1).Return(&providerv1beta1.SetLockResponse{
					Status: status.NewFailedPrecondition(ctx, nil, "lock mismatch"),
				}, nil)

				targetErr := errors.New("Something went wrong getting the lock")
				gatewayClient.On("GetLock", mock.Anything, mock.Anything).Times(1).Return(&providerv1beta1.GetLockResponse{
					Status: status.NewInternal(ctx, "lock mismatch"),
				}, targetErr)

				newLockId, err := fc.Lock(ctx, "abcdef123", "")
				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(targetErr))
				Expect(newLockId).To(Equal(""))
			})

			It("Set lock mismatches", func() {
				ctx := middleware.WopiContextToCtx(context.Background(), wopiCtx)

				gatewayClient.On("SetLock", mock.Anything, mock.Anything).Times(1).Return(&providerv1beta1.SetLockResponse{
					Status: status.NewFailedPrecondition(ctx, nil, "lock mismatch"),
				}, nil)

				gatewayClient.On("GetLock", mock.Anything, mock.Anything).Times(1).Return(&providerv1beta1.GetLockResponse{
					Status: status.NewOK(ctx),
					Lock: &providerv1beta1.Lock{
						LockId: "zzz999",
						Type:   providerv1beta1.LockType_LOCK_TYPE_WRITE,
					},
				}, nil)

				newLockId, err := fc.Lock(ctx, "abcdef123", "")
				Expect(err).To(HaveOccurred())
				conErr := err.(*connector.ConnectorError)
				Expect(conErr.HttpCodeOut).To(Equal(409))
				Expect(newLockId).To(Equal("zzz999"))
			})

			It("Set lock mismatches but get lock matches", func() {
				ctx := middleware.WopiContextToCtx(context.Background(), wopiCtx)

				gatewayClient.On("SetLock", mock.Anything, mock.Anything).Times(1).Return(&providerv1beta1.SetLockResponse{
					Status: status.NewFailedPrecondition(ctx, nil, "lock mismatch"),
				}, nil)

				gatewayClient.On("GetLock", mock.Anything, mock.Anything).Times(1).Return(&providerv1beta1.GetLockResponse{
					Status: status.NewOK(ctx),
					Lock: &providerv1beta1.Lock{
						LockId: "abcdef123",
						Type:   providerv1beta1.LockType_LOCK_TYPE_WRITE,
					},
				}, nil)

				newLockId, err := fc.Lock(ctx, "abcdef123", "")
				Expect(err).To(Succeed())
				Expect(newLockId).To(Equal("abcdef123"))
			})

			It("Set lock mismatches but get lock doesn't return lockId", func() {
				ctx := middleware.WopiContextToCtx(context.Background(), wopiCtx)

				gatewayClient.On("SetLock", mock.Anything, mock.Anything).Times(1).Return(&providerv1beta1.SetLockResponse{
					Status: status.NewFailedPrecondition(ctx, nil, "lock mismatch"),
				}, nil)

				gatewayClient.On("GetLock", mock.Anything, mock.Anything).Times(1).Return(&providerv1beta1.GetLockResponse{
					Status: status.NewOK(ctx),
				}, nil)

				newLockId, err := fc.Lock(ctx, "abcdef123", "")
				Expect(err).To(HaveOccurred())
				conErr := err.(*connector.ConnectorError)
				Expect(conErr.HttpCodeOut).To(Equal(500))
				Expect(newLockId).To(Equal(""))
			})

			It("File not found", func() {
				ctx := middleware.WopiContextToCtx(context.Background(), wopiCtx)

				gatewayClient.On("SetLock", mock.Anything, mock.Anything).Times(1).Return(&providerv1beta1.SetLockResponse{
					Status: status.NewNotFound(ctx, "file not found"),
				}, nil)

				newLockId, err := fc.Lock(ctx, "abcdef123", "")
				Expect(err).To(HaveOccurred())
				conErr := err.(*connector.ConnectorError)
				Expect(conErr.HttpCodeOut).To(Equal(404))
				Expect(newLockId).To(Equal(""))
			})

			It("Default error handling (insufficient storage)", func() {
				ctx := middleware.WopiContextToCtx(context.Background(), wopiCtx)

				gatewayClient.On("SetLock", mock.Anything, mock.Anything).Times(1).Return(&providerv1beta1.SetLockResponse{
					Status: status.NewInsufficientStorage(ctx, nil, "file too big"),
				}, nil)

				newLockId, err := fc.Lock(ctx, "abcdef123", "")
				Expect(err).To(HaveOccurred())
				conErr := err.(*connector.ConnectorError)
				Expect(conErr.HttpCodeOut).To(Equal(500))
				Expect(newLockId).To(Equal(""))
			})
		})

		Describe("Unlock and relock", func() {
			It("No valid context", func() {
				ctx := context.Background()
				newLockId, err := fc.Lock(ctx, "newLock", "oldLock")
				Expect(err).To(HaveOccurred())
				Expect(newLockId).To(Equal(""))
			})

			It("Empty lockId", func() {
				ctx := middleware.WopiContextToCtx(context.Background(), wopiCtx)

				newLockId, err := fc.Lock(ctx, "", "oldLock")
				Expect(err).To(HaveOccurred())
				conErr := err.(*connector.ConnectorError)
				Expect(conErr.HttpCodeOut).To(Equal(400))
				Expect(newLockId).To(Equal(""))
			})

			It("Refresh lock failed", func() {
				ctx := middleware.WopiContextToCtx(context.Background(), wopiCtx)

				targetErr := errors.New("Something went wrong")
				gatewayClient.On("RefreshLock", mock.Anything, mock.Anything).Times(1).Return(&providerv1beta1.RefreshLockResponse{
					Status: status.NewInternal(ctx, "Something failed"),
				}, targetErr)

				newLockId, err := fc.Lock(ctx, "abcdef123", "oldLock")
				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(targetErr))
				Expect(newLockId).To(Equal(""))
			})

			It("Refresh lock success", func() {
				ctx := middleware.WopiContextToCtx(context.Background(), wopiCtx)

				gatewayClient.On("RefreshLock", mock.Anything, mock.Anything).Times(1).Return(&providerv1beta1.RefreshLockResponse{
					Status: status.NewOK(ctx),
				}, nil)

				newLockId, err := fc.Lock(ctx, "abcdef123", "oldLock")
				Expect(err).To(Succeed())
				Expect(newLockId).To(Equal(""))
			})

			It("Refresh lock mismatches error getting lock", func() {
				ctx := middleware.WopiContextToCtx(context.Background(), wopiCtx)

				gatewayClient.On("RefreshLock", mock.Anything, mock.Anything).Times(1).Return(&providerv1beta1.RefreshLockResponse{
					Status: status.NewConflict(ctx, nil, "lock mismatch"),
				}, nil)

				targetErr := errors.New("Something went wrong getting the lock")
				gatewayClient.On("GetLock", mock.Anything, mock.Anything).Times(1).Return(&providerv1beta1.GetLockResponse{
					Status: status.NewInternal(ctx, "lock mismatch"),
				}, targetErr)

				newLockId, err := fc.Lock(ctx, "abcdef123", "112233")
				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(targetErr))
				Expect(newLockId).To(Equal(""))
			})

			It("Refresh lock mismatches", func() {
				ctx := middleware.WopiContextToCtx(context.Background(), wopiCtx)

				gatewayClient.On("RefreshLock", mock.Anything, mock.Anything).Times(1).Return(&providerv1beta1.RefreshLockResponse{
					Status: status.NewConflict(ctx, nil, "lock mismatch"),
				}, nil)

				gatewayClient.On("GetLock", mock.Anything, mock.Anything).Times(1).Return(&providerv1beta1.GetLockResponse{
					Status: status.NewOK(ctx),
					Lock: &providerv1beta1.Lock{
						LockId: "zzz999",
						Type:   providerv1beta1.LockType_LOCK_TYPE_WRITE,
					},
				}, nil)

				newLockId, err := fc.Lock(ctx, "abcdef123", "112233")
				Expect(err).To(HaveOccurred())
				conErr := err.(*connector.ConnectorError)
				Expect(conErr.HttpCodeOut).To(Equal(409))
				Expect(newLockId).To(Equal("zzz999"))
			})

			It("Refresh lock mismatches but get lock matches", func() {
				ctx := middleware.WopiContextToCtx(context.Background(), wopiCtx)

				gatewayClient.On("RefreshLock", mock.Anything, mock.Anything).Times(1).Return(&providerv1beta1.RefreshLockResponse{
					Status: status.NewConflict(ctx, nil, "lock mismatch"),
				}, nil)

				gatewayClient.On("GetLock", mock.Anything, mock.Anything).Times(1).Return(&providerv1beta1.GetLockResponse{
					Status: status.NewOK(ctx),
					Lock: &providerv1beta1.Lock{
						LockId: "abcdef123",
						Type:   providerv1beta1.LockType_LOCK_TYPE_WRITE,
					},
				}, nil)

				newLockId, err := fc.Lock(ctx, "abcdef123", "112233")
				Expect(err).To(Succeed())
				Expect(newLockId).To(Equal("abcdef123"))
			})

			It("Refresh lock mismatches but get lock doesn't return lockId", func() {
				ctx := middleware.WopiContextToCtx(context.Background(), wopiCtx)

				gatewayClient.On("RefreshLock", mock.Anything, mock.Anything).Times(1).Return(&providerv1beta1.RefreshLockResponse{
					Status: status.NewConflict(ctx, nil, "lock mismatch"),
				}, nil)

				gatewayClient.On("GetLock", mock.Anything, mock.Anything).Times(1).Return(&providerv1beta1.GetLockResponse{
					Status: status.NewOK(ctx),
				}, nil)

				newLockId, err := fc.Lock(ctx, "abcdef123", "112233")
				Expect(err).To(HaveOccurred())
				conErr := err.(*connector.ConnectorError)
				Expect(conErr.HttpCodeOut).To(Equal(500))
				Expect(newLockId).To(Equal(""))
			})

			It("File not found", func() {
				ctx := middleware.WopiContextToCtx(context.Background(), wopiCtx)

				gatewayClient.On("RefreshLock", mock.Anything, mock.Anything).Times(1).Return(&providerv1beta1.RefreshLockResponse{
					Status: status.NewNotFound(ctx, "file not found"),
				}, nil)

				newLockId, err := fc.Lock(ctx, "abcdef123", "112233")
				Expect(err).To(HaveOccurred())
				conErr := err.(*connector.ConnectorError)
				Expect(conErr.HttpCodeOut).To(Equal(404))
				Expect(newLockId).To(Equal(""))
			})

			It("Default error handling (insufficient storage)", func() {
				ctx := middleware.WopiContextToCtx(context.Background(), wopiCtx)

				gatewayClient.On("RefreshLock", mock.Anything, mock.Anything).Times(1).Return(&providerv1beta1.RefreshLockResponse{
					Status: status.NewInsufficientStorage(ctx, nil, "file too big"),
				}, nil)

				newLockId, err := fc.Lock(ctx, "abcdef123", "112233")
				Expect(err).To(HaveOccurred())
				conErr := err.(*connector.ConnectorError)
				Expect(conErr.HttpCodeOut).To(Equal(500))
				Expect(newLockId).To(Equal(""))
			})
		})
	})

	Describe("RefreshLock", func() {
		It("No valid context", func() {
			ctx := context.Background()
			newLockId, err := fc.RefreshLock(ctx, "newLock")
			Expect(err).To(HaveOccurred())
			Expect(newLockId).To(Equal(""))
		})

		It("Empty lockId", func() {
			ctx := middleware.WopiContextToCtx(context.Background(), wopiCtx)

			newLockId, err := fc.RefreshLock(ctx, "")
			Expect(err).To(HaveOccurred())
			conErr := err.(*connector.ConnectorError)
			Expect(conErr.HttpCodeOut).To(Equal(400))
			Expect(newLockId).To(Equal(""))
		})

		It("Refresh lock fails", func() {
			ctx := middleware.WopiContextToCtx(context.Background(), wopiCtx)

			targetErr := errors.New("Something went wrong")
			gatewayClient.On("RefreshLock", mock.Anything, mock.Anything).Times(1).Return(&providerv1beta1.RefreshLockResponse{
				Status: status.NewConflict(ctx, nil, "lock mismatch"),
			}, targetErr)

			newLockId, err := fc.RefreshLock(ctx, "abcdef123")
			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(targetErr))
			Expect(newLockId).To(Equal(""))
		})

		It("Refresh lock success", func() {
			ctx := middleware.WopiContextToCtx(context.Background(), wopiCtx)

			gatewayClient.On("RefreshLock", mock.Anything, mock.Anything).Times(1).Return(&providerv1beta1.RefreshLockResponse{
				Status: status.NewOK(ctx),
			}, nil)

			newLockId, err := fc.RefreshLock(ctx, "abcdef123")
			Expect(err).To(Succeed())
			Expect(newLockId).To(Equal(""))
		})

		It("Refresh lock file not found", func() {
			ctx := middleware.WopiContextToCtx(context.Background(), wopiCtx)

			gatewayClient.On("RefreshLock", mock.Anything, mock.Anything).Times(1).Return(&providerv1beta1.RefreshLockResponse{
				Status: status.NewNotFound(ctx, "file not found"),
			}, nil)

			newLockId, err := fc.RefreshLock(ctx, "abcdef123")
			Expect(err).To(HaveOccurred())
			conErr := err.(*connector.ConnectorError)
			Expect(conErr.HttpCodeOut).To(Equal(404))
			Expect(newLockId).To(Equal(""))
		})

		It("Refresh lock mismatch and get lock fails", func() {
			ctx := middleware.WopiContextToCtx(context.Background(), wopiCtx)

			gatewayClient.On("RefreshLock", mock.Anything, mock.Anything).Times(1).Return(&providerv1beta1.RefreshLockResponse{
				Status: status.NewConflict(ctx, nil, "lock mismatch"),
			}, nil)

			targetErr := errors.New("Something went wrong")
			gatewayClient.On("GetLock", mock.Anything, mock.Anything).Times(1).Return(&providerv1beta1.GetLockResponse{
				Status: status.NewConflict(ctx, nil, "lock mismatch"),
			}, targetErr)

			newLockId, err := fc.RefreshLock(ctx, "abcdef123")
			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(targetErr))
			Expect(newLockId).To(Equal(""))
		})

		It("Refresh lock mismatch and get lock status not ok", func() {
			ctx := middleware.WopiContextToCtx(context.Background(), wopiCtx)

			gatewayClient.On("RefreshLock", mock.Anything, mock.Anything).Times(1).Return(&providerv1beta1.RefreshLockResponse{
				Status: status.NewConflict(ctx, nil, "lock mismatch"),
			}, nil)

			gatewayClient.On("GetLock", mock.Anything, mock.Anything).Times(1).Return(&providerv1beta1.GetLockResponse{
				Status: status.NewInternal(ctx, "lock mismatch"),
			}, nil)

			newLockId, err := fc.RefreshLock(ctx, "abcdef123")
			Expect(err).To(HaveOccurred())
			conErr := err.(*connector.ConnectorError)
			Expect(conErr.HttpCodeOut).To(Equal(500))
			Expect(newLockId).To(Equal(""))
		})

		It("Refresh lock mismatch and no lock", func() {
			ctx := middleware.WopiContextToCtx(context.Background(), wopiCtx)

			gatewayClient.On("RefreshLock", mock.Anything, mock.Anything).Times(1).Return(&providerv1beta1.RefreshLockResponse{
				Status: status.NewConflict(ctx, nil, "lock mismatch"),
			}, nil)

			gatewayClient.On("GetLock", mock.Anything, mock.Anything).Times(1).Return(&providerv1beta1.GetLockResponse{
				Status: status.NewOK(ctx),
			}, nil)

			newLockId, err := fc.RefreshLock(ctx, "abcdef123")
			Expect(err).To(HaveOccurred())
			conErr := err.(*connector.ConnectorError)
			Expect(conErr.HttpCodeOut).To(Equal(409))
			Expect(newLockId).To(Equal(""))
		})

		It("Refresh lock mismatch", func() {
			ctx := middleware.WopiContextToCtx(context.Background(), wopiCtx)

			gatewayClient.On("RefreshLock", mock.Anything, mock.Anything).Times(1).Return(&providerv1beta1.RefreshLockResponse{
				Status: status.NewConflict(ctx, nil, "lock mismatch"),
			}, nil)

			gatewayClient.On("GetLock", mock.Anything, mock.Anything).Times(1).Return(&providerv1beta1.GetLockResponse{
				Status: status.NewOK(ctx),
				Lock: &providerv1beta1.Lock{
					LockId: "zzz999",
					Type:   providerv1beta1.LockType_LOCK_TYPE_WRITE,
				},
			}, nil)

			newLockId, err := fc.RefreshLock(ctx, "abcdef123")
			Expect(err).To(HaveOccurred())
			conErr := err.(*connector.ConnectorError)
			Expect(conErr.HttpCodeOut).To(Equal(409))
			Expect(newLockId).To(Equal("zzz999"))
		})

		It("Default error handling (insufficient storage)", func() {
			ctx := middleware.WopiContextToCtx(context.Background(), wopiCtx)

			gatewayClient.On("RefreshLock", mock.Anything, mock.Anything).Times(1).Return(&providerv1beta1.RefreshLockResponse{
				Status: status.NewInsufficientStorage(ctx, nil, "file too big"),
			}, nil)

			newLockId, err := fc.RefreshLock(ctx, "abcdef123")
			Expect(err).To(HaveOccurred())
			conErr := err.(*connector.ConnectorError)
			Expect(conErr.HttpCodeOut).To(Equal(500))
			Expect(newLockId).To(Equal(""))
		})
	})

	Describe("Unlock", func() {
		It("No valid context", func() {
			ctx := context.Background()
			newLockId, err := fc.UnLock(ctx, "newLock")
			Expect(err).To(HaveOccurred())
			Expect(newLockId).To(Equal(""))
		})

		It("Empty lockId", func() {
			ctx := middleware.WopiContextToCtx(context.Background(), wopiCtx)

			newLockId, err := fc.UnLock(ctx, "")
			Expect(err).To(HaveOccurred())
			conErr := err.(*connector.ConnectorError)
			Expect(conErr.HttpCodeOut).To(Equal(400))
			Expect(newLockId).To(Equal(""))
		})

		It("Unlock fails", func() {
			ctx := middleware.WopiContextToCtx(context.Background(), wopiCtx)

			targetErr := errors.New("Something went wrong")
			gatewayClient.On("Unlock", mock.Anything, mock.Anything).Times(1).Return(&providerv1beta1.UnlockResponse{
				Status: status.NewInternal(ctx, "something failed"),
			}, targetErr)

			newLockId, err := fc.UnLock(ctx, "abcdef123")
			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(targetErr))
			Expect(newLockId).To(Equal(""))
		})

		It("Unlock success", func() {
			ctx := middleware.WopiContextToCtx(context.Background(), wopiCtx)

			gatewayClient.On("Unlock", mock.Anything, mock.Anything).Times(1).Return(&providerv1beta1.UnlockResponse{
				Status: status.NewOK(ctx),
			}, nil)

			newLockId, err := fc.UnLock(ctx, "abcdef123")
			Expect(err).To(Succeed())
			Expect(newLockId).To(Equal(""))
		})

		It("Unlock file isn't locked", func() {
			ctx := middleware.WopiContextToCtx(context.Background(), wopiCtx)

			gatewayClient.On("Unlock", mock.Anything, mock.Anything).Times(1).Return(&providerv1beta1.UnlockResponse{
				Status: status.NewConflict(ctx, nil, "lock mismatch"),
			}, nil)

			newLockId, err := fc.UnLock(ctx, "abcdef123")
			Expect(err).To(HaveOccurred())
			conErr := err.(*connector.ConnectorError)
			Expect(conErr.HttpCodeOut).To(Equal(409))
			Expect(newLockId).To(Equal(""))
		})

		It("Unlock mismatch get lock fails", func() {
			ctx := middleware.WopiContextToCtx(context.Background(), wopiCtx)

			gatewayClient.On("Unlock", mock.Anything, mock.Anything).Times(1).Return(&providerv1beta1.UnlockResponse{
				Status: status.NewLocked(ctx, "lock mismatch"),
			}, nil)

			targetErr := errors.New("Something went wrong")
			gatewayClient.On("GetLock", mock.Anything, mock.Anything).Times(1).Return(&providerv1beta1.GetLockResponse{
				Status: status.NewInternal(ctx, "something failed"),
			}, targetErr)

			newLockId, err := fc.UnLock(ctx, "abcdef123")
			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(targetErr))
			Expect(newLockId).To(Equal(""))
		})

		It("Unlock mismatch get lock status not ok", func() {
			ctx := middleware.WopiContextToCtx(context.Background(), wopiCtx)

			gatewayClient.On("Unlock", mock.Anything, mock.Anything).Times(1).Return(&providerv1beta1.UnlockResponse{
				Status: status.NewLocked(ctx, "lock mismatch"),
			}, nil)

			gatewayClient.On("GetLock", mock.Anything, mock.Anything).Times(1).Return(&providerv1beta1.GetLockResponse{
				Status: status.NewInternal(ctx, "something failed"),
			}, nil)

			newLockId, err := fc.UnLock(ctx, "abcdef123")
			Expect(err).To(HaveOccurred())
			conErr := err.(*connector.ConnectorError)
			Expect(conErr.HttpCodeOut).To(Equal(500))
			Expect(newLockId).To(Equal(""))
		})

		It("Unlock mismatch get lock doesn't return lock", func() {
			ctx := middleware.WopiContextToCtx(context.Background(), wopiCtx)

			gatewayClient.On("Unlock", mock.Anything, mock.Anything).Times(1).Return(&providerv1beta1.UnlockResponse{
				Status: status.NewLocked(ctx, "lock mismatch"),
			}, nil)

			gatewayClient.On("GetLock", mock.Anything, mock.Anything).Times(1).Return(&providerv1beta1.GetLockResponse{
				Status: status.NewOK(ctx),
			}, nil)

			newLockId, err := fc.UnLock(ctx, "abcdef123")
			Expect(err).To(HaveOccurred())
			conErr := err.(*connector.ConnectorError)
			Expect(conErr.HttpCodeOut).To(Equal(409))
			Expect(newLockId).To(Equal(""))
		})

		It("Unlock mismatch", func() {
			ctx := middleware.WopiContextToCtx(context.Background(), wopiCtx)

			gatewayClient.On("Unlock", mock.Anything, mock.Anything).Times(1).Return(&providerv1beta1.UnlockResponse{
				Status: status.NewLocked(ctx, "lock mismatch"),
			}, nil)

			gatewayClient.On("GetLock", mock.Anything, mock.Anything).Times(1).Return(&providerv1beta1.GetLockResponse{
				Status: status.NewOK(ctx),
				Lock: &providerv1beta1.Lock{
					LockId: "zzz999",
					Type:   providerv1beta1.LockType_LOCK_TYPE_WRITE,
				},
			}, nil)

			newLockId, err := fc.UnLock(ctx, "abcdef123")
			Expect(err).To(HaveOccurred())
			conErr := err.(*connector.ConnectorError)
			Expect(conErr.HttpCodeOut).To(Equal(409))
			Expect(newLockId).To(Equal("zzz999"))
		})

		It("Default error handling (insufficient storage)", func() {
			ctx := middleware.WopiContextToCtx(context.Background(), wopiCtx)

			gatewayClient.On("Unlock", mock.Anything, mock.Anything).Times(1).Return(&providerv1beta1.UnlockResponse{
				Status: status.NewInsufficientStorage(ctx, nil, "file too big"),
			}, nil)

			newLockId, err := fc.UnLock(ctx, "abcdef123")
			Expect(err).To(HaveOccurred())
			conErr := err.(*connector.ConnectorError)
			Expect(conErr.HttpCodeOut).To(Equal(500))
			Expect(newLockId).To(Equal(""))
		})
	})

	Describe("CheckFileInfo", func() {
		It("No valid context", func() {
			ctx := context.Background()
			newFileInfo, err := fc.CheckFileInfo(ctx)
			Expect(err).To(HaveOccurred())
			Expect(newFileInfo).To(BeNil())
		})

		It("Stat fails", func() {
			ctx := middleware.WopiContextToCtx(context.Background(), wopiCtx)

			targetErr := errors.New("Something went wrong")
			gatewayClient.On("Stat", mock.Anything, mock.Anything).Times(1).Return(&providerv1beta1.StatResponse{
				Status: status.NewInternal(ctx, "something failed"),
			}, targetErr)

			newFileInfo, err := fc.CheckFileInfo(ctx)
			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(targetErr))
			Expect(newFileInfo).To(BeNil())
		})

		It("Stat fails status not ok", func() {
			ctx := middleware.WopiContextToCtx(context.Background(), wopiCtx)

			gatewayClient.On("Stat", mock.Anything, mock.Anything).Times(1).Return(&providerv1beta1.StatResponse{
				Status: status.NewInternal(ctx, "something failed"),
			}, nil)

			newFileInfo, err := fc.CheckFileInfo(ctx)
			Expect(err).To(HaveOccurred())
			conErr := err.(*connector.ConnectorError)
			Expect(conErr.HttpCodeOut).To(Equal(500))
			Expect(newFileInfo).To(BeNil())
		})

		It("Stat success", func() {
			ctx := middleware.WopiContextToCtx(context.Background(), wopiCtx)

			gatewayClient.On("Stat", mock.Anything, mock.Anything).Times(1).Return(&providerv1beta1.StatResponse{
				Status: status.NewOK(ctx),
				Info: &providerv1beta1.ResourceInfo{
					Owner: &userv1beta1.UserId{
						Idp:      "customIdp",
						OpaqueId: "aabbcc",
						Type:     userv1beta1.UserType_USER_TYPE_PRIMARY,
					},
					Size: uint64(998877),
					Mtime: &typesv1beta1.Timestamp{
						Seconds: uint64(16273849),
					},
					Path: "/path/to/test.txt",
					// Other properties aren't used for now.
				},
			}, nil)

			expectedFileInfo := &fileinfo.Microsoft{
				OwnerId:                    "61616262636340637573746f6d496470", // hex of aabbcc@customIdp
				Size:                       int64(998877),
				Version:                    "16273849.0",
				BaseFileName:               "test.txt",
				BreadcrumbDocName:          "test.txt",
				UserCanNotWriteRelative:    true,
				HostViewUrl:                "http://test.ex.prv/view",
				HostEditUrl:                "http://test.ex.prv/edit",
				SupportsExtendedLockLength: true,
				SupportsGetLock:            true,
				SupportsLocks:              true,
				SupportsUpdate:             true,
				UserCanWrite:               true,
				UserId:                     "6f7061717565496440696e6d656d6f7279", // hex of opaqueId@inmemory
				UserFriendlyName:           "Pet Shaft",
			}

			newFileInfo, err := fc.CheckFileInfo(ctx)
			Expect(err).To(Succeed())
			Expect(newFileInfo.(*fileinfo.Microsoft)).To(Equal(expectedFileInfo))
		})

		It("Stat success guests", func() {
			// add user's opaque to include public-share-role
			wopiCtx.User.Opaque = &typesv1beta1.Opaque{
				Map: map[string]*typesv1beta1.OpaqueEntry{
					"public-share-role": &typesv1beta1.OpaqueEntry{
						Decoder: "plain",
						Value:   []byte("viewer"),
					},
				},
			}
			// change view mode to view only
			wopiCtx.ViewMode = appproviderv1beta1.ViewMode_VIEW_MODE_VIEW_ONLY

			ctx := middleware.WopiContextToCtx(context.Background(), wopiCtx)

			gatewayClient.On("Stat", mock.Anything, mock.Anything).Times(1).Return(&providerv1beta1.StatResponse{
				Status: status.NewOK(ctx),
				Info: &providerv1beta1.ResourceInfo{
					Owner: &userv1beta1.UserId{
						Idp:      "customIdp",
						OpaqueId: "aabbcc",
						Type:     userv1beta1.UserType_USER_TYPE_PRIMARY,
					},
					Size: uint64(998877),
					Mtime: &typesv1beta1.Timestamp{
						Seconds: uint64(16273849),
					},
					Path: "/path/to/test.txt",
					// Other properties aren't used for now.
				},
			}, nil)

			// change wopi app provider
			cfg.WopiApp.Provider = "Collabora"

			expectedFileInfo := &fileinfo.Collabora{
				OwnerId:                 "61616262636340637573746f6d496470", // hex of aabbcc@customIdp
				Size:                    int64(998877),
				BaseFileName:            "test.txt",
				UserCanNotWriteRelative: true,
				DisableExport:           true,
				DisableCopy:             true,
				DisablePrint:            true,
				UserId:                  "guest-zzz000",
				UserFriendlyName:        "guest zzz000",
				EnableOwnerTermination:  true,
			}

			newFileInfo, err := fc.CheckFileInfo(ctx)

			// UserId and UserFriendlyName have random Ids generated which are impossible to guess
			// Check both separately
			Expect(newFileInfo.(*fileinfo.Collabora).UserId).To(HavePrefix(hex.EncodeToString([]byte("guest-"))))
			Expect(newFileInfo.(*fileinfo.Collabora).UserFriendlyName).To(HavePrefix("Guest "))
			// overwrite UserId and UserFriendlyName here for easier matching
			newFileInfo.(*fileinfo.Collabora).UserId = "guest-zzz000"
			newFileInfo.(*fileinfo.Collabora).UserFriendlyName = "guest zzz000"

			Expect(err).To(Succeed())
			Expect(newFileInfo.(*fileinfo.Collabora)).To(Equal(expectedFileInfo))
		})

		It("Stat success authenticated user", func() {
			// change view mode to view only
			wopiCtx.ViewMode = appproviderv1beta1.ViewMode_VIEW_MODE_VIEW_ONLY

			ctx := middleware.WopiContextToCtx(context.Background(), wopiCtx)

			gatewayClient.On("Stat", mock.Anything, mock.Anything).Times(1).Return(&providerv1beta1.StatResponse{
				Status: status.NewOK(ctx),
				Info: &providerv1beta1.ResourceInfo{
					Owner: &userv1beta1.UserId{
						Idp:      "customIdp",
						OpaqueId: "aabbcc",
						Type:     userv1beta1.UserType_USER_TYPE_PRIMARY,
					},
					Size: uint64(998877),
					Mtime: &typesv1beta1.Timestamp{
						Seconds: uint64(16273849),
					},
					Path: "/path/to/test.txt",
					// Other properties aren't used for now.
				},
			}, nil)

			// change wopi app provider
			cfg.WopiApp.Provider = "Collabora"

			expectedFileInfo := &fileinfo.Collabora{
				OwnerId:                 "61616262636340637573746f6d496470", // hex of aabbcc@customIdp
				Size:                    int64(998877),
				BaseFileName:            "test.txt",
				UserCanNotWriteRelative: true,
				DisableExport:           true,
				DisableCopy:             true,
				DisablePrint:            true,
				UserId:                  hex.EncodeToString([]byte("opaqueId@inmemory")),
				UserFriendlyName:        "Pet Shaft",
				EnableOwnerTermination:  true,
				WatermarkText:           "Pet Shaft shaft@example.com",
			}

			newFileInfo, err := fc.CheckFileInfo(ctx)

			Expect(err).To(Succeed())
			Expect(newFileInfo.(*fileinfo.Collabora)).To(Equal(expectedFileInfo))
		})
	})
})
