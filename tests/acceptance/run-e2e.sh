#!/usr/bin/env bash
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
OCIS_BIN="$REPO_ROOT/ocis/bin/ocis"
WRAPPER_BIN="$REPO_ROOT/tests/ociswrapper/bin/ociswrapper"
OCIS_URL="https://localhost:9200"
OCIS_CONFIG_DIR="$HOME/.ocis/config"
WEB_DIR="$REPO_ROOT/webTestRunner"

: "${E2E_ARGS:?E2E_ARGS is required, e.g. E2E_ARGS='--run-part 1' bash run-e2e.sh}"

# build ocis + ociswrapper
make -C "$REPO_ROOT/ocis" build
GOWORK=off make -C "$REPO_ROOT/tests/ociswrapper" build

# clone owncloud/web (test runner lives there)
if [ ! -d "$WEB_DIR" ]; then
  git clone --depth 1 https://github.com/owncloud/web.git "$WEB_DIR"
fi
cd "$WEB_DIR"
npm install -g pnpm
pnpm install

# init + start ocis
"$OCIS_BIN" init --insecure true
cp "$REPO_ROOT/tests/config/drone/app-registry.yaml" "$OCIS_CONFIG_DIR/app-registry.yaml"

OCIS_URL=$OCIS_URL \
OCIS_CONFIG_DIR=$OCIS_CONFIG_DIR \
STORAGE_USERS_DRIVER=ocis \
PROXY_ENABLE_BASIC_AUTH=true \
OCIS_EXCLUDE_RUN_SERVICES=idp \
OCIS_LOG_LEVEL=error \
IDM_CREATE_DEMO_USERS=true \
IDM_ADMIN_PASSWORD=admin \
OCIS_ASYNC_UPLOADS=true \
OCIS_EVENTS_ENABLE_TLS=false \
NATS_NATS_HOST=0.0.0.0 \
NATS_NATS_PORT=9233 \
OCIS_JWT_SECRET=some-ocis-jwt-secret \
WEB_UI_CONFIG_FILE="$REPO_ROOT/tests/config/drone/ocis-config.json" \
  "$WRAPPER_BIN" serve \
    --bin "$OCIS_BIN" \
    --url "$OCIS_URL" \
    --admin-username admin \
    --admin-password admin &
WRAPPER_PID=$!
trap "kill $WRAPPER_PID 2>/dev/null || true" EXIT

echo "Waiting for ocis..."
timeout 300 bash -c \
  "while [ \$(curl -sk -uadmin:admin $OCIS_URL/graph/v1.0/users/admin \
    -w %{http_code} -o /dev/null) != 200 ]; do sleep 1; done"
echo "ocis ready."

# run playwright e2e tests
cd "$WEB_DIR/tests/e2e"
echo "Running e2e: $E2E_ARGS"
BASE_URL_OCIS=$OCIS_URL \
HEADLESS=true \
RETRY=1 \
SKIP_A11Y_TESTS=true \
REPORT_TRACING=true \
  bash run-e2e.sh $E2E_ARGS
