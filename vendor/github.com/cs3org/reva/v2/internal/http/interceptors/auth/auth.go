// Copyright 2018-2021 CERN
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// In applying this license, CERN does not waive the privileges and immunities
// granted to it by virtue of its status as an Intergovernmental Organization
// or submit itself to any jurisdiction.

package auth

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/bluele/gcache"
	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	userpb "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	"github.com/cs3org/reva/v2/internal/http/interceptors/auth/credential/registry"
	tokenregistry "github.com/cs3org/reva/v2/internal/http/interceptors/auth/token/registry"
	tokenwriterregistry "github.com/cs3org/reva/v2/internal/http/interceptors/auth/tokenwriter/registry"
	"github.com/cs3org/reva/v2/pkg/appctx"
	"github.com/cs3org/reva/v2/pkg/auth"
	"github.com/cs3org/reva/v2/pkg/auth/scope"
	ctxpkg "github.com/cs3org/reva/v2/pkg/ctx"
	"github.com/cs3org/reva/v2/pkg/errtypes"
	"github.com/cs3org/reva/v2/pkg/rgrpc/status"
	"github.com/cs3org/reva/v2/pkg/rgrpc/todo/pool"
	"github.com/cs3org/reva/v2/pkg/rhttp/global"
	"github.com/cs3org/reva/v2/pkg/sharedconf"
	"github.com/cs3org/reva/v2/pkg/token"
	tokenmgr "github.com/cs3org/reva/v2/pkg/token/manager/registry"
	"github.com/cs3org/reva/v2/pkg/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/metadata"
)

// name is the Tracer name used to identify this instrumentation library.
const tracerName = "auth"

var (
	cacheOnce       sync.Once
	userGroupsCache gcache.Cache
)

type config struct {
	Priority   int    `mapstructure:"priority"`
	GatewaySvc string `mapstructure:"gatewaysvc"`
	// TODO(jdf): Realm is optional, will be filled with request host if not given?
	Realm                  string                            `mapstructure:"realm"`
	CredentialsByUserAgent map[string]string                 `mapstructure:"credentials_by_user_agent"`
	CredentialChain        []string                          `mapstructure:"credential_chain"`
	CredentialStrategies   map[string]map[string]interface{} `mapstructure:"credential_strategies"`
	TokenStrategy          string                            `mapstructure:"token_strategy"`
	TokenStrategies        map[string]map[string]interface{} `mapstructure:"token_strategies"`
	TokenManager           string                            `mapstructure:"token_manager"`
	TokenManagers          map[string]map[string]interface{} `mapstructure:"token_managers"`
	TokenWriter            string                            `mapstructure:"token_writer"`
	TokenWriters           map[string]map[string]interface{} `mapstructure:"token_writers"`
	UserGroupsCacheSize    int                               `mapstructure:"usergroups_cache_size"`
}

func parseConfig(m map[string]interface{}) (*config, error) {
	c := &config{}
	if err := mapstructure.Decode(m, c); err != nil {
		err = errors.Wrap(err, "error decoding conf")
		return nil, err
	}
	return c, nil
}

