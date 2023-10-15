package content_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"

	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"

	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	conf "github.com/owncloud/ocis/v2/services/search/pkg/config/defaults"
	"github.com/owncloud/ocis/v2/services/search/pkg/content"
	contentMocks "github.com/owncloud/ocis/v2/services/search/pkg/content/mocks"
)

var _ = Describe("Tika", func() {
	Describe("extract", func() {
		var (
			body         string
			fullResponse string
			language     string
			version      string
			srv          *httptest.Server
			tika         *content.Tika
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
				case "/language/string":
					out = language
				case "/rmeta/text":
					if fullResponse != "" {
						out = fullResponse
					} else {
						out = fmt.Sprintf(`[{"X-TIKA:content":"%s"}]`, body)
					}
				}

				_, _ = w.Write([]byte(out))
			}))

			cfg := conf.DefaultConfig()
			cfg.Extractor.Tika.TikaURL = srv.URL
			cfg.Extractor.Tika.CleanStopWords = true

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
				Size: 1,
			})
			Expect(err).ToNot(HaveOccurred())
			Expect(doc.Content).To(Equal(body))
		})

		It("adds audio content", func() {
			fullResponse = `[
				{
					"xmpDM:genre": "Some Genre",
					"xmpDM:album": "Some Album",
					"xmpDM:trackNumber": "7",
					"xmpDM:discNumber": "4",
					"xmpDM:releaseDate": "2004",
					"xmpDM:artist": "Some Artist",
					"xmpDM:albumArtist": "Some AlbumArtist",
					"xmpDM:audioCompressor": "MP3",
					"xmpDM:audioChannelType": "Stereo",
					"version": "MPEG 3 Layer III Version 1",
					"xmpDM:logComment": "some comment",
					"xmpDM:audioSampleRate": "44100",
					"channels": "2",
					"dc:title": "Some Title",
					"xmpDM:duration": "225",
					"Content-Type": "audio/mpeg",
					"samplerate": "44100"
				}
			]`
			doc, err := tika.Extract(context.TODO(), &provider.ResourceInfo{
				Type: provider.ResourceType_RESOURCE_TYPE_FILE,
				Size: 1,
			})
			Expect(err).ToNot(HaveOccurred())
			Expect(doc.Audio).ToNot(BeNil())
			Expect(*doc.Audio.Album).To(Equal("Some Album"))
			Expect(*doc.Audio.AlbumArtist).To(Equal("Some AlbumArtist"))
			Expect(*doc.Audio.Artist).To(Equal("Some Artist"))
			// Expect(*doc.Audio.Bitrate).To(Equal())
			// Expect(*doc.Audio.Composers).To(Equal())
			// Expect(*doc.Audio.Copyright).To(Equal())
			Expect(*doc.Audio.Disc).To(Equal(int32(4)))
			// Expect(*doc.Audio.DiscCount).To(Equal())
			Expect(*doc.Audio.Duration).To(Equal(int64(225000)))
			Expect(*doc.Audio.Genre).To(Equal("Some Genre"))
			// Expect(*doc.Audio.HasDrm).To(Equal())
			// Expect(*doc.Audio.IsVariableBitrate).To(Equal())
			Expect(*doc.Audio.Title).To(Equal("Some Title"))
			Expect(*doc.Audio.Track).To(Equal(int32(7)))
			// Expect(*doc.Audio.TrackCount).To(Equal())
			Expect(*doc.Audio.Year).To(Equal(int32(2004)))
		})

		It("removes stop words", func() {
			body = "body to test stop words!!! against almost everyone"
			language = "en"

			doc, err := tika.Extract(context.TODO(), &provider.ResourceInfo{
				Type: provider.ResourceType_RESOURCE_TYPE_FILE,
				Size: 1,
			})
			Expect(err).ToNot(HaveOccurred())
			Expect(doc.Content).To(Equal("body test stop words!!!"))
		})

		It("keeps stop words", func() {
			body = "body to test stop words!!! against almost everyone"
			language = "en"

			tika.CleanStopWords = false
			doc, err := tika.Extract(context.TODO(), &provider.ResourceInfo{
				Type: provider.ResourceType_RESOURCE_TYPE_FILE,
				Size: 1,
			})
			Expect(err).ToNot(HaveOccurred())
			Expect(doc.Content).To(Equal(body))
		})
	})
})
