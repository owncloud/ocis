package middleware

import "net/http"

// M undocummented
type M func(next http.Handler) http.Handler
