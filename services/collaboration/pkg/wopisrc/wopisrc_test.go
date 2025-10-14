package wopisrc_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/owncloud/ocis/v2/services/collaboration/pkg/config"
	"github.com/owncloud/ocis/v2/services/collaboration/pkg/wopisrc"
)

var _ = Describe("Wopisrc Test", func() {
	var (
		c *config.Config
	)

	Context("GenerateWopiSrc", func() {
		BeforeEach(func() {
			c = &config.Config{
				Wopi: config.Wopi{
					WopiSrc:     "https://ocis.team/wopi/files",
					ProxyURL:    "https://cloud.proxy.com",
					ProxySecret: "secret",
				},
			}
		})
		When("WopiSrc URL is incorrect", func() {
			c = &config.Config{
				Wopi: config.Wopi{
					WopiSrc: "https:&//ocis.team/wopi/files",
				},
			}
			url, err := wopisrc.GenerateWopiSrc("123456", c)
			Expect(err).To(HaveOccurred())
			Expect(url).To(BeNil())
		})
		When("proxy URL is incorrect", func() {
			c = &config.Config{
				Wopi: config.Wopi{
					WopiSrc:     "https://ocis.team/wopi/files",
					ProxyURL:    "cloud",
					ProxySecret: "secret",
				},
			}
			url, err := wopisrc.GenerateWopiSrc("123456", c)
			Expect(err).To(HaveOccurred())
			Expect(url).To(BeNil())
		})
		When("proxy URL and proxy secret are configured", func() {
			It("should generate a WOPI src URL as a jwt token", func() {
				url, err := wopisrc.GenerateWopiSrc("123456", c)
				Expect(err).ToNot(HaveOccurred())
				Expect(url.String()).To(Equal("https://cloud.proxy.com/wopi/files/eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1IjoiaHR0cHM6Ly9vY2lzLnRlYW0vd29waS9maWxlcy8iLCJmIjoiMTIzNDU2In0.6ol9PQXGKktKfAri8tsJ4X_a9rIeosJ7id6KTQW6Ui0"))
			})
		})
		When("proxy URL and proxy secret are not configured", func() {
			It("should generate a WOPI src URL as a direct URL", func() {
				c.Wopi.ProxyURL = ""
				c.Wopi.ProxySecret = ""
				url, err := wopisrc.GenerateWopiSrc("123456", c)
				Expect(err).ToNot(HaveOccurred())
				Expect(url.String()).To(Equal("https://ocis.team/wopi/files/123456"))
			})
		})
	})
})
