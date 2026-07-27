#!/bin/sh
# Seeds admin's personal space and a "Demo Space" project space with demo
# content, and adds the standard demo users as editors of the Demo Space.
# No-ops immediately unless DEMO_USERS is "true".

echo "[demo-content-init] starting"

if [ "${DEMO_USERS:-false}" != "true" ]; then
  echo "[demo-content-init] DEMO_USERS is not 'true', nothing to do."
  exit 0
fi

apk add --no-cache curl jq >/dev/null 2>&1

BASE="http://ocis:9200"
AUTH="admin:${ADMIN_PASSWORD:-admin}"
CONTENT_DIR="/demo-content"
EDITOR_ROLE="58c63c02-1d89-4572-916a-870abc5a1b7d"
DEMO_USERNAMES="einstein marie moss richard katherine"

curl_g() {
  curl -s -u "$AUTH" "$@"
}

echo "[demo-content-init] waiting for oCIS to become ready..."
i=0
while true; do
  code=$(curl -s -o /dev/null -w '%{http_code}' -u "$AUTH" "$BASE/graph/v1.0/me")
  [ "$code" = "200" ] && break
  i=$((i + 1))
  if [ "$i" -ge 60 ]; then
    echo "[demo-content-init] oCIS did not become ready in time (last status: $code), giving up."
    exit 1
  fi
  sleep 5
done
echo "[demo-content-init] oCIS is ready."

upload_tree() {
  local_dir="$1"
  space_id="$2"
  base_url="$BASE/dav/spaces/$space_id"

  find "$local_dir" -type d | sed "s|^$local_dir||" | sed '/^$/d' | sort | while read -r d; do
    code=$(curl -s -o /dev/null -w '%{http_code}' -u "$AUTH" -X MKCOL "$base_url$d")
    echo "  MKCOL $code $d"
  done

  find "$local_dir" -type f | while read -r f; do
    rel="${f#$local_dir}"
    code=$(curl -s -o /dev/null -w '%{http_code}' -u "$AUTH" -X PUT "$base_url$rel" --data-binary "@$f")
    echo "  PUT $code $rel"
  done
}

echo "[demo-content-init] resolving admin's personal space..."
PERSONAL_ID=$(curl_g "$BASE/graph/v1.0/me/drives" | jq -r '.value[] | select(.driveType=="personal") | .id')
echo "[demo-content-init] personal space id: $PERSONAL_ID"

echo "[demo-content-init] uploading personal-space demo content..."
upload_tree "$CONTENT_DIR/personal" "$PERSONAL_ID"

echo "[demo-content-init] resolving Demo Space..."
DEMO_SPACE_ID=$(curl_g "$BASE/graph/v1.0/me/drives" | jq -r '.value[] | select(.driveType=="project" and .name=="Demo Space") | .id' | head -1)

if [ -z "$DEMO_SPACE_ID" ]; then
  echo "[demo-content-init] Demo Space not found, creating it..."
  DEMO_SPACE_ID=$(curl -s -u "$AUTH" -X POST "$BASE/graph/v1beta1/drives?template=default" \
    -H 'Content-Type: application/json' --data '{"name":"Demo Space","driveType":"project"}' | jq -r '.id')
fi
echo "[demo-content-init] Demo Space id: $DEMO_SPACE_ID"

echo "[demo-content-init] uploading Demo Space content..."
upload_tree "$CONTENT_DIR/project" "$DEMO_SPACE_ID"

if [ -f "$CONTENT_DIR/space-readme.md" ]; then
  code=$(curl -s -o /dev/null -w '%{http_code}' -u "$AUTH" -X PUT "$BASE/dav/spaces/$DEMO_SPACE_ID/.space/readme.md" --data-binary "@$CONTENT_DIR/space-readme.md")
  echo "  PUT $code .space/readme.md"
fi

echo "[demo-content-init] adding demo users as Demo Space editors..."
EXISTING_PERMS=$(curl_g "$BASE/graph/v1beta1/drives/$DEMO_SPACE_ID/root/permissions")

for uname in $DEMO_USERNAMES; do
  USER_ID=""
  i=0
  while [ -z "$USER_ID" ]; do
    USER_ID=$(curl_g "$BASE/graph/v1.0/users?\$search=$uname" | jq -r --arg u "$uname" '.value[] | select(.onPremisesSamAccountName==$u) | .id' | head -1)
    if [ -z "$USER_ID" ]; then
      i=$((i + 1))
      if [ "$i" -ge 30 ]; then
        echo "[demo-content-init] user $uname never showed up after waiting, skipping."
        break
      fi
      echo "[demo-content-init] waiting for user $uname to exist..."
      sleep 5
    fi
  done
  [ -z "$USER_ID" ] && continue

  ALREADY=$(echo "$EXISTING_PERMS" | jq -r --arg id "$USER_ID" '.value[]? | select(.grantedToV2.user.id==$id) | .id')
  if [ -n "$ALREADY" ]; then
    echo "[demo-content-init] $uname is already a member."
    continue
  fi

  code=$(curl -s -o /dev/null -w '%{http_code}' -u "$AUTH" -X POST "$BASE/graph/v1beta1/drives/$DEMO_SPACE_ID/root/invite" \
    -H 'Content-Type: application/json' \
    --data "{\"recipients\":[{\"objectId\":\"$USER_ID\",\"@libre.graph.recipient.type\":\"user\"}],\"roles\":[\"$EDITOR_ROLE\"]}")
  echo "  invite $code $uname"
done

echo "[demo-content-init] done."
