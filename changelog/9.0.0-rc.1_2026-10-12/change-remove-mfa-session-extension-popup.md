Change: Remove the MFA session extension popup in vault mode

In vault mode, the web UI showed a popup shortly before the multi-factor
authentication session was about to expire, offering to extend or dismiss it.
Extending only reset a timer in the browser and did not renew anything on the
server, so the popup has been removed.

The frontend service no longer exposes the MFA session duration in its
capabilities, and its `session_duration` setting has been removed. The
`OCIS_MFA_SESSION_DURATION` environment variable is still used by the proxy
service.

https://github.com/owncloud/ocis/pull/13085
