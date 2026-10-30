#!/usr/bin/env bash
# Check the state of the NATS JWT account after an oCIS test run.
#
# The script runs all checks and prints each result. It exits with a non-zero code
# when one check fails. The checks do not need a running oCIS.
#
# Variables (see also lib.sh):
#   NATS_JWT_EXPECT_BUCKETS     KV buckets that must exist (default: the buckets of a run of run-github.py)
#   NATS_JWT_EXPECT_CONSUMERS   consumers of the stream main-queue that must exist
#   NATS_JWT_EXTRA_BUCKETS      more buckets, for the services you added
#   NATS_JWT_EXTRA_CONSUMERS    more consumers, for the services you added
#   NATS_JWT_STREAM             name of the events stream (default: main-queue)
set -uo pipefail
cd "$(dirname "$0")"
# shellcheck source=lib.sh
. ./lib.sh

: "${NATS_JWT_DIR:?set NATS_JWT_DIR}"
STREAM="${NATS_JWT_STREAM:-main-queue}"
# shellcheck disable=SC2206
BUCKETS=(${NATS_JWT_EXPECT_BUCKETS:-service-registry cache-providers cache-createhome settings-cache storage-system storage-users ids-storage-users postprocessing activitylog} ${NATS_JWT_EXTRA_BUCKETS:-})
# shellcheck disable=SC2206
CONSUMERS=(${NATS_JWT_EXPECT_CONSUMERS:-activitylog clientlog dcfs evhistory frontend graph jsoncs3sharemanager postprocessing search sse storage-users userlog} ${NATS_JWT_EXTRA_CONSUMERS:-})

failed=0
pass() { echo "ok:   $*"; }
fail() { echo "FAIL: $*"; failed=1; }

AUTH_ERRORS='authentication error|authorization violation'
# Check 2 sends a request without creds on purpose, and the server logs an error for it.
# This file keeps the number of these errors, so that check 1 also passes on the next run.
# start-nats.sh removes the file when it starts a new container.
PROBE_FILE="$NATS_JWT_DIR/probe-errors"

# count_auth_errors prints the number of authentication errors in the NATS server log.
count_auth_errors() {
  local log
  log="$(docker logs "$NATS_CONTAINER" 2>&1)" || return 1
  grep -ciE "$AUTH_ERRORS" <<<"$log" || true
}

# 1. Server log. Run this check before check 2.
# A client without valid creds makes the server log "authentication error". The tests can pass
# when the client retries in the background, so only the log shows this type of bug.
known="$(cat "$PROBE_FILE" 2>/dev/null || echo 0)"
if total="$(count_auth_errors)"; then
  violations=$((total - known))
  if [ "$violations" -le 0 ]; then
    pass "the NATS server log has no authentication error"
  else
    fail "the NATS server log has $violations authentication errors. A client connected without valid creds:"
    [ "$known" -gt 0 ] && echo "(the list can also show $known errors from the check 2 of the earlier runs)"
    docker logs "$NATS_CONTAINER" 2>&1 | grep -iE "$AUTH_ERRORS" | head -10
  fi
else
  fail "cannot read the log of the container $NATS_CONTAINER"
  total=0
fi

# 2. The server needs creds. Without this proof, a server with no auth makes the job green.
out="$(nats_as none account info 2>&1)"
if [ $? -ne 0 ] && grep -qi 'authorization violation' <<<"$out"; then
  pass "a client without creds gets Authorization Violation"
else
  fail "a client without creds can use the server. Output: $out"
fi
# The server can write the log line after the client stops. Wait up to 5 s for it.
after="$total"
for _ in $(seq 1 5); do
  after="$(count_auth_errors || echo "$total")"
  [ "$after" -gt "$total" ] && break
  sleep 1
done
echo $((known + after - total)) >"$PROBE_FILE"

# 3. KV buckets of the registry, the caches, and the stores are in the account "ocis".
buckets="$(nats_as user kv ls --names 2>&1)"
for b in "${BUCKETS[@]}"; do
  if grep -qx "$b" <<<"$buckets"; then pass "KV bucket $b exists"; else fail "KV bucket $b is missing"; fi
done

# 4. Events: the consumers exist, and they have no pending message at the end. A message is
# pending when the consumer did not get it (num_pending) or did not acknowledge it (num_ack_pending).
names="$(nats_as user consumer ls "$STREAM" --names 2>&1)"
for c in "${CONSUMERS[@]}"; do
  if grep -qE "^$c(-[0-9a-f-]+)?$" <<<"$names"; then pass "consumer $c exists"; else fail "consumer $c is missing in $STREAM"; fi
done
# A consumer can still process the last events. Wait up to 30 s.
for _ in $(seq 1 15); do
  pending=0
  while read -r n; do
    [ -z "$n" ] && continue
    p="$(nats_as user consumer info "$STREAM" "$n" --json 2>/dev/null | jq -r '(.num_pending // 0) + (.num_ack_pending // 0)')"
    pending=$((pending + p))
  done <<<"$names"
  [ "$pending" = 0 ] && break
  sleep 2
done
if [ "$pending" = 0 ]; then pass "all consumers of $STREAM have no pending message"; else fail "consumers of $STREAM have $pending pending messages"; fi

# 5. The data is in the account "ocis" only.
other="$(nats_as other kv ls --names 2>&1)"
if grep -qE '^(service-registry|storage-users|postprocessing)$' <<<"$other"; then
  fail "the account other can see the buckets of oCIS"
else
  pass "the account other sees no bucket of oCIS"
fi

if [ "$failed" = 0 ]; then log "all checks passed"; else log "checks failed"; fi
exit "$failed"
