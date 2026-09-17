#!/usr/bin/env bash
# Bootstrap an external osixia LDAP + oCIS and run the LDAP smoke test.
#
# oCIS talks to a real external LDAP (no embedded idm), with server-assigned entryUUID
# as the id attribute. This exercises the OCISDEV-1030/1032/1033 LDAP client changes
# (configurable read/write-aware retry + read-after-write) against an external
# directory — the path CI otherwise never covers (run-github.py always uses idm).
#
# Two-phase start: the id attribute is entryUUID (server-assigned), so the admin's id
# is only known after LDAP seeds. Start LDAP, read the admin entryUUID, then start oCIS
# with OCIS_ADMIN_USER_ID pinned to it — otherwise the admin can't be resolved and role
# assignment / login break.
#
# Needs: docker, curl, jq, and a built binary at ocis/bin/ocis (make -C ocis build).
set -euo pipefail
cd "$(dirname "$0")"
DIR="$(pwd)"
REPO_ROOT="$(cd ../.. && pwd)"

# Reuse the osixia seed/schema from the ocis_ldap deployment example (same demo users,
# ou structure and memberof/refint overlays) instead of duplicating the LDIFs here.
# The entryUUID / plaintext specifics that this smoke needs live in the oCIS env below,
# not in the seed, so the example's LDIFs are topology-agnostic and reused as-is.
LDAP_CONFIG="$REPO_ROOT/deployments/examples/ocis_ldap/config/ldap"

OCIS_BIN="${OCIS_BIN:-$REPO_ROOT/ocis/bin/ocis}"
CONTAINER="ldap-smoke"
ADMIN_DN="uid=admin,ou=users,dc=owncloud,dc=com"
LDAP_ADMIN_PASSWORD="${LDAP_ADMIN_PASSWORD:-admin}"
export OCIS_URL="${OCIS_URL:-https://localhost:9200}"

[ -x "$OCIS_BIN" ] || { echo "ERROR: ocis binary not found at $OCIS_BIN (run: make -C ocis build)"; exit 1; }

OCIS_PID=""
cleanup() {
  set +e
  [ -n "$OCIS_PID" ] && kill "$OCIS_PID" 2>/dev/null
  docker rm -f "$CONTAINER" >/dev/null 2>&1
}
trap cleanup EXIT

# Fresh LDAP each run so the seed (which runs only on first init) always applies.
docker rm -f "$CONTAINER" >/dev/null 2>&1 || true

echo "[*] Starting external LDAP ($CONTAINER) ..."
docker run -d --name "$CONTAINER" -p 127.0.0.1:389:389 \
  -e LDAP_TLS=false \
  -e LDAP_ORGANISATION=owncloud \
  -e LDAP_DOMAIN=owncloud.com \
  -e LDAP_ROOT="dc=owncloud,dc=com" \
  -e LDAP_ADMIN_PASSWORD="$LDAP_ADMIN_PASSWORD" \
  -e LDAP_SEED_INTERNAL_LDIF_PATH=/ldifs \
  -e LDAP_SEED_INTERNAL_SCHEMA_PATH=/schemas \
  -v "$LDAP_CONFIG/ldif:/ldifs:ro" \
  -v "$LDAP_CONFIG/schemas:/schemas:ro" \
  osixia/openldap:1.5.0 --copy-service --loglevel info >/dev/null

echo -n "[*] Waiting for LDAP seed "
UUID=""
for _ in $(seq 1 60); do
  # `|| true`: while LDAP is still seeding, ldapsearch returns non-zero (no such
  # object); with set -e + pipefail that would kill the poll on the first iteration.
  UUID=$( { docker exec "$CONTAINER" ldapsearch -LLL -o ldif-wrap=no -x \
    -H ldap://localhost:389 -D "cn=admin,dc=owncloud,dc=com" -w "$LDAP_ADMIN_PASSWORD" \
    -b "$ADMIN_DN" -s base entryUUID 2>/dev/null || true; } \
    | awk '/^entryUUID:/ {print $2; exit}')
  if [ -n "$UUID" ]; then echo " ok"; break; fi
  printf "."
  sleep 1
