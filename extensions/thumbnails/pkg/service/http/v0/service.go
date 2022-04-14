package svc

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v4"

	"github.com/owncloud/ocis/extensions/thumbnails/pkg/config"
	tjwt "github.com/owncloud/ocis/extensions/thumbnails/pkg/service/jwt"
	"github.com/owncloud/ocis/extensions/thumbnails/pkg/thumbnail"
	"github.com/owncloud/ocis/ocis-pkg/log"
)

type contextKey string

const (
	keyContextKey contextKey = "key"
)

// Service defines the extension handlers.
type Service interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
	GetThumbnail(http.ResponseWriter, *http.Request)
}

// NewService returns a service implementation for Service.
func NewService(opts ...Option) Service {
	options := newOptions(opts...)

	m := chi.NewMux()
	m.Use(options.Middleware...)

	logger := options.Logger
	resolutions, err := thumbnail.ParseResolutions(options.Config.Thumbnail.Resolutions)
	if err != nil {
		logger.Fatal().Err(err).Msg("resolutions not configured correctly")
	}
	svc := Thumbnails{
		config: options.Config,
		mux:    m,
		logger: options.Logger,
		manager: thumbnail.NewSimpleManager(
			resolutions,
			options.ThumbnailStorage,
			logger,
		),
	}

	m.Route(options.Config.HTTP.Root, func(r chi.Router) {
		r.Use(svc.TransferTokenValidator)
		r.Get("/data", svc.GetThumbnail)
	})

	return svc
}

// Thumbnails implements the business logic for Service.
type Thumbnails struct {
	config  *config.Config
	logger  log.Logger
	mux     *chi.Mux
	manager thumbnail.Manager
}

// ServeHTTP implements the Service interface.
func (s Thumbnails) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

// GetThumbnail implements the Service interface.
func (s Thumbnails) GetThumbnail(w http.ResponseWriter, r *http.Request) {
	key := r.Context().Value(keyContextKey).(string)

	thumbnail, err := s.manager.GetThumbnail(key)
	if err != nil {
		s.logger.Error().
			Err(err).
			Str("key", key).
			Msg("could not get the thumbnail")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Length", strconv.Itoa(len(thumbnail)))
	if _, err = w.Write(thumbnail); err != nil {
		s.logger.Error().
			Err(err).
			Str("key", key).
			Msg("could not write the thumbnail response")
	}
}

func (s Thumbnails) TransferTokenValidator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Transfer-Token")
		token, err := jwt.ParseWithClaims(tokenString, &tjwt.ThumbnailClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(s.config.Thumbnail.TransferTokenSecret), nil
		})
		if err != nil {
			s.logger.Error().
				Err(err).
				Str("transfer-token", tokenString).
				Msg("failed to parse transfer token")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(*tjwt.ThumbnailClaims); ok && token.Valid {
			ctx := context.WithValue(r.Context(), keyContextKey, claims.Key)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		w.WriteHeader(http.StatusUnauthorized)
	})
}
