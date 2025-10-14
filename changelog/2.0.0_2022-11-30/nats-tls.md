Enhancement: Secure the nats connection with TLS

Encrypted the connection to the event broker using TLS.
Per default TLS is not enabled but can be enabled by setting either `OCIS_EVENTS_ENABLE_TLS=true` or the respective service configs:

- `AUDIT_EVENTS_ENABLE_TLS=true`
- `GRAPH_EVENTS_ENABLE_TLS=true`
- `NATS_EVENTS_ENABLE_TLS=true`
- `NOTIFICATIONS_EVENTS_ENABLE_TLS=true`
- `SEARCH_EVENTS_ENABLE_TLS=true`
- `SHARING_EVENTS_ENABLE_TLS=true`
- `STORAGE_USERS_EVENTS_ENABLE_TLS=true`

https://github.com/owncloud/ocis/pull/4781
https://github.com/owncloud/ocis/pull/4800
https://github.com/owncloud/ocis/pull/4867
