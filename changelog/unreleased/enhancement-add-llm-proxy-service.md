Enhancement: Add LLM proxy service

oCIS now ships a native `llm` microservice, as part of the single binary,
that proxies chat-completion requests from AI web extensions to a
configured OpenAI-compatible LLM endpoint. It replaces the standalone
Node.js sidecar (`ai-llm-proxy`) previously required in
`owncloud/web-extensions`, so AI-enabled oCIS deployments no longer need a
separate process or image.

The service trusts oCIS's existing authentication (no independent OIDC
validation), rate-limits requests per user using oCIS's shared store
abstraction (correct across multiple replicas), and forwards a sanitized
request body — only the fields the LLM needs are passed through.

https://github.com/owncloud/ocis/pull/CHANGE_ME
