# LLM Proxy

The `llm` service proxies chat-completion requests from oCIS AI web
extensions to a configured, OpenAI-compatible LLM endpoint, so extensions
never see the upstream API key.

It replaces the standalone Node.js sidecar previously shipped in
`owncloud/web-extensions` (`packages/ai-llm-proxy`) — this service ships as
part of the oCIS single binary instead of a separate process/image.

The public endpoint is reached through the `proxy` service at
`POST /graph/v1beta1/extensions/org.libregraph/llm/chat/completions`.

`LLM_ENDPOINT` is a **required** environment variable; the service will not
start without it. `llm` is an **opt-in** service — it must be explicitly
enabled via `OCIS_ADD_RUN_SERVICES=llm` (the same way `antivirus`,
`auth-app`, and `invitations` are enabled) for `ocis server` to run it.

## Auth

The `llm` service does not authenticate callers itself. It trusts the
`proxy` service, which validates the caller's OIDC token and forwards a
signed identity via the `X-Access-Token` header before the request reaches
`llm`.

## Rate limiting

Requests are rate-limited per user using a sliding window backed by the
configured store (`LLM_STORE`, defaults to `nats-js-kv`), so the limit is
shared across replicas when `llm` runs multiple instances. Configure the
window and request ceiling with `LLM_RATE_LIMIT_WINDOW` /
`LLM_RATE_LIMIT_MAX_REQUESTS`.

## Configuration

See the `LLM_*` environment variables in `pkg/config` for the upstream
endpoint, API key, model override, timeout, and body-size limit.
