package opa_test

import (
	"context"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/open-policy-agent/opa/rego"
	"github.com/owncloud/ocis/v2/services/policies/pkg/engine/opa"
)

var _ = Describe("opa ocis mimetype functions", func() {
	Describe("ocis.mimetype.detect", func() {
		It("detects the mimetype", func() {
			r := rego.New(rego.Query(`ocis.mimetype.detect("")`), opa.RFMimetypeDetect)
			rs, err := r.Eval(context.Background())
			Expect(err).ToNot(HaveOccurred())
			Expect(rs[0].Expressions[0].String()).To(Equal("text/plain"))
		})
	})

	Describe("ocis.mimetype.extension_for_mimetype", func() {
		It("provides matching extensions", func() {
			r := rego.New(rego.Query(`ocis.mimetype.extensions("text/plain")`), opa.RFMimetypeExtensions)
			rs, err := r.Eval(context.Background())
			Expect(err).ToNot(HaveOccurred())
			Expect(rs[0].Expressions[0].String()).To(Equal("[.conf .def .in .list .log .text .txt]"))
		})
	})
})
