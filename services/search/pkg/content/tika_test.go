package content_test

import (
	"context"
	"fmt"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	conf "github.com/owncloud/ocis/v2/services/search/pkg/config/defaults"
	"github.com/owncloud/ocis/v2/services/search/pkg/content"
	contentMocks "github.com/owncloud/ocis/v2/services/search/pkg/content/mocks"
	"github.com/stretchr/testify/mock"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
)

var _ = Describe("Tika", func() {
	Describe("extract", func() {
		var (
			body     string
			language string
			version  string
			srv      *httptest.Server
			tika     *content.Tika
		)

		BeforeEach(func() {
			body = ""
			language = ""
			version = ""
			srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				out := ""
				switch req.URL.Path {
				case "/version":
					out = version
				case "/language/stream":
					out = language
				case "/rmeta/text":
					out = fmt.Sprintf(`[{"X-TIKA:content":"%s"}]`, body)
				}

				_, _ = w.Write([]byte(out))
			}))

			cfg := conf.DefaultConfig()
			cfg.Extractor.Tika.TikaURL = srv.URL

			var err error
			tika, err = content.NewTikaExtractor(nil, log.NewLogger(), cfg)
			Expect(err).ToNot(HaveOccurred())
			Expect(tika).ToNot(BeNil())

			retriever := &contentMocks.Retriever{}
			retriever.On("Retrieve", mock.Anything, mock.Anything, mock.Anything).Return(io.NopCloser(strings.NewReader(body)), nil)

			tika.Retriever = retriever
		})

		AfterEach(func() {
			srv.Close()
		})

		It("skips non file resources", func() {
			doc, err := tika.Extract(context.TODO(), &provider.ResourceInfo{})
			Expect(err).ToNot(HaveOccurred())
			Expect(doc.Content).To(Equal(""))
		})

		It("adds content", func() {
			body = "any body"

			doc, err := tika.Extract(context.TODO(), &provider.ResourceInfo{
				Type: provider.ResourceType_RESOURCE_TYPE_FILE,
			})
			Expect(err).ToNot(HaveOccurred())
			Expect(doc.Content).To(Equal(body))
		})

		It("removes stop words", func() {
			body = "body to test stop words!!! I, you, he, she, it, we, you, they, stay"
			language = "en"

			doc, err := tika.Extract(context.TODO(), &provider.ResourceInfo{
				Type: provider.ResourceType_RESOURCE_TYPE_FILE,
			})
			Expect(err).ToNot(HaveOccurred())
			Expect(doc.Content).To(Equal("body test stop words i stay"))
		})
	})
})
