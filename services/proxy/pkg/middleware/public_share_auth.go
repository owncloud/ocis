package middleware

import (
	"net/http"
	"strings"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
)

const (
	headerRevaAccessToken   = "x-access-token"
	headerShareToken        = "public-token"
	basicAuthPasswordPrefix = "password|"
	authenticationType      = "publicshares"
)

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

	// Currently we only want to authenticate app open request coming from public shares.
	if shareToken == "" {
		// Don't authenticate
		return nil, false
	}

	var sharePassword string
	if signature := query.Get("signature"); signature != "" {
		expiration := query.Get("expiration")
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
		a.Logger.Debug().Err(err).Str("public_share_token", shareToken).Msg("could not authenticate public share")
		// try another middleware
		return nil, false
	}

	r.Header.Add(headerRevaAccessToken, authResp.Token)
	return r, false
}
