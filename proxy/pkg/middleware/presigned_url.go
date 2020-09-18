package middleware

import (
	"context"
	"crypto/sha512"
	"encoding/hex"
	"net/http"
	"strings"
	"time"

	"github.com/owncloud/ocis-pkg/v2/log"
	ocisoidc "github.com/owncloud/ocis-pkg/v2/oidc"
	"github.com/owncloud/ocis-proxy/pkg/config"
	storepb "github.com/owncloud/ocis-store/pkg/proto/v0"
	"golang.org/x/crypto/pbkdf2"
)

const (
	iterations = 10000
	keyLen     = 32
)

// PresignedURL provides a middleware to check access secured by a presigned URL.
func PresignedURL(opts ...Option) func(next http.Handler) http.Handler {
	opt := newOptions(opts...)
	l := opt.Logger
	cfg := opt.PreSignedURLConfig

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if isSignedRequest(r) {
				if signedRequestIsValid(l, r, opt.Store, cfg) {
					// use openid claims to let the account_uuid middleware do a lookup by username
					claims := ocisoidc.StandardClaims{
						OcisID: r.URL.Query().Get("OC-Credential"),
					}

					// inject claims to the request context for the account_uuid middleware
					ctxWithClaims := ocisoidc.NewContext(r.Context(), &claims)
					r = r.WithContext(ctxWithClaims)

					next.ServeHTTP(w, r)
				} else {
					http.Error(w, "Invalid url signature", http.StatusUnauthorized)
					return
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}

func isSignedRequest(r *http.Request) bool {
	return r.URL.Query().Get("OC-Signature") != ""
}

func signedRequestIsValid(l log.Logger, r *http.Request, s storepb.StoreService, cfg config.PreSignedURL) bool {
	// TODO OC-Algorithm - defined the used algo (e.g. sha256 or sha512 - we should agree on one default algo and make this parameter optional)
	// TODO OC-Verb - defines for which http verb the request is valid - defaults to GET OPTIONAL

	return allRequiredParametersArePresent(r) &&
		requestMethodMatches(r) &&
		requestMethodIsAllowed(r.Method, cfg.AllowedHTTPMethods) &&
		!urlIsExpired(r, time.Now) &&
		signatureIsValid(l, r, s)
}

func allRequiredParametersArePresent(r *http.Request) bool {
	// OC-Credential - defines the user scope (shall we use the owncloud user id here - this might leak internal data ....) REQUIRED
	// OC-Date - defined the date the url was signed (ISO 8601 UTC) REQUIRED
	// OC-Expires - defines the expiry interval in seconds (between 1 and 604800 = 7 days) REQUIRED
	// OC-Signature - the computed signature - server will verify the request upon this REQUIRED
	return r.URL.Query().Get("OC-Signature") != "" &&
		r.URL.Query().Get("OC-Credential") != "" &&
		r.URL.Query().Get("OC-Date") != "" &&
		r.URL.Query().Get("OC-Expires") != "" &&
		r.URL.Query().Get("OC-Verb") != ""
}

func requestMethodMatches(r *http.Request) bool {
	return strings.EqualFold(r.Method, r.URL.Query().Get("OC-Verb"))
}

func requestMethodIsAllowed(m string, allowedMethods []string) bool {
	for _, allowed := range allowedMethods {
		if strings.EqualFold(m, allowed) {
			return true
		}
	}
	return false
}

func urlIsExpired(r *http.Request, now func() time.Time) bool {
	t, err := time.Parse(time.RFC3339, r.URL.Query().Get("OC-Date"))
	if err != nil {
		return true
	}
	expires, err := time.ParseDuration(r.URL.Query().Get("OC-Expires") + "s")
	if err != nil {
		return true
	}
	t.Add(expires)
	return t.After(now())
}

func signatureIsValid(l log.Logger, r *http.Request, s storepb.StoreService) bool {
	signingKey, err := getSigningKey(r.Context(), s, r.URL.Query().Get("OC-Credential"))
	if err != nil {
		l.Error().Err(err).Msg("could not retrieve signing key")
		return false
	}
	if len(signingKey) == 0 {
		l.Error().Err(err).Msg("signing key empty")
		return false
	}

	q := r.URL.Query()
	signature := q.Get("OC-Signature")
	q.Del("OC-Signature")
	r.URL.RawQuery = q.Encode()
	url := r.URL.String()
	if !r.URL.IsAbs() {
		url = "https://" + r.Host + url // TODO where do we get the scheme from
	}
	return createSignature(url, signingKey) == signature
}

func createSignature(url string, signingKey []byte) string {
	// the oc10 signature check: $hash = \hash_pbkdf2("sha512", $url, $signingKey, 10000, 64, false);
	// - sets the length of the output string to 64
	// - sets raw output to false ->  if raw_output is FALSE length corresponds to twice the byte-length of the derived key (as every byte of the key is returned as two hexits).
	// TODO change to length 128 in oc10?
	// fo golangs pbkdf2.Key we need to use 32 because it will be encoded into 64 hexits later
	hash := pbkdf2.Key([]byte(url), signingKey, iterations, keyLen, sha512.New)
	return hex.EncodeToString(hash)
}

func getSigningKey(ctx context.Context, s storepb.StoreService, credential string) ([]byte, error) {
	res, err := s.Read(ctx, &storepb.ReadRequest{
		Options: &storepb.ReadOptions{
			Database: "proxy",
			Table:    "signing-keys",
		},
		Key: credential,
	})
	if err != nil || len(res.Records) < 1 {
		return []byte{}, err
	}

	return res.Records[0].Value, nil
}
