Bugfix: Users without a matching role claim were stuck in a login loop

Users whose OIDC role claim was missing, unreadable, or matched no entry in the proxy's
`role_mapping` could not log in. The proxy answered with a bare `500 Internal Server
Error`, so the web UI showed an access denied page whose "log in again" button returned
to the same page, because the login itself had succeeded. The only way out was clearing
the browser's cookies. This is common when users are federated into the IDP from an
external user directory and simply have no role attached.

Such a request is now answered with `403 Forbidden`, and the logged error names the
claim that was read, the values found in it and the configured role mapping.

We've also added `PROXY_ROLE_ASSIGNMENT_OIDC_DEFAULT_ROLE`, which names an ocis role to
assign to those users instead of refusing them. It applies when the claim is missing,
when it cannot be read, and when it matches no mapping; a mapping that does match always
wins over it. It is empty by default, so no existing deployment changes behavior.

https://github.com/owncloud/ocis/issues/11467
