Bugfix: Wait for app registry before running WOPI cs3 validator tests

The WOPI cs3 validator occasionally failed with HTTP 500 on `OpenInApp`
calls during the `EditFlows` test group. The gateway logged
`error finding app providers` — meaning the app-provider service had not
yet finished registering its mime-type entries with the gateway, even
though the `ocis_healthy()` check (a graph API call) had already passed.

A new `app_registry_ready()` probe is added to `run-wopi.py` that calls
`GET /app/list` and verifies the response contains at least one mime-type
entry. This probe runs after `ocis_healthy()` and before the WOPI service
is started, ensuring the gateway's app provider list is populated before
the validator begins issuing `OpenInApp` requests.

https://github.com/owncloud/ocis/pull/13504
