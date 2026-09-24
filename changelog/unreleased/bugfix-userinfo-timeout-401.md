Bugfix: Return a retryable 503 when the OIDC userinfo call fails transiently

A transient failure of the OIDC userinfo call (a timeout, a network error or a
5xx/429 from the IdP) was mapped to HTTP 401. Clients read the 401 as an invalid
session and logged the user out on a brief IdP slowdown.

The proxy now distinguishes a transient IdP failure from an authentication
failure and returns a retryable 503 (with Retry-After) for the former, so clients
retry and keep their session. A genuinely invalid or expired token still returns
401.

https://kiteworks.atlassian.net/browse/OCISDEV-1411