done
[ -n "$UUID" ] || { echo " FAILED (no admin entryUUID)"; docker logs --tail=40 "$CONTAINER"; exit 1; }
echo "[*] admin entryUUID = $UUID"

echo "[*] ocis init ..."
export OCIS_CONFIG_DIR="${OCIS_CONFIG_DIR:-$HOME/.ocis/config}"
"$OCIS_BIN" init --insecure true || true

echo "[*] Starting oCIS (OCIS_ADMIN_USER_ID=$UUID) ..."
env \
  OCIS_LDAP_URI="ldap://localhost:389" \
  OCIS_LDAP_INSECURE="true" \
  OCIS_LDAP_BIND_DN="cn=admin,dc=owncloud,dc=com" \
  OCIS_LDAP_BIND_PASSWORD="$LDAP_ADMIN_PASSWORD" \
  OCIS_LDAP_GROUP_BASE_DN="ou=groups,dc=owncloud,dc=com" \
  OCIS_LDAP_GROUP_FILTER="(objectclass=owncloud)" \
  OCIS_LDAP_GROUP_OBJECTCLASS="groupOfNames" \
  OCIS_LDAP_GROUP_ADDITIONAL_OBJECTCLASSES="owncloud" \
  OCIS_LDAP_USER_BASE_DN="ou=users,dc=owncloud,dc=com" \
  OCIS_LDAP_USER_FILTER="(objectclass=owncloud)" \
  OCIS_LDAP_USER_OBJECTCLASS="inetOrgPerson" \
  LDAP_LOGIN_ATTRIBUTES="uid" \
  OCIS_LDAP_USER_SCHEMA_ID="entryUUID" \
  OCIS_LDAP_GROUP_SCHEMA_ID="entryUUID" \
  OCIS_ADMIN_USER_ID="$UUID" \
  IDP_LDAP_LOGIN_ATTRIBUTE="uid" \
  IDP_LDAP_UUID_ATTRIBUTE="entryUUID" \
  IDP_LDAP_UUID_ATTRIBUTE_TYPE="text" \
  GRAPH_LDAP_SERVER_WRITE_ENABLED="true" \
  GRAPH_LDAP_REFINT_ENABLED="true" \
  GRAPH_LDAP_SERVER_UUID="true" \
  GRAPH_LDAP_SERVER_USE_PASSWORD_MODIFY_EXOP="false" \
  OCIS_EXCLUDE_RUN_SERVICES="idm" \
  OCIS_URL="$OCIS_URL" \
  OCIS_LOG_LEVEL="${OCIS_LOG_LEVEL:-error}" \
  OCIS_LOG_COLOR="false" \
  OCIS_INSECURE="true" \
  PROXY_ENABLE_BASIC_AUTH="true" \
  "$OCIS_BIN" server &
OCIS_PID=$!

echo -n "[*] Waiting for oCIS "
# Gate on GET /graph/v1.0/users/admin -> 200 (as run-github.py does), not on
# .well-known: the latter only proves the IDP is up, while the graph admin lookup
# exercises the full proxy -> settings (roles) -> graph -> LDAP path, so it only
# succeeds once the stack — including admin role assignment — is actually ready.
ready=""
for _ in $(seq 1 180); do
  code=$(curl -sk -o /dev/null -w "%{http_code}" -u "admin:admin" \
    "$OCIS_URL/graph/v1.0/users/admin" 2>/dev/null || true)
  if [ "$code" = "200" ]; then echo " ready"; ready=1; break; fi
  # bail early if the server process died
  kill -0 "$OCIS_PID" 2>/dev/null || { echo " FAILED (ocis exited)"; exit 1; }
  printf "."
  sleep 1
done
[ -n "$ready" ] || { echo " FAILED (timeout)"; exit 1; }

echo "[*] Running smoke.sh ..."
./smoke.sh
