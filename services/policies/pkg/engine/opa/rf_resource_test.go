package opa_test

import (
	"context"
	"encoding/base64"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/open-policy-agent/opa/rego"

	"github.com/owncloud/ocis/v2/services/policies/pkg/engine/opa"
)

var _ = Describe("opa ocis resource functions", func() {
	Describe("ocis.resource.download", func() {
		It("downloads reva resources", func() {
			ts := []byte("Lorem Ipsum is simply dummy text of the printing and typesetting")
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write(ts)
			}))
			defer srv.Close()

			r := rego.New(rego.Query(`ocis.resource.download("`+srv.URL+`")`), opa.RFResourceDownload)
			rs, err := r.Eval(context.Background())
			Expect(err).ToNot(HaveOccurred())

			data, err := base64.StdEncoding.DecodeString(rs[0].Expressions[0].String())
			Expect(err).ToNot(HaveOccurred())

			Expect(data).To(Equal(ts))

		})
	})
})
