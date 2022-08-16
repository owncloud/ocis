package middleware

import (
	"net/http"
	"strings"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
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
	Logger            log.Logger
	RevaGatewayClient gateway.GatewayAPIClient
}

func (a PublicShareAuthenticator) Authenticate(r *http.Request) (*http.Request, bool) {
	if !isPublicPath(r.URL.Path) {
		return nil, false
	}

	query := r.URL.Query()
	shareToken := r.Header.Get(headerShareToken)
	if shareToken == "" {
		shareToken = query.Get(headerShareToken)
	}

	if shareToken == "" {
		// If the share token is not set then we don't need to inject the user to
		// the request context so we can just continue with the request.
		return r, true
	}

	var sharePassword string
	if signature := query.Get(_paramSignature); signature != "" {
		expiration := query.Get(_paramExpiration)
		if expiration == "" {
			a.Logger.Warn().Str("signature", signature).Msg("cannot do signature auth without the expiration")
			return nil, false
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

	authResp, err := a.RevaGatewayClient.Authenticate(r.Context(), &gateway.AuthenticateRequest{
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
		return nil, false
	}

	r.Header.Add(_headerRevaAccessToken, authResp.Token)

	a.Logger.Debug().
		Str("authenticator", "public_share").
		Str("path", r.URL.Path).
		Msg("successfully authenticated request")
	return r, true
}
