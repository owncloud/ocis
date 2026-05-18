#!/usr/bin/env bash
set -euo pipefail

# Start oCIS via ociswrapper.
# ociswrapper wraps the ocis binary and exposes HTTP API on :5200
# for dynamic reconfiguration (env var changes at runtime).
#
# Environment variables set by action.yml:
#   ADMIN_PASSWORD, LOG_LEVEL, DEMO_USERS
#   ANTIVIRUS_ENABLED, EMAIL_ENABLED, TIKA_ENABLED
#   EXTRA_SERVER_ENV (JSON object)
#   OCIS_REPO_ROOT (GITHUB_WORKSPACE — needed for config file paths)

OCIS_URL="https://localhost:9200"
CONFIG_DIR="${HOME}/.ocis/config"

# Generate a throw-away fontsMap.json pointing at the font shipped in the repo.
# Tests that render thumbnails from text files need this.
REPO_ROOT="${OCIS_REPO_ROOT:-${GITHUB_WORKSPACE}}"
FONT_PATH="${REPO_ROOT}/tests/config/ci/NotoSans.ttf"
FONTMAP=$(mktemp /tmp/fontsMap-XXXXXX.json)
echo "{\"defaultFont\": \"${FONT_PATH}\"}" > "$FONTMAP"

declare -A SERVER_ENV=(
  [OCIS_URL]="$OCIS_URL"
  [OCIS_CONFIG_DIR]="$CONFIG_DIR"
  [STORAGE_USERS_DRIVER]="ocis"
  [PROXY_ENABLE_BASIC_AUTH]="true"
  [OCIS_LOG_LEVEL]="${LOG_LEVEL:-error}"
  [IDM_CREATE_DEMO_USERS]="${DEMO_USERS:-false}"
  [IDM_ADMIN_PASSWORD]="${ADMIN_PASSWORD:-admin}"
  [FRONTEND_SEARCH_MIN_LENGTH]="2"
  [OCIS_ASYNC_UPLOADS]="true"
  [OCIS_EVENTS_ENABLE_TLS]="false"
  [NATS_NATS_HOST]="0.0.0.0"
  [NATS_NATS_PORT]="9233"
  [MICRO_REGISTRY_ADDRESS]="127.0.0.1:9233"
  [OCIS_JWT_SECRET]="some-ocis-jwt-secret"
  [EVENTHISTORY_STORE]="memory"
  [OCIS_TRANSLATION_PATH]="${REPO_ROOT}/tests/config/translations"
  [WEB_UI_CONFIG_FILE]="${REPO_ROOT}/tests/config/ci/ocis-config.json"
  [THUMBNAILS_TXT_FONTMAP_FILE]="$FONTMAP"
  [SEARCH_EXTRACTOR_TYPE]="basic"
  [FRONTEND_FULL_TEXT_SEARCH_ENABLED]="false"
  # debug addresses
  [ACTIVITYLOG_DEBUG_ADDR]="0.0.0.0:9197"
  [APP_PROVIDER_DEBUG_ADDR]="0.0.0.0:9165"
  [APP_REGISTRY_DEBUG_ADDR]="0.0.0.0:9243"
  [AUTH_BASIC_DEBUG_ADDR]="0.0.0.0:9147"
  [AUTH_MACHINE_DEBUG_ADDR]="0.0.0.0:9167"
  [AUTH_SERVICE_DEBUG_ADDR]="0.0.0.0:9198"
  [CLIENTLOG_DEBUG_ADDR]="0.0.0.0:9260"
  [EVENTHISTORY_DEBUG_ADDR]="0.0.0.0:9270"
  [FRONTEND_DEBUG_ADDR]="0.0.0.0:9141"
  [GATEWAY_DEBUG_ADDR]="0.0.0.0:9143"
  [GRAPH_DEBUG_ADDR]="0.0.0.0:9124"
  [GROUPS_DEBUG_ADDR]="0.0.0.0:9161"
  [IDM_DEBUG_ADDR]="0.0.0.0:9239"
  [IDP_DEBUG_ADDR]="0.0.0.0:9134"
  [INVITATIONS_DEBUG_ADDR]="0.0.0.0:9269"
  [NATS_DEBUG_ADDR]="0.0.0.0:9234"
  [OCDAV_DEBUG_ADDR]="0.0.0.0:9163"
  [OCM_DEBUG_ADDR]="0.0.0.0:9281"
  [OCS_DEBUG_ADDR]="0.0.0.0:9114"
  [POSTPROCESSING_DEBUG_ADDR]="0.0.0.0:9255"
  [PROXY_DEBUG_ADDR]="0.0.0.0:9205"
  [SEARCH_DEBUG_ADDR]="0.0.0.0:9224"
  [SETTINGS_DEBUG_ADDR]="0.0.0.0:9194"
  [SHARING_DEBUG_ADDR]="0.0.0.0:9151"
  [SSE_DEBUG_ADDR]="0.0.0.0:9139"
  [STORAGE_PUBLICLINK_DEBUG_ADDR]="0.0.0.0:9179"
  [STORAGE_SHARES_DEBUG_ADDR]="0.0.0.0:9156"
  [STORAGE_SYSTEM_DEBUG_ADDR]="0.0.0.0:9217"
  [STORAGE_USERS_DEBUG_ADDR]="0.0.0.0:9159"
  [THUMBNAILS_DEBUG_ADDR]="0.0.0.0:9189"
  [USERLOG_DEBUG_ADDR]="0.0.0.0:9214"
  [USERS_DEBUG_ADDR]="0.0.0.0:9145"
  [WEB_DEBUG_ADDR]="0.0.0.0:9104"
  [WEBDAV_DEBUG_ADDR]="0.0.0.0:9119"
  [WEBFINGER_DEBUG_ADDR]="0.0.0.0:9279"
)

