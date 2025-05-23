#!/usr/bin/env bash

set -e

COMMIT_SHA_SHORT=${DRONE_COMMIT:0:8}
BUILD_STATUS="✅ Success"
ROOMID="!rnWsCVUmDHDJbiSPMM:matrix.org"


if [ ${DRONE_BUILD_STATUS} == "failure" ]; then
  BUILD_STATUS="❌️ Failure"
fi

message_html='<b>'$BUILD_STATUS'</b> <a href="'${DRONE_BUILD_LINK}'">'${DRONE_REPO}'#'$COMMIT_SHA_SHORT'</a> ('${DRONE_BRANCH}') by <b>'${DRONE_COMMIT_AUTHOR}'</b>'
message_html=$(echo "$message_html" | sed 's/\\/\\\\/g' | sed 's/"/\\"/g')

response=$(curl -s -o /dev/null -X PUT -w "%{http_code}" 'https://matrix.org/_matrix/client/v3/rooms/'$ROOMID'/send/m.room.message/'$(date +%s) \
  -H "Authorization: Bearer "$MATRIX_TOKEN \
  -H 'Content-Type: application/json' \
  -d '{
    "msgtype": "m.text",
    "body": "'"$message_html"'",
    "format": "org.matrix.custom.html",
    "formatted_body": "'"$message_html"'"
  }')

if [ $status_code != 200 ]; then
  echo "❌ Error: failed sending notification to matrix"
  exit 1
fi
