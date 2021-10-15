package middleware

import (
	"context"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	revactx "github.com/cs3org/reva/pkg/ctx"
	"github.com/owncloud/ocis/ocis-pkg/log"
	"github.com/owncloud/ocis/proxy/pkg/config"
	"github.com/owncloud/ocis/proxy/pkg/user/backend"
	store "github.com/owncloud/ocis/store/pkg/proto/v0"
	"golang.org/x/crypto/pbkdf2"
)

// SignedURLAuth provides a middleware to check access secured by a signed URL.
func SignedURLAuth(optionSetters ...Option) func(next http.Handler) http.Handler {
	options := newOptions(optionSetters...)

	return func(next http.Handler) http.Handler {
		return &signedURLAuth{
			next:               next,
			logger:             options.Logger,
			preSignedURLConfig: options.PreSignedURLConfig,
			store:              options.Store,
			userProvider:       options.UserProvider,
		}
	}
}

type signedURLAuth struct {
	next               http.Handler
	logger             log.Logger
	preSignedURLConfig config.PreSignedURL
	userProvider       backend.UserBackend
	store              store.StoreService
}

func (m signedURLAuth) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if !m.shouldServe(req) {
		m.next.ServeHTTP(w, req)
		return
	}

	user, _, err := m.userProvider.GetUserByClaims(req.Context(), "username", req.URL.Query().Get("OC-Credential"), true)
	if err != nil {
		m.logger.Error().Err(err).Msg("Could not get user by claim")
		w.WriteHeader(http.StatusInternalServerError)
	}

	ctx := revactx.ContextSetUser(req.Context(), user)

	req = req.WithContext(ctx)

	if err := m.validate(req); err != nil {
		http.Error(w, "Invalid url signature", http.StatusUnauthorized)
		return
	}

	m.next.ServeHTTP(w, req)
}

func (m signedURLAuth) shouldServe(req *http.Request) bool {
	if !m.preSignedURLConfig.Enabled {
		return false
	}
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
	validFrom, err := time.Parse(time.RFC3339, query.Get("OC-Date"))
	if err != nil {
		return true, err
	}

	requestExpiry, err := time.ParseDuration(query.Get("OC-Expires") + "s")
	if err != nil {
		return true, err
	}

	validTo := validFrom.Add(requestExpiry)

	return !(now().After(validFrom) && now().Before(validTo)), nil
}

func (m signedURLAuth) signatureIsValid(req *http.Request) (ok bool, err error) {
	u := revactx.ContextMustGetUser(req.Context())
	signingKey, err := m.getSigningKey(req.Context(), u.Id.OpaqueId)
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

func (m signedURLAuth) getSigningKey(ctx context.Context, ocisID string) ([]byte, error) {
	res, err := m.store.Read(ctx, &store.ReadRequest{
		Options: &store.ReadOptions{
			Database: "proxy",
			Table:    "signing-keys",
		},
		Key: ocisID,
	})
	if err != nil || len(res.Records) < 1 {
		return []byte{}, err
	}

	return res.Records[0].Value, nil
}
