package middleware

import (
	"context"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/google/uuid"
	accounts "github.com/owncloud/ocis/accounts/pkg/proto/v0"
	"github.com/owncloud/ocis/ocis-pkg/log"
	ocisoidc "github.com/owncloud/ocis/ocis-pkg/oidc"
	"github.com/owncloud/ocis/proxy/pkg/config"
	store "github.com/owncloud/ocis/store/pkg/proto/v0"
	"golang.org/x/crypto/pbkdf2"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func SignedURLAuth(optionSetters ...Option) func(next http.Handler) http.Handler {
	options := newOptions(optionSetters...)

	return func(next http.Handler) http.Handler {
		return &signedURLAuth{
			next:               next,
			logger:             options.Logger,
			preSignedURLConfig: options.PreSignedURLConfig,
			accountsClient:     options.AccountsClient,
			store:              options.Store,
		}
	}
}

type signedURLAuth struct {
	next               http.Handler
	logger             log.Logger
	preSignedURLConfig config.PreSignedURL
	accountsClient     accounts.AccountsService
	store              store.StoreService
}

func (m signedURLAuth) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if !m.shouldServe(req) {
		m.next.ServeHTTP(w, req)
		return
	}

	if err := m.validate(req); err != nil {
		http.Error(w, "Invalid url signature", http.StatusUnauthorized)
		return
	}

	claims, err := m.claims(req.URL.Query().Get("OC-Credential"))
	if err != nil {
		http.Error(w, "Invalid url signature", http.StatusUnauthorized)
		return
	}

	m.next.ServeHTTP(w, req.WithContext(ocisoidc.NewContext(req.Context(), claims)))
}

func (m signedURLAuth) claims(credential string) (*ocisoidc.StandardClaims, error) {
	// use openid claims to let the account_uuid middleware do a lookup by username
	claims := ocisoidc.StandardClaims{
		OcisID: credential,
	}

	// OC10 username is handled as id, if we get a credantial that is not of type uuid we expect
	// that it is a PreferredUsername und we need to get the corresponding uuid
	if _, err := uuid.Parse(claims.OcisID); err != nil {
		// todo caching
		account, status := getAccount(
			m.logger,
			m.accountsClient,
			fmt.Sprintf(
				"preferred_name eq '%s'",
				strings.ReplaceAll(
					claims.OcisID,
					"'",
					"''",
				),
			),
		)

		if status != 0 || account == nil {
			return nil, fmt.Errorf("no oc-credential found for %v", claims.OcisID)
		}

		claims.OcisID = account.Id
	}

	return &claims, nil
}

func (m signedURLAuth) shouldServe(req *http.Request) bool {
	return req.URL.Query().Get("OC-Signature") != ""
}

func (m signedURLAuth) validate(req *http.Request) (err error) {
	query := req.URL.Query()

	if ok, err := m.allRequiredParametersArePresent(query); !ok {
		return err
	}

	if ok, err := m.requestMethodMatches(req.Method, query); !ok {
		return err
	}

	if ok, err := m.requestMethodIsAllowed(req.Method); !ok {
		return err
	}

	if expired, err := m.urlIsExpired(query, time.Now); expired {
		return err
	}

	if ok, err := m.signatureIsValid(req); !ok {
		return err
	}

	return nil
}

func (m signedURLAuth) allRequiredParametersArePresent(query url.Values) (ok bool, err error) {
	// check if required query parameters exist in given request query parameters
	// OC-Signature - the computed signature - server will verify the request upon this REQUIRED
	// OC-Credential - defines the user scope (shall we use the owncloud user id here - this might leak internal data ....) REQUIRED
	// OC-Date - defined the date the url was signed (ISO 8601 UTC) REQUIRED
	// OC-Expires - defines the expiry interval in seconds (between 1 and 604800 = 7 days) REQUIRED
	// TODO OC-Verb - defines for which http verb the request is valid - defaults to GET OPTIONAL
	for _, p := range []string{
		"OC-Signature",
		"OC-Credential",
		"OC-Date",
		"OC-Expires",
		"OC-Verb",
	} {
		if query.Get(p) == "" {
			return false, fmt.Errorf("required %s parameter not found", p)
		}
	}

	return true, nil
}

func (m signedURLAuth) requestMethodMatches(meth string, query url.Values) (ok bool, err error) {
	// check if given url query parameter OC-Verb matches given request method
	if !strings.EqualFold(meth, query.Get("OC-Verb")) {
		return false, errors.New("required OC-Verb parameter did not match request method")
	}

	return true, nil
}

func (m signedURLAuth) requestMethodIsAllowed(meth string) (ok bool, err error) {
	//  check if given request method is allowed
	methodIsAllowed := false
	for _, am := range m.preSignedURLConfig.AllowedHTTPMethods {
		if strings.EqualFold(meth, am) {
			methodIsAllowed = true
			break
		}
	}

	if !methodIsAllowed {
		return false, errors.New("request method is not listed in PreSignedURLConfig AllowedHTTPMethods")
	}

	return true, nil
}
func (m signedURLAuth) urlIsExpired(query url.Values, now func() time.Time) (expired bool, err error) {
	// check if url is expired by checking if given date (OC-Date) + expires in seconds (OC-Expires) is after now
	date, err := time.Parse(time.RFC3339, query.Get("OC-Date"))
	if err != nil {
		return true, err
	}

	expires, err := time.ParseDuration(query.Get("OC-Expires") + "s")
	if err != nil {
		return true, err
	}

	date.Add(expires)

	return date.After(now()), nil
}

func (m signedURLAuth) signatureIsValid(req *http.Request) (ok bool, err error) {
	signingKey, err := m.getSigningKey(req.Context(), req.URL.Query().Get("OC-Credential"))
	if err != nil {
		m.logger.Error().Err(err).Msg("could not retrieve signing key")
		return false, err
	}
	if len(signingKey) == 0 {
		m.logger.Error().Err(err).Msg("signing key empty")
		return false, err
	}
	q := req.URL.Query()
	signature := q.Get("OC-Signature")
	q.Del("OC-Signature")
	req.URL.RawQuery = q.Encode()
	url := req.URL.String()
	if !req.URL.IsAbs() {
		url = "https://" + req.Host + url // TODO where do we get the scheme from
	}

	return m.createSignature(url, signingKey) == signature, nil
}

func (m signedURLAuth) createSignature(url string, signingKey []byte) string {
	// the oc10 signature check: $hash = \hash_pbkdf2("sha512", $url, $signingKey, 10000, 64, false);
	// - sets the length of the output string to 64
	// - sets raw output to false ->  if raw_output is FALSE length corresponds to twice the byte-length of the derived key (as every byte of the key is returned as two hexits).
	// TODO change to length 128 in oc10?
	// fo golangs pbkdf2.Key we need to use 32 because it will be encoded into 64 hexits later
	hash := pbkdf2.Key([]byte(url), signingKey, 10000, 32, sha512.New)
	return hex.EncodeToString(hash)
}

func (m signedURLAuth) getSigningKey(ctx context.Context, credential string) ([]byte, error) {
	claims, err := m.claims(credential)
	if err != nil {
		return []byte{}, err
	}

	res, err := m.store.Read(ctx, &store.ReadRequest{
		Options: &store.ReadOptions{
			Database: "proxy",
			Table:    "signing-keys",
		},
		Key: claims.OcisID,
	})
	if err != nil || len(res.Records) < 1 {
		return []byte{}, err
	}

	return res.Records[0].Value, nil
}
