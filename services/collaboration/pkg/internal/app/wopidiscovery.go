package app

import (
	"context"
	"crypto/tls"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/beevik/etree"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/pkg/errors"
)

func (app *DemoApp) WopiDiscovery(ctx context.Context) error {
	res, err := getAppURLs(app.Config.WopiApp.Addr, app.Config.WopiApp.Insecure, app.Logger)
	if err != nil {
		// logging is already covered inside the `getAppURLs` function
		return err
	}

	app.AppURLs = res
	return nil
}

func getAppURLs(wopiAppUrl string, insecure bool, logger log.Logger) (map[string]map[string]string, error) {

	wopiAppUrl = wopiAppUrl + "/hosting/discovery"

	httpClient := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: insecure,
			},
		},
	}

	httpResp, err := httpClient.Get(wopiAppUrl)
	if err != nil {
		logger.Error().
			Err(err).
			Str("WopiAppUrl", wopiAppUrl).
			Msg("WopiDiscovery: failed to access wopi app url")
		return nil, err
	}

	if httpResp.StatusCode != http.StatusOK {
		logger.Error().
			Str("WopiAppUrl", wopiAppUrl).
			Int("HttpCode", httpResp.StatusCode).
			Msg("WopiDiscovery: wopi app url failed with unexpected code")
		return nil, errors.New("status code was not 200")
	}

	defer httpResp.Body.Close()

	var appURLs map[string]map[string]string

	appURLs, err = parseWopiDiscovery(httpResp.Body)
	if err != nil {
		logger.Error().
			Err(err).
			Str("WopiAppUrl", wopiAppUrl).
			Msg("WopiDiscovery: failed to parse wopi discovery response")
		return nil, errors.Wrap(err, "error parsing wopi discovery response")
	}

	// TODO: Log appUrls? not easy with the format
	// It's also a one-time call during service setup, so it's pointless
	// to use an "all-is-good" debug log
	return appURLs, nil
}

func parseWopiDiscovery(body io.Reader) (map[string]map[string]string, error) {
	appURLs := make(map[string]map[string]string)

	doc := etree.NewDocument()
	if _, err := doc.ReadFrom(body); err != nil {
		return nil, err
	}
	root := doc.SelectElement("wopi-discovery")

	for _, netzone := range root.SelectElements("net-zone") {

		if strings.Contains(netzone.SelectAttrValue("name", ""), "external") {
			for _, app := range netzone.SelectElements("app") {
				for _, action := range app.SelectElements("action") {
					access := action.SelectAttrValue("name", "")
					if access == "view" || access == "edit" {
						ext := action.SelectAttrValue("ext", "")
						urlString := action.SelectAttrValue("urlsrc", "")

						if ext == "" || urlString == "" {
							continue
						}

						u, err := url.Parse(urlString)
						if err != nil {
							continue
						}

						// remove any malformed query parameter from discovery urls
						q := u.Query()
						for k := range q {
							if strings.Contains(k, "<") || strings.Contains(k, ">") {
								q.Del(k)
							}
						}

						u.RawQuery = q.Encode()

						if _, ok := appURLs[access]; !ok {
							appURLs[access] = make(map[string]string)
						}
						appURLs[access]["."+ext] = u.String()
					}
				}
			}
		}
	}
	return appURLs, nil
}
