#!/usr/bin/env bash
# Proxy/write/read LDAP topology (one-way async syncrepl) + oCIS, run through the proxy.
# oCIS talks only to proxy-ldap: reads served from a local syncrepl replica, writes chained
# to write-ldap. Async lag exercises the reva retry / readBackAfterWrite paths
# (OCISDEV-1030/1032/1033). oCIS runs as a host binary against the proxy's published port.
# id = entryUUID (server-assigned), so admin id is only known after seed+replication.
#
# Start order:
#   1. write-ldap seeds demo data; read admin entryUUID from it.
#   2. read-ldap + proxy-ldap pull via syncrepl; wait until admin is visible through
#      proxy-ldap (replication converged before oCIS reads via it).
#   3. oCIS starts with OCIS_ADMIN_USER_ID pinned to that entryUUID, else role
#      assignment / login break.
#
# Needs: docker (compose plugin), curl, jq, and ocis/bin/ocis (make -C ocis build).
set -euo pipefail
cd "$(dirname "$0")"
DIR="$(pwd)"
REPO_ROOT="$(cd ../.. && pwd)"

OCIS_BIN="${OCIS_BIN:-$REPO_ROOT/ocis/bin/ocis}"
ADMIN_DN="uid=admin,ou=users,dc=owncloud,dc=com"
LDAP_ADMIN_PASSWORD="${LDAP_ADMIN_PASSWORD:-admin}"
export OCIS_URL="${OCIS_URL:-https://localhost:9200}"

[ -x "$OCIS_BIN" ] || { echo "ERROR: ocis binary not found at $OCIS_BIN (run: make -C ocis build)"; exit 1; }

OCIS_PID=""
cleanup() {
  set +e
  [ -n "$OCIS_PID" ] && kill "$OCIS_PID" 2>/dev/null
  docker compose down -v >/dev/null 2>&1
}
trap cleanup EXIT

# read_admin_uuid SERVICE -> admin entryUUID seen by that service (empty if absent).
# `|| true`: ldapsearch returns 32 (no such object) while seeding/replicating; set -e +
# pipefail would kill the poll on the first iteration.
read_admin_uuid() {
  local svc=$1
  { docker compose exec -T "$svc" ldapsearch -LLL -o ldif-wrap=no -x \
      -H ldap://localhost:389 -D "cn=admin,dc=owncloud,dc=com" -w "$LDAP_ADMIN_PASSWORD" \
      -b "$ADMIN_DN" -s base entryUUID 2>/dev/null || true; } \
    | awk '/^entryUUID:/ {print $2; exit}'
}

# Fresh volumes each run; the seed only runs on first init.
docker compose down -v >/dev/null 2>&1 || true

echo "[*] Starting write-ldap (provider) ..."
docker compose up -d write-ldap

echo -n "[*] Waiting for write-ldap seed "
UUID=""
for _ in $(seq 1 60); do
  UUID=$(read_admin_uuid write-ldap)
  if [ -n "$UUID" ]; then echo " ok"; break; fi
  printf "."
  sleep 1
done
[ -n "$UUID" ] || { echo " FAILED (no admin entryUUID on write-ldap)"; docker compose logs --tail=40 write-ldap; exit 1; }
echo "[*] admin entryUUID = $UUID"

echo "[*] Starting read-ldap + proxy-ldap (syncrepl consumers) ..."
docker compose up -d read-ldap proxy-ldap

echo -n "[*] Waiting for syncrepl to converge on proxy-ldap "
PUUID=""
for _ in $(seq 1 60); do
  PUUID=$(read_admin_uuid proxy-ldap)
  # Must match the provider's UUID; entryUUID is preserved across replication.
  if [ "$PUUID" = "$UUID" ]; then echo " ok"; break; fi
  printf "."
  sleep 1
done
[ "$PUUID" = "$UUID" ] || { echo " FAILED (admin not replicated to proxy-ldap; got '$PUUID')"; docker compose logs --tail=40 proxy-ldap; exit 1; }

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
  OCIS_LDAP_RETRY_MAX_COUNT="10" \
  OCIS_LDAP_RETRY_BASE_DELAY="200ms" \
  OCIS_LDAP_RETRY_MAX_DELAY="2s" \
  OCIS_EXCLUDE_RUN_SERVICES="idm" \
  OCIS_URL="$OCIS_URL" \
  OCIS_LOG_LEVEL="${OCIS_LOG_LEVEL:-error}" \
  OCIS_LOG_COLOR="false" \
  OCIS_INSECURE="true" \
  PROXY_ENABLE_BASIC_AUTH="true" \
  "$OCIS_BIN" server &
OCIS_PID=$!

echo -n "[*] Waiting for oCIS "
# Gate on GET /graph/v1.0/users/admin -> 200 (like run-github.py), not .well-known:
# this needs the full proxy -> settings -> graph -> LDAP path incl. admin role assignment.
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
# smoke.sh is topology-agnostic (graph API via OCIS_URL only); reused from ldap-smoke.
../ldap-smoke/smoke.sh
