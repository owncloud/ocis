package middleware

import (
	"context"
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/golang-jwt/jwt/v5"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/ocis-pkg/oidc"
	oidcmocks "github.com/owncloud/ocis/v2/ocis-pkg/oidc/mocks"
	"github.com/shamaton/msgpack/v2"
	"github.com/stretchr/testify/mock"
	"go-micro.dev/v4/store"
	"golang.org/x/crypto/sha3"
)

var _ = Describe("Authenticating requests", Label("OIDCAuthenticator"), func() {
	var authenticator Authenticator

	oc := oidcmocks.OIDCClient{}
	// Return a fresh claims map on every call. getClaims mutates the returned
	// map (claims["exp"] = ...) and then reads it from a fire-and-forget cache
	// goroutine, so handing every call the same map instance would let one
	// spec's goroutine race another spec's mutation. The real OIDC client
	// returns a new map per call, which this mirrors.
	oc.On("VerifyAccessToken", mock.Anything, mock.Anything).Return(
		oidc.RegClaimsWithSID{
			SessionID: "a-session-id",
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Unix(1147483647, 0)),
			},
		},
		func(_ context.Context, _ string) jwt.MapClaims {
			return jwt.MapClaims{
				"sid": "a-session-id",
				"exp": 1147483647,
			}
		},
		nil,
	)
	/*
		// to test with skipUserInfo:  true, we need to also use an interface so we can mock the UserInfo.Claim call
		oc.On("UserInfo", mock.Anything, mock.Anything).Return(
			&oidc.UserInfo{
				Subject:       "my-sub",
				EmailVerified: true,
				Email:         "test@example.org",
			},
			nil,
		)
	*/

	BeforeEach(func() {
		authenticator = &OIDCAuthenticator{
			OIDCIss:       "http://idp.example.com",
			Logger:        log.NewLogger(),
			oidcClient:    &oc,
			userInfoCache: store.NewMemoryStore(),
			skipUserInfo:  true,
		}
	})

	When("the request contains correct data", func() {
		It("should successfully authenticate", func() {
			req := httptest.NewRequest(http.MethodGet, "http://example.com/example/path", http.NoBody)
			req.Header.Set(_headerAuthorization, "Bearer jwt.token.sig")

			req2, valid := authenticator.Authenticate(req)

			Expect(valid).To(Equal(true))
			Expect(req2).ToNot(BeNil())
		})
		It("should successfully authenticate", func() {
			req := httptest.NewRequest(http.MethodGet, "http://example.com/dav/public-files", http.NoBody)
			req.Header.Set(_headerAuthorization, "Bearer jwt.token.sig")

			req2, valid := authenticator.Authenticate(req)

			Expect(valid).To(Equal(true))
			Expect(req2).ToNot(BeNil())
		})
		It("should skip authenticate if the header ShareToken is set", func() {
			req := httptest.NewRequest(http.MethodGet, "http://example.com/dav/public-files/", http.NoBody)
			req.Header.Set(_headerAuthorization, "Bearer jwt.token.sig")
			req.Header.Set(headerShareToken, "sharetoken")

			req2, valid := authenticator.Authenticate(req)

			// TODO Should the authentication of public path requests is handled by another authenticator?
			//Expect(valid).To(Equal(false))
			//Expect(req2).To(BeNil())
			Expect(valid).To(Equal(true))
			Expect(req2).ToNot(BeNil())
		})
		It("should skip authenticate if the 'public-token' is set", func() {
			req := httptest.NewRequest(http.MethodGet, "http://example.com/dav/public-files/?public-token=sharetoken", http.NoBody)
			req.Header.Set(_headerAuthorization, "Bearer jwt.token.sig")

			req2, valid := authenticator.Authenticate(req)

			// TODO Should the authentication of public path requests is handled by another authenticator?
			//Expect(valid).To(Equal(false))
			//Expect(req2).To(BeNil())
			Expect(valid).To(Equal(true))
			Expect(req2).ToNot(BeNil())
		})
	})

	When("the userinfo cache holds an entry whose cached expiry has passed", func() {
		It("re-verifies the access token instead of rejecting the request", func() {
			const token = "jwt.token.sig"

			// Reproduce the cache key the authenticator computes for the token.
			hash := make([]byte, 64)
			sha3.ShakeSum256(hash, []byte(token))
			encodedHash := base64.URLEncoding.EncodeToString(hash)

			// Seed the cache with claims whose "exp" already lies in the past,
			// but keep the store record itself readable (no Expiry) to simulate
			// the window where the entry has not been evicted yet while its
			// cached expiry has passed.
			cachedClaims := map[string]interface{}{
				"sub": "cached-subject",
				"exp": time.Now().Add(-time.Hour).Unix(),
			}
			value, err := msgpack.MarshalAsMap(cachedClaims)
			Expect(err).ToNot(HaveOccurred())

			cache := store.NewMemoryStore()
			Expect(cache.Write(&store.Record{Key: encodedHash, Value: value})).To(Succeed())

			oidcAuth := &OIDCAuthenticator{
				OIDCIss:       "http://idp.example.com",
				Logger:        log.NewLogger(),
				oidcClient:    &oc,
				userInfoCache: cache,
				skipUserInfo:  true,
				TimeFunc:      time.Now,
			}

			req := httptest.NewRequest(http.MethodGet, "http://example.com/example/path", http.NoBody)
			claims, _, err := oidcAuth.getClaims(token, req)

			// Before the fix getClaims returned jwt.ErrTokenExpired here. The
			// token is still valid (the mock re-verifies it), so authentication
			// must succeed via the re-verification fall-through.
			Expect(err).ToNot(HaveOccurred())
			Expect(claims).ToNot(BeNil())
			Expect(claims["sid"]).To(Equal("a-session-id"))
		})

		It("still rejects a token that the IDP no longer accepts", func() {
			const token = "rejected.token.sig"

			hash := make([]byte, 64)
			sha3.ShakeSum256(hash, []byte(token))
			encodedHash := base64.URLEncoding.EncodeToString(hash)

			value, err := msgpack.MarshalAsMap(map[string]interface{}{
				"sub": "cached-subject",
				"exp": time.Now().Add(-time.Hour).Unix(),
			})
			Expect(err).ToNot(HaveOccurred())

			cache := store.NewMemoryStore()
			Expect(cache.Write(&store.Record{Key: encodedHash, Value: value})).To(Succeed())

			rejectingClient := oidcmocks.OIDCClient{}
			rejectingClient.On("VerifyAccessToken", mock.Anything, mock.Anything).Return(
				oidc.RegClaimsWithSID{}, jwt.MapClaims{}, jwt.ErrTokenExpired,
			)

			oidcAuth := &OIDCAuthenticator{
				OIDCIss:       "http://idp.example.com",
				Logger:        log.NewLogger(),
				oidcClient:    &rejectingClient,
				userInfoCache: cache,
				skipUserInfo:  true,
				TimeFunc:      time.Now,
			}

			req := httptest.NewRequest(http.MethodGet, "http://example.com/example/path", http.NoBody)
			_, _, err = oidcAuth.getClaims(token, req)

			// A genuinely invalid/expired token must still be rejected by the
			// re-verification path, so the expiry fall-through does not weaken auth.
			Expect(err).To(HaveOccurred())
		})
	})
})
