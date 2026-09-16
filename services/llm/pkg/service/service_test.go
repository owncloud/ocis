package service

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	userpb "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	revactx "github.com/owncloud/reva/v2/pkg/ctx"
	microstore "go-micro.dev/v4/store"

	"github.com/owncloud/ocis/v2/services/llm/pkg/config"
	"github.com/owncloud/ocis/v2/services/llm/pkg/ratelimit"
)

func testUser(id string) *userpb.User {
	return &userpb.User{Id: &userpb.UserId{OpaqueId: id}}
}

func newTestService(t *testing.T, llmServerURL string, maxRequests int) *Service {
	t.Helper()
	cfg := &config.Config{
		LLM: config.LLM{
			Endpoint:       llmServerURL,
			MaxTokensLimit: 4096,
			Timeout:        5 * time.Second,
			MaxBodyBytes:   1024,
		},
	}
	limiter := ratelimit.New(microstore.NewMemoryStore(), time.Minute, maxRequests)
	return New(cfg, limiter)
}

func requestWithUser(body []byte, userID string) *http.Request {
	req := httptest.NewRequest(http.MethodPost, "/chat/completions", bytes.NewReader(body))
	if userID != "" {
		ctx := revactx.ContextSetUser(req.Context(), testUser(userID))
		req = req.WithContext(ctx)
	}
	return req
}

func TestHandleChatCompletions_NoUser_Returns401(t *testing.T) {
	svc := newTestService(t, "http://unused.invalid", 20)
	req := requestWithUser([]byte(`{}`), "")
	rec := httptest.NewRecorder()

	svc.HandleChatCompletions(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
}

func TestHandleChatCompletions_RateLimited_Returns429(t *testing.T) {
	svc := newTestService(t, "http://unused.invalid", 1)

	// first request consumes the only allowed slot; the upstream call will
	// fail (no server listening at unused.invalid), which is fine — we only
	// care about the rate limiter state after this call.
	first := requestWithUser([]byte(`{"model":"m","messages":[]}`), "user-1")
	svc.HandleChatCompletions(httptest.NewRecorder(), first)

	second := requestWithUser([]byte(`{"model":"m","messages":[]}`), "user-1")
	rec := httptest.NewRecorder()
	svc.HandleChatCompletions(rec, second)

	if rec.Code != http.StatusTooManyRequests {
		t.Fatalf("expected 429, got %d", rec.Code)
	}
}

func TestHandleChatCompletions_OversizedBody_Returns413(t *testing.T) {
	svc := newTestService(t, "http://unused.invalid", 20)
	oversized := bytes.Repeat([]byte("a"), 2048)
	req := requestWithUser(oversized, "user-1")
	rec := httptest.NewRecorder()

	svc.HandleChatCompletions(rec, req)

	if rec.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("expected 413, got %d", rec.Code)
	}
}

func TestHandleChatCompletions_InvalidJSON_Returns400(t *testing.T) {
	svc := newTestService(t, "http://unused.invalid", 20)
	req := requestWithUser([]byte(`not json`), "user-1")
	rec := httptest.NewRecorder()

	svc.HandleChatCompletions(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestHandleChatCompletions_MissingMessages_Returns400(t *testing.T) {
	svc := newTestService(t, "http://unused.invalid", 20)
	req := requestWithUser([]byte(`{"model":"m"}`), "user-1")
	rec := httptest.NewRecorder()

	svc.HandleChatCompletions(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestHandleChatCompletions_HappyPath_RelaysUpstreamResponse(t *testing.T) {
	var receivedPath string
	var receivedAuth string
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		receivedAuth = r.Header.Get("Authorization")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"choices":[{"message":{"content":"hi"}}]}`))
	}))
	defer upstream.Close()

	svc := newTestService(t, upstream.URL, 20)
	svc.cfg.LLM.APIKey = "test-key"

	req := requestWithUser([]byte(`{"model":"m","messages":[{"role":"user","content":"hi"}]}`), "user-1")
	rec := httptest.NewRecorder()

	svc.HandleChatCompletions(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body: %s", rec.Code, rec.Body.String())
	}
	if receivedPath != "/chat/completions" {
		t.Fatalf("expected upstream path /chat/completions, got %q", receivedPath)
	}
	if receivedAuth != "Bearer test-key" {
		t.Fatalf("expected upstream Authorization header, got %q", receivedAuth)
	}

	var got map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("expected valid JSON relayed back, got error: %v, body: %s", err, rec.Body.String())
	}
}

func TestHandleChatCompletions_UpstreamUnreachable_Returns502(t *testing.T) {
	svc := newTestService(t, "http://127.0.0.1:1", 20) // nothing listens here
	req := requestWithUser([]byte(`{"model":"m","messages":[]}`), "user-1")
	rec := httptest.NewRecorder()

	svc.HandleChatCompletions(rec, req)

	if rec.Code != http.StatusBadGateway {
		t.Fatalf("expected 502, got %d", rec.Code)
	}
}

func TestHandleChatCompletions_ClientDisconnect_CancelsUpstreamCall(t *testing.T) {
	upstreamStarted := make(chan struct{})
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Drain the request body before waiting: net/http's server only starts
		// watching the connection for a client disconnect (which cancels
		// r.Context()) once the request body has been fully read. Without this,
		// r.Context().Done() below would never fire and the test would hang -
		// this is standard net/http behavior, not something the handler under
		// test can influence, since it is the *client* whose context we cancel.
		_, _ = io.Copy(io.Discard, r.Body)
		close(upstreamStarted)
		<-r.Context().Done() // the client cancellation below should propagate here
	}))
	defer upstream.Close()

	svc := newTestService(t, upstream.URL, 20)

	ctx, cancel := context.WithCancel(context.Background())
	req := httptest.NewRequest(http.MethodPost, "/chat/completions", bytes.NewReader([]byte(`{"model":"m","messages":[]}`)))
	req = req.WithContext(revactx.ContextSetUser(ctx, testUser("user-1")))

	done := make(chan struct{})
	go func() {
		svc.HandleChatCompletions(httptest.NewRecorder(), req)
		close(done)
	}()

	<-upstreamStarted
	cancel()

	select {
	case <-done:
		// handler returned once the client context was canceled, as expected
	case <-time.After(5 * time.Second):
		t.Fatal("handler did not return after client disconnect")
	}
}
