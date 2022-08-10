package content_test

import (
	"context"
	storageProvider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	cs3Types "github.com/cs3org/go-cs3apis/cs3/types/v1beta1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/services/search/pkg/content"
)

var _ = Describe("Basic", func() {
	var (
		basic  content.Extractor
		logger = log.NewLogger()
		ctx    = context.TODO()
	)

	BeforeEach(func() {
		basic, _ = content.NewBasicExtractor(logger)
	})

	Describe("extract", func() {
		It("basic fields", func() {
			ri := &storageProvider.ResourceInfo{
				Path:     "./foo/bar.pdf",
				Size:     1024,
				MimeType: "application/pdf",
			}

			doc, err := basic.Extract(ctx, ri)

			Expect(err).To(BeNil())
			Expect(doc).ToNot(BeNil())
			Expect(doc.Name).To(Equal(ri.Path))
			Expect(doc.Size).To(Equal(ri.Size))
			Expect(doc.MimeType).To(Equal(ri.MimeType))
		})

		It("adds tags", func() {
			for _, data := range []struct {
				tags   string
				expect []string
			}{
				{tags: "", expect: []string{}},
				{tags: ",,,", expect: []string{}},
				{tags: ",foo,,", expect: []string{"foo"}},
				{tags: ",foo,,bar,", expect: []string{"foo", "bar"}},
			} {
				ri := &storageProvider.ResourceInfo{
					ArbitraryMetadata: &storageProvider.ArbitraryMetadata{
						Metadata: map[string]string{
							"tags": data.tags,
						},
					},
				}

				doc, err := basic.Extract(ctx, ri)
				Expect(err).To(BeNil())
				Expect(doc).ToNot(BeNil())
				Expect(doc.Tags).To(Equal(data.expect))
			}
		})

		It("RFC3339 mtime", func() {
			for _, data := range []struct {
				second uint64
				expect string
			}{
				{second: 4000, expect: "1970-01-01T01:06:40Z"},
				{second: 3000, expect: "1970-01-01T00:50:00Z"},
				{expect: ""},
			} {
				ri := &storageProvider.ResourceInfo{}

				if data.second != 0 {
					ri.Mtime = &cs3Types.Timestamp{Seconds: data.second}
				}

				doc, err := basic.Extract(ctx, ri)
				Expect(err).To(BeNil())
				Expect(doc).ToNot(BeNil())
				Expect(doc.Mtime).To(Equal(data.expect))
			}
		})
	})
})
