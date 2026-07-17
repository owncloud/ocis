package middleware_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path"
	"strconv"

	appprovider "github.com/cs3org/go-cs3apis/cs3/app/provider/v1beta1"
	userv1beta1 "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	providerv1beta1 "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	typesv1beta1 "github.com/cs3org/go-cs3apis/cs3/types/v1beta1"
	"github.com/golang-jwt/jwt/v5"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/owncloud/ocis/v2/services/collaboration/pkg/config"
	"github.com/owncloud/ocis/v2/services/collaboration/pkg/helpers"
	"github.com/owncloud/ocis/v2/services/collaboration/pkg/middleware"
	"github.com/owncloud/ocis/v2/services/collaboration/pkg/wopisrc"
	"github.com/owncloud/ocis/v2/services/settings/pkg/store/defaults"
	ctxpkg "github.com/owncloud/reva/v2/pkg/ctx"
	"github.com/owncloud/reva/v2/pkg/token"
	rjwt "github.com/owncloud/reva/v2/pkg/token/manager/jwt"
)

var _ = Describe("Wopi Context Middleware", func() {
	var (
		cfg     *config.Config
		ctx     context.Context
		mw      http.Handler
		rid     *providerv1beta1.ResourceId
		tknMngr token.Manager
		user    *userv1beta1.User
		src     *url.URL
	)

	BeforeEach(func() {
		var err error
		cfg = &config.Config{
			TokenManager: &config.TokenManager{JWTSecret: "jwtSecret"},
			Wopi: config.Wopi{
				Secret:  "wopiSecret",
				WopiSrc: "https://localhost:9300",
			},
		}

		ctx = context.Background()

		next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
		mw = middleware.WopiContextAuthMiddleware(cfg, nil, next)

		tknMngr, err = rjwt.New(map[string]interface{}{
			"secret":  cfg.TokenManager.JWTSecret,
			"expires": int64(24 * 60 * 60),
		})
		Expect(err).ToNot(HaveOccurred())

		user = &userv1beta1.User{
			Id: &userv1beta1.UserId{
				Idp:      "example.com",
				OpaqueId: "12345",
				Type:     userv1beta1.UserType_USER_TYPE_PRIMARY,
			},
			Username: "admin",
			Mail:     "admin@example.com",
		}

		rid = &providerv1beta1.ResourceId{
			StorageId: "storageID",
			OpaqueId:  "opaqueID",
			SpaceId:   "spaceID",
		}

		src, err = url.Parse(cfg.Wopi.WopiSrc)
		src.Path = path.Join("wopi", "files", helpers.HashResourceId(rid))
		Expect(err).ToNot(HaveOccurred())
	})
	It("Should not authorize with empty access token", func() {
		req := httptest.NewRequest("GET", src.String(), nil).WithContext(ctx)
		resp := httptest.NewRecorder()

		mw.ServeHTTP(resp, req)
		Expect(resp.Code).To(Equal(http.StatusUnauthorized))
	})
	It("Should not authorize with malformed access token", func() {
		req := httptest.NewRequest("GET", src.String(), nil).WithContext(ctx)
		q := req.URL.Query()
		q.Add("access_token", "token")
		req.URL.RawQuery = q.Encode()

		resp := httptest.NewRecorder()

		mw.ServeHTTP(resp, req)
		Expect(resp.Code).To(Equal(http.StatusUnauthorized))
	})
	It("Should not authorize when fileID mismatches", func() {
		req := httptest.NewRequest("GET", src.String(), nil).WithContext(ctx)
		// create request with different fileID in the wopi context
		token, err := tknMngr.MintToken(ctx, user, nil)
		Expect(err).ToNot(HaveOccurred())
		wopiContext := middleware.WopiContext{
			AccessToken: token,
			ViewMode:    appprovider.ViewMode_VIEW_MODE_READ_WRITE,
			FileReference: &providerv1beta1.Reference{
				ResourceId: &providerv1beta1.ResourceId{
					StorageId: "storageID",
					OpaqueId:  "opaqueID2",
					SpaceId:   "spaceID",
				},
				Path: ".",
			},
		}
		wopiToken, ttl, err := middleware.GenerateWopiToken(wopiContext, cfg, nil)
		q := req.URL.Query()
		q.Add("access_token", wopiToken)
		q.Add("access_token_ttl", strconv.FormatInt(ttl, 10))
		req.URL.RawQuery = q.Encode()
		resp := httptest.NewRecorder()

		mw.ServeHTTP(resp, req)
		Expect(resp.Code).To(Equal(http.StatusUnauthorized))
	})
	It("Should not authorize with wrong wopi secret", func() {
		src.Path = path.Join("wopi", "files", helpers.HashResourceId(rid))
		req := httptest.NewRequest("GET", src.String(), nil).WithContext(ctx)
		token, err := tknMngr.MintToken(ctx, user, nil)
		Expect(err).ToNot(HaveOccurred())

		wopiContext := middleware.WopiContext{
			AccessToken: token,
		}
		// use wrong wopi secret when generating the wopi token
		wopiToken, ttl, err := middleware.GenerateWopiToken(wopiContext, &config.Config{
			TokenManager: &config.TokenManager{JWTSecret: cfg.TokenManager.JWTSecret},
			Wopi: config.Wopi{
				Secret: "wrongSecret",
			},
		}, nil)
		q := req.URL.Query()
		q.Add("access_token", wopiToken)
		q.Add("access_token_ttl", strconv.FormatInt(ttl, 10))
		req.URL.RawQuery = q.Encode()
		resp := httptest.NewRecorder()

		mw.ServeHTTP(resp, req)
		Expect(resp.Code).To(Equal(http.StatusUnauthorized))
	})
	It("Should authorize successful", func() {
		req := httptest.NewRequest("GET", src.String(), nil).WithContext(ctx)
		token, err := tknMngr.MintToken(ctx, user, nil)
		Expect(err).ToNot(HaveOccurred())

		wopiContext := middleware.WopiContext{
			AccessToken: token,
			ViewMode:    appprovider.ViewMode_VIEW_MODE_READ_WRITE,
			FileReference: &providerv1beta1.Reference{
				ResourceId: rid,
				Path:       ".",
			},
		}
		wopiToken, ttl, err := middleware.GenerateWopiToken(wopiContext, cfg, nil)
		q := req.URL.Query()
		q.Add("access_token", wopiToken)
		q.Add("access_token_ttl", strconv.FormatInt(ttl, 10))
		req.URL.RawQuery = q.Encode()
		resp := httptest.NewRecorder()

		mw.ServeHTTP(resp, req)
		Expect(resp.Code).To(Equal(http.StatusOK))
	})
	It("Should encrypt the ViewOnlyToken inside the WOPI token and restore it on the way in", func() {
		// a known, recognizable plaintext so we can assert it is not readable in the JWT
		const viewOnlyPlaintext = "my-secret-view-only-token"

		// capture the WopiContext the middleware hands to the next handler
		var captured middleware.WopiContext
		var capturedErr error
		capturingNext := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			captured, capturedErr = middleware.WopiContextFromCtx(r.Context())
			w.WriteHeader(http.StatusOK)
		})
		capturingMw := middleware.WopiContextAuthMiddleware(cfg, nil, capturingNext)

		token, err := tknMngr.MintToken(ctx, user, nil)
		Expect(err).ToNot(HaveOccurred())

		wopiContext := middleware.WopiContext{
			AccessToken:   token,
			ViewOnlyToken: viewOnlyPlaintext,
			ViewMode:      appprovider.ViewMode_VIEW_MODE_VIEW_ONLY,
			FileReference: &providerv1beta1.Reference{
				ResourceId: rid,
				Path:       ".",
			},
		}
		wopiToken, ttl, err := middleware.GenerateWopiToken(wopiContext, cfg, nil)
		Expect(err).ToNot(HaveOccurred())

		// Security property: the plaintext view-only token must not be recoverable
		// from the (readable) JWT. Parse the JWT and assert the ViewOnlyToken claim
		// is present but is ciphertext, not the plaintext value.
		claims := &middleware.Claims{}
		_, err = jwt.ParseWithClaims(wopiToken, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(cfg.Wopi.Secret), nil
		})
		Expect(err).ToNot(HaveOccurred())
		Expect(claims.WopiContext.ViewOnlyToken).ToNot(BeEmpty())
		Expect(claims.WopiContext.ViewOnlyToken).ToNot(Equal(viewOnlyPlaintext))

		// No-regression property: the middleware must decrypt the token back to the
		// original plaintext so downstream consumers get a usable reva token.
		req := httptest.NewRequest("GET", src.String(), nil).WithContext(ctx)
		q := req.URL.Query()
		q.Add("access_token", wopiToken)
		q.Add("access_token_ttl", strconv.FormatInt(ttl, 10))
		req.URL.RawQuery = q.Encode()
		resp := httptest.NewRecorder()

		capturingMw.ServeHTTP(resp, req)

		Expect(resp.Code).To(Equal(http.StatusOK))
		Expect(capturedErr).ToNot(HaveOccurred())
		Expect(captured.ViewOnlyToken).To(Equal(viewOnlyPlaintext))
	})
	It("Should authorize successful with template reference", func() {
		req := httptest.NewRequest("GET", src.String(), nil).WithContext(ctx)
		token, err := tknMngr.MintToken(ctx, user, nil)
		Expect(err).ToNot(HaveOccurred())

		wopiContext := middleware.WopiContext{
			AccessToken: token,
			ViewMode:    appprovider.ViewMode_VIEW_MODE_READ_WRITE,
			TemplateReference: &providerv1beta1.Reference{
				ResourceId: rid,
				Path:       ".",
			},
			FileReference: &providerv1beta1.Reference{
				ResourceId: &providerv1beta1.ResourceId{
					StorageId: "storageID",
					OpaqueId:  "opaqueID2",
					SpaceId:   "spaceID",
				},
			},
		}
		wopiToken, ttl, err := middleware.GenerateWopiToken(wopiContext, cfg, nil)
		q := req.URL.Query()
		q.Add("access_token", wopiToken)
		q.Add("access_token_ttl", strconv.FormatInt(ttl, 10))
		req.URL.RawQuery = q.Encode()
		resp := httptest.NewRecorder()

		mw.ServeHTTP(resp, req)
		Expect(resp.Code).To(Equal(http.StatusOK))
	})
	It("Should not authorize when no reference matches", func() {
		req := httptest.NewRequest("GET", src.String(), nil).WithContext(ctx)
		token, err := tknMngr.MintToken(ctx, user, nil)
		Expect(err).ToNot(HaveOccurred())

		wopiContext := middleware.WopiContext{
			AccessToken: token,
			ViewMode:    appprovider.ViewMode_VIEW_MODE_READ_WRITE,
			TemplateReference: &providerv1beta1.Reference{
				ResourceId: &providerv1beta1.ResourceId{
					StorageId: "storageID",
					OpaqueId:  "opaqueID3",
					SpaceId:   "spaceID",
				},
				Path: ".",
			},
			FileReference: &providerv1beta1.Reference{
				ResourceId: &providerv1beta1.ResourceId{
					StorageId: "storageID",
					OpaqueId:  "opaqueID2",
					SpaceId:   "spaceID",
				},
			},
		}
		wopiToken, ttl, err := middleware.GenerateWopiToken(wopiContext, cfg, nil)
		q := req.URL.Query()
		q.Add("access_token", wopiToken)
		q.Add("access_token_ttl", strconv.FormatInt(ttl, 10))
		req.URL.RawQuery = q.Encode()
		resp := httptest.NewRecorder()

		mw.ServeHTTP(resp, req)
		Expect(resp.Code).To(Equal(http.StatusUnauthorized))
	})
	It("Should not authorize with proxy when fileID mismatches", func() {
		cfg.Wopi.ProxySecret = "proxySecret"
		cfg.Wopi.ProxyURL = "https://proxy"
		src, err := wopisrc.GenerateWopiSrc(helpers.HashResourceId(rid), cfg)
		Expect(err).ToNot(HaveOccurred())

		req := httptest.NewRequest("GET", src.String(), nil).WithContext(ctx)
		token, err := tknMngr.MintToken(ctx, user, nil)
		Expect(err).ToNot(HaveOccurred())
		wopiContext := middleware.WopiContext{
			AccessToken: token,
			ViewMode:    appprovider.ViewMode_VIEW_MODE_READ_WRITE,
			FileReference: &providerv1beta1.Reference{
				ResourceId: &providerv1beta1.ResourceId{
					StorageId: "storageID",
					OpaqueId:  "opaqueID3",
					SpaceId:   "spaceID",
				},
				Path: ".",
			},
		}
		wopiToken, ttl, err := middleware.GenerateWopiToken(wopiContext, cfg, nil)
		q := req.URL.Query()
		q.Add("access_token", wopiToken)
		q.Add("access_token_ttl", strconv.FormatInt(ttl, 10))
		req.URL.RawQuery = q.Encode()

		resp := httptest.NewRecorder()
		mw.ServeHTTP(resp, req)
		Expect(resp.Code).To(Equal(http.StatusUnauthorized))
	})
	It("Should authorize successful with proxy", func() {
		cfg.Wopi.ProxySecret = "proxySecret"
		cfg.Wopi.ProxyURL = "https://proxy"
		src, err := wopisrc.GenerateWopiSrc(helpers.HashResourceId(rid), cfg)
		Expect(err).ToNot(HaveOccurred())

		req := httptest.NewRequest("GET", src.String(), nil).WithContext(ctx)
		token, err := tknMngr.MintToken(ctx, user, nil)
		Expect(err).ToNot(HaveOccurred())
		wopiContext := middleware.WopiContext{
			AccessToken: token,
			ViewMode:    appprovider.ViewMode_VIEW_MODE_READ_WRITE,
			FileReference: &providerv1beta1.Reference{
				ResourceId: rid,
				Path:       ".",
			},
		}
		wopiToken, ttl, err := middleware.GenerateWopiToken(wopiContext, cfg, nil)
		q := req.URL.Query()
		q.Add("access_token", wopiToken)
		q.Add("access_token_ttl", strconv.FormatInt(ttl, 10))
		req.URL.RawQuery = q.Encode()

		resp := httptest.NewRecorder()
		mw.ServeHTTP(resp, req)
		Expect(resp.Code).To(Equal(http.StatusOK))
	})
	It("Should preserve the admin role Opaque entry through the reva token mint/dismantle round trip", func() {
		// Unlike the fileinfo/fileconnector tests, which inject the user directly
		// into the context via ctxpkg.ContextSetUser (bypassing token serialization
		// entirely), this exercises the actual seam IsAdminUser depends on in
		// production: the reva access token is minted once at file-open time,
		// then later dismantled here on every WOPI request. If the roles Opaque
		// entry didn't survive that round trip, IsAdminUser would silently never
		// fire despite every direct-injection unit test passing.
		rolesJSON, err := json.Marshal([]string{defaults.BundleUUIDRoleAdmin})
		Expect(err).ToNot(HaveOccurred())

		adminUser := &userv1beta1.User{
			Id: &userv1beta1.UserId{
				Idp:      "example.com",
				OpaqueId: "12345",
				Type:     userv1beta1.UserType_USER_TYPE_PRIMARY,
			},
			Username: "admin",
			Mail:     "admin@example.com",
			Opaque: &typesv1beta1.Opaque{
				Map: map[string]*typesv1beta1.OpaqueEntry{
					"roles": {
						Decoder: "json",
						Value:   rolesJSON,
					},
				},
			},
		}

		var captured *userv1beta1.User
		capturingNext := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			captured, _ = ctxpkg.ContextGetUser(r.Context())
			w.WriteHeader(http.StatusOK)
		})
		capturingMw := middleware.WopiContextAuthMiddleware(cfg, nil, capturingNext)

		token, err := tknMngr.MintToken(ctx, adminUser, nil)
		Expect(err).ToNot(HaveOccurred())

		wopiContext := middleware.WopiContext{
			AccessToken: token,
			ViewMode:    appprovider.ViewMode_VIEW_MODE_READ_WRITE,
			FileReference: &providerv1beta1.Reference{
				ResourceId: rid,
				Path:       ".",
			},
		}
		wopiToken, ttl, err := middleware.GenerateWopiToken(wopiContext, cfg, nil)
		Expect(err).ToNot(HaveOccurred())

		req := httptest.NewRequest("GET", src.String(), nil).WithContext(ctx)
		q := req.URL.Query()
		q.Add("access_token", wopiToken)
		q.Add("access_token_ttl", strconv.FormatInt(ttl, 10))
		req.URL.RawQuery = q.Encode()
		resp := httptest.NewRecorder()

		capturingMw.ServeHTTP(resp, req)

		Expect(resp.Code).To(Equal(http.StatusOK))
		Expect(captured).ToNot(BeNil())

		rolesEntry, ok := captured.GetOpaque().GetMap()["roles"]
		Expect(ok).To(BeTrue())

		var roleIDs []string
		Expect(json.Unmarshal(rolesEntry.GetValue(), &roleIDs)).To(Succeed())
		Expect(roleIDs).To(ContainElement(defaults.BundleUUIDRoleAdmin))
	})
})
