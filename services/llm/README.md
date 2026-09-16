# LLM Proxy

The `llm` service proxies chat-completion requests from oCIS AI web
extensions to a configured, OpenAI-compatible LLM endpoint, so extensions
never see the upstream API key.

It replaces the standalone Node.js sidecar previously shipped in
`owncloud/web-extensions` (`packages/ai-llm-proxy`) — this service ships as
part of the oCIS single binary instead of a separate process/image.

## Auth

The `llm` service does not authenticate callers itself. It trusts the
`proxy` service, which validates the caller's OIDC token and forwards a
signed identity via the `X-Access-Token` header before the request reaches
`llm`.

## Rate limiting

Requests are rate-limited per user using a sliding window backed by the
configured store (`LLM_STORE`, defaults to `nats-js-kv`), so the limit is
correct even when `llm` runs multiple replicas. Configure the window and
request ceiling with `LLM_RATE_LIMIT_WINDOW` / `LLM_RATE_LIMIT_MAX_REQUESTS`.

## Configuration

See the `LLM_*` environment variables in `pkg/config` for the upstream
endpoint, API key, model override, timeout, and body-size limit.
