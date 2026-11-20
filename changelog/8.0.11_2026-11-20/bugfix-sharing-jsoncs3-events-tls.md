Bugfix: Turn on the events TLS in the jsoncs3 share manager and in OCM

The sharing service now gives the `OCIS_EVENTS_ENABLE_TLS` setting to the
jsoncs3 share manager. The ocm service now gives this setting to the
sciencemesh and the ocmcore event connections. Before, these connections did
not turn on TLS for the connection to the events broker. Thus they did not use
the root CA certificate. The connection to a NATS server with a private CA
failed. The services tried again for up to 15 minutes.

https://github.com/owncloud/ocis/pull/13442
