package wopisrc

import (
	"errors"
	"net/url"
	"path"

	"github.com/golang-jwt/jwt/v5"
	"github.com/owncloud/ocis/v2/services/collaboration/pkg/config"
)

// GenerateWopiSrc generates a WOPI src URL for the given file reference.
// If a proxy URL and proxy secret are configured, the URL will be generated
// as a jwt token that is signed with the proxy secret and contains the file reference
// and the WOPI src URL.
// Example:
// https://cloud.proxy.com/wopi/files/eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1IjoiaHR0cHM6Ly9vY2lzLnRlYW0vd29waS9maWxlcy8iLCJmIjoiMTIzNDU2In0.6ol9PQXGKktKfAri8tsJ4X_a9rIeosJ7id6KTQW6Ui0
//
// If no proxy URL and proxy secret are configured, the URL will be generated
// as a direct URL that contains the file reference.
// Example:
// https:/ocis.team/wopi/files/12312678470610632091729803710923
func GenerateWopiSrc(fileRef string, cfg *config.Config) (*url.URL, error) {
	wopiSrcURL, err := url.Parse(cfg.Wopi.WopiSrc)
	if err != nil {
		return nil, err
	}
	if wopiSrcURL.Host == "" {
		return nil, errors.New("invalid WopiSrc URL")
	}

	if cfg.Wopi.ProxyURL != "" && cfg.Wopi.ProxySecret != "" {
		return generateProxySrc(fileRef, cfg.Wopi.ProxyURL, cfg.Wopi.ProxySecret, wopiSrcURL)
	}

	return generateDirectSrc(fileRef, wopiSrcURL)
}

func generateDirectSrc(fileRef string, wopiSrcURL *url.URL) (*url.URL, error) {
	wopiSrcURL.Path = path.Join("wopi", "files", fileRef)
	return wopiSrcURL, nil
}

func generateProxySrc(fileRef string, proxyUrl string, proxySecret string, wopiSrcURL *url.URL) (*url.URL, error) {
	proxyURL, err := url.Parse(proxyUrl)
	if err != nil {
		return nil, err
	}
	if proxyURL.Host == "" {
		return nil, errors.New("invalid proxy URL")
	}

	wopiSrcURL.Path = path.Join("wopi", "files")

	type tokenClaims struct {
		URL    string `json:"u"`
		FileID string `json:"f"`
		jwt.RegisteredClaims
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims{
		FileID: fileRef,
		// the string value from the URL package always ends with a slash
		// the office365 proxy assumes that we have a trailing slash
		URL: wopiSrcURL.String() + "/",
	})
	tokenString, err := token.SignedString([]byte(proxySecret))
	if err != nil {
		return nil, err
	}
	proxyURL.Path = path.Join("wopi", "files", tokenString)
	return proxyURL, nil
}
