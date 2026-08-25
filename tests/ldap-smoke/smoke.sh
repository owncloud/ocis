#!/usr/bin/env bash
# CRUD + login + group ops e2e against a single external LDAP, via the oCIS graph API
# with basic auth. Green = create + server-UUID read-back + login + group Modify/refint
# + member remove + delete all work with the OCISDEV-1030/1032/1033 LDAP client changes
# (configurable read/write-aware retry + read-after-write) in the binary.
set -euo pipefail
cd "$(dirname "$0")"

BASE="${OCIS_URL:-https://localhost:9200}"

pass() { echo "[PASS] $1"; }
fail() { echo "[FAIL] $1"; [ -n "${2:-}" ] && echo "$2"; exit 1; }

# reqas USER:PASS METHOD PATH [extra curl args...] -> sets $code and $json.
reqas() {
  local cred=$1 method=$2 path=$3; shift 3
  local out
  out=$(curl -sk -u "$cred" -H "Accept: application/json" -X "$method" "$@" \
    -w $'\n%{http_code}' "$BASE$path")
  code=$(printf '%s' "$out" | tail -n1)
  json=$(printf '%s' "$out" | sed '$d')
}

# req METHOD PATH [extra curl args...] -> as reqas, authed as admin.
req() { reqas "admin:admin" "$@"; }

# 1. list users -> 200 + seeded demo users
echo "[*] GET /graph/v1.0/users"
req GET /graph/v1.0/users
[ "$code" = "200" ] || fail "list users returned $code" "$json"
count=$(printf '%s' "$json" | jq '.value | length')
printf '%s' "$json" | jq -e '.value[] | select(.onPremisesSamAccountName=="einstein")' >/dev/null \
  || fail "seeded user einstein not listed" "$json"
pass "list users -> 200 ($count users, einstein present)"

# 2. create a user -> 201 with a server-assigned id (the read-back result)
NEWUSER="ldapsmoke"
echo "[*] POST /graph/v1.0/users ($NEWUSER)"
payload=$(jq -n --arg u "$NEWUSER" '{
  displayName: "LDAP Smoke",
  onPremisesSamAccountName: $u,
  mail: ($u + "@example.org"),
  accountEnabled: true,
  passwordProfile: { password: "Secret123!" }
}')
req POST /graph/v1.0/users -H "Content-Type: application/json" -d "$payload"
[ "$code" = "201" ] || fail "create user returned $code (expected 201)" "$json"
ID=$(printf '%s' "$json" | jq -r '.id // empty')
[ -n "$ID" ] || fail "create returned no id" "$json"
pass "create user -> 201, server id=$ID"

# 3. read it back -> 200
echo "[*] GET /graph/v1.0/users/$ID"
req GET "/graph/v1.0/users/$ID"
[ "$code" = "200" ] || fail "read-back returned $code" "$json"
got=$(printf '%s' "$json" | jq -r '.onPremisesSamAccountName')
[ "$got" = "$NEWUSER" ] || fail "read-back sam mismatch: $got" "$json"
pass "read back user -> 200 ($got)"

# 4. log in as the created user -> 200 (proves the LDAP account + password work)
echo "[*] GET /graph/v1.0/me (as $NEWUSER)"
reqas "$NEWUSER:Secret123!" GET /graph/v1.0/me
[ "$code" = "200" ] || fail "login as $NEWUSER returned $code" "$json"
me=$(printf '%s' "$json" | jq -r '.onPremisesSamAccountName')
[ "$me" = "$NEWUSER" ] || fail "login identity mismatch: $me" "$json"
pass "login as $NEWUSER -> 200 (/me = $me)"

# 5. create a group -> 201 with a server-assigned id
GRP="ldapsmoke-grp"
echo "[*] POST /graph/v1.0/groups ($GRP)"
req POST /graph/v1.0/groups -H "Content-Type: application/json" \
  -d "$(jq -n --arg n "$GRP" '{displayName: $n}')"
[ "$code" = "201" ] || fail "create group returned $code (expected 201)" "$json"
GID=$(printf '%s' "$json" | jq -r '.id // empty')
[ -n "$GID" ] || fail "create group returned no id" "$json"
pass "create group -> 201, server id=$GID"

# 6. add the user as a member -> 204 (LDAP Modify + refint overlay)
echo "[*] POST /graph/v1.0/groups/$GID/members/\$ref ($NEWUSER)"
req POST "/graph/v1.0/groups/$GID/members/\$ref" -H "Content-Type: application/json" \
  -d "$(jq -n --arg u "$BASE/graph/v1.0/users/$ID" '{"@odata.id": $u}')"
[ "$code" = "204" ] || fail "add member returned $code (expected 204)" "$json"
pass "add member -> 204"

# 7. list members -> 200, the user is present
echo "[*] GET /graph/v1.0/groups/$GID/members"
req GET "/graph/v1.0/groups/$GID/members"
[ "$code" = "200" ] || fail "list members returned $code" "$json"
printf '%s' "$json" | jq -e --arg id "$ID" '.[] | select(.id==$id)' >/dev/null \
  || fail "added member $ID not listed" "$json"
pass "list members -> 200 ($NEWUSER present)"

# 8. remove the member -> 204, then re-list -> gone
echo "[*] DELETE /graph/v1.0/groups/$GID/members/$ID/\$ref"
req DELETE "/graph/v1.0/groups/$GID/members/$ID/\$ref"
[ "$code" = "204" ] || fail "remove member returned $code (expected 204)" "$json"
req GET "/graph/v1.0/groups/$GID/members"
printf '%s' "$json" | jq -e --arg id "$ID" '.[] | select(.id==$id)' >/dev/null \
  && fail "removed member $ID still listed" "$json"
pass "remove member -> 204, member gone"

# 9. delete the group -> 204
echo "[*] DELETE /graph/v1.0/groups/$GID"
req DELETE "/graph/v1.0/groups/$GID"
[ "$code" = "204" ] || fail "delete group returned $code (expected 204)" "$json"
pass "delete group -> 204"

# 10. delete the user -> 204, then re-GET -> 404
echo "[*] DELETE /graph/v1.0/users/$ID"
req DELETE "/graph/v1.0/users/$ID"
[ "$code" = "204" ] || fail "delete returned $code (expected 204)"
req GET "/graph/v1.0/users/$ID"
[ "$code" = "404" ] || fail "read-after-delete returned $code (expected 404)"
pass "delete user -> 204, read-after-delete -> 404"

echo
echo "[✓] LDAP smoke green: CRUD + login + group add/list/remove + delete works against a real external LDAP"
