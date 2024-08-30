package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	appproviderv1beta1 "github.com/cs3org/go-cs3apis/cs3/app/provider/v1beta1"
	providerv1beta1 "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	ctxpkg "github.com/cs3org/reva/v2/pkg/ctx"
	rjwt "github.com/cs3org/reva/v2/pkg/token/manager/jwt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/owncloud/ocis/v2/services/collaboration/pkg/config"
	"github.com/owncloud/ocis/v2/services/collaboration/pkg/helpers"
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
	ViewMode      appproviderv1beta1.ViewMode
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
func WopiContextAuthMiddleware(cfg *config.Config, next http.Handler) http.Handler {
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

			return []byte(cfg.Wopi.Secret), nil
		})

		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		ctx := r.Context()

		wopiContextAccessToken, err := DecryptAES([]byte(cfg.Wopi.Secret), claims.WopiContext.AccessToken)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		tokenManager, err := rjwt.New(map[string]interface{}{
			"secret":  cfg.TokenManager.JWTSecret,
			"expires": int64(24 * 60 * 60),
		})
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		user, _, err := tokenManager.DismantleToken(ctx, wopiContextAccessToken)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		claims.WopiContext.AccessToken = wopiContextAccessToken

		ctx = context.WithValue(ctx, wopiContextKey, claims.WopiContext)
		// authentication for the CS3 api
		ctx = metadata.AppendToOutgoingContext(ctx, ctxpkg.TokenHeader, claims.WopiContext.AccessToken)
		ctx = ctxpkg.ContextSetUser(ctx, user)

		// include additional info in the context's logger
		// we might need to check https://learn.microsoft.com/en-us/microsoft-365/cloud-storage-partner-program/rest/common-headers
		// although some headers might not be sent depending on the client.
		logger := zerolog.Ctx(ctx)
		wopiLogger := logger.With().
			Str("WopiSessionId", r.Header.Get("X-WOPI-SessionId")).
			Str("WopiOverride", r.Header.Get("X-WOPI-Override")).
			Str("WopiProof", r.Header.Get("X-WOPI-Proof")).
			Str("WopiProofOld", r.Header.Get("X-WOPI-ProofOld")).
			Str("WopiStamp", r.Header.Get("X-WOPI-TimeStamp")).
			Str("FileReference", claims.WopiContext.FileReference.String()).
			Str("ViewMode", claims.WopiContext.ViewMode.String()).
			Str("Requester", user.GetId().String()).
			Logger()
		ctx = wopiLogger.WithContext(ctx)

		hashedRef := helpers.HashResourceId(claims.WopiContext.FileReference.GetResourceId())
		fileID := parseWopiFileID(cfg, r.URL.Path)
		if fileID != hashedRef {
			wopiLogger.Error().Msg("file reference in the URL doesn't match the one inside the access token")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

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

// The access token inside the wopiContext is expected to be decrypted.
// In order to generate the access token for WOPI, the reva token inside the
// wopiContext will be encrypted
func GenerateWopiToken(wopiContext WopiContext, cfg *config.Config) (string, int64, error) {
	cryptedReqAccessToken, err := EncryptAES([]byte(cfg.Wopi.Secret), wopiContext.AccessToken)
	if err != nil {
		return "", 0, err
	}

	cs3Claims := &jwt.RegisteredClaims{}
	cs3JWTparser := jwt.Parser{}
	_, _, err = cs3JWTparser.ParseUnverified(wopiContext.AccessToken, cs3Claims)
	if err != nil {
		return "", 0, err
	}

	wopiContext.AccessToken = cryptedReqAccessToken
	claims := &Claims{
		WopiContext: wopiContext,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: cs3Claims.ExpiresAt,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString([]byte(cfg.Wopi.Secret))

	return accessToken, claims.ExpiresAt.UnixMilli(), err
}

// parseWopiFileID extracts the file id from a wopi path
//
// If the file id is a jwt, it will be decoded and the file id will be extracted from the jwt claims.
// If the file id is not a jwt, it will be returned as is.
func parseWopiFileID(cfg *config.Config, path string) string {
	s := strings.Split(path, "/")
	if len(s) < 4 || (s[1] != "wopi" && s[2] != "files") {
		return path
	}
	// check if the fileid is a jwt
	if strings.Contains(s[3], ".") {
		token, err := jwt.Parse(s[3], func(_ *jwt.Token) (interface{}, error) {
			return []byte(cfg.Wopi.ProxySecret), nil
		})
		if err != nil {
			return s[3]
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return s[3]
		}

		f, ok := claims["f"].(string)
		if !ok {
			return s[3]
		}
		return f
	}
	// fileid is not a jwt
	return s[3]
}
