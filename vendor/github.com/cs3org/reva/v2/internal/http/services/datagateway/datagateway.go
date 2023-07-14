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

package datagateway

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"time"

	"github.com/cs3org/reva/v2/pkg/appctx"
	"github.com/cs3org/reva/v2/pkg/errtypes"
	"github.com/cs3org/reva/v2/pkg/rhttp"
	"github.com/cs3org/reva/v2/pkg/rhttp/global"
	"github.com/cs3org/reva/v2/pkg/sharedconf"
	"github.com/golang-jwt/jwt"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
	"go.opentelemetry.io/otel/trace"
)

var tracer trace.Tracer

func init() {
	tracer = otel.Tracer("github.com/cs3org/reva/pkg/storage/utils/decomposedfs/tree")
}

const (
	// TokenTransportHeader holds the header key for the reva transfer token
	TokenTransportHeader = "X-Reva-Transfer"
	// UploadExpiresHeader holds the timestamp for the transport token expiry, defined in https://tus.io/protocols/resumable-upload.html#expiration
	UploadExpiresHeader = "Upload-Expires"
)

func init() {
	global.Register("datagateway", New)
}

// TransferClaims are custom claims for a JWT token to be used between the metadata and data gateways.
type TransferClaims struct {
	jwt.StandardClaims
	Target string `json:"target"`
}
type config struct {
	Prefix               string `mapstructure:"prefix"`
	TransferSharedSecret string `mapstructure:"transfer_shared_secret"`
	Timeout              int64  `mapstructure:"timeout"`
	Insecure             bool   `mapstructure:"insecure"`
}

func (c *config) init() {
	if c.Prefix == "" {
		c.Prefix = "datagateway"
	}

	c.TransferSharedSecret = sharedconf.GetJWTSecret(c.TransferSharedSecret)
}

type svc struct {
	conf    *config
	handler http.Handler
	client  *http.Client
}

// New returns a new datagateway
func New(m map[string]interface{}, log *zerolog.Logger) (global.Service, error) {
	conf := &config{}
	if err := mapstructure.Decode(m, conf); err != nil {
		return nil, err
	}

	conf.init()

	s := &svc{
		conf: conf,
		client: rhttp.GetHTTPClient(
			rhttp.Timeout(time.Duration(conf.Timeout*int64(time.Second))),
			rhttp.Insecure(conf.Insecure),
		),
	}
	s.setHandler()
	return s, nil
}

// Close performs cleanup.
func (s *svc) Close() error {
	return nil
}

func (s *svc) Prefix() string {
	return s.conf.Prefix
}

func (s *svc) Handler() http.Handler {
	return s.handler
}

func (s *svc) Unprotected() []string {
	return []string{
		"/",
	}
}

func (s *svc) setHandler() {
	s.handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx, span := tracer.Start(ctx, "HandlerFunc")
		defer span.End()
		span.SetAttributes(
			semconv.HTTPMethodKey.String(r.Method),
			semconv.HTTPURLKey.String(r.URL.String()),
		)
		r = r.WithContext(ctx)
		switch r.Method {
		case "HEAD":
			addCorsHeader(w)
			s.doHead(w, r)
			return
		case "GET":
			s.doGet(w, r)
			return
		case "PUT":
			s.doPut(w, r)
			return
		case "PATCH":
			s.doPatch(w, r)
			return
		default:
			w.WriteHeader(http.StatusNotImplemented)
			return
		}
	})
}

func addCorsHeader(res http.ResponseWriter) {
	headers := res.Header()
	headers.Set("Access-Control-Allow-Origin", "*")
	headers.Set("Access-Control-Allow-Headers", "Content-Type, Origin, Authorization")
	headers.Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, HEAD")
}

// Verify a transfer token against the given secret
func Verify(ctx context.Context, token string, secret string) (*TransferClaims, error) {
	j, err := jwt.ParseWithClaims(token, &TransferClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return nil, errors.Wrap(err, "error parsing token")
	}

	if claims, ok := j.Claims.(*TransferClaims); ok && j.Valid {
		return claims, nil
	}
	err = errtypes.InvalidCredentials("token invalid")
	return nil, err
}

func (s *svc) verify(ctx context.Context, r *http.Request) (*TransferClaims, error) {
	// Extract transfer token from request header. If not existing, assume that it's the last path segment instead.
	token := r.Header.Get(TokenTransportHeader)
	if token == "" {
		token = path.Base(r.URL.Path)
		r.Header.Set(TokenTransportHeader, token)
	}

	return Verify(ctx, token, s.conf.TransferSharedSecret)
}

