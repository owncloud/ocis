#!/usr/bin/env bash
# Start oCIS against a stand-alone NATS server in JWT mode, with the services that are off by
# default, and check that each service connects with the creds.
#
# What the script checks:
#   1. The /readyz endpoint of each service that checks NATS returns 200. The health check
#      uses the creds, so a check without creds fails.
#   2. verify.sh: no authentication error in the NATS log, the expected KV buckets, and the
#      expected consumers, including those of audit, policies, antivirus, and notifications.
#
# No test suite uses audit, policies, antivirus, and notifications. This check shows that
# these services at least connect.
#
# Needs: docker, curl, jq, and a built binary at ocis/bin/ocis (make -C ocis build).
# Variables:
#   NATS_JWT_CREDS  "file" (default) or "content": how oCIS gets the creds
#   OCIS_BIN        path of the ocis binary
#   EXCLUDE_EXTRA   more services to exclude, for example "idp" for a build without web assets
# See also lib.sh.
set -euo pipefail
cd "$(dirname "$0")"
# shellcheck source=lib.sh
. ./lib.sh

REPO_ROOT="$(cd ../.. && pwd)"
OCIS_BIN="${OCIS_BIN:-$REPO_ROOT/ocis/bin/ocis}"
NATS_JWT_CREDS="${NATS_JWT_CREDS:-file}"
# Without NATS_JWT_DIR, the script makes a temporary folder. cleanup removes it, because it has seeds.
NATS_TMP=""
if [ -z "${NATS_JWT_DIR:-}" ]; then
  NATS_TMP="$(mktemp -d)"
  NATS_JWT_DIR="$NATS_TMP/nats-jwt"
fi
export NATS_JWT_DIR
WORK="$(mktemp -d)"
OCIS_PID=""

[ -x "$OCIS_BIN" ] || { echo "ERROR: ocis binary not found at $OCIS_BIN (run: make -C ocis build)"; exit 1; }

cleanup() {
  set +e
  [ -n "$OCIS_PID" ] && kill "$OCIS_PID" 2>/dev/null
  docker rm -f "$NATS_CONTAINER" >/dev/null 2>&1
  rm -rf "$WORK"
  [ -n "$NATS_TMP" ] && rm -rf "$NATS_TMP"
}
trap cleanup EXIT

./start-nats.sh

export OCIS_CONFIG_DIR="$WORK/config" OCIS_BASE_DATA_PATH="$WORK/data"
mkdir -p "$OCIS_BASE_DATA_PATH"
"$OCIS_BIN" init --insecure true --admin-password admin >/dev/null

export OCIS_URL="${OCIS_URL:-https://localhost:9200}"
export OCIS_INSECURE=true
export OCIS_LOG_LEVEL=warn
export PROXY_ENABLE_BASIC_AUTH=true
# The embedded nats service does not support JWT mode.
export OCIS_EXCLUDE_RUN_SERVICES="nats${EXCLUDE_EXTRA:+,$EXCLUDE_EXTRA}"
export OCIS_ADD_RUN_SERVICES=audit,policies,antivirus,notifications
# The caches and the userlog store use memory by default. These variables move them to NATS.
# The other stores use NATS by default.
export OCIS_CACHE_STORE=nats-js-kv
export OCIS_PERSISTENT_STORE=nats-js-kv
export AUDIT_DEBUG_ADDR=0.0.0.0:9229
export POLICIES_DEBUG_ADDR=0.0.0.0:9129
export ANTIVIRUS_DEBUG_ADDR=0.0.0.0:9277
export NOTIFICATIONS_DEBUG_ADDR=0.0.0.0:9174
# The notifications service stops at startup without a valid sender. No mail server is needed.
export NOTIFICATIONS_SMTP_HOST=localhost
export NOTIFICATIONS_SMTP_PORT=1025
export NOTIFICATIONS_SMTP_INSECURE=true
export NOTIFICATIONS_SMTP_SENDER="ownCloud <noreply@example.com>"
if [ "$NATS_JWT_CREDS" = content ]; then
  OCIS_NATS_CREDS="$(cat "$NATS_JWT_DIR/user.creds")"
  export OCIS_NATS_CREDS
else
  export OCIS_NATS_CREDS_FILE="$NATS_JWT_DIR/user.creds"
fi

log "starting oCIS with the creds as $NATS_JWT_CREDS"
"$OCIS_BIN" server >"$WORK/ocis.log" 2>&1 &
OCIS_PID=$!

# Services with a NATS check in /readyz. The ports are the default debug ports.
# antivirus and notifications are not in the list: their /readyz also needs a ClamAV scanner and
# an SMTP server. Their consumers in verify.sh show that they connect.
READY=(
  "sse:9139" "sharing:9151" "proxy:9205" "frontend:9141" "graph:9124" "eventhistory:9270"
  "ocs:9114" "search:9224" "clientlog:9260" "postprocessing:9255" "activitylog:9197"
  "userlog:9214" "ocm:9281" "gateway:9143" "storage-users:9159"
  "audit:9229" "policies:9129"
)

failed=0
deadline=$((SECONDS + 240))
for entry in "${READY[@]}"; do
  name="${entry%%:*}"
  port="${entry##*:}"
  code=000
  while [ $SECONDS -lt $deadline ]; do
    kill -0 "$OCIS_PID" 2>/dev/null || { echo "FAIL: oCIS stopped"; tail -20 "$WORK/ocis.log"; exit 1; }
    code="$(curl -s -o /dev/null -m 5 -w '%{http_code}' "http://127.0.0.1:$port/readyz" || true)"
    [ "$code" = 200 ] && break
    sleep 2
  done
  if [ "$code" = 200 ]; then
    echo "ok:   $name /readyz returns 200"
  else
    echo "FAIL: $name /readyz returns $code (port $port)"
    failed=1
  fi
done

# The KV buckets are created at the first use. A login, an upload, and a download create them.
log "creating traffic"
api="https://localhost:9200"
# request <label> <expected codes as a regex> <curl arguments...>
request() {
  local label="$1" want="$2" code
  shift 2
  code="$(curl -sk -u admin:admin -o /dev/null -m 30 -w '%{http_code}' "$@" || true)"
  if [[ "$code" =~ ^($want)$ ]]; then
    echo "ok:   $label returns $code"
  else
    echo "FAIL: $label returns $code (expected $want)"
    failed=1
  fi
}
request "me/drives" 200 "$api/graph/v1.0/me/drives"
request "upload" "201|204" -T /etc/hosts "$api/remote.php/dav/files/admin/hosts.txt"
sleep 5
request "download" 200 "$api/remote.php/dav/files/admin/hosts.txt"

# The consumers of the added services exist after the services started. run-github.py sets
# EVENTHISTORY_STORE=memory, so only this script has the bucket of the eventhistory service.
NATS_JWT_EXTRA_BUCKETS="${NATS_JWT_EXTRA_BUCKETS:-eventhistory}" \
NATS_JWT_EXTRA_CONSUMERS="${NATS_JWT_EXTRA_CONSUMERS:-antivirus audit notifications policies}" \
  ./verify.sh || failed=1

if [ "$failed" != 0 ]; then
  log "oCIS log (last lines):"
  tail -30 "$WORK/ocis.log"
  log "NATS log (last lines):"
  docker logs --tail 20 "$NATS_CONTAINER" 2>&1 || true
  exit 1
fi
log "startup check passed"
