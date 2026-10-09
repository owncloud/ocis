package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	revactx "github.com/owncloud/reva/v2/pkg/ctx"

	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/services/llm/pkg/config"
	"github.com/owncloud/ocis/v2/services/llm/pkg/ratelimit"
	"github.com/owncloud/ocis/v2/services/llm/pkg/sanitize"
)

// maxUpstreamResponseBytes caps how much of the upstream LLM's response body
// this service will buffer before giving up.
const maxUpstreamResponseBytes = 10 * 1024 * 1024 // 10 MiB

// Service handles chat completion requests, proxying them to the configured
// upstream LLM endpoint.
type Service struct {
	cfg     *config.Config
	limiter *ratelimit.Limiter
	client  *http.Client
	log     log.Logger
}

// New creates a new Service.
func New(cfg *config.Config, limiter *ratelimit.Limiter, logger log.Logger) *Service {
	return &Service{
		cfg:     cfg,
		limiter: limiter,
		client:  &http.Client{},
		log:     logger,
	}
}

// HandleChatCompletions handles POST requests forwarding chat completions to
// the configured LLM.
func (s *Service) HandleChatCompletions(w http.ResponseWriter, r *http.Request) {
	user, ok := revactx.ContextGetUser(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "missing authenticated user")
		return
	}

	allowed, err := s.limiter.Allow(user.GetId().GetOpaqueId())
	if err != nil {
		s.log.Warn().Err(err).Msg("rate limit check failed")
		writeError(w, http.StatusInternalServerError, "rate limit check failed")
		return
	}
	if !allowed {
		writeError(w, http.StatusTooManyRequests, "rate limit exceeded. Please slow down.")
		return
	}

	body, err := readLimited(r.Body, s.cfg.LLM.MaxBodyBytes)
	if err != nil {
		if errors.Is(err, errBodyTooLarge) {
			writeError(w, http.StatusRequestEntityTooLarge, "request body too large")
		} else {
			writeError(w, http.StatusBadRequest, "could not read request body")
		}
		return
	}

	var raw map[string]interface{}
	if err := json.Unmarshal(body, &raw); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	sanitized, err := sanitize.Body(raw, s.cfg.LLM.Model, s.cfg.LLM.MaxTokensLimit)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	payload, err := json.Marshal(sanitized)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not encode upstream request")
		return
	}

	ctx, cancel := contextWithTimeout(r.Context(), s.cfg.LLM.Timeout)
	defer cancel()

	upstreamURL := strings.TrimSuffix(s.cfg.LLM.Endpoint, "/") + "/chat/completions"
	upstreamReq, err := http.NewRequestWithContext(ctx, http.MethodPost, upstreamURL, bytes.NewReader(payload))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not build upstream request")
		return
	}
	upstreamReq.Header.Set("Content-Type", "application/json")
	if s.cfg.LLM.APIKey != "" {
		upstreamReq.Header.Set("Authorization", "Bearer "+s.cfg.LLM.APIKey)
	}

	resp, err := s.client.Do(upstreamReq)
	if err != nil {
		s.log.Warn().Err(err).Msg("could not reach upstream LLM endpoint")
		writeError(w, http.StatusBadGateway, "could not reach LLM endpoint")
		return
	}
	defer resp.Body.Close()

	respBody, err := readLimited(resp.Body, maxUpstreamResponseBytes)
	if err != nil {
		s.log.Warn().Err(err).Msg("could not read upstream LLM response")
		writeError(w, http.StatusBadGateway, "could not read LLM response")
		return
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/json"
	}
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(resp.StatusCode)
	_, _ = w.Write(respBody)
}

func writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": message})
}

// contextWithTimeout wraps context.WithTimeout; a timeout <= 0 means "no
// explicit timeout beyond whatever the parent context already has" instead
// of expiring immediately.
func contextWithTimeout(parent context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	if timeout <= 0 {
		return parent, func() {}
	}
	return context.WithTimeout(parent, timeout)
}
