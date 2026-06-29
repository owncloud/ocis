package middleware

import (
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/services/idp/pkg/config"
)

var _ = Describe("check redirect", func() {
	clients := []config.Client{
		// just set the ID and RedirectURIs, no need for more at the moment
		config.Client{
			ID:           "web1",
			RedirectURIs: []string{"https://good.server.prv/", "https://good.server.prv/oidc-callback.html"},
		},
		config.Client{
			ID:           "web2",
			RedirectURIs: []string{"https://another.server.prv:6767/", "https://another.server.prv:6767/oidc-callback.html"},
		},
		config.Client{
			ID:           "custom_scheme1",
			RedirectURIs: []string{"oc://custom.server/"},
		},
		config.Client{
			ID:           "custom_scheme2",
			RedirectURIs: []string{"ios://ios.custom.server/"},
		},
		config.Client{
			ID:           "localhost1",
			RedirectURIs: []string{"http://localhost/"},
		},
		config.Client{
			ID:           "with_opaque",
			RedirectURIs: []string{"app:open?path=/personal/folder"},
		},
	}
	cfg := &config.Config{
		Clients: clients,
	}

	DescribeTable("is URL allowed",
		func(input string, status int) {
			req, _ := http.NewRequest(http.MethodGet, "https://demo.hidden.missing.server.prv/", nil)
			values := req.URL.Query()
			values.Add("redirect_uri", input)
			req.URL.RawQuery = values.Encode()

			rr := httptest.NewRecorder()
			middle := CheckRedirect(cfg, log.NopLogger())
			handler := middle(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}))
			handler.ServeHTTP(rr, req)
			Expect(rr).To(HaveHTTPStatus(status))
		},
		Entry("https web", "https://good.server.prv/oidc-callback.html", http.StatusOK),
		Entry("https web port", "https://another.server.prv:6767/oidc-callback.html", http.StatusOK),
		Entry("http web", "http://good.server.prv/oidc-callback.html", http.StatusInternalServerError),
		Entry("wrong web", "https://nonexisting.server.prv/", http.StatusInternalServerError),
		Entry("wrong web path", "https://nonexisting.server.prv/very-wrong", http.StatusInternalServerError),
		Entry("wrong web port", "https://good.server.prv:12345/oidc-callback.html", http.StatusInternalServerError),
		Entry("android", "oc://custom.server/", http.StatusOK),
		Entry("ios", "ios://ios.custom.server/", http.StatusOK),
		Entry("localhost", "http://localhost/", http.StatusOK),
		Entry("localhost port", "http://localhost:51515/", http.StatusOK),
		Entry("localhost wrong scheme", "https://localhost:51515/", http.StatusInternalServerError),
		Entry("opaque", "app:open?path=/personal/folder", http.StatusOK),
	)
})
