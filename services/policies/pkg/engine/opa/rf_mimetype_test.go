package opa_test

import (
	"context"
	"io"
	"strings"

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
	Describe("ocis.mimetype.extensions", func() {
		DescribeTable("resolves extensions by mimetype",
			func(mimetype string, expectations []string, f io.Reader) {
				rfMimetypeExtensions, err := opa.RFMimetypeExtensions(f)
				Expect(err).ToNot(HaveOccurred())

				r := rego.New(rego.Query(`ocis.mimetype.extensions("`+mimetype+`")`), rfMimetypeExtensions)
				rs, err := r.Eval(context.Background())
				Expect(err).ToNot(HaveOccurred())

				got := rs[0].Expressions[0].String()

				if len(expectations) == 0 {
					Expect(got).To(Equal("[]"))
				}

				for i, expectation := range expectations {
					if i+1 != len(expectations) {
						expectation += " "
					}

					Expect(string(got[0])).To(Equal("["))
					Expect(strings.Contains(got, expectation)).To(BeTrue())
					Expect(string(got[len(got)-1])).To(Equal("]"))
				}
			},
			Entry("With default mimetype", "application/pdf", []string{".pdf"}, nil),
			Entry("With unknown mimetype", "ocis/with.custom.mt", []string{}, nil),
			Entry("With custom mimetype", "ocis/with.custom.mt", []string{".with.custom.mt"}, strings.NewReader("ocis/with.custom.mt    with.custom.mt")),
			Entry("With multiple custom mimetypes", "ocis/with.multiple.custom.mt", []string{".with.multiple.custom.1.mt", ".with.multiple.custom.2.mt"}, strings.NewReader("ocis/with.multiple.custom.mt                                with.multiple.custom.1.mt with.multiple.custom.2.mt")),
			Entry("With custom ignored mimetype", "ocis/with.multiple.custom.ignored.mt", []string{}, strings.NewReader("#ocis/with.multiple.custom.ignored.mt with.multiple.custom.ignored.mt")),
		)
	})
})