# Antivirus
if [[ "${ANTIVIRUS_ENABLED:-false}" == "true" ]]; then
  SERVER_ENV[ANTIVIRUS_SCANNER_TYPE]="clamav"
  SERVER_ENV[ANTIVIRUS_CLAMAV_SOCKET]="tcp://localhost:3310"
  SERVER_ENV[POSTPROCESSING_STEPS]="virusscan"
  SERVER_ENV[OCIS_ADD_RUN_SERVICES]="antivirus"
  SERVER_ENV[ANTIVIRUS_DEBUG_ADDR]="0.0.0.0:9277"
fi

# Email (notifications service)
if [[ "${EMAIL_ENABLED:-false}" == "true" ]]; then
  SERVER_ENV[OCIS_ADD_RUN_SERVICES]="${SERVER_ENV[OCIS_ADD_RUN_SERVICES]:+${SERVER_ENV[OCIS_ADD_RUN_SERVICES]},}notifications"
  SERVER_ENV[NOTIFICATIONS_SMTP_HOST]="localhost"
  SERVER_ENV[NOTIFICATIONS_SMTP_PORT]="1025"
  SERVER_ENV[NOTIFICATIONS_SMTP_INSECURE]="true"
  SERVER_ENV[NOTIFICATIONS_SMTP_SENDER]="ownCloud <noreply@example.com>"
  SERVER_ENV[NOTIFICATIONS_DEBUG_ADDR]="0.0.0.0:9174"
fi

# Tika (full-text search)
if [[ "${TIKA_ENABLED:-false}" == "true" ]]; then
  SERVER_ENV[FRONTEND_FULL_TEXT_SEARCH_ENABLED]="true"
  SERVER_ENV[SEARCH_EXTRACTOR_TYPE]="tika"
  SERVER_ENV[SEARCH_EXTRACTOR_TIKA_TIKA_URL]="http://localhost:9998"
  SERVER_ENV[SEARCH_EXTRACTOR_CS3SOURCE_INSECURE]="true"
fi

# Extra env vars from JSON input — use null-delimited records to handle values with '=' or newlines
if [[ -n "${EXTRA_SERVER_ENV:-}" && "${EXTRA_SERVER_ENV}" != "{}" ]]; then
  while IFS=$'\x01' read -r -d $'\x00' key val; do
    SERVER_ENV["$key"]="$val"
  done < <(echo "$EXTRA_SERVER_ENV" | python3 -c "
import sys, json
d = json.load(sys.stdin)
for k, v in d.items():
    sys.stdout.buffer.write(k.encode() + b'\x01' + v.encode() + b'\x00')
")
fi

# Build env for the subprocess
ENV_ARGS=()
for key in "${!SERVER_ENV[@]}"; do
  ENV_ARGS+=("${key}=${SERVER_ENV[$key]}")
done

echo "Starting ociswrapper + oCIS server..."
env "${ENV_ARGS[@]}" ociswrapper serve \
  --bin /usr/local/bin/ocis \
  --url "$OCIS_URL" \
  --admin-username admin \
  --admin-password "${ADMIN_PASSWORD:-admin}" \
  > /tmp/ocis-server.log 2>&1 &

echo $! > /tmp/ocis-wrapper.pid
echo "ociswrapper started (PID $(cat /tmp/ocis-wrapper.pid)), log: /tmp/ocis-server.log"
