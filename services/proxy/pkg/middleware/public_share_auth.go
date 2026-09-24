package middleware

import (
	"net/http"
	"strings"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/reva/v2/pkg/auth/manager/publicshares"
	"github.com/owncloud/reva/v2/pkg/rgrpc/todo/pool"
)

const (
	_headerRevaAccessToken  = "x-access-token"
	headerShareToken        = "public-token"
	basicAuthPasswordPrefix = "password|"
	authenticationType      = "publicshares"

	_paramSignature  = "signature"
	_paramExpiration = "expiration"
)

// PublicShareAuthenticator is the authenticator which can authenticate public share requests.
// It will add the share owner into the request context.
type PublicShareAuthenticator struct {
	Logger              log.Logger
	RevaGatewaySelector pool.Selectable[gateway.GatewayAPIClient]
}

// The archiver is able to create archives from public shares in which case it needs to use the
// PublicShareAuthenticator. It might however also be called using "normal" authentication or
// using signed url, which are handled by other middleware. For this reason we can't just
// handle `/archiver` with the `isPublicPath()` check.
func isPublicShareArchive(r *http.Request) bool {
	if strings.HasPrefix(r.URL.Path, "/archiver") {
		if r.URL.Query().Get(headerShareToken) != "" || r.Header.Get(headerShareToken) != "" {
			return true
		}
	}
	return false
}

// The app open requests can be made in public share contexts. For that the PublicShareAuthenticator needs to
// augment the request context.
// The app open requests can also be made in authenticated context. In these cases the PublicShareAuthenticator
// needs to ignore the request.
func isPublicShareAppOpen(r *http.Request) bool {
	return (strings.HasPrefix(r.URL.Path, "/app/open") || strings.HasPrefix(r.URL.Path, "/app/new")) &&
		(r.URL.Query().Get(headerShareToken) != "" || r.Header.Get(headerShareToken) != "")
}

// The public-files requests can also be made in authenticated context. In these cases the OIDCAuthenticator and
// the BasicAuthenticator needs to ignore the request when the headerShareToken exist.
func isPublicWithShareToken(r *http.Request) bool {
	return (strings.HasPrefix(r.URL.Path, "/dav/public-files") || strings.HasPrefix(r.URL.Path, "/remote.php/dav/public-files")) &&
		(r.URL.Query().Get(headerShareToken) != "" || r.Header.Get(headerShareToken) != "")
}

// Authenticate implements the authenticator interface to authenticate requests via public share auth.
func (a PublicShareAuthenticator) Authenticate(r *http.Request) (*http.Request, error) {
	if !isPublicPath(r.URL.Path) && !isPublicShareArchive(r) && !isPublicShareAppOpen(r) {
		return nil, ErrAuthenticationFailed
	}

	query := r.URL.Query()
	shareToken := r.Header.Get(headerShareToken)
	if shareToken == "" {
		shareToken = query.Get(headerShareToken)
	}

	if shareToken == "" {
		// If the share token is not set then we don't need to inject the user to
		// the request context so we can just continue with the request.
		return r, nil
	}

	var sharePassword string
	var signatureAuth bool
	if signature := query.Get(_paramSignature); signature != "" {
		signatureAuth = true
		expiration := query.Get(_paramExpiration)
		if expiration == "" {
			a.Logger.Warn().Str("signature", signature).Msg("cannot do signature auth without the expiration")
			return nil, ErrAuthenticationFailed
		}
		sharePassword = strings.Join([]string{"signature", signature, expiration}, "|")
	} else {
		// We can ignore the username since it is always set to "public" in public shares.
		_, password, ok := r.BasicAuth()

		sharePassword = basicAuthPasswordPrefix
		if ok {
			sharePassword += password
		}
	}

	client, err := a.RevaGatewaySelector.Next()
	if err != nil {
		a.Logger.Error().
			Err(err).
			Str("authenticator", "public_share").
			Str("public_share_token", shareToken).
			Str("path", r.URL.Path).
			Msg("could not select next gateway client")
		return nil, ErrAuthenticationFailed
	}

	// The brute force protection counts failed public link password attempts.
	// Skip counting here for requests that are re-authenticated downstream (the
	// public-files DAV paths, which ocdav authenticates again) so the same
	// attempt is not counted twice, and for signature based auth, which is a
	// different, non-guessable credential. Password attempts on the paths where
	// the proxy is the sole authenticator (e.g. /archiver, /app/open) must be
	// counted, otherwise the throttle would never trip.
	ctx := r.Context()
	if signatureAuth || isPublicWithShareToken(r) {
		ctx = publicshares.MarkSkipAttemptContext(ctx, shareToken)
	}
	authResp, err := client.Authenticate(ctx, &gateway.AuthenticateRequest{
		Type:         authenticationType,
		ClientId:     shareToken,
		ClientSecret: sharePassword,
	})

	if err != nil {
		a.Logger.Error().
			Err(err).
			Str("authenticator", "public_share").
			Str("public_share_token", shareToken).
			Str("path", r.URL.Path).
			Msg("failed to authenticate request")
		return nil, ErrAuthenticationFailed
	}

	r.Header.Add(_headerRevaAccessToken, authResp.Token)

	a.Logger.Debug().
		Str("authenticator", "public_share").
		Str("path", r.URL.Path).
		Msg("successfully authenticated request")
	return r, nil
}
