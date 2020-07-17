package middleware

import (
	"context"
	"crypto/sha512"
	"encoding/hex"
	"net/http"
	"strings"
	"time"

	storepb "github.com/owncloud/ocis-store/pkg/proto/v0"
	"golang.org/x/crypto/pbkdf2"
)

// PresignedURL provides a middleware to check access secured by a presigned URL.
func PresignedURL(opts ...Option) func(next http.Handler) http.Handler {
	opt := newOptions(opts...)

	/*tokenManager, err := jwt.New(map[string]interface{}{
		"secret":  opt.TokenManagerConfig.JWTSecret,
		"expires": int64(60),
	})
	if err != nil {
		opt.Logger.Fatal().Err(err).Msgf("Could not initialize token-manager")
	}
	*/

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			/* commented, because moving the ocs specific unmarshaling and statuscode mangling here seems wrong
			if isGetSigningKeyRequest(r) {
				claims := oidc.FromContext(r.Context())
				if claims == nil {
					http.Error(w, "No claims in context", http.StatusUnauthorized)
					return
				}

				signingKey, _ := getSigningKey(r.Context(), opt.Store, claims.Email)
				if len(signingKey) == 0 {
					http.Error(w, "No signing key", http.StatusInternalServerError)
					return
				}
				// TODO render as json or xml?
				return
			}
			*/
			if isSignedRequest(r) {
				if signedRequestIsValid(r, opt.Store) {
					// TODO store user in context, let account middleware lookup the id?
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

/*
func isGetSigningKeyRequest(r *http.Request) bool {
	return r.URL.Path == "/ocs/v1.php/cloud/user/signing-key" || r.URL.Path == "/ocs/v2.php/cloud/user/signing-key"
}
*/

func isSignedRequest(r *http.Request) bool {
	return r.URL.Query().Get("OC-Signature") != ""
}

func signedRequestIsValid(r *http.Request, s storepb.StoreService) bool {
	// cheap checks first
	// TODO OC-Algorythm - defined the used algo (e.g. sha256 or sha512 - we should agree on one default algo and make this parameter optional)
	// OC-Credential - defines the user scope (shall we use the owncloud user id here - this might leak internal data ....) REQUIRED
	// OC-Date - defined the date the url was signed (ISO 8601 UTC) REQUIRED
	// OC-Expires - defines the expiry interval in seconds (between 1 and 604800 = 7 days) REQUIRED
	// TODO OC-Verb - defines for which http verb the request is valid - defaults to GET OPTIONAL
	// OC-Signature - the computed signature - server will verify the request upon this REQUIRED
	if r.URL.Query().Get("OC-Signature") == "" || r.URL.Query().Get("OC-Credential") == "" || r.URL.Query().Get("OC-Date") == "" || r.URL.Query().Get("OC-Expires") == "" || r.URL.Query().Get("OC-Verb") == "" {
		return false
	}

	if !strings.EqualFold(r.Method, r.URL.Query().Get("OC-Verb")) {
		return false
	}

	if t, err := time.Parse(time.RFC3339, r.URL.Query().Get("OC-Date")); err != nil {
		return false
	} else if expires, err := time.ParseDuration(r.URL.Query().Get("OC-Expires") + "s"); err != nil {
		return false
	} else {
		t.Add(expires)
		if t.After(time.Now()) { // TODO now client time and server time must be in sync
			return false
		}
	}

	signingKey, _ := getSigningKey(r.Context(), s, r.URL.Query().Get("OC-Credential"))
	if len(signingKey) == 0 {
		return false
	}

	signature := r.URL.Query().Get("OC-Signature")
	r.URL.Query().Del("OC-Signature")
	url := r.URL.String()
	hash := pbkdf2.Key([]byte(url), signingKey, 10000, sha512.Size, sha512.New)
	if hex.EncodeToString(hash) != signature {
		return false
	}
	return true
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
		/* no need to create the key if that is handlead by ocs / a dedicated url-signer service
		key := make([]byte, 64)
		_, err := rand.Read(key[:])
		if err != nil {
			return []byte{}, err
		}
		_, err = s.Write(ctx, &storepb.WriteRequest{
			Options: &storepb.WriteOptions{
				Database: "proxy",
				Table:    "signing-keys",
			},
			Record: &storepb.Record{
				Key:   credential, // TODO username or id?
				Value: key,
				// TODO Expiry?
			},
		})
		if err != nil {
			return []byte{}, err
		}
		return key, nil
		*/
	}

	return res.Records[0].Value, nil
}