// New returns a new middleware with defined priority.
func New(m map[string]interface{}, unprotected []string, tp trace.TracerProvider) (global.Middleware, error) {
	conf, err := parseConfig(m)
	if err != nil {
		return nil, err
	}

	conf.GatewaySvc = sharedconf.GetGatewaySVC(conf.GatewaySvc)

	// set defaults
	if conf.TokenStrategy == "" {
		conf.TokenStrategy = "header"
	}

	if conf.TokenWriter == "" {
		conf.TokenWriter = "header"
	}

	if conf.TokenManager == "" {
		conf.TokenManager = "jwt"
	}

	if len(conf.CredentialChain) == 0 {
		conf.CredentialChain = []string{"basic", "bearer"}
	}

	if conf.CredentialsByUserAgent == nil {
		conf.CredentialsByUserAgent = map[string]string{}
	}

	if conf.UserGroupsCacheSize == 0 {
		conf.UserGroupsCacheSize = 5000
	}

	cacheOnce.Do(func() {
		userGroupsCache = gcache.New(conf.UserGroupsCacheSize).LFU().Build()
	})

	credChain := map[string]auth.CredentialStrategy{}
	for i, key := range conf.CredentialChain {
		f, ok := registry.NewCredentialFuncs[conf.CredentialChain[i]]
		if !ok {
			return nil, fmt.Errorf("credential strategy not found: %s", conf.CredentialChain[i])
		}

		credStrategy, err := f(conf.CredentialStrategies[conf.CredentialChain[i]])
		if err != nil {
			return nil, err
		}
		credChain[key] = credStrategy
	}

	g, ok := tokenregistry.NewTokenFuncs[conf.TokenStrategy]
	if !ok {
		return nil, fmt.Errorf("token strategy not found: %s", conf.TokenStrategy)
	}

	tokenStrategy, err := g(conf.TokenStrategies[conf.TokenStrategy])
	if err != nil {
		return nil, err
	}

	h, ok := tokenmgr.NewFuncs[conf.TokenManager]
	if !ok {
		return nil, fmt.Errorf("token manager not found: %s", conf.TokenStrategy)
	}

	tokenManager, err := h(conf.TokenManagers[conf.TokenManager])
	if err != nil {
		return nil, err
	}

	i, ok := tokenwriterregistry.NewTokenFuncs[conf.TokenWriter]
	if !ok {
		return nil, fmt.Errorf("token writer not found: %s", conf.TokenWriter)
	}

	tokenWriter, err := i(conf.TokenWriters[conf.TokenWriter])
	if err != nil {
		return nil, err
	}

	chain := func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// OPTION requests need to pass for preflight requests
			// TODO(labkode): this will break options for auth protected routes.
			// Maybe running the CORS middleware before auth kicks in is enough.
			ctx := r.Context()
			span := trace.SpanFromContext(ctx)
			defer span.End()
			if !span.SpanContext().HasTraceID() {
				_, span = tp.Tracer(tracerName).Start(ctx, "http auth interceptor")
			}

			if r.Method == "OPTIONS" {
				h.ServeHTTP(w, r)
				return
			}

			log := appctx.GetLogger(r.Context())
			isUnprotectedEndpoint := false

			// For unprotected URLs, we try to authenticate the request in case some service needs it,
			// but don't return any errors if it fails.
			if utils.Skip(r.URL.Path, unprotected) {
				log.Info().Msg("skipping auth check for: " + r.URL.Path)
				isUnprotectedEndpoint = true
			}

			ctx, err := authenticateUser(w, r, conf, tokenStrategy, tokenManager, tokenWriter, credChain, isUnprotectedEndpoint)
			if err != nil {
				if !isUnprotectedEndpoint {
					return
				}
			} else {
				u, ok := ctxpkg.ContextGetUser(ctx)
				if ok {
					span.SetAttributes(semconv.EnduserIDKey.String(u.Id.OpaqueId))
				}

				r = r.WithContext(ctx)
			}
			h.ServeHTTP(w, r)

		})
	}
	return chain, nil
}

