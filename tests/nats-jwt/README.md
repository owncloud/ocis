# NATS JWT tests

These scripts test oCIS against a stand-alone NATS server in operator (JWT) mode.
The embedded `nats` service does not support JWT mode. Thus oCIS runs with
`OCIS_EXCLUDE_RUN_SERVICES=nats`.

The scripts create new keys at each run. The repository never contains a seed.

## Files

| File | Purpose |
|---|---|
| `natscreds/main.go` | Creates `server.conf`, `user.creds` (account `ocis`), and `other.creds` (account `other`). |
| `lib.sh` | Image pins, variables, and the helper `nats_as`. |
| `start-nats.sh` | Creates the keys, starts the NATS container, and waits until JetStream is ready. |
| `verify.sh` | Checks the NATS state after a test run. |
| `startup-check.sh` | Starts oCIS with the services that are off by default and checks `/readyz`. |

## What the CI does

The jobs are in `.github/workflows/acceptance-tests.yml`.

| Job | Creds | What it covers |
|---|---|---|
| `natsJwt-events-stores-file` | file | `apiAsyncUpload`, `apiActivities`, `apiSearch1`: events of postprocessing, storage-users, activitylog, and search; stores of postprocessing and activitylog |
| `natsJwt-sharing-content` | content | `apiSharingNgShares`, `apiSharingNgItemLinkShare`: events of sharing; stores of storage-system and storage-publiclink |
| `natsJwt-notifications-file` | file | `apiNotification` with an email server: userlog, eventhistory, and notifications |
| `natsJwt-startup-check` | file | The `/readyz` endpoints and the connections of audit, policies, antivirus, and notifications |

All jobs also cover the service registry, the caches, and the connection of each default service.
`run-github.py` sets `EVENTHISTORY_STORE=memory`. Thus the store of the eventhistory service
uses NATS only in the startup check.
Each job sets `OCIS_CACHE_STORE=nats-js-kv` and `OCIS_PERSISTENT_STORE=nats-js-kv`. The caches
and the store of the userlog service use memory by default. These variables move them to NATS.
The other stores use NATS by default.

## What `verify.sh` checks

1. The NATS server log has no `authentication error`. A client without valid creds causes this
   line. The tests can pass when the client retries in the background, so only the log shows
   this bug type.
   Check 2 causes an error of this type on purpose. The script writes the number of these errors
   to `$NATS_JWT_DIR/probe-errors`, and check 1 ignores them. Thus you can run the script more than
   once against the same container.
2. A client without creds gets `Authorization Violation`. This proves that the server needs creds.
3. The KV buckets of the registry, the caches, and the stores exist in the account `ocis`.
4. The consumers of the stream `main-queue` exist. They have no pending message and no
   message without an acknowledgement.
5. The account `other` sees no bucket of oCIS.

The script runs all checks and exits with a non-zero code when one check fails.

## Run on a local machine

On Linux, with Docker and Go:

```bash
export NATS_JWT_DIR=/tmp/nats-jwt
tests/nats-jwt/start-nats.sh

export OCIS_EXCLUDE_RUN_SERVICES=nats
export OCIS_CACHE_STORE=nats-js-kv OCIS_PERSISTENT_STORE=nats-js-kv
export OCIS_NATS_CREDS_FILE="$NATS_JWT_DIR/user.creds"   # or: OCIS_NATS_CREDS="$(cat "$NATS_JWT_DIR/user.creds")"
BEHAT_SUITES=apiAsyncUpload python3 tests/acceptance/run-github.py

tests/nats-jwt/verify.sh
docker rm -f nats-jwt
```

The startup check needs a built binary (`make -C ocis build`) and assets (`make ci-node-generate`):

```bash
NATS_JWT_DIR=/tmp/nats-jwt tests/nats-jwt/startup-check.sh
```

On Docker Desktop (macOS), host networking does not reach the Mac. Use these variables:

```bash
export NATS_DOCKER_NETWORK=publish   # publish the ports 9233 and 8222
export NATS_CLI=nats                 # use the nats CLI of the Mac, not the nats-box image
```

`run-github.py` runs `ocis init` in `~/.ocis/config`. Run it with a separate `HOME`,
so that it does not change your own oCIS config.

## Extra services

`verify.sh` expects the buckets and the consumers of a default start. For the services you add
(for example with `OCIS_ADD_RUN_SERVICES`), set `NATS_JWT_EXTRA_BUCKETS` and
`NATS_JWT_EXTRA_CONSUMERS`. To replace the default lists, set `NATS_JWT_EXPECT_BUCKETS` and
`NATS_JWT_EXPECT_CONSUMERS`.

## Not covered

- NATS TLS together with JWT. The TLS change has its own tests.
- The OIDC userinfo cache of the proxy. The API test suites use basic authentication. Only an
  end-to-end test with an IDP login uses this cache.
- The suite `apiServiceAvailability`. It requests the embedded `nats` service on port 9234.
- Distributed setups with one process for each service.
