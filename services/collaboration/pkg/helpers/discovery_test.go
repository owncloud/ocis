package helpers_test

import (
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/services/collaboration/pkg/config"
	"github.com/owncloud/ocis/v2/services/collaboration/pkg/helpers"
)

var _ = Describe("Discovery", func() {
	var (
		discoveryContent1 string
		srv               *httptest.Server
	)

	BeforeEach(func() {
		discoveryContent1 = `
<?xml version="1.0" encoding="utf-8"?>
<wopi-discovery>
  <net-zone name="external-http">
    <app name="Word" favIconUrl="https://test.server.prv/web-apps/apps/documenteditor/main/resources/img/favicon.ico">
      <action name="view" ext="pdf" urlsrc="https://test.server.prv/hosting/wopi/word/view?&amp;&lt;rs=DC_LLCC&amp;&gt;&lt;dchat=DISABLE_CHAT&amp;&gt;&lt;embed=EMBEDDED&amp;&gt;&lt;fs=FULLSCREEN&amp;&gt;&lt;hid=HOST_SESSION_ID&amp;&gt;&lt;rec=RECORDING&amp;&gt;&lt;sc=SESSION_CONTEXT&amp;&gt;&lt;thm=THEME_ID&amp;&gt;&lt;ui=UI_LLCC&amp;&gt;&lt;wopisrc=WOPI_SOURCE&amp;&gt;&amp;"/>
      <action name="embedview" ext="pdf" urlsrc="https://test.server.prv/hosting/wopi/word/view?embed=1&amp;&lt;rs=DC_LLCC&amp;&gt;&lt;dchat=DISABLE_CHAT&amp;&gt;&lt;embed=EMBEDDED&amp;&gt;&lt;fs=FULLSCREEN&amp;&gt;&lt;hid=HOST_SESSION_ID&amp;&gt;&lt;rec=RECORDING&amp;&gt;&lt;sc=SESSION_CONTEXT&amp;&gt;&lt;thm=THEME_ID&amp;&gt;&lt;ui=UI_LLCC&amp;&gt;&lt;wopisrc=WOPI_SOURCE&amp;&gt;&amp;"/>
      <action name="view" ext="djvu" urlsrc="https://test.server.prv/hosting/wopi/word/view?&amp;&lt;rs=DC_LLCC&amp;&gt;&lt;dchat=DISABLE_CHAT&amp;&gt;&lt;embed=EMBEDDED&amp;&gt;&lt;fs=FULLSCREEN&amp;&gt;&lt;hid=HOST_SESSION_ID&amp;&gt;&lt;rec=RECORDING&amp;&gt;&lt;sc=SESSION_CONTEXT&amp;&gt;&lt;thm=THEME_ID&amp;&gt;&lt;ui=UI_LLCC&amp;&gt;&lt;wopisrc=WOPI_SOURCE&amp;&gt;&amp;"/>
      <action name="embedview" ext="djvu" urlsrc="https://test.server.prv/hosting/wopi/word/view?embed=1&amp;&lt;rs=DC_LLCC&amp;&gt;&lt;dchat=DISABLE_CHAT&amp;&gt;&lt;embed=EMBEDDED&amp;&gt;&lt;fs=FULLSCREEN&amp;&gt;&lt;hid=HOST_SESSION_ID&amp;&gt;&lt;rec=RECORDING&amp;&gt;&lt;sc=SESSION_CONTEXT&amp;&gt;&lt;thm=THEME_ID&amp;&gt;&lt;ui=UI_LLCC&amp;&gt;&lt;wopisrc=WOPI_SOURCE&amp;&gt;&amp;"/>
      <action name="view" ext="docx" urlsrc="https://test.server.prv/hosting/wopi/word/view?&amp;&lt;rs=DC_LLCC&amp;&gt;&lt;dchat=DISABLE_CHAT&amp;&gt;&lt;embed=EMBEDDED&amp;&gt;&lt;fs=FULLSCREEN&amp;&gt;&lt;hid=HOST_SESSION_ID&amp;&gt;&lt;rec=RECORDING&amp;&gt;&lt;sc=SESSION_CONTEXT&amp;&gt;&lt;thm=THEME_ID&amp;&gt;&lt;ui=UI_LLCC&amp;&gt;&lt;wopisrc=WOPI_SOURCE&amp;&gt;&amp;"/>
      <action name="embedview" ext="docx" urlsrc="https://test.server.prv/hosting/wopi/word/view?embed=1&amp;&lt;rs=DC_LLCC&amp;&gt;&lt;dchat=DISABLE_CHAT&amp;&gt;&lt;embed=EMBEDDED&amp;&gt;&lt;fs=FULLSCREEN&amp;&gt;&lt;hid=HOST_SESSION_ID&amp;&gt;&lt;rec=RECORDING&amp;&gt;&lt;sc=SESSION_CONTEXT&amp;&gt;&lt;thm=THEME_ID&amp;&gt;&lt;ui=UI_LLCC&amp;&gt;&lt;wopisrc=WOPI_SOURCE&amp;&gt;&amp;"/>
      <action name="editnew" ext="docx" requires="locks,update" urlsrc="https://test.server.prv/hosting/wopi/word/edit?&amp;&lt;rs=DC_LLCC&amp;&gt;&lt;dchat=DISABLE_CHAT&amp;&gt;&lt;embed=EMBEDDED&amp;&gt;&lt;fs=FULLSCREEN&amp;&gt;&lt;hid=HOST_SESSION_ID&amp;&gt;&lt;rec=RECORDING&amp;&gt;&lt;sc=SESSION_CONTEXT&amp;&gt;&lt;thm=THEME_ID&amp;&gt;&lt;ui=UI_LLCC&amp;&gt;&lt;wopisrc=WOPI_SOURCE&amp;&gt;&amp;"/>
      <action name="edit" ext="docx" default="true" requires="locks,update" urlsrc="https://test.server.prv/hosting/wopi/word/edit?&amp;&lt;rs=DC_LLCC&amp;&gt;&lt;dchat=DISABLE_CHAT&amp;&gt;&lt;embed=EMBEDDED&amp;&gt;&lt;fs=FULLSCREEN&amp;&gt;&lt;hid=HOST_SESSION_ID&amp;&gt;&lt;rec=RECORDING&amp;&gt;&lt;sc=SESSION_CONTEXT&amp;&gt;&lt;thm=THEME_ID&amp;&gt;&lt;ui=UI_LLCC&amp;&gt;&lt;wopisrc=WOPI_SOURCE&amp;&gt;&amp;"/>
    </app>
    <app name="Excel" favIconUrl="https://test.server.prv/web-apps/apps/spreadsheeteditor/main/resources/img/favicon.ico">
      <action name="view" ext="xls" urlsrc="https://test.server.prv/hosting/wopi/cell/view?&amp;&lt;rs=DC_LLCC&amp;&gt;&lt;dchat=DISABLE_CHAT&amp;&gt;&lt;embed=EMBEDDED&amp;&gt;&lt;fs=FULLSCREEN&amp;&gt;&lt;hid=HOST_SESSION_ID&amp;&gt;&lt;rec=RECORDING&amp;&gt;&lt;sc=SESSION_CONTEXT&amp;&gt;&lt;thm=THEME_ID&amp;&gt;&lt;ui=UI_LLCC&amp;&gt;&lt;wopisrc=WOPI_SOURCE&amp;&gt;&amp;"/>
      <action name="embedview" ext="xls" urlsrc="https://test.server.prv/hosting/wopi/cell/view?embed=1&amp;&lt;rs=DC_LLCC&amp;&gt;&lt;dchat=DISABLE_CHAT&amp;&gt;&lt;embed=EMBEDDED&amp;&gt;&lt;fs=FULLSCREEN&amp;&gt;&lt;hid=HOST_SESSION_ID&amp;&gt;&lt;rec=RECORDING&amp;&gt;&lt;sc=SESSION_CONTEXT&amp;&gt;&lt;thm=THEME_ID&amp;&gt;&lt;ui=UI_LLCC&amp;&gt;&lt;wopisrc=WOPI_SOURCE&amp;&gt;&amp;"/>
      <action name="convert" ext="xls" targetext="xlsx" requires="update" urlsrc="https://test.server.prv/hosting/wopi/convert-and-edit/xls/xlsx?&amp;&lt;rs=DC_LLCC&amp;&gt;&lt;dchat=DISABLE_CHAT&amp;&gt;&lt;embed=EMBEDDED&amp;&gt;&lt;fs=FULLSCREEN&amp;&gt;&lt;hid=HOST_SESSION_ID&amp;&gt;&lt;rec=RECORDING&amp;&gt;&lt;sc=SESSION_CONTEXT&amp;&gt;&lt;thm=THEME_ID&amp;&gt;&lt;ui=UI_LLCC&amp;&gt;&lt;wopisrc=WOPI_SOURCE&amp;&gt;&amp;"/>
      <action name="view" ext="xlsb" urlsrc="https://test.server.prv/hosting/wopi/cell/view?&amp;&lt;rs=DC_LLCC&amp;&gt;&lt;dchat=DISABLE_CHAT&amp;&gt;&lt;embed=EMBEDDED&amp;&gt;&lt;fs=FULLSCREEN&amp;&gt;&lt;hid=HOST_SESSION_ID&amp;&gt;&lt;rec=RECORDING&amp;&gt;&lt;sc=SESSION_CONTEXT&amp;&gt;&lt;thm=THEME_ID&amp;&gt;&lt;ui=UI_LLCC&amp;&gt;&lt;wopisrc=WOPI_SOURCE&amp;&gt;&amp;"/>
      <action name="embedview" ext="xlsb" urlsrc="https://test.server.prv/hosting/wopi/cell/view?embed=1&amp;&lt;rs=DC_LLCC&amp;&gt;&lt;dchat=DISABLE_CHAT&amp;&gt;&lt;embed=EMBEDDED&amp;&gt;&lt;fs=FULLSCREEN&amp;&gt;&lt;hid=HOST_SESSION_ID&amp;&gt;&lt;rec=RECORDING&amp;&gt;&lt;sc=SESSION_CONTEXT&amp;&gt;&lt;thm=THEME_ID&amp;&gt;&lt;ui=UI_LLCC&amp;&gt;&lt;wopisrc=WOPI_SOURCE&amp;&gt;&amp;"/>
      <action name="convert" ext="xlsb" targetext="xlsx" requires="update" urlsrc="https://test.server.prv/hosting/wopi/convert-and-edit/xlsb/xlsx?&amp;&lt;rs=DC_LLCC&amp;&gt;&lt;dchat=DISABLE_CHAT&amp;&gt;&lt;embed=EMBEDDED&amp;&gt;&lt;fs=FULLSCREEN&amp;&gt;&lt;hid=HOST_SESSION_ID&amp;&gt;&lt;rec=RECORDING&amp;&gt;&lt;sc=SESSION_CONTEXT&amp;&gt;&lt;thm=THEME_ID&amp;&gt;&lt;ui=UI_LLCC&amp;&gt;&lt;wopisrc=WOPI_SOURCE&amp;&gt;&amp;"/>
    </app>
    <app name="application/vnd.oasis.opendocument.presentation">
      <action name="edit" ext="" default="true" requires="locks,update" urlsrc="https://test.server.prv/hosting/wopi/slide/edit?&amp;&lt;rs=DC_LLCC&amp;&gt;&lt;dchat=DISABLE_CHAT&amp;&gt;&lt;embed=EMBEDDED&amp;&gt;&lt;fs=FULLSCREEN&amp;&gt;&lt;hid=HOST_SESSION_ID&amp;&gt;&lt;rec=RECORDING&amp;&gt;&lt;sc=SESSION_CONTEXT&amp;&gt;&lt;thm=THEME_ID&amp;&gt;&lt;ui=UI_LLCC&amp;&gt;&lt;wopisrc=WOPI_SOURCE&amp;&gt;&amp;"/>
    </app>
  </net-zone>
  <proof-key oldvalue="BgIAAACkAABSU0ExAAgAAAEAAQD/NVqekFNi8X3p6Bvdlaxm0GGuggW5kKfVEQzPGuOkGVrz6DrOMNR+k7Pq8tONY+1NHgS6Z+v3959em78qclVDuQX77Tkml0xMHAQHN4sAHF9iQJS8gOBUKSVKaHD7Z8YXch6F212YSUSc8QphpDSHWVShU7rcUeLQsd/0pkflh5+um4YKEZhm4Mou3vstp5p12NeffyK1WFZF7q4jB7jclAslYKQsP82YY3DcRwu5Tl/+W0ifVcXze0mI7v1reJ12pKn8ifRiq+0q5oJST3TRSrvmjLg9Gt3ozhVIt2HUi3La7Qh40YOAUXm0g/hUq2BepeOp1C7WSvaOFHXe6Hqq" oldmodulus="qnro3nUUjvZK1i7UqeOlXmCrVPiDtHlRgIPReAjt2nKL1GG3SBXO6N0aPbiM5rtK0XRPUoLmKu2rYvSJ/Kmkdp14a/3uiEl788VVn0hb/l9OuQtH3HBjmM0/LKRgJQuU3LgHI67uRVZYtSJ/n9fYdZqnLfveLsrgZpgRCoabrp+H5Uem9N+x0OJR3LpToVRZhzSkYQrxnERJmF3bhR5yF8Zn+3BoSiUpVOCAvJRAYl8cAIs3BwQcTEyXJjnt+wW5Q1VyKr+bXp/39+tnugQeTe1jjdPy6rOTftQwzjro81oZpOMazwwR1aeQuQWCrmHQZqyV3Rvo6X3xYlOQnlo1/w==" oldexponent="AQAB" value="BgIAAACkAABSU0ExAAgAAAEAAQD/NVqekFNi8X3p6Bvdlaxm0GGuggW5kKfVEQzPGuOkGVrz6DrOMNR+k7Pq8tONY+1NHgS6Z+v3959em78qclVDuQX77Tkml0xMHAQHN4sAHF9iQJS8gOBUKSVKaHD7Z8YXch6F212YSUSc8QphpDSHWVShU7rcUeLQsd/0pkflh5+um4YKEZhm4Mou3vstp5p12NeffyK1WFZF7q4jB7jclAslYKQsP82YY3DcRwu5Tl/+W0ifVcXze0mI7v1reJ12pKn8ifRiq+0q5oJST3TRSrvmjLg9Gt3ozhVIt2HUi3La7Qh40YOAUXm0g/hUq2BepeOp1C7WSvaOFHXe6Hqq" modulus="qnro3nUUjvZK1i7UqeOlXmCrVPiDtHlRgIPReAjt2nKL1GG3SBXO6N0aPbiM5rtK0XRPUoLmKu2rYvSJ/Kmkdp14a/3uiEl788VVn0hb/l9OuQtH3HBjmM0/LKRgJQuU3LgHI67uRVZYtSJ/n9fYdZqnLfveLsrgZpgRCoabrp+H5Uem9N+x0OJR3LpToVRZhzSkYQrxnERJmF3bhR5yF8Zn+3BoSiUpVOCAvJRAYl8cAIs3BwQcTEyXJjnt+wW5Q1VyKr+bXp/39+tnugQeTe1jjdPy6rOTftQwzjro81oZpOMazwwR1aeQuQWCrmHQZqyV3Rvo6X3xYlOQnlo1/w==" exponent="AQAB"/>
</wopi-discovery>
`
		srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			switch req.URL.Path {
			case "/bad/hosting/discovery":
				w.WriteHeader(500)
			case "/good/hosting/discovery":
				w.Write([]byte(discoveryContent1))
			case "/wrongformat/hosting/discovery":
				w.Write([]byte("Text that <can't> be XML /form<atted/"))
			}
		}))
	})

	AfterEach(func() {
		srv.Close()
	})

	Describe("GetAppURLs", func() {
		It("Good discovery URL", func() {
			cfg := &config.Config{
				App: config.App{
					Addr:     srv.URL + "/good",
					Insecure: true,
				},
			}
			logger := log.NopLogger()

			appUrls, err := helpers.GetAppURLs(cfg, logger)

			expectedAppUrls := map[string]map[string]string{
				"view": map[string]string{
					".pdf":  "https://test.server.prv/hosting/wopi/word/view",
					".djvu": "https://test.server.prv/hosting/wopi/word/view",
					".docx": "https://test.server.prv/hosting/wopi/word/view",
					".xls":  "https://test.server.prv/hosting/wopi/cell/view",
					".xlsb": "https://test.server.prv/hosting/wopi/cell/view",
				},
				"edit": map[string]string{
					".docx": "https://test.server.prv/hosting/wopi/word/edit",
				},
			}

			Expect(err).To(Succeed())
			Expect(appUrls).To(Equal(expectedAppUrls))
		})

		It("Wrong discovery URL", func() {
			cfg := &config.Config{
				App: config.App{
					Addr:     srv.URL + "/bad",
					Insecure: true,
				},
			}
			logger := log.NopLogger()

			appUrls, err := helpers.GetAppURLs(cfg, logger)
			Expect(err).To(HaveOccurred())
			Expect(appUrls).To(BeNil())
		})

		It("Not XML formatted", func() {
			cfg := &config.Config{
				App: config.App{
					Addr:     srv.URL + "/wrongformat",
					Insecure: true,
				},
			}
			logger := log.NopLogger()

			appUrls, err := helpers.GetAppURLs(cfg, logger)
			Expect(err).To(HaveOccurred())
			Expect(appUrls).To(BeNil())
		})
	})
})
