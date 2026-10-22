Security: Enforce authentication on the bundled NATS event bus

The embedded NATS server was started without any authentication, so any client
able to reach the port could subscribe to and publish forged internal events
(postprocessing/antivirus verdicts, upload finalize/revert, notifications).
The server now enforces the existing event-bus credentials
(OCIS_EVENTS_AUTH_USERNAME / OCIS_EVENTS_AUTH_PASSWORD), warns when bound to a
non-loopback address without authentication, and the example deployments no
longer expose the broker without credentials.

https://github.com/owncloud/ocis/pull/12317
