package middleware

import (
	"net/http"
	"net/url"
	"time"

	"github.com/owncloud/ocis/v2/services/collaboration/pkg/config"
	"github.com/owncloud/ocis/v2/services/collaboration/pkg/proofkeys"
	"github.com/rs/zerolog"
)

// ProofKeysMiddleware will verify the proof keys of the requests.
// This is a middleware that could be disabled / not set.
//
// Requests will fail with a 500 HTTP status if the verification fails.
// As said, this can be disabled (via configuration) if you want to skip
// the verification.
// The middleware requires hitting the "/hosting/discovery" endpoint of the
// WOPI app in order to get the keys. The keys will be cached in memory for
// 12 hours (or the configured value) before hitting the endpoint again to
// request new / updated keys.
func ProofKeysMiddleware(cfg *config.Config, next http.Handler) http.Handler {
	wopiDiscovery := cfg.App.Addr + "/hosting/discovery"
	insecure := cfg.App.Insecure
	cacheDuration, err := time.ParseDuration(cfg.App.ProofKeys.Duration)
	if err != nil {
		cacheDuration = 12 * time.Hour
	}

	pkHandler := proofkeys.NewVerifyHandler(wopiDiscovery, insecure, cacheDuration)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := zerolog.Ctx(r.Context())

		// the url we need is the one being requested, but we need the
		// scheme and host, so we'll get those from the configured WOPISrc
		wopiSrcURL, _ := url.Parse(cfg.Wopi.WopiSrc)
		if cfg.Wopi.ProxyURL != "" {
			wopiSrcURL, _ = url.Parse(cfg.Wopi.ProxyURL)
		}
		currentURL, _ := url.Parse(r.URL.String())
		currentURL.Scheme = wopiSrcURL.Scheme
		currentURL.Host = wopiSrcURL.Host

		accessToken := r.URL.Query().Get("access_token")
		stamp := r.Header.Get("X-WOPI-TimeStamp")

		err := pkHandler.Verify(
			accessToken,
			currentURL.String(),
			stamp,
			r.Header.Get("X-WOPI-Proof"),
			r.Header.Get("X-WOPI-ProofOld"),
			proofkeys.VerifyWithLogger(logger),
		)

		if err != nil {
			logger.Error().Err(err).Msg("ProofKeys verification failed")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		logger.Debug().Msg("ProofKeys verified")

		next.ServeHTTP(w, r)
	})
}
