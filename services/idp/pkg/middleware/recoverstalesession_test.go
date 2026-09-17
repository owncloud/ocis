package middleware

import (
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/services/idp/pkg/config"
)

var _ = Describe("recover stale session", func() {
	cfg := &config.Config{
		IDP: config.Settings{Iss: "https://cloud.example.test"},
	}

	newHandler := func(nextCalled *bool) http.Handler {
		middle := RecoverStaleSession(cfg, log.NopLogger())
		return middle(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			*nextCalled = true
			w.WriteHeader(http.StatusOK)
		}))
	}

	It("breaks the authorize error loop by clearing the logon cookie and restarting login", func() {
		req, _ := http.NewRequest(http.MethodGet,
			"https://cloud.example.test/signin/v1/identifier/_/authorize?error=unsupported_response_type", nil)
		req.AddCookie(&http.Cookie{Name: logonCookieName, Value: "stale-deleted-user-token"})

		nextCalled := false
		rr := httptest.NewRecorder()
		newHandler(&nextCalled).ServeHTTP(rr, req)

		Expect(nextCalled).To(BeFalse())
		Expect(rr).To(HaveHTTPStatus(http.StatusFound))
		Expect(rr.Header().Get("Location")).To(Equal(cfg.IDP.Iss))

		// the stale logon cookie must be expired so the browser drops it
		var cleared *http.Cookie
		for _, c := range rr.Result().Cookies() {
			if c.Name == logonCookieName {
				cleared = c
			}
		}
		Expect(cleared).NotTo(BeNil())
		Expect(cleared.MaxAge).To(BeNumerically("<", 0))
	})

	It("lets a regular authorization request through", func() {
		req, _ := http.NewRequest(http.MethodGet,
			"https://cloud.example.test/signin/v1/identifier/_/authorize?response_type=code&client_id=web", nil)

		nextCalled := false
		rr := httptest.NewRecorder()
		newHandler(&nextCalled).ServeHTTP(rr, req)

		Expect(nextCalled).To(BeTrue())
		Expect(rr).To(HaveHTTPStatus(http.StatusOK))
	})

	It("ignores an error parameter on a non-authorize route", func() {
		req, _ := http.NewRequest(http.MethodGet,
			"https://cloud.example.test/signin/v1/identifier?error=unsupported_response_type", nil)

		nextCalled := false
		rr := httptest.NewRecorder()
		newHandler(&nextCalled).ServeHTTP(rr, req)

		Expect(nextCalled).To(BeTrue())
		Expect(rr).To(HaveHTTPStatus(http.StatusOK))
	})
})
