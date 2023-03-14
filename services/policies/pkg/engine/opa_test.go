package engine_test

import (
	"context"
	"encoding/base64"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/open-policy-agent/opa/rego"
	"github.com/owncloud/ocis/v2/services/policies/pkg/engine"
)

var _ = Describe("Opa", func() {
	Describe("Custom OPA function", func() {
		Describe("GetResource", func() {
			It("loads reva resources", func() {
				ts := []byte("Lorem Ipsum is simply dummy text of the printing and typesetting")
				srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Write(ts)
				}))
				defer srv.Close()

				r := rego.New(rego.Query(`ocis_get_resource("`+srv.URL+`")`), engine.GetResource)
				rs, err := r.Eval(context.Background())
				Expect(err).ToNot(HaveOccurred())

				data, err := base64.StdEncoding.DecodeString(rs[0].Expressions[0].String())
				Expect(err).ToNot(HaveOccurred())

				Expect(data).To(Equal(ts))

			})
		})

		Describe("GetMimetype", func() {
			It("is defined and returns a mimetype", func() {
				r := rego.New(rego.Query(`ocis_get_mimetype("")`), engine.GetMimetype)
				rs, err := r.Eval(context.Background())
				Expect(err).ToNot(HaveOccurred())
				Expect(rs[0].Expressions[0].String()).To(Equal("text/plain"))
			})
		})
	})
})