func (s *svc) doHead(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := appctx.GetLogger(ctx)

	claims, err := s.verify(ctx, r)
	if err != nil {
		err = errors.Wrap(err, "datagateway: error validating transfer token")
		log.Error().Err(err).Str("token", r.Header.Get(TokenTransportHeader)).Msg("invalid transfer token")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	log.Debug().Str("target", claims.Target).Msg("sending request to internal data server")

	httpClient := s.client
	httpReq, err := rhttp.NewRequest(ctx, "HEAD", claims.Target, nil)
	if err != nil {
		log.Error().Err(err).Msg("wrong request")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	httpReq.Header = r.Header

	httpRes, err := httpClient.Do(httpReq)
	if err != nil {
		log.Error().Err(err).Msg("error doing HEAD request to data service")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer httpRes.Body.Close()

	copyHeader(w.Header(), httpRes.Header)

	// add upload expiry / transfer token expiry header for tus https://tus.io/protocols/resumable-upload.html#expiration
	w.Header().Set(UploadExpiresHeader, time.Unix(claims.ExpiresAt, 0).Format(time.RFC1123))

	if httpRes.StatusCode != http.StatusOK {
		// swallow the body and set content-length to 0 to prevent reverse proxies from trying to read from it
		w.Header().Set("Content-Length", "0")
		w.WriteHeader(httpRes.StatusCode)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *svc) doGet(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := appctx.GetLogger(ctx)

	claims, err := s.verify(ctx, r)
	if err != nil {
		err = errors.Wrap(err, "datagateway: error validating transfer token")
		log.Error().Err(err).Str("token", r.Header.Get(TokenTransportHeader)).Msg("invalid transfer token")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	log.Debug().Str("target", claims.Target).Msg("sending request to internal data server")

	httpClient := s.client
	httpReq, err := rhttp.NewRequest(ctx, "GET", claims.Target, nil)
	if err != nil {
		log.Error().Err(err).Msg("wrong request")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	httpReq.Header = r.Header

	httpRes, err := httpClient.Do(httpReq)
	if err != nil {
		log.Error().Err(err).Msg("error doing GET request to data service")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer httpRes.Body.Close()

	copyHeader(w.Header(), httpRes.Header)
	switch httpRes.StatusCode {
	case http.StatusOK:
	case http.StatusPartialContent:
	default:
		// swallow the body and set content-length to 0 to prevent reverse proxies from trying to read from it
		w.Header().Set("Content-Length", "0")
		w.WriteHeader(httpRes.StatusCode)
		return
	}
	w.WriteHeader(httpRes.StatusCode)

	var c int64
	c, err = io.Copy(w, httpRes.Body)
	if err != nil {
		log.Error().Err(err).Msg("error writing body after headers were sent")
	}
	if httpRes.Header.Get("Content-Length") != "" {
		i, err := strconv.ParseInt(httpRes.Header.Get("Content-Length"), 10, 64)
		if err != nil {
			log.Error().Err(err).Str("content-length", httpRes.Header.Get("Content-Length")).Msg("invalid content length in dataprovider response")
		}
		if i != c {
			log.Error().Int64("content-length", i).Int64("transferred-bytes", c).Msg("content length vs transferred bytes mismatch")
		}
	}
}

func (s *svc) doPut(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := appctx.GetLogger(ctx)

	claims, err := s.verify(ctx, r)
	if err != nil {
		err = errors.Wrap(err, "datagateway: error validating transfer token")
		log.Err(err).Str("token", r.Header.Get(TokenTransportHeader)).Msg("invalid transfer token")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	target := claims.Target
	// add query params to target, clients can send checksums and other information.
	targetURL, err := url.Parse(target)
	if err != nil {
		log.Err(err).Msg("datagateway: error parsing target url")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	targetURL.RawQuery = r.URL.RawQuery
	target = targetURL.String()

	log.Debug().Str("target", claims.Target).Msg("sending request to internal data server")

	httpClient := s.client
	httpReq, err := rhttp.NewRequest(ctx, "PUT", target, r.Body)
	if err != nil {
		log.Err(err).Msg("wrong request")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	httpReq.Header = r.Header

	httpRes, err := httpClient.Do(httpReq)
	if err != nil {
		log.Err(err).Msg("error doing PUT request to data service")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer httpRes.Body.Close()

	copyHeader(w.Header(), httpRes.Header)
	if httpRes.StatusCode != http.StatusOK {
		// swallow the body and set content-length to 0 to prevent reverse proxies from trying to read from it
		w.Header().Set("Content-Length", "0")
		w.WriteHeader(httpRes.StatusCode)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = io.Copy(w, httpRes.Body)
	if err != nil {
		log.Err(err).Msg("error writing body after header were set")
	}
}

// TODO: put and post code is pretty much the same. Should be solved in a nicer way in the long run.
func (s *svc) doPatch(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := appctx.GetLogger(ctx)

	claims, err := s.verify(ctx, r)
	if err != nil {
		err = errors.Wrap(err, "datagateway: error validating transfer token")
		log.Err(err).Str("token", r.Header.Get(TokenTransportHeader)).Msg("invalid transfer token")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	target := claims.Target
	// add query params to target, clients can send checksums and other information.
	targetURL, err := url.Parse(target)
	if err != nil {
		log.Err(err).Msg("datagateway: error parsing target url")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	targetURL.RawQuery = r.URL.RawQuery
	target = targetURL.String()

	log.Debug().Str("target", claims.Target).Msg("sending request to internal data server")

	httpClient := s.client
	httpReq, err := rhttp.NewRequest(ctx, "PATCH", target, r.Body)
	if err != nil {
		log.Err(err).Msg("wrong request")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	httpReq.Header = r.Header

	httpRes, err := httpClient.Do(httpReq)
	if err != nil {
		log.Err(err).Msg("error doing PATCH request to data service")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer httpRes.Body.Close()

	copyHeader(w.Header(), httpRes.Header)
	if httpRes.StatusCode != http.StatusOK && httpRes.StatusCode != http.StatusPartialContent {
		// swallow the body and set content-length to 0 to prevent reverse proxies from trying to read from it
		w.Header().Set("Content-Length", "0")
		w.WriteHeader(httpRes.StatusCode)
		return
	}

	w.WriteHeader(httpRes.StatusCode)
	_, err = io.Copy(w, httpRes.Body)
	if err != nil {
		log.Err(err).Msg("error writing body after header were set")
	}
}

func copyHeader(dst, src http.Header) {
	for key, values := range src {
		for i := range values {
			dst.Set(key, values[i])
		}
	}
}
