package svc_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"

	storageprovider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/cs3org/reva/v2/pkg/storagespace"
	"github.com/go-chi/chi/v5"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	libregraph "github.com/owncloud/libre-graph-api-go"
	"github.com/stretchr/testify/mock"
	"github.com/tidwall/gjson"

	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/services/graph/mocks"
	svc "github.com/owncloud/ocis/v2/services/graph/pkg/service/v0"
)

var _ = Describe("DrivesDriveItemApi", func() {
	var (
		mockProvider *mocks.DrivesDriveItemProvider
		httpAPI      svc.DrivesDriveItemApi
		rCTX         *chi.Context
	)

	BeforeEach(func() {
		logger := log.NewLogger()

		mockProvider = mocks.NewDrivesDriveItemProvider(GinkgoT())
		api, err := svc.NewDrivesDriveItemApi(mockProvider, logger)
		Expect(err).ToNot(HaveOccurred())

		httpAPI = api

		rCTX = chi.NewRouteContext()
		rCTX.URLParams.Add("driveID", "a0ca6a90-a365-4782-871e-d44447bbc668$a0ca6a90-a365-4782-871e-d44447bbc668")
		rCTX.URLParams.Add("itemID", "a0ca6a90-a365-4782-871e-d44447bbc668$a0ca6a90-a365-4782-871e-d44447bbc668!a0ca6a90-a365-4782-871e-d44447bbc668")
	})

	checkDriveIDAndItemIDValidation := func(handler http.HandlerFunc) {
		rCTX.URLParams.Add("driveID", "1$2")
		rCTX.URLParams.Add("itemID", "3$4!5")

		responseRecorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodPost, "/", nil).
			WithContext(
				context.WithValue(context.Background(), chi.RouteCtxKey, rCTX),
			)

		handler(responseRecorder, request)

		Expect(responseRecorder.Code).To(Equal(http.StatusUnprocessableEntity))

		jsonData := gjson.Get(responseRecorder.Body.String(), "error")
		Expect(jsonData.Get("message").String()).To(Equal("invalid driveID or itemID"))
	}

	Describe("CreateDriveItem", func() {
		It("validates the driveID and itemID url param", func() {
			checkDriveIDAndItemIDValidation(httpAPI.DeleteDriveItem)
		})

		It("uses the UnmountShare provider implementation", func() {
			responseRecorder := httptest.NewRecorder()

			request := httptest.NewRequest(http.MethodDelete, "/", nil).
				WithContext(
					context.WithValue(context.Background(), chi.RouteCtxKey, rCTX),
				)

			onUnmountShare := mockProvider.On("UnmountShare", mock.Anything, mock.Anything)
			onUnmountShare.
				Return(func(ctx context.Context, resourceID storageprovider.ResourceId) error {
					return errors.New("any")
				}).Once()

			httpAPI.DeleteDriveItem(responseRecorder, request)

			Expect(responseRecorder.Code).To(Equal(http.StatusFailedDependency))

			jsonData := gjson.Get(responseRecorder.Body.String(), "error")
			Expect(jsonData.Get("message").String()).To(Equal("unmounting share failed"))

			// happy path
			responseRecorder = httptest.NewRecorder()

			onUnmountShare.
				Return(func(ctx context.Context, resourceID storageprovider.ResourceId) error {
					Expect(storagespace.FormatResourceID(resourceID)).To(Equal("a0ca6a90-a365-4782-871e-d44447bbc668$a0ca6a90-a365-4782-871e-d44447bbc668!a0ca6a90-a365-4782-871e-d44447bbc668"))
					return nil
				}).Once()

			httpAPI.DeleteDriveItem(responseRecorder, request)

			Expect(responseRecorder.Code).To(Equal(http.StatusOK))
		})
	})

	Describe("CreateDriveItem", func() {
		It("validates the driveID and itemID url param", func() {
			checkDriveIDAndItemIDValidation(httpAPI.CreateDriveItem)
		})

		It("checks if the idemID and driveID is in share jail", func() {
			rCTX.URLParams.Add("driveID", "1$2")
			rCTX.URLParams.Add("itemID", "1$2!3")

			responseRecorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodPost, "/", nil).
				WithContext(
					context.WithValue(context.Background(), chi.RouteCtxKey, rCTX),
				)

			httpAPI.CreateDriveItem(responseRecorder, request)

			Expect(responseRecorder.Code).To(Equal(http.StatusUnprocessableEntity))

			jsonData := gjson.Get(responseRecorder.Body.String(), "error")
			Expect(jsonData.Get("message").String()).To(ContainSubstring("must be share jail"))
		})

		It("checks that the request body is valid", func() {
			responseRecorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodPost, "/", nil).
				WithContext(
					context.WithValue(context.Background(), chi.RouteCtxKey, rCTX),
				)

			httpAPI.CreateDriveItem(responseRecorder, request)

			Expect(responseRecorder.Code).To(Equal(http.StatusUnprocessableEntity))

			jsonData := gjson.Get(responseRecorder.Body.String(), "error")
			Expect(jsonData.Get("message").String()).To(Equal("invalid request body"))

			// valid drive item, but invalid remote item id
			driveItem := libregraph.DriveItem{}

			driveItemJson, err := json.Marshal(driveItem)
			Expect(err).ToNot(HaveOccurred())

			responseRecorder = httptest.NewRecorder()

			request = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(driveItemJson)).
				WithContext(
					context.WithValue(context.Background(), chi.RouteCtxKey, rCTX),
				)

			httpAPI.CreateDriveItem(responseRecorder, request)

			Expect(responseRecorder.Code).To(Equal(http.StatusUnprocessableEntity))

			jsonData = gjson.Get(responseRecorder.Body.String(), "error")
			Expect(jsonData.Get("message").String()).To(Equal("invalid remote item id"))
		})

		It("uses the MountShare provider implementation", func() {
			driveItemName := "a name"
			remoteItemID := "d66d28d8-3558-4f0f-ba2a-34a7185b806d$831997cf-a531-491b-ae72-9037739f04e9!c131a84c-7506-46b4-8e5e-60c56382da3b"
			driveItem := libregraph.DriveItem{
				Name: &driveItemName,
				RemoteItem: &libregraph.RemoteItem{
					Id: &remoteItemID,
				},
			}

			driveItemJson, err := json.Marshal(driveItem)
			Expect(err).ToNot(HaveOccurred())

			responseRecorder := httptest.NewRecorder()

			request := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(driveItemJson)).
				WithContext(
					context.WithValue(context.Background(), chi.RouteCtxKey, rCTX),
				)

			onMountShare := mockProvider.On("MountShare", mock.Anything, mock.Anything, mock.Anything)
			onMountShare.
				Return(func(ctx context.Context, resourceID storageprovider.ResourceId, name string) (libregraph.DriveItem, error) {
					return libregraph.DriveItem{}, errors.New("any")
				}).Once()

			httpAPI.CreateDriveItem(responseRecorder, request)

			Expect(responseRecorder.Code).To(Equal(http.StatusFailedDependency))

			jsonData := gjson.Get(responseRecorder.Body.String(), "error")
			Expect(jsonData.Get("message").String()).To(Equal("mounting share failed"))

			// happy path
			responseRecorder = httptest.NewRecorder()

			request = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(driveItemJson)).
				WithContext(
					context.WithValue(context.Background(), chi.RouteCtxKey, rCTX),
				)

			onMountShare.
				Return(func(ctx context.Context, resourceID storageprovider.ResourceId, name string) (libregraph.DriveItem, error) {
					Expect(storagespace.FormatResourceID(resourceID)).To(Equal(remoteItemID))
					Expect(driveItemName).To(Equal(name))
					return libregraph.DriveItem{}, nil
				}).Once()

			httpAPI.CreateDriveItem(responseRecorder, request)

			Expect(responseRecorder.Code).To(Equal(http.StatusCreated))
		})
	})
})
