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

	revactx "github.com/cs3org/reva/v2/pkg/ctx"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	storemsg "github.com/owncloud/ocis/v2/protogen/gen/ocis/messages/store/v0"
	storesvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/store/v0"
	"github.com/owncloud/ocis/v2/services/proxy/pkg/config"
	"github.com/owncloud/ocis/v2/services/proxy/pkg/user/backend"
	"golang.org/x/crypto/pbkdf2"
)

const (
	_paramOCSignature  = "OC-Signature"
	_paramOCCredential = "OC-Credential"
	_paramOCDate       = "OC-Date"
	_paramOCExpires    = "OC-Expires"
	_paramOCVerb       = "OC-Verb"
)

var (
	_requiredParams = [...]string{
		_paramOCSignature,
		_paramOCCredential,
		_paramOCDate,
		_paramOCExpires,
		_paramOCVerb,
	}
)

// SignedURLAuthenticator is the authenticator responsible for authenticating signed URL requests.
type SignedURLAuthenticator struct {
	Logger             log.Logger
	PreSignedURLConfig config.PreSignedURL
	UserProvider       backend.UserBackend
	Store              storesvc.StoreService
}

func (m SignedURLAuthenticator) shouldServe(req *http.Request) bool {
	if !m.PreSignedURLConfig.Enabled {
		return false
	}
	return req.URL.Query().Get(_paramOCSignature) != ""
}

func (m SignedURLAuthenticator) validate(req *http.Request) (err error) {
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

func (m SignedURLAuthenticator) allRequiredParametersArePresent(query url.Values) (ok bool, err error) {
	// check if required query parameters exist in given request query parameters
	// OC-Signature - the computed signature - server will verify the request upon this REQUIRED
	// OC-Credential - defines the user scope (shall we use the owncloud user id here - this might leak internal data ....) REQUIRED
	// OC-Date - defined the date the url was signed (ISO 8601 UTC) REQUIRED
	// OC-Expires - defines the expiry interval in seconds (between 1 and 604800 = 7 days) REQUIRED
	// TODO OC-Verb - defines for which http verb the request is valid - defaults to GET OPTIONAL
	for _, p := range _requiredParams {
		if query.Get(p) == "" {
			return false, fmt.Errorf("required %s parameter not found", p)
		}
	}

	return true, nil
}

func (m SignedURLAuthenticator) requestMethodMatches(meth string, query url.Values) (ok bool, err error) {
	// check if given url query parameter OC-Verb matches given request method
	if !strings.EqualFold(meth, query.Get(_paramOCVerb)) {
		return false, errors.New("required OC-Verb parameter did not match request method")
	}

	return true, nil
}

func (m SignedURLAuthenticator) requestMethodIsAllowed(meth string) (ok bool, err error) {
	//  check if given request method is allowed
	methodIsAllowed := false
	for _, am := range m.PreSignedURLConfig.AllowedHTTPMethods {
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

func (m SignedURLAuthenticator) urlIsExpired(query url.Values, now func() time.Time) (expired bool, err error) {
	// check if url is expired by checking if given date (OC-Date) + expires in seconds (OC-Expires) is after now
	validFrom, err := time.Parse(time.RFC3339, query.Get(_paramOCDate))
	if err != nil {
		return true, err
	}

	requestExpiry, err := time.ParseDuration(query.Get(_paramOCExpires) + "s")
	if err != nil {
		return true, err
	}

	validTo := validFrom.Add(requestExpiry)

	return !(now().After(validFrom) && now().Before(validTo)), nil
}

func (m SignedURLAuthenticator) signatureIsValid(req *http.Request) (ok bool, err error) {
	u := revactx.ContextMustGetUser(req.Context())
	signingKey, err := m.getSigningKey(req.Context(), u.Id.OpaqueId)
	if err != nil {
		m.Logger.Error().Err(err).Msg("could not retrieve signing key")
		return false, err
	}
	if len(signingKey) == 0 {
		m.Logger.Error().Err(err).Msg("signing key empty")
		return false, err
	}
	q := req.URL.Query()
	signature := q.Get(_paramOCSignature)
	q.Del(_paramOCSignature)
	req.URL.RawQuery = q.Encode()
	url := req.URL.String()
	if !req.URL.IsAbs() {
		url = "https://" + req.Host + url // TODO where do we get the scheme from
	}

	return m.createSignature(url, signingKey) == signature, nil
}

func (m SignedURLAuthenticator) createSignature(url string, signingKey []byte) string {
	// the oc10 signature check: $hash = \hash_pbkdf2("sha512", $url, $signingKey, 10000, 64, false);
	// - sets the length of the output string to 64
	// - sets raw output to false ->  if raw_output is FALSE length corresponds to twice the byte-length of the derived key (as every byte of the key is returned as two hexits).
	// TODO change to length 128 in oc10?
	// fo golangs pbkdf2.Key we need to use 32 because it will be encoded into 64 hexits later
	hash := pbkdf2.Key([]byte(url), signingKey, 10000, 32, sha512.New)
	return hex.EncodeToString(hash)
}

func (m SignedURLAuthenticator) getSigningKey(ctx context.Context, ocisID string) ([]byte, error) {
	res, err := m.Store.Read(ctx, &storesvc.ReadRequest{
		Options: &storemsg.ReadOptions{
			Database: "proxy",
			Table:    "signing-keys",
		},
		Key: ocisID,
	})
	if err != nil || len(res.Records) < 1 {
		return nil, err
	}

	return res.Records[0].Value, nil
}

func (m SignedURLAuthenticator) Authenticate(r *http.Request) (*http.Request, bool) {
	if !m.shouldServe(r) {
		return nil, false
	}

	user, _, err := m.UserProvider.GetUserByClaims(r.Context(), "username", r.URL.Query().Get(_paramOCCredential), true)
	if err != nil {
		m.Logger.Error().
			Err(err).
			Str("authenticator", "signed_url").
			Str("path", r.URL.Path).
			Msg("Could not get user by claim")
		return nil, false
	}

	ctx := revactx.ContextSetUser(r.Context(), user)

	r = r.WithContext(ctx)

	if err := m.validate(r); err != nil {
		m.Logger.Error().
			Err(err).
			Str("authenticator", "signed_url").
			Str("path", r.URL.Path).
			Msg("Could not get user by claim")
		return nil, false
	}

	m.Logger.Debug().
		Str("authenticator", "signed_url").
		Str("path", r.URL.Path).
		Msg("successfully authenticated request")
	return r, true
}
