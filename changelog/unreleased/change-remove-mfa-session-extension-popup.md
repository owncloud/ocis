Change: Remove the MFA session extension popup in vault mode

In vault mode, the web UI showed a popup shortly before the multi-factor
authentication session was about to expire, offering to extend or dismiss it.
Extending only reset a timer in the browser and did not renew anything on the
server, so the popup has been removed.

https://github.com/owncloud/ocis/pull/13085
