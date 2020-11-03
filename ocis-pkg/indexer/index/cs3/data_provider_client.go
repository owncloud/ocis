package cs3

import (
	"io"
	"net/http"
	"strings"
)

type dataProviderClient struct {
	client  http.Client
	baseURL string
}

func (d dataProviderClient) put(url string, body io.Reader, token string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPut, singleJoiningSlash(d.baseURL, url), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("x-access-token", token)
	return d.client.Do(req)
}

func (d dataProviderClient) get(url string, token string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, singleJoiningSlash(d.baseURL, url), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("x-access-token", token)
	return d.client.Do(req)
}

// TODO: this is copied from proxy. Find a better solution or move it to ocis-pkg
func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}
