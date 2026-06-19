package middleware

import (
	"net/http"

	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/services/idp/pkg/config"
)

func CheckRedirect(cfg *config.Config, logger log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			uriString := r.URL.Query().Get("redirect_uri")
			if uriString == "" {
				next.ServeHTTP(w, r)
				return
			}

			for _, client := range cfg.Clients {
				for _, redirect_uri := range client.RedirectURIs {
					if uriString == redirect_uri {
						next.ServeHTTP(w, r)
						return
					}
				}
			}

			logger.Warn().Str("redirect_uri", uriString).Msg("Unknown redirect uri")
			http.Error(w, "Unknown URI", http.StatusInternalServerError)
		})
	}
}
