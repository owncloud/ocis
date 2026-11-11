Change: Respect proxy environment variables for OIDC requests

The proxy and webfinger services now route OIDC requests through the
`HTTP_PROXY` and `HTTPS_PROXY` settings from the environment. Deployments that
should not send the OIDC issuer host through the proxy, or whose proxy cannot
reach that host, must list the issuer hostname in `NO_PROXY`, otherwise
authentication can fail when oCIS's own external hostname is reached through
the proxy.

https://github.com/owncloud/ocis/pull/12703
