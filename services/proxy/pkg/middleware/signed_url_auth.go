package middleware

import (
	"context"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/ocis-pkg/mfa"
	"github.com/owncloud/ocis/v2/services/proxy/pkg/config"
	"github.com/owncloud/ocis/v2/services/proxy/pkg/user/backend"
	"github.com/owncloud/ocis/v2/services/proxy/pkg/userroles"
	revactx "github.com/owncloud/reva/v2/pkg/ctx"
	microstore "go-micro.dev/v4/store"
	"golang.org/x/crypto/pbkdf2"
)

const (
	_paramOCSignature  = "OC-Signature"
	_paramOCCredential = "OC-Credential" // #nosec G101
	_paramOCDate       = "OC-Date"
	_paramOCExpires    = "OC-Expires"
	_paramOCVerb       = "OC-Verb"
	_paramOCAlgo       = "OC-Algo"
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

// mfaKeyBreadcrumbKey is the context key for the request-scoped breadcrumb
// that marks a request whose signed-URL signature matched the MFA-bound
// signing key (as opposed to the normal key). It is distinct from the
// final MFA flag (revactx.HasMFA): it is set by this authenticator, which
// runs before the MFA middleware's own context reset (revactx.RemoveMFA),
// and is read afterwards by that middleware to decide whether a
// claims-less request may inherit the stored per-user MFA status. It
// cannot be set by an inbound client request, since it is only ever
// derived from a successful signature match against the MFA-bound key.
type mfaKeyBreadcrumbKey struct{}

// setMFAKeyBreadcrumb marks the context as authenticated via the MFA-bound
// signing key.
func setMFAKeyBreadcrumb(ctx context.Context) context.Context {
	return context.WithValue(ctx, mfaKeyBreadcrumbKey{}, true)
}

// hasMFAKeyBreadcrumb reports whether the context carries the MFA-bound
// signing-key breadcrumb set by setMFAKeyBreadcrumb.
func hasMFAKeyBreadcrumb(ctx context.Context) bool {
	v, _ := ctx.Value(mfaKeyBreadcrumbKey{}).(bool)
	return v
}

// SignedURLAuthenticator is the authenticator responsible for authenticating signed URL requests.
type SignedURLAuthenticator struct {
	Logger             log.Logger
	PreSignedURLConfig config.PreSignedURL
	UserProvider       backend.UserBackend
	UserRoleAssigner   userroles.UserRoleAssigner
	Store              microstore.Store
	Now                func() time.Time
}

func (m SignedURLAuthenticator) shouldServe(req *http.Request) bool {
	if !m.PreSignedURLConfig.Enabled {
		return false
	}
	return req.URL.Query().Get(_paramOCSignature) != ""
}

func (m SignedURLAuthenticator) validate(req *http.Request) (err error) {
	query := req.URL.Query()

	if err := m.allRequiredParametersArePresent(query); err != nil {
		return err
	}

	if err := m.requestMethodMatches(req.Method, query); err != nil {
		return err
	}

	if err := m.requestMethodIsAllowed(req.Method); err != nil {
		return err
	}

	if err = m.urlIsExpired(query); err != nil {
		return err
	}

	if err := m.signatureIsValid(req); err != nil {
		return err
	}

	return nil
}

func (m SignedURLAuthenticator) allRequiredParametersArePresent(query url.Values) (err error) {
	// check if required query parameters exist in given request query parameters
	// OC-Signature - the computed signature - server will verify the request upon this REQUIRED
	// OC-Credential - defines the user scope (shall we use the owncloud user id here - this might leak internal data ....) REQUIRED
	// OC-Date - defined the date the url was signed (ISO 8601 UTC) REQUIRED
	// OC-Expires - defines the expiry interval in seconds (between 1 and 604800 = 7 days) REQUIRED
	// TODO OC-Verb - defines for which http verb the request is valid - defaults to GET OPTIONAL
	for _, p := range _requiredParams {
		if query.Get(p) == "" {
			return fmt.Errorf("required %s parameter not found", p)
		}
	}

	return nil
}

func (m SignedURLAuthenticator) requestMethodMatches(meth string, query url.Values) (err error) {
	// check if given url query parameter OC-Verb matches given request method
	if !strings.EqualFold(meth, query.Get(_paramOCVerb)) {
		return errors.New("required OC-Verb parameter did not match request method")
	}

	return nil
}

func (m SignedURLAuthenticator) requestMethodIsAllowed(meth string) (err error) {
	//  check if given request method is allowed
	methodIsAllowed := false
	for _, am := range m.PreSignedURLConfig.AllowedHTTPMethods {
		if strings.EqualFold(meth, am) {
			methodIsAllowed = true
			break
		}
	}

	if !methodIsAllowed {
		return errors.New("request method is not listed in PreSignedURLConfig AllowedHTTPMethods")
	}

	return nil
}

func (m SignedURLAuthenticator) urlIsExpired(query url.Values) (err error) {
	// check if url is expired by checking if given date (OC-Date) + expires in seconds (OC-Expires) is after now
	validFrom, err := time.Parse(time.RFC3339, query.Get(_paramOCDate))
	if err != nil {
		return err
	}

	requestExpiry, err := time.ParseDuration(query.Get(_paramOCExpires) + "s")
	if err != nil {
		return err
	}

	validTo := validFrom.Add(requestExpiry)
	if !(m.Now().Before(validTo)) {
		return errors.New("URL is expired")
	}

	return nil
}

// signatureIsValid tries the MFA-bound signing key first (including the
// $->%24 workaround on that attempt), then falls back to the normal key
// (with the same workaround) if the MFA-bound attempt does not match. A
// match against the MFA-bound key additionally sets a request-scoped
// breadcrumb; a match against the normal key does not. No match against
// either key fails authentication, as before this key was introduced.
func (m SignedURLAuthenticator) signatureIsValid(req *http.Request) (err error) {
	c := revactx.ContextMustGetUser(req.Context())
	userID := c.Id.OpaqueId
	u := m.buildUrlToSign(req)
	signatureInURL := req.URL.Query().Get(_paramOCSignature)

	// The MFA-bound key not existing for this user is the common case (no
	// MFA session has ever issued one) and is not itself an error; it only
	// means this particular signed URL was not produced with it. A real
	// store error, however, is logged so a broken store does not silently
	// look like "no MFA key".
	mfaKey, ok, rerr := m.readSigningKey(mfa.SigningKeyName(userID))
	if rerr != nil {
		m.Logger.Error().Err(rerr).Str("userid", userID).Msg("could not retrieve MFA-bound signing key")
	} else if ok && m.signatureMatches(u, signatureInURL, mfaKey) {
		*req = *req.WithContext(setMFAKeyBreadcrumb(req.Context()))
		return nil
	}

	signingKey, ok, err := m.readSigningKey(userID)
	if err != nil {
		m.Logger.Error().Err(err).Msg("could not retrieve signing key")
		return err
	}
	// A store may return an empty result set without an error (e.g. the "noop"
	// store, or a store that has no signing key for this user). Guard the index
	// access to avoid an out-of-range panic, and return a real error so the
	// request fails authentication instead of proceeding with an empty key.
	if !ok {
		m.Logger.Error().Str("userid", userID).Msg("signing key empty or not found")
		return errors.New("signing key not found")
	}
	if m.signatureMatches(u, signatureInURL, signingKey) {
		return nil
	}

	// Never put the computed signature into the error: it is a valid
	// credential for the URL the caller chose, and the error gets logged.
	return errors.New("signature mismatch")
}

// readSigningKey reads a signing-key record from the store. It treats
// "not found", and a store that returns an empty result set without an
// error (e.g. the "noop" store), as absence rather than an error.
func (m SignedURLAuthenticator) readSigningKey(key string) (value []byte, ok bool, err error) {
	records, err := m.Store.Read(key)
	if err != nil {
		if errors.Is(err, microstore.ErrNotFound) {
			return nil, false, nil
		}
		return nil, false, err
	}
	if len(records) == 0 || len(records[0].Value) == 0 {
		return nil, false, nil
	}
	return records[0].Value, true, nil
}

// signatureMatches checks the signature against the given key, including
// the $->%24 reverse-proxy workaround for https://github.com/owncloud/ocis/issues/10180.
func (m SignedURLAuthenticator) signatureMatches(u, signatureInURL string, signingKey []byte) bool {
	if hmac.Equal([]byte(m.createSignature(u, signingKey)), []byte(signatureInURL)) {
		return true
	}

	// try a workaround for https://github.com/owncloud/ocis/issues/10180
	// Some reverse proxies might replace $ with %24 in the URL leading to a mismatch in the signature
	uWorkaround := strings.Replace(u, "$", "%24", 1)
	return hmac.Equal([]byte(m.createSignature(uWorkaround, signingKey)), []byte(signatureInURL))
}

func (m SignedURLAuthenticator) buildUrlToSign(req *http.Request) string {
	q := req.URL.Query()

	// only params required for signing
	signParameters := make(url.Values)
	signParameters.Add(_paramOCCredential, q.Get(_paramOCCredential))
	signParameters.Add(_paramOCDate, q.Get(_paramOCDate))
	signParameters.Add(_paramOCExpires, q.Get(_paramOCExpires))
	signParameters.Add(_paramOCVerb, q.Get(_paramOCVerb))

	// remaining query params
	q.Del(_paramOCAlgo)
	q.Del(_paramOCCredential)
	q.Del(_paramOCDate)
	q.Del(_paramOCExpires)
	q.Del(_paramOCSignature)
	q.Del(_paramOCVerb)

	url := *req.URL
	if len(q) == 0 {
		url.RawQuery = signParameters.Encode()
	} else {
		url.RawQuery = strings.Join([]string{q.Encode(), signParameters.Encode()}, "&")
	}
	u := url.String()
	if !url.IsAbs() {
		u = "https://" + req.Host + u // TODO where do we get the scheme
	}
	return u
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

// Authenticate implements the authenticator interface to authenticate requests via signed URL auth.
func (m SignedURLAuthenticator) Authenticate(r *http.Request) (*http.Request, error) {
	if !m.shouldServe(r) {
		return nil, ErrAuthenticationFailed
	}

	user, _, err := m.UserProvider.GetUserByClaims(r.Context(), "username", r.URL.Query().Get(_paramOCCredential))
	if err != nil {
		m.Logger.Error().
			Err(err).
			Str("authenticator", "signed_url").
			Str("path", r.URL.Path).
			Msg("Could not get user by claim")
		return nil, ErrAuthenticationFailed
	}

	user, err = m.UserRoleAssigner.ApplyUserRole(r.Context(), user)
	if err != nil {
		m.Logger.Error().
			Err(err).
			Str("authenticator", "signed_url").
			Str("path", r.URL.Path).
			Msg("Could not get user by claim")
		return nil, ErrAuthenticationFailed
	}

	ctx := revactx.ContextSetUser(r.Context(), user)

	r = r.WithContext(ctx)

	if err := m.validate(r); err != nil {
		// Log the path only: the query carries the signature, which can still
		// be valid (e.g. when only the HTTP method did not match).
		m.Logger.Error().
			Err(err).
			Str("authenticator", "signed_url").
			Str("path", r.URL.Path).
			Msg("signed URL validation failed")
		return nil, ErrAuthenticationFailed
	}

	// TODO: set user in context
	m.Logger.Debug().
		Str("authenticator", "signed_url").
		Str("path", r.URL.Path).
		Msg("successfully authenticated request")
	return r, nil
}
