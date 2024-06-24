package connector_test

import (
	"encoding/json"
	"errors"
	"io"
	"net/http/httptest"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/owncloud/ocis/v2/services/collaboration/mocks"
	"github.com/owncloud/ocis/v2/services/collaboration/pkg/connector"
	"github.com/owncloud/ocis/v2/services/collaboration/pkg/connector/fileinfo"
	"github.com/stretchr/testify/mock"
)

var _ = Describe("HttpAdapter", func() {
	var (
		fc          *mocks.FileConnectorService
		cc          *mocks.ContentConnectorService
		con         *mocks.ConnectorService
		httpAdapter *connector.HttpAdapter
	)

	BeforeEach(func() {
		fc = &mocks.FileConnectorService{}
		cc = &mocks.ContentConnectorService{}

		con = &mocks.ConnectorService{}
		con.On("GetContentConnector").Return(cc)
		con.On("GetFileConnector").Return(fc)

		httpAdapter = connector.NewHttpAdapterWithConnector(con)
	})

	Describe("GetLock", func() {
		It("General error", func() {
			req := httptest.NewRequest("POST", "/wopi/files/abcdef", nil)
			req.Header.Set("X-WOPI-Override", "POST_LOCK")

			w := httptest.NewRecorder()

			fc.On("GetLock", mock.Anything).Times(1).Return("", errors.New("Something happened"))

			httpAdapter.GetLock(w, req)
			resp := w.Result()
			Expect(resp.StatusCode).To(Equal(500))
		})

		It("File not found", func() {
			req := httptest.NewRequest("POST", "/wopi/files/abcdef", nil)
			req.Header.Set("X-WOPI-Override", "POST_LOCK")

			w := httptest.NewRecorder()

			fc.On("GetLock", mock.Anything).Times(1).Return("", connector.NewConnectorError(404, "Couldn't get the file"))

			httpAdapter.GetLock(w, req)
			resp := w.Result()
			Expect(resp.StatusCode).To(Equal(404))
		})

		It("LockId", func() {
			req := httptest.NewRequest("POST", "/wopi/files/abcdef", nil)
			req.Header.Set("X-WOPI-Override", "POST_LOCK")

			w := httptest.NewRecorder()

			fc.On("GetLock", mock.Anything).Times(1).Return("zzz111", nil)

			httpAdapter.GetLock(w, req)
			resp := w.Result()
			Expect(resp.StatusCode).To(Equal(200))
			Expect(resp.Header.Get(connector.HeaderWopiLock)).To(Equal("zzz111"))
		})

		It("Empty LockId", func() {
			req := httptest.NewRequest("POST", "/wopi/files/abcdef", nil)
			req.Header.Set("X-WOPI-Override", "POST_LOCK")

			w := httptest.NewRecorder()

			fc.On("GetLock", mock.Anything).Times(1).Return("", nil)

			httpAdapter.GetLock(w, req)
			resp := w.Result()
			Expect(resp.StatusCode).To(Equal(200))
			Expect(resp.Header.Get(connector.HeaderWopiLock)).To(Equal(""))
		})
	})

	Describe("Lock", func() {
		Describe("Just lock", func() {
			It("General error", func() {
				req := httptest.NewRequest("POST", "/wopi/files/abcdef", nil)
				req.Header.Set("X-WOPI-Override", "LOCK")
				req.Header.Set(connector.HeaderWopiLock, "abc123")

				w := httptest.NewRecorder()

				fc.On("Lock", mock.Anything, "abc123", "").Times(1).Return("", errors.New("Something happened"))

				httpAdapter.Lock(w, req)
				resp := w.Result()
				Expect(resp.StatusCode).To(Equal(500))
			})

			It("No LockId provided", func() {
				req := httptest.NewRequest("POST", "/wopi/files/abcdef", nil)
				req.Header.Set("X-WOPI-Override", "LOCK")
				req.Header.Set(connector.HeaderWopiLock, "")

				w := httptest.NewRecorder()

				fc.On("Lock", mock.Anything, "", "").Times(1).Return("", connector.NewConnectorError(400, "No lockId"))

				httpAdapter.Lock(w, req)
				resp := w.Result()
				Expect(resp.StatusCode).To(Equal(400))
			})

			It("Conflict", func() {
				req := httptest.NewRequest("POST", "/wopi/files/abcdef", nil)
				req.Header.Set("X-WOPI-Override", "LOCK")
				req.Header.Set(connector.HeaderWopiLock, "abc123")

				w := httptest.NewRecorder()

				fc.On("Lock", mock.Anything, "abc123", "").Times(1).Return("zzz111", connector.NewConnectorError(409, "Lock conflict"))

				httpAdapter.Lock(w, req)
				resp := w.Result()
				Expect(resp.StatusCode).To(Equal(409))
				Expect(resp.Header.Get(connector.HeaderWopiLock)).To(Equal("zzz111"))
			})

			It("Success", func() {
				req := httptest.NewRequest("POST", "/wopi/files/abcdef", nil)
				req.Header.Set("X-WOPI-Override", "LOCK")
				req.Header.Set(connector.HeaderWopiLock, "abc123")

				w := httptest.NewRecorder()

				fc.On("Lock", mock.Anything, "abc123", "").Times(1).Return("", nil)

				httpAdapter.Lock(w, req)
				resp := w.Result()
				Expect(resp.StatusCode).To(Equal(200))
			})
		})

		Describe("Unlock and relock", func() {
			It("General error", func() {
				req := httptest.NewRequest("POST", "/wopi/files/abcdef", nil)
				req.Header.Set("X-WOPI-Override", "LOCK")
				req.Header.Set(connector.HeaderWopiLock, "abc123")
				req.Header.Set(connector.HeaderWopiOldLock, "qwerty")

				w := httptest.NewRecorder()

				fc.On("Lock", mock.Anything, "abc123", "qwerty").Times(1).Return("", errors.New("Something happened"))

				httpAdapter.Lock(w, req)
				resp := w.Result()
				Expect(resp.StatusCode).To(Equal(500))
			})

			It("No LockId provided", func() {
				req := httptest.NewRequest("POST", "/wopi/files/abcdef", nil)
				req.Header.Set("X-WOPI-Override", "LOCK")
				req.Header.Set(connector.HeaderWopiLock, "")
				req.Header.Set(connector.HeaderWopiOldLock, "")

				w := httptest.NewRecorder()

				fc.On("Lock", mock.Anything, "", "").Times(1).Return("", connector.NewConnectorError(400, "No lockId"))

				httpAdapter.Lock(w, req)
				resp := w.Result()
				Expect(resp.StatusCode).To(Equal(400))
			})

			It("Conflict", func() {
				req := httptest.NewRequest("POST", "/wopi/files/abcdef", nil)
				req.Header.Set("X-WOPI-Override", "LOCK")
				req.Header.Set(connector.HeaderWopiLock, "abc123")
				req.Header.Set(connector.HeaderWopiOldLock, "qwerty")

				w := httptest.NewRecorder()

				fc.On("Lock", mock.Anything, "abc123", "qwerty").Times(1).Return("zzz111", connector.NewConnectorError(409, "Lock conflict"))

				httpAdapter.Lock(w, req)
				resp := w.Result()
				Expect(resp.StatusCode).To(Equal(409))
				Expect(resp.Header.Get(connector.HeaderWopiLock)).To(Equal("zzz111"))
			})

			It("Success", func() {
				req := httptest.NewRequest("POST", "/wopi/files/abcdef", nil)
				req.Header.Set("X-WOPI-Override", "LOCK")
				req.Header.Set(connector.HeaderWopiLock, "abc123")
				req.Header.Set(connector.HeaderWopiOldLock, "qwerty")

				w := httptest.NewRecorder()

				fc.On("Lock", mock.Anything, "abc123", "qwerty").Times(1).Return("", nil)

				httpAdapter.Lock(w, req)
				resp := w.Result()
				Expect(resp.StatusCode).To(Equal(200))
			})
		})
	})

	Describe("RefreshLock", func() {
		It("General error", func() {
			req := httptest.NewRequest("POST", "/wopi/files/abcdef", nil)
			req.Header.Set("X-WOPI-Override", "REFRESH_LOCK")
			req.Header.Set(connector.HeaderWopiLock, "abc123")

			w := httptest.NewRecorder()

			fc.On("RefreshLock", mock.Anything, "abc123").Times(1).Return("", errors.New("Something happened"))

			httpAdapter.RefreshLock(w, req)
			resp := w.Result()
			Expect(resp.StatusCode).To(Equal(500))
		})

		It("No LockId provided", func() {
			req := httptest.NewRequest("POST", "/wopi/files/abcdef", nil)
			req.Header.Set("X-WOPI-Override", "REFRESH_LOCK")
			req.Header.Set(connector.HeaderWopiLock, "")

			w := httptest.NewRecorder()

			fc.On("RefreshLock", mock.Anything, "").Times(1).Return("", connector.NewConnectorError(400, "No lockId"))

			httpAdapter.RefreshLock(w, req)
			resp := w.Result()
			Expect(resp.StatusCode).To(Equal(400))
		})

		It("Conflict", func() {
			req := httptest.NewRequest("POST", "/wopi/files/abcdef", nil)
			req.Header.Set("X-WOPI-Override", "REFRESH_LOCK")
			req.Header.Set(connector.HeaderWopiLock, "abc123")

			w := httptest.NewRecorder()

			fc.On("RefreshLock", mock.Anything, "abc123").Times(1).Return("zzz111", connector.NewConnectorError(409, "Lock conflict"))

			httpAdapter.RefreshLock(w, req)
			resp := w.Result()
			Expect(resp.StatusCode).To(Equal(409))
			Expect(resp.Header.Get(connector.HeaderWopiLock)).To(Equal("zzz111"))
		})

		It("Success", func() {
			req := httptest.NewRequest("POST", "/wopi/files/abcdef", nil)
			req.Header.Set("X-WOPI-Override", "REFRESH_LOCK")
			req.Header.Set(connector.HeaderWopiLock, "abc123")

			w := httptest.NewRecorder()

			fc.On("RefreshLock", mock.Anything, "abc123").Times(1).Return("", nil)

			httpAdapter.RefreshLock(w, req)
			resp := w.Result()
			Expect(resp.StatusCode).To(Equal(200))
		})
	})

	Describe("Unlock", func() {
		It("General error", func() {
			req := httptest.NewRequest("POST", "/wopi/files/abcdef", nil)
			req.Header.Set("X-WOPI-Override", "UNLOCK")
			req.Header.Set(connector.HeaderWopiLock, "abc123")

			w := httptest.NewRecorder()

			fc.On("UnLock", mock.Anything, "abc123").Times(1).Return("", errors.New("Something happened"))

			httpAdapter.UnLock(w, req)
			resp := w.Result()
			Expect(resp.StatusCode).To(Equal(500))
		})

		It("No LockId provided", func() {
			req := httptest.NewRequest("POST", "/wopi/files/abcdef", nil)
			req.Header.Set("X-WOPI-Override", "UNLOCK")
			req.Header.Set(connector.HeaderWopiLock, "")

			w := httptest.NewRecorder()

			fc.On("UnLock", mock.Anything, "").Times(1).Return("", connector.NewConnectorError(400, "No lockId"))

			httpAdapter.UnLock(w, req)
			resp := w.Result()
			Expect(resp.StatusCode).To(Equal(400))
		})

		It("Conflict", func() {
			req := httptest.NewRequest("POST", "/wopi/files/abcdef", nil)
			req.Header.Set("X-WOPI-Override", "UNLOCK")
			req.Header.Set(connector.HeaderWopiLock, "abc123")

			w := httptest.NewRecorder()

			fc.On("UnLock", mock.Anything, "abc123").Times(1).Return("zzz111", connector.NewConnectorError(409, "Lock conflict"))

			httpAdapter.UnLock(w, req)
			resp := w.Result()
			Expect(resp.StatusCode).To(Equal(409))
			Expect(resp.Header.Get(connector.HeaderWopiLock)).To(Equal("zzz111"))
		})

		It("Success", func() {
			req := httptest.NewRequest("POST", "/wopi/files/abcdef", nil)
			req.Header.Set("X-WOPI-Override", "UNLOCK")
			req.Header.Set(connector.HeaderWopiLock, "abc123")

			w := httptest.NewRecorder()

			fc.On("UnLock", mock.Anything, "abc123").Times(1).Return("", nil)

			httpAdapter.UnLock(w, req)
			resp := w.Result()
			Expect(resp.StatusCode).To(Equal(200))
		})
	})

	Describe("CheckFileInfo", func() {
		It("General error", func() {
			req := httptest.NewRequest("GET", "/wopi/files/abcdef", nil)

			w := httptest.NewRecorder()

			fc.On("CheckFileInfo", mock.Anything).Times(1).Return(&fileinfo.Microsoft{}, errors.New("Something happened"))

			httpAdapter.CheckFileInfo(w, req)
			resp := w.Result()
			Expect(resp.StatusCode).To(Equal(500))
		})

		It("Not found", func() {
			// 404 isn't thrown at the moment. Test is here to prove it's possible to
			// throw any error code
			req := httptest.NewRequest("GET", "/wopi/files/abcdef", nil)

			w := httptest.NewRecorder()

			fc.On("CheckFileInfo", mock.Anything).Times(1).Return(&fileinfo.Microsoft{}, connector.NewConnectorError(404, "Not found"))

			httpAdapter.CheckFileInfo(w, req)
			resp := w.Result()
			Expect(resp.StatusCode).To(Equal(404))
		})

		It("Success", func() {
			req := httptest.NewRequest("GET", "/wopi/files/abcdef", nil)

			w := httptest.NewRecorder()

			// might need more info, but should be enough for the test
			finfo := &fileinfo.Microsoft{
				Size:              123456789,
				BreadcrumbDocName: "testy.docx",
			}
			fc.On("CheckFileInfo", mock.Anything).Times(1).Return(finfo, nil)

			httpAdapter.CheckFileInfo(w, req)
			resp := w.Result()
			Expect(resp.StatusCode).To(Equal(200))

			jsonInfo, _ := io.ReadAll(resp.Body)

			var responseInfo *fileinfo.Microsoft
			json.Unmarshal(jsonInfo, &responseInfo)
			Expect(responseInfo).To(Equal(finfo))
		})
	})

	Describe("GetFile", func() {
		It("General error", func() {
			req := httptest.NewRequest("GET", "/wopi/files/abcdef/contents", nil)

			w := httptest.NewRecorder()

			cc.On("GetFile", mock.Anything, mock.Anything).Times(1).Return(errors.New("Something happened"))

			httpAdapter.GetFile(w, req)
			resp := w.Result()
			Expect(resp.StatusCode).To(Equal(500))
		})

		It("Not found", func() {
			// 404 isn't thrown at the moment. Test is here to prove it's possible to
			// throw any error code
			req := httptest.NewRequest("GET", "/wopi/files/abcdef/contents", nil)

			w := httptest.NewRecorder()

			cc.On("GetFile", mock.Anything, mock.Anything).Times(1).Return(connector.NewConnectorError(404, "Not found"))

			httpAdapter.GetFile(w, req)
			resp := w.Result()
			Expect(resp.StatusCode).To(Equal(404))
		})

		It("Success", func() {
			req := httptest.NewRequest("GET", "/wopi/files/abcdef/contents", nil)

			w := httptest.NewRecorder()

			expectedContent := []byte("This is a fake content for a test file")
			cc.On("GetFile", mock.Anything, mock.Anything).Times(1).Run(func(args mock.Arguments) {
				w := args.Get(1).(io.Writer)
				w.Write(expectedContent)
			}).Return(nil)

			httpAdapter.GetFile(w, req)
			resp := w.Result()
			Expect(resp.StatusCode).To(Equal(200))

			content, _ := io.ReadAll(resp.Body)
			Expect(content).To(Equal(expectedContent))
		})
	})

	Describe("PutFile", func() {
		It("General error", func() {
			contentBody := "this is the new fake content"
			req := httptest.NewRequest("GET", "/wopi/files/abcdef/contents", strings.NewReader(contentBody))
			req.Header.Set(connector.HeaderWopiLock, "abc123")

			w := httptest.NewRecorder()

			cc.On("PutFile", mock.Anything, mock.Anything, int64(len(contentBody)), "abc123").Times(1).Return("", errors.New("Something happened"))

			httpAdapter.PutFile(w, req)
			resp := w.Result()
			Expect(resp.StatusCode).To(Equal(500))
		})

		It("Conflict", func() {
			contentBody := "this is the new fake content"
			req := httptest.NewRequest("GET", "/wopi/files/abcdef/contents", strings.NewReader(contentBody))
			req.Header.Set(connector.HeaderWopiLock, "abc123")

			w := httptest.NewRecorder()

			cc.On("PutFile", mock.Anything, mock.Anything, int64(len(contentBody)), "abc123").Times(1).Return("zzz111", connector.NewConnectorError(409, "Lock conflict"))

			httpAdapter.PutFile(w, req)
			resp := w.Result()
			Expect(resp.StatusCode).To(Equal(409))
			Expect(resp.Header.Get(connector.HeaderWopiLock)).To(Equal("zzz111"))
		})

		It("Success", func() {
			contentBody := "this is the new fake content"
			req := httptest.NewRequest("GET", "/wopi/files/abcdef/contents", strings.NewReader(contentBody))
			req.Header.Set(connector.HeaderWopiLock, "abc123")

			w := httptest.NewRecorder()

			cc.On("PutFile", mock.Anything, mock.Anything, int64(len(contentBody)), "abc123").Times(1).Return("", nil)

			httpAdapter.PutFile(w, req)
			resp := w.Result()
			Expect(resp.StatusCode).To(Equal(200))
		})
	})
})
