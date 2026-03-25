#!/usr/bin/env bash
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
OCIS_BIN="$REPO_ROOT/ocis/bin/ocis"
WRAPPER_BIN="$REPO_ROOT/tests/ociswrapper/bin/ociswrapper"
OCIS_URL="https://localhost:9200"
OCIS_CONFIG_DIR="$HOME/.ocis/config"

# suite(s) to run — set via env or passed from CI matrix
: "${BEHAT_SUITES:?BEHAT_SUITES is required, e.g. BEHAT_SUITES=apiGraph bash run-graph.sh}"

# build
make -C "$REPO_ROOT/ocis" build
GOWORK=off make -C "$REPO_ROOT/tests/ociswrapper" build

# php deps
cd "$REPO_ROOT"
composer install --no-progress
composer bin behat install --no-progress

# init ocis config
"$OCIS_BIN" init --insecure true
cp "$REPO_ROOT/tests/config/drone/app-registry.yaml" "$OCIS_CONFIG_DIR/app-registry.yaml"

# start ociswrapper in background, kill on exit
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

# wait for ocis graph API to be ready
echo "Waiting for ocis..."
timeout 300 bash -c \
  "while [ \$(curl -sk -uadmin:admin $OCIS_URL/graph/v1.0/users/admin \
    -w %{http_code} -o /dev/null) != 200 ]; do sleep 1; done"
echo "ocis ready."

# run acceptance tests for declared suites
# ACCEPTANCE_TEST_TYPE: "api" (default) or "core-api"
ACCEPTANCE_TEST_TYPE="${ACCEPTANCE_TEST_TYPE:-api}"

if [ "$ACCEPTANCE_TEST_TYPE" = "core-api" ]; then
  _FILTER_TAGS="~@skipOnGraph&&~@skipOnOcis-OCIS-Storage"
  _EXPECTED_FAILURES="${EXPECTED_FAILURES_FILE:-$REPO_ROOT/tests/acceptance/expected-failures-API-on-OCIS-storage.md}"
else
  _FILTER_TAGS="~@skip&&~@skipOnGraph&&~@skipOnOcis-OCIS-Storage"
  _EXPECTED_FAILURES="${EXPECTED_FAILURES_FILE:-$REPO_ROOT/tests/acceptance/expected-failures-localAPI-on-OCIS-storage.md}"
fi

echo "Running suites: $BEHAT_SUITES (type: $ACCEPTANCE_TEST_TYPE)"
TEST_SERVER_URL=$OCIS_URL \
OCIS_WRAPPER_URL=http://localhost:5200 \
BEHAT_SUITES=$BEHAT_SUITES \
ACCEPTANCE_TEST_TYPE=$ACCEPTANCE_TEST_TYPE \
BEHAT_FILTER_TAGS="$_FILTER_TAGS" \
EXPECTED_FAILURES_FILE="$_EXPECTED_FAILURES" \
STORAGE_DRIVER=ocis \
  make -C "$REPO_ROOT" test-acceptance-api
