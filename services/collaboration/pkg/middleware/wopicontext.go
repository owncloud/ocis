package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	appproviderv1beta1 "github.com/cs3org/go-cs3apis/cs3/app/provider/v1beta1"
	userv1beta1 "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	providerv1beta1 "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	ctxpkg "github.com/cs3org/reva/v2/pkg/ctx"
	"github.com/golang-jwt/jwt/v4"
	"github.com/rs/zerolog"
	"google.golang.org/grpc/metadata"
)

type key int

const (
	wopiContextKey key = iota
)

// WopiContext wraps all the information we need for WOPI
type WopiContext struct {
	AccessToken   string
	ViewOnlyToken string
	FileReference *providerv1beta1.Reference
	User          *userv1beta1.User
	ViewMode      appproviderv1beta1.ViewMode
	EditAppUrl    string
	ViewAppUrl    string
}

// WopiContextAuthMiddleware will prepare an HTTP handler to be used as
// middleware. The handler will create a WopiContext by parsing the
// access_token (which must be provided as part of the URL query).
// The access_token is required.
//
// This middleware will add the following to the request's context:
// * The access token as metadata for outgoing requests (for the
// authentication against the CS3 API, the "x-access-token" header).
// * The created WopiContext for the request
// * A contextual zerologger containing information about the request
// and the WopiContext
func WopiContextAuthMiddleware(jwtSecret string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessToken := r.URL.Query().Get("access_token")
		if accessToken == "" {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		claims := &Claims{}
		_, err := jwt.ParseWithClaims(accessToken, claims, func(token *jwt.Token) (interface{}, error) {

			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(jwtSecret), nil
		})

		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		if err := claims.Valid(); err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		ctx := r.Context()

		wopiContextAccessToken, err := DecryptAES([]byte(jwtSecret), claims.WopiContext.AccessToken)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		claims.WopiContext.AccessToken = wopiContextAccessToken

		ctx = context.WithValue(ctx, wopiContextKey, claims.WopiContext)
		// authentication for the CS3 api
		ctx = metadata.AppendToOutgoingContext(ctx, ctxpkg.TokenHeader, claims.WopiContext.AccessToken)

		// include additional info in the context's logger
		// we might need to check https://learn.microsoft.com/en-us/microsoft-365/cloud-storage-partner-program/rest/common-headers
		// although some headers might not be sent depending on the client.
		logger := zerolog.Ctx(ctx)
		ctx = logger.With().
			Str("WopiOverride", r.Header.Get("X-WOPI-Override")).
			Str("FileReference", claims.WopiContext.FileReference.String()).
			Str("ViewMode", claims.WopiContext.ViewMode.String()).
			Str("Requester", claims.WopiContext.User.GetId().String()).
			Logger().WithContext(ctx)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Extract a WopiContext from the context if possible. An error will be
// returned if there is no WopiContext
func WopiContextFromCtx(ctx context.Context) (WopiContext, error) {
	if wopiContext, ok := ctx.Value(wopiContextKey).(WopiContext); ok {
		return wopiContext, nil
	}
	return WopiContext{}, errors.New("no wopi context found")
}

// Set a WopiContext in the context. A new context will be returned with the
// add WopiContext inside. Note that the old one won't have the WopiContext set.
//
// This method is used for testing. The WopiContextAuthMiddleware should be
// used instead in order to provide a valid WopiContext
func WopiContextToCtx(ctx context.Context, wopiContext WopiContext) context.Context {
	return context.WithValue(ctx, wopiContextKey, wopiContext)
}