func authenticateUser(w http.ResponseWriter, r *http.Request, conf *config, tokenStrategy auth.TokenStrategy, tokenManager token.Manager, tokenWriter auth.TokenWriter, credChain map[string]auth.CredentialStrategy, isUnprotectedEndpoint bool) (context.Context, error) {
	ctx := r.Context()
	log := appctx.GetLogger(ctx)

	// Add the request user-agent to the ctx
	ctx = metadata.NewIncomingContext(ctx, metadata.New(map[string]string{ctxpkg.UserAgentHeader: r.UserAgent()}))

	client, err := pool.GetGatewayServiceClient(conf.GatewaySvc)
	if err != nil {
		logError(isUnprotectedEndpoint, log, err, "error getting the authsvc client", http.StatusUnauthorized, w)
		return nil, err
	}

	tkn := tokenStrategy.GetToken(r)
	if tkn == "" {
		log.Warn().Msg("core access token not set")

		userAgentCredKeys := getCredsForUserAgent(r.UserAgent(), conf.CredentialsByUserAgent, conf.CredentialChain)

		// obtain credentials (basic auth, bearer token, ...) based on user agent
		var creds *auth.Credentials
		for _, k := range userAgentCredKeys {
			creds, err = credChain[k].GetCredentials(w, r)
			if err != nil {
				log.Debug().Err(err).Msg("error retrieving credentials")
			}

			if creds != nil {
				log.Debug().Msgf("credentials obtained from credential strategy: type: %s, client_id: %s", creds.Type, creds.ClientID)
				break
			}
		}

		// if no credentials are found, reply with authentication challenge depending on user agent
		if creds == nil {
			if !isUnprotectedEndpoint {
				for _, key := range userAgentCredKeys {
					if cred, ok := credChain[key]; ok {
						cred.AddWWWAuthenticate(w, r, conf.Realm)
					} else {
						panic("auth credential strategy: " + key + "must have been loaded in init method")
					}
				}
				w.WriteHeader(http.StatusUnauthorized)
			}
			return nil, errtypes.PermissionDenied("no credentials found")
		}

		req := &gateway.AuthenticateRequest{
			Type:         creds.Type,
			ClientId:     creds.ClientID,
			ClientSecret: creds.ClientSecret,
		}

		log.Debug().Msgf("AuthenticateRequest: type: %s, client_id: %s against %s", req.Type, req.ClientId, conf.GatewaySvc)

		res, err := client.Authenticate(ctx, req)
		if err != nil {
			logError(isUnprotectedEndpoint, log, err, "error calling Authenticate", http.StatusUnauthorized, w)
			return nil, err
		}

		if res.Status.Code != rpc.Code_CODE_OK {
			err := status.NewErrorFromCode(res.Status.Code, "auth")
			logError(isUnprotectedEndpoint, log, err, "error generating access token from credentials", http.StatusUnauthorized, w)
			return nil, err
		}

		log.Info().Msg("core access token generated")
		// write token to response
		tkn = res.Token
		tokenWriter.WriteToken(tkn, w)
	} else {
		log.Debug().Msg("access token is already provided")
	}

	// validate token
	u, tokenScope, err := tokenManager.DismantleToken(r.Context(), tkn)
	if err != nil {
		logError(isUnprotectedEndpoint, log, err, "error dismantling token", http.StatusUnauthorized, w)
		return nil, err
	}

	if sharedconf.SkipUserGroupsInToken() {
		var groups []string
		if groupsIf, err := userGroupsCache.Get(u.Id.OpaqueId); err == nil {
			groups = groupsIf.([]string)
		} else {
			groupsRes, err := client.GetUserGroups(ctx, &userpb.GetUserGroupsRequest{UserId: u.Id})
			if err != nil {
				logError(isUnprotectedEndpoint, log, err, "error retrieving user groups", http.StatusInternalServerError, w)
				return nil, err
			}
			groups = groupsRes.Groups
			_ = userGroupsCache.SetWithExpire(u.Id.OpaqueId, groupsRes.Groups, 3600*time.Second)
		}
		u.Groups = groups
	}

	// ensure access to the resource is allowed
	ok, err := scope.VerifyScope(ctx, tokenScope, r.URL.Path)
	if err != nil {
		logError(isUnprotectedEndpoint, log, err, "error verifying scope of access token", http.StatusInternalServerError, w)
		return nil, err
	}
	if !ok {
		err := errtypes.PermissionDenied("access to resource not allowed")
		logError(isUnprotectedEndpoint, log, err, "access to resource not allowed", http.StatusUnauthorized, w)
		return nil, err
	}

	// store user and core access token in context.
	ctx = ctxpkg.ContextSetUser(ctx, u)
	ctx = ctxpkg.ContextSetToken(ctx, tkn)
	ctx = metadata.AppendToOutgoingContext(ctx, ctxpkg.TokenHeader, tkn) // TODO(jfd): hardcoded metadata key. use  PerRPCCredentials?

	ctx = metadata.AppendToOutgoingContext(ctx, ctxpkg.UserAgentHeader, r.UserAgent())

	// store scopes in context
	ctx = ctxpkg.ContextSetScopes(ctx, tokenScope)

	return ctx, nil
}

func logError(isUnprotectedEndpoint bool, log *zerolog.Logger, err error, msg string, status int, w http.ResponseWriter) {
	if !isUnprotectedEndpoint {
		log.Error().Err(err).Msg(msg)
		w.WriteHeader(status)
	}
}

// getCredsForUserAgent returns the WWW Authenticate challenges keys to use given an http request
// and available credentials.
func getCredsForUserAgent(ua string, uam map[string]string, creds []string) []string {
	if ua == "" || len(uam) == 0 {
		return creds
	}

	for u, cred := range uam {
		if strings.Contains(ua, u) {
			for _, v := range creds {
				if v == cred {
					return []string{cred}
				}
			}
			return creds

		}
	}

	return creds
}
