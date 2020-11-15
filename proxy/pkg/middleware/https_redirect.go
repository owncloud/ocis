package middleware

import (
	"fmt"
	"net/http"
)

// RedirectToHTTPS redirects insecure requests to https
func HTTPsRedirect(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		proto := req.Header.Get("x-forwarded-proto")
		if proto == "http" || proto == "HTTP" {
			http.Redirect(res, req, fmt.Sprintf("https://%s%s", req.Host, req.URL), http.StatusPermanentRedirect)
			return
		}

		next.ServeHTTP(res, req)
	})
}
